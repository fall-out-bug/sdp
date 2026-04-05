# Contributing to SDP

Thank you for contributing.

Start with [DEVELOPMENT.md](DEVELOPMENT.md) for local setup and build commands.

## Ways to Contribute

- Report bugs with a reproducible example
- Improve onboarding and reference docs
- Fix CLI behavior in `sdp-plugin/`
- Improve prompt and agent definitions in `prompts/`
- Add or document integrations for supported harnesses

## Getting Started

1. Fork the repository.
2. Clone your fork.
3. Create a branch for one coherent change.

```bash
git clone https://github.com/YOUR_USERNAME/sdp.git
cd sdp
git checkout -b feature/your-feature-name
```

## Repository Layout

```text
sdp/
├── sdp-plugin/           # Go CLI implementation
│   ├── cmd/              # CLI entry points
│   └── internal/         # CLI/runtime packages
├── src/sdp/              # Root-module Go packages
├── prompts/              # Canonical prompt and agent source
│   ├── commands/
│   ├── skills/
│   └── agents/
├── .claude/              # Claude adapter around prompts/
├── .cursor/              # Cursor adapter around prompts/
├── .opencode/            # OpenCode adapter around prompts/
├── .codex/               # Codex adapter around prompts/
├── docs/                 # Onboarding, protocol, and reference docs
├── hooks/                # Validation and git-hook support
├── schema/               # JSON schemas and contracts
└── templates/            # Project and workflow templates
```

## Build and Test

SDP currently uses Go `1.26` in both the root module and `sdp-plugin/`.

Build the CLI:

```bash
cd sdp-plugin
CGO_ENABLED=0 go build -o sdp ./cmd/sdp
```

Run the main CLI test suite:

```bash
cd sdp-plugin
go test ./...
```

Run root-module tests when your change touches root packages:

```bash
go test ./...
```

If you change installer behavior, also run:

```bash
sh scripts/test-install-project.sh
```

## What to Edit

- Edit `prompts/` for prompt or agent behavior.
- Do not hand-edit `.claude/`, `.cursor/`, `.opencode/`, or `.codex/skills/sdp` as source files.
- Edit `sdp-plugin/` for CLI behavior.
- Update docs when public behavior changes.

If docs and runtime disagree, fix docs to match shipped behavior. Do not document planned behavior as if it already exists.

## Using SDP While Contributing

Two workflows exist today:

- **Local Mode:** `sdp init`, `sdp doctor`, `sdp plan`, `sdp apply`, `sdp verify`, `sdp status`, `sdp next`
- **Operator Mode:** prompt surfaces plus Beads-backed queue management

For most contributors, Local Mode is the simplest way to exercise the current product.

If you use prompt surfaces during development, treat them as harness-specific adapters over the same stage model. Important distinction:

- `sdp deploy` records an approval event after merge
- `sdp deploy` does not merge your branch or deploy infrastructure

## Pull Requests

Before opening a PR:

- run the relevant Go tests for the code you changed
- update user-facing docs when behavior changed
- keep prompt edits in `prompts/`
- keep one feature or fix per PR
- reference the issue or problem statement in the PR description

Suggested PR titles:

```text
docs: clarify codex onboarding
feat: add init preflight check
fix: correct apply status output
```

## Code and Doc Style

- Go: standard Go formatting with `gofmt`
- Markdown: short sentences, clear headings, concrete examples
- Prompts: keep canonical source in `prompts/`

Modern Go guidance used in this repo:

- Prefer `slices.SortFunc` over `sort.Slice`
- Prefer `strings.Cut` over manual split or index logic
- Prefer `strings.CutPrefix` and `strings.CutSuffix` over trim-after-check
- Prefer `slices.Contains`, `maps.Copy`, and `maps.Clone` over handwritten loops
- Prefer `any` over `interface{}`

## Prompt Source of Truth

All prompt and agent definitions have one canonical source:

| Content | Canonical path |
|---------|----------------|
| Commands | `prompts/commands/` |
| Skills | `prompts/skills/` |
| Agents | `prompts/agents/` |

Tool-specific directories are adapters. Edit `prompts/`, not the adapters.

## License

By contributing, you agree that your contributions are licensed under the MIT License.
