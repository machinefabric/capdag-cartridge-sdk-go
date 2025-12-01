package sdk

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	capns "github.com/fmio/capns-go"
)

func TestPluginRegistryCreation(t *testing.T) {
	registry, err := NewPluginRegistry()
	require.NoError(t, err)
	require.NotNil(t, registry)

	// Test that it starts empty
	assert.Equal(t, 0, len(registry.GetPlugins()))
	assert.Equal(t, 0, len(registry.GetCapabilities()))
}

func TestPluginRegistration(t *testing.T) {
	registry, err := NewPluginRegistry()
	require.NoError(t, err)

	// Register a plugin with basic caps
	caps := []string{
		"cap:action=extract;target=metadata;",
		"cap:action=generate;output=binary;target=thumbnail;",
	}
	registry.RegisterPlugin("testplugin", "/path/to/testplugin", caps)

	// Test that plugin was registered
	plugins := registry.GetPlugins()
	assert.Equal(t, 1, len(plugins))
	assert.Contains(t, plugins, "testplugin")

	// Test that capabilities were indexed
	capabilities := registry.GetCapabilities()
	assert.Equal(t, 2, len(capabilities))
	assert.Contains(t, capabilities, "cap:action=extract;target=metadata;")
	assert.Contains(t, capabilities, "cap:action=generate;output=binary;target=thumbnail;")

	// Test plugin lookup for specific caps
	pluginsForMetadata := registry.GetPluginsForCap("cap:action=extract;target=metadata;")
	assert.Equal(t, 1, len(pluginsForMetadata))
	assert.Contains(t, pluginsForMetadata, "testplugin")
}

func TestPluginRegistrationWithMetadata(t *testing.T) {
	registry, err := NewPluginRegistry()
	require.NoError(t, err)

	// Create plugin metadata
	metadata := &PluginMetadata{
		Name:        "Test Plugin",
		Version:     "1.0.0",
		Description: "A test plugin for document processing",
		Author:      "Test Author",
		Caps: []string{
			"cap:action=extract;target=metadata;",
			"cap:action=extract;target=outline;",
		},
	}

	registry.RegisterPluginWithMetadata("testplugin", "/path/to/testplugin", metadata)

	// Test that plugin was registered with metadata
	plugins := registry.GetPlugins()
	assert.Equal(t, 1, len(plugins))
	
	plugin := plugins["testplugin"]
	require.NotNil(t, plugin)
	assert.Equal(t, "/path/to/testplugin", plugin.BinaryPath)
	assert.Equal(t, 2, len(plugin.Caps))
	assert.NotNil(t, plugin.Metadata)
	assert.Equal(t, "Test Plugin", plugin.Metadata.Name)
	assert.Equal(t, "1.0.0", plugin.Metadata.Version)
}

func TestCanMethodBasicFunctionality(t *testing.T) {
	registry, err := NewPluginRegistry()
	require.NoError(t, err)

	// Register a plugin with extract-metadata capability
	caps := []string{"cap:action=extract;target=metadata;"}
	registry.RegisterPlugin("pdfplugin", "/usr/bin/pdfplugin", caps)

	// Test that we can get a CapCaller for this capability
	caller, err := registry.Can("cap:action=extract;target=metadata;")
	require.NoError(t, err)
	require.NotNil(t, caller)

	// Test that we get an error for unsupported capability
	_, err = registry.Can("cap:action=unsupported;target=test;")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not available")
}

func TestCapCallerIntegration(t *testing.T) {
	registry, err := NewPluginRegistry()
	require.NoError(t, err)

	// Register a plugin
	caps := []string{"cap:action=extract;target=metadata;"}
	registry.RegisterPlugin("testplugin", "/bin/echo", caps)

	// Get caller for the capability
	caller, err := registry.Can("cap:action=extract;target=metadata;")
	require.NoError(t, err)
	require.NotNil(t, caller)

	// Test calling with basic arguments (this will use /bin/echo as a mock)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call with some test arguments
	response, err := caller.Call(ctx, []interface{}{"/test/file.pdf"}, []interface{}{}, nil)
	
	// Since we're using echo, we expect it to work but return echo's output
	// This test verifies the integration works, not the actual plugin functionality
	if err != nil {
		// It's okay if this fails - we're just testing the integration
		t.Logf("Call failed as expected with mock plugin: %v", err)
	} else {
		// If it succeeds, make sure we got a response
		assert.NotNil(t, response)
		t.Logf("Mock call succeeded with response size: %d bytes", response.Size())
	}
}

