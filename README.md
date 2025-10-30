# Proton Mail Export

Proton Mail Export allows you to export your emails as eml files with powerful filtering options.

## Features

- **Advanced Filtering**: Filter exports by labels, senders, recipients, domains, dates, and subject
- **Server-side and Client-side Filtering**: Leverages ProtonMail API for efficient server-side filtering where possible
- **Flexible Export Options**: Export all emails or specific subsets based on your criteria
- **Support for backup and restore operations**
- **Cross-platform support** (Linux, Mac, Windows)

## Quick Start

### Basic Export
```bash
./proton-mail-export-cli --operation backup --dir ./export
```

### Filtered Export Examples

Export only Inbox:
```bash
./proton-mail-export-cli --operation backup --label 0
```

Export emails from a specific sender:
```bash
./proton-mail-export-cli --operation backup --from user@example.com
```

Export emails from a domain:
```bash
./proton-mail-export-cli --operation backup --domain example.com
```

Export emails within a date range:
```bash
./proton-mail-export-cli --operation backup --after 2024-01-01 --before 2024-12-31
```

Combine multiple filters:
```bash
./proton-mail-export-cli --operation backup \
  --label 0,2 \
  --from @work-domain.com \
  --after 2024-01-01 \
  --subject "project"
```

## Filter Options

| Option | Description | Environment Variable | Example |
|--------|-------------|---------------------|---------|
| `--label` | Filter by folder/label IDs (comma-separated) | `ET_FILTER_LABELS` | `--label 0,2,10` |
| `--from` | Filter by sender email/domain (comma-separated) | `ET_FILTER_FROM` | `--from user@example.com,@domain.com` |
| `--to` | Filter by recipient email/domain (comma-separated) | `ET_FILTER_TO` | `--to user@example.com` |
| `--domain` | Filter by domain in sender or recipient | `ET_FILTER_DOMAIN` | `--domain example.com` |
| `--after` | Filter messages after date (YYYY-MM-DD) | `ET_FILTER_AFTER` | `--after 2024-01-01` |
| `--before` | Filter messages before date (YYYY-MM-DD) | `ET_FILTER_BEFORE` | `--before 2024-12-31` |
| `--subject` | Filter by subject substring (case-insensitive) | `ET_FILTER_SUBJECT` | `--subject "important"` |
| `--list-labels` | List available folder/label IDs | - | `--list-labels` |

### Common Label IDs
- **Inbox**: `0`
- **Drafts**: `1`
- **Sent**: `2`
- **Trash**: `3`
- **Spam**: `4`
- **All Mail**: `5`
- **Archive**: `6`
- **Starred**: `10`

Custom folders have unique IDs - use `--list-labels` to find them.

## Performance Notes

- **Server-side filtering** is used automatically for single-label and subject filters
- **Client-side filtering** is used for complex filters (multiple labels, sender/recipient, dates, domains)
- Filtering significantly reduces export time and disk space for targeted exports
- All filtering options can be combined for precise email selection

## Advanced Usage

For more detailed information on filtering, see [FILTER_EXPORT_USAGE.md](FILTER_EXPORT_USAGE.md).

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
