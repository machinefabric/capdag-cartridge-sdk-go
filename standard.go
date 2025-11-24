// Package sdk provides standard cap definitions with arguments
package sdk

import (
	"encoding/json"
	
	capdef "github.com/fmio/capdef-go"
)

// ExtractMetadataCap creates the standard extract-metadata cap with full argument definition
func ExtractMetadataCap() *capdef.Cap {
	id, _ := capdef.NewCapCardFromString("action=extract;target=metadata;")
	
	command := "extract-metadata"
	
	arguments := capdef.NewCapArguments()
	
	// Required file_path argument
	filePathValidation := &capdef.ArgumentValidation{
		Pattern:   stringPtr("^[^\\0]+$"),
		MinLength: intPtr(1),
	}
	filePathArg := capdef.CapArgument{
		Name:        "file_path",
		Type:        capdef.ArgumentTypeString,
		Description: "Path to the document file to process",
		CliFlag:     "file_path",
		Position:    intPtr(0),
		Validation:  filePathValidation,
	}
	arguments.AddRequired(filePathArg)
	
	// Optional output argument
	outputValidation := &capdef.ArgumentValidation{
		Pattern: stringPtr("^[^\\0]+$"),
	}
	outputArg := capdef.CapArgument{
		Name:        "output",
		Type:        capdef.ArgumentTypeString,
		Description: "Write output to specified file instead of stdout",
		CliFlag:     "--output",
		Validation:  outputValidation,
	}
	arguments.AddOptional(outputArg)
	
	output := &capdef.CapOutput{
		Type:        capdef.OutputTypeObject,
		SchemaRef:   stringPtr("file-metadata.json"),
		ContentType: stringPtr("application/json"),
		Validation:  &capdef.ArgumentValidation{},
		Description: "Structured metadata including file properties, document properties, and format-specific metadata",
	}
	
	cap := capdef.NewCapWithDescription(
		id,
		"1.0.0",
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
func GenerateThumbnailCap() *capdef.Cap {
	id, _ := capdef.NewCapCardFromString("action=generate;output=binary;target=thumbnail;")
	
	command := "generate-thumbnail"
	
	arguments := capdef.NewCapArguments()
	
	// Required file_path argument
	filePathValidation := &capdef.ArgumentValidation{
		Pattern:   stringPtr("^[^\\0]+$"),
		MinLength: intPtr(1),
	}
	filePathArg := capdef.CapArgument{
		Name:        "file_path",
		Type:        capdef.ArgumentTypeString,
		Description: "Path to the document file to process",
		CliFlag:     "file_path",
		Position:    intPtr(0),
		Validation:  filePathValidation,
	}
	arguments.AddRequired(filePathArg)
	
	// Optional width argument
	widthValidation := &capdef.ArgumentValidation{
		Min: float64Ptr(50.0),
		Max: float64Ptr(2000.0),
	}
	widthDefault := json.Number("200")
	widthArg := capdef.CapArgument{
		Name:        "width",
		Type:        capdef.ArgumentTypeInteger,
		Description: "Width of the thumbnail in pixels",
		CliFlag:     "--width",
		Validation:  widthValidation,
		Default:     widthDefault,
	}
	arguments.AddOptional(widthArg)
	
	// Optional height argument
	heightValidation := &capdef.ArgumentValidation{
		Min: float64Ptr(50.0),
		Max: float64Ptr(2000.0),
	}
	heightDefault := json.Number("300")
	heightArg := capdef.CapArgument{
		Name:        "height",
		Type:        capdef.ArgumentTypeInteger,
		Description: "Height of the thumbnail in pixels",
		CliFlag:     "--height",
		Validation:  heightValidation,
		Default:     heightDefault,
	}
	arguments.AddOptional(heightArg)
	
	// Optional output argument
	outputValidation := &capdef.ArgumentValidation{
		Pattern: stringPtr("\\.(png|jpg|jpeg)$"),
	}
	outputArg := capdef.CapArgument{
		Name:        "output",
		Type:        capdef.ArgumentTypeString,
		Description: "Write thumbnail to specified file instead of stdout",
		CliFlag:     "--output",
		Validation:  outputValidation,
	}
	arguments.AddOptional(outputArg)
	
	// Optional page argument
	pageValidation := &capdef.ArgumentValidation{
		Min: float64Ptr(1.0),
	}
	pageDefault := json.Number("1")
	pageArg := capdef.CapArgument{
		Name:        "page",
		Type:        capdef.ArgumentTypeInteger,
		Description: "Page number to generate thumbnail from (1-based, default: 1)",
		CliFlag:     "--page",
		Validation:  pageValidation,
		Default:     pageDefault,
	}
	arguments.AddOptional(pageArg)
	
	output := &capdef.CapOutput{
		Type:        capdef.OutputTypeBinary,
		ContentType: stringPtr("image/png"),
		Validation:  &capdef.ArgumentValidation{},
		Description: "PNG image data representing a thumbnail of the document",
	}
	
	cap := capdef.NewCapWithDescription(
		id,
		"1.0.0",
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
func ExtractOutlineCap() *capdef.Cap {
	id, _ := capdef.NewCapCardFromString("action=extract;target=outline;")
	
	command := "extract-outline"
	
	arguments := capdef.NewCapArguments()
	
	// Required file_path argument
	filePathValidation := &capdef.ArgumentValidation{
		Pattern:   stringPtr("^[^\\0]+$"),
		MinLength: intPtr(1),
	}
	filePathArg := capdef.CapArgument{
		Name:        "file_path",
		Type:        capdef.ArgumentTypeString,
		Description: "Path to the document file to process",
		CliFlag:     "file_path",
		Position:    intPtr(0),
		Validation:  filePathValidation,
	}
	arguments.AddRequired(filePathArg)
	
	// Optional max_depth argument
	maxDepthValidation := &capdef.ArgumentValidation{
		Min: float64Ptr(1.0),
		Max: float64Ptr(10.0),
	}
	maxDepthArg := capdef.CapArgument{
		Name:        "max_depth",
		Type:        capdef.ArgumentTypeInteger,
		Description: "Maximum outline depth to extract (1-10)",
		CliFlag:     "--max-depth",
		Validation:  maxDepthValidation,
	}
	arguments.AddOptional(maxDepthArg)
	
	// Optional include_page_numbers argument
	includePageNumbersArg := capdef.CapArgument{
		Name:        "include_page_numbers",
		Type:        capdef.ArgumentTypeBoolean,
		Description: "Include page numbers in the outline (default: true)",
		CliFlag:     "--include-page-numbers",
		Validation:  &capdef.ArgumentValidation{},
		Default:     true,
	}
	arguments.AddOptional(includePageNumbersArg)
	
	// Optional output argument
	outputValidation := &capdef.ArgumentValidation{
		Pattern: stringPtr("^[^\\0]+$"),
	}
	outputArg := capdef.CapArgument{
		Name:        "output",
		Type:        capdef.ArgumentTypeString,
		Description: "Write output to specified file instead of stdout",
		CliFlag:     "--output",
		Validation:  outputValidation,
	}
	arguments.AddOptional(outputArg)
	
	output := &capdef.CapOutput{
		Type:        capdef.OutputTypeObject,
		SchemaRef:   stringPtr("document-outline.json"),
		ContentType: stringPtr("application/json"),
		Validation:  &capdef.ArgumentValidation{},
		Description: "Hierarchical document outline with section titles and optional page numbers",
	}
	
	cap := capdef.NewCapWithDescription(
		id,
		"1.0.0",
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
func ExtractPagesCap() *capdef.Cap {
	id, _ := capdef.NewCapCardFromString("action=extract;target=pages;")
	
	command := "extract-pages"
	
	arguments := capdef.NewCapArguments()
	
	// Required file_path argument
	filePathValidation := &capdef.ArgumentValidation{
		Pattern:   stringPtr("^[^\\0]+$"),
		MinLength: intPtr(1),
	}
	filePathArg := capdef.CapArgument{
		Name:        "file_path",
		Type:        capdef.ArgumentTypeString,
		Description: "Path to the document file to process",
		CliFlag:     "file_path",
		Position:    intPtr(0),
		Validation:  filePathValidation,
	}
	arguments.AddRequired(filePathArg)
	
	// Optional output argument
	outputValidation := &capdef.ArgumentValidation{
		Pattern: stringPtr("^[^\\0]+$"),
	}
	outputArg := capdef.CapArgument{
		Name:        "output",
		Type:        capdef.ArgumentTypeString,
		Description: "Write output to specified file instead of stdout",
		CliFlag:     "--output",
		Validation:  outputValidation,
	}
	arguments.AddOptional(outputArg)
	
	// Optional page_range argument
	pageRangeValidation := &capdef.ArgumentValidation{
		Pattern: stringPtr("^\\d+(-\\d*)?$"),
	}
	pageRangeArg := capdef.CapArgument{
		Name:        "page_range",
		Type:        capdef.ArgumentTypeString,
		Description: "Page range to extract (e.g., '1-5' or '10-')",
		CliFlag:     "--page-range",
		Validation:  pageRangeValidation,
	}
	arguments.AddOptional(pageRangeArg)
	
	output := &capdef.CapOutput{
		Type:        capdef.OutputTypeObject,
		SchemaRef:   stringPtr("document-pages.json"),
		ContentType: stringPtr("application/json"),
		Validation:  &capdef.ArgumentValidation{},
		Description: "Structured page content extracted from the document",
	}
	
	cap := capdef.NewCapWithDescription(
		id,
		"1.0.0",
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
func GetAllStandardCaps() []*capdef.Cap {
	return []*capdef.Cap{
		ExtractMetadataCap(),
		GenerateThumbnailCap(),
		ExtractOutlineCap(),
		ExtractPagesCap(),
	}
}

// GetStandardCap returns a standard cap by name
func GetStandardCap(name string) *capdef.Cap {
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

// GetStandardCapById returns a standard cap by cap ID string
func GetStandardCapById(idStr string) *capdef.Cap {
	switch idStr {
	case "action=extract;target=metadata;":
		return ExtractMetadataCap()
	case "action=generate;output=binary;target=thumbnail;":
		return GenerateThumbnailCap()
	case "action=extract;target=outline;":
		return ExtractOutlineCap()
	case "action=extract;target=pages;":
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