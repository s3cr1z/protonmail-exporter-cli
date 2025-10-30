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

package main

import (
	"fmt"
	"os"
)

// Phase 3 - TUI entry point (scaffolded)
// TODO: Implement Bubble Tea TUI with:
// - Authentication screen
// - Export configuration screen
// - Progress screen
// - Completion screen

func main() {
	fmt.Println("ProtonMail Exporter TUI")
	fmt.Println("Phase 3 - Not yet implemented")
	fmt.Println("")
	fmt.Println("The TUI will provide an interactive terminal interface with:")
	fmt.Println("  - Secure authentication (no password echo)")
	fmt.Println("  - Interactive export configuration")
	fmt.Println("  - Live progress visualization")
	fmt.Println("  - Full keyboard navigation")
	fmt.Println("")
	fmt.Println("For now, please use the CLI version:")
	fmt.Println("  ./proton-mail-export-cli --help")
	os.Exit(0)
}
