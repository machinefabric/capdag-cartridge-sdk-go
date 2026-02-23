// Package sdk provides standard cap definitions with arguments
package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/filegrind/capns-go/cap"
	"github.com/filegrind/capns-go/standard"
	"github.com/filegrind/capns-go/urn"
)

// Spec ID constants for domain-specific media URNs
const (
	SpecIdFileMetadata    = "media:file-metadata"
	SpecIdDocumentOutline = "media:document-outline"
	SpecIdDisboundPage    = "media:disbound-page"
)

// InputSpecIdForExt returns the input media URN for a given file extension
func InputSpecIdForExt(ext string) string {
	if ext == "pdf" {
		return standard.MediaBinary
	}
	return standard.MediaString
}

// ExtractMetadataOutputSpecIdForExt returns the output media URN for extract-metadata by extension
func ExtractMetadataOutputSpecIdForExt(ext string) string {
	if ext == "pdf" {
		return SpecIdFileMetadata
	}
	return standard.MediaObject
}

// ExtractOutlineOutputSpecIdForExt returns the output media URN for extract-outline by extension
func ExtractOutlineOutputSpecIdForExt(ext string) string {
	if ext == "pdf" {
		return SpecIdDocumentOutline
	}
	return standard.MediaObject
}

// DisboundPageSpecIdForExt returns the output media URN for disbind by extension
func DisboundPageSpecIdForExt(ext string) string {
	if ext == "pdf" {
		return SpecIdDisboundPage
	}
	return standard.MediaObject
}

// ExtractMetadataUrn builds the URN for extract-metadata capability with given extension
func ExtractMetadataUrn(ext string) string {
	inSpec := InputSpecIdForExt(ext)
	outSpec := ExtractMetadataOutputSpecIdForExt(ext)
	return fmt.Sprintf("cap:ext=%s;in=%q;op=extract_metadata;out=%q", ext, inSpec, outSpec)
}

// GenerateThumbnailUrn builds the URN for generate-thumbnail capability with given extension
func GenerateThumbnailUrn(ext string) string {
	inSpec := InputSpecIdForExt(ext)
	return fmt.Sprintf("cap:ext=%s;in=%q;op=generate_thumbnail;out=%q", ext, inSpec, standard.MediaBinary)
}

// ExtractOutlineUrn builds the URN for extract-outline capability with given extension
func ExtractOutlineUrn(ext string) string {
	inSpec := InputSpecIdForExt(ext)
	outSpec := ExtractOutlineOutputSpecIdForExt(ext)
	return fmt.Sprintf("cap:ext=%s;in=%q;op=extract_outline;out=%q", ext, inSpec, outSpec)
}

// GrindUrn builds the URN for grind capability with given extension
func GrindUrn(ext string) string {
	inSpec := InputSpecIdForExt(ext)
	outSpec := DisboundPageSpecIdForExt(ext)
	return fmt.Sprintf("cap:ext=%s;in=%q;op=grind;out=%q", ext, inSpec, outSpec)
}

// ExtractMetadataCap creates the standard extract-metadata cap with full argument definition
func ExtractMetadataCap() *cap.Cap {
	id, _ := urn.NewCapUrnFromString("cap:op=extract_metadata")

	c := cap.NewCapWithDescription(
		id,
		"Extract Document Metadata",
		"extract-metadata",
		"Extract document metadata including title, author, creation date, file size, and other properties",
	)

	// Required file_path argument (positional + stdin)
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaFilePath,
		Required: true,
		Sources: []cap.ArgSource{
			{Position: intPtr(0)},
			{Stdin: stringPtr(standard.MediaBinary)},
		},
		ArgDescription: "Path to the document file to process",
	})

	// Optional output argument
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaString,
		Required: false,
		Sources: []cap.ArgSource{
			{CliFlag: stringPtr("--output")},
		},
		ArgDescription: "Write output to specified file instead of stdout",
	})

	c.SetOutput(cap.NewCapOutput(
		standard.MediaObject,
		"Structured metadata including file properties, document properties, and format-specific metadata",
	))

	return c
}

