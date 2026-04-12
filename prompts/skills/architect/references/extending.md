# Extending the Architect Module

## Adding a New Extractor

1. Implement `Extractor` interface in `internal/architect/extract/`:
   ```go
   type Extractor interface {
       Name() string
       Extract(ctx context.Context, repoRoot string) (*ProfileFragment, error)
   }
   ```
2. Return data in relevant `ProfileFragment` fields
3. Register in `internal/architect/extract/registry.go` → `DefaultExtractors()`
4. Add tests in `tests/architect/`

## Adding a Language Adapter

Language-specific analysis uses adapters in `internal/architect/extract/adapters.go`:
1. Create an adapter struct implementing the analysis interface
2. Add file extension detection
3. Build import graph clusters from language-specific imports
4. Register the adapter in the appropriate extractor

## Pipeline Architecture

```
Extract → Assemble → Filter → Enrich → Model → Output
```

1. **Extract** — 11 extractors run in parallel, each returns a `ProfileFragment`
2. **Assemble** — `ProfileAssembler` merges fragments into `CodebaseProfile` (tier-gated)
3. **Filter** — Security filter strips secrets, internal paths
4. **Enrich** — LLM enrichment (optional, `--no-llm` skips)
5. **Model** — `BuildReferenceModelFromProfile` creates C4 `ReferenceModel`
6. **Output** — Format as text/JSON, render C4 diagrams

## 11 Extractors

| Extractor | Detects |
|-----------|---------|
| `filetree` | Directory structure, file counts, extensions |
| `deps` | Package manifests, notable dependencies |
| `specs` | OpenAPI, gRPC .proto, GraphQL schemas |
| `infra` | Docker, Compose, K8s, Terraform, GitHub Actions, module boundaries |
| `go` | Go packages, interfaces, structs, import graph |
| `python` | Python modules, classes, imports |
| `java` | Java/Kotlin/Scala packages, Maven/SBT modules, import graph |
| `typescript` | TypeScript/JS modules, exports, imports |
| `git_history` | Commit patterns, active areas, contributors |
| `sql` | SQL schemas, migrations, stored procedures |
| `generated` | Auto-generated code detection |

## Code Map

```
internal/architect/
├── pipeline.go          — Pipeline orchestration, BuildReferenceModelFromProfile
├── assembler.go         — ProfileAssembler, tier-gated merging
├── profile.go           — Data model types
├── c4/
│   ├── generator.go     — Deterministic C4 model generation
│   ├── relationship.go  — Relationship inference
│   ├── render.go        — Mermaid/PlantUML rendering
│   └── scoring.go       — Confidence scoring
├── extract/
│   ├── registry.go      — DefaultExtractors() list
│   ├── adapters.go      — Language adapters
│   ├── infra.go         — InfraExtractor
│   └── ...              — Other extractors
cmd/sdp/
└── cmd_architect.go     — CLI entry point
tests/architect/
└── infra_test.go        — Integration tests
```
