# SDP Evidence Layer — Compliance Reference

Enterprise reference for the Spec-Driven Protocol (SDP) evidence layer. Use this document for security reviews, RFPs, and contract discussions.

**Evidence schema:** [schema/evidence.schema.json](../../schema/evidence.schema.json) (F054).

---

## Data Residency

- **Where evidence lives:** The evidence log lives in the customer's git repository, under `.sdp/log/` (or configured path). No evidence is stored on SDP servers by default.
- **What leaves the repo:** Nothing leaves the repository unless the customer explicitly configures export (e.g. audit log forwarding). SDP CLI and plugins operate locally against the repo.
- **Telemetry:** Anonymous usage statistics may be sent to SDP infrastructure when telemetry is enabled; opt-out is available. Telemetry does not include evidence content or prompts.
- **Code:** Raw code is never transmitted by SDP to third parties; model providers handle prompts and completions under their own terms.

**Use cases:** Banks and regulated industries can keep all evidence on-premises or in their chosen region by keeping the repo in that environment.

---

## Retention & GDPR

- **Who controls retention:** Evidence log is part of git history. Retention is controlled by the customer via branch policy, history pruning, and backup schedules.
- **Right to erasure:** Evidence can be removed using standard git tools (e.g. `git filter-branch`, BFG). Pruning the repo removes the corresponding evidence from history.
- **PII in evidence:** Evidence events do not contain PII by design. Stored fields include event type, workstream ID, timestamps, commit SHA, hashes (e.g. prompt hash), and file paths — not user names, emails, or prompt text. If your workflow records sensitive data in custom fields, define retention and redaction in `.sdp/config.yml` or equivalent.
- **Recommendation:** Document retention policy (e.g. 90 days, 1 year) and implement via automation; reference in `.sdp/config.yml` where supported.

---

## Access Control

- **Permissions:** Access to the evidence log follows git repository permissions. Whoever can read the repo can read the evidence; whoever can write can append or (if they modify history) alter evidence.
- **No separate RBAC in P0:** SDP does not implement its own role-based access control in the initial release. Enterprise RBAC (e.g. separate roles for auditors vs. developers) is on the P2 roadmap.
- **CLI:** Commands such as `sdp log trace` require read access to the repo; they do not escalate privileges.
- **Approval events:** "Who" in approval events is taken from git config (e.g. `user.name`) in P0, not from a separate authenticated identity store.

---

## Integrity Guarantees

SDP provides **corruption detection** for the evidence log, not stronger guarantees.

- **Hash chain:** Events are linked by a hash chain (e.g. `prev_hash`). This detects accidental corruption, partial writes, and lost or reordered events when validating the log.
- **What it does not provide:** The chain does **not** provide tamper-proof or non-repudiation guarantees. Anyone with write access to the repository can modify or delete evidence. There are no cryptographic signatures on evidence records in P0.
- **What it catches:** Accidental corruption, sync errors, incomplete writes.
- **What it does not catch:** Deliberate modification or deletion by a repo administrator.
- **Roadmap:** Signed evidence records for compliance-grade non-repudiation are on the P3 roadmap.

**Use cases:** Airlines and marketplaces can use the log as an honest audit trail of what happened, while security teams should treat it as integrity-checked but not cryptographically attested.

---

## Prompt Privacy

- **Raw prompts:** Raw prompts are **never** stored in the evidence log.
- **What is stored:** A hash of the prompt (`prompt_hash`, e.g. SHA-256) is stored so that runs can be matched to prompts without disclosing content. Model parameters (e.g. temperature, max_tokens) may be stored; these are not considered PII.
- **Customer code:** Evidence may reference file paths and content hashes for changed files; it does not store full file contents in the log. Verification output (e.g. test stdout) may be stored; customers should control sensitivity via test output and configuration.
- **Recommendation:** If test output or metadata could contain sensitive data, configure exclusions or redaction where supported.

---

## Regulatory Alignment

| Framework | SDP evidence relevance |
|----------|-------------------------|
| **SOC2** | Evidence log serves as an audit trail for AI-generated code changes and approvals, supporting change management and access review. |
| **HIPAA** | Evidence does not contain PHI by default. If code or tool output could include PHI, the customer must configure exclusions or handle PHI outside evidence. |
| **DORA** | Evidence supports ICT risk management and documentation of change and testing, supporting resilience and incident response. |
| **EU AI Act** | Evidence log supports transparency and record-keeping for AI-generated artifacts (e.g. code) as required by the regulation. |

Mapping is indicative; final compliance is the customer's responsibility with their legal and security teams.

---

## What SDP Does NOT Guarantee

- **Tamper-proof or immutability:** Repository admins can alter or delete evidence. The hash chain detects corruption; it does not prevent or cryptographically bind authorized users.
- **Non-repudiation:** P0 evidence is not signed; we do not guarantee that a given event was produced by a specific identity or that it was not modified after creation.
- **Access control beyond git:** There is no separate RBAC; access is governed by git permissions only in P0.
- **PHI/PII handling:** We do not store PII in evidence by design, but customers must ensure that any data they add (e.g. custom fields or tool output) is compliant.
- **Uptime or availability:** Evidence is stored in the customer's repo; availability and backup are the customer's responsibility.

---

## Reference

- Evidence schema: [schema/evidence.schema.json](../../schema/evidence.schema.json)
- F054: SDP Evidence Layer (schema and event types)
