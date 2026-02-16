# Анализ: Enforcement существующего протокола

> **Status:** Revised after user feedback
> **Date:** 2026-02-16
> **Goal:** Минимальные изменения для enforcement существующего протокола

---

## Ключевой инсайт

**Протокол УЖЕ СУЩЕСТВУЕТ.** Проблема не в отсутствии функционала, а в **enforcement** (приведении в исполнение).

### Существующий протокол

```
@oneshot F067  →  @review F067  →  @deploy F067
    │                 │                │
    ▼                 ▼                ▼
 Execute WS      APPROVED?         Merge PR
                  │
                  ├─ YES → proceed
                  └─ NO → fix loop
```

### Что было нарушено

| Шаг | Протокол | Что случилось |
|-----|----------|---------------|
| 1 | @oneshot грузит WS из backlog/ | OK |
| 2 | @review проверяет feature | **ПРОПУЩЕН** |
| 3 | @deploy только после APPROVED | **ПРОПУЩЕН** - сливал PR напрямую |
| 4 | WS status update | Делал вручную, не по протоколу |

---

## Root Cause

### Почему протокол не работал

1. **Context Loss** - После compaction забывал про roadmap
2. **Нет enforcement gate** - Ничто не блокировало merge без @review
3. **Skill invocation optional** - Мог пропустить @review и ничего не сломалось
4. **Done = PR merged** - Привычка считать "done" по PR, не по verdict

### Почему предложенные решения были неверны

| Предложение | Почему неверно |
|-------------|----------------|
| Новый @milestone skill | Протокол уже есть |
| .sdp/milestones.json | Избыточно |
| sdp status reconcile | WS manage должен быть в протоколе |
| sdp guard feature-complete | @review уже это делает |

---

## Минимальные исправления

### 1. Context Preservation (CLAUDE.md)

**Проблема:** После compaction теряется milestone context.

**Решение:** Добавить секцию в CLAUDE.md (уже загружается каждый session):

```markdown
## Milestone Context

Current milestone: **M1 "T-shirt"**

M1 Features: F054, F063, F064, F067, F068, F070, F075, F076
M2 Features: F060, F071, F073, F077, F078
M3 Features: F057, F058, F069, F072, F074, F079
M4 Features: F055, F056, F059, F061

⚠️ Only work on current milestone features unless explicitly requested.
```

**Implementation:** 5 минут, просто добавить секцию.

### 2. Session-Start Pattern

**Проблема:** Новая сессия не проверяет roadmap.

**Решение:** Создать `.claude/patterns/session-start.md`:

```markdown
# Session Start Protocol

Before any work:
1. Read current milestone from CLAUDE.md
2. Check CHANGELOG for recent changes
3. Verify: Does the work belong to current milestone?
```

**Implementation:** 5 минут, создать файл.

### 3. Review Gate (enforcement)

**Проблема:** Ничто не блокирует PR merge без @review APPROVED.

**Решение A (CI):** Добавить GitHub Action:

```yaml
# .github/workflows/review-gate.yml
name: Review Gate
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  check-review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Check review verdict
        run: |
          if [ -f .sdp/review_verdict.json ]; then
            verdict=$(jq -r '.verdict' .sdp/review_verdict.json)
            if [ "$verdict" != "APPROVED" ]; then
              echo "ERROR: Feature not reviewed. Run @review first."
              exit 1
            fi
          fi
```

**Решение B (Skill update):** Обновить @deploy skill:

```markdown
## Pre-merge Check

Before merging:
1. Check if .sdp/review_verdict.json exists
2. If not: "Run @review first"
3. If exists but not APPROVED: "Fix issues first"
4. Only proceed if APPROVED
```

**Implementation:** 15-30 минут.

### 4. @review Output (Evidence)

**Проблема:** @review вердикт не персистится.

**Решение:** @review должен создавать `.sdp/review_verdict.json`:

```json
{
  "feature": "F067",
  "verdict": "APPROVED",
  "timestamp": "2026-02-16T10:00:00Z",
  "reviewers": ["qa", "security", "devops", "sre", "techlead", "docs"],
  "summary": "All checks passed"
}
```

**Implementation:** Обновить @review skill, добавить сохранение файла.

---

## Implementation Plan

### Phase 1: Quick Wins (30 минут)

- [x] ~~Analysis document~~
- [x] Добавить Milestone Context в CLAUDE.md
- [x] Создать session-start.md pattern

### Phase 2: Enforcement (1 час)

- [x] Обновить @review skill: сохранять verdict в .sdp/review_verdict.json
- [x] Обновить @deploy skill: проверять verdict перед merge
- [ ] Добавить review-gate.yml GitHub Action

### Phase 3: Verification

- [ ] Протестировать полный flow: @oneshot → @review → @deploy
- [ ] Убедиться что @deploy блокируется без APPROVED

---

## Что НЕ делаем

| Отказались | Почему |
|------------|--------|
| @milestone skill | Избыточно, CLAUDE.md достаточно |
| .sdp/milestones.json | Дублирует ROADMAP |
| sdp status reconcile | WS manage в протоколе |
| sdp guard feature-complete | @review уже это делает |
| milestone field в Beads | Не нужно для enforcement |

---

## Success Criteria

| Metric | Before | After |
|--------|--------|-------|
| PR без @review | Возможен | Заблокирован |
| Milestone check | Ручной | В CLAUDE.md |
| Session context | Теряется | Восстанавливается |
| Review evidence | Нет | .sdp/review_verdict.json |

---

## Итог

**Проблема:** Протокол существовал, но не был enforced.

**Решение:** Минимальные изменения:
1. Context в CLAUDE.md
2. Session-start pattern
3. Review verdict gate

**Без:** Новых skills, новых команд, нового кода.

---

*Ready to implement minimal fixes.*
