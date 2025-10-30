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

package mail

// PDFMessageWriter is a placeholder interface for future PDF export functionality.
// This interface will be implemented in a later phase to support exporting messages as PDF files.
//
// Example future usage:
//   writer := NewPDFMessageWriter(config)
//   err := writer.WriteMessage(message)
type PDFMessageWriter interface {
	// WriteMessage writes a message to a PDF file.
	// Returns the path to the created PDF file and any error.
	WriteMessage(msg Message) (string, error)

	// WriteBatch writes multiple messages to a single or multiple PDF files.
	// The exact behavior (one PDF per message vs combined PDF) will be determined
	// based on the configuration.
	WriteBatch(messages []Message) ([]string, error)

	// Close finalizes any pending writes and cleans up resources.
	Close() error
}

// Message represents a complete email message for PDF export.
// This is a placeholder type that will be properly defined when PDF export is implemented.
type Message struct {
	ID          string
	Subject     string
	From        string
	To          []string
	CC          []string
	BCC         []string
	Date        int64
	Body        string
	Attachments []Attachment
}

// Attachment represents an email attachment.
type Attachment struct {
	Name string
	Data []byte
}

// PDFWriterConfig holds configuration for the PDF writer.
// This is a placeholder that will be expanded based on requirements.
type PDFWriterConfig struct {
	// OutputDir is the directory where PDF files will be written
	OutputDir string

	// IncludeAttachments determines whether to include attachment information in PDFs
	IncludeAttachments bool

	// CombineMessages determines whether to combine multiple messages into one PDF
	CombineMessages bool

	// MaxMessagesPerPDF limits the number of messages per PDF when CombineMessages is true
	MaxMessagesPerPDF int
}
