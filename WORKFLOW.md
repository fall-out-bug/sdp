# Development Workflow

How this repository is developed using the Consensus Workflow framework.

## Process

```
GitHub Issue → Local Work → Consensus Workflow → Commit
```

1. **Create GitHub Issue** - Plan and discuss publicly
2. **Work locally** - Epic directories are gitignored during development
3. **Follow Consensus** - Run agents per protocol
4. **Commit artifacts** - Push only completed work

## Working Locally

```bash
# Create epic structure
mkdir -p docs/specs/epic_XX/consensus/{artifacts,decision_log,messages/inbox/{analyst,architect,tech_lead,developer,qa,devops}}

# Work on epic (gitignored until complete)
# docs/specs/epic_*/ is in .gitignore
```

## Running Agents

Use Claude Code with appropriate model for each role:

```bash
# Strategic roles (Analyst, Architect, Security)
/model opus

# Implementation roles (Developer, QA, DevOps)
/model sonnet  # or haiku for simpler tasks
```

See [modes/multi-agent/](modes/multi-agent/) for full protocol.

## Committing Completed Work

```bash
# Stage completed epic
git add docs/specs/epic_XX/

# Commit with context
git commit -m "feat(epic XX): [title]

Artifacts: requirements.json, architecture.json, implementation.md
Agents: Analyst, Architect, Tech Lead, Developer, QA"

# Push and close GitHub issue
git push
```

## What Gets Committed

**Commit:**
- Completed epic.md
- All consensus artifacts
- Decision logs
- Implementation code and tests

**Don't commit:**
- Work in progress
- Failed experiments
- API keys or secrets

## Tips

- Document decisions in decision_log/, not just what was built
- If prompts fail, document it to improve the framework
- Update framework based on real experience
