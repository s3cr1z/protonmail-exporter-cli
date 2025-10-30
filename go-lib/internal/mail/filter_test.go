// Copyright (c) 2023 Proton AG
//
// This file is part of Proton Export Tool.
//
// Proton Export Tool is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Proton Export Tool is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Proton Export Tool.  If not, see <https://www.gnu.org/licenses/>.

package mail

import (
	"net/mail"
	"testing"
	"time"

	"github.com/ProtonMail/go-proton-api"
	"github.com/stretchr/testify/assert"
)

func TestExportFilter_IsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		filter *ExportFilter
		want   bool
	}{
		{
			name:   "empty filter",
			filter: NewExportFilter(),
			want:   true,
		},
		{
			name: "filter with labels",
			filter: &ExportFilter{
				LabelIDs: []string{"0"},
			},
			want: false,
		},
		{
			name: "filter with senders",
			filter: &ExportFilter{
				Senders: []string{"test@example.com"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.IsEmpty()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExportFilter_MatchesLabels(t *testing.T) {
	tests := []struct {
		name   string
		filter *ExportFilter
		msg    proton.MessageMetadata
		want   bool
	}{
		{
			name: "no label filter - matches all",
			filter: &ExportFilter{
				LabelIDs: []string{},
			},
			msg: proton.MessageMetadata{
				LabelIDs: []string{"0", "2"},
			},
			want: true,
		},
		{
			name: "single label matches",
			filter: &ExportFilter{
				LabelIDs: []string{"0"},
			},
			msg: proton.MessageMetadata{
				LabelIDs: []string{"0", "2"},
			},
			want: true,
		},
		{
			name: "multiple labels - one matches",
			filter: &ExportFilter{
				LabelIDs: []string{"0", "1"},
			},
			msg: proton.MessageMetadata{
				LabelIDs: []string{"2", "0"},
			},
			want: true,
		},
		{
			name: "no match",
			filter: &ExportFilter{
				LabelIDs: []string{"1"},
			},
			msg: proton.MessageMetadata{
				LabelIDs: []string{"0", "2"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.msg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExportFilter_MatchesSender(t *testing.T) {
	tests := []struct {
		name   string
		filter *ExportFilter
		msg    proton.MessageMetadata
		want   bool
	}{
		{
			name: "sender address matches",
			filter: &ExportFilter{
				Senders: []string{"alice@example.com"},
			},
			msg: proton.MessageMetadata{
				Sender: &mail.Address{
					Address: "alice@example.com",
					Name:    "Alice",
				},
			},
			want: true,
		},
		{
			name: "sender name matches",
			filter: &ExportFilter{
				Senders: []string{"alice"},
			},
			msg: proton.MessageMetadata{
				Sender: &mail.Address{
					Address: "test@example.com",
					Name:    "Alice Smith",
				},
			},
			want: true,
		},
		{
			name: "partial match on address",
			filter: &ExportFilter{
				Senders: []string{"example.com"},
			},
			msg: proton.MessageMetadata{
				Sender: &mail.Address{
					Address: "alice@example.com",
				},
			},
			want: true,
		},
		{
			name: "case insensitive match",
			filter: &ExportFilter{
				Senders: []string{"ALICE@EXAMPLE.COM"},
			},
			msg: proton.MessageMetadata{
				Sender: &mail.Address{
					Address: "alice@example.com",
				},
			},
			want: true,
		},
		{
			name: "multiple senders - one matches",
			filter: &ExportFilter{
				Senders: []string{"bob@example.com", "alice@example.com"},
			},
			msg: proton.MessageMetadata{
				Sender: &mail.Address{
					Address: "alice@example.com",
				},
			},
			want: true,
		},
		{
			name: "no match",
			filter: &ExportFilter{
				Senders: []string{"bob@example.com"},
			},
			msg: proton.MessageMetadata{
				Sender: &mail.Address{
					Address: "alice@example.com",
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.msg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExportFilter_MatchesRecipient(t *testing.T) {
	tests := []struct {
		name   string
		filter *ExportFilter
		msg    proton.MessageMetadata
		want   bool
	}{
		{
			name: "To recipient matches",
			filter: &ExportFilter{
				Recipients: []string{"bob@example.com"},
			},
			msg: proton.MessageMetadata{
				ToList: []*mail.Address{
					{Address: "bob@example.com", Name: "Bob"},
				},
			},
			want: true,
		},
		{
			name: "CC recipient matches",
			filter: &ExportFilter{
				Recipients: []string{"charlie@example.com"},
			},
			msg: proton.MessageMetadata{
				ToList: []*mail.Address{
					{Address: "bob@example.com"},
				},
				CCList: []*mail.Address{
					{Address: "charlie@example.com"},
				},
			},
			want: true,
		},
		{
			name: "BCC recipient matches",
			filter: &ExportFilter{
				Recipients: []string{"dave@example.com"},
			},
			msg: proton.MessageMetadata{
				BCCList: []*mail.Address{
					{Address: "dave@example.com"},
				},
			},
			want: true,
		},
		{
			name: "recipient name matches",
			filter: &ExportFilter{
				Recipients: []string{"Bob Smith"},
			},
			msg: proton.MessageMetadata{
				ToList: []*mail.Address{
					{Address: "bob@example.com", Name: "Bob Smith"},
				},
			},
			want: true,
		},
		{
			name: "no match",
			filter: &ExportFilter{
				Recipients: []string{"eve@example.com"},
			},
			msg: proton.MessageMetadata{
				ToList: []*mail.Address{
					{Address: "bob@example.com"},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.msg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExportFilter_MatchesDateRange(t *testing.T) {
	baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	before := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	after := time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC)

	tests := []struct {
		name   string
		filter *ExportFilter
		msg    proton.MessageMetadata
		want   bool
	}{
		{
			name: "within date range",
			filter: &ExportFilter{
				DateFrom: &before,
				DateTo:   &after,
			},
			msg: proton.MessageMetadata{
				Time: baseTime.Unix(),
			},
			want: true,
		},
		{
			name: "before date range",
			filter: &ExportFilter{
				DateFrom: &baseTime,
			},
			msg: proton.MessageMetadata{
				Time: before.Unix(),
			},
			want: false,
		},
		{
			name: "after date range",
			filter: &ExportFilter{
				DateTo: &baseTime,
			},
			msg: proton.MessageMetadata{
				Time: after.Unix(),
			},
			want: false,
		},
		{
			name: "only DateFrom set - matches",
			filter: &ExportFilter{
				DateFrom: &before,
			},
			msg: proton.MessageMetadata{
				Time: baseTime.Unix(),
			},
			want: true,
		},
		{
			name: "only DateTo set - matches",
			filter: &ExportFilter{
				DateTo: &after,
			},
			msg: proton.MessageMetadata{
				Time: baseTime.Unix(),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.msg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExportFilter_MatchesSubject(t *testing.T) {
	tests := []struct {
		name   string
		filter *ExportFilter
		msg    proton.MessageMetadata
		want   bool
	}{
		{
			name: "exact keyword match",
			filter: &ExportFilter{
				SubjectKeywords: []string{"invoice"},
			},
			msg: proton.MessageMetadata{
				Subject: "Invoice for June",
			},
			want: true,
		},
		{
			name: "case insensitive match",
			filter: &ExportFilter{
				SubjectKeywords: []string{"INVOICE"},
			},
			msg: proton.MessageMetadata{
				Subject: "invoice for June",
			},
			want: true,
		},
		{
			name: "partial match",
			filter: &ExportFilter{
				SubjectKeywords: []string{"inv"},
			},
			msg: proton.MessageMetadata{
				Subject: "Invoice for June",
			},
			want: true,
		},
		{
			name: "multiple keywords - one matches",
			filter: &ExportFilter{
				SubjectKeywords: []string{"receipt", "invoice"},
			},
			msg: proton.MessageMetadata{
				Subject: "Invoice for June",
			},
			want: true,
		},
		{
			name: "no match",
			filter: &ExportFilter{
				SubjectKeywords: []string{"receipt"},
			},
			msg: proton.MessageMetadata{
				Subject: "Invoice for June",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.msg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExportFilter_CombineLogic(t *testing.T) {
	baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	before := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	after := time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC)

	tests := []struct {
		name   string
		filter *ExportFilter
		msg    proton.MessageMetadata
		want   bool
	}{
		{
			name: "AND logic - all match",
			filter: &ExportFilter{
				LabelIDs:     []string{"0"},
				Senders:      []string{"alice@example.com"},
				DateFrom:     &before,
				DateTo:       &after,
				CombineLogic: FilterLogicAND,
			},
			msg: proton.MessageMetadata{
				LabelIDs: []string{"0"},
				Sender: &mail.Address{
					Address: "alice@example.com",
				},
				Time: baseTime.Unix(),
			},
			want: true,
		},
		{
			name: "AND logic - one fails",
			filter: &ExportFilter{
				LabelIDs:     []string{"0"},
				Senders:      []string{"bob@example.com"},
				CombineLogic: FilterLogicAND,
			},
			msg: proton.MessageMetadata{
				LabelIDs: []string{"0"},
				Sender: &mail.Address{
					Address: "alice@example.com",
				},
			},
			want: false,
		},
		{
			name: "OR logic - one matches",
			filter: &ExportFilter{
				LabelIDs:     []string{"1"},
				Senders:      []string{"alice@example.com"},
				CombineLogic: FilterLogicOR,
			},
			msg: proton.MessageMetadata{
				LabelIDs: []string{"0"},
				Sender: &mail.Address{
					Address: "alice@example.com",
				},
			},
			want: true,
		},
		{
			name: "OR logic - all fail",
			filter: &ExportFilter{
				LabelIDs:     []string{"1"},
				Senders:      []string{"bob@example.com"},
				CombineLogic: FilterLogicOR,
			},
			msg: proton.MessageMetadata{
				LabelIDs: []string{"0"},
				Sender: &mail.Address{
					Address: "alice@example.com",
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.msg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExportFilter_NeedsClientSideFiltering(t *testing.T) {
	tests := []struct {
		name   string
		filter *ExportFilter
		want   bool
	}{
		{
			name:   "empty filter",
			filter: NewExportFilter(),
			want:   false,
		},
		{
			name: "single label only",
			filter: &ExportFilter{
				LabelIDs: []string{"0"},
			},
			want: false,
		},
		{
			name: "multiple labels",
			filter: &ExportFilter{
				LabelIDs: []string{"0", "1"},
			},
			want: true,
		},
		{
			name: "sender filter",
			filter: &ExportFilter{
				Senders: []string{"alice@example.com"},
			},
			want: true,
		},
		{
			name: "date filter",
			filter: &ExportFilter{
				DateFrom: &time.Time{},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.NeedsClientSideFiltering()
			assert.Equal(t, tt.want, got)
		})
	}
}
