# ADR-0001: Cursor-Native Consensus Protocol v2.0

## Status

Proposed

## Context

The current Consensus Protocol (v1.2) relies on manual orchestration and "gentlemen's agreements" between agents.
- **Fragility**: Agents rely on textual descriptions in `PROTOCOL.md` or loose JSON schemas, leading to structure errors.
- **State Dispersion**: The state is "implicit in artifacts", making it hard to determine the current phase or blockers without reading multiple files.
- **Context Overload**: Large epics flood the context window with raw messages, increasing cost and latency.
- **Process Rigidity**: Simple bug fixes require the same full chain (Analyst -> Architect -> Tech Lead -> ...) as complex features.
- **Execution Gap**: Tech Lead plans (`implementation.md`) are often too high-level for Developers to execute reliably in one go.

We evaluated external frameworks like **LangGraph** but rejected them because they treat the agent process as a "black box" external to the IDE, breaking the "Human-in-the-loop" workflow native to Cursor/Claude Code.

## Decision

We will upgrade the protocol to a **"Cursor-Native" Operating System** that leverages the IDE's ability to read code and files.

### 1. TypeScript Interfaces as Schema ("Code is Law")
Instead of textual descriptions, we will define strict TypeScript interfaces in `consensus/schema.d.ts`.
- **Why**: LLMs are better at generating valid code matching interfaces than following abstract rules.
- **Mechanism**: Agents will have `schema.d.ts` in their context and must produce JSONs that strictly adhere to `AgentMessage` or `EpicState` types.

### 2. Centralized State File (`status.json`)
We introduce a single mutable file `consensus/status.json` that acts as the "traffic light" for the system.
- **Structure**: Defined by `EpicState` interface (phase, iteration, blockers, approvals, vetoes, active_mode).
- **Rule**: Agents must read `status.json` before acting and update it after finishing.
- **Benefit**: Immediate visibility of system state (e.g., "Blocked by Architect Veto") without parsing inbox histories.

### 3. Context Anchors (Agent Consoles)
We replace raw prompts with Markdown "Console" files in `consensus/agents/{ROLE}.md`.
- **Mechanism**: Files act as control panels with embedded context references (`<!-- @context: ... -->`) and clickable actions.
- **Usage**: The user opens `ARCHITECT.md` and uses Cursor Chat (Cmd+K) to "Execute Action 1".
- **Benefit**: Zero-setup context loading; the file itself drives the agent's behavior.

### 4. Kanban-as-Code (Micro-Tasking)
We introduce `consensus/kanban.json` to bridge the gap between High-Level Plan and Code.
- **Tech Lead Role**: Now produces a list of atomic `Task` objects (id, title, desc, dependencies) in `kanban.json` instead of just prose.
- **Developer Loop**: Selects a `todo` task -> Moves to `in_progress` -> Implements -> Verifies -> Moves to `done`.
- **Benefit**: Drastically reduces context usage (agent only sees one task at a time) and enables "Resumability".

### 5. Dynamic Topology
We introduce execution modes in `status.json`:
- **Full**: Analyst -> Architect -> Tech Lead -> Dev -> QA -> DevOps (Default)
- **Fast Track**: Dev -> QA (For minor bugs)
- **Hotfix**: Dev -> Prod (For critical incidents)

## Consequences

### Positive
- **Reliability**: Strong typing reduces invalid JSON outputs.
- **Observability**: `status.json` provides a real-time dashboard of the epic.
- **Efficiency**: "Fast Track" mode saves tokens and time for simple tasks.
- **Granularity**: Kanban allows precise control over implementation and easier rollback.
- **DX**: "Context Anchors" make running agents as simple as opening a file.

### Negative
- **Setup Overhead**: Requires creating and maintaining `schema.d.ts` and `status.json`.
- **Strictness**: Agents might fail if they cannot parse the strict schema (mitigated by TypeScript's clarity).
- **Manual State Management**: Users/Agents must ensure `status.json` is not corrupted (mitigated by file validation tools).

## References
- Replaces: Consensus Protocol v1.2 (JSON Schema descriptions)
- Rejects: External orchestration frameworks (LangGraph, CrewAI) in favor of file-based stigmergy.


