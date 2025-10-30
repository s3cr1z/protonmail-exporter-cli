package mail

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterParser_ParseCommaSeparated(t *testing.T) {
	parser := FilterParser{}

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "single value",
			input:    "value1",
			expected: []string{"value1"},
		},
		{
			name:     "multiple values",
			input:    "value1,value2,value3",
			expected: []string{"value1", "value2", "value3"},
		},
		{
			name:     "values with spaces",
			input:    "value1 , value2 , value3",
			expected: []string{"value1", "value2", "value3"},
		},
		{
			name:     "values with empty entries",
			input:    "value1,,value2,",
			expected: []string{"value1", "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.ParseCommaSeparated(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterParser_ParseDate(t *testing.T) {
	parser := FilterParser{}

	tests := []struct {
		name      string
		input     string
		expectNil bool
		wantErr   bool
	}{
		{
			name:      "empty string",
			input:     "",
			expectNil: true,
			wantErr:   false,
		},
		{
			name:      "YYYY-MM-DD format",
			input:     "2024-01-15",
			expectNil: false,
			wantErr:   false,
		},
		{
			name:      "YYYY/MM/DD format",
			input:     "2024/01/15",
			expectNil: false,
			wantErr:   false,
		},
		{
			name:      "YYYYMMDD format",
			input:     "20240115",
			expectNil: false,
			wantErr:   false,
		},
		{
			name:      "invalid format",
			input:     "15-01-2024",
			expectNil: false,
			wantErr:   true,
		},
		{
			name:      "invalid date",
			input:     "not-a-date",
			expectNil: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.ParseDate(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectNil {
					assert.Nil(t, result)
				} else {
					assert.NotNil(t, result)
				}
			}
		})
	}
}

func TestParseFilterFromStrings(t *testing.T) {
	tests := []struct {
		name       string
		labelIDs   string
		sender     string
		recipient  string
		domain     string
		after      string
		before     string
		subject    string
		expectNil  bool
		wantErr    bool
		checkLabel bool
		checkDate  bool
	}{
		{
			name:      "all empty - returns nil",
			expectNil: true,
			wantErr:   false,
		},
		{
			name:       "valid label IDs",
			labelIDs:   "0,2,10",
			expectNil:  false,
			wantErr:    false,
			checkLabel: true,
		},
		{
			name:      "valid sender",
			sender:    "user@example.com",
			expectNil: false,
			wantErr:   false,
		},
		{
			name:      "valid recipient",
			recipient: "user@example.com,@domain.com",
			expectNil: false,
			wantErr:   false,
		},
		{
			name:      "valid domain",
			domain:    "example.com,another.com",
			expectNil: false,
			wantErr:   false,
		},
		{
			name:      "valid date range",
			after:     "2024-01-01",
			before:    "2024-12-31",
			expectNil: false,
			wantErr:   false,
			checkDate: true,
		},
		{
			name:      "invalid date range",
			after:     "2024-12-31",
			before:    "2024-01-01",
			expectNil: false,
			wantErr:   true,
		},
		{
			name:      "invalid sender format",
			sender:    "not-an-email",
			expectNil: false,
			wantErr:   true,
		},
		{
			name:      "invalid after date",
			after:     "invalid-date",
			expectNil: false,
			wantErr:   true,
		},
		{
			name:      "subject only",
			subject:   "important",
			expectNil: false,
			wantErr:   false,
		},
		{
			name:       "combined filters",
			labelIDs:   "0",
			sender:     "user@example.com",
			after:      "2024-01-01",
			subject:    "test",
			expectNil:  false,
			wantErr:    false,
			checkLabel: true,
			checkDate:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseFilterFromStrings(
				tt.labelIDs,
				tt.sender,
				tt.recipient,
				tt.domain,
				tt.after,
				tt.before,
				tt.subject,
			)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)

				if tt.checkLabel {
					assert.NotEmpty(t, result.LabelIDs)
				}

				if tt.checkDate {
					if tt.after != "" {
						assert.NotNil(t, result.After)
					}
					if tt.before != "" {
						assert.NotNil(t, result.Before)
					}
				}
			}
		})
	}
}

func TestParseFilterFromStrings_DateValidation(t *testing.T) {
	// Test that parsed dates are correct
	filter, err := ParseFilterFromStrings(
		"",
		"",
		"",
		"",
		"2024-01-15",
		"2024-12-20",
		"",
	)

	require.NoError(t, err)
	require.NotNil(t, filter)
	require.NotNil(t, filter.After)
	require.NotNil(t, filter.Before)

	expectedAfter := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedBefore := time.Date(2024, 12, 20, 0, 0, 0, 0, time.UTC)

	assert.True(t, filter.After.Equal(expectedAfter), "After date should be 2024-01-15")
	assert.True(t, filter.Before.Equal(expectedBefore), "Before date should be 2024-12-20")
}
