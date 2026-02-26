# SDP Evidence Layer — Threat Model

Lightweight threat model for the Spec-Driven Protocol (SDP) evidence layer. Complements [COMPLIANCE.md](COMPLIANCE.md).

---

## Assets

What SDP helps protect or document:

- **Evidence integrity:** The evidence log (`.sdp/log/` or configured path) records workstream execution, approvals, and verification. Integrity means the log can be checked for accidental corruption (hash chain), not that it cannot be altered by an attacker with repo access.
- **Code provenance:** Evidence links code changes to workstreams, commits, and (where stored) prompt hashes. This supports audit trails for AI-generated code, not cryptographic attestation.
- **Audit trail:** A readable history of what was built, when, and by which workflow — for SOC2, DORA, and similar requirements.

---

## Threats

| Threat | Description |
|--------|-------------|
| **Tampering** | A user or process with write access to the repository can modify or delete evidence. The hash chain detects corruption after the fact; it does not prevent modification. |
| **Impersonation** | In P0, "who" in approval events comes from git config (e.g. `user.name`), not from a separate identity store. A user with repo access can set any name. |
| **Data loss** | Evidence lives in the repo. If the repo is lost, corrupted, or pruned, evidence is lost. Backups and retention are the customer's responsibility. |
| **Privacy leak** | Raw prompts are not stored; only hashes and metadata. If custom fields or tool output in evidence contain sensitive data, the customer must configure exclusions or redaction. |

---

## Mitigations

What SDP does today (P0) and on the roadmap:

| Threat | P0 (current) | P2 / P3 |
|--------|--------------|---------|
| **Tampering** | Hash chain for corruption detection; honest documentation that we do not guarantee tamper-proof. | P3: Signed evidence records for non-repudiation. |
| **Impersonation** | Git config as identity; no separate RBAC. | P2: Enterprise RBAC; authenticated identity in approval events. |
| **Data loss** | Evidence stored in customer repo; customer controls backup and retention. | No change; responsibility stays with customer. |
| **Privacy leak** | No raw prompts in evidence; prompt_hash only; customer controls tool output. | Configurable redaction/exclusions where needed. |

---

## Accepted Risks (P0)

What SDP does **not** mitigate in the initial release:

- **Repo admin can alter or delete evidence.** Anyone with write access to the repository can modify the evidence log. We detect accidental corruption; we do not prevent or attest to absence of deliberate changes.
- **No non-repudiation.** Evidence records are not cryptographically signed. We do not guarantee that a given event was produced by a specific identity or was not modified after creation.
- **No separate access control.** Access to evidence is governed by git permissions only. There is no RBAC layer within SDP in P0.
- **Identity is best-effort.** Approval events record git config identity; it can be spoofed by anyone with repo access.

---

## Future Mitigations

- **P3 — Signed evidence records:** Cryptographic signatures on evidence events for compliance-grade non-repudiation.
- **P3 — External timestamping:** Optional integration with timestamping authorities for provable ordering.
- **P2 — RBAC:** Separate roles (e.g. auditor vs. developer) and authenticated identity in approval events.

---

## Reference

- [COMPLIANCE.md](COMPLIANCE.md) — Data residency, retention, access control, regulatory alignment
- [schema/evidence.schema.json](../../schema/evidence.schema.json) — Evidence event schema
