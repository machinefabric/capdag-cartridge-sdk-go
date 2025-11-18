// Package sdk provides the main LBVR Plugin SDK for Go
package sdk

import (
	"github.com/jowharshamshiri/lbvr-plugin-sdk-go/pkg/handler"
	"github.com/jowharshamshiri/lbvr-plugin-sdk-go/pkg/metadata"
	"github.com/jowharshamshiri/lbvr-plugin-sdk-go/pkg/outline"
	"github.com/jowharshamshiri/lbvr-plugin-sdk-go/pkg/pages"
	"github.com/jowharshamshiri/lbvr-plugin-sdk-go/pkg/plugin"
	
	// Re-export cap SDK types
	capdef "github.com/lbvr/capdef-go"
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
	ProcessingResult    = handler.ProcessingResult
	HandlerRegistry     = handler.HandlerRegistry
	PluginManifest          = capdef.CapManifest
	
	// Cap types from cap SDK
	CapCard        = capdef.CapCard
	Cap          = capdef.Cap
	CapMatcher   = capdef.CapMatcher

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
	
	// Cap constructors
	NewCapCardFromString           = capdef.NewCapCardFromString
	NewCapCardFromTags             = capdef.NewCapCardFromTags
	NewCapCardBuilder              = capdef.NewCapCardBuilder
	NewCap                       = capdef.NewCap
	NewCapWithDescription        = capdef.NewCapWithDescription
	NewCapWithMetadata           = capdef.NewCapWithMetadata
	NewCapWithDescriptionAndMetadata = capdef.NewCapWithDescriptionAndMetadata
	NewPluginManifest                       = capdef.NewCapManifest
)