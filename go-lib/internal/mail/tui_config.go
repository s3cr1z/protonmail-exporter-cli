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

// TUIConfig represents configuration for the Terminal User Interface (TUI).
// This is a placeholder interface for future TUI implementation.
//
// The TUI will provide an interactive interface for:
// - Selecting filters interactively
// - Monitoring export progress with rich visual feedback
// - Managing multiple export/restore operations
// - Viewing export history and statistics
type TUIConfig struct {
	// Filter is the shared filter that will be used for the export
	Filter *Filter

	// ExportPath is the destination directory for exports
	ExportPath string

	// Theme configures the visual appearance of the TUI
	Theme TUITheme

	// EnableMouseSupport determines whether mouse interactions are enabled
	EnableMouseSupport bool

	// RefreshInterval is the interval for updating progress displays (in milliseconds)
	RefreshInterval int
}

// TUITheme defines color and style preferences for the TUI.
type TUITheme struct {
	// PrimaryColor is the main accent color
	PrimaryColor string

	// BackgroundColor is the background color
	BackgroundColor string

	// TextColor is the default text color
	TextColor string

	// ProgressBarColor is the color for progress indicators
	ProgressBarColor string
}

// DefaultTUIConfig returns a default TUI configuration.
func DefaultTUIConfig() *TUIConfig {
	return &TUIConfig{
		Filter:              NewFilter(),
		ExportPath:          "",
		EnableMouseSupport:  true,
		RefreshInterval:     100,
		Theme: TUITheme{
			PrimaryColor:     "#6D4AFF",
			BackgroundColor:  "#1A1A1A",
			TextColor:        "#FFFFFF",
			ProgressBarColor: "#6D4AFF",
		},
	}
}

// TUIScreen represents different screens in the TUI application.
// This is a placeholder for the TUI state machine.
type TUIScreen int

const (
	// TUIScreenMain is the main menu screen
	TUIScreenMain TUIScreen = iota

	// TUIScreenFilterSetup allows users to configure export filters
	TUIScreenFilterSetup

	// TUIScreenExportProgress shows the export operation progress
	TUIScreenExportProgress

	// TUIScreenHistory shows previous export operations
	TUIScreenHistory

	// TUIScreenSettings allows configuration changes
	TUIScreenSettings
)

// TUIController is a placeholder interface for the TUI controller.
// This will be implemented when TUI functionality is added.
type TUIController interface {
	// Run starts the TUI application
	Run() error

	// GetCurrentScreen returns the currently active screen
	GetCurrentScreen() TUIScreen

	// NavigateTo switches to a different screen
	NavigateTo(screen TUIScreen) error

	// GetFilter returns the currently configured filter
	GetFilter() *Filter

	// SetFilter updates the current filter configuration
	SetFilter(filter *Filter) error

	// StartExport initiates an export operation with the current configuration
	StartExport() error

	// CancelExport cancels the currently running export
	CancelExport() error

	// Shutdown gracefully shuts down the TUI
	Shutdown() error
}
