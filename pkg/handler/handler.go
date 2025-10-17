// Package handler provides the core DocumentHandler interface and related types
package handler

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/jowharshamshiri/lbvr-plugin-sdk-go/pkg/metadata"
	"github.com/jowharshamshiri/lbvr-plugin-sdk-go/pkg/outline"
	"github.com/jowharshamshiri/lbvr-plugin-sdk-go/pkg/pages"
)

// DocumentHandler is the core interface that all document file handlers must implement
type DocumentHandler interface {
	// Get the name of this handler
	Name() string

	// Get the version of this handler
	Version() string

	// Get handler capabilities with file type specificity
	GetCapabilities() *PluginCapabilities

	// Check if this handler can process the given file
	CanHandle(filePath string) bool

	// Extract document metadata
	ExtractMetadata(ctx context.Context, filePath string) (*metadata.FileMetadata, error)

	// Extract document outline/table of contents
	ExtractOutline(ctx context.Context, filePath string) (*outline.DocumentOutline, error)

	// Extract document pages with text content organized by pages and paragraphs
	ExtractPages(ctx context.Context, filePath string) (*pages.DocumentPages, error)

	// Validate that the file is not corrupted and can be processed
	ValidateFile(ctx context.Context, filePath string) (bool, error)

	// Get basic file information without full processing
	GetFileInfo(ctx context.Context, filePath string) (*FileInfo, error)

	// Generate thumbnail image for the document
	// Returns PNG image data
	GenerateThumbnail(ctx context.Context, filePath string, width, height uint32) ([]byte, error)
}

// BaseDocumentHandler provides default implementations for optional methods
type BaseDocumentHandler struct{}

// CanHandle provides default implementation based on capabilities
func (b *BaseDocumentHandler) CanHandle(filePath string, capabilities *PluginCapabilities) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:] // Remove the dot
	}

	return capabilities.CanHandleFileType(ext)
}

// GetCapabilities provides default capabilities
func (b *BaseDocumentHandler) GetCapabilities() *PluginCapabilities {
	return &PluginCapabilities{
		Capabilities: []string{
			"extract_metadata",
			"extract_outline",
			"extract_pages",
			"validate_file",
			"generate_thumbnail",
			"supports_json_output",
		},
	}
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

// PluginCapabilities represents plugin capabilities
type PluginCapabilities struct {
	Capabilities []string `json:"capabilities"`
}

// Can checks if plugin has a specific capability
func (pc *PluginCapabilities) Can(capability string) bool {
	for _, cap := range pc.Capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}

// CanHandleFileType checks if the plugin can handle a specific file type
func (pc *PluginCapabilities) CanHandleFileType(fileType string) bool {
	// Check for exact match with file type (e.g., "extract_metadata:pdf")
	for _, capability := range pc.Capabilities {
		if strings.Contains(capability, ":") {
			parts := strings.SplitN(capability, ":", 2)
			if len(parts) == 2 {
				operation, filetype := parts[0], parts[1]
				if filetype == fileType && pc.isExtractOperation(operation) {
					return true
				}
			}
		}
	}
	
	// Check for wildcard match (e.g., "extract_metadata:*")
	for _, capability := range pc.Capabilities {
		if strings.Contains(capability, ":") {
			parts := strings.SplitN(capability, ":", 2)
			if len(parts) == 2 {
				operation, filetype := parts[0], parts[1]
				if filetype == "*" && pc.isExtractOperation(operation) {
					return true
				}
			}
		}
	}
	
	return false
}

// isExtractOperation checks if an operation is a document processing operation
func (pc *PluginCapabilities) isExtractOperation(operation string) bool {
	switch operation {
	case "extract_metadata", "extract_outline", "extract_pages", 
		 "extract_text", "validate_file", "generate_thumbnail":
		return true
	default:
		return false
	}
}

// GetMostSpecificCapability gets the most specific capability for a given operation and file type
func (pc *PluginCapabilities) GetMostSpecificCapability(operation, fileType string) *string {
	// First look for exact file type match
	specific := operation + ":" + fileType
	for _, cap := range pc.Capabilities {
		if cap == specific {
			return &specific
		}
	}
	
	// Then look for wildcard match
	wildcard := operation + ":*"
	for _, cap := range pc.Capabilities {
		if cap == wildcard {
			return &wildcard
		}
	}
	
	// Finally check for operation without file type specifier (legacy support)
	for _, cap := range pc.Capabilities {
		if cap == operation {
			return &operation
		}
	}
	
	return nil
}

// CanPerformOperation checks if the plugin can perform an operation on a specific file type
func (pc *PluginCapabilities) CanPerformOperation(operation, fileType string) bool {
	return pc.GetMostSpecificCapability(operation, fileType) != nil
}


// Plugin priority levels
type PluginPriority string

const (
	PluginPriorityOptional    PluginPriority = "optional"
	PluginPriorityRecommended PluginPriority = "recommended"
	PluginPriorityCritical    PluginPriority = "critical"
)

// PluginInfo represents plugin information
type PluginInfo struct {
	// Plugin name
	Name string `json:"name"`

	// Plugin version
	Version string `json:"version"`

	// Plugin description
	Description string `json:"description"`

	// Plugin priority level
	Priority PluginPriority `json:"priority"`

	// Plugin capabilities with file type specificity
	Capabilities *PluginCapabilities `json:"capabilities"`

	// Plugin author/maintainer
	Author *string `json:"author,omitempty"`
}

// NewPluginInfo creates a new plugin info
func NewPluginInfo(name, version, description string, capabilities *PluginCapabilities, priority PluginPriority) *PluginInfo {
	return &PluginInfo{
		Name:         name,
		Version:      version,
		Description:  description,
		Priority:     priority,
		Capabilities: capabilities,
	}
}

// WithAuthor sets the author of the plugin
func (pi *PluginInfo) WithAuthor(author string) *PluginInfo {
	pi.Author = &author
	return pi
}

// PluginMetadata interface for plugins to provide metadata about themselves
type PluginMetadata interface {
	// Get plugin information
	PluginInfo() *PluginInfo

	// Get plugin capabilities
	Capabilities() *PluginCapabilities
}

// ProcessingResult represents the result of a document processing operation
type ProcessingResult struct {
	// Whether the operation was successful
	Success bool `json:"success"`

	// Processing time in milliseconds
	ProcessingTimeMs uint64 `json:"processing_time_ms"`

	// Any warnings generated during processing
	Warnings []string `json:"warnings"`

	// Error message if operation failed
	Error *string `json:"error,omitempty"`
}

// NewSuccessResult creates a successful result
func NewSuccessResult(processingTimeMs uint64) *ProcessingResult {
	return &ProcessingResult{
		Success:          true,
		ProcessingTimeMs: processingTimeMs,
		Warnings:         make([]string, 0),
	}
}

// NewFailureResult creates a failed result
func NewFailureResult(errorMsg string, processingTimeMs uint64) *ProcessingResult {
	return &ProcessingResult{
		Success:          false,
		ProcessingTimeMs: processingTimeMs,
		Warnings:         make([]string, 0),
		Error:            &errorMsg,
	}
}

// AddWarning adds a warning to the result
func (pr *ProcessingResult) AddWarning(warning string) {
	pr.Warnings = append(pr.Warnings, warning)
}

// HandlerRegistry manages registered document handlers
type HandlerRegistry struct {
	handlers []DocumentHandler
}

// NewHandlerRegistry creates a new empty registry
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make([]DocumentHandler, 0),
	}
}

