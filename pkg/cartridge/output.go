// Package cartridge provides output structures for document extraction
package cartridge

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/machinefabric/fgrnd-cartridge-sdk-go/pkg/metadata"
	"github.com/machinefabric/fgrnd-cartridge-sdk-go/pkg/outline"
	"github.com/machinefabric/fgrnd-cartridge-sdk-go/pkg/pages"
)

// OutputFormat represents output format options
type OutputFormat int

const (
	// Text format - Human-readable text format
	Text OutputFormat = iota
	// JSON format - JSON format
	JSON
	// Structured format - Structured data for programmatic use
	Structured
)

// ThumbnailInfo represents information about an extracted thumbnail
type ThumbnailInfo struct {
	// Image format (png, jpg, etc.)
	Format string `json:"format"`

	// Image dimensions if known [width, height]
	Dimensions *[2]uint32 `json:"dimensions,omitempty"`

	// File size of the image data
	SizeBytes uint64 `json:"size_bytes"`

	// Suggested filename
	Filename string `json:"filename"`
}

// ExtractionSummary represents summary of what was extracted
type ExtractionSummary struct {
	// Source file path
	SourceFile string `json:"source_file"`

	// Handler that performed the extraction
	HandlerName string `json:"handler_name"`

	// Extraction timestamp
	ExtractedAt time.Time `json:"extracted_at"`

	// What was successfully extracted
	ExtractedComponents []string `json:"extracted_components"`

	// Any errors or warnings
	Warnings []string `json:"warnings"`

	// Total processing time in milliseconds
	ProcessingTimeMs uint64 `json:"processing_time_ms"`
}

// NewExtractionSummary creates a new extraction summary
func NewExtractionSummary(sourceFile, handlerName string) *ExtractionSummary {
	return &ExtractionSummary{
		SourceFile:          sourceFile,
		HandlerName:         handlerName,
		ExtractedAt:         time.Now(),
		ExtractedComponents: make([]string, 0),
		Warnings:            make([]string, 0),
		ProcessingTimeMs:    0,
	}
}

// AddComponent adds a successfully extracted component
func (es *ExtractionSummary) AddComponent(component string) {
	es.ExtractedComponents = append(es.ExtractedComponents, component)
}

// AddWarning adds a warning
func (es *ExtractionSummary) AddWarning(warning string) {
	es.Warnings = append(es.Warnings, warning)
}

// SetProcessingTime sets the processing time
func (es *ExtractionSummary) SetProcessingTime(timeMs uint64) {
	es.ProcessingTimeMs = timeMs
}

// ExtractedData represents combined output containing multiple types of extracted data
type ExtractedData struct {
	// Document metadata
	Metadata *metadata.FileMetadata `json:"metadata,omitempty"`

	// Document outline
	Outline *outline.DocumentOutline `json:"outline,omitempty"`

	// Disbound pages with text content
	Pages []pages.DisboundPage `json:"pages,omitempty"`

	// Extracted text content (alternative to pages structure)
	TextContent *string `json:"text_content,omitempty"`

	// Thumbnail information
	Thumbnail *ThumbnailInfo `json:"thumbnail,omitempty"`

	// Extraction summary
	ExtractionSummary ExtractionSummary `json:"extraction_summary"`
}

// NewExtractedData creates new extracted data
func NewExtractedData(sourceFile, handlerName string) *ExtractedData {
	return &ExtractedData{
		ExtractionSummary: *NewExtractionSummary(sourceFile, handlerName),
	}
}

// WithMetadata adds metadata
func (ed *ExtractedData) WithMetadata(metadata *metadata.FileMetadata) *ExtractedData {
	ed.Metadata = metadata
	ed.ExtractionSummary.AddComponent("metadata")
	return ed
}

// WithOutline adds outline
func (ed *ExtractedData) WithOutline(outline *outline.DocumentOutline) *ExtractedData {
	ed.Outline = outline
	ed.ExtractionSummary.AddComponent("outline")
	return ed
}

