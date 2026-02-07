// Package sdk provides standard cap definitions with arguments
package sdk

import (
	"encoding/json"
	"fmt"

	capns "github.com/macina/cap-sdk-go"
)

// Spec ID constants
const (
	SpecIdStr    = "media:string"
	SpecIdInt    = "media:integer"
	SpecIdNum    = "media:number"
	SpecIdBool   = "media:boolean"
	SpecIdObj    = "media:object"
	SpecIdBinary = "media:binary"

	// Output spec IDs for PDF document processing (with full schemas)
	SpecIdFileMetadata    = "media:file-metadata"
	SpecIdDocumentOutline = "media:document-outline"
	SpecIdDisboundPage    = "media:disbound-page"
)

// InputSpecIdForExt returns the input spec ID for a given file extension
// - PDF files: media:bytes
// - Text files (md, rst, log, txt): media:string
func InputSpecIdForExt(ext string) string {
	if ext == "pdf" {
		return SpecIdBinary
	}
	return SpecIdStr
}

// ExtractMetadataOutputSpecIdForExt returns the output spec ID for extract-metadata by extension
// - PDF files: media:extract-metadata-output (has full schema)
// - Text files: media:object (generic JSON object)
func ExtractMetadataOutputSpecIdForExt(ext string) string {
	if ext == "pdf" {
		return SpecIdFileMetadata
	}
	return SpecIdObj
}

// ExtractOutlineOutputSpecIdForExt returns the output spec ID for extract-outline by extension
// - PDF files: media:extract-outline-output (has full schema)
// - Text files: media:object (generic JSON object)
func ExtractOutlineOutputSpecIdForExt(ext string) string {
	if ext == "pdf" {
		return SpecIdDocumentOutline
	}
	return SpecIdObj
}

// DisboundPageSpecIdForExt returns the output spec ID for disbind by extension
// - PDF files: media:disbound-page (has full schema, output is array)
// - Text files: media:object (generic JSON object)
func DisboundPageSpecIdForExt(ext string) string {
	if ext == "pdf" {
		return SpecIdDisboundPage
	}
	return SpecIdObj
}

// ExtractMetadataUrn builds the URN for extract-metadata capability with given extension
func ExtractMetadataUrn(ext string) string {
	inSpec := InputSpecIdForExt(ext)
	outSpec := ExtractMetadataOutputSpecIdForExt(ext)
	return fmt.Sprintf("cap:ext=%s;in=%s;op=extract_metadata;out=%s", ext, inSpec, outSpec)
}

// GenerateThumbnailUrn builds the URN for generate-thumbnail capability with given extension
func GenerateThumbnailUrn(ext string) string {
	inSpec := InputSpecIdForExt(ext)
	return fmt.Sprintf("cap:ext=%s;in=%s;op=generate_thumbnail;out=media:binary", ext, inSpec)
}

// ExtractOutlineUrn builds the URN for extract-outline capability with given extension
func ExtractOutlineUrn(ext string) string {
	inSpec := InputSpecIdForExt(ext)
	outSpec := ExtractOutlineOutputSpecIdForExt(ext)
	return fmt.Sprintf("cap:ext=%s;in=%s;op=extract_outline;out=%s", ext, inSpec, outSpec)
}

// GrindUrn builds the URN for grind capability with given extension
func GrindUrn(ext string) string {
	inSpec := InputSpecIdForExt(ext)
	outSpec := DisboundPageSpecIdForExt(ext)
	return fmt.Sprintf("cap:ext=%s;in=%s;op=grind;out=%s", ext, inSpec, outSpec)
}

