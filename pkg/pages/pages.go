// Package pages provides document page and paragraph structures
package pages

import (
	"encoding/json"
	"strings"
	"time"
	"unicode"
)

// DocumentParagraph represents a single paragraph within a page
type DocumentParagraph struct {
	// Paragraph number within the page (1-indexed)
	ParagraphNumber uint `json:"paragraph_number"`

	// Text content of this paragraph
	TextContent string `json:"text_content"`

	// Optional source reference (filename, section, etc.)
	SourceRef *string `json:"source_ref,omitempty"`

	// Word count for this paragraph
	WordCount *uint `json:"word_count,omitempty"`

	// Character count for this paragraph
	CharacterCount *uint `json:"character_count,omitempty"`
}

// NewDocumentParagraph creates a new document paragraph with automatic word/character counting
func NewDocumentParagraph(paragraphNumber uint, textContent string) *DocumentParagraph {
	wordCount := uint(countWords(textContent))
	charCount := uint(len(textContent))

	return &DocumentParagraph{
		ParagraphNumber: paragraphNumber,
		TextContent:     textContent,
		WordCount:       &wordCount,
		CharacterCount:  &charCount,
	}
}

// WithSourceRef sets the source reference
func (p *DocumentParagraph) WithSourceRef(sourceRef string) *DocumentParagraph {
	p.SourceRef = &sourceRef
	return p
}

// IsEmpty checks if paragraph is empty
func (p *DocumentParagraph) IsEmpty() bool {
	return strings.TrimSpace(p.TextContent) == ""
}

// DocumentPage represents a single page within a document
type DocumentPage struct {
	// Page number (1-indexed)
	PageNumber uint `json:"page_number"`

	// Text content of this page
	TextContent string `json:"text_content"`

	// Optional source reference (filename, section, etc.)
	SourceRef *string `json:"source_ref,omitempty"`

	// Word count for this page
	WordCount *uint `json:"word_count,omitempty"`

	// Character count for this page
	CharacterCount *uint `json:"character_count,omitempty"`
}

// NewDocumentPage creates a new document page
func NewDocumentPage(pageNumber uint) *DocumentPage {
	return &DocumentPage{
		PageNumber:  pageNumber,
		TextContent: "",
	}
}

// NewDocumentPageWithText creates a new document page with text content
func NewDocumentPageWithText(pageNumber uint, textContent string) *DocumentPage {
	wordCount := uint(countWords(textContent))
	charCount := uint(len(textContent))

	return &DocumentPage{
		PageNumber:     pageNumber,
		TextContent:    textContent,
		WordCount:      &wordCount,
		CharacterCount: &charCount,
	}
}

// WithSourceRef sets the source reference
func (p *DocumentPage) WithSourceRef(sourceRef string) *DocumentPage {
	p.SourceRef = &sourceRef
	return p
}

// SetTextContent sets the text content and updates word/character counts
func (p *DocumentPage) SetTextContent(textContent string) {
	p.TextContent = textContent
	wordCount := uint(countWords(textContent))
	charCount := uint(len(textContent))
	p.WordCount = &wordCount
	p.CharacterCount = &charCount
}

// GetTextContent gets the text content of this page
func (p *DocumentPage) GetTextContent() string {
	return p.TextContent
}

// GetWordCount gets word count for this page
func (p *DocumentPage) GetWordCount() uint {
	if p.WordCount != nil {
		return *p.WordCount
	}
	return uint(countWords(p.TextContent))
}

// GetCharacterCount gets character count for this page
func (p *DocumentPage) GetCharacterCount() uint {
	if p.CharacterCount != nil {
		return *p.CharacterCount
	}
	return uint(len(p.TextContent))
}

// IsEmpty checks if page is empty
func (p *DocumentPage) IsEmpty() bool {
	return strings.TrimSpace(p.TextContent) == ""
}

// ExtractionInfo contains information about the extraction process
type ExtractionInfo struct {
	// Tool that performed the extraction
	ExtractorName string `json:"extractor_name"`

	// Version of the extraction tool
	ExtractorVersion string `json:"extractor_version"`

	// Timestamp of extraction
	ExtractedAt *time.Time `json:"extracted_at,omitempty"`

	// Any warnings or notes about the extraction
	Warnings []string `json:"warnings"`
}

// NewExtractionInfo creates extraction info for a specific tool
func NewExtractionInfo(extractorName, version string) *ExtractionInfo {
	now := time.Now()
	return &ExtractionInfo{
		ExtractorName:    extractorName,
		ExtractorVersion: version,
		ExtractedAt:      &now,
		Warnings:         make([]string, 0),
	}
}

// AddWarning adds a warning message
func (e *ExtractionInfo) AddWarning(warning string) {
	e.Warnings = append(e.Warnings, warning)
}

// DocumentPages represents the complete document with pages
type DocumentPages struct {
	// Source file path
	SourceFile string `json:"source_file"`

	// Document title (from metadata if available)
	DocumentTitle *string `json:"document_title,omitempty"`

	// Document format/type (PDF, EPUB, etc.)
	DocumentType string `json:"document_type"`

	// Total number of pages in document
	TotalPages uint `json:"total_pages"`

	// All pages in the document
	Pages []DocumentPage `json:"pages"`

	// Metadata about the extraction process
	ExtractionInfo ExtractionInfo `json:"extraction_info"`
}

// NewDocumentPages creates a new document pages structure
func NewDocumentPages(sourceFile, documentType string, totalPages uint) *DocumentPages {
	return &DocumentPages{
		SourceFile:     sourceFile,
		DocumentType:   documentType,
		TotalPages:     totalPages,
		Pages:          make([]DocumentPage, 0),
		ExtractionInfo: *NewExtractionInfo("unknown", "unknown"),
	}
}

// WithTitle sets the document title
func (d *DocumentPages) WithTitle(title string) *DocumentPages {
	d.DocumentTitle = &title
	return d
}

// AddPage adds a page to the document
func (d *DocumentPages) AddPage(page DocumentPage) {
	d.Pages = append(d.Pages, page)
}

// GetPage gets a specific page by number (1-indexed)
func (d *DocumentPages) GetPage(pageNumber uint) *DocumentPage {
	for i := range d.Pages {
		if d.Pages[i].PageNumber == pageNumber {
			return &d.Pages[i]
		}
	}
	return nil
}

// GetAllText gets all text content concatenated
func (d *DocumentPages) GetAllText() string {
	var parts []string
	for _, page := range d.Pages {
		pageText := page.GetTextContent()
		if pageText != "" {
			parts = append(parts, pageText)
		}
	}
	return strings.Join(parts, "\n\n")
}

// TotalWordCount gets total word count across all pages
func (d *DocumentPages) TotalWordCount() uint {
	total := uint(0)
	for _, page := range d.Pages {
		total += page.GetWordCount()
	}
	return total
}

// TotalCharacterCount gets total character count across all pages
func (d *DocumentPages) TotalCharacterCount() uint {
	total := uint(0)
	for _, page := range d.Pages {
		total += page.GetCharacterCount()
	}
	return total
}

// IsEmpty checks if document is empty
func (d *DocumentPages) IsEmpty() bool {
	return len(d.Pages) == 0
}

// ToJSON converts document pages to JSON string
func (d *DocumentPages) ToJSON() (string, error) {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// countWords counts words in a string using simple whitespace splitting
func countWords(text string) int {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}

	// Split by unicode whitespace
	fields := strings.FieldsFunc(text, unicode.IsSpace)
	return len(fields)
}