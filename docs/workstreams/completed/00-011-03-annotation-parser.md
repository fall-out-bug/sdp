---
ws_id: 00-195-03
project_id: 00
feature: F011
status: backlog
size: MEDIUM
github_issue: 1032
assignee: null
started: null
completed: null
blocked_reason: null
---

## 02-195-03: Annotation Parser

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Парсер извлекает `@prd_flow` и `@prd_step` декораторы из Python
- Парсер извлекает `# @prd:` комментарии из bash/yaml
- Парсер собирает flow steps в структурированный формат

**Acceptance Criteria:**
- [ ] AC1: `parse_python_annotations(path)` возвращает list[FlowStep]
- [ ] AC2: `parse_bash_annotations(path)` возвращает list[FlowStep]
- [ ] AC3: Поддержка multi-file parsing с glob patterns
- [ ] AC4: Вывод FlowStep содержит: flow_name, step_number, description, source_file, line_number
- [ ] AC5: Корректная обработка edge cases (отсутствие аннотаций, malformed syntax)

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

Аннотации в коде — это single source of truth для диаграмм. Парсер извлекает эти аннотации для последующей генерации Mermaid/PlantUML.

### Зависимость

Независимый (может работать параллельно с 00--01)

### Входные файлы

- `tools/hw_checker/docs/drafts/idea-prd-driven-project-maps.md` — формат аннотаций
- Python files с @prd декораторами (примеры)
- Bash/YAML files с # @prd: комментариями (примеры)

### Шаги

1. Создать `sdp/src/sdp/prd/annotations.py` — data classes для FlowStep
2. Создать `sdp/src/sdp/prd/parser_python.py` — Python parser
3. Создать `sdp/src/sdp/prd/parser_bash.py` — Bash/YAML parser
4. Создать `sdp/src/sdp/prd/decorators.py` — @prd_flow, @prd_step декораторы
5. Написать тесты с fixtures

### Код

```python
# sdp/src/sdp/prd/annotations.py
from dataclasses import dataclass
from pathlib import Path

@dataclass
class FlowStep:
    flow_name: str
    step_number: int
    description: str
    source_file: Path
    line_number: int
    participant: str | None = None  # For sequence diagrams

@dataclass
class Flow:
    name: str
    steps: list[FlowStep]
```

```python
# sdp/src/sdp/prd/decorators.py
"""PRD annotation decorators for Python code."""
from functools import wraps
from typing import Callable, TypeVar

F = TypeVar('F', bound=Callable)

def prd_flow(flow_name: str) -> Callable[[F], F]:
    """Mark function as part of a PRD flow."""
    def decorator(func: F) -> F:
        func._prd_flow = flow_name  # type: ignore
        return func
    return decorator

def prd_step(step_number: int, description: str) -> Callable[[F], F]:
    """Mark function as a step in the PRD flow."""
    def decorator(func: F) -> F:
        func._prd_step = step_number  # type: ignore
        func._prd_step_desc = description  # type: ignore
        return func
    return decorator
```

```python
# sdp/src/sdp/prd/parser_python.py
import ast
import re
from pathlib import Path
from .annotations import FlowStep

def parse_python_annotations(path: Path) -> list[FlowStep]:
    """Parse @prd_flow and @prd_step decorators from Python file."""
    content = path.read_text()
    steps = []
    
    # Use regex to find decorator patterns
    # Pattern: @prd_flow("name") or @prd_step(N, "desc")
    flow_pattern = re.compile(
        r'@prd_flow\(["\']([^"\']+)["\']\)\s*\n'
        r'(?:@prd_step\((\d+),\s*["\']([^"\']+)["\']\)\s*\n)?'
        r'(?:async\s+)?def\s+(\w+)',
        re.MULTILINE
    )
    
    for match in flow_pattern.finditer(content):
        flow_name = match.group(1)
        step_num = int(match.group(2)) if match.group(2) else 0
        step_desc = match.group(3) or match.group(4)  # fallback to func name
        line_number = content[:match.start()].count('\n') + 1
        
        steps.append(FlowStep(
            flow_name=flow_name,
            step_number=step_num,
            description=step_desc,
            source_file=path,
            line_number=line_number
        ))
    
    return steps

def parse_directory(directory: Path, pattern: str = "**/*.py") -> list[FlowStep]:
    """Parse all Python files in directory."""
    all_steps = []
    for file in directory.glob(pattern):
        all_steps.extend(parse_python_annotations(file))
    return all_steps
```

```python
# sdp/src/sdp/prd/parser_bash.py
import re
from pathlib import Path
from .annotations import FlowStep

def parse_bash_annotations(path: Path) -> list[FlowStep]:
    """Parse # @prd: comments from bash/yaml files."""
    content = path.read_text()
    steps = []
    
    # Pattern: # @prd: flow=name, step=N, desc=description
    pattern = re.compile(
        r'^#\s*@prd:\s*flow=([^,]+),\s*step=(\d+)(?:,\s*desc=(.+))?$',
        re.MULTILINE
    )
    
    for match in pattern.finditer(content):
        flow_name = match.group(1).strip()
        step_num = int(match.group(2))
        description = match.group(3).strip() if match.group(3) else ""
        line_number = content[:match.start()].count('\n') + 1
        
        steps.append(FlowStep(
            flow_name=flow_name,
            step_number=step_num,
            description=description,
            source_file=path,
            line_number=line_number
        ))
    
    return steps
```

### Ожидаемый результат

```
sdp/src/sdp/prd/
├── annotations.py      # Data classes
├── decorators.py       # @prd_flow, @prd_step
├── parser_python.py    # Python parser
└── parser_bash.py      # Bash/YAML parser

sdp/tests/unit/prd/
├── test_parser_python.py
├── test_parser_bash.py
└── fixtures/
    ├── sample_annotated.py
    └── sample_annotated.sh
```

### Scope Estimate

- Файлов: ~7 создано
- Строк кода: ~350 (annotations: 30, decorators: 30, parsers: 200, tests: 90)
- Токенов: ~1800

**Оценка размера:** MEDIUM

### Критерий завершения

```bash
# Unit tests
pytest sdp/tests/unit/prd/test_parser_python.py -v
pytest sdp/tests/unit/prd/test_parser_bash.py -v

# Coverage ≥ 80%
pytest sdp/tests/unit/prd/test_parser*.py -v \
  --cov=sdp/src/sdp/prd \
  --cov-report=term-missing \
  --cov-fail-under=80

# Import check
python -c "from sdp.prd import parse_python_annotations, parse_bash_annotations, prd_flow, prd_step"
```

### Ограничения

- НЕ генерировать диаграммы (00--04)
- НЕ использовать AST для сложных случаев (regex fallback для v1.0)
- НЕ поддерживать nested functions
