// Package handler provides the unified cap-based plugin interface
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	
	capdef "github.com/lbvr/capdef-go"
)

// PluginRegistry provides unified cap-based access to plugins
type PluginRegistry struct {
	plugins map[string]*PluginEntry
	capIndex map[string][]string // cap -> plugin names
}

// PluginEntry represents a registered plugin
type PluginEntry struct {
	BinaryPath string
	Caps []string
}

// CapCaller provides the unified interface for calling plugin caps
type CapCaller struct {
	PluginName string
	Cap string
	BinaryPath string
}

// ResponseWrapper provides type-safe deserialization of plugin output
type ResponseWrapper struct {
	data []byte
}

// NewPluginRegistry creates a new empty plugin registry
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[string]*PluginEntry),
		capIndex: make(map[string][]string),
	}
}

// RegisterPlugin registers a plugin with its caps
func (pr *PluginRegistry) RegisterPlugin(name, binaryPath string, caps []string) {
	entry := &PluginEntry{
		BinaryPath: binaryPath,
		Caps: caps,
	}
	
	// Update cap index
	for _, cap := range caps {
		if _, exists := pr.capIndex[cap]; !exists {
			pr.capIndex[cap] = make([]string, 0)
		}
		pr.capIndex[cap] = append(pr.capIndex[cap], name)
	}
	
	pr.plugins[name] = entry
}

// Can checks if a cap is available and returns a caller
func (pr *PluginRegistry) Can(cap string) (*CapCaller, error) {
	// Find the best plugin for this cap
	pluginName := pr.findBestPluginForCap(cap)
	if pluginName == "" {
		return nil, fmt.Errorf("cap '%s' is not available in any registered plugin", cap)
	}
	
	plugin, exists := pr.plugins[pluginName]
	if !exists {
		return nil, fmt.Errorf("plugin '%s' not found in registry", pluginName)
	}
	
	return &CapCaller{
		PluginName: pluginName,
		Cap: cap,
		BinaryPath: plugin.BinaryPath,
	}, nil
}

// CallWithStdin executes the cap with the given arguments and optional stdin data
func (cc *CapCaller) CallWithStdin(ctx context.Context, args []interface{}, stdinData []byte) (*ResponseWrapper, error) {
	// Convert cap to CLI flag
	operation := strings.SplitN(cc.Cap, ":", 2)[0]
	command := "--" + operation
	
	// Build command arguments
	cmdArgs := []string{command}
	for _, arg := range args {
		cmdArgs = append(cmdArgs, fmt.Sprintf("%v", arg))
	}
	
	// Execute the plugin
	cmd := exec.CommandContext(ctx, cc.BinaryPath, cmdArgs...)
	
	// Set stdin if provided
	if stdinData != nil {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
		}
		
		go func() {
			defer stdin.Close()
			stdin.Write(stdinData)
		}()
	}
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("plugin execution failed: %w", err)
	}
	
	return &ResponseWrapper{data: output}, nil
}

// AsType deserializes the response to a specific type
func (rw *ResponseWrapper) AsType(v interface{}) error {
	return json.Unmarshal(rw.data, v)
}

// AsString converts the response to a string
func (rw *ResponseWrapper) AsString() (string, error) {
	return string(rw.data), nil
}

// AsBytes returns the raw byte array
func (rw *ResponseWrapper) AsBytes() []byte {
	return rw.data
}

// AsInt converts the response to an integer
func (rw *ResponseWrapper) AsInt() (int64, error) {
	var result int64
	err := json.Unmarshal(rw.data, &result)
	return result, err
}

// AsBool converts the response to a boolean
func (rw *ResponseWrapper) AsBool() (bool, error) {
	var result bool
	err := json.Unmarshal(rw.data, &result)
	return result, err
}

// findBestPluginForCap finds the best plugin for a cap
func (pr *PluginRegistry) findBestPluginForCap(cap string) string {
	candidates := pr.getCapCandidates(cap)
	if len(candidates) == 0 {
		return ""
	}
	
	// Find the candidate with the highest specificity
	bestPlugin := ""
	bestScore := -1
	
	for _, pluginName := range candidates {
		plugin := pr.plugins[pluginName]
		score := pr.calculateCapScore(plugin, cap)
		if score > bestScore {
			bestPlugin = pluginName
			bestScore = score
		}
	}
	
	return bestPlugin
}

// getCapCandidates returns plugins that might support the cap
func (pr *PluginRegistry) getCapCandidates(cap string) []string {
	// Direct match
	if plugins, exists := pr.capIndex[cap]; exists {
		return plugins
	}
	
	// Try wildcard variations
	if strings.Contains(cap, ":") {
		parts := strings.SplitN(cap, ":", 2)
		if len(parts) == 2 {
			wildcardCap := parts[0] + ":*"
			if plugins, exists := pr.capIndex[wildcardCap]; exists {
				return plugins
			}
		}
	}
	
	return []string{}
}

// calculateCapScore calculates specificity score for a plugin cap match
func (pr *PluginRegistry) calculateCapScore(plugin *PluginEntry, cap string) int {
	score := 0
	
	// Add specificity score
	for _, cap := range plugin.Caps {
		if cap == cap {
			if strings.Contains(cap, ":") && !strings.HasSuffix(cap, ":*") {
				score += 20 // Exact file type match
			} else if strings.HasSuffix(cap, ":*") {
				score += 10 // Wildcard match
			} else {
				score += 5 // Operation-only match
			}
			break
		}
	}
	
	return score
}

// ListCaps returns all available caps
func (pr *PluginRegistry) ListCaps() []string {
	caps := make([]string, 0, len(pr.capIndex))
	for cap := range pr.capIndex {
		caps = append(caps, cap)
	}
	return caps
}


// Re-export CapManifest from capdef as PluginManifest for backward compatibility
type PluginManifest = capdef.CapManifest

// NewPluginManifest creates a new plugin manifest
func NewPluginManifest(name, version, description string, caps []capdef.Cap) *capdef.CapManifest {
	return capdef.NewCapManifest(name, version, description, caps)
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