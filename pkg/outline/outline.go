// Package outline provides document outline and table of contents structures
package outline

import (
	"encoding/json"
	"time"
)

// TocEntry represents a single entry in a document's table of contents
type TocEntry struct {
	// The title/label of this TOC entry
	Title string `json:"title"`

	// Hierarchical level (0 = top level, 1 = first sublevel, etc.)
	Level uint `json:"level"`

	// Page/section number this entry points to (1-indexed, nil if no destination)
	Page *uint `json:"page,omitempty"`

	// Optional source reference (filename, anchor, etc.)
	SourceRef *string `json:"source_ref,omitempty"`

	// Child entries under this one
	Children []TocEntry `json:"children"`
}

// NewTocEntry creates a new TOC entry
func NewTocEntry(title string, level uint) *TocEntry {
	return &TocEntry{
		Title:    title,
		Level:    level,
		Children: make([]TocEntry, 0),
	}
}

// WithPage sets the page/section number
func (t *TocEntry) WithPage(page uint) *TocEntry {
	t.Page = &page
	return t
}

// WithSourceRef sets the source reference
func (t *TocEntry) WithSourceRef(sourceRef string) *TocEntry {
	t.SourceRef = &sourceRef
	return t
}

// AddChild adds a child entry
func (t *TocEntry) AddChild(child TocEntry) {
	t.Children = append(t.Children, child)
}

// CountAllEntries gets total number of entries (including children)
func (t *TocEntry) CountAllEntries() uint {
	count := uint(1) // This entry
	for _, child := range t.Children {
		count += child.CountAllEntries()
	}
	return count
}

// FindByTitle finds entry by title (depth-first search)
func (t *TocEntry) FindByTitle(title string) *TocEntry {
	if t.Title == title {
		return t
	}

	for i := range t.Children {
		if found := t.Children[i].FindByTitle(title); found != nil {
			return found
		}
	}

	return nil
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

// DocumentOutline represents the complete document outline structure
type DocumentOutline struct {
	// Source file path
	SourceFile string `json:"source_file"`

	// Document title (from metadata if available)
	DocumentTitle *string `json:"document_title,omitempty"`

	// Document format/type (PDF, EPUB, etc.)
	DocumentType string `json:"document_type"`

	// Total number of pages/sections in document
	TotalPages uint `json:"total_pages"`

	// Root-level TOC entries
	Entries []TocEntry `json:"entries"`

	// Whether this document has any outline/bookmarks
	HasOutline bool `json:"has_outline"`

	// Metadata about the extraction process
	ExtractionInfo ExtractionInfo `json:"extraction_info"`
}

// NewDocumentOutline creates a new document outline
func NewDocumentOutline(sourceFile, documentType string, totalPages uint) *DocumentOutline {
	return &DocumentOutline{
		SourceFile:     sourceFile,
		DocumentType:   documentType,
		TotalPages:     totalPages,
		Entries:        make([]TocEntry, 0),
		HasOutline:     false,
		ExtractionInfo: *NewExtractionInfo("unknown", "unknown"),
	}
}

// WithTitle sets the document title
func (d *DocumentOutline) WithTitle(title string) *DocumentOutline {
	d.DocumentTitle = &title
	return d
}

// AddEntry adds a root-level TOC entry
func (d *DocumentOutline) AddEntry(entry TocEntry) {
	d.Entries = append(d.Entries, entry)
	d.HasOutline = true
}

// TotalEntries gets total number of TOC entries (including all children)
func (d *DocumentOutline) TotalEntries() uint {
	total := uint(0)
	for _, entry := range d.Entries {
		total += entry.CountAllEntries()
	}
	return total
}

// IsEmpty checks if outline is empty
func (d *DocumentOutline) IsEmpty() bool {
	return len(d.Entries) == 0
}

// ToJSON converts outline to JSON string
func (d *DocumentOutline) ToJSON() (string, error) {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}