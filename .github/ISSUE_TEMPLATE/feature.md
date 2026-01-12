---
name: Feature
about: New feature for consensus workflow development
title: 'Feature XX: [Title]'
labels: feature
assignees: ''
---

## Goal

[Clear one-line description of what this epic achieves]

## Context

[Why is this needed? What problem does it solve?]

## Scope

### In Scope
- Item 1
- Item 2

### Out of Scope
- Item 1
- Item 2

## Success Criteria

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

## User Stories

1. **As a [role]**, I want [feature], so that [benefit]
2. **As a [role]**, I want [feature], so that [benefit]

## Technical Constraints

- Constraint 1
- Constraint 2

## Consensus Workflow

This epic will follow the full Consensus Workflow protocol:

- [ ] **Analyst**: Create requirements.json
- [ ] **Architect**: Design architecture.json (may veto)
- [ ] **Tech Lead**: Create implementation.md
- [ ] **Developer**: Implement with TDD
- [ ] **QA**: Verify and create test_results.md
- [ ] **DevOps**: Package/deploy (if applicable)

## Models to Use

Per [MODELS.md](../../../MODELS.md):
- **Strategic** (Analyst, Architect, Security): Claude Opus 4.5
- **Implementation** (Tech Lead, Developer, QA): Gemini 3 Flash
- **Automation** (DevOps): Qwen3-Coder (free via Ollama)

## Deliverables

When epic is complete, the following will be committed:

- [ ] `docs/specs/epic_XX_name/epic.md`
- [ ] `docs/specs/epic_XX_name/consensus/artifacts/*.json`
- [ ] `docs/specs/epic_XX_name/consensus/decision_log/*.md`
- [ ] Code/tests/docs created during implementation
- [ ] Lessons learned document

## Notes

[Any additional context, links, references]
