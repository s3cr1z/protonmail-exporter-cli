# PDF Export Feature - Design Document

## Status: Phase 2 - Scaffolded (Not Yet Implemented)

This document outlines the design for PDF email export functionality. The types and interfaces are scaffolded, but full implementation is pending.

## Overview

The PDF export feature will allow users to export emails as PDF documents with configurable attachment handling and rendering quality.

## Architecture

### Directory Structure

```
go-lib/internal/mail/writer/pdf/
├── writer.go                    # Main PDFMessageWriter implementation
├── render/
│   ├── basic/
│   │   └── basic.go            # Basic HTML-to-text renderer
│   └── advanced/
│       └── advanced.go          # External tool renderer (feature-flagged)
└── attachment_handler.go        # Attachment processing logic
```

### Core Components

#### 1. PDFMessageWriter

Implements the `MessageWriter` interface for PDF output.

```go
type PDFMessageWriter struct {
    msg            proton.FullMessage
    pdfData        bytes.Buffer
    attachmentMode AttachmentMode
}
```

**Methods:**
- `WriteMessage()`: Generates PDF and writes to disk
- `GetMetadata()`: Returns message metadata with PDF type

#### 2. Attachment Modes

Three strategies for handling attachments:

```go
type AttachmentMode int

const (
    AttachmentModeListOnly  // Default: list in PDF, save separately
    AttachmentModeEmbed     // Embed small attachments (using pdfcpu)
    AttachmentModeZip       // Create companion ZIP file
)
```

**Mode Details:**

- **List-Only (Default)**
  - Attachments listed in PDF with filename, size, MIME type
  - Actual files saved alongside PDF in export directory
  - Most compatible, works with all PDF readers

- **Embed**
  - Small files (<10MB default) embedded in PDF
  - Uses pdfcpu's attachment feature (PDF 1.4+)
  - Size limits configurable
  - Fallback to list-only for large files

- **Zip**
  - All attachments bundled in separate .zip
  - PDF references the zip file
  - Good for emails with many attachments

#### 3. Renderer Interface

```go
type Renderer interface {
    RenderEmail(msg proton.FullMessage) ([]byte, error)
}
```

Two implementations planned:

**Basic Renderer** (Pure Go)
- HTML to plain text conversion
- Preserves links as footnotes
- Simple formatting (bold, italic, lists)
- Uses gofpdf for PDF generation
- No external dependencies

**Advanced Renderer** (Optional, Feature-Flagged)
- High-fidelity HTML rendering
- CSS support
- Complex layouts preserved
- Requires external tool:
  - wkhtmltopdf (recommended)
  - chromedp (alternative)
- Automatic fallback to basic renderer if unavailable

### PDF Structure

Each PDF will contain:

```
┌─────────────────────────────────────┐
│ HEADER SECTION                      │
│  Subject: [Email subject]           │
│  From:    [Sender]                  │
│  To:      [Recipients]              │
│  CC:      [CC recipients]           │
│  Date:    [ISO timestamp]           │
├─────────────────────────────────────┤
│ BODY SECTION                        │
│  [Rendered email content]           │
│  [Inline images as captions]        │
├─────────────────────────────────────┤
│ ATTACHMENTS SECTION                 │
│  □ filename1.pdf (1.2 MB)          │
│  □ image.jpg (345 KB)              │
│  [Embedded or referenced]           │
└─────────────────────────────────────┘
```

## Implementation Plan

### Library Selection

#### PDF Generation
**Recommended: gofpdf**
- Pure Go, no CGO
- Well-maintained
- Good documentation
- Limited HTML support (acceptable for basic renderer)

**Alternative: pdfcpu**
- Pure Go
- Excellent for attachment embedding
- Lower-level API
- Can be used alongside gofpdf

#### HTML Processing
**For Basic Renderer:**
- `jaytaylor/html2text` (already in dependencies!)
- `bluemonday` for sanitization

**For Advanced Renderer:**
- `chromedp` for Chrome DevTools Protocol
- Or shell out to `wkhtmltopdf`

### Configuration

```go
type RendererConfig struct {
    UseAdvancedRenderer bool    // Enable external renderer
    MaxAttachmentSize   int64   // Max size for embedding (bytes)
    PreserveStyling     bool    // Attempt to maintain HTML styling
}
```

### File Naming

Follow existing export structure:

```
mail_2024_06_15_120000/
├── labels.json
├── message-id-1.pdf                    # PDF file
├── message-id-1.meta.json              # Metadata
├── message-id-1_att1_document.pdf      # Attachment (list-only mode)
├── message-id-1_att2_image.jpg         # Attachment (list-only mode)
└── message-id-1_attachments.zip        # Attachments (zip mode)
```

## Usage Examples

### CLI (Future)

```bash
# Export as PDF with default settings
./proton-mail-export-cli --format pdf --operation backup

# Export with embedded attachments
./proton-mail-export-cli --format pdf --pdf-attachments embed --operation backup

# Export with advanced renderer (if available)
./proton-mail-export-cli --format pdf --pdf-renderer advanced --operation backup

# Export with ZIP attachments
./proton-mail-export-cli --format pdf --pdf-attachments zip --operation backup
```

