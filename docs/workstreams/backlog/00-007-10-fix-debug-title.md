---
id: WS-201-10
title: Fix WS-201-04 title mismatch (4-phase vs 5-phase)
feature: F007
status: backlog
size: TINY
github_issue: TBD
dependencies:
  - WS-201-04 # /debug command for Cursor and OpenCode
---

## 02-201-10: Fix WS-201-04 title mismatch (4-phase vs 5-phase)

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Title WS-201-04 соответствует реальной реализации (5-phase workflow)
- Нет путаницы между названием и содержанием
- Документация согласована

**Acceptance Criteria:**
- [ ] Title WS-201-04 изменен с "4-phase debugging workflow" → "5-phase debugging workflow"
- [ ] Описание соответствует 5-phase workflow (Symptom → Hypothesis → Test → Root Cause → Impact)
- [ ] Другие части WS (шаги, код) не изменены (только title)
- [ ] Code review note обновлен (описать что title исправлен)

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

**Проблема:**
- Title WS-201-04: "/debug command for Cursor and OpenCode"
- Описание title: "Systematic debugging workflow (4-phase)"
- Реализация использует **5-phase workflow** (не 4-phase)

**Расхождение:**
- Spec говорит: "4-phase debugging workflow"
- Master prompt (`sdp/prompts/commands/issue.md` Section 4.0) описывает: **5 фаз**
  1. Symptom Documentation
  2. Hypothesis Formation
  3. Systematic Elimination
  4. Root Cause Isolation
  5. Impact Chain Analysis

**Code Review Note (из WS-201-04):**
> WS title mentions "4-phase debugging workflow" but implementation uses "5-phase workflow" (Symptom → Hypothesis → Test → Root Cause → Impact). This is a documentation discrepancy only; functionality is correct.

**Решение:**
- Исправить title/описание на "5-phase debugging workflow"
- Оставить реализацию без изменений (она правильная)
- Добавить note в code review что title исправлен

---

### Зависимость

WS-201-04 (команда /debug создана и работает)

---

### Входные файлы

- `tools/hw_checker/docs/workstreams/backlog/WS-201-04-debug-command.md`
- Code review notes в WS-201-04 (Execution Report)

---

### Шаги

1. **Найти где упоминается "4-phase":**
   - Грепнуть WS-201-04.md для "4-phase"
   - Проверить title, описание, шаги

2. **Заменить на "5-phase":**
   - Заменить все упоминания "4-phase debugging workflow" → "5-phase debugging workflow"
   - Оставить остальной контент без изменений
   - Убедиться что список фаз правильный (5 фаз)

3. **Обновить code review note:**
   - Добавить note что title исправлен
   - Указать что это чисто документационное исправление
   - Упомянуть что functionality неизменна

4. **Верификация:**
   - Проверить что "4-phase" больше нет в WS-201-04.md
   - Проверить что "5-phase" используется корректно
   - Проверить что список фаз соответствует master prompt

---

### Ожидаемый результат

- WS-201-04.md: title и описание обновлены (5-phase workflow)
- Code review note: добавлен комментарий об исправлении
- Функциональность: неизменна (команда /debug работает)

### Scope Estimate

- Файлов: 1 изменено
- Строк: ~5 (TINY)
- Токенов: ~15

---

### Критерий завершения

```bash
# Нет упоминаний "4-phase"
! grep -i "4-phase" tools/hw_checker/docs/workstreams/backlog/WS-201-04-debug-command.md

# Есть упоминания "5-phase"
grep -q "5-phase" tools/hw_checker/docs/workstreams/backlog/WS-201-04-debug-command.md

# Функциональность неизменна (команды существуют)
test -f .cursor/commands/debug.md
test -f .opencode/commands/debug.md
```

---

### Ограничения

- НЕ менять: реализацию команд (`.cursor/commands/debug.md`, `.opencode/commands/debug.md`)
- НЕ трогать: мастер-промпт `sdp/prompts/commands/issue.md`
- НЕ делать: изменений в функциональности
- ТОЛЬКО исправить: документацию (title/описание)

---

### Пример исправления

**Было:**
```markdown
## WS-201-04: /debug command for Cursor and OpenCode

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Команда `/debug` доступна в Cursor и OpenCode
- /debug использует мастер-промпт из `sdp/prompts/commands/debug.md`
- Systematic debugging workflow работает (4-phase: Gather, Analyze, Fix, Verify)
- Failsafe rule соблюден (3 strikes → escalate)

**Acceptance Criteria:**
- [x] `/debug` работает в Cursor (тестовый сценарий)
- [x] /debug работает в OpenCode (тестовый сценарий)
- [x] 4-phase debugging workflow используется
- [x] Failsafe rule (3 strikes) соблюден
```

**Стало:**
```markdown
## WS-201-04: /debug command for Cursor and OpenCode (5-phase workflow)

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Команда `/debug` доступна в Cursor и OpenCode
- /debug использует мастер-промпт из `sdp/prompts/commands/issue.md` Section 4.0
- Systematic debugging workflow работает (5-phase: Symptom → Hypothesis → Test → Root Cause → Impact)
- Failsafe rule соблюден (3 strikes → escalate)

**Acceptance Criteria:**
- [x] `/debug` работает в Cursor (тестовый сценарий)
- [x] /debug работает в OpenCode (тестовый сценарий)
- [x] 5-phase debugging workflow используется
- [x] Failsafe rule (3 strikes) соблюден
```

### Note for Code Review

Добавить в конце WS-201-04.md (после Code Review Results):

```markdown
---

## Post-Review Fix (2026-01-23)

**Issue:** Title mentioned "4-phase debugging workflow" but implementation uses 5 phases

**Fix Applied:**
- Updated title: "/debug command for Cursor and OpenCode (5-phase workflow)"
- Updated goal description to reflect 5-phase workflow
- Updated AC3: "5-phase debugging workflow используется"
- Implementation unchanged (still uses correct 5-phase workflow from issue.md Section 4.0)

**Verification:**
- Title now matches implementation
- No functional changes
- Code review note added for traceability
```
