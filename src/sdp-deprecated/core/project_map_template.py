"""Project map template generation for SDP projects.

This module provides template generation functionality for PROJECT_MAP.md files.
"""

from pathlib import Path


def create_project_map_template(file_path: Path, project_name: str) -> None:
    """Generate PROJECT_MAP.md template for new project.

    Args:
        file_path: Path where template will be created
        project_name: Project name to use in template
    """
    template = f"""# Project Map: {project_name}

**Назначение:** Сквозная карта решений по проекту. Читается агентами перед началом WS.

**Обновление:** После каждого значимого WS — добавить запись о принятых решениях.

---

## Как использовать

### Для агентов (перед началом WS):

1. **ОБЯЗАТЕЛЬНО прочитай** перед планированием/выполнением WS
2. Проверь что твой WS не противоречит существующим решениям
3. Если нужно принять новое архитектурное решение → создай ADR

### Структура:

```
PROJECT_MAP.md (этот файл)
    ↓ ссылки на
docs/architecture/decisions/ (ADR — Architecture Decision Records)
```

---

## Ключевые решения (Quick Reference)

| Область | Решение | ADR | Дата |
|---------|---------|-----|------|
| _(Добавь решения здесь)_ | | | |

---

## Current State (Live)

### Production Services
- _(Опиши текущие сервисы)_

### Domains (L2)
- _(Опиши домены)_

---

## Patterns & Conventions

### Naming
- Files: `snake_case.py`
- Classes: `PascalCase`
- Functions: `snake_case`
- Constants: `UPPER_SNAKE_CASE`

### Import Order
1. stdlib
2. third-party
3. local imports
4. Relative imports

### Type Hints
- Python 3.10+ syntax: `list[str]`, `dict[str, int]`, `str | None`
- Всегда `-> None` для void functions

### Testing
- Unit tests: `tests/unit/` (marker: `@pytest.mark.fast`)
- Integration: `tests/integration/`
- Coverage ≥ 80% для нового кода

---

## Active Constraints

### AI-Readiness
- Files < 200 LOC
- Complexity < 10 (CC)
- No nesting > 3 levels

### Clean Architecture
- Domain: NO imports from app/infra/presentation
- Application: NO direct infrastructure imports (use ports)
- Infrastructure: implements ports
- Presentation: uses application layer only

### Security
- _(Добавь security constraints)_

---

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | _(укажи язык)_ |
| Framework | _(укажи framework)_ |
| Database | _(укажи БД)_ |
| Testing | _(укажи testing tools)_ |

---

## How to Add Decision

Когда принимаешь значимое решение в WS:

1. **Создай ADR** в `docs/architecture/decisions/`:
   ```bash
   YYYY-MM-DD-{{short-title}}.md
   ```

2. **Формат:** (см. `sdp/PROTOCOL.md` → ADR Template)
   - Status: Proposed / Accepted / Deprecated
   - Context
   - Decision
   - Alternatives Considered
   - Consequences

3. **Обнови PROJECT_MAP.md** (этот файл):
   - Добавь строку в "Ключевые решения"
   - Обнови "Current State" если меняется архитектура

4. **Укажи в WS Review:**
   - "ADR-XXX created: {{short description}}"

---

## Deprecated / Superseded

_(Решения, от которых отказались)_

| Было | Почему отказались | Когда | ADR |
|------|-------------------|-------|-----|
| _(Добавь deprecated решения)_ | | | |

---

## Roadmap Alignment

**Current Release:** _(укажи release)_
**Current Features in progress:**
- _(список features)_

---

**Last updated:** _(YYYY-MM-DD)_
**Maintained by:** _(team/agent name)_
"""

    file_path.write_text(template, encoding="utf-8")
