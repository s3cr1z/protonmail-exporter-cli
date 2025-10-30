package mail

import (
	"net/mail"
	"testing"

	"github.com/ProtonMail/go-proton-api"
	"github.com/stretchr/testify/require"
)

func TestFilterOptions_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		expected bool
	}{
		{
			name:     "completely empty",
			filter:   FilterOptions{},
			expected: true,
		},
		{
			name: "has label filter",
			filter: FilterOptions{
				LabelIDs: []string{"0"},
			},
			expected: false,
		},
		{
			name: "has after timestamp",
			filter: FilterOptions{
				After: 1234567890,
			},
			expected: false,
		},
		{
			name: "has before timestamp",
			filter: FilterOptions{
				Before: 1234567890,
			},
			expected: false,
		},
		{
			name: "has from filter",
			filter: FilterOptions{
				From: []string{"test@example.com"},
			},
			expected: false,
		},
		{
			name: "has to filter",
			filter: FilterOptions{
				To: []string{"test@example.com"},
			},
			expected: false,
		},
		{
			name: "has from domain filter",
			filter: FilterOptions{
				FromDomains: []string{"example.com"},
			},
			expected: false,
		},
		{
			name: "has to domain filter",
			filter: FilterOptions{
				ToDomains: []string{"example.com"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.filter.IsEmpty())
		})
	}
}

func TestFilterOptions_MatchesLabels(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		labels   []string
		expected bool
	}{
		{
			name:     "no filter matches any labels",
			filter:   FilterOptions{},
			labels:   []string{"0", "5"},
			expected: true,
		},
		{
			name: "single label matches",
			filter: FilterOptions{
				LabelIDs: []string{"0"},
			},
			labels:   []string{"0", "5"},
			expected: true,
		},
		{
			name: "one of multiple labels matches",
			filter: FilterOptions{
				LabelIDs: []string{"0", "2"},
			},
			labels:   []string{"2", "5"},
			expected: true,
		},
		{
			name: "no label matches",
			filter: FilterOptions{
				LabelIDs: []string{"0", "2"},
			},
			labels:   []string{"3", "5"},
			expected: false,
		},
		{
			name: "empty message labels with filter",
			filter: FilterOptions{
				LabelIDs: []string{"0"},
			},
			labels:   []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.filter.matchesLabels(tt.labels))
		})
	}
}

func TestFilterOptions_MatchesDateRange(t *testing.T) {
	tests := []struct {
		name      string
		filter    FilterOptions
		timestamp int64
		expected  bool
	}{
		{
			name:      "no filter matches any timestamp",
			filter:    FilterOptions{},
			timestamp: 1000,
			expected:  true,
		},
		{
			name: "after filter - message is after",
			filter: FilterOptions{
				After: 1000,
			},
			timestamp: 2000,
			expected:  true,
		},
		{
			name: "after filter - message is exactly at boundary",
			filter: FilterOptions{
				After: 1000,
			},
			timestamp: 1000,
			expected:  true,
		},
		{
			name: "after filter - message is before",
			filter: FilterOptions{
				After: 1000,
			},
			timestamp: 500,
			expected:  false,
		},
		{
			name: "before filter - message is before",
			filter: FilterOptions{
				Before: 1000,
			},
			timestamp: 500,
			expected:  true,
		},
		{
			name: "before filter - message is at boundary",
			filter: FilterOptions{
				Before: 1000,
			},
			timestamp: 1000,
			expected:  false,
		},
		{
			name: "before filter - message is after",
			filter: FilterOptions{
				Before: 1000,
			},
			timestamp: 2000,
			expected:  false,
		},
		{
			name: "both filters - message in range",
			filter: FilterOptions{
				After:  1000,
				Before: 2000,
			},
			timestamp: 1500,
			expected:  true,
		},
		{
			name: "both filters - message before range",
			filter: FilterOptions{
				After:  1000,
				Before: 2000,
			},
			timestamp: 500,
			expected:  false,
		},
		{
			name: "both filters - message after range",
			filter: FilterOptions{
				After:  1000,
				Before: 2000,
			},
			timestamp: 3000,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.filter.matchesDateRange(tt.timestamp))
		})
	}
}