### Programmatic (Future)

```go
import "github.com/ProtonMail/export-tool/internal/mail/writer/pdf"

// Create PDF writer with config
writer := pdf.NewPDFMessageWriter(
    fullMessage,
    pdf.AttachmentModeListOnly,
)

// Write to export directory
err := writer.WriteMessage(dir, tempDir, log, integrityChecker)
```

## Testing Strategy

### Unit Tests

1. **PDF Structure Tests**
   - Headers rendered correctly
   - Body content formatted properly
   - Attachment section lists files

2. **Attachment Mode Tests**
   - List-only: files saved separately
   - Embed: small files embedded
   - Zip: companion zip created

3. **Renderer Tests**
   - Basic: HTML to text conversion
   - Advanced: fallback behavior

### Golden File Tests

Representative emails:
- Plain text email
- Simple HTML email
- Multi-part mixed
- Inline images
- Many attachments (>10)
- Complex HTML with CSS
- Non-Latin characters

### Integration Tests

- End-to-end export to PDF
- Verify PDF can be opened
- Attachment extraction works

## Security Considerations

### HTML Sanitization

All HTML must be sanitized before rendering:
- Remove `<script>` tags
- Remove `<iframe>` and `<embed>`
- Strip dangerous attributes (`onclick`, `onerror`, etc.)
- Whitelist safe tags and attributes
- Use bluemonday with strict policy

### Path Traversal

Attachment filenames must be sanitized:
- Remove `../` sequences
- Strip absolute paths
- Limit filename length
- Validate against whitelist of safe characters

### External Renderer Safety

When using advanced renderer:
- Only activate with explicit opt-in flag
- Validate tool path
- Run with minimal privileges
- Timeout long-running renders
- Catch and handle crashes gracefully

## Performance Considerations

### Memory

- Stream PDF generation where possible
- Don't load all attachments in memory
- Process large emails in chunks

### Disk

- Use temporary directory for intermediate files
- Clean up on success/failure
- Atomic writes (write to temp, then move)

### Speed

Estimated times (for reference):
- Basic renderer: 5-10 emails/second
- Advanced renderer: 1-2 emails/second
- Bottleneck: HTML rendering, not PDF generation

## Known Limitations

1. **Complex CSS**
   - Basic renderer: very limited CSS support
   - Advanced renderer: depends on external tool capabilities

2. **Fonts**
   - Embedded fonts may not be preserved
   - Fallback to standard PDF fonts

3. **Interactive Elements**
   - Forms and buttons become static
   - JavaScript doesn't execute

4. **Large Attachments**
   - Embedding has size limits
   - Very large emails may timeout

5. **PDF/A Compliance**
   - Not PDF/A compliant (archival standard)
   - Future enhancement if needed

## Dependencies

To be added to go.mod:

```go
require (
    github.com/jung-kurt/gofpdf v1.16.2        // PDF generation
    github.com/microcosm-cc/bluemonday v1.0.25  // HTML sanitization
    // html2text already in dependencies
)
```

Optional (for advanced renderer):
```go
require (
    github.com/chromedp/chromedp v0.9.5  // Chrome DevTools Protocol
)
```

## Future Enhancements

1. **PDF/A Support**
   - Archival-quality PDFs
   - Long-term preservation

2. **Encryption**
   - Password-protected PDFs
   - Inherit from email encryption status

3. **Metadata**
   - Embed email metadata in PDF properties
   - Searchable by date, sender, etc.

4. **Bookmarks**
   - PDF bookmarks for navigation
   - Useful for long email threads

5. **Portfolio/Bundle**
   - Multiple emails in single PDF
   - Table of contents
   - Cross-email search

## Migration Notes

Existing exports (EML/MBOX) remain default and unchanged. PDF is an opt-in format:

```bash
# Old (still works)
./proton-mail-export-cli --operation backup

# New (PDF format)
./proton-mail-export-cli --format pdf --operation backup
```

## Acceptance Criteria

✅ Scaffolded (Current):
- [x] Type definitions created
- [x] Interface defined
- [x] Directory structure in place
- [x] Design documented

⏳ TODO (Phase 2 Implementation):
- [ ] Basic renderer implementation
- [ ] Attachment handling for all modes
- [ ] Golden file tests pass
- [ ] Integration with export pipeline
- [ ] CLI flags added
- [ ] Documentation updated
- [ ] Performance benchmarks
- [ ] Security review

## References

- [gofpdf Documentation](https://pkg.go.dev/github.com/jung-kurt/gofpdf)
- [pdfcpu Documentation](https://pkg.go.dev/github.com/pdfcpu/pdfcpu)
- [bluemonday Documentation](https://pkg.go.dev/github.com/microcosm-cc/bluemonday)
- [PDF Specification](https://www.adobe.com/devnet/pdf/pdf_reference.html)
