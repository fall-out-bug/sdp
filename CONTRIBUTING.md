# Contributing to Spec-Driven Protocol

Thank you for your interest in contributing!

## Ways to Contribute

- **Report bugs** — Open an issue describing the problem
- **Suggest features** — Open an issue with your idea
- **Improve documentation** — Fix typos, add examples, clarify explanations
- **Add skills** — Create new agent skills in `.claude/skills/`
- **Add agents** — Create new agent definitions in `.claude/agents/`
- **Share integrations** — Document how you use SDP with other tools

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/sdp.git
   cd sdp
   ```
3. Create a branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Project Structure

```
sdp/
├── sdp-plugin/           # Go implementation (CLI + agents)
│   ├── cmd/              # CLI commands
│   └── internal/         # Core logic
├── src/sdp/              # Go source (graph, monitoring, synthesis)
├── tests/                # Go test suite
├── .claude/
│   ├── skills/           # AI agent skill definitions
│   └── agents/           # Multi-agent definitions
├── .cursor/              # Cursor IDE integration
├── .opencode/            # OpenCode integration
├── docs/
│   ├── PROTOCOL.md       # Core specification
│   ├── reference/        # API and command reference
│   ├── vision/           # Strategic vision documents
│   ├── drafts/           # Feature specifications
│   └── workstreams/      # Backlog and completed WS
├── hooks/                # Git hooks and validators
├── templates/            # Workstream templates
├── PRODUCT_VISION.md     # Product vision v3.0
├── CLAUDE.md             # Claude Code integration guide
├── AGENTS.md             # Agent instructions
└── go.mod                # Go module definition
```

## Using SDP for Contributions

For larger changes, use the SDP workflow:

1. **Requirements** — Run `@idea "description"` to create draft
2. **Design** — Run `@design idea-{slug}` to create workstreams
3. **Implement** — Run `@build 00-FFF-SS` for each workstream
4. **Review** — Run `@review F{FF}` to verify quality
5. **Deploy** — Run `@deploy F{FF}` when ready

## Pull Request Process

1. **Update documentation** if your change affects usage
2. **Write clear commit messages** (conventional commits)
3. **One feature per PR**
4. **Reference issues** in PR description

### PR Title Format

```
type: brief description

Examples:
- docs: add integration example
- feat: add @refactor skill
- fix: correct dependency resolution
```

## Code Style

- **Go** — Follow standard Go conventions, `gofmt`
- **Markdown** — Consistent formatting, no trailing whitespace
- **Skills** — Follow `.claude/skills/` SKILL.md format

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Version:** 0.9.0
