# /prd — Product Requirements Document Command

Generate and maintain PROJECT_MAP.md PRD documents with automatic diagram generation.

## Usage

```bash
/prd {project}                    # Create new PRD via interactive dialog
/prd {project} --update           # Update diagrams and frontmatter
/prd {project} --type {service|library|cli}  # Force project type
```

## Workflow

### Initial Creation

```markdown
User: /prd hw-checker

Claude: Let me help you create a PRD for hw-checker.

1. First, I'll analyze the project structure...
   Detected type: service (docker-compose.yml found)

2. I'll scaffold the PRD with 7 sections:
   - 1. Назначение
   - 2. Глоссарий
   - 3. Внешний API
   - 4. Модель БД
   - 5. Sequence Flows
   - 6. Внешние зависимости
   - 7. Мониторинги

3. For each section, I'll ask targeted questions...

Let's start with section 1 (Назначение):
- What is the primary purpose of hw-checker?
- What are the key capabilities?

[Continue interactive dialog for all sections]
```

### Update Mode

```markdown
User: /prd hw-checker --update

Claude: Updating PRD for hw-checker...

1. Parsing @prd annotations from code...
   Found 12 flow steps across 3 files

2. Generating diagrams...
   - sequence-submission-processing.mmd
   - sequence-submission-processing.puml
   - component-overview.mmd
   - deployment-production.puml

3. Updating frontmatter...
   diagrams_hash: abc123def456

4. Running validation...
   All sections within limits

PRD updated successfully!
```

## Project Types

### Service (docker-compose.yml detected)
7 sections: Purpose, Glossary, API, DB, Sequence Flows, Dependencies, Monitoring

### Library (default)
7 sections: Purpose, Glossary, Public API, Data Structures, Usage Examples, Dependencies, Error Handling

### CLI (cli.py with Click/Typer detected)
7 sections: Purpose, Glossary, Command Reference, Configuration, Usage Examples, Exit Codes, Error Handling

## Section Limits

- "Назначение" max 500 characters (enforced)
- Other sections have format-specific limits

## Diagrams

Diagrams are auto-generated from code annotations:

```python
from sdp.prd import prd_flow, prd_step

@prd_flow("submission-processing")
@prd_step(1, "Receive submission from queue")
async def process_submission(self, job: Job) -> RunResult:
    ...
```

Or in bash:

```bash
# @prd: flow=submission-processing, step=2, desc=Clone repository
git clone "$url"
```

Generated diagrams:
- Mermaid (.mmd) - for GitHub rendering
- PlantUML (.puml) - for external tools

## Validation

```bash
sdp-prd validate PROJECT_MAP.md    # Check section limits
```

## Integration

The `/codereview` command automatically checks PRD freshness by comparing diagrams_hash.

If mismatch detected:
```
❌ Диаграммы устарели
   Run: /prd hw-checker --update
```

## See Also

- `.claude/skills/prd/SKILL.md` - Claude Code skill definition
- `sdp/src/sdp/prd/` - Implementation
- `docs/workstreams/backlog/00-011-*.md` - Related workstreams
