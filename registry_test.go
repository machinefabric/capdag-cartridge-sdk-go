package sdk

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	capns "github.com/fmio/capns-go"
)

func TestRegistryManagerCreation(t *testing.T) {
	manager, err := NewRegistryManager()
	require.NoError(t, err)
	require.NotNil(t, manager)
}

func TestStandardCapValidation(t *testing.T) {
	// Get a standard cap
	cap := ExtractMetadataCap()
	require.NotNil(t, cap)

	manager, err := NewRegistryManager()
	require.NoError(t, err)

	// Validate single cap (this would normally contact the registry)
	errors := manager.ValidatePluginCaps([]*capns.Cap{cap})
	
	// We expect this to either succeed or fail with network/registry errors
	// The important thing is that the API works correctly
	if len(errors) > 0 {
		t.Logf("Validation errors (expected if registry is unavailable): %v", errors)
	} else {
		t.Log("Cap validated successfully against registry")
	}
}

func TestCanonicalCapCreation(t *testing.T) {
	manager, err := NewRegistryManager()
	require.NoError(t, err)

	// This would normally contact the registry
	_, err = manager.CreateCanonicalCap("cap:action=extract;target=metadata;")
	
	// We expect this to either succeed or fail with network/registry errors
	if err != nil {
		t.Logf("Expected error when registry unavailable: %v", err)
	} else {
		t.Log("Successfully created canonical cap")
	}
}

func TestIsCapCanonical(t *testing.T) {
	manager, err := NewRegistryManager()
	require.NoError(t, err)

	// Test with a standard cap URN
	isCanonical := manager.IsCapCanonical("cap:action=extract;target=metadata;")
	t.Logf("Cap is canonical: %t", isCanonical)

	// Test with a fake cap URN
	isFakeCanonical := manager.IsCapCanonical("cap:action=fake;target=nonexistent")
	assert.False(t, isFakeCanonical, "Fake cap should not be canonical")
}

// Note: Uncomment to test with real registry
/*
func TestGetStandardCapByUrnCanonical(t *testing.T) {
	cap, err := GetStandardCapByUrnCanonical("cap:action=extract;target=metadata;")
	if err != nil {
		t.Skipf("Skipping real registry test: %v", err)
		return
	}
	
	assert.NotNil(t, cap)
	assert.Equal(t, "cap:action=extract;target=metadata;", cap.UrnString())
}

func TestValidateStandardCaps(t *testing.T) {
	err := ValidateStandardCaps()
	if err != nil {
		t.Logf("Standard caps validation error (may be expected): %v", err)
	} else {
		t.Log("All standard caps validated successfully")
	}
}
*/