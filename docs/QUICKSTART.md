# SDP Quick Start

Get from zero to a real first success in one repo.

Use this guide when you are adopting SDP in your own project. If you are contributing to SDP itself, start with `DEVELOPMENT.md` and `CONTRIBUTING.md`.

> This guide follows Local Mode from [PRODUCT_CONTRACT.md](PRODUCT_CONTRACT.md). For the full product definition, see that document.

## 0. Pick the Right Mode

| Mode | Start here when | Requires |
|------|-----------------|----------|
| Local Mode | You want the shortest path to a working setup | Supported IDE integration plus the `sdp` CLI |
| Operator Mode | You already want a queue-backed workflow across sessions | Beads plus prompt-driven workstreams |

**Start with Local Mode unless you already run Beads.** That is the current public happy path.

## 1. Install

Default install adds prompts, hooks, and the optional CLI:

```bash
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh
```

CLI only:

```bash
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only
```

Vendored as a submodule:

```bash
git submodule add https://github.com/fall-out-bug/sdp.git sdp
```

Supported integrations:

- `Claude Code`
- `Cursor`
- `OpenCode` / `Windsurf`
- `Codex`

If auto-detect misses your tool, rerun with `SDP_IDE=claude|cursor|opencode|codex`.

Use the public GitHub URL as the canonical submodule source. Do not point `.gitmodules` at a local sibling checkout such as `../sdp`.

## 2. Initialize the Project

```bash
sdp init --auto
```

Interactive setup is also available:

```bash
sdp init --guided
```

`sdp init` creates `.sdp/config.yml`, `.sdp/guard-rules.yml`, and refreshes the IDE integration already present in the repo. If no supported integration exists yet, it falls back to `.claude/`.

If your CLI is too old to support `--auto`, reinstall the binary:

```bash
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only
```

Then verify the environment:

```bash
sdp doctor
```

## 3. Get a First Success

The fastest guided proof that your install works:

```bash
sdp demo
```

The demo creates a temporary project and walks through:

1. `sdp init --guided`
2. `sdp doctor`
3. `sdp status --text`

If you want to skip the demo and use your real repo immediately, continue with the local CLI flow below.

## 4. Plan and Execute in Local Mode

Create workstreams from a feature description:

```bash
sdp plan "Add OAuth2"
```

Preview without writing:

```bash
sdp plan "Add OAuth2" --dry-run
```

Execute ready workstreams:

```bash
sdp apply --dry-run
sdp apply
```

Execute one workstream directly:

```bash
sdp build 00-001-01
```

Verify and inspect:

```bash
sdp verify 00-001-01
sdp status --text
sdp next
sdp log show
```

## 5. Advanced: Operator Mode with Beads

Use this only if you want a live queue across sessions.

Install Beads:

```bash
brew tap beads-dev/tap && brew install beads
```

Common queue commands:

```bash
bd ready
bd create --title="..."
bd close <id>
```

Once Beads is in place, SDP also installs prompt surfaces such as `/feature`, `/build`, `/review`, `/oneshot`, and `/strataudit`. That mode assumes workstreams and operator discipline already exist; it is not required for a first run.

Important distinction:

- `/deploy` is a prompt surface in the prompt bundle.
- `sdp deploy` is a CLI command that records an approval event after merge.
- `sdp deploy` does not merge branches or deploy infrastructure.

## 6. What to Use Next

- Use [CLI_REFERENCE.md](CLI_REFERENCE.md) for current command behavior.
- Use [PROTOCOL.md](PROTOCOL.md) for the current protocol overview.
- Use [PRODUCT_CONTRACT.md](PRODUCT_CONTRACT.md) for product definition and mode policy.
- Use [reference/README.md](reference/README.md) for the reference index and legacy-doc status.
- Use [MANIFESTO.md](MANIFESTO.md) for vision and rationale.
