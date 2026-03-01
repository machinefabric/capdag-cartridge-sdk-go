// Example usage of the updated MACINA Plugin SDK with capdag-go integration
package main

import (
	"context"
	"fmt"
	"log"

	sdk "github.com/machinefabric/fgrnd-plugin-sdk-go"
	"github.com/machinefabric/capdag-go/cap"
	"github.com/machinefabric/capdag-go/standard"
)

func main() {
	fmt.Println("MACINA Plugin SDK with CapCaller Integration Example")
	fmt.Println("=================================================")

	registry, err := sdk.NewPluginRegistry()
	if err != nil {
		log.Fatal("Failed to create plugin registry:", err)
	}

	pluginCaps := []string{
		"cap:op=extract;target=metadata;",
		"cap:op=generate;output=binary;target=thumbnail;",
		"cap:op=extract;target=outline;",
		"cap:op=extract;target=pages",
	}

	registry.RegisterPlugin("samplePdfPlugin", "/usr/local/bin/pdfplugin", pluginCaps)

	registry.RegisterPlugin("genericPlugin", "/usr/local/bin/generic", []string{
		"cap:op=extract;target=*;",
	})

	fmt.Printf("\nRegistered capabilities:\n")
	for _, c := range registry.GetCapabilities() {
		fmt.Printf("  - %s\n", c)
	}

	fmt.Printf("\n1. Testing capability: extract metadata\n")
	if _, err := registry.Can("cap:op=extract;target=metadata;"); err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Found capability, got CapCaller\n")
		fmt.Printf("   CapCaller ready for execution\n")
	}

	fmt.Printf("\n2. Testing wildcard capability: extract text\n")
	if _, err := registry.Can("cap:op=extract;target=text;"); err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Wildcard match found, got CapCaller\n")
	}

	fmt.Printf("\n3. Testing unsupported capability\n")
	if _, err := registry.Can("cap:op=unsupported;target=test;"); err != nil {
		fmt.Printf("   Expected error: %v\n", err)
	} else {
		fmt.Printf("   Unexpected: Should have failed\n")
	}

	fmt.Printf("\n4. Example caller usage (mock execution)\n")
	if caller, err := registry.Can("cap:op=extract;target=metadata;"); err == nil {
		fmt.Printf("   CapCaller created for metadata extraction\n")
		fmt.Printf("   Would call: caller.Call(ctx, []CapArgumentValue{...}, registry)\n")
		_ = caller
	}

	fmt.Printf("\n5. Standard capabilities available:\n")
	standardCaps := []string{
		"cap:op=extract_metadata",
		"cap:op=generate_thumbnail",
		"cap:op=extract_outline",
		"cap:op=grind",
	}

	for _, capUrn := range standardCaps {
		if c := sdk.GetStandardCapByUrn(capUrn); c != nil {
			fmt.Printf("   - %s: %s\n", c.UrnString(), c.GetCommand())
		}
	}

	fmt.Printf("\n6. Registered plugins:\n")
	for name, plugin := range registry.GetPlugins() {
		fmt.Printf("   - %s (%s)\n", name, plugin.BinaryPath)
		fmt.Printf("     Capabilities: %v\n", plugin.Caps)
		if plugin.Metadata != nil {
			fmt.Printf("     Metadata: %s v%s\n", plugin.Metadata.Name, plugin.Metadata.Version)
		}
	}

	fmt.Printf("\nIntegration complete - CapCaller system working with plugin registry!\n")
}

// demonstrateCapCallerUsage shows how to use a CapCaller in practice
func demonstrateCapCallerUsage(caller *sdk.CapCaller) {
	ctx := context.Background()

	// Arguments are now identified by media URN
	args := []cap.CapArgumentValue{
		cap.NewCapArgumentValueFromStr(standard.MediaFilePath, "/path/to/document.pdf"),
	}

	response, err := caller.Call(ctx, args, nil)
	if err != nil {
		log.Printf("Execution failed: %v", err)
		return
	}

	if response.IsBinary() {
		fmt.Printf("Received binary response: %d bytes\n", response.Size())
		data := response.AsBytes()
		_ = data
	} else if response.IsJSON() {
		var metadata map[string]interface{}
		if err := response.AsType(&metadata); err != nil {
			log.Printf("Failed to parse JSON: %v", err)
		} else {
			fmt.Printf("Received metadata: %+v\n", metadata)
		}
	} else {
		text, err := response.AsString()
		if err != nil {
			log.Printf("Failed to get text: %v", err)
		} else {
			fmt.Printf("Received text: %s\n", text)
		}
	}
}
