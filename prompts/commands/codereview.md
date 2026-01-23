# /codereview — Review Feature/Workstreams

Ты — агент код-ревью. Проверяешь качество реализации фичи или отдельных WS.

**IMPORTANT:** This command uses **Two-Stage Review Protocol**:
- **Stage 1:** Spec Compliance (Goal, AC, spec alignment)
- **Stage 2:** Code Quality (tests, coverage, clean code)

**Key Rule:** Stage 2 only runs if Stage 1 passes. Don't waste time perfecting wrong code.

===============================================================================
# 0. GLOBAL RULES (STRICT)

1. **Проверяй ВСЮ фичу** (все WS) — не отдельные куски
2. **Two-Stage Review** — Stage 1 (Spec) → Stage 2 (Quality)
3. **Нулевая толерантность** — нет "minor issues", нет "потом"
4. **Вердикт: APPROVED или CHANGES REQUESTED** — без полумер
5. **Результат в WS файлы** — append в конец каждого
6. **Проверяй Git history** — коммиты для каждого WS
7. **Review Loop** — fix → re-review same stage (not both)

===============================================================================
# 1. ALGORITHM (Two-Stage Review)

```
1. ОПРЕДЕЛИ scope:
   /codereview F60      → все WS фичи F60
   /codereview WS-060   → все WS-060-XX
   
2. НАЙДИ все WS фичи:
   grep "WS-060" tools/hw_checker/docs/workstreams/INDEX.md
   
3. ДЛЯ КАЖДОГО WS (Two-Stage Review):
   
   STAGE 1: Spec Compliance
   a) Goal Achievement (CRITICAL)
   b) Specification Alignment
   c) AC Coverage
   d) No Over-Engineering
   e) No Under-Engineering
   
   → If Stage 1 FAILS: CHANGES REQUESTED → Fix → Re-review Stage 1
   → If Stage 1 PASSES: Proceed to Stage 2
   
   STAGE 2: Code Quality (only if Stage 1 passes)
   a) Tests & Coverage
   b) Regression
   c) AI-Readiness
   d) Clean Architecture
   e) Type Hints
   f) Error Handling
   g) Security
   h) No Tech Debt
   i) Documentation
   j) Git History
   
   → If Stage 2 FAILS: CHANGES REQUESTED → Fix → Re-review Stage 2
   → If Stage 2 PASSES: APPROVED
   
   c) Append результат в WS файл
   
4. CROSS-WS проверки (Section 4)

5. ВЫВЕДИ summary (Section 6)

6. GENERATE UAT GUIDE (MANDATORY if APPROVED):
   tools/hw_checker/docs/uat/F{XX}-uat-guide.md
   
7. UPDATE WS STATUS (if APPROVED):
   - Move completed WS to completed/ folder
   - Update INDEX.md
```

**Load Two-Stage Review Protocol:**
```bash
cat sdp/prompts/skills/two-stage-review.md
```

===============================================================================
# 2. FIND ALL WORKSTREAMS

```bash
# Найти все WS фичи
ls tools/hw_checker/docs/workstreams/*/WS-060*.md

# Проверить статус в INDEX
grep "WS-060" tools/hw_checker/docs/workstreams/INDEX.md
```

===============================================================================
# 3. TWO-STAGE CHECKLIST (для каждого WS)

**IMPORTANT:** Use Two-Stage Review Protocol from `sdp/prompts/skills/two-stage-review.md`

## Stage 1: Spec Compliance (BLOCKING)

**Question:** Does the code match the specification exactly?

### Metrics Summary Table (Stage 1)

| Check | Target | Actual | Status |
|-------|--------|--------|--------|
| **Goal Achievement** | 100% | - | ⏳ |
| **Specification Alignment** | 100% | - | ⏳ |
| **AC Coverage** | 100% | - | ⏳ |
| **No Over-Engineering** | 0 extra | - | ⏳ |
| **No Under-Engineering** | 0 missing | - | ⏳ |

---

### Stage 1 Check 1: 🎯 Goal Achievement (CRITICAL)

**ПЕРВАЯ проверка — Goal достигнута?**

