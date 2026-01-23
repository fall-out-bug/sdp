# Release v{X.Y.Z}

**Date:** {YYYY-MM-DD}
**Feature:** {Feature ID} - {Feature Name}

---

## Overview

{Краткое описание что добавлено в этом релизе — 2-3 предложения}

---

## New Features

### {Feature Name}

{Описание функциональности для пользователей}

**Что нового:**
- {Пункт 1}
- {Пункт 2}
- {Пункт 3}

**Использование:**

```bash
# Пример команды или использования
hwc {command} {args}
```

**API (если применимо):**

```bash
# Пример API запроса
curl -X POST http://localhost:8000/api/endpoint \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"key": "value"}'
```

---

## Improvements

- {Улучшение 1}
- {Улучшение 2}

---

## Bug Fixes

- {Исправление 1}
- {Исправление 2}

---

## Breaking Changes

{Если нет breaking changes, написать "None"}

### {Breaking Change 1}

**Было:**
```python
# Старый способ
old_function(arg1, arg2)
```

**Стало:**
```python
# Новый способ
new_function(arg1, arg2, arg3)
```

**Миграция:**
1. Замените `old_function` на `new_function`
2. Добавьте третий аргумент

---

## Migration Guide

{Если не требуется миграция, написать "No migration required"}

### Database Migrations

```bash
# Запуск миграций
cd tools/hw_checker
alembic upgrade head
```

### Configuration Changes

{Если изменился формат конфигурации}

```yaml
# Было
old_config: value

# Стало
new_config:
  nested: value
```

---

## Known Issues

{Если нет known issues, написать "None"}

- {Issue 1}: {описание} — workaround: {как обойти}
- {Issue 2}: {описание}

---

## Dependencies

### Updated
- {Library 1}: v{old} → v{new}
- {Library 2}: v{old} → v{new}

### Added
- {New library}: v{version} — {для чего}

### Removed
- {Removed library} — {почему удалили}

---

## Contributors

- {Contributor 1}
- {Contributor 2}

---

## Full Changelog

See [CHANGELOG.md](../CHANGELOG.md) for full history.

**Workstreams in this release:**
- WS-{ID1}: {title}
- WS-{ID2}: {title}
- WS-{ID3}: {title}
