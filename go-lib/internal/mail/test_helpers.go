package mail

import "time"

// Test helper functions shared across test files

// timePtr is a helper function for creating time pointers in tests.
func timePtr(t time.Time) *time.Time {
	return &t
}