// GenerateThumbnailCap creates the standard generate-thumbnail cap with full argument definition
func GenerateThumbnailCap() *cap.Cap {
	id, _ := urn.NewCapUrnFromString("cap:op=generate_thumbnail")

	c := cap.NewCapWithDescription(
		id,
		"Generate Thumbnail",
		"generate-thumbnail",
		"Generate a thumbnail image preview of the document",
	)

	// Required file_path argument (positional + stdin)
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaFilePath,
		Required: true,
		Sources: []cap.ArgSource{
			{Position: intPtr(0)},
			{Stdin: stringPtr(standard.MediaBinary)},
		},
		ArgDescription: "Path to the document file to process",
	})

	// Optional width argument
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaInteger,
		Required: false,
		Sources: []cap.ArgSource{
			{CliFlag: stringPtr("--width")},
		},
		ArgDescription:  "Width of the thumbnail in pixels",
		DefaultValue: json.Number("200"),
	})

	// Optional height argument
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaInteger,
		Required: false,
		Sources: []cap.ArgSource{
			{CliFlag: stringPtr("--height")},
		},
		ArgDescription:  "Height of the thumbnail in pixels",
		DefaultValue: json.Number("300"),
	})

	// Optional output argument
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaString,
		Required: false,
		Sources: []cap.ArgSource{
			{CliFlag: stringPtr("--output")},
		},
		ArgDescription: "Write thumbnail to specified file instead of stdout",
	})

	// Optional page argument
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaInteger,
		Required: false,
		Sources: []cap.ArgSource{
			{CliFlag: stringPtr("--page")},
		},
		ArgDescription:  "Page number to generate thumbnail from (1-based, default: 1)",
		DefaultValue: json.Number("1"),
	})

	c.SetOutput(cap.NewCapOutput(
		standard.MediaBinary,
		"PNG image data representing a thumbnail of the document",
	))

	return c
}

// ExtractOutlineCap creates the standard extract-outline cap with full argument definition
func ExtractOutlineCap() *cap.Cap {
	id, _ := urn.NewCapUrnFromString("cap:op=extract_outline")

	c := cap.NewCapWithDescription(
		id,
		"Extract Document Outline",
		"extract-outline",
		"Extract document outline/table of contents with hierarchical structure",
	)

	// Required file_path argument (positional + stdin)
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaFilePath,
		Required: true,
		Sources: []cap.ArgSource{
			{Position: intPtr(0)},
			{Stdin: stringPtr(standard.MediaBinary)},
		},
		ArgDescription: "Path to the document file to process",
	})

	// Optional max_depth argument
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaInteger,
		Required: false,
		Sources: []cap.ArgSource{
			{CliFlag: stringPtr("--max-depth")},
		},
		ArgDescription: "Maximum outline depth to extract (1-10)",
	})

	// Optional include_order_indexes argument
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaBoolean,
		Required: false,
		Sources: []cap.ArgSource{
			{CliFlag: stringPtr("--include-order-indexes")},
		},
		ArgDescription:  "Include page numbers in the outline (default: true)",
		DefaultValue: true,
	})

	// Optional output argument
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaString,
		Required: false,
		Sources: []cap.ArgSource{
			{CliFlag: stringPtr("--output")},
		},
		ArgDescription: "Write output to specified file instead of stdout",
	})

	c.SetOutput(cap.NewCapOutput(
		standard.MediaObject,
		"Hierarchical document outline with section titles and optional page numbers",
	))

	return c
}

// DisbindCap creates the standard grind cap with full argument definition
func DisbindCap() *cap.Cap {
	id, _ := urn.NewCapUrnFromString("cap:op=grind")

	c := cap.NewCapWithDescription(
		id,
		"Extract File Chips",
		"grind",
		"Extract structured page content from the document",
	)

	// Required file_path argument (positional + stdin)
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaFilePath,
		Required: true,
		Sources: []cap.ArgSource{
			{Position: intPtr(0)},
			{Stdin: stringPtr(standard.MediaBinary)},
		},
		ArgDescription: "Path to the document file to process",
	})

	// Optional output argument
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaString,
		Required: false,
		Sources: []cap.ArgSource{
			{CliFlag: stringPtr("--output")},
		},
		ArgDescription: "Write output to specified file instead of stdout",
	})

	// Optional index_range argument
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaString,
		Required: false,
		Sources: []cap.ArgSource{
			{CliFlag: stringPtr("--index-range")},
		},
		ArgDescription: "Index Range to extract (e.g., '1-5' or '10-')",
	})

	c.SetOutput(cap.NewCapOutput(
		standard.MediaObjectArray,
		"Array of structured page content extracted from the document",
	))

	return c
}

// GetAllStandardCaps returns all standard plugin caps
func GetAllStandardCaps() []*cap.Cap {
	return []*cap.Cap{
		ExtractMetadataCap(),
		GenerateThumbnailCap(),
		ExtractOutlineCap(),
		DisbindCap(),
	}
}

// GetStandardCap returns a standard cap by name
func GetStandardCap(name string) *cap.Cap {
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
func GetStandardCapByUrn(urnStr string) *cap.Cap {
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

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
