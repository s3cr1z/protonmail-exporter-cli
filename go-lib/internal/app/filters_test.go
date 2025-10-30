package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int64
		expectError bool
	}{
		{
			name:     "Unix timestamp",
			input:    "1609459200",
			expected: 1609459200,
		},
		{
			name:     "RFC3339 format",
			input:    "2021-01-01T00:00:00Z",
			expected: 1609459200,
		},
		{
			name:     "YYYY-MM-DD format",
			input:    "2021-01-01",
			expected: 1609459200,
		},
		{
			name:     "RFC3339 with timezone",
			input:    "2021-01-01T00:00:00+00:00",
			expected: 1609459200,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "invalid format",
			input:       "invalid-date",
			expectError: true,
		},
		{
			name:        "whitespace only",
			input:       "   ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTimestamp(tt.input)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseTimestamp_YYYYMMDD(t *testing.T) {
	// Test that YYYY-MM-DD is parsed as start of day UTC
	result, err := parseTimestamp("2023-06-15")
	require.NoError(t, err)

	// Verify it's the start of the day
	parsedTime := time.Unix(result, 0).UTC()
	require.Equal(t, 2023, parsedTime.Year())
	require.Equal(t, time.June, parsedTime.Month())
	require.Equal(t, 15, parsedTime.Day())
	require.Equal(t, 0, parsedTime.Hour())
	require.Equal(t, 0, parsedTime.Minute())
	require.Equal(t, 0, parsedTime.Second())
}
