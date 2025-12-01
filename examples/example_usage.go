// Example usage of the updated FMIO Plugin SDK with capns-go integration
package main

import (
	"context"
	"fmt"
	"log"

	sdk "github.com/jowharshamshiri/fmio-plugin-sdk-go"
)

func main() {
	fmt.Println("FMIO Plugin SDK with CapCaller Integration Example")
	fmt.Println("=================================================")

	// Create a new plugin registry
	registry, err := sdk.NewPluginRegistry()
	if err != nil {
		log.Fatal("Failed to create plugin registry:", err)
	}

	// Register a sample plugin with capabilities
	pluginCaps := []string{
		"cap:action=extract;target=metadata;",
		"cap:action=generate;output=binary;target=thumbnail;",
		"cap:action=extract;target=outline;",
		"cap:action=extract;target=pages",
	}

	registry.RegisterPlugin("samplePdfPlugin", "/usr/local/bin/pdfplugin", pluginCaps)

	// Register another plugin with wildcard capability
	registry.RegisterPlugin("genericPlugin", "/usr/local/bin/generic", []string{
		"cap:action=extract;target=*;", // Wildcard - can handle any target
	})

	// Show registered capabilities
	fmt.Printf("\nRegistered capabilities:\n")
	for _, cap := range registry.GetCapabilities() {
		fmt.Printf("  - %s\n", cap)
	}

	// Test capability discovery and caller creation
	fmt.Printf("\n1. Testing capability: extract metadata\n")
	if _, err := registry.Can("cap:action=extract;target=metadata;"); err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Found capability, got CapCaller\n")
		fmt.Printf("   CapCaller ready for execution\n")
	}

	// Test wildcard matching
	fmt.Printf("\n2. Testing wildcard capability: extract text\n")
	if _, err := registry.Can("cap:action=extract;target=text;"); err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Wildcard match found, got CapCaller\n")
	}

	// Test unsupported capability
	fmt.Printf("\n3. Testing unsupported capability\n")
	if _, err := registry.Can("cap:action=unsupported;target=test;"); err != nil {
		fmt.Printf("   ✓ Expected error: %v\n", err)
	} else {
		fmt.Printf("   ✗ Unexpected: Should have failed\n")
	}

	// Demonstrate actual calling (would normally execute real plugin)
	fmt.Printf("\n4. Example caller usage (mock execution)\n")
	if caller, err := registry.Can("cap:action=extract;target=metadata;"); err == nil {
		fmt.Printf("   CapCaller created for metadata extraction\n")
		fmt.Printf("   Would call: caller.Call(ctx, args, namedArgs, stdinData)\n")
		_ = caller // Use the variable
		
		// Example call (commented out since we don't have real plugin)
		/*
		ctx := context.Background()
		response, err := caller.Call(ctx, 
			[]interface{}{"/path/to/document.pdf"}, // positional args
			[]interface{}{},                        // named args  
			nil,                                    // stdin data
		)
		if err != nil {
			fmt.Printf("   Execution error: %v\n", err)
		} else {
			fmt.Printf("   ✓ Response received: %d bytes\n", response.Size())
		}
		*/
	}

	// Show standard capabilities
	fmt.Printf("\n5. Standard capabilities available:\n")
	standardCaps := []string{
		"cap:action=extract;target=metadata;",
		"cap:action=generate;output=binary;target=thumbnail;", 
		"cap:action=extract;target=outline;",
		"cap:action=extract;target=pages",
	}
	
	for _, capUrn := range standardCaps {
		if cap := sdk.GetStandardCapByUrn(capUrn); cap != nil {
			fmt.Printf("   - %s: %s\n", cap.UrnString(), cap.Command)
		}
	}

	// Plugin information
	fmt.Printf("\n6. Registered plugins:\n")
	for name, plugin := range registry.GetPlugins() {
		fmt.Printf("   - %s (%s)\n", name, plugin.BinaryPath)
		fmt.Printf("     Capabilities: %v\n", plugin.Caps)
		if plugin.Metadata != nil {
			fmt.Printf("     Metadata: %s v%s\n", plugin.Metadata.Name, plugin.Metadata.Version)
		}
	}

	fmt.Printf("\n✓ Integration complete - CapCaller system working with plugin registry!\n")
}

// demonstrateCapCallerUsage shows how to use a CapCaller in practice
func demonstrateCapCallerUsage(caller *sdk.CapCaller) {
	ctx := context.Background()
	
	// Example: Extract metadata from a PDF
	positionalArgs := []interface{}{"/path/to/document.pdf"}
	namedArgs := []interface{}{
		map[string]interface{}{"name": "output", "value": "/tmp/metadata.json"},
	}
	
	response, err := caller.Call(ctx, positionalArgs, namedArgs, nil)
	if err != nil {
		log.Printf("Execution failed: %v", err)
		return
	}
	
	// Handle different response types
	if response.IsBinary() {
		fmt.Printf("Received binary response: %d bytes\n", response.Size())
		data := response.AsBytes()
		// Process binary data...
		_ = data
	} else if response.IsJSON() {
		// Parse as structured data
		var metadata map[string]interface{}
		if err := response.AsType(&metadata); err != nil {
			log.Printf("Failed to parse JSON: %v", err)
		} else {
			fmt.Printf("Received metadata: %+v\n", metadata)
		}
	} else {
		// Handle as text
		text, err := response.AsString()
		if err != nil {
			log.Printf("Failed to get text: %v", err)
		} else {
			fmt.Printf("Received text: %s\n", text)
		}
	}
}