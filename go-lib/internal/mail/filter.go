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
	"fmt"
	"strings"
	"time"

	"github.com/ProtonMail/go-proton-api"
)

// Filter represents all possible export filtering options.
// This struct is designed to be reusable by both CLI and future TUI implementations.
type Filter struct {
	// LabelIDs filters messages by label/folder IDs (OR logic)
	LabelIDs []string

	// Sender filters messages from specific email addresses or domains
	// Supports exact match (user@domain.com) or domain match (@domain.com)
	Sender []string

	// Recipient filters messages to specific email addresses or domains
	// Applies to To, CC, and BCC fields
	// Supports exact match (user@domain.com) or domain match (@domain.com)
	Recipient []string

	// Domain filters messages by sender or recipient domain
	// More convenient than specifying @domain.com for sender/recipient
	Domain []string

	// After filters messages sent after this date (inclusive)
	After *time.Time

	// Before filters messages sent before this date (inclusive)
	Before *time.Time

	// Subject filters messages by subject (substring match, case-insensitive)
	Subject string
}

// NewFilter creates a new empty filter.
func NewFilter() *Filter {
	return &Filter{
		LabelIDs:  make([]string, 0),
		Sender:    make([]string, 0),
		Recipient: make([]string, 0),
		Domain:    make([]string, 0),
	}
}

// IsEmpty returns true if no filters are set.
func (f *Filter) IsEmpty() bool {
	return len(f.LabelIDs) == 0 &&
		len(f.Sender) == 0 &&
		len(f.Recipient) == 0 &&
		len(f.Domain) == 0 &&
		f.After == nil &&
		f.Before == nil &&
		f.Subject == ""
}

// Validate checks if the filter configuration is valid.
func (f *Filter) Validate() error {
	if f.After != nil && f.Before != nil && f.After.After(*f.Before) {
		return fmt.Errorf("after date must be before or equal to before date")
	}

	// Validate email formats for sender/recipient
	for _, sender := range f.Sender {
		if err := validateEmailOrDomain(sender); err != nil {
			return fmt.Errorf("invalid sender format %q: %w", sender, err)
		}
	}

	for _, recipient := range f.Recipient {
		if err := validateEmailOrDomain(recipient); err != nil {
			return fmt.Errorf("invalid recipient format %q: %w", recipient, err)
		}
	}

	for _, domain := range f.Domain {
		if err := validateDomain(domain); err != nil {
			return fmt.Errorf("invalid domain format %q: %w", domain, err)
		}
	}

	return nil
}

// ToServerFilter converts this filter to a proton.MessageFilter for server-side filtering.
// Returns nil if no server-side filters can be applied.
// Note: Only LabelID and Subject are supported server-side.
func (f *Filter) ToServerFilter() *proton.MessageFilter {
	filter := &proton.MessageFilter{
		Desc: true, // Always fetch in descending order
	}

	hasServerFilter := false

	// Server-side label filtering (single label only)
	if len(f.LabelIDs) == 1 {
		filter.LabelID = f.LabelIDs[0]
		hasServerFilter = true
	}

	// Server-side subject filtering
	if f.Subject != "" {
		filter.Subject = f.Subject
		hasServerFilter = true
	}

	if !hasServerFilter {
		return nil
	}

	return filter
}

// NeedsClientFiltering returns true if any client-side filters need to be applied.
func (f *Filter) NeedsClientFiltering() bool {
	return len(f.LabelIDs) > 1 || // Multiple labels require client-side OR
		len(f.Sender) > 0 ||
		len(f.Recipient) > 0 ||
		len(f.Domain) > 0 ||
		f.After != nil ||
		f.Before != nil ||
		(f.Subject != "" && len(f.LabelIDs) > 0) // Subject + labels requires client-side
}

// MatchesMetadata checks if a message metadata matches all filter criteria.
// This is used for client-side filtering.
func (f *Filter) MatchesMetadata(metadata proton.MessageMetadata) bool {
	// Check label filter
	if len(f.LabelIDs) > 0 {
		if !f.matchesLabel(metadata) {
			return false
		}
	}

	// Check sender filter
	if len(f.Sender) > 0 {
		if !f.matchesSender(metadata) {
			return false
		}
	}

	// Check recipient filter
	if len(f.Recipient) > 0 {
		if !f.matchesRecipient(metadata) {
			return false
		}
	}

	// Check domain filter (applies to both sender and recipient)
	if len(f.Domain) > 0 {
		if !f.matchesDomain(metadata) {
			return false
		}
	}

	// Check date filters
	if f.After != nil || f.Before != nil {
		if !f.matchesDate(metadata) {
			return false
		}
	}

	// Check subject filter (case-insensitive substring match)
	if f.Subject != "" {
		if !f.matchesSubject(metadata) {
			return false
		}
	}

	return true
}

