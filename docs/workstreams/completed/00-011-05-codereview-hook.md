---
ws_id: 00-195-05
project_id: 00
feature: F011
status: backlog
size: SMALL
github_issue: 1036
assignee: null
started: null
completed: null
blocked_reason: null
---

## 02-195-05: Codereview Hook Integration

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- `/codereview` проверяет актуальность диаграмм через hash comparison
- При расхождении hash выводит CHANGES REQUESTED с инструкцией
- Hash хранится в frontmatter PROJECT_MAP.md (`diagrams_hash: abc123`)

**Acceptance Criteria:**
- [ ] AC1: `calculate_diagrams_hash(path)` возвращает SHA256 hash всех аннотаций
- [ ] AC2: post-codereview.sh содержит Check N: PRD Diagrams
- [ ] AC3: При mismatch выводит "❌ Диаграммы устарели. Run: /prd {project} --update"
- [ ] AC4: `--skip-prd` флаг пропускает проверку
- [ ] AC5: Если PROJECT_MAP.md отсутствует — skip с warning

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

Hook гарантирует что диаграммы синхронизированы с кодом. Это финальный enforcement — review не проходит если диаграммы устарели.

### Зависимость

00--04 (diagram generator, для hash calculation)

### Входные файлы

- `sdp/hooks/post-codereview.sh` — существующий hook
- `sdp/src/sdp/prd/parser_python.py` — для извлечения аннотаций

### Шаги

1. Создать `sdp/src/sdp/prd/hash.py` — hash calculation
2. Добавить Check в `post-codereview.sh`
3. Обновить `codereview.md` — документация проверки
4. Написать тесты

### Код

```python
# sdp/src/sdp/prd/hash.py
import hashlib
from pathlib import Path
from .parser_python import parse_directory
from .parser_bash import parse_bash_annotations

def calculate_diagrams_hash(project_path: Path) -> str:
    """Calculate SHA256 hash of all @prd annotations in project."""
    all_steps = []
    
    # Collect Python annotations
    for py_file in project_path.rglob("*.py"):
        from .parser_python import parse_python_annotations
        all_steps.extend(parse_python_annotations(py_file))
    
    # Collect bash/yaml annotations
    for ext in ["*.sh", "*.yml", "*.yaml"]:
        for file in project_path.rglob(ext):
            all_steps.extend(parse_bash_annotations(file))
    
    # Sort for deterministic hash
    sorted_steps = sorted(all_steps, key=lambda s: (s.flow_name, s.step_number, str(s.source_file)))
    
    # Create hash from normalized content
    content = "\n".join(
        f"{s.flow_name}|{s.step_number}|{s.description}"
        for s in sorted_steps
    )
    
    return hashlib.sha256(content.encode()).hexdigest()[:12]

def get_stored_hash(project_map_path: Path) -> str | None:
    """Extract diagrams_hash from PROJECT_MAP.md frontmatter."""
    import re
    if not project_map_path.exists():
        return None
    
    content = project_map_path.read_text()
    match = re.search(r'^diagrams_hash:\s*(\w+)', content, re.MULTILINE)
    return match.group(1) if match else None

def update_stored_hash(project_map_path: Path, new_hash: str) -> None:
    """Update diagrams_hash in PROJECT_MAP.md frontmatter."""
    import re
    content = project_map_path.read_text()
    
    if "diagrams_hash:" in content:
        content = re.sub(
            r'^diagrams_hash:\s*\w*',
            f'diagrams_hash: {new_hash}',
            content,
            flags=re.MULTILINE
        )
    else:
        # Add after prd_version line
        content = re.sub(
            r'^(prd_version:\s*".+")',
            f'\\1\ndiagrams_hash: {new_hash}',
            content,
            flags=re.MULTILINE
        )
    
    project_map_path.write_text(content)
```

```bash
# Addition to post-codereview.sh

echo ""
echo "Check N: PRD Diagrams Актуальны"

# Skip if --skip-prd flag
if [[ " $* " == *" --skip-prd "* ]]; then
    echo "  ⚠️ Skipped (--skip-prd)"
else
    PROJECT_MAP="tools/hw_checker/docs/PROJECT_MAP.md"
    
    if [ ! -f "$PROJECT_MAP" ]; then
        echo "  ⚠️ PROJECT_MAP.md not found, skipping PRD check"
    else
        # Get stored hash
        STORED_HASH=$(grep "^diagrams_hash:" "$PROJECT_MAP" | cut -d: -f2 | tr -d ' ')
        
        if [ -z "$STORED_HASH" ]; then
            echo "  ⚠️ No diagrams_hash in PROJECT_MAP.md, skipping"
        else
            # Calculate current hash
            cd sdp
            CURRENT_HASH=$(poetry run python -c "
from sdp.prd.hash import calculate_diagrams_hash
from pathlib import Path
print(calculate_diagrams_hash(Path('../tools/hw_checker')))
")
            cd ..
            
            if [ "$CURRENT_HASH" != "$STORED_HASH" ]; then
                echo "❌ Диаграммы устарели"
                echo "   Stored:  $STORED_HASH"
                echo "   Current: $CURRENT_HASH"
                echo "   Run: /prd hw-checker --update"
                exit 1
            else
                echo "✓ PRD diagrams up-to-date (hash: $STORED_HASH)"
            fi
        fi
    fi
fi
```

### Ожидаемый результат

```
sdp/src/sdp/prd/
└── hash.py               # Hash calculation + storage

sdp/hooks/post-codereview.sh  # Updated with PRD check

sdp/tests/unit/prd/
└── test_hash.py
```

### Scope Estimate

- Файлов: ~3 создано/изменено
- Строк кода: ~150 (hash: 80, hook: 30, tests: 40)
- Токенов: ~800

**Оценка размера:** SMALL

### Критерий завершения

```bash
# Unit tests
pytest sdp/tests/unit/prd/test_hash.py -v

# Integration test
cd sdp && poetry run python -c "
from sdp.prd.hash import calculate_diagrams_hash
from pathlib import Path
h = calculate_diagrams_hash(Path('../tools/hw_checker'))
print(f'Hash: {h}')
"

# Hook test (manual)
./sdp/hooks/post-codereview.sh F195
```

### Ограничения

- НЕ автоматически обновлять диаграммы (только detection)
- НЕ блокировать review если hash отсутствует (warning only)
- НЕ парсить полный AST (regex fallback)