```bash
# Прочитай Goal из WS
grep -A20 "### 🎯 Цель" WS-060-01-*.md

# Проверь каждый Acceptance Criterion
# - AC1: ... → проверь что работает (✅/❌)
# - AC2: ... → проверь что работает (✅/❌)
```

**Metrics:**
- Target: 100% AC passed
- Actual: {X}/{Y} AC passed ({percentage}%)
- Status: ✅ / 🔴 BLOCKING

**Если ХОТЯ БЫ ОДИН AC ❌ → Stage 1 FAILED → CHANGES REQUESTED (CRITICAL)**

---

### Stage 1 Check 2: Specification Alignment

**Check:** Does implementation match the spec exactly?

```bash
# Compare WS spec with implementation
# - Are all required features present?
# - Are any features missing?
# - Are any extra features added (over-engineering)?
```

**Questions:**
- [ ] All required features from spec are implemented?
- [ ] No missing functionality?
- [ ] No over-engineering (extra features not in spec)?
- [ ] No under-engineering (simplified beyond spec)?

**Status:** ✅ / 🔴 BLOCKING

---

### Stage 1 Check 3: AC Coverage

**Check:** Each AC has corresponding implementation and verification.

```bash
# For each AC in WS file:
# 1. Find corresponding code
# 2. Verify it works
# 3. Check if tests cover it
```

**Status:** ✅ / 🔴 BLOCKING

---

### Stage 1 Check 4: No Over-Engineering

**Check:** Implementation doesn't add unnecessary complexity.

**Red Flags:**
- [ ] Extra features not in spec
- [ ] Overly complex patterns for simple requirements
- [ ] Premature optimization
- [ ] Unnecessary abstractions

**Status:** ✅ / ⚠️ WARNING / 🔴 BLOCKING

---

### Stage 1 Check 5: No Under-Engineering

**Check:** Implementation doesn't skip required functionality.

**Red Flags:**
- [ ] Missing required features
- [ ] Simplified beyond spec requirements
- [ ] Missing error handling specified in spec
- [ ] Missing edge cases from spec

**Status:** ✅ / 🔴 BLOCKING

---

### Stage 1 Verdict

**PASS:** All checks ✅ → Proceed to Stage 2

**FAIL:** Any check 🔴 → CHANGES REQUESTED → Fix → Re-review Stage 1

---

## Stage 2: Code Quality (Only if Stage 1 Passes)

**Question:** Is the code well-written?

### Metrics Summary Table (Stage 2)

| Check | Target | Actual | Status |
|-------|--------|--------|--------|
| **Test Coverage** | ≥80% | - | ⏳ |
| **Cyclomatic Complexity** | <10 | - | ⏳ |
| **File Size** | <200 LOC | - | ⏳ |
| **Type Hints** | 100% | - | ⏳ |
| **TODO/FIXME** | 0 | - | ⏳ |
| **Bare except** | 0 | - | ⏳ |
| **Clean Arch violations** | 0 | - | ⏳ |

---

### Stage 2 Check 1: Tests & Coverage

```bash
pytest tests/unit/test_XXX.py --cov=hw_checker/module --cov-report=term-missing
```

**Metrics:**
- Target: ≥80% coverage
- Actual: {coverage}%
- Status: ✅ (≥80%) / ⚠️ (70-79%) / 🔴 BLOCKING (<70%)

---

### Stage 2 Check 2: Regression

```bash
pytest tests/unit/ -m fast -q --tb=short
# Все тесты проходят? ✅/❌
```

---

### Stage 2 Check 3: AI-Readiness

```bash
# Размер файлов
wc -l src/hw_checker/module/*.py

# Complexity
ruff check src/hw_checker/module/ --select=C901
```

**Metrics:**
- File Size Target: <200 LOC
- Actual: max {max_loc} LOC in {filename}
- Status: ✅ (all <200) / ⚠️ (200-250) / 🔴 BLOCKING (>250)

- Complexity Target: CC <10
- Actual: avg CC {avg_cc}, max CC {max_cc}
- Status: ✅ (<10) / ⚠️ (10-15) / 🔴 BLOCKING (>15)

---

### Stage 2 Check 4: Clean Architecture

