# Proton Mail Export

Proton Mail Export allows you to export your emails in various formats with comprehensive filtering options.

## Features

- **Export Formats**: EML (default), MBOX, PDF (Phase 2 - in development)
- **Comprehensive Filtering**: Filter by labels, senders, recipients, date ranges, and subject keywords
- **Advanced Filter Logic**: Combine filters with AND/OR logic for precise exports
- **Cross-platform Support**: Linux, Mac, Windows
- **Interactive TUI**: Terminal User Interface for guided exports (Phase 3 - in development)
- **Backup and Restore Operations**: Export and re-import your emails

For comprehensive information on using filters, see:
- [ENHANCED_FILTERING_GUIDE.md](ENHANCED_FILTERING_GUIDE.md) - Complete filtering documentation
- [FILTER_EXPORT_USAGE.md](FILTER_EXPORT_USAGE.md) - Basic label filtering guide

For technical design documentation:
- [PDF_EXPORT_DESIGN.md](PDF_EXPORT_DESIGN.md) - PDF export feature design (Phase 2)
- [TUI_DESIGN.md](TUI_DESIGN.md) - Terminal UI design (Phase 3)

# Building

## Requirements

- C++ 17 compatible compiler
  - GCC/Clang (Linux/Mac)
  - MSVC 2022 (Windows)
- CMake >= 3.23
- Go >= 1.24

## Fetch submodules

```
git submodule update --init --recursive
```

## Linux/Mac

```
cmake -S. -B $BUILD_DIR -G <Insert favorite Generator>
cmake --build $BUILD_DIR
```

## Windows

```
cmake -S. -B $BUILD_DIR -G "Visual Studio 17 2022" -DVCPKG_TARGET_TRIPLET=x64-windows-static
cmake --build $BUILD_DIR --config Release
```

**Note:** An active internet connection is required in order to download a standalone MingW compiler in order to compile
the CGO module.

## Layout

- [go-lib](go-lib): CGO Shared library implementation
- [lib](lib): C++ shared library over the exported C interface from [go-lib](go-lib)
- [cli](cli): CLI application
