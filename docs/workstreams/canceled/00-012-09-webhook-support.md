---
ws_id: 00-012-09
project_id: 00
feature: F012
status: superseded
superseded_by: 00-032-01
supersede_reason: "F012 daemon/agent framework superseded by F032 Guard + Beads integration"
size: MEDIUM
github_issue: null
assignee: null
started: null
completed: null
blocked_reason: null
---

## 00-012-09: Webhook Support

### 🎯 Goal

**What must WORK after this WS is complete:**
- `sdp webhook server` starts HTTP server on port 8080
- Webhook receives `issues` and `project_v2` events
- Webhook validates GitHub signature (X-Hub-Signature-256)
- Webhook triggers `sync_service.sync_workstream()` on issue update
- Webhook logs all events to `.sdp/webhook.log`
- `--smee-url` flag for tunneling (local dev)
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `sdp webhook server` starts HTTP server on port 8080
- [ ] AC2: Webhook receives `issues` and `project_v2` events
- [ ] AC3: Webhook validates GitHub signature (X-Hub-Signature-256)
- [ ] AC4: Webhook triggers `sync_service.sync_workstream()` on issue update
- [ ] AC5: Webhook logs all events to `.sdp/webhook.log`
- [ ] AC6: `--smee-url` flag for tunneling (local dev)
- [ ] AC7: Coverage ≥ 80%
- [ ] AC8: mypy --strict passes

---

### Context

Real-time sync requires webhook support for immediate updates when issues change on GitHub. This eliminates polling delay.

---

### Dependencies

00-012-07 (GitHub Project Fields Integration)

---

### Steps

1. Create `src/sdp/webhook/__init__.py` with exports
2. Create `src/sdp/webhook/signature.py` for signature validation
3. Create `src/sdp/webhook/handler.py` for webhook handler
4. Create `src/sdp/webhook/server.py` for HTTP server
5. Modify `pyproject.toml` to add `starlette` or `aiohttp` dependency
6. Add `@main.group() def webhook()` to `src/sdp/cli.py` with `server` subcommand
7. Create unit tests for webhook
8. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/webhook/
├── __init__.py
├── server.py         (~300 LOC)
├── handler.py        (~250 LOC)
└── signature.py      (~100 LOC)

tests/unit/webhook/
├── test_handler.py    (~150 LOC)
└── test_signature.py  (~100 LOC)
```

**Modified Files:**
- `pyproject.toml` - Add starlette/aiohttp dependency (~5 LOC)
- `src/sdp/cli.py` - Add webhook command group (~50 LOC)

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/webhook/ -v
pytest --cov=sdp.webhook --cov-fail-under=80

# Type checking
mypy src/sdp/webhook/ --strict

# Manual test
sdp webhook server --port 8080
```

---

### Constraints

- USING `SyncService` from WS-00-012-03
- USING existing GitHub client from `github/client.py`
- FOLLOWING webhook patterns from GitHub docs
- POLLING fallback if webhook fails

---

### Scope Estimate

- **Files:** 8 created/modified
- **Lines:** ~900 LOC
- **Size:** MEDIUM (500-1500 LOC)
