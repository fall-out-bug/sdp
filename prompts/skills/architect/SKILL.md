---
name: architect
description: Reverse-architecture analysis of source code repositories. Produces C4 architecture models (L1/L2/L3), codebase profiles, and architecture reports. Use this skill whenever the user mentions architecture analysis, C4 diagrams, codebase structure extraction, reverse engineering architecture, or wants to understand how a repository is organized. Also use when debugging or extending the architect module itself.
version: 1.0.0
---

# @architect — Reverse-Architecture Analysis

**Extract C4 architecture models from source code repositories automatically.**

Analyzes a repo's file tree, dependencies, infrastructure, import graphs, and language-specific constructs to produce a structured architecture profile and C4 reference model — no manual diagramming required.

---

## Quick Start

```bash
# Full analysis (deterministic, no LLM)
sdp architect analyze /path/to/repo --no-llm --format json

# C4 diagrams
sdp architect c4 /path/to/repo --level L2 --format mermaid

# Evaluate against ground truth
sdp architect eval /path/to/repo --ground-truth ground-truth.json
```

---

## CLI Reference

### `sdp architect analyze <repo-path>`

Full pipeline run. Produces a codebase profile, architecture report, and C4 reference model.

| Flag | Default | Description |
|------|---------|-------------|
| `--tier` | `2` | Detail level: 1 (~2K tokens), 2 (~5-15K), 3 (on-demand) |
| `--format` | `text` | Output: `text`, `json` |
| `--extractors` | all | Comma-separated list to run specific extractors only |
| `--no-llm` | `false` | Deterministic-only mode (no LLM enrichment) |
| `--allow-external-llm` | `false` | Allow LLM calls (overridden by `--no-llm`) |
| `--timeout` | `120s` | Pipeline timeout |
| `--verbose` | `false` | Show per-extractor timing |
| `--skip-git` | `false` | Skip git history extractor |
| `--language` | auto | Force language: `go`, `java`, `python`, `typescript` |

### `sdp architect c4 <repo-path>`

Generate C4 diagrams from analysis results.

| Flag | Default | Description |
|------|---------|-------------|
| `--level` | `L2` | Diagram level: `L1` (context), `L2` (container), `L3` (component) |
| `--output` | stdout | Output file path |
| `--format` | `mermaid` | Diagram format: `mermaid`, `plantuml`, `json` |

### `sdp architect eval <repo-path>`

Evaluate analysis quality against a ground-truth file.

| Flag | Description |
|------|-------------|
| `--ground-truth` | Path to ground-truth JSON |

---

## Architecture

### 6-Stage Pipeline

```
Extract → Assemble → Filter → Enrich → Model → Output
```

1. **Extract** — 11 extractors run in parallel, each returns a `ProfileFragment`
2. **Assemble** — `ProfileAssembler` merges fragments into `CodebaseProfile` (tier-gated detail)
3. **Filter** — Security filter strips secrets, internal paths, sensitive data
4. **Enrich** — LLM enrichment (Phase 2, optional — skipped with `--no-llm`)
5. **Model** — `BuildReferenceModelFromProfile` creates C4 `ReferenceModel`
6. **Output** — Format as text/JSON, render C4 diagrams

### 11 Extractors

| Extractor | What it detects |
|-----------|----------------|
| `filetree` | Directory structure, file counts, extension distribution |
| `deps` | Package manifests (go.mod, pom.xml, package.json, etc.), notable dependencies |
| `specs` | API specs (OpenAPI, gRPC .proto, GraphQL schemas) |
| `infra` | Docker, Compose, Kubernetes, Terraform, GitHub Actions, module boundaries (Maven/Gradle/SBT/npm) |
| `go` | Go packages, interfaces, structs, import graph |
| `python` | Python modules, classes, imports |
| `java` | Java/Kotlin/Scala packages, classes, Maven/SBT modules, import graph |
| `typescript` | TypeScript/JavaScript modules, exports, imports |
| `git_history` | Commit patterns, active areas, contributor stats |
| `sql` | SQL schemas, migrations, stored procedures |
| `generated` | Auto-generated code detection (protobuf, codegen markers) |

