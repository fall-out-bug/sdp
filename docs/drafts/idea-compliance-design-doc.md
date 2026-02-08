# F055: Compliance Design Doc

> Beads: sdp-raho | Priority: P0

---

## Problem

Enterprise customers (bank, airline, marketplace evaluating) need answers before signing contracts:
- Where does the evidence live?
- Who has access?
- What does the hash chain actually guarantee?
- Are raw prompts ever stored?
- How does this align with SOC2/HIPAA/DORA?

Without a clear document, every sales conversation reinvents these answers.

## Solution

A single document answering all enterprise compliance questions. Not code — a reference document that sales, legal, and security teams can use.

## Sections

### 1. Data Residency
- Evidence log lives in the customer's git repo (`.sdp/log/`)
- No data leaves the repo unless customer configures export
- Telemetry (anonymous usage stats) goes to SDP servers — opt-out available
- Raw code never transmitted by SDP itself (model providers handle that separately)

### 2. Retention & GDPR
- Evidence log = git history = customer controls retention
- No personal data in evidence (model IDs, hashes, pass/fail — no PII)
- GDPR right to erasure: evidence log can be pruned with `git filter-branch`
- Recommended: retention policy in `.sdp/config.yml`

### 3. RBAC & Access Control
- Evidence log inherits git repo permissions
- No separate access layer in P0 (enterprise RBAC in P2)
- `sdp log trace` requires repo read access — no escalation
- Approval events record "who" from git config (not authenticated identity in P0)

### 4. Integrity Guarantees (Honest Labeling)
- Hash chain = corruption detection
- Hash chain ≠ tamper-proof (anyone with repo write can modify)
- Hash chain ≠ non-repudiation (no cryptographic signatures in P0)
- What it catches: accidental corruption, partial writes, lost events
- What it doesn't catch: deliberate modification by repo admin
- P3 roadmap: signed evidence records for compliance-grade non-repudiation

### 5. Prompt Privacy
- Raw prompts NEVER stored in evidence log
- `prompt_hash` (SHA-256) stored — allows matching without disclosure
- Model parameters (temperature, max_tokens) stored — no privacy concern
- Customer's code in evidence: only file paths and code hashes, not code content
- Verification output: actual tool stdout (test results) — customer controls sensitivity

### 6. Regulatory Alignment
- **SOC2**: evidence log = audit trail for AI-generated code changes
- **HIPAA**: no PHI in evidence by default; customer must configure exclusions if code touches PHI
- **DORA**: evidence supports ICT risk management documentation
- **EU AI Act**: evidence log supports transparency requirements for AI-generated artifacts

## Users

- Enterprise security teams (pre-contract evaluation)
- Legal/compliance reviewers
- Sales engineering (answering RFPs)
- SDP developers (ensuring promises match implementation)

## Success Metrics

- Document used in at least 2 enterprise conversations
- No compliance question from prospects that isn't covered
- Security team at bank signs off on evidence approach

## Dependencies

- F054 evidence schema (must be designed to write accurate compliance doc)

## Deliverables

1. `docs/compliance/COMPLIANCE.md` — the reference document
2. `docs/compliance/THREAT-MODEL.md` — what SDP protects against and what it doesn't
3. Update CLAUDE.md with compliance doc reference
