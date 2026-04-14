// Package sdk provides standard cap definitions with arguments
package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/machinefabric/capdag-go/cap"
	"github.com/machinefabric/capdag-go/standard"
	"github.com/machinefabric/capdag-go/urn"
)

// InputSpecIdForExt returns the input media URN for a given file extension
func InputSpecIdForExt(ext string) string {
	if ext == "pdf" {
		return standard.MediaIdentity
	}
	return standard.MediaString
}

// RenderPageImageUrn builds the URN for render-page-image capability with given extension.
// Output is always a PNG page image.
func RenderPageImageUrn(ext string) string {
	inSpec := InputSpecIdForExt(ext)
	return fmt.Sprintf("cap:ext=%s;in=%q;op=render_page_image;out=%q", ext, inSpec, standard.MediaPNG)
}

// GrindUrn builds the URN for grind capability with given extension
func GrindUrn(ext string) string {
	inSpec := InputSpecIdForExt(ext)
	return fmt.Sprintf("cap:ext=%s;in=%q;op=grind;out=%q", ext, inSpec, standard.MediaTextablePage)
}

// RenderPageImageCap creates the standard render-page-image cap with full argument definition.
// Output is always a PNG page image.
func RenderPageImageCap() *cap.Cap {
	id, _ := urn.NewCapUrnFromString("cap:op=render_page_image")

	c := cap.NewCapWithDescription(
		id,
		"Render Page Image",
		"render-page-image",
		"Render a page of the document as a PNG image",
	)

	// Required file_path argument (positional + stdin)
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaFilePath,
		Required: true,
		Sources: []cap.ArgSource{
			{Position: intPtr(0)},
			{Stdin: stringPtr(standard.MediaIdentity)},
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
		ArgDescription: "Write rendered image to specified file instead of stdout",
	})

	// Optional page argument
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaInteger,
		Required: false,
		Sources: []cap.ArgSource{
			{CliFlag: stringPtr("--page")},
		},
		ArgDescription:  "Page number to render (1-based, default: 1)",
		DefaultValue: json.Number("1"),
	})

	// Optional width argument
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaInteger,
		Required: false,
		Sources: []cap.ArgSource{
			{CliFlag: stringPtr("--width")},
		},
		ArgDescription:  "Width of the output image in pixels",
		DefaultValue: json.Number("200"),
	})

	// Optional height argument
	c.AddArg(cap.CapArg{
		MediaUrn: standard.MediaInteger,
		Required: false,
		Sources: []cap.ArgSource{
			{CliFlag: stringPtr("--height")},
		},
		ArgDescription:  "Height of the output image in pixels",
		DefaultValue: json.Number("300"),
	})

	c.SetOutput(cap.NewCapOutput(
		standard.MediaPNG,
		"PNG image data of the rendered page",
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
			{Stdin: stringPtr(standard.MediaIdentity)},
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
		standard.MediaTextablePage,
		"Sequence of extracted page content from the document",
	))

	return c
}

// GetAllStandardCaps returns all standard cartridge caps
func GetAllStandardCaps() []*cap.Cap {
	return []*cap.Cap{
		RenderPageImageCap(),
		DisbindCap(),
	}
}

// GetStandardCap returns a standard cap by name
func GetStandardCap(name string) *cap.Cap {
	switch name {
	case "render-page-image":
		return RenderPageImageCap()
	case "grind":
		return DisbindCap()
	default:
		return nil
	}
}

// GetStandardCapByUrn returns a standard cap by cap URN string
func GetStandardCapByUrn(urnStr string) *cap.Cap {
	switch urnStr {
	case "cap:op=render_page_image":
		return RenderPageImageCap()
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