// ExtractMetadataCap creates the standard extract-metadata cap with full argument definition
// Note: This creates a generic cap without extension-specific URN. Use ExtractMetadataUrn for extension-specific URNs.
func ExtractMetadataCap() *capns.Cap {
	id, _ := capns.NewCapUrnFromString("cap:op=extract_metadata")
	
	command := "extract-metadata"
	
	arguments := capns.NewCapArguments()
	
	// Required file_path argument
	filePathValidation := &capns.MediaValidation{
		Pattern:   stringPtr("^[^\\0]+$"),
		MinLength: intPtr(1),
	}
	filePathArg := capns.CapArgument{
		Name:        "file_path",
		ArgType:        capns.ArgumentTypeString,
		ArgDescription: "Path to the document file to process",
		CliFlag:     "file_path",
		Position:    intPtr(0),
		Validation:  filePathValidation,
	}
	arguments.AddRequired(filePathArg)
	
	// Optional output argument
	outputValidation := &capns.MediaValidation{
		Pattern: stringPtr("^[^\\0]+$"),
	}
	outputArg := capns.CapArgument{
		Name:        "output",
		ArgType:        capns.ArgumentTypeString,
		ArgDescription: "Write output to specified file instead of stdout",
		CliFlag:     "--output",
		Validation:  outputValidation,
	}
	arguments.AddOptional(outputArg)
	
	output := &capns.CapOutput{
		OutputType:        capns.OutputTypeObject,
		SchemaRef:   stringPtr("file-metadata.json"),
		ContentType: stringPtr("application/json"),
		Validation:  &capns.MediaValidation{},
		OutputDescription: "Structured metadata including file properties, document properties, and format-specific metadata",
	}
	
	cap := capns.NewCapWithDescription(
		id,
		"Extract Document Metadata",
		command,
		"Extract document metadata including title, author, creation date, file size, and other properties",
	)
	cap.SetArguments(arguments)
	cap.SetOutput(output)
	
	// Metadata extraction can accept stdin for direct file content processing
	cap.AcceptsStdin = true
	
	return cap
}

// GenerateThumbnailCap creates the standard generate-thumbnail cap with full argument definition
// Note: This creates a generic cap without extension-specific URN. Use GenerateThumbnailUrn for extension-specific URNs.
func GenerateThumbnailCap() *capns.Cap {
	id, _ := capns.NewCapUrnFromString("cap:op=generate_thumbnail")
	
	command := "generate-thumbnail"
	
	arguments := capns.NewCapArguments()
	
	// Required file_path argument
	filePathValidation := &capns.MediaValidation{
		Pattern:   stringPtr("^[^\\0]+$"),
		MinLength: intPtr(1),
	}
	filePathArg := capns.CapArgument{
		Name:        "file_path",
		ArgType:        capns.ArgumentTypeString,
		ArgDescription: "Path to the document file to process",
		CliFlag:     "file_path",
		Position:    intPtr(0),
		Validation:  filePathValidation,
	}
	arguments.AddRequired(filePathArg)
	
	// Optional width argument
	widthValidation := &capns.MediaValidation{
		Min: float64Ptr(50.0),
		Max: float64Ptr(2000.0),
	}
	widthDefault := json.Number("200")
	widthArg := capns.CapArgument{
		Name:        "width",
		ArgType:        capns.ArgumentTypeInteger,
		ArgDescription: "Width of the thumbnail in pixels",
		CliFlag:     "--width",
		Validation:  widthValidation,
		DefaultValue:     widthDefault,
	}
	arguments.AddOptional(widthArg)
	
	// Optional height argument
	heightValidation := &capns.MediaValidation{
		Min: float64Ptr(50.0),
		Max: float64Ptr(2000.0),
	}
	heightDefault := json.Number("300")
	heightArg := capns.CapArgument{
		Name:        "height",
		ArgType:        capns.ArgumentTypeInteger,
		ArgDescription: "Height of the thumbnail in pixels",
		CliFlag:     "--height",
		Validation:  heightValidation,
		DefaultValue:     heightDefault,
	}
	arguments.AddOptional(heightArg)
	
	// Optional output argument
	outputValidation := &capns.MediaValidation{
		Pattern: stringPtr("\\.(png|jpg|jpeg)$"),
	}
	outputArg := capns.CapArgument{
		Name:        "output",
		ArgType:        capns.ArgumentTypeString,
		ArgDescription: "Write thumbnail to specified file instead of stdout",
		CliFlag:     "--output",
		Validation:  outputValidation,
	}
	arguments.AddOptional(outputArg)
	
	// Optional page argument
	pageValidation := &capns.MediaValidation{
		Min: float64Ptr(1.0),
	}
	pageDefault := json.Number("1")
	pageArg := capns.CapArgument{
		Name:        "page",
		ArgType:        capns.ArgumentTypeInteger,
		ArgDescription: "Page number to generate thumbnail from (1-based, default: 1)",
		CliFlag:     "--page",
		Validation:  pageValidation,
		DefaultValue:     pageDefault,
	}
	arguments.AddOptional(pageArg)
	
	output := &capns.CapOutput{
		OutputType:        capns.OutputTypeBinary,
		ContentType: stringPtr("image/png"),
		Validation:  &capns.MediaValidation{},
		OutputDescription: "PNG image data representing a thumbnail of the document",
	}
	
	cap := capns.NewCapWithDescription(
		id,
		"Generate Thumbnail",
		command,
		"Generate a thumbnail image preview of the document",
	)
	cap.SetArguments(arguments)
	cap.SetOutput(output)
	
	// Thumbnail generation can accept stdin for direct file content processing
	cap.AcceptsStdin = true
	
	return cap
}

