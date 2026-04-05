# Reality OSS Reference

Public OSS `reality` is a local, evidence-first baseline for one repository. It emits open artifacts and validates them against the published schema contract.

## Commands

```bash
sdp reality emit-oss [--quick|--deep|--bootstrap-sdp] [--focus=architecture|quality|testing|docs|security]
sdp reality validate
```

## Outputs

Machine-readable outputs:

- `.sdp/reality/reality-summary.json`
- `.sdp/reality/feature-inventory.json`
- `.sdp/reality/architecture-map.json`
- `.sdp/reality/integration-map.json`
- `.sdp/reality/quality-report.json`
- `.sdp/reality/drift-report.json`
- `.sdp/reality/readiness-report.json`

Human-readable outputs:

- `docs/reality/summary.md`
- `docs/reality/architecture.md`
- `docs/reality/quality.md`
- `docs/reality/bootstrap.md`

## Published Contract

- [../specs/reality/OSS-SPEC.md](../specs/reality/OSS-SPEC.md)
- [../specs/reality/ARTIFACT-CONTRACT.md](../specs/reality/ARTIFACT-CONTRACT.md)
- `schema/reality/*.schema.json`

## Publish Checklist

1. Run `sdp reality emit-oss` in a representative repository.
2. Run `sdp reality validate`.
3. Check that `docs/reality/` and `.sdp/reality/` changed deterministically on rerun.
4. Confirm docs and help text describe OSS `reality` as a single-repo baseline, not a consulting-grade multi-agent mesh.
5. If the public contract changes, update `schema/reality/`, `docs/specs/reality/`, and `prompts/skills/reality/SKILL.md` in the same PR.