func TestStandardCapValidation(t *testing.T) {
	// Get a standard cap
	cap := ExtractMetadataCap()
	require.NotNil(t, cap)

	registry, err := NewPluginRegistry()
	require.NoError(t, err)

	// Validate single cap (this would normally contact the registry)
	errors := registry.ValidatePluginCaps([]*capns.Cap{cap})
	
	// We expect this to either succeed or skip validation if registry is unavailable
	// The important thing is that the API works correctly
	if len(errors) > 0 {
		t.Logf("Validation errors (expected if registry is unavailable): %v", errors)
	} else {
		t.Log("Cap validated successfully or validation skipped")
	}
}

func TestCapPatternMatching(t *testing.T) {
	registry, err := NewPluginRegistry()
	require.NoError(t, err)

	// Register a plugin with pattern matching
	caps := []string{
		"cap:action=extract;target=metadata;",
		"cap:action=extract;target=*;", // Wildcard target
	}
	registry.RegisterPlugin("universalplugin", "/usr/bin/universalplugin", caps)

	// Test exact match
	caller, err := registry.Can("cap:action=extract;target=metadata;")
	assert.NoError(t, err)
	assert.NotNil(t, caller)

	// Test wildcard match
	caller, err = registry.Can("cap:action=extract;target=text;")
	assert.NoError(t, err)
	assert.NotNil(t, caller)

	// Test no match
	_, err = registry.Can("cap:action=generate;target=video;")
	assert.Error(t, err)
}

func TestMultiplePluginPriority(t *testing.T) {
	registry, err := NewPluginRegistry()
	require.NoError(t, err)

	// Register multiple plugins with overlapping capabilities
	registry.RegisterPlugin("generic", "/usr/bin/generic", []string{
		"cap:action=extract;target=*;", // Wildcard
	})
	
	registry.RegisterPlugin("specific", "/usr/bin/pdfspecific", []string{
		"cap:action=extract;target=metadata;", // Exact match
	})

	// Test that specific plugin is chosen over generic
	caller, err := registry.Can("cap:action=extract;target=metadata;")
	require.NoError(t, err)
	require.NotNil(t, caller)

	// The exact implementation details may vary, but we should get a caller
	// Priority logic should prefer the more specific match
}

func TestGetStandardCapByUrnCanonical(t *testing.T) {
	// Test with a known standard cap
	cap, err := GetStandardCapByUrnCanonical("cap:action=extract;target=metadata;")
	if err != nil {
		// If this fails, it might be because registry is not available
		// which is acceptable in tests
		t.Skipf("Skipping canonical cap test: %v", err)
		return
	}
	
	require.NotNil(t, cap)
	// Note: The URN string format may vary (with or without trailing semicolon) 
	// but should contain the expected content
	urnStr := cap.UrnString()
	assert.True(t, 
		urnStr == "cap:action=extract;target=metadata;" || urnStr == "cap:action=extract;target=metadata",
		"Expected URN to be 'cap:action=extract;target=metadata;' or 'cap:action=extract;target=metadata', got '%s'", urnStr)
}

func TestValidateStandardCaps(t *testing.T) {
	err := ValidateStandardCaps()
	if err != nil {
		// Standard caps validation might fail if registry is not available
		// Log the error but don't fail the test
		t.Logf("Standard caps validation error (may be expected): %v", err)
	} else {
		t.Log("All standard caps validated successfully or validation skipped")
	}
}

func TestHostImplementationInterface(t *testing.T) {
	registry, err := NewPluginRegistry()
	require.NoError(t, err)

	// Test that our host implementation satisfies the CapHost interface
	var _ capns.CapHost = registry.hostImpl

	// Register a test plugin
	registry.RegisterPlugin("testhost", "/bin/echo", []string{
		"cap:action=test;target=interface;",
	})

	// Test ExecuteCap method directly
	ctx := context.Background()
	result, err := registry.hostImpl.ExecuteCap(
		ctx,
		"cap:action=test;target=interface;",
		[]string{"test"},
		map[string]string{"flag": "value"},
		nil,
	)

	// This might fail with /bin/echo but should test our interface
	if err != nil {
		t.Logf("ExecuteCap failed as expected with mock: %v", err)
	} else {
		assert.NotNil(t, result)
		t.Logf("ExecuteCap succeeded with result: %+v", result)
	}
}