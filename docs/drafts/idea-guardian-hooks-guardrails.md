# Idea: Guardian-style Hooks and Guard Rails for SDP

## Problem

SDP currently has split enforcement surfaces: repo hook scripts, `sdp hooks install` template, and CLI checks that are not fully aligned. This creates drift between local behavior, CI behavior, and documentation.

## Goal

Introduce Guardian-style guard-rail patterns into SDP while preserving `sdp guard` as the workstream scope gate.

## Non-Goals

- Replacing `sdp guard` ownership and scope model.
- Large workflow redesign outside hooks/guard rails.
- Full policy governance redesign in one step.

## User Outcomes

- Developers get consistent local checks on staged changes.
- Hook installation is safe, idempotent, and does not clobber user-managed hooks.
- CI and local checks use the same rule evaluation and exit codes.
- Exceptions are explicit, expiring, and auditable.

## Requirements (from interview)

1. **Rollout mode:** hybrid blocking.
   - `ERROR` blocks on pre-commit/pre-push/CI.
   - `WARNING` is advisory.
2. **MVP scope:** full hook surface.
   - `pre-commit`, `pre-push`, `post-merge`, `post-checkout`, and CI integration.
3. **Rules config:** separate file.
   - `.sdp/guard-rules.yml` as canonical rule definition.
4. **Exceptions + governance timing:** MVP now.
   - TTL-based exceptions and governance/meta-check in the first release.
5. **Inspiration source link is mandatory.**
   - Feature and workstream descriptions must include an explicit source link (`Source of inspiration`) to the external/internal reference used (e.g., guardian-cli doc/file URL or repo path).

## Source of inspiration

- [guardian-cli README](https://github.com/AlexGladkov/guardian-cli/blob/main/README.md)
- [guardian-cli hook implementation](https://github.com/AlexGladkov/guardian-cli/blob/main/internal/git/hooks.go)
- [guardian-cli guard check pipeline](https://github.com/AlexGladkov/guardian-cli/blob/main/internal/cli/check.go)
- [guardian-cli rule engine registry](https://github.com/AlexGladkov/guardian-cli/blob/main/internal/engine/checker.go)

## Proposed Functional Capabilities

1. Marker-based hook ownership and idempotent install/uninstall.
2. `sdp guard check --staged` for diff-based rule checks.
3. Stable output contracts:
   - human output for local runs,
   - JSON output for CI integrations.
4. CI diff range detection (`base..head`) with safe defaults.
5. Exception model with expiry (`expires_at`) and path/rule scoping.
6. Governance meta-check for guarded policy files.

## Risks

- Regression risk from consolidating duplicate hook paths.
- False positives if staged diff parsing is too strict early.
- Adoption friction if `ERROR` set is too broad initially.

## Success Metrics

- Hook install success rate and idempotency checks.
- Reduced drift incidents between docs/CI/local checks.
- Reduction in CI rejections caused by local/CI mismatch.
- Exception inventory remains low and actively expires.

## Next Step

Design and decompose into `F063` workstreams with phased dependencies and explicit scope files.
