---
id: WS-201-09
title: Fix WS-201-ep30-statefulsets misclassification
feature: F007
status: backlog
size: SMALL
github_issue: TBD
---

## WS-201-09: Fix WS-201-ep30-statefulsets misclassification

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- WS-201-ep30-statefulsets.md перемещен в правильную категорию (F30)
- INDEX.md обновлен (WS-201-ep30-statefulsets в секции F30, не F201)
- Нет путаницы в организации workstreams

**Acceptance Criteria:**
- [ ] WS-201-ep30-statefulsets.md перемещен из `backlog/` в `backlog/` под F30 секцией
- [ ] INDEX.md обновлен (WS-201-ep30-statefulsets в F30, не F201)
- [ ] F201 INDEX section не содержит WS-201-ep30-statefulsets
- [ ] F30 INDEX section содержит WS-201-ep30-statefulsets
- [ ] Правильный frontmatter: `feature: F30` (не F201)

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

**Текущее состояние:**
- WS-201-ep30-statefulsets.md находится в backlog/
- WS incorrectly categorized under F201 (Multi-IDE SDP Parity)
- Frontmatter: `feature: F007` (неправильно!)

**Проблема:**
- WS-201-ep30-statefulsets относится к F30 (Kubernetes Migration)
- WS описывает создание StatefulSets для Nexus, PostgreSQL, Redis в k3s
- Это часть F30 feature, не F201

**Решение:**
- Исправить frontmatter: `feature: F30`
- Переместить или правильно отобразить в INDEX.md (в секции F30)
- Удалить из секции F201 в INDEX.md

---

### Зависимость

Независный (чисто организационное изменение)

---

### Входные файлы

- `tools/hw_checker/docs/workstreams/backlog/WS-201-ep30-statefulsets.md`
- `tools/hw_checker/docs/workstreams/INDEX.md`

---

### Шаги

1. **Исправить frontmatter WS:**
   - Откыть `WS-201-ep30-statefulsets.md`
   - Изменить `feature: F007` → `feature: F30`
   - Сохранить файл

2. **Обновить INDEX.md:**
   - Найти WS-201-ep30-statefulsets в секции F201
   - Удалить из секции F201
   - Добавить в секцию P3: Kubernetes F30
   - Убедиться что формат соответствует другим WS в F30 секции

3. **Верификация:**
   - Проверить что WS-201-ep30-statefulsets в секции F30
   - Проверить что WS-201-ep30-statefulsets НЕ в секции F201
   - Проверить что frontmatter содержит `feature: F30`

---

### Ожидаемый результат

- WS-201-ep30-statefulsets.md: обновлен (feature: F30)
- INDEX.md: WS-201-ep30-statefulsets в секции P3: Kubernetes F30
- INDEX.md: WS-201-ep30-statefulsets НЕ в секции P0: F201 Multi-IDE SDP Parity

### Scope Estimate

- Файлов: 2 изменено (WS + INDEX)
- Строк: ~30 (SMALL)
- Токенов: ~90

---

### Критерий завершения

```bash
# Frontmatter исправлен
grep -q "feature: F30" tools/hw_checker/docs/workstreams/backlog/WS-201-ep30-statefulsets.md
! grep -q "feature: F007" tools/hw_checker/docs/workstreams/backlog/WS-201-ep30-statefulsets.md

# INDEX.md: в F30 секции
grep -A5 "### P3: Kubernetes F30" tools/hw_checker/docs/workstreams/INDEX.md | grep -q "WS-201"

# INDEX.md: НЕ в F201 секции
! grep -A20 "### P0: F201 Multi-IDE SDP Parity" tools/hw_checker/docs/workstreams/INDEX.md | grep -q "WS-201-ep30"
```

---

### Ограничения

- НЕ менять: контент WS (только frontmatter и классификацию)
- НЕ трогать: другие WS
- НЕ делать: изменений в кодовой базе (только документация)

---

### Пример исправления frontmatter

**Было:**
```markdown
---
ws_id: 00-201
project_id: 00
feature: F007
status: backlog
size: MEDIUM
github_issue: 231
assignee: null
started: null
completed: null
blocked_reason: null
---
```

**Стало:**
```markdown
---
ws_id: 00-201
feature: F30
status: backlog
size: MEDIUM
github_issue: 231
assignee: null
started: null
completed: null
blocked_reason: null
---
```

### Пример изменения в INDEX.md

**Удалить из F201:**
```markdown
### P0: F201 Multi-IDE SDP Parity

| ID | Title | Зависимость | Файл |
|----|-------|-------------|------|
| WS-201-01 | Validate /oneshot in Cursor and OpenCode | — | [→](backlog/WS-201-01-validate-oneshot.md) |
| WS-201-02 | Cross-platform Git hooks for SDP | — | [→](backlog/WS-201-02-git-hooks.md) |
| WS-201-03 | Cursor agents parity + OpenCode integration | — | [→](backlog/WS-201-03-cursor-agents.md) |
| WS-201-04 | /debug command for Cursor and OpenCode | — | [→](backlog/WS-201-04-debug-command.md) |
| WS-201-05 | /test command for Cursor and OpenCode (after F194) | WS-410-01..05 | [→](backlog/WS-201-05-test-command.md) |
| WS-201-06 | Documentation & runbooks for multi-ide parity | WS-201-01..05 | [→](backlog/WS-201-06-documentation.md) |
| WS-201-ep30-statefulsets | StatefulSets (Nexus, PG, Redis) | WS-200 | [→](backlog/WS-201-ep30-statefulsets.md) | ← УДАЛИТЬ ЭТУ СТРОКУ
```

**Добавить в F30:**
```markdown
### P3: Kubernetes F30

| ID | Title | Source | Зависимость | Файл |
|----|-------|--------|-------------|------|
| WS-200 | k3s Cluster Setup | F30 | — | [→](backlog/WS-200-ep30-k3s-setup.md) |
| WS-201 | StatefulSets (Nexus, PG, Redis) | F30 | WS-200 | [→](backlog/WS-201-ep30-statefulsets.md) | ← ДОБАВИТЬ ЭТУ СТРОКУ
| WS-202 | Kubernetes Executor Adapter | F30 | WS-201 | [→](backlog/WS-202-ep30-k8s-executor.md) |
| WS-203 | Network Policy | F30 | WS-201 | [→](backlog/WS-203-ep30-network-policy.md) |
```