func (f *Filter) matchesLabel(metadata proton.MessageMetadata) bool {
	for _, requestedLabel := range f.LabelIDs {
		for _, msgLabel := range metadata.LabelIDs {
			if msgLabel == requestedLabel {
				return true
			}
		}
	}
	return false
}

func (f *Filter) matchesSender(metadata proton.MessageMetadata) bool {
	if metadata.Sender == nil {
		return false
	}

	senderEmail := strings.ToLower(metadata.Sender.Address)

	// Check explicit sender filters
	for _, sender := range f.Sender {
		if matchesEmailPattern(senderEmail, strings.ToLower(sender)) {
			return true
		}
	}

	return false
}

func (f *Filter) matchesRecipient(metadata proton.MessageMetadata) bool {
	// Collect all recipient emails
	recipients := make([]string, 0)
	for _, addr := range metadata.ToList {
		if addr != nil {
			recipients = append(recipients, strings.ToLower(addr.Address))
		}
	}
	for _, addr := range metadata.CCList {
		if addr != nil {
			recipients = append(recipients, strings.ToLower(addr.Address))
		}
	}
	for _, addr := range metadata.BCCList {
		if addr != nil {
			recipients = append(recipients, strings.ToLower(addr.Address))
		}
	}

	if len(recipients) == 0 {
		return false
	}

	// Check if any recipient matches filters
	for _, recipEmail := range recipients {
		// Check explicit recipient filters
		for _, recip := range f.Recipient {
			if matchesEmailPattern(recipEmail, strings.ToLower(recip)) {
				return true
			}
		}
	}

	return false
}

func (f *Filter) matchesDomain(metadata proton.MessageMetadata) bool {
	// Check sender domain
	if metadata.Sender != nil {
		senderEmail := strings.ToLower(metadata.Sender.Address)
		for _, domain := range f.Domain {
			if matchesDomain(senderEmail, strings.ToLower(domain)) {
				return true
			}
		}
	}

	// Check recipient domains
	recipients := make([]string, 0)
	for _, addr := range metadata.ToList {
		if addr != nil {
			recipients = append(recipients, strings.ToLower(addr.Address))
		}
	}
	for _, addr := range metadata.CCList {
		if addr != nil {
			recipients = append(recipients, strings.ToLower(addr.Address))
		}
	}
	for _, addr := range metadata.BCCList {
		if addr != nil {
			recipients = append(recipients, strings.ToLower(addr.Address))
		}
	}

	for _, recipEmail := range recipients {
		for _, domain := range f.Domain {
			if matchesDomain(recipEmail, strings.ToLower(domain)) {
				return true
			}
		}
	}

	return false
}

func (f *Filter) matchesDate(metadata proton.MessageMetadata) bool {
	msgTime := time.Unix(metadata.Time, 0)

	if f.After != nil && msgTime.Before(*f.After) {
		return false
	}

	if f.Before != nil && msgTime.After(*f.Before) {
		return false
	}

	return true
}

func (f *Filter) matchesSubject(metadata proton.MessageMetadata) bool {
	return strings.Contains(
		strings.ToLower(metadata.Subject),
		strings.ToLower(f.Subject),
	)
}

// Helper functions

func validateEmailOrDomain(s string) error {
	if s == "" {
		return fmt.Errorf("empty value")
	}

	// Domain pattern: @domain.com
	if strings.HasPrefix(s, "@") {
		return validateDomain(strings.TrimPrefix(s, "@"))
	}

	// Email pattern: must contain @
	if !strings.Contains(s, "@") {
		return fmt.Errorf("must be an email address or domain (starting with @)")
	}

	parts := strings.Split(s, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func validateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("empty domain")
	}

	if strings.Contains(domain, "@") {
		return fmt.Errorf("domain should not contain @")
	}

	if !strings.Contains(domain, ".") {
		return fmt.Errorf("invalid domain format")
	}

	return nil
}

func matchesEmailPattern(email, pattern string) bool {
	// Domain pattern: @domain.com
	if strings.HasPrefix(pattern, "@") {
		return matchesDomain(email, strings.TrimPrefix(pattern, "@"))
	}

	// Exact email match
	return email == pattern
}

func matchesDomain(email, domain string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	return parts[1] == domain
}
