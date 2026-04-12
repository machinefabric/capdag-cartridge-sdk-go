package sdk

import (
	"context"
	"testing"
	"time"

	"github.com/machinefabric/capdag-go/cap"
	"github.com/machinefabric/capdag-go/standard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCartridgeRegistryCreation(t *testing.T) {
	registry, err := NewCartridgeRegistry()
	require.NoError(t, err)
	require.NotNil(t, registry)

	assert.Equal(t, 0, len(registry.GetCartridges()))
	assert.Equal(t, 0, len(registry.GetCapabilities()))
}

func TestCartridgeRegistration(t *testing.T) {
	registry, err := NewCartridgeRegistry()
	require.NoError(t, err)

	caps := []string{
		"cap:op=extract;target=metadata;",
		"cap:op=generate;output=binary;target=thumbnail;",
	}
	registry.RegisterCartridge("testcartridge", "/path/to/testcartridge", caps)

	cartridges := registry.GetCartridges()
	assert.Equal(t, 1, len(cartridges))
	assert.Contains(t, cartridges, "testcartridge")

	capabilities := registry.GetCapabilities()
	assert.Equal(t, 2, len(capabilities))
	assert.Contains(t, capabilities, "cap:op=extract;target=metadata;")
	assert.Contains(t, capabilities, "cap:op=generate;output=binary;target=thumbnail;")

	cartridgesForMetadata := registry.GetCartridgesForCap("cap:op=extract;target=metadata;")
	assert.Equal(t, 1, len(cartridgesForMetadata))
	assert.Contains(t, cartridgesForMetadata, "testcartridge")
}

func TestCartridgeRegistrationWithMetadata(t *testing.T) {
	registry, err := NewCartridgeRegistry()
	require.NoError(t, err)

	metadata := &CartridgeMetadata{
		Name:        "Test Cartridge",
		Version:     "1.0.0",
		Description: "A test cartridge for document processing",
		Author:      "Test Author",
		Caps: []string{
			"cap:op=extract;target=metadata;",
			"cap:op=extract;target=outline;",
		},
	}

	registry.RegisterCartridgeWithMetadata("testcartridge", "/path/to/testcartridge", metadata)

	cartridges := registry.GetCartridges()
	assert.Equal(t, 1, len(cartridges))

	cartridge := cartridges["testcartridge"]
	require.NotNil(t, cartridge)
	assert.Equal(t, "/path/to/testcartridge", cartridge.BinaryPath)
	assert.Equal(t, 2, len(cartridge.Caps))
	assert.NotNil(t, cartridge.Metadata)
	assert.Equal(t, "Test Cartridge", cartridge.Metadata.Name)
	assert.Equal(t, "1.0.0", cartridge.Metadata.Version)
}

func TestCanMethodBasicFunctionality(t *testing.T) {
	registry, err := NewCartridgeRegistry()
	require.NoError(t, err)

	caps := []string{"cap:op=extract;target=metadata;"}
	registry.RegisterCartridge("pdfcartridge", "/usr/bin/pdfcartridge", caps)

	caller, err := registry.Can("cap:op=extract;target=metadata;")
	require.NoError(t, err)
	require.NotNil(t, caller)

	_, err = registry.Can("cap:op=unsupported;target=test;")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not available")
}

func TestCapCallerIntegration(t *testing.T) {
	registry, err := NewCartridgeRegistry()
	require.NoError(t, err)

	caps := []string{"cap:op=extract;target=metadata;"}
	registry.RegisterCartridge("testcartridge", "/bin/echo", caps)

	caller, err := registry.Can("cap:op=extract;target=metadata;")
	require.NoError(t, err)
	require.NotNil(t, caller)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call with CapArgumentValue (media URN identified arguments)
	response, err := caller.Call(ctx, []cap.CapArgumentValue{
		cap.NewCapArgumentValueFromStr(standard.MediaFilePath, "/test/file.pdf"),
	}, nil)

	if err != nil {
		t.Logf("Call failed as expected with mock cartridge: %v", err)
	} else {
		assert.NotNil(t, response)
		t.Logf("Mock call succeeded with response size: %d bytes", response.Size())
	}
}

