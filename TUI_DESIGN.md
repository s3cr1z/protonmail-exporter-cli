# TUI (Terminal User Interface) - Design Document

## Status: Phase 3 - Scaffolded (Not Yet Implemented)

This document outlines the design for an interactive terminal UI using Bubble Tea framework.

## Overview

The TUI will provide a user-friendly, interactive alternative to the CLI with:
- Secure credential input
- Visual progress tracking
- Guided configuration
- Full keyboard navigation

## Technology Stack

### Core Framework
- **Bubble Tea**: TUI framework (The Elm Architecture for Go)
- **Bubbles**: Pre-built components (inputs, spinners, progress bars)
- **Lip Gloss**: Styling and layout

### Dependencies
```go
require (
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/bubbles v0.17.1
    github.com/charmbracelet/lipgloss v0.9.1
)
```

## Architecture

### Directory Structure

```
go-lib/cmd/pm-exporter-tui/
├── main.go                      # Entry point
└── internal/tui/
    ├── models/
    │   ├── app.go              # Root application model
    │   ├── auth.go             # Authentication screen model
    │   ├── config.go           # Configuration screen model
    │   ├── progress.go         # Progress screen model
    │   └── complete.go         # Completion screen model
    ├── views/
    │   ├── auth.go             # Authentication view
    │   ├── config.go           # Configuration view
    │   ├── progress.go         # Progress view
    │   └── complete.go         # Completion view
    └── components/
        ├── filter.go           # Filter configuration component
        ├── progress_bar.go     # Custom progress visualization
        └── help.go             # Help/keybindings component
```

### State Machine

```
┌─────────────┐
│   Initial   │
└──────┬──────┘
       │
       v
┌─────────────┐
│    Auth     │ ← username, password, 2FA
└──────┬──────┘
       │ success
       v
┌─────────────┐
│   Config    │ ← format, filters, output path
└──────┬──────┘
       │ start
       v
┌─────────────┐
│  Progress   │ ← live updates, cancellable
└──────┬──────┘
       │ done
       v
┌─────────────┐
│  Complete   │ ← summary, next steps
└─────────────┘
```

## Screen Designs

### 1. Authentication Screen

```
┌─────────────────────────────────────────────────────────────┐
│                   ProtonMail Exporter                       │
│                                                             │
│  Please enter your ProtonMail credentials                   │
│                                                             │
│  Username: ▓user@proton.me__________________________       │
│                                                             │
│  Password: ▓••••••••_____________________________           │
│                                                             │
│  2FA Code: ▓123456______________________________            │
│            (Press Enter to skip if not enabled)            │
│                                                             │
│  [Tab] Next field  [Enter] Login  [Ctrl+C] Quit           │
└─────────────────────────────────────────────────────────────┘
```

**Features:**
- Password field never shows plaintext
- Support paste (Ctrl+V)
- 2FA field optional
- Validation feedback
- Error messages on failure
- Retry on incorrect credentials

**Security:**
- No password echo
- Clear password from memory after use
- Optional keychain integration (future)

### 2. Export Configuration Screen

```
┌─────────────────────────────────────────────────────────────┐
│                   Export Configuration                       │
│                                                             │
│  Export Format:     [EML] MBOX PDF                         │
│                     ▲▼ to change                            │
│                                                             │
│  Output Path:       ▓/home/user/exports_______________     │
│                                                             │
│  ┌─── Filters ────────────────────────────────────────┐   │
│  │                                                     │   │
│  │  □ Label Filter:    [ Inbox, Sent          ]      │   │
│  │  □ Sender Filter:   [ alice@example.com    ]      │   │
│  │  □ Date Range:      [ 2024-01-01 to        ]      │   │
│  │  □ Subject Keywords:[ invoice              ]      │   │
│  │                                                     │   │
│  │  Combine Mode: ◉ AND  ○ OR                        │   │
│  │                                                     │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌─── Advanced ───────────────────────────────────────┐   │
│  │  PDF Attachments: [List Only] Embed Zip           │   │
│  │  Concurrency:     [10]                             │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  [Enter] Start Export  [Esc] Cancel  [?] Help              │
└─────────────────────────────────────────────────────────────┘
```

**Features:**
- Visual filter configuration
- Toggle filters on/off
- Auto-suggest for labels (fetch from account)
- Path validation and autocomplete
- Preview of estimated message count
- Format-specific options (PDF, MBOX, etc.)

### 3. Progress Screen

