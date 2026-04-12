# StratAudit Runtime Policy

StratAudit is runtime-neutral. The skill must choose a runtime by capability and
trust needs, not by vendor preference.

## Runtime Order

1. harness-injected host-native runtime
2. configured OpenAI-compatible runtime
3. OpenRouter as the default network enhancer/fallback
4. no runtime only for artifact-only modes

## Capability Requirements By Mode

| Mode | Required capability |
|------|---------------------|
| `corpus-audit` | text extraction and enough reasoning to summarize corpus quality; can reuse artifacts if they already exist |
| `traceability-audit` | structured extraction, embeddings, and conservative verification |
| `coverage-audit` | inspectable coverage artifacts or enough runtime support to produce them |
| `evidence-pack` | no runtime if artifacts already exist; otherwise same floor as the audit that produced them |
| `report-redraft` | no runtime if rewriting from existing artifacts only; must not introduce new claims |

## Selection Rules

- prefer host-native models when they meet the capability bar
- prefer artifact reuse over recomputation when the user is redrafting or packaging evidence
- do not silently fall back from a stronger mode to a weaker one
- if the runtime cannot support the requested mode, fail closed with an explicit explanation

## CLI Boundary

`sdp-strataudit` can resolve configured network runtimes. It cannot create a
host-native runtime by itself.

That means:

- harnesses such as Cursor, Codex, Claude, or OpenCode can inject native runtimes when available
- OpenRouter remains useful as a capability amplifier, not as the only valid execution path

## Must Be Reported

Every run should state:

- selected mode
- runtime class used
- whether artifacts were reused or regenerated
- any trust caveat caused by runtime limits
