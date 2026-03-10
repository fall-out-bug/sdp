# Idea: Audit Demo Surface Polish

Status: draft
Type: polish
Scope: existing SDP CLI and evidence surfaces

## Why This Exists

SDP already has the right primitives for an audit walkthrough: evidence logs, `sdp log show`, `sdp status`, `sdp next`, install flow, and a documented `sdp demo` surface. What it lacks is a tight first-success path that makes the audit story obvious to a new user in one sitting.

The comparison with `dcl-eval-pipeline-demo` is useful here. That repo is not stronger than SDP as a system, but it is easier to understand as a public-facing audit demo. We should borrow that packaging quality without copying its thinner architecture.

## Problem

Today SDP is stronger than lightweight eval demos in protocol depth, evidence quality, and workflow discipline, but weaker in demo clarity.

Visible gaps:

- `sdp` docs advertise `sdp demo` as a guided first-success walkthrough.
- Shell completion also advertises `demo`.
- The repo does not appear to contain a corresponding `cmd/sdp` implementation yet.
- The audit story is spread across `sdp init`, `sdp status`, `sdp next`, `sdp log show`, docs, and evidence schemas rather than presented as one coherent walkthrough.

This creates an adoption gap, not a protocol gap.

## Goal

Polish the existing audit and evidence surfaces so a new user can understand, run, and inspect an SDP workflow as an auditable system without learning the whole protocol first.

## Non-Goals

- No new protocol layer
- No new evidence schema version
- No new review or orchestration subsystem
- No dashboard product
- No Langfuse/Phoenix exporter requirement for initial scope

## Existing Building Blocks

Already present in SDP:

- `sdp init`
- `sdp status`
- `sdp next`
- `sdp log show`
- evidence envelope and validation flow
- docs for quickstart, protocol, manifesto, and CLI reference

This idea is about packaging and alignment, not inventing a new capability.

## Proposed Scope

### 1. Make the first-success walkthrough honest

Pick one of these paths and make docs plus implementation match:

- implement a minimal `sdp demo` command as a thin guided walkthrough over existing commands, or
- remove `sdp demo` from public surfaces until it exists and document the walkthrough explicitly with existing commands

The first option is better if the implementation stays thin and does not create a second workflow.

### 2. Package the audit story around existing evidence

Add one canonical walkthrough that shows:

1. initialize SDP
2. run a small flow
3. inspect project state with `sdp status`
4. inspect next action with `sdp next`
5. inspect evidence with `sdp log show`

This should end in a concrete answer to: "what proof did SDP produce?"

### 3. Add a sample evidence story, not a fake platform layer

Ship one small sample or fixture-backed walkthrough that demonstrates:

- a realistic evidence log entry
- the shape of a verification trail
- the difference between workflow state and evidence output

This can be docs, test fixture reuse, or a sample report path. It should not become a separate service.

### 4. Keep observability optional and layered

Document external export and tracing as integration points, not as core SDP requirements. The right shape is:

- SDP remains the source of truth for local evidence
- external observability tools are adapters on top

## Acceptance Criteria

- A newcomer can follow one documented path from install to evidence inspection in under 10 minutes.
- Public docs and CLI surfaces no longer advertise a walkthrough command that does not exist.
- If `sdp demo` ships, it is a thin wrapper over existing flows rather than a parallel subsystem.
- At least one canonical example shows how `sdp status`, `sdp next`, and `sdp log show` fit together.
- The resulting materials strengthen the audit story without changing protocol semantics.

## Why This Fits Polish Scope

This is polish because it improves packaging, onboarding, and surface honesty around capabilities SDP already claims or already has in pieces.

It does not require:

- a new protocol primitive
- a new evidence model
- a new platform service
- a new roadmap feature line

It is best treated as adoption polish for the existing evidence layer.

## Recommended First Cut

Smallest useful slice:

1. decide whether `sdp demo` is implemented now or removed from references
2. add one walkthrough doc for audit-first onboarding
3. add one sample evidence inspection example around `sdp log show`

If that lands well, optional follow-up polish can add richer examples or external tracing adapters.