```
┌─────────────────────────────────────────────────────────────┐
│                   Export in Progress                         │
│                                                             │
│  Exporting Inbox (Filtered)                                │
│                                                             │
│  Overall Progress                                           │
│  ████████████████████████████░░░░░░░░░░ 72% (1,234/1,700) │
│                                                             │
│  Current Phase: Downloading messages                        │
│  ████████████████████████░░░░░░░░░░░░░░ 58% (98/170)      │
│                                                             │
│  Speed:    15.2 messages/sec                               │
│  ETA:      ~2m 35s remaining                               │
│  Errors:   0                                                │
│                                                             │
│  ┌─── Recent Activity ─────────────────────────────────┐  │
│  │ ✓ Downloaded message ABC123...                      │  │
│  │ ✓ Downloaded message DEF456...                      │  │
│  │ ⚠ Skipped message GHI789 (already exists)          │  │
│  │ ✓ Downloaded message JKL012...                      │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                             │
│  [Ctrl+C] Cancel Export  [L] View Logs                     │
└─────────────────────────────────────────────────────────────┘
```

**Features:**
- Two-level progress (overall + current phase)
- Live updates (smooth animation)
- Throughput and ETA calculation
- Error counter with expandable log
- Activity feed (last N operations)
- Cancellable with confirmation
- Responsive to terminal resize

**Updates:**
- Progress updated on each batch completion
- Activity log scrolls automatically
- Errors highlighted in red
- Warnings in yellow

### 4. Completion Screen

```
┌─────────────────────────────────────────────────────────────┐
│                   Export Complete!                           │
│                                                             │
│  ✓ Successfully exported 1,234 of 1,700 messages           │
│                                                             │
│  Summary:                                                   │
│    • Processed: 1,234 messages                             │
│    • Skipped:   466 (already existed)                      │
│    • Failed:    0                                           │
│    • Duration:  5m 42s                                      │
│    • Output:    /home/user/exports/mail_2024_06_15_120000 │
│                                                             │
│  Next Steps:                                                │
│    • View exported files in output directory               │
│    • Import to email client using MBOX format              │
│    • Backup the export to external storage                 │
│                                                             │
│  [O] Open Folder  [N] New Export  [Q] Quit                 │
└─────────────────────────────────────────────────────────────┘
```

**Features:**
- Success/partial success indication
- Detailed statistics
- Clickable output path (if terminal supports)
- Quick actions (open folder, new export)
- Error log button if failures occurred

## Implementation Details

### Bubble Tea Model Structure

```go
type AppModel struct {
    currentScreen Screen
    auth          AuthModel
    config        ConfigModel
    progress      ProgressModel
    complete      CompleteModel
    
    session       *session.Session
    exportTask    *mail.ExportTask
    
    width, height int
    err           error
}

// Init, Update, View methods implement tea.Model
```

### Message Passing

Use Bubble Tea's message system:

```go
// Authentication messages
type AuthSuccessMsg struct {
    Session *session.Session
}

type AuthFailedMsg struct {
    Error error
}

// Export messages
type ExportStartMsg struct{}

type ExportProgressMsg struct {
    Current, Total uint64
    Phase         string
}

type ExportCompleteMsg struct {
    Success bool
    Stats   ExportStats
}
```

### Component Reusability

Shared components:

```go
// Filter configuration component
type FilterComponent struct {
    filter *mail.ExportFilter
    // Bubble Tea sub-model
}

// Progress bar with custom styling
type ProgressBar struct {
    current, total int
    width         int
    // Bubble Tea sub-model
}
```

## Keyboard Navigation

### Global Keys
- `Ctrl+C`: Quit (with confirmation if export in progress)
- `?`: Show help overlay
- `Esc`: Go back / Cancel

### Screen-Specific Keys

**Authentication:**
- `Tab` / `Shift+Tab`: Navigate fields
- `Enter`: Submit
- `Ctrl+V`: Paste

**Configuration:**
- `↑` / `↓`: Navigate options
- `Space`: Toggle checkboxes
- `Enter`: Activate / Submit
- `Tab`: Next field

**Progress:**
- `L`: Toggle log view
- `Ctrl+C`: Cancel (with confirmation)

**Complete:**
- `O`: Open folder
- `N`: New export
- `Q`: Quit

## Accessibility

### Color Themes

Support multiple themes:
- **Default**: Blue/white/gray
- **High Contrast**: Black/white only
- **Colorblind Safe**: Distinguishable without color

Detect terminal capabilities:
```go
if !term.IsColorSupported() {
    // Use ASCII-only, no colors
}
```

### Screen Reader Compatibility

- Clear labels for all inputs
- Status announcements
- Progress updates in text

### Non-TTY Mode

