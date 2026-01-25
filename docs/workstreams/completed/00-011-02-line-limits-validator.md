---
ws_id: 00-195-02
project_id: 00
feature: F011
status: backlog
size: SMALL
github_issue: 1030
assignee: null
started: null
completed: null
blocked_reason: null
---

## 02-195-02: Line Limits Validator

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Валидатор проверяет line limits для каждой секции PRD
- Warning/Error режимы в зависимости от enforcement level
- CLI команда `sdp-prd validate {path}` проверяет PRD файл

**Acceptance Criteria:**
- [ ] AC1: Валидатор проверяет max_chars для "Назначение" (≤500)
- [ ] AC2: Валидатор проверяет формат "Модель БД" (1 строка на поле)
- [ ] AC3: Warning выводится при превышении soft limits
- [ ] AC4: Error и exit 1 при превышении hard limits
- [ ] AC5: `sdp-prd validate PROJECT_MAP.md` работает из CLI

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

"ПишиСокращай" — защита от verbose документации. Line limits enforcement гарантирует что PRD остаётся concise.

### Зависимость

00--01 (profiles.py с определением limits)

### Входные файлы

- `sdp/src/sdp/prd/profiles.py` — определение секций с limits
- `tools/hw_checker/docs/drafts/idea-prd-driven-project-maps.md` — спецификация limits

### Шаги

1. Создать `sdp/src/sdp/prd/validator.py` — валидатор секций
2. Создать `sdp/src/sdp/prd/parser.py` — парсер PROJECT_MAP.md
3. Добавить CLI команду в `sdp/src/sdp/cli.py`
4. Написать тесты

### Код

```python
# sdp/src/sdp/prd/validator.py
from dataclasses import dataclass
from enum import Enum
from pathlib import Path

class Severity(Enum):
    WARNING = "warning"
    ERROR = "error"

@dataclass
class ValidationIssue:
    section: str
    message: str
    severity: Severity
    current: int
    limit: int

def validate_prd(content: str) -> list[ValidationIssue]:
    """Validate PRD content against line limits."""
    issues = []
    
    # Parse sections
    sections = parse_prd_sections(content)
    
    # Check "Назначение" (max 500 chars)
    if "Назначение" in sections:
        text = sections["Назначение"]
        if len(text) > 500:
            issues.append(ValidationIssue(
                section="Назначение",
                message=f"Превышен лимит: {len(text)}/500 chars",
                severity=Severity.WARNING,
                current=len(text),
                limit=500
            ))
    
    # Check "Модель БД" (1 line per field)
    if "Модель БД" in sections:
        db_section = sections["Модель БД"]
        for line in db_section.strip().split("\n"):
            # Skip headers and empty lines
            if line.startswith("#") or not line.strip():
                continue
            # Each field should be on single line
            if len(line) > 120:
                issues.append(ValidationIssue(
                    section="Модель БД",
                    message=f"Поле превышает 120 chars: {line[:50]}...",
                    severity=Severity.ERROR,
                    current=len(line),
                    limit=120
                ))
    
    return issues

def parse_prd_sections(content: str) -> dict[str, str]:
    """Parse PRD content into sections."""
    import re
    sections = {}
    current_section = None
    current_content = []
    
    for line in content.split("\n"):
        if match := re.match(r"^## \d+\. (.+)$", line):
            if current_section:
                sections[current_section] = "\n".join(current_content)
            current_section = match.group(1)
            current_content = []
        elif current_section:
            current_content.append(line)
    
    if current_section:
        sections[current_section] = "\n".join(current_content)
    
    return sections
```

### Ожидаемый результат

```
sdp/src/sdp/prd/
├── validator.py      # Validation logic
└── parser.py         # PRD section parser

sdp/tests/unit/prd/
├── test_validator.py
└── test_parser.py
```

### Scope Estimate

- Файлов: ~4 создано/изменено
- Строк кода: ~200
- Токенов: ~1000

**Оценка размера:** SMALL

### Критерий завершения

```bash
# Unit tests
pytest sdp/tests/unit/prd/test_validator.py -v
pytest sdp/tests/unit/prd/test_parser.py -v

# CLI check
cd sdp
poetry run sdp-prd validate ../tools/hw_checker/docs/PROJECT_MAP.md
```

### Ограничения

- НЕ парсить аннотации (00--03)
- НЕ генерировать диаграммы (00--04)
- Только validation, не auto-fix
