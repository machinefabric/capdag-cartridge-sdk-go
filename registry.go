// Package sdk provides registry integration for plugin execution with caller system
package sdk

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/machinefabric/capns-go/cap"
	"github.com/machinefabric/capns-go/urn"
)

// PluginRegistry provides cap-based access to plugins with caller support
type PluginRegistry struct {
	plugins  map[string]*PluginEntry
	capIndex map[string][]string
	registry *cap.CapRegistry
	hostImpl *PluginCapSet
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

// PluginCapSet implements CapSet interface for plugin execution
type PluginCapSet struct {
	registry *PluginRegistry
}

// ExecuteCap implements the CapSet interface for plugin execution.
// Arguments are identified by media_urn and converted to CLI arguments
// based on the cap definition's argument sources.
func (pch *PluginCapSet) ExecuteCap(
	ctx context.Context,
	capUrn string,
	arguments []cap.CapArgumentValue,
) (*cap.HostResult, error) {
	pluginName := pch.registry.findBestPluginForCap(capUrn)
	if pluginName == "" {
		return nil, fmt.Errorf("no plugin found for cap: %s", capUrn)
	}

	plugin, exists := pch.registry.plugins[pluginName]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginName)
	}

	capUrnObj, err := urn.NewCapUrnFromString(capUrn)
	if err != nil {
		return nil, fmt.Errorf("invalid cap URN: %w", err)
	}

	// Build command from cap op
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

	// Look up cap definition to map arguments to CLI args
	var capDef *cap.Cap
	if standardCap := GetStandardCapByUrn(capUrn); standardCap != nil {
		capDef = standardCap
	}

	cmdArgs := []string{command}
	var stdinData []byte

	if capDef != nil {
		// Map each CapArgumentValue to its source using the cap definition
		for _, argVal := range arguments {
			argDef := capDef.FindArgByMediaUrn(argVal.MediaUrn)
			if argDef == nil {
				// No definition found; treat as positional
				valStr, err := argVal.ValueAsStr()
				if err != nil {
					cmdArgs = append(cmdArgs, string(argVal.Value))
				} else {
					cmdArgs = append(cmdArgs, valStr)
				}
				continue
			}

			// Use the first non-stdin source to determine CLI form
			placed := false
			for _, src := range argDef.Sources {
				if src.IsStdin() {
					stdinData = argVal.Value
					placed = true
					break
				}
				if src.IsCliFlag() {
					flag := src.GetCliFlag()
					if flag != nil {
						valStr, err := argVal.ValueAsStr()
						if err != nil {
							valStr = string(argVal.Value)
						}
						cmdArgs = append(cmdArgs, *flag, valStr)
						placed = true
						break
					}
				}
				if src.IsPosition() {
					valStr, err := argVal.ValueAsStr()
					if err != nil {
						valStr = string(argVal.Value)
					}
					cmdArgs = append(cmdArgs, valStr)
					placed = true
					break
				}
			}
			if !placed {
				valStr, err := argVal.ValueAsStr()
				if err != nil {
					cmdArgs = append(cmdArgs, string(argVal.Value))
				} else {
					cmdArgs = append(cmdArgs, valStr)
				}
			}
		}
	} else {
		// No cap definition: pass all argument values as positional
		for _, argVal := range arguments {
			valStr, err := argVal.ValueAsStr()
			if err != nil {
				cmdArgs = append(cmdArgs, string(argVal.Value))
			} else {
				cmdArgs = append(cmdArgs, valStr)
			}
		}
	}

	cmd := exec.CommandContext(ctx, plugin.BinaryPath, cmdArgs...)

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

	isBinary := false
	if outputType, exists := capUrnObj.GetTag("output"); exists && outputType == "binary" {
		isBinary = true
	}

	result := &cap.HostResult{}
	if isBinary {
		result.BinaryOutput = output
	} else {
		result.TextOutput = string(output)
	}

	return result, nil
}