Detect non-interactive terminal:
```go
if !term.IsTerminal(os.Stdin.Fd()) {
    // Fall back to CLI mode
}
```

## Testing Strategy

### Model Tests

Test state transitions:

```go
func TestAuthModel_SuccessfulLogin(t *testing.T) {
    model := NewAuthModel()
    
    // Simulate input
    model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("user")})
    // ... enter password, 2FA
    model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
    
    // Check auth initiated
    assert.NotNil(t, cmd)
}
```

### View Tests

Snapshot testing for views:

```go
func TestConfigView_Render(t *testing.T) {
    model := NewConfigModel()
    view := model.View()
    
    // Compare with golden file
    golden.AssertMatch(t, "config_view.txt", view)
}
```

### Integration Tests

End-to-end TUI testing with teatest:

```go
func TestTUI_FullExport(t *testing.T) {
    model := NewAppModel()
    
    tm := teatest.New(t, model)
    defer tm.Quit()
    
    // Simulate user interaction
    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
    // ... more interactions
    
    // Assert final state
    finalModel := tm.FinalModel()
    assert.Equal(t, ScreenComplete, finalModel.currentScreen)
}
```

## Logging

### Structured Logs

Write to file, not stdout:

```go
logFile, _ := os.Create("pm-exporter-tui.log")
logger := logrus.New()
logger.SetOutput(logFile)
```

### Debug Mode

Toggle with environment variable:

```bash
PMEXPORTER_DEBUG=1 ./pm-exporter-tui
```

Shows:
- Model state changes
- Network requests
- Performance metrics

## Error Handling

### Graceful Degradation

```go
// Terminal resize errors
if err := term.GetSize(); err != nil {
    // Use default size, continue
}

// Network errors
if err := fetchMessages(); err != nil {
    // Show error, offer retry
    return retryPrompt
}
```

### Panic Recovery

```go
defer func() {
    if r := recover(); r != nil {
        // Log panic
        // Show error screen
        // Safe exit
    }
}()
```

## Performance

### Responsive UI

- Max 60 FPS for animations
- Debounce rapid updates
- Virtual scrolling for long lists

### Memory

- Don't store full message list in memory
- Stream progress updates
- Periodic GC hints for large exports

## Platform Compatibility

### Tested Platforms

- Linux (Ubuntu, Fedora, Arch)
- macOS (10.15+)
- Windows 10+ (with Windows Terminal recommended)

### Terminal Emulators

- iTerm2 (macOS)
- Terminal.app (macOS)
- GNOME Terminal (Linux)
- Konsole (Linux)
- Windows Terminal (Windows)
- Command Prompt (Windows, limited)

## Build and Distribution

### Build Targets

```bash
# Build TUI binary
go build -o pm-exporter-tui ./cmd/pm-exporter-tui

# Cross-compile
GOOS=linux GOARCH=amd64 go build ...
GOOS=darwin GOARCH=arm64 go build ...
GOOS=windows GOARCH=amd64 go build ...
```

### Packaging

- Standalone binary (no dependencies)
- Include in main CLI package
- Separate download for users who prefer TUI

## Migration Path

### Coexistence with CLI

Both CLI and TUI will be available:

```bash
# CLI (existing)
./proton-mail-export-cli --operation backup

# TUI (new)
./pm-exporter-tui
```

### Deprecation Plan

1. **Phase 1**: TUI released as beta
2. **Phase 2**: TUI becomes recommended for interactive use
3. **Phase 3**: CLI remains for scripting and automation
4. No deprecation: both maintained long-term

## Documentation

### User Guide

- Screenshots/GIFs of each screen
- Step-by-step walkthrough
- Keyboard shortcuts reference
- Troubleshooting section

### Developer Guide

- Architecture overview
- Adding new screens
- Custom components
- Testing guidelines

## Acceptance Criteria

✅ Scaffolded (Current):
- [x] Project structure created
- [x] Entry point defined
- [x] Design documented

⏳ TODO (Phase 3 Implementation):
- [ ] Bubble Tea dependencies added
- [ ] Authentication screen implemented
- [ ] Configuration screen implemented
- [ ] Progress screen implemented
- [ ] Completion screen implemented
- [ ] Keyboard navigation working
- [ ] Terminal resize handled
- [ ] Model tests pass
- [ ] Integration tests pass
- [ ] Documentation with screenshots
- [ ] Cross-platform testing complete

## References

- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Lip Gloss Styling](https://github.com/charmbracelet/lipgloss)
- [The Elm Architecture](https://guide.elm-lang.org/architecture/)
