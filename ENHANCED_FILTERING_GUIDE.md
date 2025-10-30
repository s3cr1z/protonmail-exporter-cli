# Enhanced Filtering Guide

## Overview

The ProtonMail Exporter supports comprehensive filtering options to export specific subsets of your email based on multiple criteria including labels, senders, recipients, dates, and subject keywords.

## Filter Capabilities

### Server-Side vs Client-Side Filtering

The exporter intelligently uses:
- **Server-side filtering** when supported by the ProtonMail API (single label filtering)
- **Client-side filtering** for advanced criteria (multiple labels, senders, recipients, dates, subject keywords)

Client-side filtering processes messages incrementally during download, minimizing memory usage even for large mailboxes.

### Supported Filter Types

#### 1. Label/Folder Filtering

Filter emails by one or more labels or folders. Multiple labels use OR logic (messages matching ANY label are included).

**Common System Labels:**
- Inbox: `0`
- Drafts: `1`
- Sent: `2`
- Trash: `3`
- Spam: `4`
- All Mail: `5`
- Archive: `6`
- Starred: `10`

**Examples:**
```bash
# Export Inbox only
--filter-labels "0"

# Export Inbox and Sent
--filter-labels "0,2"

# Export custom folder (get ID from --list-labels)
--filter-labels "your-custom-label-id"
```

#### 2. Sender Filtering

Filter by sender email address or name. Matching is case-insensitive and supports partial matches.

**Examples:**
```bash
# Export emails from specific address
--filter-sender "alice@example.com"

# Export from multiple senders (OR logic)
--filter-sender "alice@example.com,bob@example.com"

# Export from domain
--filter-sender "@company.com"

# Export by sender name
--filter-sender "Alice Smith"
```

#### 3. Recipient Filtering

Filter by any recipient (To, CC, or BCC). Matching works on both email addresses and names.

**Examples:**
```bash
# Export emails to specific recipient
--filter-recipient "team@company.com"

# Export to multiple recipients
--filter-recipient "alice@example.com,bob@example.com"

# Export by recipient name
--filter-recipient "Project Team"
```

#### 4. Date Range Filtering

Filter emails within a specific time period. Dates should be in ISO 8601 format (YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS).

**Examples:**
```bash
# Export emails from 2024 onwards
--filter-date-from "2024-01-01"

# Export emails until end of 2023
--filter-date-to "2023-12-31"

# Export emails in specific month
--filter-date-from "2024-06-01" --filter-date-to "2024-06-30"

# Export with time precision
--filter-date-from "2024-06-15T09:00:00" --filter-date-to "2024-06-15T17:00:00"
```

**Timezone Handling:**
- Dates without timezone are treated as UTC
- Include timezone for precision: `2024-06-15T09:00:00-05:00`
- ProtonMail stores message times in UTC

#### 5. Subject Keyword Filtering

Filter by keywords in the subject line. Multiple keywords use OR logic (case-insensitive).

**Examples:**
```bash
# Export invoices
--filter-subject "invoice"

# Export invoices or receipts
--filter-subject "invoice,receipt"

# Partial match works
--filter-subject "meeting" # matches "Team Meeting", "Meeting Notes", etc.
```

### Combining Filter Criteria

#### AND Logic (Default)

By default, all specified filter types must match (AND logic). A message must satisfy ALL criteria to be exported.

**Example:**
```bash
# Export Inbox emails from Alice in June 2024
./proton-mail-export-cli \
  --filter-labels "0" \
  --filter-sender "alice@example.com" \
  --filter-date-from "2024-06-01" \
  --filter-date-to "2024-06-30" \
  --operation backup
```

This exports messages that are:
- IN the Inbox AND
- FROM alice@example.com AND
- SENT between June 1-30, 2024

#### OR Logic

Use `--filter-logic OR` to export messages matching ANY criterion.

**Example:**
```bash
# Export Inbox OR emails from Alice (regardless of folder)
./proton-mail-export-cli \
  --filter-labels "0" \
  --filter-sender "alice@example.com" \
  --filter-logic OR \
  --operation backup
```

This exports messages that are:
- IN the Inbox OR
- FROM alice@example.com

**Note:** Within each filter type, OR logic is always used (e.g., multiple labels, multiple senders).

## Complete Usage Examples

### Use Case 1: Export Work Emails from 2024

```bash
./proton-mail-export-cli \
  --filter-sender "@company.com" \
  --filter-date-from "2024-01-01" \
  --operation backup \
  --dir ./work-backup-2024
```

### Use Case 2: Export Project-Related Emails

```bash
./proton-mail-export-cli \
  --filter-subject "ProjectX,Project X" \
  --filter-recipient "projectx@company.com" \
  --filter-logic OR \
  --operation backup \
  --dir ./projectx-backup
```

### Use Case 3: Export Recent Important Emails

```bash
# Get custom label ID first
./proton-mail-export-cli --list-labels | grep "Important"
# Output: ID: abc123xyz | Name: Important

# Export important emails from last 30 days
./proton-mail-export-cli \
  --filter-labels "abc123xyz" \
  --filter-date-from "$(date -d '30 days ago' +%Y-%m-%d)" \
  --operation backup \
  --dir ./important-recent
```

### Use Case 4: Export All Communication with Specific Person

