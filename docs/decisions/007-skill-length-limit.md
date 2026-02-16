# ADR-007: Skill Length Limit

## Status

Accepted

## Context

Current skills are 400-500 lines. Evidence shows:
- Agents cherry-pick rules from long prompts
- Later sections get ignored
- Duplicate instructions lead to confusion

Research on LLM prompt effectiveness suggests:
- Shorter prompts have higher compliance
- Clear structure improves following
- References work better than inline detail

## Decision

Limit skills to **100 lines** (excluding code examples in docs/).

### Structure

1. **Header** (10 lines): name, description, tools
2. **Purpose** (5 lines): what and when
3. **Quick Reference** (10 lines): table of steps
4. **Workflow** (50 lines): step-by-step, 5-10 lines each
5. **Quality Gates** (5 lines): reference to docs
6. **Errors** (10 lines): common errors table
7. **See Also** (10 lines): links to detailed docs

### Details go to docs/

- `docs/reference/{skill}-spec.md` — Full specification
- `docs/examples/{skill}/` — Code examples
- `docs/reference/quality-gates.md` — Shared quality gates

## Consequences

### Positive
- Higher agent compliance
- Consistent skill structure
- Easier maintenance
- Clear separation of concerns

### Negative
- Need to create docs/reference/ files
- Migration effort for existing skills
- Agents need to follow references

## Validation

`sdp skill validate` checks:
- Line count ≤ 100 (warning), ≤ 150 (error)
- Required sections present
- References resolve
