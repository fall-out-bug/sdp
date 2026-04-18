# StratAudit Output Modes

The portable skill should not behave like one vague "run the audit" instruction.
Different user intents need different outputs and different trust boundaries.

## Modes

| Mode | Primary question | Inputs | Must emit | Must not claim |
|------|------------------|--------|-----------|----------------|
| `corpus-audit` | "What corpus do we actually have and is it usable?" | corpus root or source dirs | corpus inventory, exclusions, per-level coverage, trust caveats | final strategic alignment verdict |
| `traceability-audit` | "What links to what and where are the gaps?" | corpus or reusable artifacts | entities, traces, findings, trust caveats | verified traces without inspectable support |
| `coverage-audit` | "What is covered, by level and by document?" | corpus or reusable artifacts | coverage table with explicit denominators | percentages without denominator or caveat |
| `evidence-pack` | "Show me the proof behind the claims." | existing `.strataudit` artifacts preferred | quote/span references, trace rows, grouped findings, caveats | executive confidence beyond the evidence pack |
| `report-redraft` | "Make the report better without changing truth." | existing evidence pack or prior artifacts | rewritten sections and explicit unchanged trust boundaries | new entities, new traces, or upgraded certainty |

## Default Routing

- unknown corpus quality -> `corpus-audit`
- direct alignment or trace question -> `traceability-audit`
- explicit coverage question -> `coverage-audit`
- request for proof or trace drill-down -> `evidence-pack`
- request to rewrite an existing audit report -> `report-redraft`

## Output Discipline

Every mode should expose:

- artifact paths or artifact references
- selected runtime or artifact-only path
- key caveats
- what is not claimed
