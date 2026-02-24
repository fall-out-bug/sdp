---
name: readme
description: Agent index for SDP multi-agent coordination.
tools:
  read: true
---

# SDP Agent Index

12 agents for feature development. Each role has one clear purpose.

| Agent | Purpose |
|-------|---------|
| orchestrator | @oneshot — autonomous feature execution |
| implementer | @build — TDD workstream execution |
| spec-reviewer | @build — specification compliance |
| reviewer | @review — quality validation |
| planner | @feature — workstream decomposition |
| deployer | @deploy — deployment orchestration |
| qa | @review — quality assurance |
| security | @review — security review |
| devops | @review — CI/CD review |
| sre | @review — reliability review |
| tech-lead | @review — technical leadership |
| architect | @design, @feature — system design |

**See:** `AGENTS.md` for workflow. Each agent file: `.opencode/agents/{name}.md`
