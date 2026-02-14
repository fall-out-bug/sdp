# SDP Layered Adoption Model (L0-L2, Public OSS)

> Last updated: 2026-02-14

## Purpose

Define a progressive adoption path so teams can use SDP in slices without adopting the full stack.

This document covers **public OSS scope only**.

## Level Matrix

| Level | Scope | Install Surface | Distribution | License |
|-------|-------|-----------------|--------------|---------|
| `L0` | Protocol only: prompts, guides, templates, schemas | Claude plugin / prompt pack | Plugin-style distribution | MIT |
| `L1` | Safety and evidence: hooks, guard, traces, provenance | CLI safety bundle | Homebrew package | MIT |
| `L2` | Orchestrator core: plan/apply/log, dispatcher, checkpoints | CLI core bundle | Homebrew package | MIT |

## Design Rules

- Each level must be useful as a standalone product.
- Upgrade path is additive: `L0 -> L1 -> L2`.
- Public layers must not require non-public dependencies.
- Public distributions remain MIT.

## Packaging Targets

## L0 Package Target

- Deliver protocol starter for Claude plugin workflows.
- Include onboarding quickstart for "protocol-only" teams.
- Include migration guide to enable `L1/L2` later.

## L1-L2 Package Target

- Publish brew-installable profiles:
  - `protocol`
  - `safety`
  - `core`
- Ensure deterministic upgrade and rollback between profiles.
- Keep profile outputs consistent with protocol contracts.

## Implementation Mapping

- Public strategy: `/Users/fall_out_bug/projects/vibe_coding/sdp/PRODUCT_VISION.md`
- Public execution roadmap: `/Users/fall_out_bug/projects/vibe_coding/sdp/docs/vision/ROADMAP.md`
- Public adoption workstreams: `/Users/fall_out_bug/projects/vibe_coding/sdp/docs/workstreams/backlog/F074-WORKSTREAMS-SUMMARY.md`
