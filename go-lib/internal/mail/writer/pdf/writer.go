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

package pdf

import (
	"bytes"

	"github.com/ProtonMail/export-tool/internal/mail"
	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
)

// AttachmentMode defines how attachments are handled in PDF export
type AttachmentMode int

const (
	// AttachmentModeListOnly lists attachments in PDF, saves separately
	AttachmentModeListOnly AttachmentMode = iota
	// AttachmentModeEmbed embeds small attachments in PDF (when supported)
	AttachmentModeEmbed
	// AttachmentModeZip creates companion ZIP file of attachments
	AttachmentModeZip
)

// PDFMessageWriter writes email messages as PDF files
// Phase 2 - Implementation scaffolded, TODO: implement full functionality
type PDFMessageWriter struct {
	msg            proton.FullMessage
	pdfData        bytes.Buffer
	attachmentMode AttachmentMode
}

// NewPDFMessageWriter creates a new PDF message writer
// TODO: Implement PDF generation logic
func NewPDFMessageWriter(msg proton.FullMessage, mode AttachmentMode) *PDFMessageWriter {
	return &PDFMessageWriter{
		msg:            msg,
		attachmentMode: mode,
	}
}

// WriteMessage implements the MessageWriter interface
// TODO: Implement actual PDF writing
func (p *PDFMessageWriter) WriteMessage(dir string, tempDir string, log *logrus.Entry, checker utils.IntegrityChecker) error {
	log.Warn("PDF export is not yet implemented (Phase 2)")
	// TODO: Implement PDF generation and writing
	// 1. Render email content to PDF using selected renderer
	// 2. Handle attachments according to attachmentMode
	// 3. Write PDF file to dir
	return nil
}

// GetMetadata returns the message metadata
func (p *PDFMessageWriter) GetMetadata() mail.MessageMetadata {
	return mail.NewMessageMetadata(mail.MessageWriterTypePDF, &p.msg.Message)
}

// Renderer defines the interface for HTML-to-PDF rendering
type Renderer interface {
	// RenderEmail converts email content to PDF format
	RenderEmail(msg proton.FullMessage) ([]byte, error)
}

// RendererConfig configures PDF rendering behavior
type RendererConfig struct {
	// UseAdvancedRenderer enables external renderer if available
	UseAdvancedRenderer bool
	// MaxAttachmentSize for embedding (bytes)
	MaxAttachmentSize int64
	// PreserveStyling attempts to maintain HTML styling
	PreserveStyling bool
}

// DefaultRendererConfig returns sensible defaults
func DefaultRendererConfig() *RendererConfig {
	return &RendererConfig{
		UseAdvancedRenderer: false,
		MaxAttachmentSize:   10 * 1024 * 1024, // 10MB
		PreserveStyling:     false,
	}
}
