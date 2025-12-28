# Contributing to Consensus Workflow

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Ways to Contribute

- **Report bugs** - Open an issue describing the problem
- **Suggest features** - Open an issue with your idea
- **Improve documentation** - Fix typos, add examples, clarify explanations
- **Add new agent roles** - Create prompts for specialized domains
- **Share integrations** - Document how you use Consensus with other tools

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/consensus.git
   cd consensus
   ```
3. Create a branch for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Project Structure

```
consensus/
├── prompts/           # Agent prompt templates
│   └── quick/         # Token-optimized prompts
├── docs/
│   ├── guides/        # Integration guides
│   └── examples/      # Worked examples
├── PROTOCOL.md        # Consensus protocol spec
├── RULES_COMMON.md    # Shared agent rules
└── MODELS.md          # Model recommendations
```

## Contribution Guidelines

### For Documentation

- Write in clear, concise English
- Include examples where helpful
- Keep formatting consistent with existing docs
- Test any code examples you include

### For Prompts

When adding or modifying agent prompts:

1. **Follow the existing structure** - Use the JSON format with `meta`, `context`, `mission`, `workflow`, etc.
2. **Keep language-agnostic** - Don't hardcode specific technologies (use placeholders like `{database}`)
3. **Include all required sections**:
   - `meta` - Role metadata and veto powers
   - `mission` - Clear one-line purpose
   - `workflow` - Step-by-step instructions
   - `boundaries` - Must do / must not do
   - `veto` - Triggers and actions
4. **Create both full and quick versions** - Full prompt in `prompts/`, quick version in `prompts/quick/`
5. **Test with actual AI tools** - Verify the prompt works with Claude Code or Cursor

### For New Agent Roles

To add a new agent role:

1. Create `prompts/{role}_prompt.md` (full version)
2. Create `prompts/quick/{role}_quick.md` (quick version)
3. Add role to `consensus_architecture.json`
4. Update `README.md` with role description
5. Update `MODELS.md` with model recommendation

### Code Style

- **English only** - All content must be in English
- **Consistent formatting** - Follow existing Markdown/JSON style
- **No trailing whitespace**
- **End files with newline**

## Using Consensus for Contributions (Recommended)

This repository itself is maintained using the Consensus Workflow with AI models. While not required, you're encouraged to use this approach for larger contributions:

### For Small Changes (bug fixes, typos, minor updates)
Regular PR workflow is perfectly fine - just submit your changes directly.

### For Larger Changes (new features, major refactors)
Consider using the consensus approach:

1. **Define requirements** - Use `prompts/analyst_prompt.md` to clarify scope
2. **Validate design** - Use `prompts/architect_prompt.md` for architectural decisions
3. **Plan implementation** - Use `prompts/tech_lead_prompt.md` for task breakdown
4. **Implement** - Use `prompts/developer_prompt.md` for coding with TDD
5. **Verify quality** - Use `prompts/qa_prompt.md` for testing

**Benefits:**
- Clear documentation of design decisions
- Better architecture consistency
- Reduced back-and-forth during PR review
- Can include consensus artifacts (`requirements.json`, `architecture.json`) in PR for transparency

**Recommended tools:**
- [Claude Code CLI](docs/guides/CLAUDE_CODE.md) - Multi-provider support (Gemini, Claude, OpenAI, Ollama)
- [Cursor IDE](docs/guides/CURSOR.md) - Visual multi-agent mode with parallel execution

See [MODELS.md](MODELS.md) for cost-effective model selection (Gemini 3 Flash recommended for most tasks).

## Pull Request Process

1. **Update documentation** - If your change affects usage, update relevant docs
2. **Write clear commit messages** - Describe what and why
3. **One feature per PR** - Keep changes focused
4. **Reference issues** - Link to related issues in PR description
5. **(Optional) Include consensus artifacts** - If you used the workflow, attach artifacts for transparency

### PR Title Format

```
type: brief description

Examples:
- docs: add example for Python project
- prompt: add database architect role
- fix: correct inbox path in developer prompt
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

We welcome examples showing Consensus in action:

1. Add to `docs/examples/`
2. Include complete epic with all artifacts
3. Add README explaining the example
4. Keep examples generic (no proprietary code)

## Reporting Bugs

When reporting bugs, include:

- What you expected to happen
- What actually happened
- Steps to reproduce
- Which AI tool you're using (Claude Code, Cursor, etc.)
- Relevant prompt or configuration

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
