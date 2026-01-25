"""PRD profiles for different project types.

This module defines the structure and templates for PRD documents
based on project type (service, library, cli).
"""

from dataclasses import dataclass
from enum import Enum


class ProjectType(Enum):
    """Project type enumeration."""
    SERVICE = "service"
    LIBRARY = "library"
    CLI = "cli"


@dataclass
class PRDSection:
    """A section in the PRD document.

    Attributes:
        name: Section name
        required: Whether this section is required
        template: Template content for the section
        max_chars: Maximum character limit (None for unlimited)
    """
    name: str
    required: bool
    template: str
    max_chars: int | None = None


@dataclass
class PRDProfile:
    """Profile for a specific project type.

    Attributes:
        project_type: The project type
        sections: List of sections for this profile
    """
    project_type: ProjectType
    sections: list[PRDSection]


# Service profile: all 7 sections
SERVICE_SECTIONS = [
    PRDSection(
        "Назначение",
        True,
        "## 1. Назначение\n\n"
        "{Краткое описание проекта 1-2 абзаца}\n\n"
        "**Ключевые возможности:**\n"
        "- {Возможность 1}\n"
        "- {Возможность 2}\n"
        "- {Возможность 3}\n",
        500,
    ),
    PRDSection(
        "Глоссарий",
        True,
        "## 2. Глоссарий\n\n"
        "| Термин | Описание |\n"
        "|--------|----------|\n"
        "| {Термин1} | {Описание1} |\n"
        "| {Термин2} | {Описание2} |\n",
    ),
    PRDSection(
        "Внешний API",
        True,
        "## 3. Внешний API\n\n"
        "### REST Endpoints\n\n"
        "#### POST /api/v1/{resource}\n\n"
        "**Описание:** {Описание операции}\n\n"
        "**Request:**\n"
        "```json\n"
        '{\n'
        '  "field1": "value1",\n'
        '  "field2": "value2"\n'
        "}\n"
        "```\n\n"
        "**Response:** `200 OK`\n"
        "```json\n"
        '{\n'
        '  "result": "success"\n'
        "}\n"
        "```\n",
    ),
    PRDSection(
        "Модель БД",
        True,
        "## 4. Модель БД\n\n"
        "### Таблица {table_name}\n\n"
        "| Поле | Тип | Описание |\n"
        "|------|-----|----------|\n"
        "| id | UUID | Primary key |\n"
        "| field1 | VARCHAR(255) | {Описание} |\n",
    ),
    PRDSection(
        "Sequence Flows",
        True,
        "## 5. Sequence Flows\n\n"
        "### {flow_name}\n\n"
        "```mermaid\n"
        "sequenceDiagram\n"
        "    participant Client\n"
        "    participant API\n"
        "    participant Service\n\n"
        "    Client->>API: {Operation}\n"
        "    API->>Service: Process\n"
        "    Service-->>API: Result\n"
        "    API-->>Client: Response\n"
        "```\n",
    ),
    PRDSection(
        "Внешние зависимости",
        True,
        "## 6. Внешние зависимости\n\n"
        "| Зависимость | Версия | Назначение |\n"
        "|-------------|--------|------------|\n"
        "| {dependency1} | {version} | {purpose} |\n"
        "| {dependency2} | {version} | {purpose} |\n",
    ),
    PRDSection(
        "Мониторинги",
        True,
        "## 7. Мониторинги\n\n"
        "### Metrics\n\n"
        "- `{metric1}`: {Описание}\n"
        "- `{metric2}`: {Описание}\n\n"
        "### Alerts\n\n"
        "- `{alert_name}`: {Условие} → {Действие}\n",
    ),
]

