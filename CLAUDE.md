# SDP: Spec-Driven Protocol

AI-native dev with workstreams, gates, TDD.

## Decision Tree

```
New? → @init → @vision → @reality → @feature
Demo? → sdp demo
No → State? → @reality --quick
WS? → @oneshot <id>
No → @feature "X"
    Choose mode:
    ├── --default (full interactive)
    ├── --quick (@design only, 0 questions)
    └── --auto (from roadmap, 0 questions)
```

## Path

- **Try?** → [Try section](#try)
- **Solo?** → [Adopt section](#adopt)
- **Team?** → [Scale section](#scale)

## Try

```bash
sdp demo                    # Guided walkthrough
# or install manually:
# go install github.com/fall-out-bug/sdp/sdp-plugin/cmd/sdp@latest
@init
@feature "X"
@build 00-001-01
```

## Adopt (Prereq: [PROTOCOL.md](docs/PROTOCOL.md) | [Skills](docs/reference/skills.md) | [CLI](docs/CLI_REFERENCE.md))

```bash
sdp init
@reality --quick
@feature "X"
@oneshot <id>
```

## Scale (Prereq: Team + [PROTOCOL.md](docs/PROTOCOL.md) | [Agents](.claude/agents/README.md) | [Design](docs/reference/design-spec.md) | [Principles](docs/reference/PRINCIPLES.md))

```bash
brew install beads
@vision "X"
@feature "X"
@oneshot <id>
@review <id>
@deploy <id>
```

**Key:** Aggregate=container, Leaf=executable, Feature=5-30 WS | **Format:** `PP-FFF-SS` | **Done:** @review APPROVED + @deploy

**Commands:** @vision @reality @feature @oneshot @build @review @deploy

**v0.9.8** | [Protocol](docs/PROTOCOL.md) | [Ref](docs/reference/README.md)
