// Package sdk provides standard capability definitions with arguments
package sdk

import (
	"encoding/json"
	
	capdef "github.com/lbvr/capdef-go"
)

// ExtractMetadataCapability creates the standard extract-metadata capability with full argument definition
func ExtractMetadataCapability() *capdef.Capability {
	id, _ := capdef.NewCapabilityIdFromString("document:extract:metadata")
	
	command := "extract-metadata"
	
	arguments := capdef.NewCapabilityArguments()
	
	// Required file_path argument
	filePathValidation := &capdef.ArgumentValidation{
		Pattern:   stringPtr("^[^\\0]+$"),
		MinLength: intPtr(1),
	}
	filePathArg := capdef.CapabilityArgument{
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
	outputArg := capdef.CapabilityArgument{
		Name:        "output",
		Type:        capdef.ArgumentTypeString,
		Description: "Write output to specified file instead of stdout",
		CliFlag:     "--output",
		Validation:  outputValidation,
	}
	arguments.AddOptional(outputArg)
	
	output := &capdef.CapabilityOutput{
		Type:        capdef.OutputTypeObject,
		SchemaRef:   stringPtr("file-metadata.json"),
		ContentType: stringPtr("application/json"),
		Validation:  &capdef.ArgumentValidation{},
		Description: "Structured metadata including file properties, document properties, and format-specific metadata",
	}
	
	capability := capdef.NewCapabilityWithDescription(
		id,
		"1.0.0",
		command,
		"Extract document metadata including title, author, creation date, file size, and other properties",
	)
	capability.SetArguments(arguments)
	capability.SetOutput(output)
	
	return capability
}

// GenerateThumbnailCapability creates the standard generate-thumbnail capability with full argument definition
func GenerateThumbnailCapability() *capdef.Capability {
	id, _ := capdef.NewCapabilityIdFromString("document:generate:thumbnail")
	
	command := "generate-thumbnail"
	
	arguments := capdef.NewCapabilityArguments()
	
	// Required file_path argument
	filePathValidation := &capdef.ArgumentValidation{
		Pattern:   stringPtr("^[^\\0]+$"),
		MinLength: intPtr(1),
	}
	filePathArg := capdef.CapabilityArgument{
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
	widthArg := capdef.CapabilityArgument{
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
	heightArg := capdef.CapabilityArgument{
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
	outputArg := capdef.CapabilityArgument{
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
	pageArg := capdef.CapabilityArgument{
		Name:        "page",
		Type:        capdef.ArgumentTypeInteger,
		Description: "Page number to generate thumbnail from (1-based, default: 1)",
		CliFlag:     "--page",
		Validation:  pageValidation,
		Default:     pageDefault,
	}
	arguments.AddOptional(pageArg)
	
	output := &capdef.CapabilityOutput{
		Type:        capdef.OutputTypeBinary,
		ContentType: stringPtr("image/png"),
		Validation:  &capdef.ArgumentValidation{},
		Description: "PNG image data representing a thumbnail of the document",
	}
	
	capability := capdef.NewCapabilityWithDescription(
		id,
		"1.0.0",
		command,
		"Generate a thumbnail image preview of the document",
	)
	capability.SetArguments(arguments)
	capability.SetOutput(output)
	
	return capability
}

// ExtractOutlineCapability creates the standard extract-outline capability with full argument definition
func ExtractOutlineCapability() *capdef.Capability {
	id, _ := capdef.NewCapabilityIdFromString("document:extract:outline")
	
	command := "extract-outline"
	
	arguments := capdef.NewCapabilityArguments()
	
	// Required file_path argument
	filePathValidation := &capdef.ArgumentValidation{
		Pattern:   stringPtr("^[^\\0]+$"),
		MinLength: intPtr(1),
	}
	filePathArg := capdef.CapabilityArgument{
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
	maxDepthArg := capdef.CapabilityArgument{
		Name:        "max_depth",
		Type:        capdef.ArgumentTypeInteger,
		Description: "Maximum outline depth to extract (1-10)",
		CliFlag:     "--max-depth",
		Validation:  maxDepthValidation,
	}
	arguments.AddOptional(maxDepthArg)
	
	// Optional include_page_numbers argument
	includePageNumbersArg := capdef.CapabilityArgument{
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
	outputArg := capdef.CapabilityArgument{
		Name:        "output",
		Type:        capdef.ArgumentTypeString,
		Description: "Write output to specified file instead of stdout",
		CliFlag:     "--output",
		Validation:  outputValidation,
	}
	arguments.AddOptional(outputArg)
	
	output := &capdef.CapabilityOutput{
		Type:        capdef.OutputTypeObject,
		SchemaRef:   stringPtr("document-outline.json"),
		ContentType: stringPtr("application/json"),
		Validation:  &capdef.ArgumentValidation{},
		Description: "Hierarchical document outline with section titles and optional page numbers",
	}
	
	capability := capdef.NewCapabilityWithDescription(
		id,
		"1.0.0",
		command,
		"Extract document outline/table of contents with hierarchical structure",
	)
	capability.SetArguments(arguments)
	capability.SetOutput(output)
	
	return capability
}

// ExtractPagesCapability creates the standard extract-pages capability with full argument definition
func ExtractPagesCapability() *capdef.Capability {
	id, _ := capdef.NewCapabilityIdFromString("document:extract:pages")
	
	command := "extract-pages"
	
	arguments := capdef.NewCapabilityArguments()
	
	// Required file_path argument
	filePathValidation := &capdef.ArgumentValidation{
		Pattern:   stringPtr("^[^\\0]+$"),
		MinLength: intPtr(1),
	}
	filePathArg := capdef.CapabilityArgument{
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
	outputArg := capdef.CapabilityArgument{
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
	pageRangeArg := capdef.CapabilityArgument{
		Name:        "page_range",
		Type:        capdef.ArgumentTypeString,
		Description: "Page range to extract (e.g., '1-5' or '10-')",
		CliFlag:     "--page-range",
		Validation:  pageRangeValidation,
	}
	arguments.AddOptional(pageRangeArg)
	
	output := &capdef.CapabilityOutput{
		Type:        capdef.OutputTypeObject,
		SchemaRef:   stringPtr("document-pages.json"),
		ContentType: stringPtr("application/json"),
		Validation:  &capdef.ArgumentValidation{},
		Description: "Structured page content extracted from the document",
	}
	
	capability := capdef.NewCapabilityWithDescription(
		id,
		"1.0.0",
		command,
		"Extract structured page content from the document",
	)
	capability.SetArguments(arguments)
	capability.SetOutput(output)
	
	return capability
}

// GetAllStandardCapabilities returns all standard plugin capabilities
func GetAllStandardCapabilities() *capdef.PluginCapabilities {
	capabilities := capdef.NewPluginCapabilities()
	capabilities.AddCapability(ExtractMetadataCapability())
	capabilities.AddCapability(GenerateThumbnailCapability())
	capabilities.AddCapability(ExtractOutlineCapability())
	capabilities.AddCapability(ExtractPagesCapability())
	return capabilities
}

// GetStandardCapability returns a standard capability by name
func GetStandardCapability(name string) *capdef.Capability {
	switch name {
	case "extract-metadata":
		return ExtractMetadataCapability()
	case "generate-thumbnail":
		return GenerateThumbnailCapability()
	case "extract-outline":
		return ExtractOutlineCapability()
	case "extract-pages":
		return ExtractPagesCapability()
	default:
		return nil
	}
}

// GetStandardCapabilityById returns a standard capability by capability ID string
func GetStandardCapabilityById(idStr string) *capdef.Capability {
	switch idStr {
	case "document:extract:metadata":
		return ExtractMetadataCapability()
	case "document:generate:thumbnail":
		return GenerateThumbnailCapability()
	case "document:extract:outline":
		return ExtractOutlineCapability()
	case "document:extract:pages":
		return ExtractPagesCapability()
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