# Library profile: no DB, no monitoring
LIBRARY_SECTIONS = [
    PRDSection(
        "Назначение",
        True,
        "## 1. Назначение\n\n"
        "{Краткое описание библиотеки}\n\n"
        "**Ключевые возможности:**\n"
        "- {Возможность 1}\n"
        "- {Возможность 2}\n",
        500,
    ),
    PRDSection(
        "Глоссарий",
        True,
        "## 2. Глоссарий\n\n"
        "| Термин | Описание |\n"
        "|--------|----------|\n"
        "| {Термин1} | {Описание1} |\n",
    ),
    PRDSection(
        "Public API",
        True,
        "## 3. Public API\n\n"
        "### Class {ClassName}\n\n"
        "```python\n"
        "class {ClassName}:\n"
        '    """{Описание}"""\n\n'
        "    def method_name(self, param: str) -> Result:\n"
        '        """{Описание метода}."""\n'
        "```\n",
    ),
    PRDSection(
        "Data Structures",
        True,
        "## 4. Data Structures\n\n"
        "### {StructName}\n\n"
        "```python\n"
        "@dataclass\n"
        "class {StructName}:\n"
        '    """{Описание}"""\n'
        "    field1: str\n"
        "    field2: int\n"
        "```\n",
    ),
    PRDSection(
        "Usage Examples",
        True,
        "## 5. Usage Examples\n\n"
        "### Basic Usage\n\n"
        "```python\n"
        "from {library} import {Class}\n\n"
        "instance = {Class}()\n"
        "result = instance.method()\n"
        "```\n",
    ),
    PRDSection(
        "Внешние зависимости",
        True,
        "## 6. Внешние зависимости\n\n"
        "| Зависимость | Версия | Назначение |\n"
        "|-------------|--------|------------|\n"
        "| {dependency1} | {version} | {purpose} |\n",
    ),
    PRDSection(
        "Error Handling",
        True,
        "## 7. Error Handling\n\n"
        "### Exceptions\n\n"
        "- `{Error1}`: {Когда возникает}\n"
        "- `{Error2}`: {Когда возникает}\n\n"
        "### Error Recovery\n\n"
        "{Стратегия обработки ошибок}\n",
    ),
]

# CLI profile: command reference instead of API
CLI_SECTIONS = [
    PRDSection(
        "Назначение",
        True,
        "## 1. Назначение\n\n"
        "{Краткое описание CLI-инструмента}\n\n"
        "**Ключевые команды:**\n"
        "- `{cmd1}`: {Описание}\n"
        "- `{cmd2}`: {Описание}\n",
        500,
    ),
    PRDSection(
        "Глоссарий",
        True,
        "## 2. Глоссарий\n\n"
        "| Термин | Описание |\n"
        "|--------|----------|\n"
        "| {Термин1} | {Описание1} |\n",
    ),
    PRDSection(
        "Command Reference",
        True,
        "## 3. Command Reference\n\n"
        "### {command}\n\n"
        "```bash\n"
        "{tool} {command} [OPTIONS]\n\n"
        "**Options:**\n"
        "  --option1 TEXT  {Описание}\n"
        "  --option2 INT   {Описание}\n"
        "  --help          Show help\n"
        "```\n",
    ),
    PRDSection(
        "Configuration",
        True,
        "## 4. Configuration\n\n"
        "### Config File\n\n"
        "```yaml\n"
        "# {config_path}\n"
        "setting1: value1\n"
        "setting2: value2\n"
        "```\n\n"
        "### Environment Variables\n\n"
        "- `{VAR1}`: {Описание}\n"
        "- `{VAR2}`: {Описание}\n",
    ),
    PRDSection(
        "Usage Examples",
        True,
        "## 5. Usage Examples\n\n"
        "### Basic Usage\n\n"
        "```bash\n"
        "# {Пример 1}\n"
        "{tool} {command} --arg value\n\n"
        "# {Пример 2}\n"
        "{tool} {command} < input.txt\n"
        "```\n",
    ),
    PRDSection(
        "Exit Codes",
        True,
        "## 6. Exit Codes\n\n"
        "| Code | Meaning |\n"
        "|------|---------|\n"
        "| 0 | Success |\n"
        "| 1 | General error |\n"
        "| 2 | Invalid input |\n",
    ),
    PRDSection(
        "Error Handling",
        True,
        "## 7. Error Handling\n\n"
        "### Error Messages\n\n"
        "```\n"
        "Error: {Описание ошибки}\n\n"
        "Cause: {Причина}\n"
        "Solution: {Решение}\n"
        "```\n",
    ),
]

PROFILES: dict[ProjectType, PRDProfile] = {
    ProjectType.SERVICE: PRDProfile(ProjectType.SERVICE, SERVICE_SECTIONS),
    ProjectType.LIBRARY: PRDProfile(ProjectType.LIBRARY, LIBRARY_SECTIONS),
    ProjectType.CLI: PRDProfile(ProjectType.CLI, CLI_SECTIONS),
}


def get_profile(project_type: ProjectType) -> PRDProfile:
    """Get profile for a project type.

    Args:
        project_type: The project type

    Returns:
        The PRD profile for the project type
    """
    return PROFILES[project_type]


def get_section_limit(project_type: ProjectType, section_name: str) -> int | None:
    """Get character limit for a section.

    Args:
        project_type: The project type
        section_name: Name of the section

    Returns:
        Character limit or None if unlimited
    """
    profile = PROFILES[project_type]
    for section in profile.sections:
        if section.name == section_name:
            return section.max_chars
    return None
