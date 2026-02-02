# FGND Plugin SDK for Go

A Go SDK for building document processing plugins for the FGND (FileGrind) system. This SDK provides unified data structures and interfaces for extracting metadata, outlines, and text content from various document formats.

## Features

- **Unified Interface**: Consistent `DocumentHandler` interface for all document types
- **Rich Metadata**: Comprehensive metadata extraction with format-specific fields
- **Hierarchical Outlines**: Support for nested table of contents structures
- **Page-based Content**: Organized text extraction with pages and paragraphs
- **JSON Compatible**: All data structures support JSON serialization
- **Extensible**: Plugin-specific metadata through extended fields

## Installation

```bash
go get github.com/jowharshamshiri/fgnd-plugin-sdk-go
```

## Quick Start

### Implementing a Document Handler

```go
package main

import (
    "context"
    "fmt"
    
    sdk "github.com/jowharshamshiri/fgnd-plugin-sdk-go"
)

// MyDocumentHandler implements the DocumentHandler interface
type MyDocumentHandler struct {
    sdk.BaseDocumentHandler
}

func (h *MyDocumentHandler) Name() string {
    return "my-handler"
}

func (h *MyDocumentHandler) Version() string {
    return "1.0.0"
}

func (h *MyDocumentHandler) SupportedExtensions() []string {
    return []string{"txt", "md"}
}

func (h *MyDocumentHandler) CanHandle(filePath string) bool {
    return h.BaseDocumentHandler.CanHandle(filePath, h.SupportedExtensions())
}

func (h *MyDocumentHandler) ExtractMetadata(ctx context.Context, filePath string) (*sdk.FileMetadata, error) {
    // Get file info
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        return nil, err
    }
    
    // Create metadata
    metadata := sdk.NewFileMetadata(filePath, "text", uint64(fileInfo.Size()))
    metadata.SetExtended("extractor", "my-handler")
    
    return metadata, nil
}

func (h *MyDocumentHandler) ExtractOutline(ctx context.Context, filePath string) (*sdk.DocumentOutline, error) {
    outline := sdk.NewDocumentOutline(filePath, "text", 1)
    outline.ExtractionInfo = *sdk.NewExtractionInfo("my-handler", "1.0.0")
    return outline, nil
}

func (h *MyDocumentHandler) Disbind(ctx context.Context, filePath string) ([]sdk.DisboundPage, error) {
    // Read file content
    content, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    // Create pages array
    var pages []sdk.DisboundPage
    page := sdk.NewDisboundPageWithText(1, string(content))
    pages = append(pages, *page)

    return pages, nil
}

// ... implement other required methods
```

### Using the Handler Registry

```go
func main() {
    // Create registry and register handler
    registry := sdk.NewHandlerRegistry()
    registry.Register(&MyDocumentHandler{})
    
    // Find handler for a file
    handler := registry.FindHandler("document.txt")
    if handler != nil {
        // Extract metadata
        metadata, err := handler.ExtractMetadata(context.Background(), "document.txt")
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }
        
        // Convert to JSON
        jsonStr, _ := metadata.ToJSON()
        fmt.Println(jsonStr)
    }
}
```

## Data Structures

### FileMetadata

Consolidated metadata structure supporting:
- Common document fields (title, authors, dates)
- Format-specific fields (PDF, EPUB specific)
- Extensible metadata via `ExtendedMetadata` map
- Word/character/page counts

```go
metadata := sdk.NewFileMetadata("/path/to/file.pdf", "pdf", 1024*1024)
metadata.Title = &"Document Title"
metadata.AddAuthor("John Doe")
metadata.SetExtended("custom_field", "custom_value")
```

### DocumentOutline

Hierarchical table of contents:
- Nested Outline entries with unlimited depth
- Page/section references
- Source references (filenames, anchors)

```go
outline := sdk.NewDocumentOutline("/path/to/file.pdf", "pdf", 100)
entry := sdk.NewOutlineEntry("Chapter 1", 0).WithPage(5)
outline.AddEntry(*entry)
```

### DisboundPage (Array)

Page-based text content organization using a simple slice:
- Pages are 1-indexed via order_index
- Automatic word/character counting per page
- Output is a JSON array of page objects

```go
var pages []sdk.DisboundPage
page := sdk.NewDisboundPageWithText(1, "Page content here...")
pages = append(pages, *page)
```

## Plugin Caps

The SDK supports a cap-based system:

```go
caps := &sdk.PluginCaps{
    Caps: []string{
        "extract_metadata",
        "extract_outline", 
        "grind",
        "validate_file",
        "generate_thumbnail",
        "supports_json_output",
    },
}
```

## Error Handling

All handler methods return errors following Go conventions:

```go
metadata, err := handler.ExtractMetadata(ctx, filePath)
if err != nil {
    // Handle error
    return fmt.Errorf("failed to extract metadata: %w", err)
}
```

## JSON Schemas

This SDK implements the JSON schemas defined in the `fgnd/plugin-schemas/` directory:
- `file-metadata.json` - FileMetadata structure
- `document-outline.json` - DocumentOutline structure  
- `disbound-page.json` - DisboundPage structure (output is array of these)
- `manifest.json` - Plugin manifest
- `handler-interface.json` - DocumentHandler interface

## Testing

```bash
go test ./...
```

## Examples

See the `examples/` directory for complete implementation examples:
- Text file handler
- Markdown processor
- Plugin registration patterns

## Contributing

1. Ensure all data structures match the JSON schemas in `../fgnd/plugin-schemas/`
2. Implement the complete `DocumentHandler` interface
3. Add appropriate tests
4. Follow Go naming conventions

## License

MIT License - see LICENSE file for details.

## Related Projects

- **fgnd-plugin-sdk** - Rust SDK implementation
- **fgnd** - Main FGND engine
- **txtcartridge** - Text and Markdown processor
- **pdfcartridge** - PDF processor  
- **epubcartridge** - EPUB processor