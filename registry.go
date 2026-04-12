// Package sdk provides registry integration for cartridge execution with caller system
package sdk

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/machinefabric/capdag-go/cap"
	"github.com/machinefabric/capdag-go/urn"
)

// CartridgeRegistry provides cap-based access to cartridges with caller support
type CartridgeRegistry struct {
	cartridges  map[string]*CartridgeEntry
	capIndex map[string][]string
	registry *cap.CapRegistry
	hostImpl *CartridgeCapSet
}

// CartridgeEntry represents a registered cartridge with its capabilities
type CartridgeEntry struct {
	BinaryPath string
	Caps       []string
	Metadata   *CartridgeMetadata
}

// CartridgeMetadata contains metadata about a cartridge
type CartridgeMetadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author,omitempty"`
	Caps        []string `json:"caps"`
}

// CartridgeCapSet implements CapSet interface for cartridge execution
type CartridgeCapSet struct {
	registry *CartridgeRegistry
}

// ExecuteCap implements the CapSet interface for cartridge execution.
// Arguments are identified by media_urn and converted to CLI arguments
// based on the cap definition's argument sources.
func (pch *CartridgeCapSet) ExecuteCap(
	ctx context.Context,
	capUrn string,
	arguments []cap.CapArgumentValue,
) (*cap.HostResult, error) {
	cartridgeName := pch.registry.findBestCartridgeForCap(capUrn)
	if cartridgeName == "" {
		return nil, fmt.Errorf("no cartridge found for cap: %s", capUrn)
	}

	cartridge, exists := pch.registry.cartridges[cartridgeName]
	if !exists {
		return nil, fmt.Errorf("cartridge %s not found", cartridgeName)
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

	cmd := exec.CommandContext(ctx, cartridge.BinaryPath, cmdArgs...)

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
			return nil, fmt.Errorf("cartridge execution failed with stderr: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("cartridge execution failed: %w", err)
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

// NewCartridgeRegistry creates a new cartridge registry with cap support
func NewCartridgeRegistry() (*CartridgeRegistry, error) {
	registry, err := cap.NewCapRegistry()
	if err != nil {
		registry = nil
	}

	pr := &CartridgeRegistry{
		cartridges:  make(map[string]*CartridgeEntry),
		capIndex: make(map[string][]string),
		registry: registry,
	}

	pr.hostImpl = &CartridgeCapSet{registry: pr}
	return pr, nil
}

// RegisterCartridge registers a cartridge with its capabilities
func (pr *CartridgeRegistry) RegisterCartridge(name, binaryPath string, caps []string) {
	entry := &CartridgeEntry{
		BinaryPath: binaryPath,
		Caps:       caps,
		Metadata: &CartridgeMetadata{
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

	pr.cartridges[name] = entry
}

// RegisterCartridgeWithMetadata registers a cartridge with full metadata
func (pr *CartridgeRegistry) RegisterCartridgeWithMetadata(name, binaryPath string, metadata *CartridgeMetadata) {
	entry := &CartridgeEntry{
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

	pr.cartridges[name] = entry
}

// Can checks if a cap is available and returns a CapCaller instance
func (pr *CartridgeRegistry) Can(capUrn string) (*cap.CapCaller, error) {
	cartridgeName := pr.findBestCartridgeForCap(capUrn)
	if cartridgeName == "" {
		return nil, fmt.Errorf("cap '%s' is not available in any registered cartridge", capUrn)
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

		capDefinition = cap.NewCap(capUrnObj, "Cartridge Capability", command)
	}

	caller := cap.NewCapCaller(capUrn, pr.hostImpl, capDefinition)
	return caller, nil
}

// ValidateCartridgeCaps validates all caps in a cartridge against canonical definitions
func (pr *CartridgeRegistry) ValidateCartridgeCaps(caps []*cap.Cap) []error {
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

// GetCartridges returns all registered cartridges
func (pr *CartridgeRegistry) GetCartridges() map[string]*CartridgeEntry {
	result := make(map[string]*CartridgeEntry)
	for name, entry := range pr.cartridges {
		result[name] = entry
	}
	return result
}

// GetCapabilities returns all available capabilities
func (pr *CartridgeRegistry) GetCapabilities() []string {
	caps := make([]string, 0, len(pr.capIndex))
	for c := range pr.capIndex {
		caps = append(caps, c)
	}
	return caps
}

// GetCartridgesForCap returns all cartridges that support a given cap
func (pr *CartridgeRegistry) GetCartridgesForCap(c string) []string {
	if cartridges, exists := pr.capIndex[c]; exists {
		result := make([]string, len(cartridges))
		copy(result, cartridges)
		return result
	}
	return []string{}
}

// findBestCartridgeForCap finds the best cartridge to handle a specific cap
func (pr *CartridgeRegistry) findBestCartridgeForCap(c string) string {
	candidates := pr.getCapCandidates(c)
	if len(candidates) == 0 {
		return ""
	}

	bestCartridge := ""
	bestScore := -1

	for _, cartridgeName := range candidates {
		cartridge := pr.cartridges[cartridgeName]
		score := pr.calculateCapScore(cartridge, c)
		if score > bestScore {
			bestCartridge = cartridgeName
			bestScore = score
		}
	}

	return bestCartridge
}

// getCapCandidates returns cartridges that might support the cap
func (pr *CartridgeRegistry) getCapCandidates(c string) []string {
	if cartridges, exists := pr.capIndex[c]; exists {
		return cartridges
	}

	candidates := make([]string, 0)

	for registeredCap, cartridges := range pr.capIndex {
		if pr.isCapMatch(registeredCap, c) {
			candidates = append(candidates, cartridges...)
		}
	}

	return candidates
}

// isCapMatch checks if a registered cap pattern matches the requested cap
func (pr *CartridgeRegistry) isCapMatch(registeredCap, requestedCap string) bool {
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

// calculateCapScore calculates specificity score for a cartridge cap match
func (pr *CartridgeRegistry) calculateCapScore(cartridge *CartridgeEntry, c string) int {
	score := 0

	for _, cartridgeCap := range cartridge.Caps {
		if cartridgeCap == c {
			score += 100
			break
		} else if pr.isCapMatch(cartridgeCap, c) {
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

