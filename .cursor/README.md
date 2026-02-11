# Cursor Integration

**Canonical skill source:** `../prompts/skills/`

All SDP command prompts are unified in `prompts/skills/` with symlink adapters in tool folders.

## Usage

Use `@` prefix to invoke skills:

```bash
@feature "description"  # Unified entry point
@idea "description"     # Requirements gathering
@design {id}            # Plan workstreams
@build {ws-id}          # Execute workstream
@review {feature}       # Quality review
@deploy {feature}       # Production deployment
@debug "issue"          # Systematic debugging
@hotfix "critical"      # Emergency fix
@bugfix "issue"         # Quality fix
```

## See Also

- [CLAUDE.md](../CLAUDE.md) - Full protocol
- [prompts/skills/](../prompts/skills/) - Canonical skill definitions
- [.claude/skills/](../.claude/skills/) - Claude compatibility symlink