### Supported Languages & Build Systems

| Language | Build Systems | Import Graph |
|----------|--------------|--------------|
| Go | go.mod | Yes — packages, interfaces |
| Java/Kotlin/Scala | Maven (pom.xml), Gradle (build.gradle), SBT (build.sbt) | Yes — packages, module slugs |
| Python | pip (requirements.txt), poetry, setuptools | Yes — modules, classes |
| TypeScript/JS | npm (package.json), yarn, workspaces | Yes — modules, exports |

### C4 Model Levels

- **L1 — System Context**: Actors, the system, external systems, and their relationships
- **L2 — Container**: Deployable units (services, databases, message brokers) and connections
- **L3 — Component**: Internal modules within each container, derived from import graph clusters

### Container Detection Priority

The pipeline discovers containers in this order (first match wins per unit):

1. **Docker/Compose services** — from `docker-compose.yml`, Dockerfiles
2. **Kubernetes workloads** — Deployments, StatefulSets, DaemonSets
3. **Service dependencies** — `depends_on` edges in Compose
4. **Module boundaries** — Maven modules, SBT subprojects, npm workspaces, Go `cmd/` dirs
5. **Import graph clusters** — fallback when no infra files exist
6. **Single "Application"** — last resort if nothing else detected

### Phantom Container Filtering

After container detection, insignificant containers are removed. A container is kept if:
- It comes from a Dockerfile or Compose service
- It's a Maven/Gradle/SBT module
- It has at least one relationship edge (service dep, import)
- It has components assigned from import clusters

---

## Key Data Types

```
CodebaseProfile
├── Name, Metrics (files, LOC, languages)
├── FileTree (structure, extensions)
├── Dependencies (manifests, notable deps, language)
├── Specs[] (OpenAPI, gRPC, GraphQL artifacts)
├── Infra
│   ├── Containers[] (name, type, image, source)
│   ├── Services[] (from→to dependency edges)
│   ├── ModuleBoundaries[] (build system, children)
│   ├── Resources[] (Terraform resources)
│   ├── ExposedPorts[], BaseImages[]
│   └── DeploymentType (kubernetes | docker-compose | bare)
├── ImportGraph
│   └── Clusters[] (id, packages, internal/external edges)
└── Git history, SQL schemas, Generated code info

ReferenceModel (C4)
├── System (name, description)
├── Containers[] (id, name, technology, components[])
├── Relationships[] (from, to, description, technology)
├── Actors[] (id, description)
├── ExternalSystems[] (id, description, technology)
└── Confidence scores per element
```

---

## For AI Agents

### Running Analysis

```bash
# Deterministic analysis (recommended for automated pipelines)
sdp architect analyze /path/to/repo --no-llm --tier 2 --format json > profile.json

# Quick scan for CI
sdp architect analyze /path/to/repo --no-llm --tier 1 --timeout 30s
```

### Interpreting Results

The JSON output has three top-level keys:
- `profile` — raw codebase data (files, deps, infra, imports)
- `report` — architecture report with findings
- `reference_model` — C4 model with containers, relationships, actors

Check `reference_model.containers` for deployable units and `reference_model.relationships` for how they connect. Confidence scores (0.0–1.0) indicate extraction certainty.

### Common Patterns

**Monorepo with build modules** (Maven/SBT/Gradle):
- Containers come from module boundaries
- Import graph clusters map to components within those containers
- Module slugs (e.g., "spark-sql-core") link clusters to containers

**Microservices with Docker Compose**:
- Each Compose service becomes a container
- `depends_on` creates relationship edges
- Database/cache images get typed as "database"/"cache"

**Single-app repos**:
- One "Application" container with components from import clusters
- External systems inferred from cloud SDK dependencies

