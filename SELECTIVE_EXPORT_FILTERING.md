# Selective Export Filtering

The Proton Mail Export Tool supports comprehensive filtering to export only specific subsets of your mailbox. Filters can be combined to create precise export criteria.

## Available Filters

### Label/Folder Filtering

Filter messages by label or folder ID.

```bash
# Export only Inbox
proton-mail-export-cli --operation backup --label "0"

# Export Inbox and Sent (multiple labels use OR logic)
proton-mail-export-cli --operation backup --label "0" --label "2"
```

**Common System Label IDs:**
- Inbox: `0`
- Drafts: `1`
- Sent: `2`
- Starred: `10`
- Archive: `6`
- Spam: `4`
- Trash: `3`
- All Mail: `5`

Custom folders and labels have alphanumeric IDs that can be found in a previous export's `labels.json` file.

### Date Range Filtering

Filter messages by date using `--after` and `--before` flags.

**Supported Date Formats:**
- Unix timestamp (seconds since epoch): `1609459200`
- RFC3339: `2021-01-01T00:00:00Z`
- YYYY-MM-DD: `2021-01-01`

```bash
# Export messages from 2023 onwards
proton-mail-export-cli --operation backup --after "2023-01-01"

# Export messages before a specific date
proton-mail-export-cli --operation backup --before "2023-12-31"

# Export messages in a specific date range
proton-mail-export-cli --operation backup --after "2023-01-01" --before "2023-12-31"

# Using Unix timestamp
proton-mail-export-cli --operation backup --after "1672531200"

# Using RFC3339
proton-mail-export-cli --operation backup --after "2023-01-01T00:00:00Z"
```

### Sender Email Filtering

Filter messages by sender email address.

```bash
# Export messages from a specific sender
proton-mail-export-cli --operation backup --from "alice@example.com"

# Export messages from multiple senders (OR logic)
proton-mail-export-cli --operation backup --from "alice@example.com" --from "bob@example.com"
```

Email matching is case-insensitive and whitespace is trimmed.

### Recipient Email Filtering

Filter messages by recipient email address (matches To, CC, or BCC fields).

```bash
# Export messages to a specific recipient
proton-mail-export-cli --operation backup --to "alice@example.com"

# Export messages to multiple recipients (OR logic)
proton-mail-export-cli --operation backup --to "alice@example.com" --to "bob@example.com"
```

### Domain Filtering

Filter messages by sender or recipient domain.

```bash
# Export messages from a specific domain
proton-mail-export-cli --operation backup --from-domain "example.com"

# Export messages to a specific domain
proton-mail-export-cli --operation backup --to-domain "example.com"

# Export messages between specific domains
proton-mail-export-cli --operation backup --from-domain "company.com" --to-domain "client.com"

# Multiple domains (OR logic)
proton-mail-export-cli --operation backup --from-domain "example.com" --from-domain "proton.me"
```

## Combining Filters

Filters of different types use AND logic - messages must match all filter types to be exported. Within the same filter type (e.g., multiple `--from` values), OR logic is used.

### Examples

**Export Inbox messages from 2023:**
```bash
proton-mail-export-cli --operation backup \
  --label "0" \
  --after "2023-01-01" \
  --before "2024-01-01"
```

**Export messages from specific sender in a date range:**
```bash
proton-mail-export-cli --operation backup \
  --from "boss@company.com" \
  --after "2023-06-01"
```

**Export messages between specific domains in Sent folder:**
```bash
proton-mail-export-cli --operation backup \
  --label "2" \
  --from-domain "mycompany.com" \
  --to-domain "client.com"
```

**Complex filter combining multiple criteria:**
```bash
proton-mail-export-cli --operation backup \
  --label "0" --label "2" \          # Inbox OR Sent
  --after "2023-01-01" \              # AND after Jan 1, 2023
  --from "alice@example.com" \        # AND from alice
  --to-domain "company.com"           # AND to company.com domain
```

## Environment Variables

All filter flags can also be set via environment variables:

```bash
export ET_FILTER_LABELS="0,2"
export ET_FILTER_AFTER="2023-01-01"
export ET_FILTER_BEFORE="2023-12-31"
export ET_FILTER_FROM="alice@example.com,bob@example.com"
export ET_FILTER_TO="recipient@example.com"
export ET_FILTER_FROM_DOMAIN="example.com"
export ET_FILTER_TO_DOMAIN="proton.me"

proton-mail-export-cli --operation backup
```

## Filter Semantics

### AND Logic (Between Filter Types)

When you specify filters of different types, a message must satisfy ALL of them:

```bash
# Message must be in Inbox AND from alice AND after 2023
--label "0" --from "alice@example.com" --after "2023-01-01"
```

### OR Logic (Within Same Filter Type)

When you specify multiple values for the same filter type, a message can match ANY of them:

```bash
# Message can be in Inbox OR Sent
--label "0" --label "2"

# Message can be from alice OR bob
--from "alice@example.com" --from "bob@example.com"
```

## Performance Benefits

Using filters significantly improves export performance:

- **Reduced Download Time**: Only downloads filtered messages
- **Lower Disk Usage**: Exports only the subset you need
- **Faster Processing**: Processes fewer messages overall
- **Lower Memory Usage**: Streaming filter application

**Example**: If you have 50,000 emails but only need 2,000 from a specific label, filtering can reduce export time by 95%.

## Backwards Compatibility

- If no filters are specified, all messages are exported (default behavior)
- All filters are optional
- Existing scripts without filters continue to work unchanged

## Filter Application

Filters are applied during the metadata fetching stage, before downloading message content. This means:

1. Messages are filtered as early as possible in the pipeline
2. Filtered messages are never downloaded, saving bandwidth
3. Progress reporting accounts for filtered messages correctly
4. Export resumes work correctly with filters

## Logging

When filters are active, the tool logs:

- A summary of active filters at the start of export (console and log file)
- The number of messages filtered vs. exported
- Filter criteria in the session log for troubleshooting

## Troubleshooting

### No messages exported

- Verify filter values are correct
- Check that date formats are valid
- Ensure email addresses and domains are spelled correctly
- Review the session log for filter details

### Unexpected messages in export

- Remember that OR logic applies within filter types
- Check that all filter criteria are specified correctly
- Verify label IDs from a previous export's `labels.json`

### Date filtering not working as expected

- Ensure date format is correct (YYYY-MM-DD, RFC3339, or Unix timestamp)
- Remember `--after` is inclusive, `--before` is exclusive
- Dates in YYYY-MM-DD format use UTC midnight

### Performance not improved

- Ensure filters actually reduce the message set significantly
- Check that filters are actually being applied (look for filter summary in output)
- Verify there are no typos in filter values

## Examples by Use Case

### Export Recent Messages Only

```bash
# Last 30 days (adjust date accordingly)
proton-mail-export-cli --operation backup --after "2024-10-01"
```

### Export Work Communications

```bash
# Messages from/to work domain
proton-mail-export-cli --operation backup \
  --from-domain "mycompany.com" \
  --to-domain "mycompany.com"
```

### Export Specific Conversation

```bash
# Messages with a specific person
proton-mail-export-cli --operation backup \
  --from "person@example.com" \
  --to "person@example.com"
```

### Export Project-Related Emails

```bash
# Assuming you have a project label/folder
proton-mail-export-cli --operation backup --label "PROJECT_LABEL_ID"
```

### Compliance Export

```bash
# Export all business emails from 2023 for audit
proton-mail-export-cli --operation backup \
  --from-domain "company.com" \
  --after "2023-01-01" \
  --before "2024-01-01"
```

## Notes

- Client-side filtering is used (no server-side API support currently)
- All metadata fields are available for filtering
- Filters are logged in the session log for audit trails
- Progress reporting shows total messages vs. filtered messages
- Resume functionality works correctly with filters
- The `labels.json` file in exports contains all labels regardless of filters