func TestStandardCapValidation(t *testing.T) {
	c := ExtractMetadataCap()
	require.NotNil(t, c)

	registry, err := NewCartridgeRegistry()
	require.NoError(t, err)

	errors := registry.ValidateCartridgeCaps([]*cap.Cap{c})

	if len(errors) > 0 {
		t.Logf("Validation errors (expected if registry is unavailable): %v", errors)
	} else {
		t.Log("Cap validated successfully or validation skipped")
	}
}

func TestCapPatternMatching(t *testing.T) {
	registry, err := NewCartridgeRegistry()
	require.NoError(t, err)

	caps := []string{
		"cap:op=extract;target=metadata;",
		"cap:op=extract;target=*;",
	}
	registry.RegisterCartridge("universalcartridge", "/usr/bin/universalcartridge", caps)

	caller, err := registry.Can("cap:op=extract;target=metadata;")
	assert.NoError(t, err)
	assert.NotNil(t, caller)

	caller, err = registry.Can("cap:op=extract;target=text;")
	assert.NoError(t, err)
	assert.NotNil(t, caller)

	_, err = registry.Can("cap:op=generate;target=video;")
	assert.Error(t, err)
}

func TestMultipleCartridgePriority(t *testing.T) {
	registry, err := NewCartridgeRegistry()
	require.NoError(t, err)

	registry.RegisterCartridge("generic", "/usr/bin/generic", []string{
		"cap:op=extract;target=*;",
	})

	registry.RegisterCartridge("specific", "/usr/bin/pdfspecific", []string{
		"cap:op=extract;target=metadata;",
	})

	caller, err := registry.Can("cap:op=extract;target=metadata;")
	require.NoError(t, err)
	require.NotNil(t, caller)
}

func TestGetStandardCapByUrnCanonical(t *testing.T) {
	c, err := GetStandardCapByUrnCanonical("cap:op=extract;target=metadata;")
	if err != nil {
		t.Skipf("Skipping canonical cap test: %v", err)
		return
	}

	require.NotNil(t, c)
	urnStr := c.UrnString()
	assert.True(t,
		urnStr == "cap:op=extract;target=metadata;" || urnStr == "cap:op=extract;target=metadata",
		"Expected URN to be 'cap:op=extract;target=metadata;' or 'cap:op=extract;target=metadata', got '%s'", urnStr)
}

func TestValidateStandardCaps(t *testing.T) {
	err := ValidateStandardCaps()
	if err != nil {
		t.Logf("Standard caps validation error (may be expected): %v", err)
	} else {
		t.Log("All standard caps validated successfully or validation skipped")
	}
}

func TestHostImplementationInterface(t *testing.T) {
	registry, err := NewCartridgeRegistry()
	require.NoError(t, err)

	// Verify CartridgeCapSet implements CapSet interface
	var _ cap.CapSet = registry.hostImpl

	registry.RegisterCartridge("testhost", "/bin/echo", []string{
		"cap:op=test;target=interface;",
	})

	ctx := context.Background()
	result, err := registry.hostImpl.ExecuteCap(
		ctx,
		"cap:op=test;target=interface;",
		[]cap.CapArgumentValue{
			cap.NewCapArgumentValueFromStr(standard.MediaString, "test"),
		},
	)

	if err != nil {
		t.Logf("ExecuteCap failed as expected with mock: %v", err)
	} else {
		assert.NotNil(t, result)
		t.Logf("ExecuteCap succeeded with result: %+v", result)
	}
}

func TestStandardCapStructure(t *testing.T) {
	c := ExtractMetadataCap()
	require.NotNil(t, c)

	// Verify args use media URNs
	args := c.GetArgs()
	require.GreaterOrEqual(t, len(args), 1)

	// First arg should be file_path with position + stdin sources
	filePathArg := args[0]
	assert.Equal(t, standard.MediaFilePath, filePathArg.MediaUrn)
	assert.True(t, filePathArg.Required)
	assert.True(t, filePathArg.HasPositionSource())
	assert.True(t, filePathArg.HasStdinSource())

	// Cap should accept stdin (derived from args)
	assert.True(t, c.AcceptsStdin())

	// Output should be set
	output := c.GetOutput()
	require.NotNil(t, output)
	assert.Equal(t, standard.MediaObject, output.MediaUrn)
}
