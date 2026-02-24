# SDP Coding Workflow Predicate (v1)

**Predicate type:** `https://sdp.dev/attestation/coding-workflow/v1`  
**Format:** in-toto Statement v0.1  
**Schema:** [schema/coding-workflow-predicate.schema.json](../schema/coding-workflow-predicate.schema.json)

## Overview

This predicate attests that an AI coding agent (or human) followed the SDP protocol: planned workstreams, stayed within declared scope, passed verification (tests, lint), and completed review. It answers: *"Did the agent actually do what it claimed, or did it wing it?"*

## When to Use

- **Light mode:** CI auto-generates attestations from observation (git diff, test results, lint). No agent action required.
- **Full mode:** `sdp-orchestrate` emits attestations at each phase transition. Agent + CI together produce the chain.

## Statement Structure

```json
{
  "_type": "https://in-toto.io/Statement/v0.1",
  "predicateType": "https://sdp.dev/attestation/coding-workflow/v1",
  "subject": [{ "name": "PR URL or branch", "digest": { "sha256": "commit SHA" } }],
  "predicate": {
    "intent":       { "issue_id", "trigger", ... },
    "plan":         { "workstreams", "ordering_rationale" },
    "execution":    { "branch", "changed_files", "claimed_issue_ids" },
    "verification": { "tests", "lint", "coverage" },
    "review":       { "self_review", "adversarial_review" },
    "risk_notes":   { "residual_risks", "out_of_scope" },
    "boundary":     { "declared", "observed", "compliance" },
    "provenance":   { "run_id", "orchestrator", "captured_at", ... },
    "trace":        { "beads_ids", "commits", "pr_url" }
  }
}
```

## Validation

Use `sdp-evidence` from [sdp_lab](https://github.com/fall-out-bug/sdp_lab):

```bash
sdp-evidence validate .sdp/evidence/F028.json
sdp-evidence inspect .sdp/evidence/F028.json
```

## Signing

Attestations can be signed with Sigstore (keyless) for tamper-evidence:

```bash
cosign sign-blob --yes --bundle attestation.bundle attestation.json
```

## See Also

- [Getting Started (sdp_lab)](https://github.com/fall-out-bug/sdp_lab/blob/master/docs/getting-started.md)
- [in-toto attestation format](https://github.com/in-toto/attestation)
- [ADR-002: Standards Pivot](https://github.com/fall-out-bug/sdp_lab/blob/master/docs/decisions/ADR-002-standards-pivot.md)
