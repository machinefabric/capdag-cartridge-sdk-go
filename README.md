# CapDAG Cartridge SDK for Go

This public Go module mirrors the Rust `capdag-cartridge-sdk` LLM contract.
It is for Go cartridge authors and library consumers who need the canonical
request, stream, vocabulary, model-information, media-URN, and cap-URN types.

## Layout

- `llm/` — canonical LLM request/stream/vocab/model-info types, media URNs,
  cap URNs, and `BackendForModelSpec` classification.

## Usage

```go
import "github.com/machinefabric/capdag-cartridge-sdk-go/llm"

req := llm.NewGenerationRequestWithDefaults("Hello", "hf:meta-llama/Llama-3.1-8B-Instruct")
fmt.Println(req.ToJSON())

switch llm.BackendForModelSpec(req.ModelSpec) {
case llm.BackendGGUF:  // dispatch to llm.CapLLMInferenceGGUF
case llm.BackendMLX:   // dispatch to llm.CapLLMInferenceMLX
case llm.BackendCandle: // dispatch to llm.CapLLMInferenceCandle
}
```

## Verify changes

```bash
go test ./...
```

The Rust SDK is authoritative for CapDAG-specific wire types. Changes
must be mirrored here with the same substantive numbered tests. Language-neutral
runtime behavior belongs to the [CapDAG specification](../../capdag/docs/01-overview.md).
