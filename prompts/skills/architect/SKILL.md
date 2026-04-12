---
name: architect
description: Reverse-architecture analysis of brownfield repositories. Use this skill whenever the user wants to understand a codebase's architecture, generate C4 diagrams, write architecture reports, or analyze how a repo is organized. This is not just a CLI wrapper — it guides you through a multi-step process of automated extraction + manual deep-dive + synthesis into a report that serves business, tech leads, and new developers. Also use when extending the sdp architect module itself.
version: 2.0.0
---

# @architect — Brownfield Architecture Analysis

**Understand a codebase like a real architect would: examine, diagram, critique, recommend.**

This skill produces architecture reports that answer questions for three audiences:
- **Business**: what does this system do, what are the risks, what are the alternatives
- **Tech lead**: how is it built, where is the tech debt, what's dangerous to change
- **New developer**: where to start, what's where, how things connect

The `sdp architect` CLI does the automated extraction. **You do the thinking.**

---

## Workflow

When the user asks to analyze a repository, follow these steps in order.

### Step 1: Automated Extraction

Run the CLI to get structured data. Use `--section` to avoid huge JSON blobs.

```bash
# Summary first — orient yourself (compact text, <5KB)
sdp architect analyze --no-llm --tier 2 --section summary /path/to/repo

# C4 diagrams (mermaid)
sdp architect c4 --level 1 /path/to/repo    # System context
sdp architect c4 --level 2 /path/to/repo    # Containers

# Deeper data if needed
sdp architect analyze --no-llm --tier 2 --section model --format json /path/to/repo
sdp architect analyze --no-llm --tier 2 --section report --format json /path/to/repo
```

**Important**: flags go BEFORE the repo path (Go `flag` requirement).

What you get from the CLI:
- File/LOC/language stats
- Module boundaries (Maven, Gradle, SBT, npm, Go)
- Import graph (clusters, edges)
- Container detection (Docker, K8s, Compose, modules)
- Dependency graph between modules
- External systems (from deps and infra)
- C4 diagrams (L1, L2, L3)

What the CLI **cannot** give you:
- Architectural patterns and design decisions
- Business context and purpose
- Tech debt analysis and risk assessment
- Quality of code and test coverage insights
- How execution flows through the system
- Critical paths and bottlenecks
- Recommendations

### Step 2: Manual Deep-Dive

The CLI gives you the skeleton. Now you need to understand the body.

**Read these files** (use Read tool or spawn Explore agents in parallel):

1. **README.md / CONTRIBUTING.md** — project purpose, architecture notes
2. **Entry points** — `main()`, `Application`, `SparkContext`, etc. Follow the startup path
3. **Core abstractions** — the 3-5 key classes/interfaces that define the system's model
4. **Config files** — what's configurable tells you what's important
5. **Key directories** — the 2-3 most-imported packages in the import graph

**Search for signals**:
- `grep -r "TODO\|FIXME\|HACK\|XXX" --include="*.{go,java,scala,py,ts}" | wc -l` — tech debt volume
- Look at the biggest files (LOC) — god objects
- Check test directories — what's tested well, what isn't

**Ask yourself**:
- What's the main execution model? (request/response? pipeline? event-driven? batch?)
- What are the key architectural decisions and why were they made?
- What's the most dangerous thing to change?
- What's deprecated but not removed?
- Where are the boundaries between modules? Are they clean?

### Step 3: Write the Report

Use the template below. **Every section is mandatory.** The report should be 500-1500 lines depending on repo size. Write in the language the user communicates in.

---

## Report Template

```markdown
# {Project Name} — Архитектурный отчёт

> **Дата**: YYYY-MM-DD | **Версия**: X.Y | **Инструмент**: sdp architect
> **Для кого**: техлид, архитектор, новый разработчик

---

## Зачем читать этот отчёт

{2-3 sentences: what repo is this, how big, why you'd need this report}

---

## 1. Что это такое — в одном абзаце

{Plain language: what does this system do, who uses it, how is it deployed.
NOT a list of modules. A human explanation.}

---

## 2. Ландшафт — кто с кем разговаривает (L1)

{Mermaid diagram: system context. Users, the system, external systems.
Show data flows, not just boxes. Use arrows with labels.}

```mermaid
graph TB
    ...
```

{Brief explanation of the diagram — who are the actors, what are the external systems}

---

## 3. Из чего состоит — модульная карта (L2)

{Mermaid diagram: containers/modules grouped by layer or domain}

```mermaid
graph TD
    ...
