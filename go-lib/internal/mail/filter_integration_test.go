package mail

import (
	"context"
	"net/mail"
	"testing"
	"time"

	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Helper function for creating time pointers
func timePtr(t time.Time) *time.Time {
	return &t
}

// TestMetadataStage_WithFilter tests the metadata stage with various filters
func TestMetadataStage_WithFilter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := apiclient.NewMockClient(mockCtrl)
	errReporter := NewMockStageErrorReporter(mockCtrl)
	fileChecker := NewMockMetadataFileChecker(mockCtrl)
	reporter := NewMockReporter(mockCtrl)

	const pageSize = 2

	// Create test messages with different properties
	testMessages := []proton.MessageMetadata{
		{
			ID:       "msg1",
			LabelIDs: []string{"0", "5"}, // Inbox and All Mail
			Sender:   &mail.Address{Address: "alice@example.com"},
			Time:     time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC).Unix(),
			Subject:  "Important meeting",
		},
		{
			ID:       "msg2",
			LabelIDs: []string{"2", "5"}, // Sent and All Mail
			Sender:   &mail.Address{Address: "bob@work.com"},
			Time:     time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC).Unix(),
			Subject:  "Project update",
		},
	}

	tests := []struct {
		name           string
		filter         *Filter
		expectedMsgIDs []string
	}{
		{
			name:           "no filter - all messages",
			filter:         nil,
			expectedMsgIDs: []string{"msg1", "msg2"},
		},
		{
			name: "label filter - inbox only",
			filter: &Filter{
				LabelIDs: []string{"0"},
			},
			expectedMsgIDs: []string{"msg1"},
		},
		{
			name: "sender filter",
			filter: &Filter{
				Sender: []string{"alice@example.com"},
			},
			expectedMsgIDs: []string{"msg1"},
		},
		{
			name: "domain filter",
			filter: &Filter{
				Domain: []string{"work.com"},
			},
			expectedMsgIDs: []string{"msg2"},
		},
		{
			name: "date filter - after",
			filter: &Filter{
				After: timePtr(time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectedMsgIDs: []string{"msg1"},
		},
		{
			name: "subject filter",
			filter: &Filter{
				Subject: "meeting",
			},
			expectedMsgIDs: []string{"msg1"},
		},
		{
			name: "combined filters",
			filter: &Filter{
				LabelIDs: []string{"0"},
				After:    timePtr(time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectedMsgIDs: []string{"msg1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup expectations
			client.EXPECT().GetMessageMetadataPage(gomock.Any(), gomock.Eq(0), gomock.Eq(pageSize), gomock.Any()).
				Return(testMessages, nil).Times(1)

			// Empty result on next call to stop pagination
			client.EXPECT().GetMessageMetadataPage(gomock.Any(), gomock.Eq(0), gomock.Eq(pageSize), gomock.Any()).
				Return([]proton.MessageMetadata{}, nil).AnyTimes()

			fileChecker.EXPECT().HasMessage(gomock.Any()).Return(false, nil).AnyTimes()
			reporter.EXPECT().OnProgress(gomock.Any()).AnyTimes()

			// Create metadata stage with filter
			metadata := NewMetadataStage(client, logrus.WithField("test", "test"), pageSize, 1, tt.filter)

			// Run metadata stage
			go func() {
				metadata.Run(context.Background(), errReporter, fileChecker, reporter)
			}()

			// Collect results
			result := make([]proton.MessageMetadata, 0)
			for out := range metadata.outputCh {
				result = append(result, out...)
			}

			// Verify results
			require.Len(t, result, len(tt.expectedMsgIDs))
			for i, expected := range tt.expectedMsgIDs {
				assert.Equal(t, expected, result[i].ID, "Message ID mismatch at index %d", i)
			}
		})
	}
}

// TestFilter_Integration tests the complete filter integration
func TestFilter_Integration(t *testing.T) {
	// Test filter creation from strings
	filter, err := ParseFilterFromStrings(
		"0,2",                   // labels
		"user@example.com",      // sender
		"recipient@test.com",    // recipient
		"work.com",              // domain
		"2024-01-01",            // after
		"2024-12-31",            // before
		"important",             // subject
	)

	require.NoError(t, err)
	require.NotNil(t, filter)

	// Verify filter was parsed correctly
	assert.Equal(t, []string{"0", "2"}, filter.LabelIDs)
	assert.Equal(t, []string{"user@example.com"}, filter.Sender)
	assert.Equal(t, []string{"recipient@test.com"}, filter.Recipient)
	assert.Equal(t, []string{"work.com"}, filter.Domain)
	assert.Equal(t, "important", filter.Subject)
	assert.NotNil(t, filter.After)
	assert.NotNil(t, filter.Before)

	// Verify server-side filter is not created (multiple labels)
	serverFilter := filter.ToServerFilter()
	assert.Nil(t, serverFilter, "Server-side filter should be nil for multiple labels")

	// Verify client-side filtering is needed
	assert.True(t, filter.NeedsClientFiltering())

	// Test matching
	matchingMsg := proton.MessageMetadata{
		ID:       "match1",
		LabelIDs: []string{"0", "5"},
		Sender:   &mail.Address{Address: "user@example.com"},
		ToList: []*mail.Address{
			{Address: "recipient@test.com"},
		},
		Time:    time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC).Unix(),
		Subject: "This is an important message",
	}

	assert.True(t, filter.MatchesMetadata(matchingMsg))

	nonMatchingMsg := proton.MessageMetadata{
		ID:       "nomatch1",
		LabelIDs: []string{"3"}, // Wrong label
		Sender:   &mail.Address{Address: "other@example.com"},
		Time:     time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC).Unix(),
		Subject:  "Different subject",
	}

	assert.False(t, filter.MatchesMetadata(nonMatchingMsg))
}
