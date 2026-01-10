// Package sdk provides registry integration for plugin execution with caller system
package sdk

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	capns "github.com/fgnd/cap-sdk-go"
)

// PluginRegistry provides cap-based access to plugins with caller support
type PluginRegistry struct {
	plugins   map[string]*PluginEntry
	capIndex  map[string][]string
	registry  *capns.CapRegistry
	hostImpl  *PluginCapHost
}

// PluginEntry represents a registered plugin with its capabilities
type PluginEntry struct {
	BinaryPath string
	Caps       []string
	Metadata   *PluginMetadata
}

// PluginMetadata contains metadata about a plugin
type PluginMetadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author,omitempty"`
	Caps        []string `json:"caps"`
}

// PluginCapHost implements CapHost interface for plugin execution
type PluginCapHost struct {
	registry *PluginRegistry
}

// ExecuteCap implements the CapHost interface for plugin execution
func (pch *PluginCapHost) ExecuteCap(
	ctx context.Context,
	capUrn string,
	positionalArgs []string,
	namedArgs map[string]string,
	stdinData []byte,
) (*capns.HostResult, error) {
	// Find the plugin that can handle this cap
	pluginName := pch.registry.findBestPluginForCap(capUrn)
	if pluginName == "" {
		return nil, fmt.Errorf("no plugin found for cap: %s", capUrn)
	}

	plugin, exists := pch.registry.plugins[pluginName]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginName)
	}

	// Parse cap URN to extract command
	capUrnObj, err := capns.NewCapUrnFromString(capUrn)
	if err != nil {
		return nil, fmt.Errorf("invalid cap URN: %w", err)
	}

	// Build command based on cap op
	var command string
	if op, exists := capUrnObj.GetTag("op"); exists {
		if target, targetExists := capUrnObj.GetTag("target"); targetExists {
			command = fmt.Sprintf("--%s-%s", op, target)
		} else {
			command = fmt.Sprintf("--%s", op)
		}
	} else {
		return nil, fmt.Errorf("cap URN missing op tag: %s", capUrn)
	}

	// Build full command arguments
	cmdArgs := []string{command}

	// Add positional arguments
	cmdArgs = append(cmdArgs, positionalArgs...)

	// Add named arguments
	for name, value := range namedArgs {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--%s", name), value)
	}

	// Execute the plugin
	cmd := exec.CommandContext(ctx, plugin.BinaryPath, cmdArgs...)

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
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("plugin execution failed with stderr: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("plugin execution failed: %w", err)
	}

	// Determine if output is binary based on cap URN
	isBinary := false
	if outputType, exists := capUrnObj.GetTag("output"); exists && outputType == "binary" {
		isBinary = true
	}

	result := &capns.HostResult{}
	if isBinary {
		result.BinaryOutput = output
	} else {
		result.TextOutput = string(output)
	}

	return result, nil
}

// NewPluginRegistry creates a new plugin registry with cap support
func NewPluginRegistry() (*PluginRegistry, error) {
	registry, err := capns.NewCapRegistry()
	if err != nil {
		// Allow registry creation to fail - we can still work with local definitions
		registry = nil
	}

	pr := &PluginRegistry{
		plugins:  make(map[string]*PluginEntry),
		capIndex: make(map[string][]string),
		registry: registry,
	}

	pr.hostImpl = &PluginCapHost{registry: pr}
	return pr, nil
}