// WithPages adds pages
func (ed *ExtractedData) WithPages(pageList []pages.DisboundPage) *ExtractedData {
	ed.Pages = pageList
	ed.ExtractionSummary.AddComponent("pages")
	return ed
}

// WithText adds text content
func (ed *ExtractedData) WithText(textContent string) *ExtractedData {
	ed.TextContent = &textContent
	ed.ExtractionSummary.AddComponent("text_content")
	return ed
}

// WithThumbnail adds thumbnail
func (ed *ExtractedData) WithThumbnail(thumbnail *ThumbnailInfo) *ExtractedData {
	ed.Thumbnail = thumbnail
	ed.ExtractionSummary.AddComponent("thumbnail")
	return ed
}

// AddWarning adds warning
func (ed *ExtractedData) AddWarning(warning string) {
	ed.ExtractionSummary.AddWarning(warning)
}

// SetProcessingTime sets processing time
func (ed *ExtractedData) SetProcessingTime(timeMs uint64) {
	ed.ExtractionSummary.SetProcessingTime(timeMs)
}


// OutlineFormatter provides formatting for document outlines
type OutlineFormatter struct{}

// FormatText formats outline as human-readable text
func (of *OutlineFormatter) FormatText(outline *outline.DocumentOutline) string {
	// Simple text formatting implementation
	result := "Document Outline\n"
	result += "================\n\n"

	if outline.DocumentTitle != nil {
		result += "Title: " + *outline.DocumentTitle + "\n"
	}
	result += "Type: " + outline.DocumentType + "\n"
	result += "Total Pages: " + fmt.Sprintf("%d", outline.TotalPages) + "\n\n"

	if outline.HasOutline {
		result += "Table of Contents:\n"
		for _, entry := range outline.Entries {
			result += of.formatOutlineEntry(&entry, 0)
		}
	} else {
		result += "No table of contents available.\n"
	}

	return result
}

// FormatJSON formats outline as JSON
func (of *OutlineFormatter) FormatJSON(outline *outline.DocumentOutline) (string, error) {
	data, err := json.MarshalIndent(outline, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (of *OutlineFormatter) formatOutlineEntry(entry *outline.OutlineEntry, depth int) string {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	result := indent + "- " + entry.Title
	if entry.Page != nil {
		result += fmt.Sprintf(" (Page %d)", *entry.Page)
	}
	result += "\n"

	for _, child := range entry.Children {
		result += of.formatOutlineEntry(&child, depth+1)
	}

	return result
}

// MetadataFormatter provides formatting for document metadata
type MetadataFormatter struct{}

// FormatText formats metadata as human-readable text
func (mf *MetadataFormatter) FormatText(metadata *metadata.FileMetadata) string {
	result := "Document Metadata\n"
	result += "=================\n\n"

	result += "File: " + metadata.FilePath + "\n"
	result += "Type: " + metadata.DocumentType + "\n"
	result += fmt.Sprintf("Size: %d bytes\n", metadata.FileSizeBytes)

	if metadata.Title != nil {
		result += "Title: " + *metadata.Title + "\n"
	}

	if len(metadata.Authors) > 0 {
		result += "Authors: " + strings.Join(metadata.Authors, ", ") + "\n"
	}

	if metadata.CreationDate != nil {
		result += "Created: " + metadata.CreationDate.Format("2006-01-02 15:04:05") + "\n"
	}

	if metadata.PageCount != nil {
		result += fmt.Sprintf("Pages: %d\n", *metadata.PageCount)
	}

	if metadata.WordCount != nil {
		result += fmt.Sprintf("Words: %d\n", *metadata.WordCount)
	}

	return result
}

// FormatJSON formats metadata as JSON
func (mf *MetadataFormatter) FormatJSON(metadata *metadata.FileMetadata) (string, error) {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

