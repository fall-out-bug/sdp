"""Project map parsing and querying for SDP projects.

This module provides ProjectMap abstraction that contains project-level
decisions, constraints, patterns, and current state information.
"""

import re
from dataclasses import dataclass, field
from pathlib import Path
from typing import Optional


@dataclass
class Decision:
    """Single architectural decision record.

    Attributes:
        area: Decision area (e.g., "Architecture", "Storage")
        decision: Decision description
        adr: ADR identifier (e.g., "ADR-001")
        date: Decision date (YYYY-MM-DD)
    """

    area: str
    decision: str
    adr: str
    date: str


@dataclass
class Constraint:
    """Single project constraint.

    Attributes:
        category: Constraint category (e.g., "AI-Readiness", "Clean Architecture")
        description: Constraint description text
    """

    category: str
    description: str


@dataclass
class TechStackItem:
    """Single tech stack entry.

    Attributes:
        layer: Layer name (e.g., "Language", "API")
        technology: Technology name (e.g., "Python 3.11+", "FastAPI")
    """

    layer: str
    technology: str


@dataclass
class ProjectMap:
    """Parsed project map specification.

    Attributes:
        project_name: Project name extracted from title
        decisions: List of architectural decisions
        constraints: List of active constraints
        current_state: Current state section content (raw markdown)
        patterns: Patterns & conventions section content (raw markdown)
        tech_stack: List of tech stack items
        file_path: Path to source PROJECT_MAP.md file
    """

    project_name: str
    decisions: list[Decision] = field(default_factory=list)
    constraints: list[Constraint] = field(default_factory=list)
    current_state: str = ""
    patterns: str = ""
    tech_stack: list[TechStackItem] = field(default_factory=list)
    file_path: Optional[Path] = None


class ProjectMapParseError(Exception):
    """Error parsing project map file."""

    pass


def parse_project_map(file_path: Path) -> ProjectMap:
    """Parse PROJECT_MAP.md file.

    Args:
        file_path: Path to PROJECT_MAP.md file

    Returns:
        Parsed ProjectMap instance

    Raises:
        FileNotFoundError: If file doesn't exist
        ProjectMapParseError: If file format is invalid
    """
    if not file_path.exists():
        raise FileNotFoundError(f"Project map file not found: {file_path}")

    content = file_path.read_text(encoding="utf-8")

    # Extract project name from title
    project_name = _extract_project_name(content)
    if not project_name:
        raise ProjectMapParseError("Could not extract project name from title")

    # Parse decisions table
    decisions = _parse_decisions_table(content)

    # Parse constraints
    constraints = _parse_constraints(content)

    # Extract current state section
    current_state = _extract_section(content, "Current State")

    # Extract patterns section
    patterns = _extract_section(content, "Patterns & Conventions")

    # Parse tech stack table
    tech_stack = _parse_tech_stack_table(content)

    return ProjectMap(
        project_name=project_name,
        decisions=decisions,
        constraints=constraints,
        current_state=current_state,
        patterns=patterns,
        tech_stack=tech_stack,
        file_path=file_path,
    )


def get_decision(
    project_map: ProjectMap,
    *,
    area: Optional[str] = None,
    adr: Optional[str] = None,
) -> Optional[Decision]:
    """Query decision by area or ADR.

    Args:
        project_map: ProjectMap instance
        area: Decision area to search for
        adr: ADR identifier to search for

    Returns:
        Decision if found, None otherwise

    Raises:
        ValueError: If neither area nor adr is provided
    """
    if area is None and adr is None:
        raise ValueError("Must provide either 'area' or 'adr' parameter")

    for decision in project_map.decisions:
        if area is not None and decision.area == area:
            return decision
        if adr is not None and decision.adr == adr:
            return decision

    return None


def get_constraint(
    project_map: ProjectMap,
    *,
    category: Optional[str] = None,
    keyword: Optional[str] = None,
) -> list[Constraint]:
    """Query constraints by category or keyword.

    Args:
        project_map: ProjectMap instance
        category: Constraint category to filter by
        keyword: Keyword to search in constraint descriptions

    Returns:
        List of matching constraints (empty if none found)

    Raises:
        ValueError: If neither category nor keyword is provided
    """
    if category is None and keyword is None:
        raise ValueError("Must provide either 'category' or 'keyword' parameter")

    results: list[Constraint] = []

    for constraint in project_map.constraints:
        if category is not None and constraint.category == category:
            results.append(constraint)
        elif keyword is not None and keyword.lower() in constraint.description.lower():
            results.append(constraint)

    return results


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


# =============================================================================
# PRIVATE HELPER FUNCTIONS
# =============================================================================


def _extract_project_name(content: str) -> Optional[str]:
    """Extract project name from title line.

    Args:
        content: Markdown content

    Returns:
        Project name or None if not found
    """
    match = re.search(r"^# Project Map:\s*(.+)$", content, re.MULTILINE)
    if match:
        return match.group(1).strip()
    return None