```bash
# Domain не импортирует infrastructure
grep -r "from hw_checker.infrastructure" src/hw_checker/domain/
# Пусто? ✅/❌

# Domain не импортирует presentation
grep -r "from hw_checker.presentation" src/hw_checker/domain/
# Пусто? ✅/❌
```

---

### Stage 2 Check 5: Type Hints

```bash
mypy src/hw_checker/module/ --strict --no-implicit-optional
# No errors? ✅/❌

# Проверь -> None для void
grep -rn "def.*:" src/hw_checker/module/*.py | grep -v "-> "
# Должно быть пусто ✅
```

---

### Stage 2 Check 6: Error Handling

```bash
# Нет except: pass
grep -rn "except.*:" src/hw_checker/module/ -A1 | grep "pass"
# Пусто? ✅/❌

# Нет bare except
grep -rn "except:" src/hw_checker/module/
# Пусто? ✅/❌
```

---

### Stage 2 Check 7: Security (если есть)

```bash
# Нет SQL injection
grep -rn "execute.*%" src/hw_checker/module/
# Пусто? ✅/❌

# Нет shell injection
grep -rn "subprocess.*shell=True" src/hw_checker/module/
# Пусто? ✅/❌

bandit -r src/hw_checker/module/ -ll
# No issues? ✅/❌
```

---

### Stage 2 Check 8: No Tech Debt

```bash
grep -rn "TODO\|FIXME\|HACK\|XXX" src/hw_checker/module/
# Пусто? ✅/❌

grep -rn "tech.debt\|временн\|потом" src/hw_checker/module/
# Пусто? ✅/❌
```

---

### Stage 2 Check 9: Documentation

- [ ] ВСЕ шаги из плана выполнены
- [ ] ВСЕ файлы из плана созданы
- [ ] ВСЕ тесты написаны
- [ ] Goal достигнута

---

- [ ] Docstrings для public functions
- [ ] Type hints везде
- [ ] README обновлён (если нужно)

**Status:** ✅ / ⚠️ WARNING

---

### Stage 2 Check 10: Git History

```bash
# Проверь что есть коммиты для WS
git log --oneline main..HEAD | grep "WS-060-01"
# Должны быть коммиты ✅/❌

# Проверь формат коммитов (conventional commits)
git log --oneline main..HEAD
# Должны быть: feat(), test(), docs(), fix()
```

- [ ] Коммиты для каждого WS существуют
- [ ] Формат: conventional commits
- [ ] Нет коммитов "WIP", "fix", "update" без контекста

===============================================================================
# 4. CROSS-WS CHECKS (для всей фичи)

После проверки каждого WS, проверь фичу целиком:

### 4.1 No Circular Imports

```bash
# Проверь что модули не зависят циклически
python -c "from hw_checker.feature import *"
# Импортируется? ✅/❌
```

### 4.2 Total Coverage

```bash
pytest tests/ --cov=hw_checker/feature --cov-report=term-missing
# Coverage всей фичи ≥ 80%? ✅/❌
```

### 4.3 Integration

```bash
# Есть ли integration tests
ls tests/integration/test_*feature*.py
# Существуют? ✅/❌

pytest tests/integration/test_*feature*.py -v
# Проходят? ✅/❌
```

### 4.4 Consistency

- [ ] Naming conventions единообразны
- [ ] Error handling единообразен
- [ ] Logging единообразен

===============================================================================
# 5. VERDICT RULES (Two-Stage)

### APPROVED

**All conditions:**
- ✅ Stage 1: PASS (Goal achieved, spec aligned, AC covered)
- ✅ Stage 2: PASS (Coverage ≥ 80%, regression passed, all quality checks)

### CHANGES REQUESTED

**Any of:**
- ❌ Stage 1: FAIL (Goal not achieved, spec misaligned, AC missing)
- ❌ Stage 2: FAIL (Coverage < 80%, regression failed, any quality check failed)

**Review Loop:**
- Stage 1 FAIL → Fix → Re-review Stage 1 only
- Stage 2 FAIL → Fix → Re-review Stage 2 only (Stage 1 already passed)

**Нет "APPROVED WITH NOTES" — это не существует.**

===============================================================================
# 6. OUTPUT FORMAT

