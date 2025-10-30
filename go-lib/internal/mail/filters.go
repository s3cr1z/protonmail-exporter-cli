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
	"strings"

	"github.com/ProtonMail/go-proton-api"
)

// FilterOptions contains all filtering criteria for message export.
// If all fields are empty/zero, no filtering is applied (export all messages).
// Filters are combined with AND semantics across different types.
// Multiple values within the same field type are combined with OR semantics.
type FilterOptions struct {
	// LabelIDs filters messages by label/folder IDs.
	// Empty slice = no label filtering.
	// Non-empty = message must have at least one of these labels.
	LabelIDs []string

	// After filters messages by timestamp (Unix epoch seconds).
	// Zero = no lower bound. Messages must have Time >= After.
	After int64

	// Before filters messages by timestamp (Unix epoch seconds).
	// Zero = no upper bound. Messages must have Time < Before.
	Before int64

	// From filters by sender email addresses.
	// Empty slice = no sender filtering.
	// Non-empty = sender must match at least one address (case-insensitive).
	From []string

	// To filters by recipient email addresses (To/CC/BCC).
	// Empty slice = no recipient filtering.
	// Non-empty = at least one recipient must match (case-insensitive).
	To []string

	// FromDomains filters by sender domain.
	// Empty slice = no sender domain filtering.
	// Non-empty = sender domain must match at least one (case-insensitive).
	FromDomains []string

	// ToDomains filters by recipient domains (To/CC/BCC).
	// Empty slice = no recipient domain filtering.
	// Non-empty = at least one recipient domain must match (case-insensitive).
	ToDomains []string
}

// IsEmpty returns true if no filters are configured.
func (f *FilterOptions) IsEmpty() bool {
	return len(f.LabelIDs) == 0 &&
		f.After == 0 &&
		f.Before == 0 &&
		len(f.From) == 0 &&
		len(f.To) == 0 &&
		len(f.FromDomains) == 0 &&
		len(f.ToDomains) == 0
}

// MatchesMessage returns true if the message matches all configured filters.
func (f *FilterOptions) MatchesMessage(msg proton.MessageMetadata) bool {
	// If no filters, all messages match
	if f.IsEmpty() {
		return true
	}

	// Check label filter (OR within labels)
	if !f.matchesLabels(msg.LabelIDs) {
		return false
	}

	// Check date range filters (AND)
	if !f.matchesDateRange(msg.Time) {
		return false
	}

	// Check sender filters (OR within senders)
	if !f.matchesSender(msg.Sender) {
		return false
	}

	// Check recipient filters (OR within recipients)
	if !f.matchesRecipients(msg.ToList, msg.CCList, msg.BCCList) {
		return false
	}

	// Check sender domain filters (OR within domains)
	if !f.matchesSenderDomain(msg.Sender) {
		return false
	}

	// Check recipient domain filters (OR within domains)
	if !f.matchesRecipientDomains(msg.ToList, msg.CCList, msg.BCCList) {
		return false
	}

	return true
}

// matchesLabels checks if message has at least one of the requested labels.
// Returns true if no label filter is configured.
func (f *FilterOptions) matchesLabels(messageLabelIDs []string) bool {
	if len(f.LabelIDs) == 0 {
		return true // No filter
	}

	for _, requestedLabel := range f.LabelIDs {
		for _, msgLabel := range messageLabelIDs {
			if msgLabel == requestedLabel {
				return true
			}
		}
	}

	return false
}

// matchesDateRange checks if message timestamp is within configured range.
// Returns true if no date range filter is configured.
func (f *FilterOptions) matchesDateRange(timestamp int64) bool {
	if f.After != 0 && timestamp < f.After {
		return false
	}
	if f.Before != 0 && timestamp >= f.Before {
		return false
	}
	return true
}

// matchesSender checks if sender email matches any configured sender filter.
// Returns true if no sender filter is configured.
func (f *FilterOptions) matchesSender(sender *mail.Address) bool {
	if len(f.From) == 0 {
		return true // No filter
	}

	if sender == nil {
		return false
	}

	senderEmail := normalizeEmail(sender.Address)
	for _, filterEmail := range f.From {
		if senderEmail == normalizeEmail(filterEmail) {
			return true
		}
	}

	return false
}

// matchesRecipients checks if any recipient matches configured recipient filters.
// Returns true if no recipient filter is configured.
func (f *FilterOptions) matchesRecipients(toList, ccList, bccList []*mail.Address) bool {
	if len(f.To) == 0 {
		return true // No filter
	}

	// Combine all recipients
	allRecipients := make([]*mail.Address, 0, len(toList)+len(ccList)+len(bccList))
	allRecipients = append(allRecipients, toList...)
	allRecipients = append(allRecipients, ccList...)
	allRecipients = append(allRecipients, bccList...)

	// Check if any recipient matches any filter
	for _, recipient := range allRecipients {
		if recipient == nil {
			continue
		}
		recipientEmail := normalizeEmail(recipient.Address)
		for _, filterEmail := range f.To {
			if recipientEmail == normalizeEmail(filterEmail) {
				return true
			}
		}
	}

	return false
}

// matchesSenderDomain checks if sender domain matches any configured domain filter.
// Returns true if no sender domain filter is configured.
func (f *FilterOptions) matchesSenderDomain(sender *mail.Address) bool {
	if len(f.FromDomains) == 0 {
		return true // No filter
	}

	if sender == nil {
		return false
	}

	senderDomain := extractDomain(sender.Address)
	if senderDomain == "" {
		return false
	}

	for _, filterDomain := range f.FromDomains {
		if senderDomain == normalizeDomain(filterDomain) {
			return true
		}
	}

	return false
}

// matchesRecipientDomains checks if any recipient domain matches configured domain filters.
// Returns true if no recipient domain filter is configured.
func (f *FilterOptions) matchesRecipientDomains(toList, ccList, bccList []*mail.Address) bool {
	if len(f.ToDomains) == 0 {
		return true // No filter
	}

	// Combine all recipients
	allRecipients := make([]*mail.Address, 0, len(toList)+len(ccList)+len(bccList))
	allRecipients = append(allRecipients, toList...)
	allRecipients = append(allRecipients, ccList...)
	allRecipients = append(allRecipients, bccList...)

	// Check if any recipient domain matches any filter
	for _, recipient := range allRecipients {
		if recipient == nil {
			continue
		}
		recipientDomain := extractDomain(recipient.Address)
		if recipientDomain == "" {
			continue
		}
		for _, filterDomain := range f.ToDomains {
			if recipientDomain == normalizeDomain(filterDomain) {
				return true
			}
		}
	}

	return false
}

// normalizeEmail normalizes an email address for comparison (lowercase, trimmed).
func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// extractDomain extracts the domain portion from an email address.
// Returns empty string if email is invalid or has no domain.
func extractDomain(email string) string {
	email = strings.TrimSpace(email)
	atIndex := strings.LastIndex(email, "@")
	if atIndex == -1 || atIndex == len(email)-1 {
		return ""
	}
	return strings.ToLower(email[atIndex+1:])
}

// normalizeDomain normalizes a domain for comparison (lowercase, trimmed).
func normalizeDomain(domain string) string {
	return strings.ToLower(strings.TrimSpace(domain))
}
