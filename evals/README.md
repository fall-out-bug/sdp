# SDP Prompt Evals

Evaluation framework for SDP protocol prompts (skills). Test cases define input → expected output structure for key skills.

## Structure

- `review/` — @review skill evals
  - `test_cases.jsonl` — one JSON object per line: `{"input": "...", "expected_keys": ["feature", "verdict", "reviewers", ...]}`
  - `rubric.md` — scoring criteria
  - `results/` — score outputs per run (e.g. `v14.0.0.json`)
- `build/` — @build skill evals
  - `test_cases.jsonl` — input ws-id, expected_keys for ws-verdict (ws_id, feature_id, verdict, quality_gates, timestamp, existing_work_summary)
- `idea/` — @idea skill evals
  - `test_cases.jsonl` — input problem description, expected_keys for intent (problem, users, success_criteria)
- `rubric.md` (per skill) — scoring criteria (optional)

## How to run

1. **Review evals:** For each test case in `review/test_cases.jsonl`, invoke @review with the given input and assert the output JSON has the expected structure (e.g. all 7 reviewer keys).
2. **Schema validation:** Validate output against `schema/review-verdict.schema.json`.
3. **Record results:** Write scores to `review/results/<version>.json`.

## Format: test_cases.jsonl

Each line is a JSON object:

```json
{"id": "review-001", "input": "feature-example", "expected_keys": ["feature", "verdict", "reviewers", "round", "timestamp"], "reviewer_keys": ["qa", "security", "devops", "sre", "techlead", "docs", "promptops"]}
```

No trailing newline after the last line.