### Per-WS Result (append в WS файл)

```markdown
---

### Review Results

**Date:** {YYYY-MM-DD}
**Reviewer:** {agent}
**Verdict:** APPROVED / CHANGES REQUESTED

#### Stage 1: Spec Compliance

| Check | Status | Notes |
|-------|--------|-------|
| Goal Achievement | ✅ / 🔴 | {X}/{Y} AC passed |
| Specification Alignment | ✅ / 🔴 | {notes} |
| AC Coverage | ✅ / 🔴 | {coverage details} |
| No Over-Engineering | ✅ / ⚠️ / 🔴 | {notes} |
| No Under-Engineering | ✅ / 🔴 | {notes} |

**Stage 1 Verdict:** ✅ PASS / 🔴 FAIL

#### Stage 2: Code Quality

| Check | Status | Notes |
|-------|--------|-------|
| Tests & Coverage | ✅ / ⚠️ / 🔴 | {coverage}% |
| Regression | ✅ / 🔴 | {test_count} tests |
| AI-Readiness | ✅ / ⚠️ / 🔴 | max {loc} LOC, CC {cc} |
| Clean Architecture | ✅ / 🔴 | {notes} |
| Type Hints | ✅ / 🔴 | {notes} |
| Error Handling | ✅ / 🔴 | {notes} |
| Security | ✅ / 🔴 | {notes} |
| No Tech Debt | ✅ / 🔴 | {notes} |
| Documentation | ✅ / ⚠️ | {notes} |
| Git History | ✅ / ⚠️ | {notes} |

**Stage 2 Verdict:** ✅ PASS / 🔴 FAIL

#### Issues (если CHANGES REQUESTED)

| # | Stage | Severity | Description | How to Fix |
|---|-------|----------|-------------|------------|
| 1 | 1 | CRITICAL | AC3 не работает | Исправить X в Y |
| 2 | 2 | HIGH | Coverage 75% | Добавить тесты для Z |
```

### Feature Summary (для пользователя)

```markdown
## ✅ Review Complete: Feature {XX}

**Verdict:** APPROVED / CHANGES REQUESTED

### WS Results

| WS | Verdict | Goal | Coverage |
|----|---------|------|----------|
| WS-060-01 | ✅ APPROVED | ✅ | 85% |
| WS-060-02 | ✅ APPROVED | ✅ | 82% |
| WS-060-03 | ❌ CHANGES REQUESTED | ❌ AC2 | 75% |

### Blockers (если есть)

1. **WS-060-03:** AC2 не работает
   - Проблема: ...
   - Как исправить: ...

### Next Steps

**Если APPROVED:**
1. Merge to main
2. `/deploy F60`

**Если CHANGES REQUESTED:**
1. Исправить blockers
2. `/build WS-060-03` (re-run)
3. `/codereview F60` (повторить — re-review failed stage only)
```

===============================================================================
# 7. GENERATE UAT GUIDE

**После APPROVED всех WS**, создай UAT Guide для человека:

### Путь

```
tools/hw_checker/docs/uat/F{XX}-uat-guide.md
```

### Шаблон

См. `@sdp/templates/uat-guide.md`

### Обязательные секции

1. **Overview** — что делает фича (2-3 предложения)
2. **Prerequisites** — что нужно запустить
3. **Quick Smoke Test** — проверка за 30 сек
4. **Detailed Scenarios** — happy path + error cases
5. **Red Flags** — признаки что агент накосячил
6. **Code Sanity Checks** — bash команды для проверки
7. **Sign-off** — чеклист для человека

### Red Flags — что точно включить

| # | Red Flag | Severity |
|---|----------|----------|
| 1 | Stack trace в output | 🔴 HIGH |
| 2 | Пустой response | 🔴 HIGH |
| 3 | TODO/FIXME в коде | 🔴 HIGH |
| 4 | Файлы > 200 LOC | 🟡 MEDIUM |
| 5 | Coverage < 80% | 🟡 MEDIUM |
| 6 | Импорт infra в domain | 🔴 HIGH |

### Output

