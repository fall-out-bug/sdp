---
ws_id: 00-410-01
project_id: 00
feature: F008
status: backlog
size: SMALL
github_issue: 806
assignee: null
started: null
completed: null
blocked_reason: null
---

## 02-410-01: Contract-Driven WS v2 spec + template

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Введена спецификация Contract-Driven WS v2.0 на основе idea draft (capability tiers, contract sections, verification).
- Шаблон WS обновлён так, чтобы поддерживать `capability_tier` и секции Contract/Verification без нарушения текущего SDP формата.

**Acceptance Criteria:**
- [x] AC1: Создана спецификация (docs) с описанием Contract-Driven WS v2.0 и capability tiers T0–T3.
- [x] AC2: `tools/hw_checker/docs/workstreams/TEMPLATE.md` обновлён: добавлен `capability_tier` и секции Contract/Verification.
- [x] AC3: В спецификации явно указан запрет на модификацию Interface/Tests для T2/T3.

---

### Контекст

Нужно закрепить модель-агностичный формат WS как стабильный capability-tier протокол и встроить его в существующие SDP шаблоны без ломки парсера (`sdp/src/sdp/core/workstream.py`).

### Зависимость

Независимый.

### Входные файлы

- `docs/drafts/idea-model-agnostic-ws-protocol.md` — исходный draft с требованиями.
- `tools/hw_checker/docs/workstreams/TEMPLATE.md` — текущий шаблон WS.
- `sdp/PROTOCOL.md` — правила SDP (формат, gate’ы).

### Шаги

1. Оформить спецификацию Contract-Driven WS v2.0 как отдельный документ (в `tools/hw_checker/docs/specs/feature_194/feature.md`).
2. Обновить `TEMPLATE.md`: добавить `capability_tier`, секции Contract/Verification и ограничения для T2/T3.
3. Добавить явные правила “contract is read-only” для T2/T3 в спецификации.

### Код

```markdown
# Нет кода. Только документация и шаблоны.
```

### Ожидаемый результат

- Новый документ `tools/hw_checker/docs/specs/feature_194/feature.md`.
- Обновлённый `tools/hw_checker/docs/workstreams/TEMPLATE.md`.

### Scope Estimate

- Файлов: ~1 создано + ~1 изменено
- Строк: ~150-250 (SMALL)
- Токенов: ~800-1200

### Критерий завершения

```bash
# Doc lint (если используется)
python sdp/scripts/validate.py

# Regression
pytest tools/hw_checker/tests/unit/ -m fast -v
```

### Ограничения

- НЕ менять обязательные поля frontmatter: `ws_id`, `feature`, `status`, `size`.
- НЕ добавлять новые требования к существующему парсеру WS.

---

### Execution Report

**Изменённые файлы:**
- `tools/hw_checker/docs/specs/feature_194/feature.md` — создан, ~270 строк (спецификация Contract-Driven WS v2.0)
- `tools/hw_checker/docs/workstreams/TEMPLATE.md` — изменён, добавлены секции Contract/Verification и поле `capability_tier`

**Выполненные шаги:**
- [x] Шаг 1: Создана спецификация Contract-Driven WS v2.0 в `tools/hw_checker/docs/specs/feature_194/feature.md`
- [x] Шаг 2: Обновлён `TEMPLATE.md` с добавлением `capability_tier` и секций Contract/Verification
- [x] Шаг 3: Добавлены явные правила "contract is read-only" для T2/T3 в спецификации

**Критерии завершения:**
```bash
# Проверка структуры файлов
$ ls -la tools/hw_checker/docs/specs/feature_194/
feature.md

$ grep -c "capability_tier" tools/hw_checker/docs/workstreams/TEMPLATE.md
1

$ grep -c "Contract" tools/hw_checker/docs/workstreams/TEMPLATE.md
3

$ grep -c "Read-Only\|read-only\|DO NOT MODIFY" tools/hw_checker/docs/specs/feature_194/feature.md
9
```

**Проблемы:**
- Нет

**Acceptance Criteria Status:**
- ✅ AC1: Спецификация создана с полным описанием Contract-Driven WS v2.0, capability tiers T0–T3, структурой WS, правилами read-only для T2/T3
- ✅ AC2: TEMPLATE.md обновлён: добавлено поле `capability_tier` в frontmatter, секции Contract (с Interface и Tests), Verification, и ограничения для T2/T3
- ✅ AC3: В спецификации явно указан раздел "Read-Only Contract Rules for T2/T3" с критическими правилами о запрете модификации Interface/Tests секций