```

{Table: modules grouped by architectural layer, with LOC estimates and purpose}

| Слой | Модули | LOC (≈) | Назначение |
|------|--------|---------|------------|
| ... | ... | ... | ... |

---

## 4. Как работает — архитектура выполнения

### 4.1 Основной путь выполнения

{ASCII or mermaid diagram showing the main execution flow.
From user input → through processing layers → to output.
This is the MOST IMPORTANT section — it explains how the system actually works.}

### 4.2 Ключевые подсистемы

{For each major subsystem (2-4): what it does, how it works,
key files, and why it matters. Not just "Optimizer — optimizes queries"
but HOW it optimizes: what rules, what strategy, what tradeoffs.}

### 4.3 Новые / стратегические компоненты

{Components that represent the project's future direction.
E.g. "Spark Connect" for Spark, "Server Components" for React.}

---

## 5. Архитектурные паттерны

### 5.1 Что сделано хорошо

{Table: pattern | where used | why it's good.
Be specific — name files, classes, interfaces.}

### 5.2 Ключевые архитектурные решения

{List of 3-5 ADRs (Architecture Decision Records) you can infer:
what decision was made, why (or your best guess), and what it costs.}

---

## 6. Где болит — технический долг и риски

### 6.1 Критические риски

{Table: # | problem | severity | location | impact.
Be specific — file names, LOC counts, concrete consequences.}

### 6.2 Архитектурный долг

{Things that are wrong structurally: god objects, deprecated-but-not-removed,
duplicated abstractions, missing boundaries.}

### 6.3 Метрики tech debt

{TODO/FIXME/HACK counts, biggest files, recurring themes in comments.}

---

## 7. Граф зависимостей

```mermaid
graph TD
    ...
```

{Highlight the dependency hubs — modules that everything depends on.
These are the most dangerous to change.}

---

## 8. Тестирование

{Table: test category | count | format | what it covers.
What's well-tested, what's not. Strong and weak sides.}

---

## 9. API и контракты

{What public APIs exist (REST, gRPC, SDK, CLI).
What's stable, what's evolving.}

---

## 10. Рекомендации

### Для техлида

{Table: # | recommendation | priority | reasoning.
Concrete, actionable. "Don't touch X" is a valid recommendation.}

### Для нового разработчика

{Ordered list: read X first, then Y, then Z. With file paths.}

### Для бизнеса

{3-4 bullet points: maturity, risks, strategic direction, alternatives to consider.}

---

## Приложение: Статистика

{Raw numbers from sdp architect output}
```

---

## Quality Checklist

Before delivering the report, verify:

- [ ] **Has mermaid diagrams** — at least L1 (system context), L2 (containers), and dependency graph
- [ ] **Explains HOW it works** — not just WHAT modules exist, but the execution flow
- [ ] **Identifies god objects** — files >2000 LOC that do too much
- [ ] **Counts tech debt** — TODO/FIXME/HACK with actual numbers
- [ ] **Names specific files** — recommendations reference real paths, not abstractions
- [ ] **Serves three audiences** — business gets risks/maturity, tech lead gets patterns/debt, dev gets entry points
- [ ] **Has recommendations** — not just description, but opinions and advice
- [ ] **Goes beyond static analysis** — you read key source files, not just CLI output
- [ ] **Diagrams have arrows with labels** — data flow, not just boxes
- [ ] **Written in user's language** — if user speaks Russian, report is in Russian

---

## CLI Reference

### `sdp architect analyze <repo-path>`

| Flag | Default | Description |
|------|---------|-------------|
| `--tier` | `2` | Detail level: 1 (~2K tokens), 2 (~5-15K), 3 (on-demand) |
| `--format` | `json` | Output: `text`, `json`, `mermaid` |
| `--section` | all | Output only: `profile`, `report`, `model`, `diagrams`, `summary` |
| `--extractors` | all | Comma-separated list to run specific extractors only |
| `--no-llm` | `false` | Deterministic-only mode (no LLM enrichment) |
| `--allow-external-llm` | `false` | Allow LLM calls (overridden by `--no-llm`) |
| `--timeout` | `5m` | Pipeline timeout |
| `--verbose` | `false` | Show per-extractor timing |
| `--skip-git` | `false` | Skip git history extractor |
| `--language` | auto | Force language: `go`, `java`, `python`, `typescript` |

### `sdp architect c4 <repo-path>`

| Flag | Default | Description |
|------|---------|-------------|
| `--level` | all | Diagram level: `1` (context), `2` (container), `3` (component) |
| `--output` | stdout | Output directory for .mmd files |
| `--format` | `mermaid` | Output format: `mermaid`, `json` |

### `sdp architect eval <repo-path>`

| Flag | Description |
|------|-------------|
| `--ground-truth` | Path to ground-truth JSON |

**Important**: all flags must come **before** the repo path.

---

## Extending the Module

### Adding a New Extractor

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

### Code Map

```
internal/architect/
├── pipeline.go          — Pipeline orchestration, BuildReferenceModelFromProfile
├── assembler.go         — ProfileAssembler, tier-gated merging
├── profile.go           — Data model types (CodebaseProfile, InfraInfo, etc.)
├── c4/
│   ├── generator.go     — Deterministic C4 model generation
│   ├── relationship.go  — Relationship inference, package→container mapping
│   ├── render.go        — Mermaid/PlantUML rendering (L1/L2/L3)
│   └── scoring.go       — Confidence scoring
├── extract/
│   ├── registry.go      — DefaultExtractors() list
│   ├── adapters.go      — Language adapters (Go, Java, Python, TS)
│   ├── infra.go         — InfraExtractor (Docker, K8s, Terraform, SBT, Maven)
│   └── ...              — Other extractors (deps, specs, filetree, sql, git, generated)
cmd/sdp/
└── cmd_architect.go     — CLI entry point (analyze, c4, eval subcommands)
tests/architect/
└── infra_test.go        — Integration tests
```
