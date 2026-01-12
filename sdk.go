// Package sdk provides the main FGND Plugin SDK for Go
package sdk

import (
	"github.com/jowharshamshiri/fgnd-plugin-sdk-go/pkg/handler"
	"github.com/jowharshamshiri/fgnd-plugin-sdk-go/pkg/metadata"
	"github.com/jowharshamshiri/fgnd-plugin-sdk-go/pkg/outline"
	"github.com/jowharshamshiri/fgnd-plugin-sdk-go/pkg/pages"
	"github.com/jowharshamshiri/fgnd-plugin-sdk-go/pkg/plugin"
	
	// Re-export cap SDK types
	capns "github.com/fgnd/cap-sdk-go"
)

// Version of the FGND Plugin SDK
const Version = "0.1.0"

// Re-export argument type constants
const (
	ArgumentTypeString  = capns.ArgumentTypeString
	ArgumentTypeInteger = capns.ArgumentTypeInteger
	ArgumentTypeNumber  = capns.ArgumentTypeNumber
	ArgumentTypeBoolean = capns.ArgumentTypeBoolean
	ArgumentTypeArray   = capns.ArgumentTypeArray
	ArgumentTypeObject  = capns.ArgumentTypeObject
	ArgumentTypeBinary  = capns.ArgumentTypeBinary
)

// Re-export output type constants
const (
	OutputTypeString  = capns.OutputTypeString
	OutputTypeInteger = capns.OutputTypeInteger
	OutputTypeNumber  = capns.OutputTypeNumber
	OutputTypeBoolean = capns.OutputTypeBoolean
	OutputTypeArray   = capns.OutputTypeArray
	OutputTypeObject  = capns.OutputTypeObject
	OutputTypeBinary  = capns.OutputTypeBinary
)

// Re-export commonly used types for convenience
type (
	// Core interfaces
	DocumentHandler = handler.DocumentHandler

	// Data structures
	FileMetadata     = metadata.FileMetadata
	DocumentOutline  = outline.DocumentOutline
	GroundChips    = pages.GroundChips
	FileChip     = pages.FileChip
	OutlineEntry     = outline.OutlineEntry
	ExtractionInfo   = outline.ExtractionInfo
	
	// Handler types
	FileInfo            = handler.FileInfo
	QuickMetadata       = handler.QuickMetadata
	ProcessingResult    = handler.ProcessingResult
	HandlerRegistry     = handler.HandlerRegistry
	PluginManifest      = capns.CapManifest
	
	// Cap types from cap SDK
	CapUrn         = capns.CapUrn
	Cap           = capns.Cap
	CapMatcher    = capns.CapMatcher
	CapCaller     = capns.CapCaller
	ResponseWrapper = capns.ResponseWrapper
	CapSet       = capns.CapSet
	HostResult    = capns.HostResult
	
	// Argument and output types
	ArgumentType = capns.ArgumentType
	OutputType   = capns.OutputType
	
	// Schema validation types
	SchemaValidator            = capns.SchemaValidator
	SchemaValidationError      = capns.SchemaValidationError
	SchemaResolver             = capns.SchemaResolver
	FileSchemaResolver         = capns.FileSchemaResolver
	CapValidationCoordinator   = capns.CapValidationCoordinator
	InputValidator             = capns.InputValidator
	OutputValidator            = capns.OutputValidator

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
	NewOutlineEntry        = outline.NewOutlineEntry
	NewExtractionInfo  = outline.NewExtractionInfo

	// Pages constructors
	NewGroundChips        = pages.NewGroundChips
	NewFileChip         = pages.NewFileChip
	NewFileChipWithText = pages.NewFileChipWithText

	// Handler constructors
	NewHandlerRegistry = handler.NewHandlerRegistry
	NewSuccessResult   = handler.NewSuccessResult
	NewFailureResult   = handler.NewFailureResult

	// Output constructors
	NewExtractedData     = plugin.NewExtractedData
	NewExtractionSummary = plugin.NewExtractionSummary
	
	// Cap constructors
	NewCapUrnFromString           = capns.NewCapUrnFromString
	NewCapUrnFromTags             = capns.NewCapUrnFromTags
	NewCapUrnBuilder              = capns.NewCapUrnBuilder
	NewCap                       = capns.NewCap
	NewCapWithDescription        = capns.NewCapWithDescription
	NewCapWithMetadata           = capns.NewCapWithMetadata
	NewCapWithDescriptionAndMetadata = capns.NewCapWithDescriptionAndMetadata
	NewPluginManifest                       = capns.NewCapManifest
	
	// Schema-enabled constructors
	NewCapArgumentWithSchema      = capns.NewCapArgumentWithSchema
	NewCapArgumentWithSchemaRef   = capns.NewCapArgumentWithSchemaRef
	NewCapOutputWithSchemaRef     = capns.NewCapOutputWithSchemaRef
	NewCapOutputWithEmbeddedSchema = capns.NewCapOutputWithEmbeddedSchema
	
	// Validation constructors
	NewSchemaValidator                     = capns.NewSchemaValidator
	NewSchemaValidatorWithResolver         = capns.NewSchemaValidatorWithResolver
	NewFileSchemaResolver                  = capns.NewFileSchemaResolver
	NewCapValidationCoordinator            = capns.NewCapValidationCoordinator
	NewCapValidationCoordinatorWithSchemaResolver = capns.NewCapValidationCoordinatorWithSchemaResolver
	NewInputValidator                      = capns.NewInputValidator
	NewInputValidatorWithSchemaResolver    = capns.NewInputValidatorWithSchemaResolver
	NewOutputValidator                     = capns.NewOutputValidator
	NewOutputValidatorWithSchemaResolver   = capns.NewOutputValidatorWithSchemaResolver
	
	// Caller/Response constructors
	NewCapCaller               = capns.NewCapCaller
	NewResponseWrapperFromJSON = capns.NewResponseWrapperFromJSON
	NewResponseWrapperFromText = capns.NewResponseWrapperFromText
	NewResponseWrapperFromBinary = capns.NewResponseWrapperFromBinary
	
)