// Package sdk provides the main MACINA Plugin SDK for Go
package sdk

import (
	"github.com/machinefabric/fgrnd-plugin-sdk-go/pkg/handler"
	"github.com/machinefabric/fgrnd-plugin-sdk-go/pkg/metadata"
	"github.com/machinefabric/fgrnd-plugin-sdk-go/pkg/outline"
	"github.com/machinefabric/fgrnd-plugin-sdk-go/pkg/pages"
	"github.com/machinefabric/fgrnd-plugin-sdk-go/pkg/plugin"

	"github.com/machinefabric/capdag-go/cap"
	"github.com/machinefabric/capdag-go/urn"
)

// Version of the MACINA Plugin SDK
const Version = "0.2.0"

// Re-export commonly used types for convenience
type (
	// Core interfaces
	DocumentHandler = handler.DocumentHandler

	// Data structures
	FileMetadata    = metadata.FileMetadata
	DocumentOutline = outline.DocumentOutline
	DisboundPage    = pages.DisboundPage
	OutlineEntry    = outline.OutlineEntry
	ExtractionInfo  = outline.ExtractionInfo

	// Handler types
	FileInfo         = handler.FileInfo
	QuickMetadata    = handler.QuickMetadata
	ProcessingResult = handler.ProcessingResult
	HandlerRegistry  = handler.HandlerRegistry

	// Cap types from cap SDK
	CapUrn          = urn.CapUrn
	CapMatcher      = urn.CapMatcher
	Cap             = cap.Cap
	CapArg          = cap.CapArg
	ArgSource       = cap.ArgSource
	CapOutput       = cap.CapOutput
	CapCaller       = cap.CapCaller
	ResponseWrapper = cap.ResponseWrapper
	CapSet          = cap.CapSet
	HostResult      = cap.HostResult

	// Argument value types
	CapArgumentValue = cap.CapArgumentValue

	// Schema validation types
	SchemaValidator          = cap.SchemaValidator
	SchemaValidationError    = cap.SchemaValidationError
	SchemaResolver           = cap.SchemaResolver
	FileSchemaResolver       = cap.FileSchemaResolver
	CapValidationCoordinator = cap.CapValidationCoordinator
	InputValidator           = cap.InputValidator
	OutputValidator          = cap.OutputValidator

	// Output types
	ExtractedData     = plugin.ExtractedData
	ThumbnailInfo     = plugin.ThumbnailInfo
	ExtractionSummary = plugin.ExtractionSummary
	OutputFormat      = plugin.OutputFormat
	OutlineFormatter  = plugin.OutlineFormatter
	MetadataFormatter = plugin.MetadataFormatter
)

// Constructor functions for convenience
var (
	// Metadata constructors
	NewFileMetadata        = metadata.NewFileMetadata
	NewMinimalFileMetadata = metadata.NewMinimalFileMetadata

	// Outline constructors
	NewDocumentOutline = outline.NewDocumentOutline
	NewOutlineEntry    = outline.NewOutlineEntry
	NewExtractionInfo  = outline.NewExtractionInfo

	// Pages constructors
	NewDisboundPage         = pages.NewDisboundPage
	NewDisboundPageWithText = pages.NewDisboundPageWithText

	// Handler constructors
	NewHandlerRegistry = handler.NewHandlerRegistry
	NewSuccessResult   = handler.NewSuccessResult
	NewFailureResult   = handler.NewFailureResult

	// Output constructors
	NewExtractedData     = plugin.NewExtractedData
	NewExtractionSummary = plugin.NewExtractionSummary

	// URN constructors
	NewCapUrnFromString = urn.NewCapUrnFromString
	NewCapUrnFromTags   = urn.NewCapUrnFromTags
	NewCapUrnBuilder    = urn.NewCapUrnBuilder

	// Cap constructors
	NewCap                = cap.NewCap
	NewCapWithDescription = cap.NewCapWithDescription
	NewCapWithMetadata    = cap.NewCapWithMetadata

	// Argument constructors
	NewCapArg                = cap.NewCapArg
	NewCapArgWithDescription = cap.NewCapArgWithDescription
	NewCapOutput             = cap.NewCapOutput

	// Argument value constructors
	NewCapArgumentValue        = cap.NewCapArgumentValue
	NewCapArgumentValueFromStr = cap.NewCapArgumentValueFromStr

	// Validation constructors
	NewSchemaValidator                                = cap.NewSchemaValidator
	NewSchemaValidatorWithResolver                    = cap.NewSchemaValidatorWithResolver
	NewFileSchemaResolver                             = cap.NewFileSchemaResolver
	NewCapValidationCoordinator                       = cap.NewCapValidationCoordinator
	NewCapValidationCoordinatorWithSchemaResolver     = cap.NewCapValidationCoordinatorWithSchemaResolver
	NewInputValidator                                 = cap.NewInputValidator
	NewInputValidatorWithSchemaResolver               = cap.NewInputValidatorWithSchemaResolver
	NewOutputValidator                                = cap.NewOutputValidator
	NewOutputValidatorWithSchemaResolver              = cap.NewOutputValidatorWithSchemaResolver

	// Caller/Response constructors
	NewCapCaller               = cap.NewCapCaller
	NewResponseWrapperFromJSON   = cap.NewResponseWrapperFromJSON
	NewResponseWrapperFromText   = cap.NewResponseWrapperFromText
	NewResponseWrapperFromBinary = cap.NewResponseWrapperFromBinary
)