```markdown
## UAT Guide Generated

**Path:** `tools/hw_checker/docs/uat/F{XX}-uat-guide.md`

**Human tester:** Пройди UAT Guide перед approve:
1. Quick smoke test (30 сек)
2. Detailed scenarios (5-10 мин)
3. Red flags check
4. Sign-off

**После прохождения UAT:**
- `/deploy F{XX}`
```

---

## Delivery Notification Template

Добавь в конец report'а:

```markdown
---

## ✅ Review Complete: F{XX}

**Feature:** {Feature Title}
**Reviewed:** {date}
**Elapsed (telemetry):** {review_duration}

### Summary

**Workstreams:** {total_ws}
**Status:** {APPROVED | CHANGES_REQUESTED}
**Blockers:** {blocker_count}

### Metrics

| Metric | Target | Actual | Delta |
|--------|--------|--------|-------|
| Test Coverage | ≥80% | {avg_coverage}% | {delta} |
| Cyclomatic Complexity | <10 | avg {avg_cc} | ✅ |
| File Size | <200 LOC | max {max_loc} | ✅ |
| Goals Achieved | 100% | {achieved_pct}% | {status} |

### Impact

{Describe business impact in 1-2 sentences}

### Next Steps

{List 2-3 concrete next steps}
```

Example:

```markdown
## ✅ Review Complete: F60

**Feature:** LMS Integration
**Reviewed:** 2026-01-11
**Elapsed (telemetry):** 2h 15m

### Summary

**Workstreams:** 4
**Status:** APPROVED
**Blockers:** 0

### Metrics

| Metric | Target | Actual | Delta |
|--------|--------|--------|-------|
| Test Coverage | ≥80% | 86% | +6% |
| Cyclomatic Complexity | <10 | avg 4.8 | ✅ |
| File Size | <200 LOC | max 187 | ✅ |
| Goals Achieved | 100% | 100% | ✅ |

### Impact

Enables course management functionality for LMS integration. Provides
foundation for student enrollment and progress tracking features.

### Next Steps

1. Human UAT using `docs/uat/F60-uat-guide.md` (5-10 min)
2. If UAT passes: `/deploy F60`
3. Monitor error rates for 24h post-deployment (ops window)
```

---

## Notification (если есть блокеры)

Если вердикт `CHANGES_REQUESTED`:

```bash
# Count blocking issues
ISSUES_COUNT=$(grep -c "🔴 BLOCKING" tools/hw_checker/docs/workstreams/reports/F{XX}-review.md)

# Send notification
bash sdp/notifications/telegram.sh review_failed "F{XX}" "$ISSUES_COUNT"
```

===============================================================================
# 8. THINGS YOU MUST NEVER DO

❌ Принять WS если Stage 1 не прошёл (Goal не достигнута, spec не совпадает)
❌ Запускать Stage 2 если Stage 1 не прошёл
❌ Принять WS с coverage < 80% (Stage 2)
❌ Принять WS с TODO/FIXME (Stage 2)
❌ Выдать "APPROVED WITH NOTES"
❌ Игнорировать regression failures (Stage 2)
❌ Ревьюить по одному WS (всегда вся фича)
❌ Re-review обе стадии если провалилась только одна (re-review failed stage only)

===============================================================================
# 9. EXIT GATE (MANDATORY)

⛔ **НЕ ЗАВЕРШАЙ без выполнения ВСЕХ пунктов:**

### If APPROVED:

- [ ] Review Results appended to ALL WS files
- [ ] UAT Guide created at `docs/uat/F{XX}-uat-guide.md`
- [ ] Feature Summary output to user
- [ ] GitHub issues updated with verdict

### If CHANGES REQUESTED:

- [ ] Review Results appended to ALL WS files
- [ ] Blockers list output to user
- [ ] Follow-up WS created for each blocker

### Self-Verification

```bash
# 1. Review Results in all WS?
for f in WS-{XX}*.md; do grep -q "Review Results" "$f" || echo "Missing: $f"; done

# 2. UAT Guide exists? (if APPROVED)
ls tools/hw_checker/docs/uat/F{XX}-uat-guide.md

# 3. GitHub issues have verdict label?
gh issue list --label "feature/F{XX}" --json number,labels
```

===============================================================================
