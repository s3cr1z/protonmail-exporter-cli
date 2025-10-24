# Folder/Label Filter Export Feature

## Overview

The Proton Mail Export Tool now supports filtering exports by specific folders or labels. This feature allows you to export only emails from selected folders instead of exporting all emails, significantly improving efficiency when you need specific subsets of your mailbox.

## Usage

### Command Line

Use the `--filter` or `-f` option to specify label IDs (comma-separated):

```bash
./proton-mail-export-cli --filter "label-id-1,label-id-2" --operation backup
```

### Environment Variable

Alternatively, set the `ET_FILTER_LABELS` environment variable:

```bash
export ET_FILTER_LABELS="label-id-1,label-id-2"
./proton-mail-export-cli --operation backup
```

## Finding Label IDs

To find the label IDs for your folders/labels:

### Method 1: Using --list-labels (Recommended)

Simply run the export tool with the `--list-labels` flag to see all available folders and labels:

```bash
./proton-mail-export-cli --list-labels
```

This will log you in and display all your folders, labels, and their IDs without performing any export.

### Method 2: From Previous Export

1. Check an existing export if you have one
2. Open the `labels.json` file in your export directory
3. Find the label IDs for the folders/labels you want to filter

The `labels.json` file structure:
```json
{
  "Version": 1,
  "Data": [
    {
      "ID": "label-id-here",
      "Name": "Folder/Label Name",
      "Path": "Folder/Label Path",
      "Color": "#color",
      "Type": 1,
      ...
    }
  ]
}
```

### Method 3: Common System Labels

ProtonMail has several system labels with standard IDs:

- **All Mail**: `5` (default - exports all messages)
- **Inbox**: `0`
- **Drafts**: `1`
- **Sent**: `2`
- **Starred**: `10`
- **Archive**: `6`
- **Spam**: `4`
- **Trash**: `3`

Note: Custom folders and labels have unique IDs that you'll need to find using Method 1 or Method 2.

## Examples

### Export only Inbox messages
```bash
./proton-mail-export-cli -f "0" -o backup
```

### Export Inbox and Sent messages
```bash
./proton-mail-export-cli -f "0,2" -o backup
```

### Export specific custom folder
```bash
# First, find the label ID from labels.json in a previous export
./proton-mail-export-cli -f "your-custom-label-id" -o backup
```

### Combined with other options
```bash
./proton-mail-export-cli \
  --user user@proton.me \
  --filter "0,2" \
  --dir /path/to/export \
  --operation backup
```

## Quick Start Workflow

Here's the typical workflow to export a specific folder:

1. **List all your folders and labels**:
   ```bash
   ./proton-mail-export-cli --list-labels
   ```

2. **Find your folder** in the output. For example, if you see:
   ```
   Folders:
     ID: AbCdEfGhIjKlMnOp | Name: law_office | Path: law_office
   ```

3. **Export that folder**:
   ```bash
   ./proton-mail-export-cli --filter "AbCdEfGhIjKlMnOp" --operation backup
   ```

That's it! No need to export everything first.

## How It Works

1. **Filter Application**: The filter is applied during the metadata fetching stage
2. **Message Selection**: Only messages that have at least one of the specified label IDs are included
3. **Multiple Labels**: When multiple label IDs are provided, messages matching ANY of them are exported (OR logic)
4. **No Filter**: If no filter is specified (empty string or not provided), all messages are exported (default behavior)

## Performance Benefits

Using filters can significantly improve export performance:

- **Reduced Download Time**: Only downloads messages from specified folders
- **Lower Disk Usage**: Exports only the subset you need
- **Faster Processing**: Processes fewer messages overall

Example: If you have 50,000 total emails but only need 2,000 from your Inbox, using `--filter "0"` will export only those 2,000 messages instead of all 50,000.

## Notes

- The filter applies to the backup/export operation only (not restore)
- All labels and metadata are still exported for reference
- Messages can have multiple labels, so a message in multiple filtered folders will only be exported once
- The feature maintains backward compatibility - existing scripts without the filter will continue to work as before

## Troubleshooting

### No messages exported
- Verify the label IDs are correct (check labels.json from a previous export)
- Ensure the label IDs are comma-separated without spaces: `"0,2"` not `"0, 2"`

### Incorrect messages exported
- Double-check the label ID from labels.json
- Remember that system labels use numeric IDs (0-10)
- Custom labels use alphanumeric string IDs

### Performance not improved
- Ensure you're filtering to a smaller subset of your total emails
- Check that the filter string is actually being passed (look for the "Filtering export by label IDs:" message in output)