// ExtractOutlineCap creates the standard extract-outline cap with full argument definition
// Note: This creates a generic cap without extension-specific URN. Use ExtractOutlineUrn for extension-specific URNs.
func ExtractOutlineCap() *capns.Cap {
	id, _ := capns.NewCapUrnFromString("cap:op=extract_outline")
	
	command := "extract-outline"
	
	arguments := capns.NewCapArguments()
	
	// Required file_path argument
	filePathValidation := &capns.MediaValidation{
		Pattern:   stringPtr("^[^\\0]+$"),
		MinLength: intPtr(1),
	}
	filePathArg := capns.CapArgument{
		Name:        "file_path",
		ArgType:        capns.ArgumentTypeString,
		ArgDescription: "Path to the document file to process",
		CliFlag:     "file_path",
		Position:    intPtr(0),
		Validation:  filePathValidation,
	}
	arguments.AddRequired(filePathArg)
	
	// Optional max_depth argument
	maxDepthValidation := &capns.MediaValidation{
		Min: float64Ptr(1.0),
		Max: float64Ptr(10.0),
	}
	maxDepthArg := capns.CapArgument{
		Name:        "max_depth",
		ArgType:        capns.ArgumentTypeInteger,
		ArgDescription: "Maximum outline depth to extract (1-10)",
		CliFlag:     "--max-depth",
		Validation:  maxDepthValidation,
	}
	arguments.AddOptional(maxDepthArg)
	
	// Optional include_order_indexes argument
	includeOrderIndexesArg := capns.CapArgument{
		Name:        "include_order_indexes",
		ArgType:        capns.ArgumentTypeBoolean,
		ArgDescription: "Include page numbers in the outline (default: true)",
		CliFlag:     "--include-order-indexes",
		Validation:  &capns.MediaValidation{},
		DefaultValue:     true,
	}
	arguments.AddOptional(includeOrderIndexesArg)
	
	// Optional output argument
	outputValidation := &capns.MediaValidation{
		Pattern: stringPtr("^[^\\0]+$"),
	}
	outputArg := capns.CapArgument{
		Name:        "output",
		ArgType:        capns.ArgumentTypeString,
		ArgDescription: "Write output to specified file instead of stdout",
		CliFlag:     "--output",
		Validation:  outputValidation,
	}
	arguments.AddOptional(outputArg)
	
	output := &capns.CapOutput{
		OutputType:        capns.OutputTypeObject,
		SchemaRef:   stringPtr("document-outline.json"),
		ContentType: stringPtr("application/json"),
		Validation:  &capns.MediaValidation{},
		OutputDescription: "Hierarchical document outline with section titles and optional page numbers",
	}
	
	cap := capns.NewCapWithDescription(
		id,
		"Extract Document Outline",
		command,
		"Extract document outline/table of contents with hierarchical structure",
	)
	cap.SetArguments(arguments)
	cap.SetOutput(output)
	
	// Outline extraction can accept stdin for direct file content processing
	cap.AcceptsStdin = true
	
	return cap
}

