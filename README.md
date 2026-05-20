# MachFab Cartridge SDK for Go

Canonical Go types for MachFab LLM cartridges. Mirrors the Rust
`machfab-cartridge-sdk`: all types and constants in the `llm` subpackage
track the capdag media defs one-for-one.

## Layout

- `llm/` — canonical LLM request/stream/vocab/model-info types, media URNs,
  cap URNs, and `BackendForModelSpec` classification.

## Usage

```go
import "github.com/machinefabric/machfab-cartridge-sdk-go/llm"

req := llm.NewGenerationRequestWithDefaults("Hello", "hf:meta-llama/Llama-3.1-8B-Instruct")
fmt.Println(req.ToJSON())

switch llm.BackendForModelSpec(req.ModelSpec) {
case llm.BackendGGUF:  // dispatch to llm.CapLLMInferenceGGUF
case llm.BackendMLX:   // dispatch to llm.CapLLMInferenceMLX
case llm.BackendCandle: // dispatch to llm.CapLLMInferenceCandle
}
```

No Go cartridges exist in the monorepo today — the SDK is kept as a scaffold
so future Go cartridges align with the Rust and Swift SDKs.
