// Package sdk provides the main LBVR Plugin SDK for Go
package sdk

import (
	"github.com/jowharshamshiri/lbvr-plugin-sdk-go/pkg/handler"
	"github.com/jowharshamshiri/lbvr-plugin-sdk-go/pkg/metadata"
	"github.com/jowharshamshiri/lbvr-plugin-sdk-go/pkg/outline"
	"github.com/jowharshamshiri/lbvr-plugin-sdk-go/pkg/pages"
	"github.com/jowharshamshiri/lbvr-plugin-sdk-go/pkg/plugin"
)

// Version of the LBVR Plugin SDK
const Version = "0.1.0"

// Re-export commonly used types for convenience
type (
	// Core interfaces
	DocumentHandler = handler.DocumentHandler
	PluginMetadata  = handler.PluginMetadata

	// Data structures
	FileMetadata     = metadata.FileMetadata
	DocumentOutline  = outline.DocumentOutline
	DocumentPages    = pages.DocumentPages
	DocumentPage     = pages.DocumentPage
	DocumentParagraph = pages.DocumentParagraph
	TocEntry         = outline.TocEntry
	ExtractionInfo   = outline.ExtractionInfo
	
	// Handler types
	FileInfo            = handler.FileInfo
	QuickMetadata       = handler.QuickMetadata
	PluginCapabilities  = handler.PluginCapabilities
	PluginInfo          = handler.PluginInfo
	ProcessingResult    = handler.ProcessingResult
	HandlerRegistry     = handler.HandlerRegistry

	// Output types
	ExtractedData      = plugin.ExtractedData
	ThumbnailInfo      = plugin.ThumbnailInfo
	ExtractionSummary  = plugin.ExtractionSummary
	OutputFormat       = plugin.OutputFormat
	OutlineFormatter   = plugin.OutlineFormatter
	MetadataFormatter  = plugin.MetadataFormatter
)

// Constructor functions for convenience
var (
	// Metadata constructors
	NewFileMetadata        = metadata.NewFileMetadata
	NewMinimalFileMetadata = metadata.NewMinimalFileMetadata

	// Outline constructors
	NewDocumentOutline = outline.NewDocumentOutline
	NewTocEntry        = outline.NewTocEntry
	NewExtractionInfo  = outline.NewExtractionInfo

	// Pages constructors
	NewDocumentPages         = pages.NewDocumentPages
	NewDocumentPage          = pages.NewDocumentPage
	NewDocumentPageWithText  = pages.NewDocumentPageWithText
	NewDocumentParagraph     = pages.NewDocumentParagraph

	// Handler constructors
	NewHandlerRegistry = handler.NewHandlerRegistry
	NewSuccessResult   = handler.NewSuccessResult
	NewFailureResult   = handler.NewFailureResult

	// Output constructors
	NewExtractedData     = plugin.NewExtractedData
	NewExtractionSummary = plugin.NewExtractionSummary
)