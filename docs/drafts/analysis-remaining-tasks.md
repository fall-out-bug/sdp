# Remaining Tasks from SDP Analysis Design

> **Date:** 2026-01-28
> **Status:** Planning
> **Based on:** docs/plans/2025-01-26-sdp-analysis-design.md

---

## ✅ Completed

1. **F014: Workflow Efficiency** (100% done)
   - @oneshot execution modes (--auto-approve, --sandbox, --dry-run)
   - @idea two-round interview (3-5 critical questions + optional deep-dive)
   - Destructive operations detection
   - Audit logging

2. **PRODUCT_VISION.md** (exists)
   - Mission, Users, Success Metrics defined

3. **Intent Schema** (exists)
   - docs/schema/intent.schema.json created

---

## ❌ Not Started - Phase 1: Quick Wins (1-2 weeks)

### F015: AI-Human Communication Enhancement
**Priority:** HIGH
**Impact:** Helps AI understand "why" behind features

**Tasks:**
1. [ ] Update @idea skill to require PRODUCT_VISION.md check
2. [ ] Auto-load PRODUCT_VISION.md in all skills
3. [ ] Add vision validation to @idea workflow
4. [ ] Create decision log mechanism
5. [ ] Link WS execution reports to intent schema

**Expected Impact:**
- AI has context from product vision
- Better alignment with user goals
- Decision audit trail

---

### F016: Developer Experience - Unified Dashboard
**Priority:** HIGH
**Impact:** Single source of truth for project status

**Tasks:**
1. [ ] Create `sdp status` CLI command
2. [ ] Implement TUI dashboard (rich/textual)
3. [ ] Show IDEAS, WORKSTREAMS, FEATURES status
4. [ ] Add keyboard shortcuts ([n]ew, [d]esign, [b]uild, etc.)
5. [ ] Auto-refresh on file changes

**Expected Impact:**
- No more manual file navigation
- Clear visibility into project state
- Faster decision making

---

### F017: Documentation - English Translation & Tutorial
**Priority:** MEDIUM
**Impact:** Removes language barrier, improves onboarding

**Tasks:**
1. [ ] Create PROTOCOL_EN.md (English translation)
2. [ ] Create TUTORIAL.md (15-minute walkthrough)
3. [ ] Use minimal example (fix bug in SDP)
4. [ ] Add checkpoints with expected output
5. [ ] Celebrate completion

**Expected Impact:**
- International users can use SDP
- New users onboard in <30 min
- "I can do this!" moment

---

### F018: Documentation Consolidation
**Priority:** LOW
**Impact:** Easier navigation

**Tasks:**
1. [ ] Create NAVIGATION.md
2. [ ] Consolidate L1-L4 documentation
3. [ ] Add decision trees (@build vs @oneshot)
4. [ ] Simplify "Quick Start" (3 commands)

**Expected Impact:**
- Single entry point for docs
- Progressive disclosure
- Reduced cognitive load

---

### F019: Quality Gates Evolution
**Priority:** LOW
**Impact:** Less blocking, more practical

**Tasks:**
1. [ ] Change 200 LOC limit from hard block to warning
2. [ ] Allow TODO with WS-ID references
3. [ ] Update pre-build hook
4. [ ] Update post-build hook

**Expected Impact:**
- Fewers false positives
- TODO markers trackable to workstreams
- More pragmatic quality enforcement

---

## ❌ Not Started - Phase 2-5: Longer Term

These require more investment (3-6 weeks each):

### Phase 2: DX Foundation (3-4 weeks)
- [ ] Auto-file management in @build skill
- [ ] Enhanced intent validation
- [ ] Update @idea skill with schema validation

### Phase 3: Fast Feedback (2-3 weeks)
- [ ] Implement `sdp test --watch` command
- [ ] Enhance pre-commit hook with incremental validation
- [ ] Add fast test markers

### Phase 4: Quality Evolution (4-6 weeks)
- [ ] `sdp metrics` command (Change Failure Rate)
- [ ] `sdp hotspots` command
- [ ] Language profile system (Rust, TypeScript)

### Phase 5: Robustness (3-4 weeks)
- [ ] F013: Transactional workstream execution
- [ ] Branch-per-attempt model
- [ ] `sdp rollback WS-ID` command
- [ ] Team coordination primitives

---

## Recommendation

**Start with F015: AI-Human Communication Enhancement**

**Why:**
1. Builds on existing PRODUCT_VISION.md and intent schema
2. High impact (AI understanding "why")
3. Medium complexity (1-2 weeks)
4. Unblocker for other features (skills need vision context)

**Next:** F016 (Unified Dashboard) or F017 (English + Tutorial)

---

## Success Metrics (from analysis)

| Metric | Baseline | Target | Status |
|--------|----------|--------|--------|
| Time to first running code | ~3h 45m | <1 hour | 🔄 F014 done, F015-019 pending |
| Test feedback speed | Manual (30-60s) | <2 seconds | ❌ Watch mode not done |
| New user onboarding | Overwhelmed | <30 min | ❌ Tutorial not done |
| Cognitive load | 10 commands | 3 commands basic | ❌ Docs not consolidated |
| Cycle time (idea → deployed) | 3h 45m | <45 min | ✅ F014 done (5x faster) |
| Change Failure Rate | Not measured | <5% | ❌ Metrics not implemented |

---

**Next Action:** Choose feature to implement (F015-F019)
