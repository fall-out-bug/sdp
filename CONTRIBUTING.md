# Contributing to Spec-Driven Protocol

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Ways to Contribute

- **Report bugs** - Open an issue describing the problem
- **Suggest features** - Open an issue with your idea
- **Improve documentation** - Fix typos, add examples, clarify explanations
- **Add command templates** - Enhance existing slash commands
- **Share integrations** - Document how you use SDP with other tools

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/sdp.git
   cd sdp
   ```
3. Create a branch for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Project Structure

```
consensus/
├── prompts/
│   ├── commands/       # Slash command prompts (/idea, /design, /build, etc.)
│   └── structured/     # Phase-based workflow prompts
├── docs/
│   ├── guides/         # Integration guides (Claude Code, Cursor)
│   ├── concepts/       # Core concepts (Clean Architecture, Artifacts, Roles)
│   ├── adr/            # Architecture decision records
│   └── specs/          # Feature specifications
├── .cursor/            # Cursor IDE slash commands
├── .cursorrules        # Cursor IDE rules
├── .claudecode/        # Claude Code configuration
├── hooks/              # Git hooks and validators
├── templates/          # Document templates
├── PROTOCOL.md         # SDP specification
├── RULES_COMMON.md     # Shared rules
└── MODELS.md           # Model recommendations
```

## Contribution Guidelines

### For Documentation

- Write in clear, concise English
- Include examples where helpful
- Keep formatting consistent with existing docs
- Test any code examples you include

### For Command Prompts

When adding or modifying command prompts in `prompts/commands/`:

1. **Follow the existing structure** - Use the format with sections like ALGORITHM, PRE-FLIGHT CHECKS, etc.
2. **Keep language-agnostic** - Don't hardcode specific technologies (use placeholders like `{language}`, `{framework}`)
3. **Include all required sections**:
   - GLOBAL RULES - Core principles
   - ALGORITHM - Step-by-step workflow
   - OUTPUT FORMAT - What to display to user
   - THINGS YOU MUST NEVER DO - Hard constraints
4. **Test with actual AI tools** - Verify the prompt works with Claude Code or Cursor

### For New Commands

To add a new slash command:

1. Create `prompts/commands/{command}.md` (full prompt)
2. Create `.cursor/commands/{command}.md` (quick reference for Cursor IDE)
3. Add skill to `.claudecode/skills/{command}/SKILL.md` (Claude Code integration)
4. Update `README.md` with command description
5. Update `MODELS.md` with model recommendation

### Code Style

- **English only** - All content must be in English (except README_RU.md)
- **Consistent formatting** - Follow existing Markdown style
- **No trailing whitespace**
- **End files with newline**

## Using SDP for Your Contributions

You're welcome to use SDP workflow for larger contributions:

### For Larger Changes (new features, major refactors)

1. **Requirements** - Run `/idea "{description}"` to create draft
2. **Design** - Run `/design idea-{slug}` to create workstreams
3. **Implement** - Run `/build WS-XXX-XX` for each workstream
4. **Review** - Run `/review F{XX}` to verify quality
5. **Deploy** - Run `/deploy F{XX}` when ready

**Recommended tools:**
- [Claude Code](docs/guides/CLAUDE_CODE.md) - CLI with multiple providers
- [Cursor IDE](docs/guides/CURSOR.md) - Visual IDE with slash commands

See [MODELS.md](MODELS.md) for model selection.

## Pull Request Process

1. **Update documentation** - If your change affects usage, update relevant docs
2. **Write clear commit messages** - Describe what and why
3. **One feature per PR** - Keep changes focused
4. **Reference issues** - Link to related issues in PR description

### PR Title Format

```
type: brief description

Examples:
- docs: add Python project example
- feat: add /refactor command
- fix: correct path in /build prompt
```

### PR Description Template

```markdown
## Summary
Brief description of changes

## Changes
- Change 1
- Change 2

## Testing
How you tested the changes

## Related Issues
Fixes #123
```

## Review Process

1. Maintainers will review your PR
2. Address any requested changes
3. Once approved, your PR will be merged

## Adding Examples

We welcome examples showing SDP in action:

1. Add to `examples/`
2. Include complete feature with workstreams
3. Add README explaining the example
4. Keep examples generic (no proprietary code)

## Reporting Bugs

When reporting bugs, include:

- What you expected to happen
- What actually happened
- Steps to reproduce
- Which AI tool you're using (Claude Code, Cursor, etc.)
- Relevant command or configuration

## Suggesting Features

When suggesting features:

- Describe the use case
- Explain why existing features don't solve it
- Propose a solution (optional)

## Questions?

- Check existing issues and documentation first
- Open a discussion for general questions
- Open an issue for specific problems

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing!
