# StratAudit Evidence Policy

StratAudit must distinguish verified evidence from derived and inferred claims.

## Authority Order

From strongest to weakest:

1. exact document text, quote, or span
2. document metadata and local context such as section or level
3. derived but inspectable artifacts such as coverage tables or saved traces
4. model-assisted inference over inspected evidence
5. analyst prose and executive summary

Lower layers may summarize higher layers. Lower layers may not overrule them.

## Claim Classes

| Class | Meaning | Allowed wording |
|------|---------|-----------------|
| `verified` | backed by inspectable source text or span plus context | "verified", "shown in", "documented in" |
| `derived` | computed from inspectable artifacts with a clear denominator | "derived from artifacts", "coverage shows" |
| `inferred` | model or analyst synthesis built on evidence but not directly quoted | "inferred", "likely", "suggests" |
| `unsupported` | evidence is missing, ambiguous, or below trust bar | do not make the claim |

## Hard Rules

- never fabricate quotes, spans, or source locations
- never call similarity alone a verified trace
- never translate entity names in the evidence layer unless explicitly requested
- never publish a coverage percentage without saying what the denominator is
- never let the executive summary sound stronger than the strongest underlying claim class

## Minimum Evidence Bundle

Every substantive claim should expose or be traceable to:

- source document or artifact
- quote/span, trace row, or coverage table
- claim class: `verified`, `derived`, or `inferred`
- caveat when provenance is partial

## Downgrade and Refusal Rules

- if source text is missing, downgrade `verified` to `inferred` or refuse
- if the denominator is unclear, refuse to claim numeric coverage
- if the runtime produced plausible prose but no inspectable evidence, mark it `unsupported`
- if the user requests broader certainty than the evidence allows, say so directly
