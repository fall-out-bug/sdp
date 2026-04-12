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

### Step 1: Structural Extraction (2-5 min)

Try `sdp architect` CLI first. If it's unavailable, use the manual fallback — both produce equivalent data.

**Option A — CLI (if `sdp` binary is available):**
```bash
sdp architect analyze --no-llm --tier 2 --section summary /path/to/repo
sdp architect c4 --level 1 /path/to/repo
sdp architect c4 --level 2 /path/to/repo
```
Flags go BEFORE the repo path.

**Option B — Manual fallback (works everywhere):**
Run these commands via Bash to collect the same structural data:
```bash
# File count and LOC
find /path/to/repo -type f \( -name "*.go" -o -name "*.java" -o -name "*.scala" -o -name "*.py" -o -name "*.ts" -o -name "*.js" -o -name "*.rs" -o -name "*.rb" -o -name "*.c" -o -name "*.cpp" -o -name "*.h" \) | head -20000 | wc -l
find /path/to/repo -type f \( -name "*.go" -o -name "*.java" -o -name "*.scala" -o -name "*.py" -o -name "*.ts" \) | head -20000 | xargs wc -l 2>/dev/null | tail -1

# Language breakdown (top extensions by file count)
find /path/to/repo -type f -name "*.*" | grep -v node_modules | grep -v vendor | grep -v .git | sed 's/.*\.//' | sort | uniq -c | sort -rn | head -15

# Build system detection
ls /path/to/repo/{pom.xml,build.gradle,build.sbt,go.mod,package.json,Cargo.toml,Makefile,CMakeLists.txt,pyproject.toml,setup.py} 2>/dev/null

# Module boundaries (Maven/Gradle subprojects)
find /path/to/repo -name "pom.xml" -not -path "*/target/*" | head -60
find /path/to/repo -name "build.gradle" -o -name "build.gradle.kts" | head -30
find /path/to/repo -name "build.sbt" | head -5

# Infrastructure files
find /path/to/repo -name "Dockerfile" -o -name "docker-compose*.yml" -o -name "*.tf" | head -20
find /path/to/repo -path "*/.github/workflows/*.yml" | head -20
find /path/to/repo -name "*.proto" | head -20

# API specs
find /path/to/repo -name "openapi*.yaml" -o -name "openapi*.json" -o -name "swagger*.yaml" | head -10
```

**Gate**: You have file count, LOC, languages, build system, and module list. If you're missing any, the deep-dive agents in Step 2 will fill gaps.

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

**Agent B — Execution Architecture & Structure:**
```
In /path/to/repo:
1. Find main entry point(s): main(), Application, App, Server, cli, cmd/ directories
2. Trace the startup path: what gets initialized, in what order
3. Identify core abstractions (3-5 key interfaces/classes that define the domain model)
4. What's the execution model: request/response, pipeline, event-driven, batch, actor?
5. Find the 5 largest source files (not generated/vendored) — these are likely god objects.
   Use: find . -name "*.{go,java,scala,py,ts}" -not -path "*/vendor/*" -not -path "*/node_modules/*" -exec wc -l {} + | sort -rn | head -10
6. If there are submodules/subprojects, identify inter-module dependencies.
Report: entry points, core abstractions with file paths, execution model, top 5 largest files with LOC.
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
