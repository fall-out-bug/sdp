# Installation Profiles (Layered Public OSS)

> Status: transitional. This document defines the target installation model for progressive SDP adoption in the public MIT repository.

## Profile Matrix

| Profile | Level | Includes | Distribution | License |
|---------|-------|----------|--------------|---------|
| `protocol` | `L0` | Prompts, guides, schemas, protocol templates | Claude plugin/prompt distribution | MIT |
| `safety` | `L1` | Hooks, guard, traces, provenance layer | Homebrew package | MIT |
| `core` | `L2` | Orchestrator and core SDP CLI tools | Homebrew package | MIT |

## Upgrade Path

- Start with `protocol` for low-friction adoption.
- Add `safety` when teams need stronger control and evidence.
- Add `core` when teams need orchestration and scale.

## Boundary Rules

- `L0-L2` are MIT and independently adoptable.
- Public profiles must not require non-public dependencies.
- Profile upgrades must be deterministic and reversible.

## Packaging Targets

- `L0`: plugin/prompt package for protocol-only operation.
- `L1-L2`: brew-installable profile packages.

## Related Docs

- `/Users/fall_out_bug/projects/vibe_coding/sdp/PRODUCT_VISION.md`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/docs/vision/ROADMAP.md`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/docs/vision/LAYERED-ADOPTION.md`
