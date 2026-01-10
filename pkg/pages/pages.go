// Package pages provides file chip and paragraph structures
package pages

import (
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

// FileChip represents a single page within a document
type FileChip struct {
	// Page number (1-indexed)
	OrderIndex uint `json:"order_index"`

	// Text content of this page
	TextContent string `json:"text_content"`

	// Optional source reference (filename, section, etc.)
	SourceRef *string `json:"source_ref,omitempty"`

	// Word count for this page
	WordCount *uint `json:"word_count,omitempty"`

	// Character count for this page
	CharacterCount *uint `json:"character_count,omitempty"`
}

// NewFileChip creates a new file chip
func NewFileChip(orderIndex uint) *FileChip {
	return &FileChip{
		OrderIndex:  orderIndex,
		TextContent: "",
	}
}

// NewFileChipWithText creates a new file chip with text content
func NewFileChipWithText(orderIndex uint, textContent string) *FileChip {
	wordCount := uint(countWords(textContent))
	charCount := uint(len(textContent))

	return &FileChip{
		OrderIndex:     orderIndex,
		TextContent:    textContent,
		WordCount:      &wordCount,
		CharacterCount: &charCount,
	}
}

// WithSourceRef sets the source reference
func (p *FileChip) WithSourceRef(sourceRef string) *FileChip {
	p.SourceRef = &sourceRef
	return p
}

// SetTextContent sets the text content and updates word/character counts
func (p *FileChip) SetTextContent(textContent string) {
	p.TextContent = textContent
	wordCount := uint(countWords(textContent))
	charCount := uint(len(textContent))
	p.WordCount = &wordCount
	p.CharacterCount = &charCount
}

// GetTextContent gets the text content of this page
func (p *FileChip) GetTextContent() string {
	return p.TextContent
}

// GetWordCount gets word count for this page
func (p *FileChip) GetWordCount() uint {
	if p.WordCount != nil {
		return *p.WordCount
	}
	return uint(countWords(p.TextContent))
}

// GetCharacterCount gets character count for this page
func (p *FileChip) GetCharacterCount() uint {
	if p.CharacterCount != nil {
		return *p.CharacterCount
	}
	return uint(len(p.TextContent))
}

// IsEmpty checks if page is empty
func (p *FileChip) IsEmpty() bool {
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

// GroundChips represents the complete document with pages
type GroundChips struct {
	// Source file path
	SourceFile string `json:"source_file"`

	// Document title (from metadata if available)
	DocumentTitle *string `json:"document_title,omitempty"`

	// Document format/type (PDF, EPUB, etc.)
	DocumentType string `json:"document_type"`

	// Total number of pages in document
	TotalPages uint `json:"total_pages"`

	// All pages in the document
	Pages []FileChip `json:"pages"`

	// Metadata about the extraction process
	ExtractionInfo ExtractionInfo `json:"extraction_info"`
}

// NewGroundChips creates a new file chips structure
func NewGroundChips(sourceFile, documentType string, totalPages uint) *GroundChips {
	return &GroundChips{
		SourceFile:     sourceFile,
		DocumentType:   documentType,
		TotalPages:     totalPages,
		Pages:          make([]FileChip, 0),
		ExtractionInfo: *NewExtractionInfo("unknown", "unknown"),
	}
}

// WithTitle sets the document title
func (d *GroundChips) WithTitle(title string) *GroundChips {
	d.DocumentTitle = &title
	return d
}

// AddPage adds a page to the document
func (d *GroundChips) AddPage(page FileChip) {
	d.Pages = append(d.Pages, page)
}

// GetChip gets a specific page by number (1-indexed)
func (d *GroundChips) GetChip(orderIndex uint) *FileChip {
	for i := range d.Pages {
		if d.Pages[i].OrderIndex == orderIndex {
			return &d.Pages[i]
		}
	}
	return nil
}

// GetAllText gets all text content concatenated
func (d *GroundChips) GetAllText() string {
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
func (d *GroundChips) TotalWordCount() uint {
	total := uint(0)
	for _, page := range d.Pages {
		total += page.GetWordCount()
	}
	return total
}

// TotalCharacterCount gets total character count across all pages
func (d *GroundChips) TotalCharacterCount() uint {
	total := uint(0)
	for _, page := range d.Pages {
		total += page.GetCharacterCount()
	}
	return total
}

// IsEmpty checks if document is empty
func (d *GroundChips) IsEmpty() bool {
	return len(d.Pages) == 0
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