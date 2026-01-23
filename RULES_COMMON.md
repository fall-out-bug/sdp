# Common Rules for All Agents

**Version:** 1.0  
**Applies to:** All consensus agents (development + education)

---

## 1. Language & Format

- **ALL messages MUST be in English** — no Russian or other languages in any field
- **JSON format only** for inbox messages
- **Naming scheme:** `{YYYY-MM-DD}-{subject}.json`
- **Compact keys:** d, st, r, epic, sm, nx, artifacts, answers

## 2. Inbox Rules

| Action | Rule |
|--------|------|
| READ | Only your own inbox: `messages/inbox/{your_role}/` |
| WRITE | Only to OTHER agents' inboxes |
| Self-reference | Status summaries allowed in own inbox |

**Never:**
- Read other agents' inboxes
- Write operational messages to your own inbox

## 3. Message Schema

```json
{
  "d": "2025-12-09",          // date (required)
  "st": "status|request|veto|approval|handoff",  // type (required)
  "r": "developer",           // role (required)
  "epic": "EP08",             // epic ID (required)
  "sm": ["summary points"],   // summary (required, array)
  "nx": ["next actions"],     // next steps (optional)
  "artifacts": ["paths"],     // artifact paths (optional)
  "answers": {"Q1": "A1"}     // Q&A responses (optional)
}
```

## 4. Quality Gates

| Gate | Rule |
|------|------|
| No silent fallbacks | All errors must be explicit, logged, or raised |
| No layer violations | Dependencies point inward (Domain ← App ← Infra ← Presentation) |
| Test coverage | ≥80% in touched areas |
| Documentation | Update relevant docs after code changes |

## 5. Self-Verification Checklist

Before completing any work:
- [ ] Clean Architecture boundaries respected
- [ ] Engineering principles followed (DRY, SOLID, Clean Code)
- [ ] No fallbacks hiding errors
- [ ] Documentation updated
- [ ] All messages in English

## 6. Veto Protocol

When issuing a veto:
```json
{
  "d": "2025-12-09",
  "st": "veto",
  "r": "architect",
  "epic": "EP08",
  "sm": [
    "Violation: layer_violation",
    "Location: application/foo.py imports infrastructure",
    "Impact: Breaks dependency inversion"
  ],
  "nx": ["Remove direct import", "Use port injection"]
}
```

**Cannot override vetoes for:**
- Architecture violations (architect)
- Security issues (security)
- Missing rollback plan (devops)
- Code review violations (tech_lead, qa)

## 7. Cross-Epic Context

At epic completion, preserve context for next epic:
- Key decisions
- Technical debt identified
- Lessons learned
- Foundation for next work

---

**Reference this file in prompts:**
```json
"common_rules": "See docs/roles/RULES_COMMON.md"
```
