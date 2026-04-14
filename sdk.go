// Package sdk provides the main MACINA Cartridge SDK for Go
package sdk

import (
	"github.com/machinefabric/fgrnd-cartridge-sdk-go/pkg/handler"
	"github.com/machinefabric/fgrnd-cartridge-sdk-go/pkg/metadata"
	"github.com/machinefabric/fgrnd-cartridge-sdk-go/pkg/outline"
	"github.com/machinefabric/fgrnd-cartridge-sdk-go/pkg/pages"
	"github.com/machinefabric/fgrnd-cartridge-sdk-go/pkg/cartridge"

	"github.com/machinefabric/capdag-go/cap"
	"github.com/machinefabric/capdag-go/urn"
)

// Version of the MACINA Cartridge SDK
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
	CapResult       = cap.CapResult
	CapResultKind   = cap.CapResultKind

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
	ExtractedData     = cartridge.ExtractedData
	ThumbnailInfo     = cartridge.ThumbnailInfo
	ExtractionSummary = cartridge.ExtractionSummary
	OutputFormat      = cartridge.OutputFormat
	OutlineFormatter  = cartridge.OutlineFormatter
	MetadataFormatter = cartridge.MetadataFormatter
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
	NewExtractedData     = cartridge.NewExtractedData
	NewExtractionSummary = cartridge.NewExtractionSummary

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