// DisbindCap creates the standard grind cap with full argument definition
// Note: This creates a generic cap without extension-specific URN. Use GrindUrn for extension-specific URNs.
func DisbindCap() *capns.Cap {
	id, _ := capns.NewCapUrnFromString("cap:op=grind")
	
	command := "grind"
	
	arguments := capns.NewCapArguments()
	
	// Required file_path argument
	filePathValidation := &capns.MediaValidation{
		Pattern:   stringPtr("^[^\\0]+$"),
		MinLength: intPtr(1),
	}
	filePathArg := capns.CapArgument{
		Name:        "file_path",
		ArgType:        capns.ArgumentTypeString,
		ArgDescription: "Path to the document file to process",
		CliFlag:     "file_path",
		Position:    intPtr(0),
		Validation:  filePathValidation,
	}
	arguments.AddRequired(filePathArg)
	
	// Optional output argument
	outputValidation := &capns.MediaValidation{
		Pattern: stringPtr("^[^\\0]+$"),
	}
	outputArg := capns.CapArgument{
		Name:        "output",
		ArgType:        capns.ArgumentTypeString,
		ArgDescription: "Write output to specified file instead of stdout",
		CliFlag:     "--output",
		Validation:  outputValidation,
	}
	arguments.AddOptional(outputArg)
	
	// Optional index_range argument
	indexRangeValidation := &capns.MediaValidation{
		Pattern: stringPtr("^\\d+(-\\d*)?$"),
	}
	indexRangeArg := capns.CapArgument{
		Name:        "index_range",
		ArgType:        capns.ArgumentTypeString,
		ArgDescription: "Index Range to extract (e.g., '1-5' or '10-')",
		CliFlag:     "--index-range",
		Validation:  indexRangeValidation,
	}
	arguments.AddOptional(indexRangeArg)
	
	output := &capns.CapOutput{
		OutputType:        capns.OutputTypeArray,
		SchemaRef:         stringPtr("disbound-page.json"),
		ContentType:       stringPtr("application/json"),
		Validation:        &capns.MediaValidation{},
		OutputDescription: "Array of structured page content extracted from the document",
	}
	
	cap := capns.NewCapWithDescription(
		id,
		"Extract File Chips",
		command,
		"Extract structured page content from the document",
	)
	cap.SetArguments(arguments)
	cap.SetOutput(output)
	
	// Page extraction can accept stdin for direct file content processing
	cap.AcceptsStdin = true
	
	return cap
}

// GetAllStandardCaps returns all standard plugin caps
func GetAllStandardCaps() []*capns.Cap {
	return []*capns.Cap{
		ExtractMetadataCap(),
		GenerateThumbnailCap(),
		ExtractOutlineCap(),
		DisbindCap(),
	}
}

// GetStandardCap returns a standard cap by name
func GetStandardCap(name string) *capns.Cap {
	switch name {
	case "extract-metadata":
		return ExtractMetadataCap()
	case "generate-thumbnail":
		return GenerateThumbnailCap()
	case "extract-outline":
		return ExtractOutlineCap()
	case "grind":
		return DisbindCap()
	default:
		return nil
	}
}

// GetStandardCapByUrn returns a standard cap by cap URN string
// Note: This matches generic URN format (without extension). For extension-specific URNs, use the URN builder functions.
func GetStandardCapByUrn(urnStr string) *capns.Cap {
	switch urnStr {
	case "cap:op=extract_metadata":
		return ExtractMetadataCap()
	case "cap:op=generate_thumbnail":
		return GenerateThumbnailCap()
	case "cap:op=extract_outline":
		return ExtractOutlineCap()
	case "cap:op=grind":
		return DisbindCap()
	default:
		return nil
	}
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}