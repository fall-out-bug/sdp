# SDP: Spec-Driven Protocol

**Workstream-driven development** for AI agents with multi-language support.

**Plugin Version:** Language-agnostic (Python, Java, Go)

---

## Quick Start

```bash
# Install plugin (no Python required)
git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp
cp -r ~/.claude/sdp/prompts/* .claude/

# Create feature (interactive)
@feature "Add user authentication"

# Plan workstreams
@design feature-auth

# Execute workstream
@build 00-001-01

# Review quality
@review feature-auth

# Deploy to production
@deploy feature-auth
```

---

## Core Concepts

### Hierarchy

| Level | Scope | Size | Example |
|-------|-------|------|---------|
| **Release** | Product milestone | 10-30 Features | R1: Submissions E2E |
| **Feature** | Major feature | 5-30 Workstreams | F24: Unified Workflow |
| **Workstream** | Atomic task | SMALL/MEDIUM/LARGE | WS-060: Domain Model |

### Workstream Size

- **SMALL**: < 500 LOC, < 1500 tokens
- **MEDIUM**: 500-1500 LOC, 1500-5000 tokens
- **LARGE**: > 1500 LOC → split into 2+ WS

⚠️ **NO TIME-BASED ESTIMATES** - Use scope metrics (LOC/tokens) only.

---

## Workstream Flow

```
┌────────────┐    ┌────────────┐    ┌────────────┐    ┌────────────┐
│  ANALYZE   │───→│    PLAN    │───→│  EXECUTE   │───→│   REVIEW   │
│  (Sonnet)  │    │  (Sonnet)  │    │   (Auto)   │    │  (Sonnet)  │
└────────────┘    └────────────┘    └────────────┘    └────────────┘
     │                  │                  │                  │
     ▼                  ▼                  ▼                  ▼
  Map WS           Plan WS            Code           APPROVED/FIX
```

---

## Quality Gates

Every workstream must pass quality gates. **Commands are language-specific:**

### Test Coverage ≥80%

| Language | Command |
|----------|---------|
| **Python** | `pytest tests/unit/ --cov=src/ --cov-fail-under=80` |
| **Java** | `mvn verify` (JaCoCo report) or `gradle test jacocoTestReport` |
| **Go** | `go test -coverprofile=coverage.out && go tool cover -func=coverage.out` |

### Type Checking

| Language | Command |
|----------|---------|
| **Python** | `mypy src/ --strict` |
| **Java** | `javac -Xlint:all` (compiler checks) |
| **Go** | `go vet ./...` |

### Linting

| Language | Command |
|----------|---------|
| **Python** | `ruff check src/` |
| **Java** | `mvn checkstyle:check` or `gradle checkstyleMain` |
| **Go** | `golint ./...` or `staticcheck ./...` |

### File Size <200 LOC

| Language | Command |
|----------|---------|
| **All** | Count lines per file (language-agnostic) |

**Forbidden:**
- ❌ Bare exceptions (`except:`, `catch(Exception e)`, `recover()` without check)
- ❌ Files > 200 LOC
- ❌ Coverage < 80%
- ❌ Time-based estimates
- ❌ TODO without followup WS

---

## Language-Specific Workflows

### Python Projects

```bash
# Prerequisites
pip install pytest pytest-cov mypy ruff

# Workflow
@feature "Add REST API"
@design feature-rest-api
@build 00-001-01
# Runs: pytest tests/ -v
# Quality: pytest --cov=src/ --cov-fail-under=80, mypy src/ --strict, ruff check src/
```

### Java Projects

```bash
# Prerequisites
# Maven: mvn verify (runs JaCoCo)
# Gradle: gradle test jacocoTestReport

# Workflow
@feature "Add REST API"
@design feature-rest-api
@build 00-001-01
# Runs: mvn test
# Quality: JaCoCo coverage ≥80%, javac -Xlint:all, checkstyle
```

### Go Projects

```bash
# Prerequisites
# Go 1.21+ with go tool cover, go vet, golint

# Workflow
@feature "Add REST API"
@design feature-rest-api
@build 00-001-01
# Runs: go test ./...
# Quality: go tool cover -func=coverage.out (≥80%), go vet ./..., golint ./...
```

---

## Project Type Detection

SDP automatically detects project type:

1. **Python**: `pyproject.toml`, `setup.py`, or `requirements.txt` present
2. **Java**: `pom.xml` or `build.gradle` present
3. **Go**: `go.mod` present
4. **Node.js**: `package.json` present
5. **Rust**: `Cargo.toml` present

If multiple build files exist, SDP prompts user to specify.

---

## Multi-Agent Coordination

SDP integrates multi-agent coordination with role-based message routing.

### Agent Spawning

```python
from sdp.unified.agent.spawner import AgentSpawner, AgentConfig

# Spawn agents
spawner = AgentSpawner()
builder = spawner.spawn_agent(AgentConfig(
    name="builder",
    prompt="Execute workstreams with TDD...",
))
```

### Message Routing

```python
from sdp.unified.agent.router import SendMessageRouter, Message

router = SendMessageRouter()
router.send_message(Message(
    sender="orchestrator",
    content="Execute 00-060-01",
    recipient=builder,
))
```

---

## Workstream Format

Workstreams use PP-FFF-SS format:
- **PP**: Project ID (00 = SDP core, 01-99 = custom)
- **FFF**: Feature ID (001-999)
- **SS**: Sequence (01-99)

Example: `00-001-01` = Project 00, Feature 001, Workstream 01

---

## Quality Gate Enforcement

### AI-Based Validation

SDP plugin uses AI validators instead of static analysis tools:

1. **Coverage Validator**: Analyzes test coverage by reading code
2. **Architecture Validator**: Checks Clean Architecture layer separation
3. **Error Validator**: Finds unsafe exception handling
4. **Complexity Validator**: Identifies overly complex code

**Gate:** All validators must PASS for workstream approval.

### Tool-Based Validation (Optional)

If language tools are available, SDP uses them:

```bash
# Python: pytest, mypy, ruff
# Java: Maven/Gradle, JaCoCo, checkstyle
# Go: go test, go tool cover, go vet
```

---

## Documentation

- [Tutorial](TUTORIAL.md) - Full workflow guide
- [Python Examples](examples/python/) - Python-specific guides
- [Java Examples](examples/java/) - Java-specific guides
- [Go Examples](examples/go/) - Go-specific guides
- [Migration Guide](MIGRATION.md) - From Python SDP to plugin
