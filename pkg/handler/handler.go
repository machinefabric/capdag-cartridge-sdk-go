// Package handler provides the unified cap-based plugin interface  
package handler

import (
	"encoding/json"
	
	capns "github.com/fmio/cap-sdk-go"
)

// Re-export CapManifest from capns as PluginManifest for backward compatibility
type PluginManifest = capns.CapManifest

// NewPluginManifest creates a new plugin manifest
func NewPluginManifest(name, version, description string, caps []capns.Cap) *capns.CapManifest {
	return capns.NewCapManifest(name, version, description, caps)
}

// FileInfo represents basic file information
type FileInfo struct {
	// File path
	Path string `json:"path"`

	// File size in bytes
	Size uint64 `json:"size"`

	// Document type detected
	DocumentType string `json:"document_type"`

	// Whether the file appears to be valid
	IsValid bool `json:"is_valid"`

	// Quick metadata (title, author if easily accessible)
	QuickMetadata *QuickMetadata `json:"quick_metadata,omitempty"`
}

// QuickMetadata represents quick metadata that can be extracted without full processing
type QuickMetadata struct {
	// Document title
	Title *string `json:"title,omitempty"`

	// Primary author
	Author *string `json:"author,omitempty"`

	// Page/section count
	PageCount *uint64 `json:"page_count,omitempty"`
}


// StandardizedCaps are the standard cap names
var StandardizedCaps = struct {
	ExtractMetadata   string
	ExtractOutline    string
	ExtractPages      string
	GenerateThumbnail string
	ValidateFile      string
}{
	ExtractMetadata:   "extract-metadata",
	ExtractOutline:    "extract-outline", 
	ExtractPages:      "extract-pages",
	GenerateThumbnail: "generate-thumbnail",
	ValidateFile:      "validate-file",
}

// ToJSON converts any struct to JSON string
func ToJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DocumentHandler interface defines the contract for document processing plugins
type DocumentHandler interface {
	// GetPluginManifest returns plugin manifest including caps
	GetPluginManifest() *PluginManifest
	
	// ExtractMetadata extracts metadata from a document
	ExtractMetadata(filePath string) (*ProcessingResult, error)
	
	// ExtractOutline extracts the outline/table of contents from a document
	ExtractOutline(filePath string) (*ProcessingResult, error)
	
	// ExtractText extracts plain text from a document
	ExtractText(filePath string) (*ProcessingResult, error)
	
	// GenerateThumbnail generates a thumbnail image from a document
	GenerateThumbnail(filePath string, width, height int, page int) (*ProcessingResult, error)
}

// PluginMetadata represents metadata about the plugin itself
type PluginMetadata struct {
	// Plugin name
	Name string `json:"name"`
	
	// Plugin version
	Version string `json:"version"`
	
	// Plugin description
	Description string `json:"description"`
	
	// Supported file types
	SupportedTypes []string `json:"supported_types"`
	
	// Plugin caps
	Caps []string `json:"caps"`
	
	// Plugin author
	Author *string `json:"author,omitempty"`
}

// ProcessingResult represents the result of a document processing operation
type ProcessingResult struct {
	// Whether the operation was successful
	Success bool `json:"success"`
	
	// Result data (can be any type)
	Data interface{} `json:"data,omitempty"`
	
	// Error message if operation failed
	Error *string `json:"error,omitempty"`
	
	// Processing time in milliseconds
	ProcessingTimeMs *int64 `json:"processing_time_ms,omitempty"`
	
	// File information
	FileInfo *FileInfo `json:"file_info,omitempty"`
}

// HandlerRegistry manages document handlers for different file types
type HandlerRegistry struct {
	handlers map[string]DocumentHandler
}

// NewHandlerRegistry creates a new handler registry
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]DocumentHandler),
	}
}

// RegisterHandler registers a handler for specific file types
func (hr *HandlerRegistry) RegisterHandler(fileTypes []string, handler DocumentHandler) {
	for _, fileType := range fileTypes {
		hr.handlers[fileType] = handler
	}
}

// GetHandler gets a handler for a specific file type
func (hr *HandlerRegistry) GetHandler(fileType string) (DocumentHandler, bool) {
	handler, exists := hr.handlers[fileType]
	return handler, exists
}

// NewSuccessResult creates a successful processing result
func NewSuccessResult(data interface{}) *ProcessingResult {
	return &ProcessingResult{
		Success: true,
		Data:    data,
	}
}

// NewFailureResult creates a failed processing result
func NewFailureResult(errorMsg string) *ProcessingResult {
	return &ProcessingResult{
		Success: false,
		Error:   &errorMsg,
	}
}