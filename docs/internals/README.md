# Internals Documentation

Developer and maintainer documentation for SDP architecture and extensibility.

---

## Contents

- [Architecture](#architecture)
- [Extending SDP](#extending-sdp)
- [Contributing](#contributing)
- [Development Setup](#development-setup)

---

## Architecture

### System Overview

SDP is a workstream-driven development framework built on:

1. **Skill System** - Claude Code command handlers
2. **Quality Gates** - Automated quality enforcement
3. **Multi-Agent Coord** - Orchestrated agent execution
4. **Task Tracking** - Beads CLI integration
5. **Git Integration** - GitHub sync and deployment

### Component Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        SDP Framework                         │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Skills     │  │  Quality     │  │    Agents    │      │
│  │  (@feature)  │  │   Gates      │  │ (Orchestrator)│      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│         │                  │                  │              │
│         └──────────────────┴──────────────────┘              │
│                            │                                 │
│                    ┌───────────────┐                        │
│                    │   Core Layer  │                        │
│                    │  (workstream, │                        │
│                    │   feature,    │                        │
│                    │    design)    │                        │
│                    └───────────────┘                        │
│                            │                                 │
│         ┌──────────────────┼──────────────────┐              │
│         │                  │                  │              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    Beads     │  │    GitHub    │  │   Telegram   │      │
│  │  (Tasks)     │  │  (Sync/CI)   │  │ (Notify)     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### Module Structure

**Core Modules:**
- `sdp.core` - Workstream, feature, project map parsing
- `sdp.quality` - Quality gate validation
- `sdp.design` - Dependency graph and design
- `sdp.feature` - Product vision management
- `sdp.tdd` - TDD cycle runner

**Integration Modules:**
- `sdp.beads` - Beads CLI integration
- `sdp.github` - GitHub issues and projects
- `sdp.unified` - Multi-agent coordination
- `sdp.adapters` - AI-IDE adapters (Claude, Cursor, OpenCode)

**Support Modules:**
- `sdp.schema` - JSON schema validation
- `sdp.validators` - Workstream validation
- `sdp.prd` - PRD parsing and generation
- `sdp.errors` - Error framework

---

## Extending SDP

### Creating Custom Skills

**Directory:**
```
.claude/skills/{skill_name}/
├── SKILL.md          # Skill definition
└── examples/         # Usage examples
```

**Example:**
```markdown
# @myskill

Custom skill for my workflow.

## Usage
\`\`\`bash
@myskill "argument"
\`\`\`

## Process
1. Parse input
2. Do work
3. Output result
```

### Creating Custom Quality Gates

**File:** `quality-gate.toml`

```toml
[my_custom_check]
enabled = true
threshold = 100
```

**Python Implementation:**
```python
from sdp.quality.validator import QualityGateValidator

class MyCustomCheck:
    def validate(self, file_path: Path) -> CheckResult:
        # Custom validation logic
        pass
```

### Creating Custom Validators

**Location:** `src/sdp/validators/`

**Example:**
```python
from pathlib import Path
from sdp.validators.base import Validator

class MyValidator(Validator):
    """Custom validator for my use case."""

    def validate(self, file_path: Path) -> ValidationResult:
        # Validation logic
        return ValidationResult(
            passed=True,
            issues=[]
        )
```

---

## Contributing

### Contribution Workflow

1. **Fork** the repository
2. **Create feature branch** (`git checkout -b feature/my-change`)
3. **Make changes** following quality gates
4. **Write tests** (≥80% coverage)
5. **Run quality checks** (`@review`)
6. **Submit PR** with description

### Code Standards

**Quality Gates:**
- Coverage ≥80%
- mypy --strict
- ruff clean
- Files <200 LOC
- No bare exceptions

**Style Guide:**
- Follow PEP 8
- Use type hints everywhere
- Docstrings for modules/classes
- Max complexity CC=10

### Testing

**Test Structure:**
```
tests/
├── unit/           # Unit tests (fast)
├── integration/    # Integration tests
└── e2e/           # End-to-end tests
```

**Test Commands:**
```bash
# Run all tests
pytest

# Run fast tests only
pytest -m fast

# Run with coverage
pytest --cov=src/sdp --cov-report=term-missing

# Run specific test
pytest tests/unit/test_module.py -v
```

---

## Development Setup

### Local Development

**1. Clone Repository:**
```bash
git clone https://github.com/fall-out-bug/sdp.git
cd sdp
```

**2. Create Virtual Environment:**
```bash
python3 -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
```

**3. Install Dependencies:**
```bash
pip install -e .
pip install -e ".[dev]"  # Development dependencies
```

**4. Run Health Checks:**
```bash
sdp doctor
```

**5. Run Tests:**
```bash
pytest tests/ -v
```

### Development Tools

**Required:**
- Python 3.14+
- mypy (type checking)
- ruff (linting)
- pytest (testing)
- git (version control)

**Optional:**
- Beads CLI (task tracking)
- Claude Code (AI-IDE)
- Telegram CLI (notifications)

### Configuration

**1. Claude Code Settings:**
```json
{
  "skills": ["feature", "design", "build", "review", "deploy"],
  "hooks": {
    "pre-commit": "hooks/pre-commit.sh",
    "post-build": "hooks/post-build.sh"
  }
}
```

**2. Quality Gates:**
```toml
[coverage]
minimum = 80

[complexity]
max_cc = 10
```

**3. Git Hooks:**
```bash
ln -sf ../../hooks/pre-commit.sh .git/hooks/pre-commit
```

---

## Release Process

### Version Bumping

**File:** `src/sdp/__init__.py`

```python
__version__ = "0.5.0"  # Bump this
```

### Changelog

**File:** `CHANGELOG.md`

```markdown
## [0.5.0] - 2026-01-29

### Added
- New feature X
- New command Y

### Changed
- Improved Z
```

### Tagging

```bash
git tag -a v0.5.0 -m "Release v0.5.0"
git push origin v0.5.0
```

---

## Architecture Decisions

### Key ADRs

1. **[ADR-0001]** File-native consensus protocol
2. **[ADR-0002]** Reliability improvements
3. **[ADR-0003]** Progressive universal protocol
4. **[ADR-0004]** Unified progressive consensus

**Location:** `docs/adr/`

---

## Performance Considerations

### Optimization Targets

- **Test Suite:** < 30 seconds for full suite
- **Validation:** < 5 seconds per file
- **Skill Execution:** < 2 seconds overhead
- **Agent Coordination:** < 100ms message latency

### Profiling

```bash
# Profile test execution
pytest --profile

# Profile validation
python -m cProfile -s -m sdp.quality.validator check
```

---

## Security Considerations

### Secrets Management

**Never Commit:**
- `.env` files with secrets
- API keys or tokens
- Passwords or credentials

**Use Environment Variables:**
```bash
export GITHUB_TOKEN=ghp_*
export BEADS_SERVER_URL=http://localhost:8080
```

### Code Security

**Quality Gates:**
- No hardcoded secrets (enforced)
- No SQL injection patterns (enforced)
- No eval() usage (enforced)
- HTTPS URLs required (enforced)

---

## Troubleshooting

### Common Issues

**Import errors:**
```bash
export PYTHONPATH=src
```

**Mypy errors:**
```bash
mypy src/sdp/ --strict --no-error-summary
```

**Test failures:**
```bash
pytest tests/unit/test_failing.py -v --tb=short
```

**See:** [reference/error-handling.md](../reference/error-handling.md)

---

## See Also

- [architecture.md](architecture.md) - Detailed architecture
- [extending.md](extending.md) - Extension guide
- [contributing.md](contributing.md) - Contribution guide
- [development.md](development.md) - Development setup
- [reference/](../reference/) - User-facing documentation

---

**Version:** SDP v0.5.0
**Updated:** 2026-01-29
**Maintainers:** SDP Protocol Team
