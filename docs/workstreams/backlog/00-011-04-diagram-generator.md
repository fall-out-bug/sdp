---
ws_id: 00-195-04
project_id: 00
feature: F011
status: backlog
size: MEDIUM
github_issue: 1033
assignee: null
started: null
completed: null
blocked_reason: null
---

## 02-195-04: Diagram Generator

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Генератор создаёт Mermaid sequence diagrams из FlowSteps
- Генератор создаёт PlantUML sequence diagrams из FlowSteps
- Component и deployment diagram templates
- Хранение диаграмм как код (.mmd, .puml)

**Acceptance Criteria:**
- [ ] AC1: `generate_mermaid_sequence(flow)` возвращает валидный Mermaid код
- [ ] AC2: `generate_plantuml_sequence(flow)` возвращает валидный PlantUML код
- [ ] AC3: Диаграммы сохраняются в `docs/diagrams/` как .mmd/.puml файлы
- [ ] AC4: Component diagram template генерируется для service profile
- [ ] AC5: Deployment diagram template генерируется для service profile

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

Диаграммы генерируются из аннотаций в коде. Это обеспечивает их актуальность — single source of truth.

### Зависимость

00--03 (annotation parser)

### Входные файлы

- `sdp/src/sdp/prd/annotations.py` — FlowStep data class
- `tools/hw_checker/docs/drafts/idea-prd-driven-project-maps.md` — примеры диаграмм

### Шаги

1. Создать `sdp/src/sdp/prd/generator_mermaid.py` — Mermaid generator
2. Создать `sdp/src/sdp/prd/generator_plantuml.py` — PlantUML generator
3. Создать `sdp/src/sdp/prd/generator.py` — unified interface
4. Создать templates для component/deployment diagrams
5. Добавить CLI команду для генерации
6. Написать тесты

### Код

```python
# sdp/src/sdp/prd/generator_mermaid.py
from .annotations import Flow, FlowStep

def generate_mermaid_sequence(flow: Flow) -> str:
    """Generate Mermaid sequence diagram from flow steps."""
    lines = [
        "sequenceDiagram",
    ]
    
    # Collect participants
    participants = set()
    for step in flow.steps:
        if step.participant:
            participants.add(step.participant)
    
    for p in sorted(participants):
        lines.append(f"    participant {p}")
    
    lines.append("")
    
    # Generate sequence
    sorted_steps = sorted(flow.steps, key=lambda s: s.step_number)
    for step in sorted_steps:
        # Format: participant->>other: description
        desc = step.description
        source = step.source_file.stem
        lines.append(f"    Note over {source}: Step {step.step_number}")
        lines.append(f"    Note over {source}: {desc}")
    
    return "\n".join(lines)

def generate_mermaid_component() -> str:
    """Generate component diagram template."""
    return """flowchart TB
    subgraph Presentation
        API[FastAPI]
        CLI[Click CLI]
    end
    
    subgraph Application
        UseCase[Use Cases]
        Ports[Ports/Interfaces]
    end
    
    subgraph Domain
        Entities[Entities]
        Services[Domain Services]
    end
    
    subgraph Infrastructure
        DB[(PostgreSQL)]
        Queue[(Redis)]
        External[External APIs]
    end
    
    API --> UseCase
    CLI --> UseCase
    UseCase --> Entities
    UseCase --> Ports
    Ports --> DB
    Ports --> Queue
    Ports --> External
"""
```

```python
# sdp/src/sdp/prd/generator_plantuml.py
from .annotations import Flow

def generate_plantuml_sequence(flow: Flow) -> str:
    """Generate PlantUML sequence diagram from flow steps."""
    lines = [
        "@startuml",
        f"title {flow.name}",
        "",
    ]
    
    sorted_steps = sorted(flow.steps, key=lambda s: s.step_number)
    for step in sorted_steps:
        source = step.source_file.stem
        lines.append(f"note over {source}: Step {step.step_number}: {step.description}")
    
    lines.append("")
    lines.append("@enduml")
    return "\n".join(lines)

def generate_plantuml_deployment() -> str:
    """Generate deployment diagram template."""
    return """@startuml
!include <C4/C4_Deployment>

title Deployment Diagram

Deployment_Node(docker, "Docker Compose", "docker-compose.yml") {
    Container(api, "API", "FastAPI", "REST endpoints")
    Container(worker, "Worker", "Python", "Job processing")
    ContainerDb(pg, "PostgreSQL", "Database")
    ContainerDb(redis, "Redis", "Job queue")
}

Rel(api, pg, "reads/writes")
Rel(api, redis, "enqueue")
Rel(worker, redis, "dequeue")
Rel(worker, pg, "writes")

@enduml
"""
```

```python
# sdp/src/sdp/prd/generator.py
from pathlib import Path
from .annotations import Flow
from .generator_mermaid import generate_mermaid_sequence, generate_mermaid_component
from .generator_plantuml import generate_plantuml_sequence, generate_plantuml_deployment

def generate_diagrams(flows: list[Flow], output_dir: Path) -> list[Path]:
    """Generate all diagrams for flows and save to output directory."""
    output_dir.mkdir(parents=True, exist_ok=True)
    created_files = []
    
    for flow in flows:
        # Mermaid sequence
        mmd_path = output_dir / f"sequence-{flow.name}.mmd"
        mmd_path.write_text(generate_mermaid_sequence(flow))
        created_files.append(mmd_path)
        
        # PlantUML sequence
        puml_path = output_dir / f"sequence-{flow.name}.puml"
        puml_path.write_text(generate_plantuml_sequence(flow))
        created_files.append(puml_path)
    
    # Component diagram (template)
    comp_path = output_dir / "component-overview.mmd"
    comp_path.write_text(generate_mermaid_component())
    created_files.append(comp_path)
    
    # Deployment diagram (template)
    deploy_path = output_dir / "deployment-production.puml"
    deploy_path.write_text(generate_plantuml_deployment())
    created_files.append(deploy_path)
    
    return created_files
```

### Ожидаемый результат

```
sdp/src/sdp/prd/
├── generator.py           # Unified interface
├── generator_mermaid.py   # Mermaid generation
└── generator_plantuml.py  # PlantUML generation

sdp/tests/unit/prd/
├── test_generator_mermaid.py
├── test_generator_plantuml.py
└── test_generator.py

# Example output:
docs/diagrams/
├── sequence-message-processing.mmd
├── sequence-message-processing.puml
├── component-overview.mmd
└── deployment-production.puml
```

### Scope Estimate

- Файлов: ~6 создано
- Строк кода: ~400 (mermaid: 150, plantuml: 150, unified: 50, tests: 50)
- Токенов: ~2000

**Оценка размера:** MEDIUM

### Критерий завершения

```bash
# Unit tests
pytest sdp/tests/unit/prd/test_generator*.py -v

# Coverage ≥ 80%
pytest sdp/tests/unit/prd/test_generator*.py -v \
  --cov=sdp/src/sdp/prd \
  --cov-report=term-missing \
  --cov-fail-under=80

# Manual validation
poetry run python -c "
from sdp.prd.generator_mermaid import generate_mermaid_component
print(generate_mermaid_component())
"
# Должен вывести валидный Mermaid
```

### Ограничения

- НЕ рендерить диаграммы в PNG/SVG (только код)
- НЕ интегрировать PlantUML server (только генерация кода)
- НЕ интегрировать в /codereview (00--05)