---

## Extending the Module

### Adding a New Extractor

1. Implement the `Extractor` interface in `internal/architect/extract/`:
   ```go
   type Extractor interface {
       Name() string
       Extract(ctx context.Context, repoRoot string) (*ProfileFragment, error)
   }
   ```
2. Return data in the relevant `ProfileFragment` fields
3. Register in `internal/architect/extract/registry.go` → `DefaultExtractors()`
4. Add tests in `tests/architect/`

### Adding a Language Adapter

Language-specific analysis uses adapters in `internal/architect/extract/adapters.go`:
1. Create an adapter struct implementing the analysis interface
2. Add file extension detection
3. Build import graph clusters from language-specific imports
4. Register the adapter in the appropriate extractor

### Testing

```bash
# Run all architect tests
go test ./internal/architect/... ./tests/architect/... -v

# Run specific test
go test ./tests/architect/ -run TestInfraExtractor_SBTModules -v

# Run with race detector
go test ./internal/architect/... -race
```

---

## Debugging

### Empty C4 Diagram

If L2/L3 diagrams are empty:
1. Check `--tier` level (Tier 1 may omit detail)
2. Run with `--verbose` to see extractor timing — a zero-duration extractor may have skipped
3. Check if the repo has recognizable infra files (Dockerfile, compose, k8s YAML, pom.xml)
4. For L3: verify import graph clusters exist — components come from cluster→container mapping

### Wrong Container Detection

Container detection follows priority order. If wrong containers appear:
1. Check if CI-only containers leak through (they should be filtered by `isCIContainer`)
2. Verify module boundary detection matches the build system
3. Run with `--extractors infra` to isolate infrastructure extraction

### Missing Language Analysis

If a language isn't detected:
1. Check `--language` flag to force detection
2. Verify the language adapter exists in `extract/adapters.go`
3. Check file extension mapping in `assembler.go` → `extToLanguage`

---

## Evaluation & Quality

### Council Review

For major changes, run a council review with multiple LLM evaluators:
1. Generate analysis output: `sdp architect analyze <repo> --no-llm --format json > output.json`
2. Create ground truth or compare against known architecture
3. Have reviewers score: container accuracy, relationship completeness, technology detection

### Quality Metrics

| Metric | Target |
|--------|--------|
| Container detection precision | > 80% |
| Relationship recall | > 70% |
| Technology identification | > 90% |
| Zero false-positive external systems | Yes |

---

## Code Map

```
internal/architect/
├── pipeline.go          — Pipeline orchestration, BuildReferenceModelFromProfile
├── assembler.go         — ProfileAssembler, tier-gated merging
├── profile.go           — Data model types (CodebaseProfile, InfraInfo, etc.)
├── extract_types.go     — ExtractionResult, ProfileFragment
├── security_filter.go   — Secret/path stripping
├── c4/
│   ├── generator.go     — C4 Generate() (deterministic container/component/relationship creation)
│   ├── relationship.go  — Relationship inference, package→container mapping
│   ├── render.go        — Mermaid/PlantUML rendering (L1/L2/L3)
│   └── scoring.go       — Confidence scoring
├── extract/
│   ├── registry.go      — DefaultExtractors() list
│   ├── adapters.go      — Language adapters (Go, Java, Python, TS)
│   ├── infra.go         — InfraExtractor (Docker, K8s, Terraform, SBT, Maven)
│   ├── deps.go          — DependencyExtractor
│   ├── specs.go         — SpecExtractor
│   ├── filetree.go      — FileTreeExtractor
│   ├── sql_extract.go   — SQLExtractor
│   ├── generated.go     — GeneratedCodeExtractor
│   └── git_history.go   — GitHistoryExtractor
cmd/sdp/
└── cmd_architect.go     — CLI entry point (analyze, c4, eval subcommands)
tests/architect/
└── infra_test.go        — Integration tests for extractors
```
