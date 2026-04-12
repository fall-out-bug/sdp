---
name: architect
description: Brownfield architecture analysis — understand a codebase like a senior architect would. Produces architecture reports with C4 diagrams, execution flow analysis, tech debt assessment, and actionable recommendations for business, tech leads, and developers. Use this skill whenever the user mentions architecture analysis, codebase understanding, reverse engineering, C4 diagrams, or "what is this repo". This is NOT a CLI wrapper — it orchestrates automated extraction + parallel deep-dive + synthesis.
version: 3.0.0
---

# @architect — Brownfield Architecture Analysis

You are a senior architect who just walked into a new codebase. Your job: examine it, understand it, diagram it, critique it, and explain it to three audiences — business, tech lead, and new developer.

The `sdp architect` CLI extracts structural data. **You provide the understanding.**

---

## Mandatory Steps

Follow these steps in order. Do not skip steps. Each step has a verification gate — do not proceed until the gate passes.

### Step 1: Automated Extraction (2-5 min)

Run CLI commands to get structural data. Use `--section` to keep output manageable.

```bash
# 1a. Summary — orient yourself first
sdp architect analyze --no-llm --tier 2 --section summary /path/to/repo

# 1b. C4 diagrams
sdp architect c4 --level 1 /path/to/repo   # System context (L1)
sdp architect c4 --level 2 /path/to/repo   # Containers (L2)

# 1c. Model data (only if you need container/relationship details)
sdp architect analyze --no-llm --tier 2 --section model --format json /path/to/repo

# 1d. Report data (patterns, risks, styles)
sdp architect analyze --no-llm --tier 2 --section report --format json /path/to/repo
```

**Flags go BEFORE the repo path** (Go `flag` requirement).

**Gate**: You have file count, LOC, languages, module count, container count, and at least L1 diagram. If CLI fails or returns empty, fall back to manual exploration in Step 2.

**Evidence tag**: Everything from this step is `[EXTRACTED]` — machine-verified structural data.

---

### Step 2: Parallel Deep-Dive (the step agents skip)

This is where the real architecture analysis happens. The CLI gives you bones — now you find the muscle, nerves, and scars.

**Spawn 3-4 Explore agents in the SAME message** (parallel, not sequential):

**Agent A — Identity & Purpose:**
```
Read README.md, CONTRIBUTING.md, and any docs/ folder in /path/to/repo.
What is this project? Who uses it? What problem does it solve?
What's the deployment model (library, service, CLI, platform)?
Report in <10 lines.
```

**Agent B — Execution Architecture:**
```
Find the main entry point(s) in /path/to/repo:
- Look for main(), Application, App, Server, cli, cmd/ directories
- Trace the startup path: what gets initialized, in what order
- Identify the core abstractions (3-5 key interfaces/classes that define the domain model)
- What's the execution model: request/response, pipeline, event-driven, batch, actor?
Find the 5 largest files by code (not generated): these are likely god objects.
Report: entry points, core abstractions with file paths, execution model, top 5 largest files.
```

**Agent C — Patterns & Decisions:**
```
In /path/to/repo, search for architectural patterns:
- Plugin/SPI patterns: grep for "interface.*Plugin\|trait.*Provider\|abstract.*Factory"
- Configuration: what's configurable? (properties files, env vars, CLI flags)
- Serialization: how does data cross boundaries?
- Error handling: centralized or scattered?
- Design patterns: Strategy, Observer, Builder, etc. — evidence, not guesses
Count TODO/FIXME/HACK/XXX markers: `grep -r "TODO\|FIXME\|HACK\|XXX" --include="*.{go,java,scala,py,ts,js,rs}" /path/to/repo | wc -l`
Report: patterns found with file paths, tech debt count, top themes.
```

**Agent D — Testing & API Surface:**
```
In /path/to/repo:
- Count test files by type: unit (*_test.go, *Test.java, *Suite.scala, test_*.py, *.test.ts)
- What testing frameworks? (testify, JUnit, ScalaTest, pytest, Jest, etc.)
- Are there integration/e2e tests? Where?
- What public APIs exist? (REST endpoints, gRPC .proto files, CLI commands, SDK exports)
- What's the CI setup? (.github/workflows/, Jenkinsfile, .gitlab-ci.yml)
Report: test counts by category, frameworks, API surface, CI pipeline.
```

**Gate**: You have answers from at least 3 of 4 agents. You know: what the project IS, how execution flows, what patterns are used, and what the test/API surface looks like.

**Evidence tags**:
- Facts from README, config files, test counts → `[EXTRACTED]`
- Patterns inferred from code structure → `[INFERRED]` (state confidence: high/medium/low)
- Guesses about purpose or intent → `[AMBIGUOUS]` (flag for reader)

---

### Step 3: Synthesis — Write the Report

