# CLAUDE.md - AI Assistant Guide for ProtonMail Exporter CLI

> **Last Updated**: 2025-11-16
> **Purpose**: This document helps AI assistants understand the codebase structure, development workflows, and conventions for the ProtonMail Exporter CLI project.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture at a Glance](#architecture-at-a-glance)
3. [Directory Structure](#directory-structure)
4. [Development Environment Setup](#development-environment-setup)
5. [Build System](#build-system)
6. [Code Conventions](#code-conventions)
7. [Testing Guidelines](#testing-guidelines)
8. [Common Development Workflows](#common-development-workflows)
9. [Key Components Reference](#key-components-reference)
10. [Important Constraints & Gotchas](#important-constraints--gotchas)
11. [Making Changes Safely](#making-changes-safely)

---

## Project Overview

### What This Project Does

ProtonMail Exporter CLI is a **cross-platform command-line tool** that allows users to:
- Export emails from ProtonMail accounts as RFC 2822 EML files
- Filter exports by labels, senders, recipients, domains, dates, and subject
- Restore previously exported emails back to ProtonMail
- List available labels/folders in their account

### Technology Stack

- **Go 1.24**: Core business logic (mail processing, API client, encryption)
- **C++17**: CLI frontend and user interface
- **C/CGO**: Interoperability layer between Go and C++
- **CMake 3.23+**: Cross-platform build orchestration
- **vcpkg**: C++ dependency management

### Supported Platforms

- Linux (x86_64, ARM64)
- macOS (x86_64, ARM64, Universal binaries)
- Windows (x86_64 with MinGW for Go compilation)

---

## Architecture at a Glance

### Layer Model

```
User Input
    ↓
[C++ CLI Layer]          cli/bin/main.cpp (cxxopts, interactive prompts)
    ↓
[C++ Wrapper Library]    lib/ (etcpp namespace, exception-based error handling)
    ↓
[CGO C Interface]        go-lib/cmd/lib/ (C function exports, handle-based object passing)
    ↓
[Go Core Library]        go-lib/internal/ (mail processing, API client, session management)
    ↓
[ProtonMail API + Filesystem]
```

### Key Design Principles

1. **Separation of Concerns**: Each layer has distinct responsibilities
2. **Pipeline Architecture**: Export/restore operations use multi-stage async pipelines
3. **Cross-platform Compatibility**: Platform-specific code isolated in utilities
4. **Type Safety**: CGO boundary carefully managed with explicit memory handling
5. **Resumable Operations**: Export/restore can resume from interruptions

---

## Directory Structure

```
protonmail-exporter-cli/
│
├── go-lib/                          # Go shared library (CGO-enabled)
│   ├── cmd/
│   │   ├── lib/                     # CGO C interface exports
│   │   │   ├── export_backup.go     # Backup/export C functions
│   │   │   ├── export_restore.go    # Restore/import C functions
│   │   │   ├── export_globals.go    # Global state management
│   │   │   └── cgo_headers/         # C header files
│   │   └── proton-mail-export/      # Standalone Go CLI (optional)
│   │       └── main.go
│   └── internal/                    # Go implementation packages
│       ├── mail/                    # ⭐ Core email processing logic
│       │   ├── export.go            # ExportTask orchestration
│       │   ├── export_stage_*.go    # Pipeline stages (metadata, download, build, write)
│       │   ├── filter.go            # Filter struct and matching
│       │   ├── filter_parser.go     # CLI filter argument parsing
│       │   ├── restore.go           # RestoreTask orchestration
│       │   └── *_test.go            # Unit tests
│       ├── app/                     # CLI application layer
│       │   ├── app.go               # Main CLI app (urfave/cli)
│       │   └── operation.go         # Operation enum (backup/restore)
│       ├── session/                 # ProtonMail session management
│       ├── apiclient/               # ProtonMail API client wrapper
│       ├── reporter/                # Progress/error reporting
│       ├── telemetry/               # Sentry integration
│       ├── utils/                   # Utilities (JSON versioning, safe I/O)
│       └── hv/                      # Hardware/virtualization detection
│
├── lib/                             # C++ static library wrapper
│   ├── include/                     # Public C++ API headers
│   │   ├── et.hpp                   # Main header (includes all below)
│   │   ├── etsession.hpp            # Session class
│   │   ├── etbackup.hpp             # BackupTask class
│   │   └── etrestore.hpp            # RestoreTask class
│   ├── lib/                         # C++ implementation
│   │   ├── etsession.cpp
│   │   ├── etbackup.cpp
│   │   ├── etrestore.cpp
│   │   └── etutil*.cpp              # Platform-specific utilities
│   └── tests/                       # C++ unit tests (Catch2)
│
├── cli/                             # C++ CLI executable
│   ├── bin/
│   │   ├── main.cpp                 # ⭐ CLI entry point (32KB)
│   │   ├── tasks/                   # Task implementations
│   │   │   ├── backup_task.hpp/cpp
│   │   │   ├── restore_task.hpp/cpp
│   │   │   └── session_task.hpp
│   │   ├── tui_util.cpp/hpp         # Terminal UI utilities
│   │   └── operation.cpp/h          # Operation handling
│   └── cmake/                       # Platform-specific build templates
│
├── cmake/                           # Build system configuration
│   ├── mingw.cmake                  # Windows MinGW compiler setup
│   ├── vcpkg_setup.cmake            # vcpkg integration
│   ├── clang_format.cmake           # Code formatting
│   ├── clang_tidy.cmake             # Static analysis
│   └── compile_options.cmake        # Compiler flags
│
├── ci/                              # GitLab CI configuration
│   ├── setup.yml
│   ├── lint.yml
│   ├── build.yml
│   └── deploy.yml
│
├── .github/workflows/               # GitHub Actions
│   └── codacy.yml                   # Security scanning
│
├── CMakeLists.txt                   # Root build configuration
├── vcpkg.json                       # C++ dependencies
├── go.mod                           # Go dependencies
├── build.sh                         # Quick build script
├── README.md                        # User documentation
├── FILTER_EXPORT_USAGE.md          # Filter feature guide
├── IMPLEMENTATION_SUMMARY.md        # Technical implementation docs
└── SECURITY.md                      # Security policy
```

### Where to Find Things

| What You're Looking For | Where to Look |
|------------------------|---------------|
| Email export logic | `go-lib/internal/mail/export*.go` |
| Email restore logic | `go-lib/internal/mail/restore*.go` |
| Filter implementation | `go-lib/internal/mail/filter*.go` |
| ProtonMail API calls | `go-lib/internal/apiclient/` |
| CLI argument parsing | `cli/bin/main.cpp` + `go-lib/internal/app/app.go` |
| CGO interface | `go-lib/cmd/lib/` |
| C++ wrapper API | `lib/include/` |
| Tests | `go-lib/internal/*/\*_test.go` + `lib/tests/` |
| Build configuration | `CMakeLists.txt` files in each directory |
| CI/CD pipelines | `.gitlab-ci.yml` + `ci/*.yml` |

---

## Development Environment Setup

### Prerequisites

**Required:**
- **CMake** 3.23 or higher
- **Go** 1.24 or higher
- **C++17 compiler**:
  - Linux: GCC 7+ or Clang 5+
  - macOS: Apple Clang (Xcode Command Line Tools)
  - Windows: MSVC 2022
- **Git** (with submodules)

**Recommended:**
- **Ninja** build system (faster than Make)
- **golangci-lint** 1.64.8+ for linting
- **clang-format** for C++ formatting

### Initial Setup

```bash
# Clone with submodules
git clone --recursive <repo-url>
cd protonmail-exporter-cli

# If already cloned, fetch submodules
git submodule update --init --recursive

# Quick build (Unix/macOS with Ninja)
./build.sh

# Or use CMake manually
cmake -S. -B build -G Ninja
cmake --build build
```

### Platform-Specific Setup

**Linux/macOS:**
```bash
cmake -S. -B build -G "Unix Makefiles"
cmake --build build --config Release
```

**Windows:**
```bash
cmake -S. -B build -G "Visual Studio 17 2022" -DVCPKG_TARGET_TRIPLET=x64-windows-static
cmake --build build --config Release
```

**Note**: Windows build requires internet connection to download MinGW compiler for CGO compilation.

---

## Build System

### CMake Structure

The build is hierarchical with three main components:

1. **go-lib** (198-line CMakeLists.txt)
   - Compiles Go code with CGO enabled
   - Creates shared library: `libproton-mail-export.so/.dll/.dylib`
   - Handles cross-compilation for different architectures
   - Most complex build configuration

2. **lib** (59-line CMakeLists.txt)
   - Builds C++ static library `etcpp`
   - Links against Go library and platform libs
   - Includes test subdirectory

3. **cli** (102-line CMakeLists.txt)
   - Builds executable `proton-mail-export-cli`
   - Links against etcpp library
   - Platform-specific: macOS app bundles, Windows resources

### Build Targets

```bash
# Full build
cmake --build build

# Run tests
cmake --build build --target test

# Run linters
cmake --build build --target go-lib-lint
cmake --build build --target clang-format-check

# Clean build
cmake --build build --target clean
```

### Dependencies

**C++ (managed by vcpkg):**
- `catch2` - Unit testing framework
- `fmt` - Fast string formatting
- `cxxopts` - CLI argument parsing

**Go (managed by go.mod):**
- `github.com/ProtonMail/gluon` - Mail library
- `github.com/ProtonMail/go-proton-api` - Official API client
- `github.com/ProtonMail/gopenpgp/v2` - PGP encryption
- `github.com/urfave/cli/v2` - CLI framework
- `github.com/sirupsen/logrus` - Logging
- `github.com/getsentry/sentry-go` - Error reporting
- 83+ transitive dependencies

---

## Code Conventions

### Go Code Style

**Package Organization:**
- One primary type/concept per file
- Internal packages under `go-lib/internal/`
- Tests in `*_test.go` files alongside implementation

**Naming:**
- **Packages**: lowercase, single word (e.g., `mail`, `session`, `apiclient`)
- **Types**: PascalCase (e.g., `ExportTask`, `Filter`, `Session`)
- **Functions**: camelCase for private, PascalCase for exported
- **Constants**: PascalCase (e.g., `NumParallelDownloads`, `MetadataPageSize`)
- **Interfaces**: Usually named `{Concept}` or `{Concept}Interface`

**Error Handling:**
- Always return errors explicitly
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Log at appropriate levels (info, warn, error)

**Comments:**
- Document all exported types and functions
- Use `//` for line comments
- Use `/* */` for block comments
- License header in all files (GNU GPL v3)

**Example:**
```go
// ExportTask manages the email export operation.
type ExportTask struct {
    client   *apiclient.Client
    filter   *Filter
    reporter Reporter
}

// Run executes the export pipeline with the given context.
func (t *ExportTask) Run(ctx context.Context) error {
    if err := t.validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    // ...
}
```

### C++ Code Style

**Namespace:**
- All code in `etcpp` namespace (export tool C++)

**Naming:**
- **Classes**: PascalCase (e.g., `Session`, `BackupTask`)
- **Methods**: camelCase (e.g., `newBackup`, `setCallbacks`)
- **Private members**: `m` prefix (e.g., `mPtr`, `mCallbacks`, `mUsername`)
- **Exceptions**: `*Exception` suffix (e.g., `SessionException`, `CancelledException`)

**File Organization:**
- Header/implementation split: `.hpp`/`.cpp`
- Public headers in `lib/include/`
- Implementation in `lib/lib/`
- One class per file (generally)

**Error Handling:**
- Use exceptions for errors
- Throw specific exception types
- Catch at appropriate boundaries

**Example:**
```cpp
namespace etcpp {

class Session {
public:
    Session(const std::string& username, const std::string& password);
    BackupTask newBackup(const std::string& dir, const std::vector<std::string>& filters);

private:
    void* mPtr;  // Go handle (cgo.Handle)
    std::string mUsername;
};

} // namespace etcpp
```

### CGO Interface Conventions

**Memory Management:**
- All Go-allocated strings returned to C must be freed with `etFree()`
- Use `cgo.Handle` for passing Go objects to C (not raw pointers)
- Explicit cleanup functions for all resources

**Naming:**
- C functions: snake_case with `et` prefix (e.g., `etSessionNew`, `etSessionLogin`)
- Always use `extern "C"` linkage

**Example:**
```go
//export etSessionNew
func etSessionNew(username, password *C.char) C.uintptr_t {
    user := C.GoString(username)
    pass := C.GoString(password)

    session := newSession(user, pass)
    handle := cgo.NewHandle(session)
    return C.uintptr_t(handle)
}

//export etFree
func etFree(ptr unsafe.Pointer) {
    C.free(ptr)
}
```

### File Headers

All source files must include the GNU GPL v3 license header:

```
// Copyright (C) 2024 Proton AG
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
```

---

## Testing Guidelines

### Go Testing

**Framework**: Standard `testing` package + `testify/assert`

**Test File Organization:**
- Tests in `*_test.go` files
- Parallel structure to implementation
- Test package: `package mail_test` for black-box, `package mail` for white-box

**Running Tests:**
```bash
# All tests
go test ./...

# Specific package
go test ./go-lib/internal/mail

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

**Test Naming:**
- Functions: `Test{Function}_{Scenario}` (e.g., `TestFilter_MatchesSender`)
- Tables: Use `t.Run()` for subtests

**Example:**
```go
func TestFilter_MatchesSender(t *testing.T) {
    tests := []struct {
        name     string
        filter   Filter
        message  *proton.Message
        expected bool
    }{
        {
            name:     "exact match",
            filter:   Filter{Sender: []string{"test@example.com"}},
            message:  &proton.Message{Sender: &mail.Address{Address: "test@example.com"}},
            expected: true,
        },
        // ...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := tt.filter.Matches(tt.message)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### C++ Testing

**Framework**: Catch2

**Running Tests:**
```bash
# Build and run tests
cmake --build build --target test

# Or directly
./build/lib/tests/etcpp-tests
```

**Test File Organization:**
- Tests in `lib/tests/test_*.cpp`
- Test data in `lib/tests/test_data/`

**Example:**
```cpp
#include <catch2/catch.hpp>
#include "etsession.hpp"

TEST_CASE("Session creation", "[session]") {
    SECTION("valid credentials") {
        etcpp::Session session("user@proton.me", "password");
        REQUIRE_NOTHROW(session.login());
    }

    SECTION("invalid credentials") {
        etcpp::Session session("invalid", "wrong");
        REQUIRE_THROWS_AS(session.login(), etcpp::SessionException);
    }
}
```

### Test Coverage

**Current Coverage:**
- `go-lib/internal/mail`: 16 test files covering:
  - Filter validation and matching
  - Filter parsing (CLI arguments)
  - Export pipeline stages
  - Restore operations
  - Integration tests

**Key Test Files:**
- `filter_test.go` - Filter logic unit tests
- `filter_parser_test.go` - CLI parsing tests
- `filter_integration_test.go` - End-to-end filtering
- `export_stage_*_test.go` - Pipeline stage tests
- `export_test.go` - Full export workflow tests

---

## Common Development Workflows

### Adding a New CLI Flag

**Steps:**

1. **Add Go flag definition** in `go-lib/internal/app/app.go`:
```go
&cli.StringFlag{
    Name:    "my-flag",
    Usage:   "Description of flag",
    EnvVars: []string{"ET_MY_FLAG"},
},
```

2. **Parse flag value** in app logic:
```go
myValue := c.String("my-flag")
```

3. **Update CGO interface** in `go-lib/cmd/lib/export_backup.go`:
```go
//export etSessionNewBackup
func etSessionNewBackup(..., myValue *C.char) C.uintptr_t {
    val := C.GoString(myValue)
    // ...
}
```

4. **Update C++ wrapper** in `lib/include/etsession.hpp` and `lib/lib/etsession.cpp`:
```cpp
BackupTask newBackup(..., const std::string& myValue);
```

5. **Update CLI** in `cli/bin/main.cpp`:
```cpp
result.add_options()
    ("my-flag", "Description", cxxopt::value<std::string>());
```

6. **Update documentation** in `README.md` and usage examples.

### Adding a New Filter Type

**Example: Adding a "Has Attachment" filter**

1. **Update Filter struct** in `go-lib/internal/mail/filter.go`:
```go
type Filter struct {
    // ... existing fields
    HasAttachment *bool  // nil = don't filter, true/false = filter
}

func (f *Filter) Matches(msg *proton.Message) bool {
    // ... existing logic

    if f.HasAttachment != nil {
        hasAttach := len(msg.Attachments) > 0
        if hasAttach != *f.HasAttachment {
            return false
        }
    }

    return true
}
```

2. **Add parser logic** in `go-lib/internal/mail/filter_parser.go`:
```go
func ParseFilter(ctx *cli.Context) (*Filter, error) {
    // ... existing parsing

    var hasAttach *bool
    if ctx.IsSet("has-attachment") {
        val := ctx.Bool("has-attachment")
        hasAttach = &val
    }

    return &Filter{
        // ... existing fields
        HasAttachment: hasAttach,
    }, nil
}
```

3. **Add tests** in `go-lib/internal/mail/filter_test.go`:
```go
func TestFilter_MatchesAttachment(t *testing.T) {
    trueVal := true
    falseVal := false

    tests := []struct {
        name     string
        filter   Filter
        message  *proton.Message
        expected bool
    }{
        {
            name:     "has attachment - matches",
            filter:   Filter{HasAttachment: &trueVal},
            message:  &proton.Message{Attachments: []*proton.Attachment{{}}},
            expected: true,
        },
        // ... more cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.expected, tt.filter.Matches(tt.message))
        })
    }
}
```

4. **Follow "Adding a New CLI Flag" workflow** to expose to CLI.

### Debugging CGO Issues

**Common Issues:**

1. **Memory leaks**:
   - Check all `C.CString()` calls are paired with `C.free()`
   - Ensure Go strings returned to C are freed with `etFree()`

2. **Crashes at CGO boundary**:
   - Verify handle validity before dereferencing
   - Check for nil pointers before passing to C
   - Use `defer recover()` in exported functions

3. **String encoding issues**:
   - Always use `C.GoString()` to convert `*C.char` to Go string
   - Use `C.CString()` to convert Go string to `*C.char`

**Debugging Tools:**
```bash
# Enable CGO debug output
export GODEBUG=cgocheck=1

# Run with race detector
go test -race ./...

# Use valgrind (Linux)
valgrind --leak-check=full ./build/proton-mail-export-cli

# Use AddressSanitizer
CGO_CFLAGS="-fsanitize=address" go build
```

### Running Linters

```bash
# Go linting (golangci-lint)
cmake --build build --target go-lib-lint

# Or directly
golangci-lint run ./go-lib/...

# C++ formatting check
cmake --build build --target clang-format-check

# Auto-format C++ code
clang-format -i lib/**/*.cpp lib/**/*.hpp cli/**/*.cpp cli/**/*.hpp
```

---

## Key Components Reference

### Export Pipeline

**Location**: `go-lib/internal/mail/export*.go`

**Flow**:
1. **MetadataStage** - Fetch message metadata from API
   - Server-side filtering (single label + subject)
   - Client-side filtering (complex queries)
   - Pagination with resume support

2. **DownloadStage** - Download and decrypt messages
   - Parallel downloads (10 workers)
   - Decryption with user's private key

3. **BuildStage** - Construct RFC 2822 EML format
   - Parallel building (4 workers)
   - Headers, body, attachments

4. **WriteStage** - Write to filesystem
   - Parallel writing (4 workers)
   - Creates `.eml` + `.meta.json` files
   - Resume capability

**Configuration Constants** (in `export.go`):
```go
const (
    NumParallelDownloads = 10
    NumParallelBuilds    = 4
    NumParallelWrites    = 4
    MetadataPageSize     = 150
)
```

### Filter System

**Location**: `go-lib/internal/mail/filter*.go`

**Filter Struct**:
```go
type Filter struct {
    LabelIDs  []string    // Folder/label IDs (OR logic)
    Sender    []string    // Email addresses or domains (@domain.com)
    Recipient []string    // To/CC/BCC addresses
    Domain    []string    // Domain matching (sender OR recipient)
    After     *time.Time  // Date range start
    Before    *time.Time  // Date range end
    Subject   string      // Case-insensitive substring
}
```

**Server-side vs Client-side**:
- **Server-side**: Single label + optional subject (most efficient)
- **Client-side**: All other combinations (streaming filter during pagination)

**Supported Date Formats**:
- `YYYY-MM-DD` (ISO 8601)
- `YYYY/MM/DD`
- `YYYYMMDD`

### Session Management

**Location**: `go-lib/internal/session/session.go`

**Features**:
- Login with username/password
- TOTP 2FA support
- Mailbox password handling
- Network loss detection & recovery
- Token refresh
- User information caching

**Lifecycle**:
```go
session := session.New(username, password)
session.Login(totpCode)
session.SetMailboxPassword(mboxPass)
// ... use session
session.Logout()
```

### Progress Reporting

**Location**: `go-lib/internal/reporter/reporter.go`

**Interface**:
```go
type Reporter interface {
    OnProgress(stage string, current, total int64)
    OnError(stage string, err error)
    OnComplete(stage string)
}
```

**Stages**:
- `metadata` - Fetching message list
- `download` - Downloading messages
- `build` - Building EML files
- `write` - Writing to disk

---

## Important Constraints & Gotchas

### 1. CGO Memory Management

**CRITICAL**: Every Go string returned to C must be explicitly freed.

**Correct Pattern**:
```go
//export etGetString
func etGetString() *C.char {
    return C.CString("value")  // Caller MUST call etFree()
}
```

**C++ Side**:
```cpp
char* str = etGetString();
std::string value(str);
etFree(str);  // MUST FREE
```

### 2. Platform-Specific Code

**Windows Considerations**:
- MinGW compiler required for CGO (auto-downloaded by CMake)
- Path separators: Use `filepath.Join()` in Go
- Line endings: Git should normalize (`core.autocrlf=true`)

**macOS Considerations**:
- Universal binaries combine x86_64 + ARM64
- App bundle creation in `cli/cmake/`
- Foundation framework required for native APIs

**Linux Considerations**:
- GLIBC version compatibility
- Static linking may be needed for portability

### 3. API Rate Limiting

ProtonMail API has rate limits. The client implements auto-retry with exponential backoff:

**Location**: `go-lib/internal/apiclient/autoretryclient.go`

**Configuration**:
- Initial delay: 1 second
- Max delay: 32 seconds
- Max retries: 5

### 4. Large Mailbox Handling

**Memory Considerations**:
- Metadata is paginated (150 messages per page)
- Filters applied during pagination (streaming)
- Download/build/write stages use worker pools to limit concurrency

**Best Practices**:
- Don't load all messages into memory at once
- Use channels to stream between stages
- Monitor goroutine count

### 5. Test Data Isolation

**DO NOT**:
- Use real ProtonMail credentials in tests
- Commit test data to repository (use `test_data/` with `.gitignore`)
- Make real API calls in unit tests (use mocks)

**DO**:
- Use mock clients (`apiclient/mocks.go`)
- Generate test messages programmatically
- Use table-driven tests for comprehensive coverage

### 6. Versioning and Compatibility

**Export Format Versioning**:
- `labels.json` and metadata files include version numbers
- Use `internal/utils/versioned_json.go` for reading/writing
- Maintain backward compatibility when changing format

**API Version**:
- ProtonMail API versioning handled by `go-proton-api`
- Monitor for breaking changes in dependencies

---

## Making Changes Safely

### Pre-commit Checklist

Before committing changes:

- [ ] Run tests: `go test ./... && cmake --build build --target test`
- [ ] Run linters: `golangci-lint run ./go-lib/...`
- [ ] Format code: `clang-format -i` for C++
- [ ] Update documentation if adding features
- [ ] Check CGO memory management (no leaks)
- [ ] Verify cross-platform compatibility (if touching platform code)
- [ ] Update `IMPLEMENTATION_SUMMARY.md` for significant features

### Branching Strategy

**Main Branches**:
- `master` - Stable release branch
- `dev-windsurf` - Active development branch

**Feature Branches**:
- Create from `dev-windsurf`
- Name pattern: `feature/{description}` or `fix/{issue-number}`
- Merge back via pull request

**AI Assistant Branches**:
- Prefix with `claude/` (e.g., `claude/add-attachment-filter`)
- Push to remote when complete
- Create PR for human review

### Commit Message Format

Follow conventional commits:

```
type(scope): brief description

Longer explanation if needed.

Fixes #123
```

**Types**:
- `feat` - New feature
- `fix` - Bug fix
- `refactor` - Code restructuring
- `test` - Adding tests
- `docs` - Documentation changes
- `chore` - Maintenance tasks

**Examples**:
```
feat(mail): add attachment filtering support

Implements has-attachment filter to allow users to filter
emails based on attachment presence.

Fixes #42
```

```
fix(cgo): prevent memory leak in etSessionNewBackup

Added missing etFree() call for filter string parameters.
```

### CI/CD Pipeline

**GitLab CI Stages** (`.gitlab-ci.yml`):
1. `analyse` - Static analysis
2. `lint` - Code linting (Go + C++)
3. `build` - Multi-platform builds
4. `installer` - Package creation
5. `deploy` - Release distribution

**GitHub Actions**:
- `codacy.yml` - Security scanning (weekly + on push)

**All CI must pass** before merging to `dev-windsurf` or `master`.

### Debugging Tips

**Go Debugging**:
```bash
# Print debug logs
export ET_LOG_LEVEL=debug
./proton-mail-export-cli --operation backup

# Use delve debugger
dlv debug ./go-lib/cmd/proton-mail-export
```

**C++ Debugging**:
```bash
# GDB
gdb ./build/proton-mail-export-cli
(gdb) run --operation backup

# LLDB (macOS)
lldb ./build/proton-mail-export-cli
(lldb) run --operation backup
```

**CGO Debugging**:
```bash
# Enable CGO checks
export GODEBUG=cgocheck=2

# Verbose CGO output
export CGO_CFLAGS="-g -O0"
export CGO_LDFLAGS="-g"
```

---

## Additional Resources

### Documentation Files

- `README.md` - User-facing documentation
- `FILTER_EXPORT_USAGE.md` - Detailed filter usage guide
- `IMPLEMENTATION_SUMMARY.md` - Technical implementation details
- `IMPLEMENTATION_PHASE1_SUMMARY.md` - Phase 1 feature summary
- `SECURITY.md` - Security policy and vulnerability reporting

### External Documentation

- [ProtonMail API Docs](https://github.com/ProtonMail/go-proton-api)
- [CMake Documentation](https://cmake.org/documentation/)
- [CGO Documentation](https://pkg.go.dev/cmd/cgo)
- [vcpkg Documentation](https://vcpkg.io/en/getting-started.html)

### Getting Help

- **Issues**: Check existing issues in the repository
- **CI Logs**: Review GitLab CI pipeline logs for build failures
- **Code Search**: Use `grep -r "pattern" go-lib/` to find examples

---

## Summary for AI Assistants

When working on this codebase:

1. **Understand the layer you're modifying**: Go core, CGO interface, C++ wrapper, or CLI
2. **Follow the data flow**: User → CLI → C++ → CGO → Go → API
3. **Test thoroughly**: Both Go and C++ have comprehensive test suites
4. **Mind the CGO boundary**: Memory management is critical
5. **Maintain cross-platform compatibility**: Test on Linux, macOS, and Windows if possible
6. **Document changes**: Update relevant MD files and code comments
7. **Use the pipeline pattern**: Export/restore use multi-stage async processing
8. **Respect existing conventions**: Naming, file organization, error handling
9. **Check CI before merging**: All linters and tests must pass

**Most Important**: This is a tool that handles user data (emails). Security, correctness, and data integrity are paramount. When in doubt, ask for clarification rather than making assumptions.

---

**End of CLAUDE.md**
