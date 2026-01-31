---
feature: F034
status: completed
title: A+ Quality Initiative
---

# F034: A+ Quality Initiative

**Goal:** Довести SDP до уровня A+ по коду, документации, архитектуре и UX.

**Source:** Repository analysis (2026-01-31) — Expert panel review

## Current State Assessment

| Aspect | Current | Target | Gap |
|--------|---------|--------|-----|
| **Code Quality** | B (68% cov, 25 files >200 LOC) | A+ (80%+ cov, 0 violations) | HIGH |
| **Documentation** | B+ (inconsistencies) | A+ (consistent, no broken links) | MEDIUM |
| **Architecture** | B+ (beads→core violation) | A+ (Clean Architecture) | MEDIUM |
| **UX** | B- (naming confusion) | A+ (clear, discoverable) | HIGH |

## Workstreams

| WS | Title | Complexity | Dependencies |
|----|-------|------------|--------------|
| 00-034-01 | Split Large Files (Phase 1: Core) | MEDIUM | None |
| 00-034-02 | Split Large Files (Phase 2: Beads/Unified) | MEDIUM | 00-034-01 |
| 00-034-03 | Increase Test Coverage to 80%+ | LARGE | 00-034-02 |
| 00-034-04 | Documentation Consistency | SMALL | None |
| 00-034-05 | Extract Domain Layer | MEDIUM | 00-034-01 |
| 00-034-06 | Add `sdp status` Command | MEDIUM | None |
| 00-034-07 | Add Skill Discovery | SMALL | None |

## Dependency Graph

```
00-034-04 (docs) ────────────────────────────┐
                                              │
00-034-01 (split core) ──┬── 00-034-02 (split beads) ──┬── 00-034-03 (coverage)
                         │                              │
                         └── 00-034-05 (domain layer) ──┘
                                              
00-034-06 (sdp status) ──────────────────────┤
                                              │
00-034-07 (skill discovery) ─────────────────┘
```

## Success Criteria

- [ ] 0 files >200 LOC in src/sdp/
- [ ] Test coverage ≥80% (verified by pytest --cov)
- [ ] All docs use PP-FFF-SS format consistently
- [ ] No broken links in documentation
- [ ] Version aligned to v0.6.0 everywhere
- [ ] Clean Architecture: domain layer extracted, no beads→core imports
- [ ] `sdp status` command works
- [ ] Skill discovery available via `@help` or `sdp skills`

## Risk Assessment

| Risk | Mitigation |
|------|------------|
| Refactoring breaks existing tests | Run full test suite after each WS |
| Coverage increase is time-consuming | Focus on low-coverage critical paths first |
| Domain extraction changes APIs | Maintain backward compatibility, add deprecation warnings |

## Priority Order

1. **00-034-04** (docs) — Quick win, no code changes
2. **00-034-01** (split core) — Foundation for other refactoring
3. **00-034-05** (domain layer) — Architecture fix
4. **00-034-02** (split beads) — Continue refactoring
5. **00-034-03** (coverage) — After files are smaller
6. **00-034-06** (sdp status) — UX improvement
7. **00-034-07** (skill discovery) — Final polish

---

**Created:** 2026-01-31
**Owner:** AI Agent