func TestFilterOptions_MatchesSender(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		sender   *mail.Address
		expected bool
	}{
		{
			name:     "no filter matches any sender",
			filter:   FilterOptions{},
			sender:   &mail.Address{Address: "test@example.com"},
			expected: true,
		},
		{
			name: "sender matches exactly",
			filter: FilterOptions{
				From: []string{"test@example.com"},
			},
			sender:   &mail.Address{Address: "test@example.com"},
			expected: true,
		},
		{
			name: "sender matches case-insensitive",
			filter: FilterOptions{
				From: []string{"Test@Example.COM"},
			},
			sender:   &mail.Address{Address: "test@example.com"},
			expected: true,
		},
		{
			name: "sender matches with whitespace",
			filter: FilterOptions{
				From: []string{" test@example.com "},
			},
			sender:   &mail.Address{Address: "test@example.com"},
			expected: true,
		},
		{
			name: "sender doesn't match",
			filter: FilterOptions{
				From: []string{"other@example.com"},
			},
			sender:   &mail.Address{Address: "test@example.com"},
			expected: false,
		},
		{
			name: "sender matches one of multiple",
			filter: FilterOptions{
				From: []string{"other@example.com", "test@example.com"},
			},
			sender:   &mail.Address{Address: "test@example.com"},
			expected: true,
		},
		{
			name: "nil sender with filter",
			filter: FilterOptions{
				From: []string{"test@example.com"},
			},
			sender:   nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.filter.matchesSender(tt.sender))
		})
	}
}

func TestFilterOptions_MatchesRecipients(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		toList   []*mail.Address
		ccList   []*mail.Address
		bccList  []*mail.Address
		expected bool
	}{
		{
			name:   "no filter matches any recipient",
			filter: FilterOptions{},
			toList: []*mail.Address{
				{Address: "test@example.com"},
			},
			expected: true,
		},
		{
			name: "to recipient matches",
			filter: FilterOptions{
				To: []string{"test@example.com"},
			},
			toList: []*mail.Address{
				{Address: "test@example.com"},
			},
			expected: true,
		},
		{
			name: "cc recipient matches",
			filter: FilterOptions{
				To: []string{"test@example.com"},
			},
			ccList: []*mail.Address{
				{Address: "test@example.com"},
			},
			expected: true,
		},
		{
			name: "bcc recipient matches",
			filter: FilterOptions{
				To: []string{"test@example.com"},
			},
			bccList: []*mail.Address{
				{Address: "test@example.com"},
			},
			expected: true,
		},
		{
			name: "one of many recipients matches",
			filter: FilterOptions{
				To: []string{"test@example.com"},
			},
			toList: []*mail.Address{
				{Address: "other1@example.com"},
				{Address: "test@example.com"},
				{Address: "other2@example.com"},
			},
			expected: true,
		},
		{
			name: "recipient matches case-insensitive",
			filter: FilterOptions{
				To: []string{"Test@Example.COM"},
			},
			toList: []*mail.Address{
				{Address: "test@example.com"},
			},
			expected: true,
		},
		{
			name: "no recipient matches",
			filter: FilterOptions{
				To: []string{"missing@example.com"},
			},
			toList: []*mail.Address{
				{Address: "test@example.com"},
			},
			ccList: []*mail.Address{
				{Address: "other@example.com"},
			},
			expected: false,
		},
		{
			name: "empty recipients with filter",
			filter: FilterOptions{
				To: []string{"test@example.com"},
			},
			toList:  []*mail.Address{},
			ccList:  []*mail.Address{},
			bccList: []*mail.Address{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.filter.matchesRecipients(tt.toList, tt.ccList, tt.bccList))
		})
	}
}

