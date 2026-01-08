// Package sdk provides standard cap definitions with arguments
package sdk

import (
	"encoding/json"
	"fmt"

	capns "github.com/fgrnd/cap-sdk-go"
)

// Spec ID constants
const (
	SpecIdStr    = "std:str.v1"
	SpecIdInt    = "std:int.v1"
	SpecIdNum    = "std:num.v1"
	SpecIdBool   = "std:bool.v1"
	SpecIdObj    = "std:obj.v1"
	SpecIdBinary = "std:binary.v1"

	// Output spec IDs for PDF document processing (with full schemas)
	SpecIdExtractMetadataOutput = "capns:extract-metadata-output.v1"
	SpecIdExtractOutlineOutput  = "capns:extract-outline-output.v1"
	SpecIdExtractPagesOutput    = "capns:extract-pages-output.v1"
)

// InputSpecIdForExt returns the input spec ID for a given file extension
// - PDF files: std:binary.v1
// - Text files (md, rst, log, txt): std:str.v1
func InputSpecIdForExt(ext string) string {
	if ext == "pdf" {
		return SpecIdBinary
	}
	return SpecIdStr
}

// ExtractMetadataOutputSpecIdForExt returns the output spec ID for extract-metadata by extension
// - PDF files: capns:extract-metadata-output.v1 (has full schema)
// - Text files: std:obj.v1 (generic JSON object)
func ExtractMetadataOutputSpecIdForExt(ext string) string {
	if ext == "pdf" {
		return SpecIdExtractMetadataOutput
	}
	return SpecIdObj
}

// ExtractOutlineOutputSpecIdForExt returns the output spec ID for extract-outline by extension
// - PDF files: capns:extract-outline-output.v1 (has full schema)
// - Text files: std:obj.v1 (generic JSON object)
func ExtractOutlineOutputSpecIdForExt(ext string) string {
	if ext == "pdf" {
		return SpecIdExtractOutlineOutput
	}
	return SpecIdObj
}

// ExtractPagesOutputSpecIdForExt returns the output spec ID for extract-pages by extension
// - PDF files: capns:extract-pages-output.v1 (has full schema)
// - Text files: std:obj.v1 (generic JSON object)
func ExtractPagesOutputSpecIdForExt(ext string) string {
	if ext == "pdf" {
		return SpecIdExtractPagesOutput
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
	return fmt.Sprintf("cap:ext=%s;in=%s;op=generate_thumbnail;out=std:binary.v1", ext, inSpec)
}

// ExtractOutlineUrn builds the URN for extract-outline capability with given extension
func ExtractOutlineUrn(ext string) string {
	inSpec := InputSpecIdForExt(ext)
	outSpec := ExtractOutlineOutputSpecIdForExt(ext)
	return fmt.Sprintf("cap:ext=%s;in=%s;op=extract_outline;out=%s", ext, inSpec, outSpec)
}

// ExtractPagesUrn builds the URN for extract-pages capability with given extension
func ExtractPagesUrn(ext string) string {
	inSpec := InputSpecIdForExt(ext)
	outSpec := ExtractPagesOutputSpecIdForExt(ext)
	return fmt.Sprintf("cap:ext=%s;in=%s;op=extract_pages;out=%s", ext, inSpec, outSpec)
}

// ExtractMetadataCap creates the standard extract-metadata cap with full argument definition
// Note: This creates a generic cap without extension-specific URN. Use ExtractMetadataUrn for extension-specific URNs.
func ExtractMetadataCap() *capns.Cap {
	id, _ := capns.NewCapUrnFromString("cap:op=extract_metadata")
	
	command := "extract-metadata"
	
	arguments := capns.NewCapArguments()
	
	// Required file_path argument
	filePathValidation := &capns.ArgumentValidation{
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
	outputValidation := &capns.ArgumentValidation{
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
		Validation:  &capns.ArgumentValidation{},
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
	filePathValidation := &capns.ArgumentValidation{
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
	widthValidation := &capns.ArgumentValidation{
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
	heightValidation := &capns.ArgumentValidation{
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
	outputValidation := &capns.ArgumentValidation{
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
	pageValidation := &capns.ArgumentValidation{
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
		Validation:  &capns.ArgumentValidation{},
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
	filePathValidation := &capns.ArgumentValidation{
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
	maxDepthValidation := &capns.ArgumentValidation{
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
		Validation:  &capns.ArgumentValidation{},
		DefaultValue:     true,
	}
	arguments.AddOptional(includeOrderIndexesArg)
	
	// Optional output argument
	outputValidation := &capns.ArgumentValidation{
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
		Validation:  &capns.ArgumentValidation{},
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

// ExtractPagesCap creates the standard extract-pages cap with full argument definition
// Note: This creates a generic cap without extension-specific URN. Use ExtractPagesUrn for extension-specific URNs.
func ExtractPagesCap() *capns.Cap {
	id, _ := capns.NewCapUrnFromString("cap:op=extract_pages")
	
	command := "extract-pages"
	
	arguments := capns.NewCapArguments()
	
	// Required file_path argument
	filePathValidation := &capns.ArgumentValidation{
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
	outputValidation := &capns.ArgumentValidation{
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
	
	// Optional page_range argument
	pageRangeValidation := &capns.ArgumentValidation{
		Pattern: stringPtr("^\\d+(-\\d*)?$"),
	}
	pageRangeArg := capns.CapArgument{
		Name:        "page_range",
		ArgType:        capns.ArgumentTypeString,
		ArgDescription: "Page range to extract (e.g., '1-5' or '10-')",
		CliFlag:     "--page-range",
		Validation:  pageRangeValidation,
	}
	arguments.AddOptional(pageRangeArg)
	
	output := &capns.CapOutput{
		OutputType:        capns.OutputTypeObject,
		SchemaRef:   stringPtr("document-pages.json"),
		ContentType: stringPtr("application/json"),
		Validation:  &capns.ArgumentValidation{},
		OutputDescription: "Structured page content extracted from the document",
	}
	
	cap := capns.NewCapWithDescription(
		id,
		"Extract Document Pages",
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
		ExtractPagesCap(),
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
	case "extract-pages":
		return ExtractPagesCap()
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
	case "cap:op=extract_pages":
		return ExtractPagesCap()
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