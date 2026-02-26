# @review Eval Rubric

Scoring criteria for review skill output validation.

## Required structure

- **verdict**: One of PASS, FAIL, CHANGES_REQUESTED. Must be present.
- **reviewers**: Object with exactly 7 keys: qa, security, devops, sre, techlead, docs, promptops. Missing role = fail.
- **feature**: Feature ID. Must match input.
- **round**, **timestamp**: Present for traceability.

## Severity consistency

- P0 = exploitable / security-critical. Verdict must be FAIL if any P0.
- P1 = edge case / correctness. Verdict must be FAIL if any P1.
- P2 = debt / maintainability. PASS allowed if max severity is P2.
- P3 = style / nit. PASS allowed.

## Synthesis (if present)

- **conflicts**: Two reviewers disagree on same finding — expect escalation or both positions noted.
- **rubber_stamps**: Reviewer with 0 findings — flag for human review.

## Score (0–100)

- 50: Output is valid JSON and has required top-level keys.
- +20: All 7 reviewer keys present.
- +15: Verdict matches finding severities (no P0/P1 with PASS).
- +15: Schema validates against review-verdict.schema.json.
