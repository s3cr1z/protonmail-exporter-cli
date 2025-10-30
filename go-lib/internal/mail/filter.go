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
	"strings"
	"time"

	"github.com/ProtonMail/go-proton-api"
)

// FilterLogic defines how multiple filter criteria are combined
type FilterLogic int

const (
	// FilterLogicAND requires all criteria to match
	FilterLogicAND FilterLogic = iota
	// FilterLogicOR requires at least one criterion to match
	FilterLogicOR
)

// ExportFilter defines comprehensive filtering options for email export
type ExportFilter struct {
	// Label/Folder filtering (OR logic within labels)
	LabelIDs []string

	// Sender filtering - matches if any sender matches (OR logic)
	Senders []string

	// Recipient filtering - matches To, CC, or BCC (OR logic)
	Recipients []string

	// Date range filtering
	DateFrom *time.Time
	DateTo   *time.Time

	// Subject keyword filtering (substring match, case-insensitive)
	SubjectKeywords []string

	// Logic for combining different filter types (labels AND dates AND senders, etc.)
	// Note: Within each filter type (e.g., multiple senders), OR logic is used
	CombineLogic FilterLogic
}

// NewExportFilter creates a new ExportFilter with default values
func NewExportFilter() *ExportFilter {
	return &ExportFilter{
		CombineLogic: FilterLogicAND,
	}
}

// IsEmpty returns true if no filters are set
func (f *ExportFilter) IsEmpty() bool {
	return len(f.LabelIDs) == 0 &&
		len(f.Senders) == 0 &&
		len(f.Recipients) == 0 &&
		f.DateFrom == nil &&
		f.DateTo == nil &&
		len(f.SubjectKeywords) == 0
}

// Matches determines if a message matches the filter criteria
func (f *ExportFilter) Matches(msg proton.MessageMetadata) bool {
	if f.IsEmpty() {
		return true // No filter means match all
	}

	// When using AND logic, all criteria that are set must match
	// When using OR logic, at least one criterion must match

	criteriaResults := make([]bool, 0)

	// Check label filter if set
	if len(f.LabelIDs) > 0 {
		criteriaResults = append(criteriaResults, f.matchesLabels(msg))
	}

	// Check sender filter if set
	if len(f.Senders) > 0 {
		criteriaResults = append(criteriaResults, f.matchesSender(msg))
	}

	// Check recipient filter if set
	if len(f.Recipients) > 0 {
		criteriaResults = append(criteriaResults, f.matchesRecipient(msg))
	}

	// Check date range filter if set
	if f.DateFrom != nil || f.DateTo != nil {
		criteriaResults = append(criteriaResults, f.matchesDateRange(msg))
	}

	// Check subject keywords if set
	if len(f.SubjectKeywords) > 0 {
		criteriaResults = append(criteriaResults, f.matchesSubject(msg))
	}

	if len(criteriaResults) == 0 {
		return true // No criteria set
	}

	// Apply combine logic
	if f.CombineLogic == FilterLogicAND {
		// All criteria must match
		for _, result := range criteriaResults {
			if !result {
				return false
			}
		}
		return true
	} else {
		// At least one criterion must match
		for _, result := range criteriaResults {
			if result {
				return true
			}
		}
		return false
	}
}

// matchesLabels checks if message has any of the requested labels (OR logic)
func (f *ExportFilter) matchesLabels(msg proton.MessageMetadata) bool {
	for _, requestedLabel := range f.LabelIDs {
		for _, msgLabel := range msg.LabelIDs {
			if msgLabel == requestedLabel {
				return true
			}
		}
	}
	return false
}

// matchesSender checks if message sender matches any of the filter senders (OR logic)
func (f *ExportFilter) matchesSender(msg proton.MessageMetadata) bool {
	senderAddr := strings.ToLower(msg.Sender.Address)
	senderName := strings.ToLower(msg.Sender.Name)

	for _, filterSender := range f.Senders {
		filterSender = strings.ToLower(strings.TrimSpace(filterSender))
		if filterSender == "" {
			continue
		}

		// Match against address or name
		if strings.Contains(senderAddr, filterSender) || strings.Contains(senderName, filterSender) {
			return true
		}
	}
	return false
}

// matchesRecipient checks if message has any recipient matching the filter (OR logic)
func (f *ExportFilter) matchesRecipient(msg proton.MessageMetadata) bool {
	// Collect all recipients (To, CC, BCC)
	allRecipients := make([]string, 0)

	for _, to := range msg.ToList {
		allRecipients = append(allRecipients, strings.ToLower(to.Address), strings.ToLower(to.Name))
	}
	for _, cc := range msg.CCList {
		allRecipients = append(allRecipients, strings.ToLower(cc.Address), strings.ToLower(cc.Name))
	}
	for _, bcc := range msg.BCCList {
		allRecipients = append(allRecipients, strings.ToLower(bcc.Address), strings.ToLower(bcc.Name))
	}

	for _, filterRecipient := range f.Recipients {
		filterRecipient = strings.ToLower(strings.TrimSpace(filterRecipient))
		if filterRecipient == "" {
			continue
		}

		for _, recipient := range allRecipients {
			if strings.Contains(recipient, filterRecipient) {
				return true
			}
		}
	}
	return false
}

// matchesDateRange checks if message falls within the specified date range
func (f *ExportFilter) matchesDateRange(msg proton.MessageMetadata) bool {
	msgTime := time.Unix(msg.Time, 0)

	if f.DateFrom != nil && msgTime.Before(*f.DateFrom) {
		return false
	}

	if f.DateTo != nil && msgTime.After(*f.DateTo) {
		return false
	}

	return true
}

// matchesSubject checks if message subject contains any of the keywords (OR logic, case-insensitive)
func (f *ExportFilter) matchesSubject(msg proton.MessageMetadata) bool {
	subject := strings.ToLower(msg.Subject)

	for _, keyword := range f.SubjectKeywords {
		keyword = strings.ToLower(strings.TrimSpace(keyword))
		if keyword == "" {
			continue
		}

		if strings.Contains(subject, keyword) {
			return true
		}
	}
	return false
}

// GetServerSideFilter extracts parameters that can be sent to the ProtonMail API
// Returns a proton.MessageFilter for server-side filtering where supported
func (f *ExportFilter) GetServerSideFilter() proton.MessageFilter {
	filter := proton.MessageFilter{
		Desc: true, // Always sort descending by time
	}

	// Server-side label filtering (single label only)
	if len(f.LabelIDs) == 1 {
		filter.LabelID = f.LabelIDs[0]
	}
	// Note: Multiple labels, senders, recipients, dates require client-side filtering

	return filter
}

// NeedsClientSideFiltering returns true if any filters require client-side processing
func (f *ExportFilter) NeedsClientSideFiltering() bool {
	// Multiple labels require client-side filtering
	if len(f.LabelIDs) > 1 {
		return true
	}

	// Any other filter type requires client-side filtering
	return len(f.Senders) > 0 ||
		len(f.Recipients) > 0 ||
		f.DateFrom != nil ||
		f.DateTo != nil ||
		len(f.SubjectKeywords) > 0
}
