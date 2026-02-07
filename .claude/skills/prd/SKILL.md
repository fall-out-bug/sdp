# @prd — PRD Generation Skill

Generate and maintain PROJECT_MAP.md PRD documents with automatic diagram generation.

## Usage

```
@prd "hw-checker"
```

## What This Skill Does

1. **Detects project type** (service/library/cli) from file structure
2. **Scaffolds PRD** with appropriate sections for the project type
3. **Generates diagrams** from @prd annotations in code
4. **Validates** section limits and format
5. **Updates frontmatter** with diagrams_hash

## Workflow

### Initial Creation (`@prd "project-name"`)

1. Analyze project structure for type detection
2. Present detected type and allow override
3. Guide through filling each PRD section interactively
4. Create `docs/PROJECT_MAP.md` with frontmatter
5. Generate initial diagram templates
6. Validate output

### Update Mode (`@prd "project-name" --update`)

1. Parse all @prd annotations from code
2. Regenerate diagrams from annotations
3. Calculate new diagrams_hash
4. Update PROJECT_MAP.md frontmatter
5. Run validation checks
6. Report changes

## Project Types

| Type | Trigger | Sections |
|------|---------|----------|
| **service** | docker-compose.yml exists | 7 sections with API, DB, Monitoring |
| **library** | default (no docker/cli) | 7 sections with Public API, Usage Examples |
| **cli** | cli.py with Click/Typer | 7 sections with Command Reference, Exit Codes |

## Section Limits

The following limits are enforced during validation:

- **"Purpose"**: max 500 characters
- **"Data Model"**: 1 line per field (max 120 chars per line)
- Other sections: format-specific limits

## Diagram Generation

Diagrams are generated from code annotations:

### Python

```python
from sdp.prd import prd_flow, prd_step

@prd_flow("submission-processing")
@prd_step(1, "Receive submission from queue")
async def process_submission(self, job: Job) -> RunResult:
    """Process single submission through SAGA orchestrator."""
    ...
```

### Bash/YAML

```bash
# @prd: flow=submission-processing, step=2, desc=Clone repository
git clone "$url" "$workspace"
```

### Generated Files

- `docs/diagrams/sequence-{flow_name}.mmd` - Mermaid diagram
- `docs/diagrams/sequence-{flow_name}.puml` - PlantUML diagram
- `docs/diagrams/component-overview.mmd` - Component template
- `docs/diagrams/deployment-production.puml` - Deployment template

## Validation Rules

The following validation checks are performed:

1. **Frontmatter completeness**: project_type, prd_version, last_updated
2. **Section limits**: Character counts and format rules
3. **Diagram freshness**: diagrams_hash matches current annotations

## Implementation

The skill uses the following modules:

- `sdp.src.sdp.prd.profiles` - PRD profiles and sections
- `sdp.src.sdp.prd.detector` - Project type auto-detection
- `sdp.src.sdp.prd.scaffold` - PRD template generation
- `sdp.src.sdp.prd.validator` - Section validation
- `sdp.src.sdp.prd.parser_python` - Python annotation parser
- `sdp.src.sdp.prd.parser_bash` - Bash annotation parser
- `sdp.src.sdp.prd.generator_mermaid` - Mermaid diagram generator
- `sdp.src.sdp.prd.generator_plantuml` - PlantUML diagram generator

## Quality Gates

Before considering a PRD complete:

- [ ] All 7 sections filled for detected project type
- [ ] "Purpose" section ≤ 500 characters
- [ ] Frontmatter contains project_type, prd_version, last_updated
- [ ] diagrams_hash set (after diagram generation)
- [ ] All diagrams generated and saved to docs/diagrams/
- [ ] Validation passes without errors

## Related Commands

- `/codereview` - Checks PRD freshness via diagrams_hash
- `/design` - Creates feature design (can use PRD as input)
- `/build` - Implements workstreams (can reference PRD sections)

## Example Session

```
User: @prd "hw-checker"

Assistant: Let me help you create a PRD for hw-checker.

Analyzing project structure...
✓ Detected: service (docker-compose.yml found)

I'll create a PRD with 7 sections for a service profile:

1. Purpose (max 500 chars)
2. Glossary
3. External API
4. Data Model
5. Sequence Flows
6. External Dependencies
7. Monitoring

Let's start with section 1: Purpose

What is the primary purpose of hw-checker? Please describe in 1-2 sentences.

[Interactive dialog continues for all sections...]

✓ PRD created: tools/hw_checker/docs/PROJECT_MAP.md
✓ Diagrams generated: 4 files in docs/diagrams/
✓ Validation passed

Next steps:
1. Add @prd annotations to key code files
2. Run @prd "hw-checker" --update to regenerate diagrams
```

## Context Files

- `docs/PROJECT_MAP.md` - The PRD document
- `docs/diagrams/*.mmd` - Mermaid diagrams
- `docs/diagrams/*.puml` - PlantUML diagrams
- `docs/workstreams/backlog/00-011-*.md` - Feature F011 workstreams