def _parse_decisions_table(content: str) -> list[Decision]:
    """Parse decisions from Quick Reference table.

    Args:
        content: Markdown content

    Returns:
        List of Decision objects
    """
    decisions: list[Decision] = []

    # Find the decisions table section
    table_start = content.find("## Ключевые решения")
    if table_start == -1:
        return decisions

    # Find the table (starts after header row)
    # Note: There may be blank lines between section header and table
    table_match = re.search(
        r"\|\s*Область\s*\|[^\n]*\n\|[-\s|]+\|\s*\n(.*?)(?=\n\s*\n|\n\s*##|\Z)",
        content[table_start:],
        re.DOTALL,
    )

    if not table_match:
        return decisions

    table_rows = table_match.group(1).strip()

    # Parse each row
    for line in table_rows.split("\n"):
        line = line.strip()
        if not line.startswith("|"):
            continue

        # Split by | and remove empty strings
        parts = [p.strip() for p in line.split("|") if p.strip()]
        
        if len(parts) >= 4:
            area = parts[0].replace("**", "").strip()
            decision = parts[1].strip()
            adr_raw = parts[2].strip()
            date = parts[3].strip()

            # Extract ADR from markdown link if present: [ADR-001](path) -> ADR-001
            adr_match = re.search(r"\[(ADR-\d+)\]", adr_raw)
            adr = adr_match.group(1) if adr_match else adr_raw

            # Skip empty rows (check for placeholder text)
            if area and decision and adr and date and not area.startswith("_"):
                decisions.append(Decision(area=area, decision=decision, adr=adr, date=date))

    return decisions


def _parse_constraints(content: str) -> list[Constraint]:
    """Parse constraints from Active Constraints section.

    Args:
        content: Markdown content

    Returns:
        List of Constraint objects
    """
    constraints: list[Constraint] = []

    # Find Active Constraints section
    section_start = content.find("## Active Constraints")
    if section_start == -1:
        return constraints

    # Extract section content until next ## or ---
    section_match = re.search(
        r"## Active Constraints\s*\n(.*?)(?=\n\s*##|\n\s*---|\Z)",
        content,
        re.DOTALL,
    )

    if not section_match:
        return constraints

    section_content = section_match.group(1)

    # Parse subsections (### Category)
    current_category = ""
    for line in section_content.split("\n"):
        line = line.strip()

        # Check for category header
        category_match = re.match(r"^### (.+)$", line)
        if category_match:
            current_category = category_match.group(1).strip()
            continue

        # Check for constraint item (starts with -)
        if line.startswith("-") and current_category:
            # Remove markdown list marker and extract description
            description = re.sub(r"^-\s*", "", line).strip()
            if description:
                constraints.append(Constraint(category=current_category, description=description))

    return constraints


def _extract_section(content: str, section_name: str) -> str:
    """Extract section content by name.

    Args:
        content: Markdown content
        section_name: Section name to extract

    Returns:
        Section content (raw markdown) or empty string
    """
    # Try different header patterns (## and ###)
    # Allow optional trailing content after section name (e.g., "Current State (Live)")
    escaped_name = re.escape(section_name)
    
    # Find section header (## or ###)
    for header_level in ["##", "###"]:
        # Pattern to match header with section name and optional trailing content
        header_pattern = rf"^{re.escape(header_level)}\s+{escaped_name}(?:[^\n]*)?$"
        header_match = re.search(header_pattern, content, re.MULTILINE)
        
        if not header_match:
            continue
        
        # Find the start of content (after the header line and its newline)
        # header_match.end() is at the end of the matched line
        # Skip any newline character to get to the actual content
        content_start = header_match.end()
        # Skip the newline if present
        if content_start < len(content) and content[content_start] == '\n':
            content_start += 1
        
        # Find the next section header at the same or higher level
        # For ##, stop at next ##
        # For ###, stop at next ### or ##
        if header_level == "##":
            next_section_pattern = r"^\s*##\s+"
        else:  # ###
            next_section_pattern = r"^\s*(?:###|##)\s+"
        
        # Find the next section header
        next_match = re.search(next_section_pattern, content[content_start:], re.MULTILINE)
        
        if next_match:
            # Extract content up to the next section (excluding the newline before it)
            content_end = content_start + next_match.start()
            section_content = content[content_start:content_end]
        else:
            # No next section, extract to end of file
            section_content = content[content_start:]
        
        return section_content.strip()
    
    return ""


def _parse_tech_stack_table(content: str) -> list[TechStackItem]:
    """Parse tech stack from table.

    Args:
        content: Markdown content

    Returns:
        List of TechStackItem objects
    """
    items: list[TechStackItem] = []

    # Find Tech Stack section
    section_start = content.find("## Tech Stack")
    if section_start == -1:
        return items

    # Find the table
    table_match = re.search(
        r"\|\s*Layer\s*\|[^\n]*\n\|[-\s|]+\|\s*\n(.*?)(?=\n\s*\n|\n\s*##|\Z)",
        content[section_start:],
        re.DOTALL,
    )

    if not table_match:
        return items

    table_rows = table_match.group(1).strip()

    # Parse each row
    for line in table_rows.split("\n"):
        line = line.strip()
        if not line.startswith("|"):
            continue

        # Split by | and remove empty strings
        parts = [p.strip() for p in line.split("|") if p.strip()]

        if len(parts) >= 2:
            layer = parts[0].strip()
            technology = parts[1].strip()

            # Skip placeholder rows
            if layer and technology and not layer.startswith("_"):
                items.append(TechStackItem(layer=layer, technology=technology))

    return items
