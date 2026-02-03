# Spec-Driven Protocol (SDP) Plugin

**Workstream-driven development for AI agents with multi-language support.**

## Features

✅ **TDD Discipline** - Red → Green → Refactor cycle enforced by prompts
✅ **Clean Architecture** - Layer separation validated by AI
✅ **Quality Gates** - Coverage ≥80%, type safety, error handling
✅ **Multi-Language** - Python, Java, Go support
✅ **No Installation Required** - Prompts work standalone
✅ **Optional Binary** - Go CLI for init, doctor, hooks

## Quick Start

### Option 1: Manual Installation (No Python Required)

```bash
# 1. Clone plugin repository
git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp

# 2. Copy prompts to your project
cp -r ~/.claude/sdp/prompts/* .claude/

# 3. Start development
@feature "Add user authentication"
@design feature-auth
@build 00-001-01
```

### Option 2: With Optional Go Binary

```bash
# Download binary (macOS arm64 example)
curl -L https://github.com/ai-masters/sdp/releases/latest/download/sdp-darwin-arm64 -o sdp
chmod +x sdp

# Initialize project
./sdp init
./sdp doctor
```

## What's Included

### Skills (18 total)

Core workflow skills:
- `@feature` - Progressive vision/requirements gathering
- `@design` - Workstream planning with dependencies
- `@build` - Execute workstream with TDD cycle
- `@review` - Quality check with AI validators
- `@deploy` - Deployment workflow

Support skills:
- `@idea` - Requirements gathering
- `@issue` - Bug classification
- `@debug` - Systematic debugging
- `/help` - Skill discovery

### Agents (11 total)

Multi-agent coordination:
- `planner` - Workstream decomposition
- `builder` - Workstream execution
- `reviewer` - Quality validation
- `tester` - Test strategy
- `architect` - System design
- And more...

### Validators (4 total)

AI-based quality validation:
- `/coverage-validator` - Test coverage analysis
- `/architecture-validator` - Clean Architecture checks
- `/error-validator` - Error handling audit
- `/complexity-validator` - Complexity analysis

## Language Support

| Language | Tests | Coverage | Type Check | Lint |
|----------|-------|----------|------------|------|
| Python   | pytest | pytest-cov | mypy | ruff |
| Java     | Maven/Gradle | JaCoCo | javac | checkstyle |
| Go       | go test | go tool cover | go vet | golint |

## Documentation

- [Full Tutorial](docs/TUTORIAL.md)
- [Python Examples](docs/examples/python/)
- [Java Examples](docs/examples/java/)
- [Go Examples](docs/examples/go/)

## Migration from Python SDP

If you're using the Python `sdp` CLI tool:

✅ **Your existing workstreams still work** (prompts are compatible)
✅ **Git hooks continue to work** (use Go binary for convenience)
⚠️ **Quality checks now use AI validation** (no Python required)
📖 **See [MIGRATION.md](MIGRATION.md)** for details

## Directory Structure

```
sdp-plugin/
├── plugin.json           # Plugin manifest
├── README.md             # This file
├── prompts/
│   ├── skills/           # 18 workflow skills
│   ├── agents/           # 11 agent definitions
│   └── validators/       # 4 AI validators
└── docs/
    ├── TUTORIAL.md       # Full tutorial
    ├── MIGRATION.md      # Migration guide
    └── examples/
        ├── python/       # Python quick start
        ├── java/         # Java quick start
        └── go/           # Go quick start
```

## License

MIT © MSU AI Masters

## Version

1.0.0 (Claude Plugin Distribution)
