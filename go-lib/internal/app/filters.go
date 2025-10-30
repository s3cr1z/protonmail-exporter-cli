package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/ProtonMail/export-tool/internal/mail"
	"github.com/urfave/cli/v2"
)

// parseFilterOptionsFromCLI parses all filter flags from the CLI context
func parseFilterOptionsFromCLI(ctx *cli.Context) (mail.FilterOptions, error) {
	filters := mail.FilterOptions{}

	// Parse label IDs
	if ctx.IsSet("label") {
		filters.LabelIDs = ctx.StringSlice("label")
	}

	// Parse after date
	if ctx.IsSet("after") {
		after, err := parseTimestamp(ctx.String("after"))
		if err != nil {
			return filters, fmt.Errorf("failed to parse --after: %w", err)
		}
		filters.After = after
	}

	// Parse before date
	if ctx.IsSet("before") {
		before, err := parseTimestamp(ctx.String("before"))
		if err != nil {
			return filters, fmt.Errorf("failed to parse --before: %w", err)
		}
		filters.Before = before
	}

	// Parse from addresses
	if ctx.IsSet("from") {
		filters.From = ctx.StringSlice("from")
	}

	// Parse to addresses
	if ctx.IsSet("to") {
		filters.To = ctx.StringSlice("to")
	}

	// Parse from domains
	if ctx.IsSet("from-domain") {
		filters.FromDomains = ctx.StringSlice("from-domain")
	}

	// Parse to domains
	if ctx.IsSet("to-domain") {
		filters.ToDomains = ctx.StringSlice("to-domain")
	}

	return filters, nil
}

// parseTimestamp parses a timestamp string in various formats and returns Unix epoch seconds.
// Supports:
// - Unix timestamp (integer seconds)
// - RFC3339 format (2006-01-02T15:04:05Z07:00)
// - YYYY-MM-DD format (converted to start of day UTC)
func parseTimestamp(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty timestamp")
	}

	// Try parsing as YYYY-MM-DD first (before Unix timestamp to avoid partial matches)
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.Unix(), nil
	}

	// Try parsing as RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.Unix(), nil
	}

	// Try parsing as Unix timestamp
	var unixTS int64
	n, err := fmt.Sscanf(s, "%d", &unixTS)
	if err == nil && n == 1 && fmt.Sprintf("%d", unixTS) == s {
		return unixTS, nil
	}

	return 0, fmt.Errorf("invalid timestamp format (expected Unix epoch, RFC3339, or YYYY-MM-DD): %s", s)
}
