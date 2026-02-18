# Agent Versions Index (WS-067-13: AC5)

Track agent version compatibility and changes.

## Core Agents

| Agent | Version | Last Updated | Purpose |
|-------|---------|--------------|---------|
| orchestrator | 2.2.0 | 2026-02-12 | Feature orchestration |
| supervisor | 1.1.0 | 2026-02-08 | Hierarchical coordination |
| builder | 1.1.0 | 2026-02-12 | TDD workstream execution |
| reviewer | 1.0.0 | 2026-01-30 | Code review (17-point) |
| spec-reviewer | 1.0.0 | 2026-02-07 | Spec compliance verification |
| implementer | 1.0.0 | 2026-02-11 | TDD implementation |
| tester | 1.0.0 | 2026-01-31 | QA testing |
| planner | 1.0.0 | 2026-01-30 | Workstream planning |

## Specialist Agents

| Agent | Version | Last Updated | Purpose |
|-------|---------|--------------|---------|
| architect | 1.0.0 | 2026-01-31 | System design |
| analyst | 1.0.0 | 2026-01-31 | Business analysis |
| developer | 1.0.0 | 2026-02-07 | Code implementation |
| deployer | 1.0.0 | 2026-01-30 | Deployment automation |
| ci-reviewer | 1.0.0 | 2026-02-11 | CI/CD review |
| workflow-auditor | 1.0.0 | 2026-02-11 | Process drift audit |

## Agents Without Version (Need Update)

The following agents lack version frontmatter and should be updated:

- business-analyst.md
- devops.md
- implementer.md (has content but no frontmatter)
- product-manager.md
- qa.md
- security.md
- spec-reviewer.md (has content but no frontmatter)
- sre.md
- system-architect.md
- systems-analyst.md
- tech-lead.md
- technical-decomposition.md

## Version Format

```yaml
---
name: agent-name
description: Brief description
model: inherit
version: X.Y.Z
changes:
  - Change description for this version
---
```

## Version Compatibility

- **Major (X)**: Breaking changes to agent behavior
- **Minor (Y)**: New capabilities, backward compatible
- **Patch (Z)**: Bug fixes, documentation updates

---

**Last Updated:** 2026-02-13
**Workstream:** WS-067-13