// NewPluginRegistry creates a new plugin registry with cap support
func NewPluginRegistry() (*PluginRegistry, error) {
	registry, err := cap.NewCapRegistry()
	if err != nil {
		registry = nil
	}

	pr := &PluginRegistry{
		plugins:  make(map[string]*PluginEntry),
		capIndex: make(map[string][]string),
		registry: registry,
	}

	pr.hostImpl = &PluginCapSet{registry: pr}
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

	for _, c := range caps {
		if _, exists := pr.capIndex[c]; !exists {
			pr.capIndex[c] = make([]string, 0)
		}
		pr.capIndex[c] = append(pr.capIndex[c], name)
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

	for _, c := range metadata.Caps {
		if _, exists := pr.capIndex[c]; !exists {
			pr.capIndex[c] = make([]string, 0)
		}
		pr.capIndex[c] = append(pr.capIndex[c], name)
	}

	pr.plugins[name] = entry
}

// Can checks if a cap is available and returns a CapCaller instance
func (pr *PluginRegistry) Can(capUrn string) (*cap.CapCaller, error) {
	pluginName := pr.findBestPluginForCap(capUrn)
	if pluginName == "" {
		return nil, fmt.Errorf("cap '%s' is not available in any registered plugin", capUrn)
	}

	var capDefinition *cap.Cap
	if pr.registry != nil {
		if registryCap, err := pr.registry.GetCap(capUrn); err == nil {
			capDefinition = registryCap
		}
	}

	if capDefinition == nil {
		if standardCap := GetStandardCapByUrn(capUrn); standardCap != nil {
			capDefinition = standardCap
		}
	}

	if capDefinition == nil {
		capUrnObj, err := urn.NewCapUrnFromString(capUrn)
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

		capDefinition = cap.NewCap(capUrnObj, "Plugin Capability", command)
	}

	caller := cap.NewCapCaller(capUrn, pr.hostImpl, capDefinition)
	return caller, nil
}

// ValidatePluginCaps validates all caps in a plugin against canonical definitions
func (pr *PluginRegistry) ValidatePluginCaps(caps []*cap.Cap) []error {
	if pr.registry == nil {
		return nil
	}

	var errors []error
	for _, c := range caps {
		if err := cap.ValidateCapCanonical(pr.registry, c); err != nil {
			errors = append(errors, fmt.Errorf("cap %s validation failed: %w", c.UrnString(), err))
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
	for c := range pr.capIndex {
		caps = append(caps, c)
	}
	return caps
}

// GetPluginsForCap returns all plugins that support a given cap
func (pr *PluginRegistry) GetPluginsForCap(c string) []string {
	if plugins, exists := pr.capIndex[c]; exists {
		result := make([]string, len(plugins))
		copy(result, plugins)
		return result
	}
	return []string{}
}

// findBestPluginForCap finds the best plugin to handle a specific cap
func (pr *PluginRegistry) findBestPluginForCap(c string) string {
	candidates := pr.getCapCandidates(c)
	if len(candidates) == 0 {
		return ""
	}

	bestPlugin := ""
	bestScore := -1

	for _, pluginName := range candidates {
		plugin := pr.plugins[pluginName]
		score := pr.calculateCapScore(plugin, c)
		if score > bestScore {
			bestPlugin = pluginName
			bestScore = score
		}
	}

	return bestPlugin
}

// getCapCandidates returns plugins that might support the cap
func (pr *PluginRegistry) getCapCandidates(c string) []string {
	if plugins, exists := pr.capIndex[c]; exists {
		return plugins
	}

	candidates := make([]string, 0)

	for registeredCap, plugins := range pr.capIndex {
		if pr.isCapMatch(registeredCap, c) {
			candidates = append(candidates, plugins...)
		}
	}

	return candidates
}

// isCapMatch checks if a registered cap pattern matches the requested cap
func (pr *PluginRegistry) isCapMatch(registeredCap, requestedCap string) bool {
	if registeredCap == requestedCap {
		return true
	}

	regUrn, err1 := urn.NewCapUrnFromString(registeredCap)
	reqUrn, err2 := urn.NewCapUrnFromString(requestedCap)
	if err1 != nil || err2 != nil {
		return false
	}

	// Use directional matching: request.Accepts(registered) — request is pattern, registered is instance
	return reqUrn.Accepts(regUrn)
}

// calculateCapScore calculates specificity score for a plugin cap match
func (pr *PluginRegistry) calculateCapScore(plugin *PluginEntry, c string) int {
	score := 0

	for _, pluginCap := range plugin.Caps {
		if pluginCap == c {
			score += 100
			break
		} else if pr.isCapMatch(pluginCap, c) {
			score += 50
		}
	}

	return score
}

// GetStandardCapByUrnCanonical returns a standard cap by fetching from registry if available
func GetStandardCapByUrnCanonical(urnStr string) (*cap.Cap, error) {
	if localCap := GetStandardCapByUrn(urnStr); localCap != nil {
		registry, err := cap.NewCapRegistry()
		if err == nil {
			if err := cap.ValidateCapCanonical(registry, localCap); err == nil {
				return localCap, nil
			}
		}
		return localCap, nil
	}

	registry, err := cap.NewCapRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to create registry: %w", err)
	}

	return registry.GetCap(urnStr)
}

// ValidateStandardCaps validates all standard caps against the registry
func ValidateStandardCaps() error {
	registry, err := cap.NewCapRegistry()
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	standardUrns := []string{
		"cap:op=extract_metadata",
		"cap:op=generate_thumbnail",
		"cap:op=extract_outline",
		"cap:op=grind",
	}

	for _, u := range standardUrns {
		if localCap := GetStandardCapByUrn(u); localCap != nil {
			if err := cap.ValidateCapCanonical(registry, localCap); err != nil {
				return fmt.Errorf("standard cap %s validation failed: %w", u, err)
			}
		}
	}

	return nil
}

