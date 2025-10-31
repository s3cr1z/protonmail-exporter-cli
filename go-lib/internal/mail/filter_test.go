package mail

import (
	"net/mail"
	"testing"
	"time"

	"github.com/ProtonMail/go-proton-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		filter   *Filter
		expected bool
	}{
		{
			name:     "empty filter",
			filter:   NewFilter(),
			expected: true,
		},
		{
			name: "filter with labels",
			filter: &Filter{
				LabelIDs: []string{"0"},
			},
			expected: false,
		},
		{
			name: "filter with sender",
			filter: &Filter{
				Sender: []string{"user@example.com"},
			},
			expected: false,
		},
		{
			name: "filter with date",
			filter: &Filter{
				After: timePtr(time.Now()),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.filter.IsEmpty())
		})
	}
}

func TestFilter_Validate(t *testing.T) {
	tests := []struct {
		name    string
		filter  *Filter
		wantErr bool
	}{
		{
			name:    "empty filter",
			filter:  NewFilter(),
			wantErr: false,
		},
		{
			name: "valid date range",
			filter: &Filter{
				After:  timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				Before: timePtr(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "invalid date range",
			filter: &Filter{
				After:  timePtr(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)),
				Before: timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			wantErr: true,
		},
		{
			name: "valid sender email",
			filter: &Filter{
				Sender: []string{"user@example.com"},
			},
			wantErr: false,
		},
		{
			name: "valid sender domain",
			filter: &Filter{
				Sender: []string{"@example.com"},
			},
			wantErr: false,
		},
		{
			name: "invalid sender",
			filter: &Filter{
				Sender: []string{"invalid"},
			},
			wantErr: true,
		},
		{
			name: "valid domain",
			filter: &Filter{
				Domain: []string{"example.com"},
			},
			wantErr: false,
		},
		{
			name: "invalid domain",
			filter: &Filter{
				Domain: []string{"@example.com"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.filter.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFilter_ToServerFilter(t *testing.T) {
	tests := []struct {
		name           string
		filter         *Filter
		expectNil      bool
		checkLabelID   bool
		expectedLabel  string
		checkSubject   bool
		expectedSubj   string
	}{
		{
			name:      "empty filter",
			filter:    NewFilter(),
			expectNil: true,
		},
		{
			name: "single label - server-side",
			filter: &Filter{
				LabelIDs: []string{"0"},
			},
			expectNil:     false,
			checkLabelID:  true,
			expectedLabel: "0",
		},
		{
			name: "multiple labels - no server filter",
			filter: &Filter{
				LabelIDs: []string{"0", "2"},
			},
			expectNil: true,
		},
		{
			name: "subject only - server-side",
			filter: &Filter{
				Subject: "test",
			},
			expectNil:    false,
			checkSubject: true,
			expectedSubj: "test",
		},
		{
			name: "sender only - no server filter",
			filter: &Filter{
				Sender: []string{"user@example.com"},
			},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.ToServerFilter()

			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				if tt.checkLabelID {
					assert.Equal(t, tt.expectedLabel, result.LabelID)
				}
				if tt.checkSubject {
					assert.Equal(t, tt.expectedSubj, result.Subject)
				}
			}
		})
	}
}

func TestFilter_NeedsClientFiltering(t *testing.T) {
	tests := []struct {
		name     string
		filter   *Filter
		expected bool
	}{
		{
			name:     "empty filter",
			filter:   NewFilter(),
			expected: false,
		},
		{
			name: "single label only - no client filtering",
			filter: &Filter{
				LabelIDs: []string{"0"},
			},
			expected: false,
		},
		{
			name: "multiple labels - needs client filtering",
			filter: &Filter{
				LabelIDs: []string{"0", "2"},
			},
			expected: true,
		},
		{
			name: "sender filter - needs client filtering",
			filter: &Filter{
				Sender: []string{"user@example.com"},
			},
			expected: true,
		},
		{
			name: "date filter - needs client filtering",
			filter: &Filter{
				After: timePtr(time.Now()),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.filter.NeedsClientFiltering())
		})
	}
}

func TestFilter_MatchesMetadata(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	tests := []struct {
		name     string
		filter   *Filter
		metadata proton.MessageMetadata
		expected bool
	}{
		{
			name:   "empty filter matches all",
			filter: NewFilter(),
			metadata: proton.MessageMetadata{
				ID: "msg1",
			},
			expected: true,
		},
		{
			name: "label filter matches",
			filter: &Filter{
				LabelIDs: []string{"0", "2"},
			},
			metadata: proton.MessageMetadata{
				ID:       "msg1",
				LabelIDs: []string{"0", "5"},
			},
			expected: true,
		},
		{
			name: "label filter no match",
			filter: &Filter{
				LabelIDs: []string{"0", "2"},
			},
			metadata: proton.MessageMetadata{
				ID:       "msg1",
				LabelIDs: []string{"3", "4"},
			},
			expected: false,
		},
		{
			name: "sender filter exact match",
			filter: &Filter{
				Sender: []string{"alice@example.com"},
			},
			metadata: proton.MessageMetadata{
				ID:     "msg1",
				Sender: &mail.Address{Address: "alice@example.com"},
			},
			expected: true,
		},
		{
			name: "sender filter domain match",
			filter: &Filter{
				Sender: []string{"@example.com"},
			},
			metadata: proton.MessageMetadata{
				ID:     "msg1",
				Sender: &mail.Address{Address: "alice@example.com"},
			},
			expected: true,
		},
		{
			name: "sender filter no match",
			filter: &Filter{
				Sender: []string{"bob@other.com"},
			},
			metadata: proton.MessageMetadata{
				ID:     "msg1",
				Sender: &mail.Address{Address: "alice@example.com"},
			},
			expected: false,
		},
		{
			name: "recipient filter matches To",
			filter: &Filter{
				Recipient: []string{"bob@example.com"},
			},
			metadata: proton.MessageMetadata{
				ID: "msg1",
				ToList: []*mail.Address{
					{Address: "bob@example.com"},
				},
			},
			expected: true,
		},
		{
			name: "recipient filter matches CC",
			filter: &Filter{
				Recipient: []string{"bob@example.com"},
			},
			metadata: proton.MessageMetadata{
				ID: "msg1",
				CCList: []*mail.Address{
					{Address: "bob@example.com"},
				},
			},
			expected: true,
		},
		{
			name: "domain filter matches sender",
			filter: &Filter{
				Domain: []string{"example.com"},
			},
			metadata: proton.MessageMetadata{
				ID:     "msg1",
				Sender: &mail.Address{Address: "alice@example.com"},
			},
			expected: true,
		},
		{
			name: "domain filter matches recipient",
			filter: &Filter{
				Domain: []string{"example.com"},
			},
			metadata: proton.MessageMetadata{
				ID:     "msg1",
				Sender: &mail.Address{Address: "alice@other.com"},
				ToList: []*mail.Address{
					{Address: "bob@example.com"},
				},
			},
			expected: true,
		},
		{
			name: "after date filter matches",
			filter: &Filter{
				After: &yesterday,
			},
			metadata: proton.MessageMetadata{
				ID:   "msg1",
				Time: now.Unix(),
			},
			expected: true,
		},
		{
			name: "after date filter no match",
			filter: &Filter{
				After: &tomorrow,
			},
			metadata: proton.MessageMetadata{
				ID:   "msg1",
				Time: now.Unix(),
			},
			expected: false,
		},
		{
			name: "before date filter matches",
			filter: &Filter{
				Before: &tomorrow,
			},
			metadata: proton.MessageMetadata{
				ID:   "msg1",
				Time: now.Unix(),
			},
			expected: true,
		},
		{
			name: "before date filter no match",
			filter: &Filter{
				Before: &yesterday,
			},
			metadata: proton.MessageMetadata{
				ID:   "msg1",
				Time: now.Unix(),
			},
			expected: false,
		},
		{
			name: "subject filter matches",
			filter: &Filter{
				Subject: "important",
			},
			metadata: proton.MessageMetadata{
				ID:      "msg1",
				Subject: "This is an IMPORTANT message",
			},
			expected: true,
		},
		{
			name: "subject filter no match",
			filter: &Filter{
				Subject: "urgent",
			},
			metadata: proton.MessageMetadata{
				ID:      "msg1",
				Subject: "This is an important message",
			},
			expected: false,
		},
		{
			name: "combined filters all match",
			filter: &Filter{
				LabelIDs: []string{"0"},
				Sender:   []string{"alice@example.com"},
				After:    &yesterday,
			},
			metadata: proton.MessageMetadata{
				ID:       "msg1",
				LabelIDs: []string{"0"},
				Sender:   &mail.Address{Address: "alice@example.com"},
				Time:     now.Unix(),
			},
			expected: true,
		},
		{
			name: "combined filters one fails",
			filter: &Filter{
				LabelIDs: []string{"0"},
				Sender:   []string{"alice@example.com"},
				After:    &tomorrow,
			},
			metadata: proton.MessageMetadata{
				ID:       "msg1",
				LabelIDs: []string{"0"},
				Sender:   &mail.Address{Address: "alice@example.com"},
				Time:     now.Unix(),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.filter.MatchesMetadata(tt.metadata))
		})
	}
}

func TestValidateEmailOrDomain(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"user@example.com", false},
		{"@example.com", false},
		{"invalid", true},
		{"@", true},
		{"", true},
		{"user@", true},
		{"@example", true}, // Invalid - domain must have a dot
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := validateEmailOrDomain(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDomain(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"example.com", false},
		{"sub.example.com", false},
		{"@example.com", true},
		{"invalid", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := validateDomain(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
