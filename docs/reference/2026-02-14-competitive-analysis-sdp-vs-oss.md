# SDP Competitive Analysis vs OSS SDD/Agent Frameworks (2026-02-14)

## Scope and Method

Compared SDP against these repositories:

- `/Users/fall_out_bug/projects/vibe_coding/BMAD-METHOD`
- `/Users/fall_out_bug/projects/vibe_coding/OpenSpec`
- `/Users/fall_out_bug/projects/vibe_coding/ccpm`
- `/Users/fall_out_bug/projects/vibe_coding/claude-flow`
- `/Users/fall_out_bug/projects/vibe_coding/claude-task-master`
- `/Users/fall_out_bug/projects/vibe_coding/spec-kit`
- `/Users/fall_out_bug/projects/vibe_coding/claude-brain`

Evidence sources: each repo `README.md` + core docs/command references and SDP roadmap/workstreams.

---

## Executive Decision

**Recommendation: continue with SDP as the base, do not rebase on any single competitor repo.**

Why:
- SDP already has protocol + evidence + traceability direction that most competitors do not provide together.
- Competitors have strong isolated advantages (UX onboarding, action-oriented flow, task ergonomics, memory, collaboration), but each has major gaps for SDP goals (auditability, protocol rigor, multi-layer roadmap, governance depth).
- Best strategy is **selective adoption** into F068-F073.

---

## Repo-by-Repo Assessment

## 1) BMAD-METHOD

### Strong vs SDP
- Clear workflow map and phase framing (`docs/reference/workflow-map.md`).
- Strong guided navigation via `/bmad-help` and command generation (`docs/reference/commands.md`).
- Practical quick-flow with escalation guardrails (`docs/explanation/quick-flow.md`).

### Weak vs SDP
- Less explicit protocol/evidence chain focus.
- Heavier methodology surface (many workflows/agents) can increase learning curve.

### Take into SDP
- Adaptive “what next” assistant behavior (`/bmad-help` style).
- Explicit quick-flow to full-flow escalation rules.
- Better command discovery by intent.

### Roadmap mapping
- `00-068-03` Help/Status IA
- `00-069-02` Deterministic recommendation engine
- `00-069-04` Drive-mode guided loop UX

---

## 2) OpenSpec

### Strong vs SDP
- Action-oriented workflow (`/opsx:new`, `/opsx:continue`, `/opsx:ff`, `/opsx:apply`, `/opsx:verify`, `/opsx:archive`).
- Strong artifact graph with readiness/blocking model (`openspec/specs/artifact-graph/spec.md`).
- CLI support for agent-friendly status/instructions/templates (`docs/cli.md`).
- Brownfield and low-ceremony positioning (`docs/concepts.md`).

### Weak vs SDP
- Lighter on enforcement and quality/evidence rigor than SDP target.
- Less emphasis on policy/gating for high-assurance workflows.

### Take into SDP
- “Actions not phases” UX in CLI surfaces.
- First-class “ready/blocked/next artifact” model in status.
- Explicit instructions endpoint semantics for agent tooling.

### Roadmap mapping
- `00-068-03` Help/Status IA
- `00-069-01` Next-step contract and state model
- `00-069-03` Next-step in status/help/error outputs
- `00-072-02` Import pipeline (OpenSpec adapter)

---

## 3) ccpm

### Strong vs SDP
- GitHub-native collaboration loop, strong task/project visibility (`README.md`, `COMMANDS.md`).
- Clear team handoff and multi-agent parallel execution narrative.
- Good context hygiene mindset and local mode fallback (`LOCAL_MODE.md`, `CONTEXT_ACCURACY.md`).

### Weak vs SDP
- GitHub-centric model can be too opinionated for all SDP deployments.
- Protocol and evidence schema are less explicit than SDP direction.

### Take into SDP
- Team workflow surfaces: standup/blocked/in-progress style summaries.
- Structured handoff payload discipline.
- Scope collision and coordination UX for parallel work.

### Roadmap mapping
- `00-071-01` Team operating model
- `00-071-02` Handoff package contract
- `00-071-03` Scope collision collaboration UX
- `00-071-05` Team adoption playbook
- `00-072-02` Import adapter for ccpm task structures

---

## 4) claude-flow

### Strong vs SDP
- Rich orchestration toolbox (agent swarms, claims/handoff, routing, memory/autopilot concepts).
- Advanced context/memory continuity narrative (autopilot, compaction interception).
- Strong operational thinking around coordination/load balancing in docs.

### Weak vs SDP
- Very high complexity and operational surface area.
- Risk of over-engineering and hard-to-control product scope.
- Claims are broad; adoption should be evidence-first and incremental.

### Take into SDP (narrow slice only)
- Claims/handoff lifecycle concepts for team coordination.
- Lightweight resume/context continuity for failed/interrupted flows.
- Do not import full swarm/intelligence stack into near-term SDP scope.

### Roadmap mapping
- `00-070-03` Resume/checkpoint UX hardening
- `00-071-02` Handoff package and session continuity
- `00-071-03` Collision + coordination UX

---

## 5) claude-task-master

### Strong vs SDP
- Excellent task-level ergonomics: `next`, dependencies, subtasks, tags, metadata (`docs/command-reference.md`, `docs/task-structure.md`).
- Strong MCP and multi-provider operational integration.
- Useful loop pattern for bounded autonomous iterations (`apps/docs/capabilities/loop.mdx`).

### Weak vs SDP
- Heavy API/model configuration burden can hurt first-run UX.
- Task engine depth can overshadow protocol-level clarity if copied wholesale.