// Register registers a document handler
func (hr *HandlerRegistry) Register(handler DocumentHandler) {
	hr.handlers = append(hr.handlers, handler)
}

// FindHandler finds a handler for the given file
func (hr *HandlerRegistry) FindHandler(filePath string) DocumentHandler {
	for _, handler := range hr.handlers {
		if handler.CanHandle(filePath) {
			return handler
		}
	}
	return nil
}

// Handlers gets all registered handlers
func (hr *HandlerRegistry) Handlers() []DocumentHandler {
	return hr.handlers
}

// HandlersForFileType gets handlers that support a specific file type
func (hr *HandlerRegistry) HandlersForFileType(fileType string) []DocumentHandler {
	var result []DocumentHandler
	for _, handler := range hr.handlers {
		if handler.GetCapabilities().CanHandleFileType(fileType) {
			result = append(result, handler)
		}
	}
	return result
}

// HandlerCount gets the number of registered handlers
func (hr *HandlerRegistry) HandlerCount() int {
	return len(hr.handlers)
}

// SupportedFileTypes gets all supported file types
func (hr *HandlerRegistry) SupportedFileTypes() []string {
	fileTypeSet := make(map[string]bool)
	for _, handler := range hr.handlers {
		for _, capability := range handler.GetCapabilities().Capabilities {
			if strings.Contains(capability, ":") {
				parts := strings.SplitN(capability, ":", 2)
				if len(parts) == 2 {
					fileType := parts[1]
					if fileType != "*" {
						fileTypeSet[strings.ToLower(fileType)] = true
					}
				}
			}
		}
	}

	var fileTypes []string
	for fileType := range fileTypeSet {
		fileTypes = append(fileTypes, fileType)
	}
	return fileTypes
}

// FindBestHandler finds the best handler for a specific operation and file type
func (hr *HandlerRegistry) FindBestHandler(operation, fileType string) DocumentHandler {
	var bestHandler DocumentHandler
	bestSpecificity := 0
	
	for _, handler := range hr.handlers {
		capabilities := handler.GetCapabilities()
		if capability := capabilities.GetMostSpecificCapability(operation, fileType); capability != nil {
			specificity := 0
			if strings.Contains(*capability, ":"+fileType) {
				specificity = 2 // Exact file type match
			} else if strings.Contains(*capability, ":*") {
				specificity = 1 // Wildcard match
			} else {
				specificity = 0 // Legacy operation-only match
			}
			
			if specificity > bestSpecificity {
				bestHandler = handler
				bestSpecificity = specificity
			}
		}
	}
	
	return bestHandler
}

// IsSupported checks if a file is supported by any handler
func (hr *HandlerRegistry) IsSupported(filePath string) bool {
	return hr.FindHandler(filePath) != nil
}

// ToJSON converts any struct to JSON string
func ToJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}