func TestFilterOptions_MatchesSenderDomain(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		sender   *mail.Address
		expected bool
	}{
		{
			name:     "no filter matches any domain",
			filter:   FilterOptions{},
			sender:   &mail.Address{Address: "test@example.com"},
			expected: true,
		},
		{
			name: "sender domain matches",
			filter: FilterOptions{
				FromDomains: []string{"example.com"},
			},
			sender:   &mail.Address{Address: "test@example.com"},
			expected: true,
		},
		{
			name: "sender domain matches case-insensitive",
			filter: FilterOptions{
				FromDomains: []string{"Example.COM"},
			},
			sender:   &mail.Address{Address: "test@example.com"},
			expected: true,
		},
		{
			name: "sender domain matches with whitespace",
			filter: FilterOptions{
				FromDomains: []string{" example.com "},
			},
			sender:   &mail.Address{Address: "test@example.com"},
			expected: true,
		},
		{
			name: "sender domain doesn't match",
			filter: FilterOptions{
				FromDomains: []string{"other.com"},
			},
			sender:   &mail.Address{Address: "test@example.com"},
			expected: false,
		},
		{
			name: "sender domain matches one of multiple",
			filter: FilterOptions{
				FromDomains: []string{"other.com", "example.com", "third.com"},
			},
			sender:   &mail.Address{Address: "test@example.com"},
			expected: true,
		},
		{
			name: "nil sender with filter",
			filter: FilterOptions{
				FromDomains: []string{"example.com"},
			},
			sender:   nil,
			expected: false,
		},
		{
			name: "invalid email with filter",
			filter: FilterOptions{
				FromDomains: []string{"example.com"},
			},
			sender:   &mail.Address{Address: "invalid-email"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.filter.matchesSenderDomain(tt.sender))
		})
	}
}

func TestFilterOptions_MatchesRecipientDomains(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		toList   []*mail.Address
		ccList   []*mail.Address
		bccList  []*mail.Address
		expected bool
	}{
		{
			name:   "no filter matches any domain",
			filter: FilterOptions{},
			toList: []*mail.Address{
				{Address: "test@example.com"},
			},
			expected: true,
		},
		{
			name: "to recipient domain matches",
			filter: FilterOptions{
				ToDomains: []string{"example.com"},
			},
			toList: []*mail.Address{
				{Address: "test@example.com"},
			},
			expected: true,
		},
		{
			name: "cc recipient domain matches",
			filter: FilterOptions{
				ToDomains: []string{"example.com"},
			},
			ccList: []*mail.Address{
				{Address: "test@example.com"},
			},
			expected: true,
		},
		{
			name: "bcc recipient domain matches",
			filter: FilterOptions{
				ToDomains: []string{"example.com"},
			},
			bccList: []*mail.Address{
				{Address: "test@example.com"},
			},
			expected: true,
		},
		{
			name: "one of many recipient domains matches",
			filter: FilterOptions{
				ToDomains: []string{"example.com"},
			},
			toList: []*mail.Address{
				{Address: "test1@other.com"},
				{Address: "test2@example.com"},
				{Address: "test3@third.com"},
			},
			expected: true,
		},
		{
			name: "recipient domain matches case-insensitive",
			filter: FilterOptions{
				ToDomains: []string{"Example.COM"},
			},
			toList: []*mail.Address{
				{Address: "test@example.com"},
			},
			expected: true,
		},
		{
			name: "no recipient domain matches",
			filter: FilterOptions{
				ToDomains: []string{"missing.com"},
			},
			toList: []*mail.Address{
				{Address: "test@example.com"},
			},
			ccList: []*mail.Address{
				{Address: "other@other.com"},
			},
			expected: false,
		},
		{
			name: "empty recipients with filter",
			filter: FilterOptions{
				ToDomains: []string{"example.com"},
			},
			toList:   []*mail.Address{},
			ccList:   []*mail.Address{},
			bccList:  []*mail.Address{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.filter.matchesRecipientDomains(tt.toList, tt.ccList, tt.bccList))
		})
	}
}

