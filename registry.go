// Package sdk provides registry integration for plugin validation
package sdk

import (
	"fmt"
	capns "github.com/fmio/capns-go"
)

// RegistryManager handles cap validation against the canonical registry
type RegistryManager struct {
	registry *capns.CapRegistry
}

// NewRegistryManager creates a new registry manager for plugin validation
func NewRegistryManager() (*RegistryManager, error) {
	registry, err := capns.NewCapRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to create registry: %w", err)
	}
	
	return &RegistryManager{
		registry: registry,
	}, nil
}

// ValidatePluginCaps validates all caps in a plugin against canonical definitions
func (rm *RegistryManager) ValidatePluginCaps(caps []*capns.Cap) []error {
	var errors []error
	
	for _, cap := range caps {
		if err := capns.ValidateCapCanonical(rm.registry, cap); err != nil {
			errors = append(errors, fmt.Errorf("cap %s validation failed: %w", cap.UrnString(), err))
		}
	}
	
	return errors
}

// CreateCanonicalCap creates a cap from its canonical registry definition
func (rm *RegistryManager) CreateCanonicalCap(urn string) (*capns.Cap, error) {
	return rm.registry.GetCap(urn)
}

// IsCapCanonical checks if a cap URN has a canonical definition
func (rm *RegistryManager) IsCapCanonical(urn string) bool {
	return rm.registry.CapExists(urn)
}

// GetStandardCapByUrnCanonical returns a standard cap by fetching from registry if available
func GetStandardCapByUrnCanonical(urnStr string) (*capns.Cap, error) {
	// First try to get from local standard caps
	if localCap := GetStandardCapByUrn(urnStr); localCap != nil {
		// Validate against registry if available
		rm, err := NewRegistryManager()
		if err == nil {
			if err := capns.ValidateCapCanonical(rm.registry, localCap); err == nil {
				return localCap, nil
			}
		}
		// Return local cap even if registry validation fails
		return localCap, nil
	}
	
	// Try to get from registry
	rm, err := NewRegistryManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create registry manager: %w", err)
	}
	
	return rm.CreateCanonicalCap(urnStr)
}

// ValidateStandardCaps validates all standard caps against the registry
func ValidateStandardCaps() error {
	rm, err := NewRegistryManager()
	if err != nil {
		return fmt.Errorf("failed to create registry manager: %w", err)
	}
	
	standardUrns := []string{
		"cap:action=extract;target=metadata;",
		"cap:action=generate;output=binary;target=thumbnail;",
		"cap:action=extract;target=outline;",
		"cap:action=extract;target=pages",
	}
	
	for _, urn := range standardUrns {
		if localCap := GetStandardCapByUrn(urn); localCap != nil {
			if err := capns.ValidateCapCanonical(rm.registry, localCap); err != nil {
				return fmt.Errorf("standard cap %s validation failed: %w", urn, err)
			}
		}
	}
	
	return nil
}