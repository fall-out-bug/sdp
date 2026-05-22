# Reality OSS Spec

## Status

Draft.

## Purpose

`@reality` in OSS is a local, evidence-first reverse-engineering workflow for a single repository. Its job is to recover a usable baseline of the system as it exists today and convert that baseline into SDP-ready artifacts.

It must answer:

- What is implemented?
- How is the code organized?
- What are the major boundaries, dependencies, and entrypoints?
- Where do code and documentation drift?
- Is this repository ready for agent-driven work under SDP?

It must not pretend to replace a consulting-grade technical audit.

## Goals

1. Analyze one repository without external services.
2. Reconstruct a baseline feature inventory from code, tests, configs, and in-repo docs.
3. Recover a practical architecture map: modules, boundaries, entrypoints, integrations, data stores.
4. Produce a code-quality and test-quality baseline.
5. Detect obvious documentation drift and unsupported claims.
6. Prepare bootstrap artifacts for future SDP workstreams.

## Non-Goals

- Multi-repo system reconstruction.
- Native connectors to Confluence, Jira, Notion, or wiki systems.
- Persistent cross-run repository memory.
- Consulting-grade intent recovery from historical and external evidence.
- Heavy adversarial review orchestration.
- Automatic refactoring or remediation during analysis.

## Inputs

Required inputs:

- repository working tree
- source code
- tests
- local configs and manifests
- local git history
- in-repo documentation

Optional inputs:

- manually provided local exports of external docs
- user-provided product focus or analysis scope

## Source Priority Model

OSS `@reality` must rank sources in this order:

1. code and executable configuration
2. tests and fixtures
3. runtime manifests and deployment configuration
4. local repository docs
5. user notes or attached exported docs

Every conclusion must be tagged as one of:

- `observed`
- `documented`
- `inferred`
- `conflicted`
- `unknown`

## Modes

| Mode | Purpose | Expected Depth |
|---|---|---|
| `--quick` | fast baseline scan | same artifact families, reduced finding detail |
| `--deep` | full single-repo baseline | feature inventory, architecture map, quality, drift, readiness |
| `--focus=architecture|quality|testing|docs|security` | domain-specific reporting emphasis | one area, deeper local review |
| `--bootstrap-sdp` | generate agent-readiness outputs | workstream and scope bootstrap with starter recommendations |

## Workflow

### Phase 1: Intake

- detect language, framework, build tool, test tool, runtime shape
- detect repository scale: file count, LOC, package/module count
- classify repo type: app, library, protocol, infra, mixed

### Phase 2: Local Index Build

Build a temporary local index for the current run:

- file inventory
- import/dependency graph
- test inventory
- config inventory
- docs inventory
- entrypoint candidates

The OSS index is run-scoped, not persistent across sessions.

### Phase 3: Baseline Analysis

Required analyzers:

- structure analyzer
- architecture analyzer
- integration analyzer
- quality analyzer
- testing analyzer
- documentation drift analyzer

### Phase 4: Limited Cross-Check Review

OSS must still avoid single-agent overconfidence.

Minimum review structure:

1. primary source-first pass generates findings from code and executable config
2. secondary heuristic review checks findings against tests, manifests, and docs
3. synthesis keeps supported findings and downgrades weak claims

If the review pass cannot confirm a claim, the claim must be downgraded to `inferred`, `conflicted`, or `unknown`.

Spawning multiple LLM agents is optional and not required for OSS.

### Phase 5: Synthesis

Produce:

- feature baseline
- architecture baseline
- risk baseline
- readiness baseline
- suggested first SDP workstreams

## Required Outputs

Human-readable outputs:

- `docs/reality/summary.md`
- `docs/reality/architecture.md`
- `docs/reality/quality.md`
- `docs/reality/bootstrap.md`

Machine-readable outputs:

- `.sdp/reality/reality-summary.json`
- `.sdp/reality/feature-inventory.json`
- `.sdp/reality/architecture-map.json`
- `.sdp/reality/integration-map.json`
- `.sdp/reality/quality-report.json`
- `.sdp/reality/drift-report.json`
- `.sdp/reality/readiness-report.json`

## Minimum Artifact Semantics

### Feature Inventory

Each feature candidate should include:

- feature id or synthetic id
- title
- description
- evidence paths
- status: implemented, partial, candidate, dead
- confidence

### Architecture Map

Must include at least:

- top-level modules
- entrypoints
- internal boundaries
- external integrations
- data stores
- high-coupling zones

### Readiness Report

Must answer:

- can safe scopes be defined?
- are there tests that protect future changes?
- what hotspots block autonomous work?
- what should be isolated first?

## Done Criteria

OSS `@reality` is done when it can, for one repository:

1. identify top-level code areas and entrypoints
2. produce a baseline feature inventory with evidence and confidence
3. produce an architecture map with integrations and boundaries
4. report major code-quality and documentation-drift findings
5. produce a repository readiness verdict:
   - `ready`
   - `ready_with_constraints`
   - `not_ready`
6. emit bootstrap recommendations for SDP workstreams

## Agent Readiness Rules

OSS readiness focuses on future SDP execution. It must evaluate:

- safe scope extraction
- module size and coupling
- existence of verification surfaces
- hidden dependencies
- large files and low-test zones
- unstable boundaries

## Open-Core Boundary

OSS owns:

- single-repo baseline analysis
- open artifact formats
- local analyzers
- limited cross-check review inside one run
- SDP bootstrap starter outputs

Anything requiring persistent memory, native external connectors, multi-repo orchestration, or consulting-grade adversarial synthesis belongs in `reality-pro`.
