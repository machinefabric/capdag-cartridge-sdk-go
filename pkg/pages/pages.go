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

// DisboundPage represents a single page within a document
type DisboundPage struct {
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

// NewDisboundPage creates a new file chip
func NewDisboundPage(orderIndex uint) *DisboundPage {
	return &DisboundPage{
		OrderIndex:  orderIndex,
		TextContent: "",
	}
}

// NewDisboundPageWithText creates a new file chip with text content
func NewDisboundPageWithText(orderIndex uint, textContent string) *DisboundPage {
	wordCount := uint(countWords(textContent))
	charCount := uint(len(textContent))

	return &DisboundPage{
		OrderIndex:     orderIndex,
		TextContent:    textContent,
		WordCount:      &wordCount,
		CharacterCount: &charCount,
	}
}

// WithSourceRef sets the source reference
func (p *DisboundPage) WithSourceRef(sourceRef string) *DisboundPage {
	p.SourceRef = &sourceRef
	return p
}

// SetTextContent sets the text content and updates word/character counts
func (p *DisboundPage) SetTextContent(textContent string) {
	p.TextContent = textContent
	wordCount := uint(countWords(textContent))
	charCount := uint(len(textContent))
	p.WordCount = &wordCount
	p.CharacterCount = &charCount
}

// GetTextContent gets the text content of this page
func (p *DisboundPage) GetTextContent() string {
	return p.TextContent
}

// GetWordCount gets word count for this page
func (p *DisboundPage) GetWordCount() uint {
	if p.WordCount != nil {
		return *p.WordCount
	}
	return uint(countWords(p.TextContent))
}

// GetCharacterCount gets character count for this page
func (p *DisboundPage) GetCharacterCount() uint {
	if p.CharacterCount != nil {
		return *p.CharacterCount
	}
	return uint(len(p.TextContent))
}

// IsEmpty checks if page is empty
func (p *DisboundPage) IsEmpty() bool {
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