```bash
./proton-mail-export-cli \
  --filter-sender "alice@example.com" \
  --filter-recipient "alice@example.com" \
  --filter-logic OR \
  --operation backup \
  --dir ./alice-correspondence
```

### Use Case 5: Archive Old Sent Emails

```bash
./proton-mail-export-cli \
  --filter-labels "2" \
  --filter-date-to "2020-12-31" \
  --operation backup \
  --dir ./sent-archive-pre-2021
```

## Performance Considerations

### Optimization Tips

1. **Use Server-Side Filtering When Possible**
   - Single label filters are processed server-side (fastest)
   - Multiple criteria require client-side processing

2. **Date Range Filters Are Efficient**
   - Date filtering happens early in the pipeline
   - Reduces data transfer and processing time

3. **Combine Filters Strategically**
   - Use the most restrictive filter first
   - Example: Date range + sender is faster than sender + date range

4. **Monitor Progress**
   - Progress bar shows actual filtered message count
   - Large mailboxes with restrictive filters process quickly

### Performance Examples

| Mailbox Size | Filter Type | Approx. Time | Data Transfer |
|--------------|-------------|--------------|---------------|
| 50,000 emails | No filter | ~2 hours | 10 GB |
| 50,000 emails | Single label (Inbox: 2,000) | ~10 minutes | 400 MB |
| 50,000 emails | Date range (last month: 500) | ~5 minutes | 100 MB |
| 50,000 emails | Sender + date (50 emails) | ~2 minutes | 10 MB |

## Troubleshooting

### No Messages Exported

**Problem:** Filter is too restrictive or incorrect

**Solutions:**
1. Verify label IDs with `--list-labels`
2. Check date format (use ISO 8601: YYYY-MM-DD)
3. Test sender/recipient email addresses (case doesn't matter, but spelling does)
4. Try relaxing one criterion at a time to isolate the issue

### Unexpected Messages Included

**Problem:** Filter logic not working as expected

**Solutions:**
1. Check `--filter-logic` setting (AND vs OR)
2. Remember: Multiple values within one filter type always use OR logic
3. Verify sender/recipient partial matching behavior
4. Check timezone for date filters

### Slow Performance

**Problem:** Filter still processes many messages

**Solutions:**
1. Add date range to reduce search space
2. Use single label filter if possible (server-side optimization)
3. Ensure disk has sufficient speed (SSD recommended)
4. Check network connection quality

### Filter Not Applied

**Problem:** All emails being exported despite filter flags

**Solutions:**
1. Verify flags are passed correctly: `--filter-labels "0"` (with equals sign and quotes)
2. Check for typos in filter flag names
3. Ensure operation is `backup` not `restore`
4. Check CLI output for "Applying filters" message

## Advanced Topics

### Finding Custom Label IDs

```bash
# List all labels with IDs
./proton-mail-export-cli --list-labels

# Output example:
# System Labels:
#   ID: 0  | Name: Inbox
#   ID: 2  | Name: Sent
# Custom Labels:
#   ID: abc123 | Name: Work
#   ID: xyz789 | Name: Personal
```

### Programmatic Filter Creation

When using the Go API directly:

```go
import "github.com/ProtonMail/export-tool/internal/mail"

// Create comprehensive filter
filter := mail.NewExportFilter()
filter.LabelIDs = []string{"0", "2"}           // Inbox and Sent
filter.Senders = []string{"@company.com"}       // From company
filter.DateFrom = &time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
filter.SubjectKeywords = []string{"invoice"}
filter.CombineLogic = mail.FilterLogicAND      // All must match

// Use in export task
task := mail.NewExportTask(ctx, exportPath, session, filter)
```

### Filter Testing

Test your filters before running full export:

```bash
# Add --dry-run flag (if implemented)
./proton-mail-export-cli \
  --filter-sender "alice@example.com" \
  --filter-date-from "2024-01-01" \
  --dry-run \
  --operation backup
# Output: Would export 42 messages
```

## Migration from Simple Label Filtering

If you're currently using the basic `--filter` or `-f` flag:

**Old syntax (still works):**
```bash
./proton-mail-export-cli -f "0,2" -o backup
```

**New equivalent:**
```bash
./proton-mail-export-cli --filter-labels "0,2" --operation backup
```

Both syntaxes are supported for backward compatibility.

## See Also

- [FILTER_EXPORT_USAGE.md](FILTER_EXPORT_USAGE.md) - Basic filtering guide
- [README.md](README.md) - General usage and features
- [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md) - Technical implementation details

## API Filter Support Reference

Based on go-proton-api MessageFilter capabilities:

| Filter Type | Server-Side | Client-Side | Notes |
|-------------|-------------|-------------|-------|
| Single Label | ✅ | ✅ | Fastest option |
| Multiple Labels | ❌ | ✅ | Requires client-side processing |
| Sender | ❌ | ✅ | Full text matching |
| Recipient | ❌ | ✅ | Searches To/CC/BCC |
| Date Range | ❌ | ✅ | Efficient early filtering |
| Subject Keywords | ❌ | ✅ | Case-insensitive search |

**Legend:**
- ✅ Supported
- ❌ Not supported
- Server-side: Processed by ProtonMail API (faster, less bandwidth)
- Client-side: Processed locally during export (more flexible, still efficient)
