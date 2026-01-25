---
ws_id: 00-195-01
project_id: 00
feature: F011
status: backlog
size: MEDIUM
github_issue: 1028
assignee: null
started: null
completed: null
blocked_reason: null
---

## 02-195-01: PRD Command + Profiles

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Команда `/prd {project}` создаёт PROJECT_MAP.md v2.0 через интерактивный диалог
- Команда `/prd {project} --update` обновляет существующий PRD
- Автоопределение профиля проекта (service/library/cli)
- Скаффолдинг 7 секций PRD по выбранному профилю

**Acceptance Criteria:**
- [ ] AC1: `/prd hw-checker` запускает scaffold-диалог и создаёт PROJECT_MAP.md
- [ ] AC2: Auto-detect определяет профиль по docker-compose.yml/pyproject.toml/cli.py
- [ ] AC3: Все 3 профиля (service, library, cli) имеют корректные шаблоны секций
- [ ] AC4: `--update` флаг обновляет диаграммы без перезаписи ручных правок
- [ ] AC5: Frontmatter с project_type, prd_version, last_updated создаётся корректно

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

F195 расширяет PROJECT_MAP.md в полноценный PRD с 7 стандартными секциями.
Этот WS создаёт фундамент: команду `/prd` и систему профилей.

### Зависимость

Независимый

### Входные файлы

- `sdp/prompts/commands/design.md` — пример структуры команды
- `tools/hw_checker/docs/PROJECT_MAP.md` — текущий формат
- `tools/hw_checker/docs/drafts/idea-prd-driven-project-maps.md` — спецификация

### Шаги

1. Создать `sdp/prompts/commands/prd.md` — основной промпт команды
2. Создать `sdp/src/sdp/prd/__init__.py` — модуль PRD
3. Создать `sdp/src/sdp/prd/profiles.py` — 3 профиля (service/library/cli)
4. Создать `sdp/src/sdp/prd/detector.py` — автодетект профиля
5. Создать `sdp/src/sdp/prd/scaffold.py` — генератор шаблона
6. Создать `.claude/skills/prd/SKILL.md` — skill для Claude Code
7. Написать тесты для detector и scaffold

### Код

```python
# sdp/src/sdp/prd/profiles.py
from dataclasses import dataclass
from enum import Enum

class ProjectType(Enum):
    SERVICE = "service"
    LIBRARY = "library"
    CLI = "cli"

@dataclass
class PRDSection:
    name: str
    required: bool
    template: str
    max_chars: int | None = None

@dataclass
class PRDProfile:
    project_type: ProjectType
    sections: list[PRDSection]

# Service profile: all 7 sections
SERVICE_SECTIONS = [
    PRDSection("Назначение", True, "...", 500),
    PRDSection("Глоссарий", True, "...", None),
    PRDSection("Внешний API", True, "...", None),
    PRDSection("Модель БД", True, "...", None),
    PRDSection("Sequence Flows", True, "...", None),
    PRDSection("Внешние зависимости", True, "...", None),
    PRDSection("Мониторинги", True, "...", None),
]

# Library profile: no DB, no monitoring
LIBRARY_SECTIONS = [
    PRDSection("Назначение", True, "...", 500),
    PRDSection("Глоссарий", True, "...", None),
    PRDSection("Public API", True, "...", None),
    PRDSection("Data Structures", True, "...", None),
    PRDSection("Usage Examples", True, "...", None),
    PRDSection("Внешние зависимости", True, "...", None),
    PRDSection("Error Handling", True, "...", None),
]

# CLI profile: command reference instead of API
CLI_SECTIONS = [
    PRDSection("Назначение", True, "...", 500),
    PRDSection("Глоссарий", True, "...", None),
    PRDSection("Command Reference", True, "...", None),
    PRDSection("Configuration", True, "...", None),
    PRDSection("Usage Examples", True, "...", None),
    PRDSection("Exit Codes", True, "...", None),
    PRDSection("Error Handling", True, "...", None),
]

PROFILES = {
    ProjectType.SERVICE: PRDProfile(ProjectType.SERVICE, SERVICE_SECTIONS),
    ProjectType.LIBRARY: PRDProfile(ProjectType.LIBRARY, LIBRARY_SECTIONS),
    ProjectType.CLI: PRDProfile(ProjectType.CLI, CLI_SECTIONS),
}
```

```python
# sdp/src/sdp/prd/detector.py
from pathlib import Path
from .profiles import ProjectType

def detect_project_type(project_path: Path) -> ProjectType:
    """Auto-detect project type from file structure."""
    # Check for docker-compose.yml + API endpoints → service
    if (project_path / "docker-compose.yml").exists():
        return ProjectType.SERVICE
    
    # Check for cli.py with Click/Typer → cli
    cli_files = list(project_path.glob("**/cli.py"))
    for cli_file in cli_files:
        content = cli_file.read_text()
        if "click" in content.lower() or "typer" in content.lower():
            return ProjectType.CLI
    
    # Default: library
    return ProjectType.LIBRARY
```

### Ожидаемый результат

```
sdp/
├── prompts/commands/prd.md        # Command prompt
├── src/sdp/prd/
│   ├── __init__.py
│   ├── profiles.py                # 3 profiles
│   ├── detector.py                # Auto-detect
│   └── scaffold.py                # Template generator
└── tests/unit/prd/
    ├── __init__.py
    ├── test_profiles.py
    ├── test_detector.py
    └── test_scaffold.py

.claude/skills/prd/SKILL.md        # Claude Code skill
```

### Scope Estimate

- Файлов: ~8 создано
- Строк кода: ~400 (profiles: 100, detector: 50, scaffold: 150, prompt: 100)
- Токенов: ~2000

**Оценка размера:** MEDIUM

### Критерий завершения

```bash
# Unit tests
pytest sdp/tests/unit/prd/ -v

# Coverage ≥ 80%
pytest sdp/tests/unit/prd/ -v \
  --cov=sdp/src/sdp/prd \
  --cov-report=term-missing \
  --cov-fail-under=80

# Import check
python -c "from sdp.prd import detect_project_type, PROFILES"
```

### Ограничения

- НЕ реализовывать парсинг аннотаций (00--03)
- НЕ реализовывать генерацию диаграмм (00--04)
- НЕ интегрировать в /codereview (00--05)

---

### Human Verification (UAT)

#### 🚀 Quick Smoke Test (30 секунд)

```bash
# Проверить автодетект для hw-checker
cd sdp
poetry run python -c "
from sdp.prd.detector import detect_project_type
from pathlib import Path
result = detect_project_type(Path('../tools/hw_checker'))
print(f'hw-checker detected as: {result.value}')
"
# Ожидаемый результат: hw-checker detected as: service
```

#### 📋 Manual Test Scenarios

| # | Сценарий | Шаги | Ожидание | ✅/❌ |
|---|----------|------|----------|------|
| 1 | Service detection | Проект с docker-compose.yml | service |  |
| 2 | CLI detection | Проект с cli.py + Click | cli |  |
| 3 | Library fallback | Проект без docker/cli | library |  |