// RegisterPlugin registers a plugin with its capabilities
func (pr *PluginRegistry) RegisterPlugin(name, binaryPath string, caps []string) {
	entry := &PluginEntry{
		BinaryPath: binaryPath,
		Caps:       caps,
		Metadata: &PluginMetadata{
			Name: name,
			Caps: caps,
		},
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

// RegisterPluginWithMetadata registers a plugin with full metadata
func (pr *PluginRegistry) RegisterPluginWithMetadata(name, binaryPath string, metadata *PluginMetadata) {
	entry := &PluginEntry{
		BinaryPath: binaryPath,
		Caps:       metadata.Caps,
		Metadata:   metadata,
	}

	// Update cap index
	for _, cap := range metadata.Caps {
		if _, exists := pr.capIndex[cap]; !exists {
			pr.capIndex[cap] = make([]string, 0)
		}
		pr.capIndex[cap] = append(pr.capIndex[cap], name)
	}

	pr.plugins[name] = entry
}

// Can checks if a cap is available and returns a CapCaller instance
func (pr *PluginRegistry) Can(capUrn string) (*capns.CapCaller, error) {
	// Find the best plugin for this cap
	pluginName := pr.findBestPluginForCap(capUrn)
	if pluginName == "" {
		return nil, fmt.Errorf("cap '%s' is not available in any registered plugin", capUrn)
	}

	// Get cap definition (from registry if available, or create basic definition)
	var capDefinition *capns.Cap
	if pr.registry != nil {
		if registryCap, err := pr.registry.GetCap(capUrn); err == nil {
			capDefinition = registryCap
		}
	}

	// If no registry definition, try to get from standard caps
	if capDefinition == nil {
		if standardCap := GetStandardCapByUrn(capUrn); standardCap != nil {
			capDefinition = standardCap
		}
	}

	// If still no definition, create a basic one
	if capDefinition == nil {
		capUrnObj, err := capns.NewCapUrnFromString(capUrn)
		if err != nil {
			return nil, fmt.Errorf("invalid cap URN: %w", err)
		}

		var command string
		if op, exists := capUrnObj.GetTag("op"); exists {
			if target, targetExists := capUrnObj.GetTag("target"); targetExists {
				command = fmt.Sprintf("%s-%s", op, target)
			} else {
				command = op
			}
		} else {
			command = "unknown"
		}

		capDefinition = capns.NewCap(capUrnObj, "Plugin Capability", command)
	}

	// Create and return CapCaller
	caller := capns.NewCapCaller(capUrn, pr.hostImpl, capDefinition)
	return caller, nil
}

// ValidatePluginCaps validates all caps in a plugin against canonical definitions
func (pr *PluginRegistry) ValidatePluginCaps(caps []*capns.Cap) []error {
	if pr.registry == nil {
		return nil // Skip validation if no registry available
	}

	var errors []error

	for _, cap := range caps {
		if err := capns.ValidateCapCanonical(pr.registry, cap); err != nil {
			errors = append(errors, fmt.Errorf("cap %s validation failed: %w", cap.UrnString(), err))
		}
	}

	return errors
}

// GetPlugins returns all registered plugins
func (pr *PluginRegistry) GetPlugins() map[string]*PluginEntry {
	result := make(map[string]*PluginEntry)
	for name, entry := range pr.plugins {
		result[name] = entry
	}
	return result
}

// GetCapabilities returns all available capabilities
func (pr *PluginRegistry) GetCapabilities() []string {
	caps := make([]string, 0, len(pr.capIndex))
	for cap := range pr.capIndex {
		caps = append(caps, cap)
	}
	return caps
}

// GetPluginsForCap returns all plugins that support a given cap
func (pr *PluginRegistry) GetPluginsForCap(cap string) []string {
	if plugins, exists := pr.capIndex[cap]; exists {
		result := make([]string, len(plugins))
		copy(result, plugins)
		return result
	}
	return []string{}
}

// Internal helper methods

// findBestPluginForCap finds the best plugin to handle a specific cap
func (pr *PluginRegistry) findBestPluginForCap(cap string) string {
	candidates := pr.getCapCandidates(cap)
	if len(candidates) == 0 {
		return ""
	}

	// Find the candidate with the highest specificity score
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

	// Try variations and partial matches
	candidates := make([]string, 0)

	// For each registered cap, check if it's a match or pattern
	for registeredCap, plugins := range pr.capIndex {
		if pr.isCapMatch(registeredCap, cap) {
			candidates = append(candidates, plugins...)
		}
	}

	return candidates
}

// isCapMatch checks if a registered cap pattern matches the requested cap
func (pr *PluginRegistry) isCapMatch(registeredCap, requestedCap string) bool {
	// Exact match
	if registeredCap == requestedCap {
		return true
	}

	// Parse both cap URNs for pattern matching
	_, err1 := capns.NewCapUrnFromString(registeredCap)
	_, err2 := capns.NewCapUrnFromString(requestedCap)
	if err1 != nil || err2 != nil {
		return false
	}

	// For pattern matching, we need to check if the registered cap pattern matches the requested cap
	// Since we don't have direct access to tags, we'll use the URN string representation
	// and implement a simpler pattern matching approach
	
	// Check if registered cap is a wildcard pattern
	if strings.Contains(registeredCap, "*") {
		// For simple wildcard matching, check if the non-wildcard parts match
		// Convert both to comparable forms by removing wildcards and comparing structure
		registeredParts := strings.Split(strings.TrimSuffix(registeredCap, ";"), ";")
		requestedParts := strings.Split(strings.TrimSuffix(requestedCap, ";"), ";")
		
		if len(registeredParts) != len(requestedParts) {
			return false
		}
		
		for i, regPart := range registeredParts {
			if !strings.Contains(regPart, "=") || !strings.Contains(requestedParts[i], "=") {
				continue
			}
			
			regKey := strings.Split(regPart, "=")[0]
			regValue := strings.Split(regPart, "=")[1]
			reqKey := strings.Split(requestedParts[i], "=")[0]
			reqValue := strings.Split(requestedParts[i], "=")[1]
			
			// Keys must match exactly
			if regKey != reqKey {
				return false
			}
			
			// Values can be wildcards
			if regValue != "*" && regValue != reqValue {
				return false
			}
		}
		
		return true
	}

	return false
}

// calculateCapScore calculates specificity score for a plugin cap match
func (pr *PluginRegistry) calculateCapScore(plugin *PluginEntry, cap string) int {
	score := 0

	for _, pluginCap := range plugin.Caps {
		if pluginCap == cap {
			// Exact match gets highest score
			score += 100
			break
		} else if pr.isCapMatch(pluginCap, cap) {
			// Pattern match gets lower score
			score += 50
		}
	}

	return score
}

// GetStandardCapByUrnCanonical returns a standard cap by fetching from registry if available
func GetStandardCapByUrnCanonical(urnStr string) (*capns.Cap, error) {
	// First try to get from local standard caps
	if localCap := GetStandardCapByUrn(urnStr); localCap != nil {
		// Try to validate against registry if available
		registry, err := capns.NewCapRegistry()
		if err == nil {
			if err := capns.ValidateCapCanonical(registry, localCap); err == nil {
				return localCap, nil
			}
		}
		// Return local cap even if registry validation fails
		return localCap, nil
	}

	// Try to get from registry
	registry, err := capns.NewCapRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to create registry: %w", err)
	}

	return registry.GetCap(urnStr)
}

// ValidateStandardCaps validates all standard caps against the registry
func ValidateStandardCaps() error {
	registry, err := capns.NewCapRegistry()
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	standardUrns := []string{
		"cap:op=extract_metadata",
		"cap:op=generate_thumbnail",
		"cap:op=extract_outline",
		"cap:op=grind",
	}

	for _, urn := range standardUrns {
		if localCap := GetStandardCapByUrn(urn); localCap != nil {
			if err := capns.ValidateCapCanonical(registry, localCap); err != nil {
				return fmt.Errorf("standard cap %s validation failed: %w", urn, err)
			}
		}
	}

	return nil
}