func TestFilterOptions_MatchesMessage_Combined(t *testing.T) {
	// Test combined filters (AND semantics across different types)
	tests := []struct {
		name     string
		filter   FilterOptions
		message  proton.MessageMetadata
		expected bool
	}{
		{
			name:   "empty filter matches all",
			filter: FilterOptions{},
			message: proton.MessageMetadata{
				ID:       "msg1",
				LabelIDs: []string{"0"},
				Time:     1000,
				Sender:   &mail.Address{Address: "test@example.com"},
			},
			expected: true,
		},
		{
			name: "all filters match",
			filter: FilterOptions{
				LabelIDs: []string{"0"},
				After:    1000,
				Before:   2000,
				From:     []string{"test@example.com"},
			},
			message: proton.MessageMetadata{
				ID:       "msg1",
				LabelIDs: []string{"0"},
				Time:     1500,
				Sender:   &mail.Address{Address: "test@example.com"},
			},
			expected: true,
		},
		{
			name: "label doesn't match - fail",
			filter: FilterOptions{
				LabelIDs: []string{"2"},
				After:    1000,
				From:     []string{"test@example.com"},
			},
			message: proton.MessageMetadata{
				ID:       "msg1",
				LabelIDs: []string{"0"},
				Time:     1500,
				Sender:   &mail.Address{Address: "test@example.com"},
			},
			expected: false,
		},
		{
			name: "date out of range - fail",
			filter: FilterOptions{
				LabelIDs: []string{"0"},
				After:    1000,
				Before:   1200,
			},
			message: proton.MessageMetadata{
				ID:       "msg1",
				LabelIDs: []string{"0"},
				Time:     1500,
			},
			expected: false,
		},
		{
			name: "sender doesn't match - fail",
			filter: FilterOptions{
				LabelIDs: []string{"0"},
				From:     []string{"other@example.com"},
			},
			message: proton.MessageMetadata{
				ID:       "msg1",
				LabelIDs: []string{"0"},
				Time:     1500,
				Sender:   &mail.Address{Address: "test@example.com"},
			},
			expected: false,
		},
		{
			name: "complex filter with domains",
			filter: FilterOptions{
				FromDomains: []string{"example.com"},
				ToDomains:   []string{"proton.me"},
			},
			message: proton.MessageMetadata{
				ID:     "msg1",
				Time:   1500,
				Sender: &mail.Address{Address: "test@example.com"},
				ToList: []*mail.Address{
					{Address: "user@proton.me"},
				},
			},
			expected: true,
		},
		{
			name: "recipient domain doesn't match - fail",
			filter: FilterOptions{
				ToDomains: []string{"gmail.com"},
			},
			message: proton.MessageMetadata{
				ID:   "msg1",
				Time: 1500,
				ToList: []*mail.Address{
					{Address: "user@proton.me"},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.filter.MatchesMessage(tt.message))
		})
	}
}

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "valid email",
			email:    "test@example.com",
			expected: "example.com",
		},
		{
			name:     "valid email with subdomain",
			email:    "test@mail.example.com",
			expected: "mail.example.com",
		},
		{
			name:     "email with capital letters",
			email:    "Test@Example.COM",
			expected: "example.com",
		},
		{
			name:     "email with whitespace",
			email:    " test@example.com ",
			expected: "example.com",
		},
		{
			name:     "no @ sign",
			email:    "testexample.com",
			expected: "",
		},
		{
			name:     "@ at end",
			email:    "test@",
			expected: "",
		},
		{
			name:     "empty string",
			email:    "",
			expected: "",
		},
		{
			name:     "multiple @ signs",
			email:    "test@domain@example.com",
			expected: "example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, extractDomain(tt.email))
		})
	}
}

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "lowercase email",
			email:    "test@example.com",
			expected: "test@example.com",
		},
		{
			name:     "uppercase email",
			email:    "TEST@EXAMPLE.COM",
			expected: "test@example.com",
		},
		{
			name:     "mixed case email",
			email:    "Test@Example.Com",
			expected: "test@example.com",
		},
		{
			name:     "email with whitespace",
			email:    "  test@example.com  ",
			expected: "test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, normalizeEmail(tt.email))
		})
	}
}