Now combine CLI data + deep-dive findings into a coherent architecture report. Use the template in `references/report-template.md` (read it now if you haven't).

**Key rules for synthesis:**

1. **Tell a story, not a data dump.** Start with "what is this" in plain language. Then zoom in.
2. **Every diagram needs narration.** A mermaid graph without explanation is noise.
3. **Every claim needs a file path.** "The optimizer uses rule-based approach" → WHERE? Which file?
4. **Label your evidence.** Use `[EXTRACTED]`, `[INFERRED]`, `[AMBIGUOUS]` so readers know what's certain.
5. **Recommendations must be actionable.** "Improve test coverage" is useless. "Add integration tests for `sql/catalyst/optimizer/` — currently 0 tests for 44 rule files" is actionable.

**Mermaid diagrams — minimum 3:**
- L1 system context (actors + system + external systems)
- L2 module map (containers grouped by architectural layer)
- Dependency graph (who depends on whom, highlight hubs)

Write additional diagrams for execution flow if the system is complex enough.

**Gate**: Report has all 10 sections from the template, at least 3 mermaid diagrams, identifies god objects, counts tech debt, and has specific recommendations with file paths.

---

## Anti-Rationalization Table

Your natural tendency is to take shortcuts. Here's why each shortcut produces a bad report:

| What you'll want to do | Why it seems reasonable | Why it produces garbage |
|------------------------|----------------------|----------------------|
| Skip Step 2, just format CLI output | "The CLI already extracted everything" | CLI gives structure, not understanding. You'll produce a module list, not an architecture report. The user can run the CLI themselves. |
| Read only README, skip code exploration | "README describes the architecture" | READMEs are aspirational. Code is truth. README says "clean architecture", code has 5000-line god objects. |
| Write "Catalyst optimizer optimizes queries" | "That's what it does" | This tells the reader nothing. HOW does it optimize? Rule-based? Cost-based? What rules? Which file? |
| Skip mermaid diagrams | "Text descriptions are enough" | Humans process diagrams 60,000x faster than text. A report without diagrams is a wall of text nobody will read. |
| List all 48 modules in a flat table | "Completeness is important" | Nobody reads a 48-row table. Group by layer, show relationships. 6 layers with 48 modules > 48 rows. |
| Write recommendations like "improve test coverage" | "It's true and helpful" | It's true and useless. Which tests? Which modules? What's the current coverage? What specifically should be tested? |
| Skip tech debt section | "I didn't find any issues" | Every codebase >50K LOC has tech debt. If you found zero, you didn't look. Count TODO/FIXME/HACK. Find the largest files. Check for deprecated APIs. |
| Use [EXTRACTED] for everything | "I'm confident in my analysis" | If you didn't read the source code, it's [INFERRED]. Be honest — it builds trust. |

---

## Red Flags — You're Going Wrong If:

- Your report is under 200 lines → you skipped the deep-dive
- You have zero mermaid diagrams → you're writing a text dump
- No file paths appear in your analysis → you're describing from imagination
- Your "recommendations" section says "consider" or "could" → too vague
- All your evidence is [EXTRACTED] → you didn't think, you just reformatted
- The "how it works" section lists modules instead of explaining execution flow
- You spent <5 minutes on the whole report → you skipped something

---

## Evidence Tagging

Tag every significant claim in the report:

| Tag | Meaning | Source | Example |
|-----|---------|--------|---------|
| `[EXTRACTED]` | Machine-verified fact | CLI output, file counts, grep results | "48 Maven modules [EXTRACTED]" |
| `[INFERRED]` | Reasoned from evidence | Code reading, pattern recognition | "Rule-based optimizer [INFERRED from 44 rule files in sql/catalyst/optimizer/]" |
| `[AMBIGUOUS]` | Uncertain, flagged for review | Partial evidence, README claims | "Performance tested via benchmarks [AMBIGUOUS — no benchmark code found in repo]" |

Readers trust transparent analysis more than confident-sounding guesses.

---

## CLI Reference

### `sdp architect analyze <repo-path>`

| Flag | Default | Description |
|------|---------|-------------|
| `--tier` | `2` | Detail: 1 (~2K tokens), 2 (~5-15K), 3 (on-demand) |
| `--format` | `json` | Output: `text`, `json`, `mermaid` |
| `--section` | all | Only: `profile`, `report`, `model`, `diagrams`, `summary` |
| `--no-llm` | `false` | Deterministic-only (no LLM enrichment) |
| `--timeout` | `5m` | Pipeline timeout |
| `--verbose` | `false` | Per-extractor timing |
| `--skip-git` | `false` | Skip git history extractor |
| `--language` | auto | Force: `go`, `java`, `python`, `typescript` |

### `sdp architect c4 <repo-path>`

| Flag | Default | Description |
|------|---------|-------------|
| `--level` | all | Level: `1` (context), `2` (container), `3` (component) |
| `--output` | stdout | Output directory for .mmd files |
| `--format` | `mermaid` | Format: `mermaid`, `json` |

### `sdp architect eval <repo-path>`

| Flag | Description |
|------|-------------|
| `--ground-truth` | Path to ground-truth JSON |

---

## For More

- **Report template with section details**: read `references/report-template.md`
- **Extending the module (adding extractors, adapters)**: read `references/extending.md`
- **Debugging empty diagrams, wrong containers**: read `references/debugging.md`