### Take into SDP
- Better “next task” prioritization semantics and dependency clarity.
- Metadata model for integration IDs in workstream/task objects.
- Controlled loop mode pattern for repeated recovery/quality actions.

### Roadmap mapping
- `00-069-01` Next-step contract model
- `00-069-02` Recommendation engine logic
- `00-070-02` Recovery playbook engine (loop-assisted retries)
- `00-072-02` Taskmaster import adapter

---

## 6) spec-kit

### Strong vs SDP
- High-quality command grammar for spec lifecycle (`/speckit.constitution`, `specify`, `clarify`, `plan`, `tasks`, `implement`).
- Excellent requirement-quality discipline (“checklist as unit tests for English”) (`templates/commands/checklist.md`).
- Strong ambiguity/clarification gates before implementation (`docs/quickstart.md`, `spec-driven.md`).

### Weak vs SDP
- Can feel heavyweight for fast tactical changes.
- Less explicit runtime evidence/provenance stack than SDP’s direction.

### Take into SDP
- Stronger clarify/checklist requirement gates.
- Better spec quality scoring and ambiguity detection.
- Requirement-quality checks before implementation recommendations.

### Roadmap mapping
- `00-068-01` UX baseline (learning metrics tied to specification clarity)
- `00-069-04` Drive-mode guided loop (clarify flow)
- `00-073-01` Explainability model levels
- `00-073-04` Policy transparency and gate reasoning
- `00-073-05` Trust pack criteria

---

## 7) claude-brain

### Strong vs SDP
- Clear memory portability model (`.claude/mind.mv2`) and simple mental model.
- Practical hook-based capture on session/tool lifecycle (`hooks/hooks.json`, `src/hooks/*.ts`).
- Strong continuity ergonomics (search/recent/ask/stats memory commands).

### Weak vs SDP
- Memory-first product, not full protocol/task orchestration.
- Needs careful governance when mapped to strict audit contexts.

### Take into SDP
- Lightweight session continuity artifacts for handoff and resume.
- Local searchable memory index for operational context.
- Keep evidence log as canonical truth; memory is assistive layer, not audit source.

### Roadmap mapping
- `00-070-03` Resume/checkpoint UX
- `00-071-02` Handoff package continuity
- `00-073-03` Evidence trace UX unification

---

## What SDP Does Better (Current/Planned Moat)

- Protocol-centric roadmap with explicit schema/evidence trace direction ([docs/roadmap/ROADMAP.md](../roadmap/ROADMAP.md)).
- Quality + provenance narrative in one platform (plan/generation/verification/approval chain).
- Multi-layer product direction already defined (`README.md`, `PRODUCT_VISION.md`): protocol, CLI, orchestrator.
- Existing workstream discipline and dependency modeling in backlog system.

---

## What SDP Is Currently Worse At

- First-run UX and discoverability are behind OpenSpec/BMAD.
- Team collaboration ergonomics behind ccpm/claude-flow in day-to-day handoff visibility.
- Task ergonomics (`next`, tags, metadata, loop) behind task-master.
- Session continuity ergonomics behind claude-brain/claude-flow memory UX.

---

## Roadmap Adoption Matrix (F068-F073)

| Source Repo | Capability to Adopt | SDP Target WS | Implementation Mechanism | KPI |
|-------------|---------------------|---------------|--------------------------|-----|
| BMAD-METHOD | Contextual next-step helper | 00-068-03, 00-069-02 | Intent-based help/status recommender | Higher first-run completion |
| OpenSpec | Artifact readiness/blocked model | 00-069-01, 00-069-03 | State graph + deterministic “next action” blocks | Lower “what now?” friction |
| OpenSpec | CLI agent instructions endpoint pattern | 00-069-03 | JSON guidance payload in status/help | Better automation compatibility |
| ccpm | Team handoff discipline | 00-071-02 | Structured handoff package with ownership and risks | Reduced handoff loss |
| ccpm | Team operational views | 00-071-04, 00-071-05 | Status slices: blocked/in-progress/next | Faster team alignment |
| claude-flow | Claims lifecycle concepts | 00-071-03 | Non-blocking coordination signals + ownership transitions | Fewer parallel collisions |
| claude-flow | Resume continuity mindset | 00-070-03 | Explicit checkpoint state and resume options | Lower recovery time |
| task-master | Task dependency + next heuristics | 00-069-01, 00-069-02 | Priority/dependency scoring model | Better recommendation acceptance |
| task-master | Metadata/tag interoperability | 00-072-02 | Import mapping for task metadata/dependencies | Smoother migration |
| spec-kit | Clarify + requirement quality gates | 00-069-04, 00-073-04 | Pre-implementation requirement quality checks | Fewer downstream reworks |
| spec-kit | Explainability through structured quality criteria | 00-073-01, 00-073-05 | Level-based rationale and gate explanations | Higher trust score |
| claude-brain | Portable local memory for continuity | 00-070-03, 00-071-02 | Local memory snapshots linked to handoff/resume | Faster context restore |

---

## Guardrails: What NOT to Copy

- Do not copy full claude-flow orchestration complexity into P0/P1 UX track.
- Do not import task-master API/provider surface area into default SDP first-run path.
- Do not replace SDP evidence schema/log with memory store abstractions.
- Do not lock SDP to GitHub-only collaboration model.

---

## Final Recommendation

1. **Keep SDP core strategy** (protocol + evidence + multi-layer architecture).
2. **Adopt UX and team ergonomics aggressively** from OpenSpec/BMAD/ccpm/task-master.
3. **Adopt memory continuity narrowly** from claude-brain/claude-flow into recovery/handoff flows.
4. **Treat interop as first-class** in F072 so migration from competitor ecosystems is easy.
