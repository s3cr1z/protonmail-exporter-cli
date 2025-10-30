package mail

import (
	"context"
	"net/mail"
	"testing"

	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestMetadataStage_WithFilters tests filtering integration with MetadataStage
func TestMetadataStage_WithFilters(t *testing.T) {
	// Helper to setup common mock expectations
	setupMockExpectations := func(client *apiclient.MockClient, fileChecker *MockMetadataFileChecker, messages []proton.MessageMetadata) {
		// First call returns all messages, second call returns empty (end of pagination)
		client.EXPECT().GetMessageMetadataPage(gomock.Any(), gomock.Eq(0), gomock.Eq(len(messages)), gomock.Any()).
			Return(messages, nil)
		client.EXPECT().GetMessageMetadataPage(gomock.Any(), gomock.Eq(0), gomock.Eq(len(messages)), gomock.Any()).
			Return([]proton.MessageMetadata{}, nil)

		fileChecker.EXPECT().HasMessage(gomock.Any()).AnyTimes().Return(false, nil)
	}

	t.Run("label filter", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		client := apiclient.NewMockClient(mockCtrl)
		errReporter := NewMockStageErrorReporter(mockCtrl)
		fileChecker := NewMockMetadataFileChecker(mockCtrl)
		reporter := NewMockReporter(mockCtrl)

		// Create test metadata with different labels
		allMessages := []proton.MessageMetadata{
			{ID: "msg1", LabelIDs: []string{"0"}},      // Inbox
			{ID: "msg2", LabelIDs: []string{"2"}},      // Sent
			{ID: "msg3", LabelIDs: []string{"0", "5"}}, // Inbox + All Mail
			{ID: "msg4", LabelIDs: []string{"3"}},      // Trash
		}

		setupMockExpectations(client, fileChecker, allMessages)
		
		// We expect 2 messages to be filtered out (msg2 and msg4)
		const expectedFilteredCount = 2
		reporter.EXPECT().OnProgress(expectedFilteredCount)

		// Filter for Inbox only
		filters := FilterOptions{
			LabelIDs: []string{"0"},
		}
		metadata := NewMetadataStage(client, logrus.WithField("test", "test"), len(allMessages), 10, filters)

		go func() {
			metadata.Run(context.Background(), errReporter, fileChecker, reporter)
		}()

		result := make([]proton.MessageMetadata, 0)
		for out := range metadata.outputCh {
			result = append(result, out...)
		}

		// Should only get msg1 and msg3 (both have label "0")
		require.Len(t, result, 2)
		require.Equal(t, "msg1", result[0].ID)
		require.Equal(t, "msg3", result[1].ID)
	})

	t.Run("date range filter", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		client := apiclient.NewMockClient(mockCtrl)
		errReporter := NewMockStageErrorReporter(mockCtrl)
		fileChecker := NewMockMetadataFileChecker(mockCtrl)
		reporter := NewMockReporter(mockCtrl)

		// Create test metadata with different timestamps
		allMessages := []proton.MessageMetadata{
			{ID: "msg1", Time: 1000},
			{ID: "msg2", Time: 1500},
			{ID: "msg3", Time: 2000},
			{ID: "msg4", Time: 2500},
		}

		setupMockExpectations(client, fileChecker, allMessages)
		
		// We expect 2 messages to be filtered out (msg1 at 1000 and msg4 at 2500)
		const expectedFilteredCount = 2
		reporter.EXPECT().OnProgress(expectedFilteredCount)

		// Filter for messages between 1200 and 2200
		filters := FilterOptions{
			After:  1200,
			Before: 2200,
		}
		metadata := NewMetadataStage(client, logrus.WithField("test", "test"), len(allMessages), 10, filters)

		go func() {
			metadata.Run(context.Background(), errReporter, fileChecker, reporter)
		}()

		result := make([]proton.MessageMetadata, 0)
		for out := range metadata.outputCh {
			result = append(result, out...)
		}

		// Should get msg2 and msg3 (timestamps 1500 and 2000)
		require.Len(t, result, 2)
		require.Equal(t, "msg2", result[0].ID)
		require.Equal(t, "msg3", result[1].ID)
	})

	t.Run("sender filter", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		client := apiclient.NewMockClient(mockCtrl)
		errReporter := NewMockStageErrorReporter(mockCtrl)
		fileChecker := NewMockMetadataFileChecker(mockCtrl)
		reporter := NewMockReporter(mockCtrl)

		// Create test metadata with different senders
		allMessages := []proton.MessageMetadata{
			{ID: "msg1", Sender: &mail.Address{Address: "alice@example.com"}},
			{ID: "msg2", Sender: &mail.Address{Address: "bob@example.com"}},
			{ID: "msg3", Sender: &mail.Address{Address: "charlie@example.com"}},
		}

		setupMockExpectations(client, fileChecker, allMessages)
		
		// We expect 2 messages to be filtered out (bob and charlie)
		const expectedFilteredCount = 2
		reporter.EXPECT().OnProgress(expectedFilteredCount)

		// Filter for alice only
		filters := FilterOptions{
			From: []string{"alice@example.com"},
		}
		metadata := NewMetadataStage(client, logrus.WithField("test", "test"), len(allMessages), 10, filters)

		go func() {
			metadata.Run(context.Background(), errReporter, fileChecker, reporter)
		}()

		result := make([]proton.MessageMetadata, 0)
		for out := range metadata.outputCh {
			result = append(result, out...)
		}

		// Should only get msg1
		require.Len(t, result, 1)
		require.Equal(t, "msg1", result[0].ID)
	})

	t.Run("domain filter", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		client := apiclient.NewMockClient(mockCtrl)
		errReporter := NewMockStageErrorReporter(mockCtrl)
		fileChecker := NewMockMetadataFileChecker(mockCtrl)
		reporter := NewMockReporter(mockCtrl)

		// Create test metadata with different domains
		allMessages := []proton.MessageMetadata{
			{ID: "msg1", Sender: &mail.Address{Address: "user@example.com"}},
			{ID: "msg2", Sender: &mail.Address{Address: "user@proton.me"}},
			{ID: "msg3", Sender: &mail.Address{Address: "admin@example.com"}},
		}

		setupMockExpectations(client, fileChecker, allMessages)
		
		// We expect 1 message to be filtered out (proton.me)
		const expectedFilteredCount = 1
		reporter.EXPECT().OnProgress(expectedFilteredCount)

		// Filter for example.com domain
		filters := FilterOptions{
			FromDomains: []string{"example.com"},
		}
		metadata := NewMetadataStage(client, logrus.WithField("test", "test"), len(allMessages), 10, filters)

		go func() {
			metadata.Run(context.Background(), errReporter, fileChecker, reporter)
		}()

		result := make([]proton.MessageMetadata, 0)
		for out := range metadata.outputCh {
			result = append(result, out...)
		}

		// Should get msg1 and msg3
		require.Len(t, result, 2)
		require.Equal(t, "msg1", result[0].ID)
		require.Equal(t, "msg3", result[1].ID)
	})

	t.Run("combined filters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		client := apiclient.NewMockClient(mockCtrl)
		errReporter := NewMockStageErrorReporter(mockCtrl)
		fileChecker := NewMockMetadataFileChecker(mockCtrl)
		reporter := NewMockReporter(mockCtrl)

		// Create test metadata with multiple attributes
		allMessages := []proton.MessageMetadata{
			{
				ID:       "msg1",
				LabelIDs: []string{"0"},
				Time:     1500,
				Sender:   &mail.Address{Address: "alice@example.com"},
			},
			{
				ID:       "msg2",
				LabelIDs: []string{"0"},
				Time:     2500,
				Sender:   &mail.Address{Address: "alice@example.com"},
			},
			{
				ID:       "msg3",
				LabelIDs: []string{"0"},
				Time:     1500,
				Sender:   &mail.Address{Address: "bob@example.com"},
			},
			{
				ID:       "msg4",
				LabelIDs: []string{"2"},
				Time:     1500,
				Sender:   &mail.Address{Address: "alice@example.com"},
			},
		}

		setupMockExpectations(client, fileChecker, allMessages)
		
		// We expect 3 messages to be filtered out (msg2 fails time, msg3 fails sender, msg4 fails label)
		const expectedFilteredCount = 3
		reporter.EXPECT().OnProgress(expectedFilteredCount)

		// Filter: Inbox AND before 2000 AND from alice
		filters := FilterOptions{
			LabelIDs: []string{"0"},
			Before:   2000,
			From:     []string{"alice@example.com"},
		}
		metadata := NewMetadataStage(client, logrus.WithField("test", "test"), len(allMessages), 10, filters)

		go func() {
			metadata.Run(context.Background(), errReporter, fileChecker, reporter)
		}()

		result := make([]proton.MessageMetadata, 0)
		for out := range metadata.outputCh {
			result = append(result, out...)
		}

		// Should only get msg1 (matches all criteria)
		require.Len(t, result, 1)
		require.Equal(t, "msg1", result[0].ID)
	})

	t.Run("no filters - all messages pass", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		client := apiclient.NewMockClient(mockCtrl)
		errReporter := NewMockStageErrorReporter(mockCtrl)
		fileChecker := NewMockMetadataFileChecker(mockCtrl)
		reporter := NewMockReporter(mockCtrl)

		allMessages := []proton.MessageMetadata{
			{ID: "msg1"},
			{ID: "msg2"},
			{ID: "msg3"},
		}

		setupMockExpectations(client, fileChecker, allMessages)

		// Empty filters - should export all (no OnProgress call for filtering)
		filters := FilterOptions{}
		metadata := NewMetadataStage(client, logrus.WithField("test", "test"), len(allMessages), 10, filters)

		go func() {
			metadata.Run(context.Background(), errReporter, fileChecker, reporter)
		}()

		result := make([]proton.MessageMetadata, 0)
		for out := range metadata.outputCh {
			result = append(result, out...)
		}

		// Should get all messages
		require.Len(t, result, 3)
	})
}
