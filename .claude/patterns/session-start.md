# Session Start Protocol

Before starting any work, follow this checklist:

## 1. Check Current Milestone

Read CLAUDE.md "Milestone Context" section.

**Current milestone:** M1 "T-shirt"

**M1 Features:** Evidence layer, guard, repository hardening, UX foundation, error recovery, self-healing doctor, guided onboarding

## 2. Check Recent Changes

```bash
cat CHANGELOG.md | head -50
```

What was done recently? What's pending?

## 3. Verify Alignment

Before working on any feature or workstream:

- Does this belong to the current milestone?
- If NO: Are you explicitly requested to work on it?
- If NO: Ask user or work on current milestone instead

## 4. Protocol Flow Check

Remember the correct flow:

```
@oneshot → @review → @deploy
```

**NOT:** @oneshot → merge PR (skipping @review)

## 5. "Done" Definition

Feature is "done" when:
- [ ] @review returns APPROVED
- [ ] @deploy completed successfully
- [ ] All workstreams have status: completed

**NOT when:** PR is merged (that's just one step)

---

## Quick Reference

| Check | Command/Action |
|-------|---------------|
| Current milestone | Read CLAUDE.md header |
| Recent changes | `head -50 CHANGELOG.md` |
| Ready work | `bd ready` |
| Blocked work | `bd blocked` |
| Active WS | `ls docs/workstreams/backlog/*.md` |
