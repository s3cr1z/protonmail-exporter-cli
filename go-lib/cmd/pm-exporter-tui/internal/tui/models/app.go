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

package models

// Phase 3 - TUI Models (scaffolded)
// TODO: Implement Bubble Tea models for:
// - Authentication state
// - Export configuration state
// - Progress tracking state
// - Completion state

// AppModel represents the root application model
type AppModel struct {
	// TODO: Add Bubble Tea model fields
	CurrentScreen Screen
}

// Screen represents different TUI screens
type Screen int

const (
	ScreenAuth Screen = iota
	ScreenConfig
	ScreenProgress
	ScreenComplete
)

// TODO: Implement Bubble Tea Update, View, and Init methods
