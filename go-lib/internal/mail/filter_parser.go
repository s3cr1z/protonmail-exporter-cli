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
)

// FilterParser provides utilities for parsing filter parameters from strings.
type FilterParser struct{}

// ParseCommaSeparated parses a comma-separated string into a slice of trimmed strings.
// Empty strings are filtered out.
func (FilterParser) ParseCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}

	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// ParseDate parses a date string in various common formats.
// Supported formats: YYYY-MM-DD, YYYY/MM/DD, YYYYMMDD
func (FilterParser) ParseDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}

	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"20060102",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("invalid date format: %s (expected YYYY-MM-DD, YYYY/MM/DD, or YYYYMMDD)", s)
}

// ParseFilterFromStrings creates a Filter from string parameters.
// This is the main entry point for CLI and CGO interfaces.
func ParseFilterFromStrings(
	labelIDs string,
	sender string,
	recipient string,
	domain string,
	after string,
	before string,
	subject string,
) (*Filter, error) {
	parser := FilterParser{}
	filter := NewFilter()

	filter.LabelIDs = parser.ParseCommaSeparated(labelIDs)
	filter.Sender = parser.ParseCommaSeparated(sender)
	filter.Recipient = parser.ParseCommaSeparated(recipient)
	filter.Domain = parser.ParseCommaSeparated(domain)
	filter.Subject = subject

	if after != "" {
		afterTime, err := parser.ParseDate(after)
		if err != nil {
			return nil, fmt.Errorf("invalid after date: %w", err)
		}
		filter.After = afterTime
	}

	if before != "" {
		beforeTime, err := parser.ParseDate(before)
		if err != nil {
			return nil, fmt.Errorf("invalid before date: %w", err)
		}
		filter.Before = beforeTime
	}

	// Validate the filter
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	// Return nil if filter is empty (no filters specified)
	if filter.IsEmpty() {
		return nil, nil
	}

	return filter, nil
}
