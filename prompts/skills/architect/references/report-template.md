# Architecture Report Template

Use this template for every architecture report. All 10 sections are mandatory.
Write in the language the user communicates in.

---

## Section Guide

### Header Block

```markdown
# {Project Name} — Architecture Report

> **Date**: YYYY-MM-DD | **Version**: X.Y | **Tool**: sdp architect
> **Audience**: tech lead, architect, new developer
```

### Section 1: Why Read This Report

2-3 sentences: what repo, how big, what you'll learn. Not a summary — a pitch.

**GOOD**: "You inherited 2.9M lines of distributed compute engine with 15 years of history and 48 modules. This report maps the terrain so you don't spend 3 weeks figuring out what I figured out in 30 minutes."

**BAD**: "This report contains the results of an architecture analysis performed by the sdp architect tool."

### Section 2: What Is This — One Paragraph

Plain language explanation for a smart person who's never seen this codebase. What it does, who uses it, how it's deployed.

**GOOD**: "Spark is a distributed data processing engine. It takes data from any source (HDFS, S3, Kafka), processes it across a cluster of machines, and returns results. It's a framework, not a service — you embed it in your application or submit jobs to a cluster."

**BAD**: "Apache Spark is a multi-module Maven project written primarily in Scala with Java and Python components."

### Section 3: Landscape — Who Talks to Whom (L1)

Mermaid diagram showing: users → the system → external systems. Include data flow labels on arrows.

After the diagram: 2-3 sentences explaining the actors and external systems. Not just "uses Kafka" but WHY — "Kafka provides streaming data ingestion for real-time analytics pipelines."

### Section 4: Module Map — What's Inside (L2)

Mermaid diagram grouping modules by architectural layer (not a flat list).

Table: layer → modules → approximate LOC → purpose. Group by FUNCTION, not by directory.

**GOOD**: 6 rows (Core Engine, SQL Engine, ML, Streaming, Connectors, Infra) each containing multiple modules.

**BAD**: 48 rows, one per Maven module, alphabetically sorted.

### Section 5: How It Works — Execution Architecture

**This is the most important section.** It's what separates a real architect's report from a CLI dump.

Show the main execution path as an ASCII or mermaid flow:
- Input → processing stages → output
- Name the key classes/files at each stage
- Explain the execution model (pipeline? event loop? scheduler? actor model?)

Then cover 2-4 key subsystems in depth:
- Not "Optimizer optimizes" but HOW it optimizes: what algorithm, what tradeoffs, what files
- Include file paths so the reader can go look

If there's a new/strategic component (e.g. a new API, a rewrite), give it its own subsection.

### Section 6: Architectural Patterns

Table: pattern → where used → why it's good/interesting.

Then 3-5 inferred ADRs (Architecture Decision Records):
- What decision was made
- Why (or your best guess)
- What it costs (tradeoffs)

**GOOD**: "Decision: Lazy evaluation for RDD. Why: enables whole-graph optimization before execution. Cost: debugging is harder because nothing happens until an action is called."

**BAD**: "Uses lazy evaluation pattern."

### Section 7: Where It Hurts — Tech Debt & Risks

Three subsections:

**Critical risks** — table with severity, location (file path!), impact. These are things that could cause outages or make changes dangerous.

**Architectural debt** — structural problems: god objects (name them, give LOC), deprecated-but-not-removed code, duplicated abstractions, missing boundaries.

**Tech debt metrics** — TODO/FIXME/HACK counts (actual numbers from grep), biggest files, recurring themes.

**GOOD**: "DAGScheduler.scala — 3700 LOC, manages stage lifecycle, shuffle tracking, retry, speculative execution. Any change is high-risk. [EXTRACTED: wc -l]"

**BAD**: "There is some technical debt in the codebase."

### Section 8: Dependency Graph

Mermaid diagram showing module dependencies. Highlight hubs — modules that many others depend on.

Explain: which modules are the most dangerous to change (highest fan-in). Which modules are isolated and safe to modify.

### Section 9: Testing & API Surface

Table: test category → count → framework → what it covers.

Call out: what's well-tested, what's NOT tested, what the testing strategy is (unit-heavy? integration-heavy? golden tests?).

API surface: what public interfaces exist (REST, gRPC, SDK, CLI). What's stable, what's evolving.

### Section 10: Recommendations

Three subsections for three audiences:

**For tech lead**: table with priority (red/yellow/green), recommendation, reasoning. Be specific: file paths, module names, concrete actions.

**For new developer**: ordered reading list with file paths. "Read X first, then Y, then Z." Not modules — specific files.

**For business**: 3-4 bullets on maturity, risks, strategic direction, alternatives.

---

## Few-Shot Examples

### GOOD Execution Flow (Section 5):

```markdown
### How a SQL query becomes results

SQL string → Parser (ANTLR) → Unresolved Logical Plan
→ Analyzer (115 files, resolves names/types) → Resolved Plan
→ Catalyst Optimizer (44 rule files, batch iteration to fixpoint)
→ Physical Planner (chooses join strategy: BroadcastHash vs SortMerge)
→ DAGScheduler (splits at shuffle boundaries into Stages)
→ TaskScheduler (distributes Tasks to Executors)
→ Executors (TaskRunner + BlockManager, 3-tier storage)
→ Results back to Driver

Key files:
- Parser: sql/catalyst/.../parser/SqlBaseParser.g4
- Optimizer: sql/catalyst/.../optimizer/Optimizer.scala (2947 LOC)
- DAGScheduler: core/.../scheduler/DAGScheduler.scala (3700 LOC)
```

### BAD Execution Flow:

```markdown
### Architecture

Spark has a SQL engine that processes queries. It uses the Catalyst optimizer
and DAG scheduler. Results are computed on executors.
```

(No flow, no file paths, no details — reader learns nothing they couldn't get from Wikipedia.)

### GOOD Tech Debt (Section 7):

```markdown
| # | Problem | Severity | Location | Impact |
|---|---------|----------|----------|--------|
| 1 | God Object: DAGScheduler | 🔴 | core/.../DAGScheduler.scala (3700 LOC) | Stage lifecycle, shuffle tracking, retry, speculative execution all in one file. Any change risks regression |
| 2 | Deprecated DStreams not removed | 🟡 | streaming/ (~150K LOC) | Dead weight — Structured Streaming replaced it in 2016, but backward compat keeps it alive |

Tech debt markers: 142 files with TODO/FIXME/HACK [EXTRACTED: grep -r count]
Top themes: serialization compat (25%), memory management (20%), thread safety (15%)
```

### BAD Tech Debt:

```markdown
The codebase has some technical debt that should be addressed.
Consider improving code quality and test coverage.
```

(No specifics, no numbers, no file paths — completely useless.)
