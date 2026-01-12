# Prompts

Command-based workflow for SDP.

## Slash Commands

Use these for one-shot execution:

```
/idea "{description}"           # Requirements gathering
/design idea-{slug}              # Workstream planning
/build WS-XXX-XX                 # Execute workstream
/review F{XX}                    # Quality review
/deploy F{XX}                    # Deployment
/issue "{description}"           # Debug and route
/hotfix "{description}"          # Emergency fix
/bugfix "{description}"          # Quality fix
/oneshot F{XX}                   # Autonomous execution (Task orchestrator)
/oneshot F{XX} --background      # Background execution
/oneshot F{XX} --resume {id}     # Resume from checkpoint
```

See `prompts/commands/*.md` for full prompts.

## When to Use Which Command

| Task | Command | Model |
|------|---------|-------|
| Start new feature | `/idea` | Sonnet |
| Plan workstreams | `/design` | Opus |
| Implement workstream | `/build` | Haiku |
| Review feature | `/review` | Opus |
| Deploy to production | `/deploy` | Haiku |
| Debug issue | `/issue` | Sonnet |
| Emergency fix | `/hotfix` | Haiku |
| Quality fix | `/bugfix` | Haiku |
| Autonomous execution (Task-based) | `/oneshot` | Opus |

## Guardrails and Quality Gates

See `PROTOCOL.md` — everything in one place.
