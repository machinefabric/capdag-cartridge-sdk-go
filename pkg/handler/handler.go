// Package handler provides the unified capability-based plugin interface
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// PluginRegistry provides unified capability-based access to plugins
type PluginRegistry struct {
	plugins map[string]*PluginEntry
	capabilityIndex map[string][]string // capability -> plugin names
}

// PluginEntry represents a registered plugin
type PluginEntry struct {
	BinaryPath string
	Capabilities []string
}

// CapabilityCaller provides the unified interface for calling plugin capabilities
type CapabilityCaller struct {
	PluginName string
	Capability string
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
		capabilityIndex: make(map[string][]string),
	}
}

// RegisterPlugin registers a plugin with its capabilities
func (pr *PluginRegistry) RegisterPlugin(name, binaryPath string, capabilities []string) {
	entry := &PluginEntry{
		BinaryPath: binaryPath,
		Capabilities: capabilities,
	}
	
	// Update capability index
	for _, capability := range capabilities {
		if _, exists := pr.capabilityIndex[capability]; !exists {
			pr.capabilityIndex[capability] = make([]string, 0)
		}
		pr.capabilityIndex[capability] = append(pr.capabilityIndex[capability], name)
	}
	
	pr.plugins[name] = entry
}

// Can checks if a capability is available and returns a caller
func (pr *PluginRegistry) Can(capability string) (*CapabilityCaller, error) {
	// Find the best plugin for this capability
	pluginName := pr.findBestPluginForCapability(capability)
	if pluginName == "" {
		return nil, fmt.Errorf("capability '%s' is not available in any registered plugin", capability)
	}
	
	plugin, exists := pr.plugins[pluginName]
	if !exists {
		return nil, fmt.Errorf("plugin '%s' not found in registry", pluginName)
	}
	
	return &CapabilityCaller{
		PluginName: pluginName,
		Capability: capability,
		BinaryPath: plugin.BinaryPath,
	}, nil
}

// Call executes the capability with the given arguments
func (cc *CapabilityCaller) Call(ctx context.Context, args []interface{}) (*ResponseWrapper, error) {
	// Convert capability to CLI flag
	operation := strings.SplitN(cc.Capability, ":", 2)[0]
	cliFlag := "--" + operation
	
	// Build command arguments
	cmdArgs := []string{cliFlag}
	for _, arg := range args {
		cmdArgs = append(cmdArgs, fmt.Sprintf("%v", arg))
	}
	
	// Execute the plugin
	cmd := exec.CommandContext(ctx, cc.BinaryPath, cmdArgs...)
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

// findBestPluginForCapability finds the best plugin for a capability
func (pr *PluginRegistry) findBestPluginForCapability(capability string) string {
	candidates := pr.getCapabilityCandidates(capability)
	if len(candidates) == 0 {
		return ""
	}
	
	// Find the candidate with the highest specificity
	bestPlugin := ""
	bestScore := -1
	
	for _, pluginName := range candidates {
		plugin := pr.plugins[pluginName]
		score := pr.calculateCapabilityScore(plugin, capability)
		if score > bestScore {
			bestPlugin = pluginName
			bestScore = score
		}
	}
	
	return bestPlugin
}

// getCapabilityCandidates returns plugins that might support the capability
func (pr *PluginRegistry) getCapabilityCandidates(capability string) []string {
	// Direct match
	if plugins, exists := pr.capabilityIndex[capability]; exists {
		return plugins
	}
	
	// Try wildcard variations
	if strings.Contains(capability, ":") {
		parts := strings.SplitN(capability, ":", 2)
		if len(parts) == 2 {
			wildcardCapability := parts[0] + ":*"
			if plugins, exists := pr.capabilityIndex[wildcardCapability]; exists {
				return plugins
			}
		}
	}
	
	return []string{}
}

// calculateCapabilityScore calculates specificity score for a plugin capability match
func (pr *PluginRegistry) calculateCapabilityScore(plugin *PluginEntry, capability string) int {
	score := 0
	
	// Add specificity score
	for _, cap := range plugin.Capabilities {
		if cap == capability {
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

// ListCapabilities returns all available capabilities
func (pr *PluginRegistry) ListCapabilities() []string {
	capabilities := make([]string, 0, len(pr.capabilityIndex))
	for capability := range pr.capabilityIndex {
		capabilities = append(capabilities, capability)
	}
	return capabilities
}

// PluginCapabilities represents plugin capabilities
type PluginCapabilities struct {
	Capabilities []string `json:"capabilities"`
}

// PluginInfo represents plugin information for --plugin-info output
type PluginInfo struct {
	// Plugin name
	Name string `json:"name"`

	// Plugin version
	Version string `json:"version"`

	// Plugin description
	Description string `json:"description"`

	// Plugin capabilities with file type specificity
	Capabilities *PluginCapabilities `json:"capabilities"`

	// Plugin author/maintainer
	Author *string `json:"author,omitempty"`
}

// NewPluginInfo creates a new plugin info
func NewPluginInfo(name, version, description string, capabilities []string) *PluginInfo {
	return &PluginInfo{
		Name:         name,
		Version:      version,
		Description:  description,
		Capabilities: &PluginCapabilities{Capabilities: capabilities},
	}
}

// WithAuthor sets the author of the plugin
func (pi *PluginInfo) WithAuthor(author string) *PluginInfo {
	pi.Author = &author
	return pi
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

// Can checks if plugin has a specific capability
func (pc *PluginCapabilities) Can(capability string) bool {
	for _, cap := range pc.Capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}

// StandardizedCapabilities are the standard capability names
var StandardizedCapabilities = struct {
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