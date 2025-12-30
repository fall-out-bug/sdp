# Cursor IDE Integration

This guide explains how to use Consensus Workflow with [Cursor](https://cursor.com).

## Basic Setup

1. Open your project in Cursor
2. Cursor will automatically read `.cursorrules` or `CLAUDE.md` for context

## Using Agent Prompts

Copy prompts from `prompts/` directory into Cursor chat:

```
# In Cursor Chat (Cmd/Ctrl + L):

@prompts/analyst_prompt.md
@docs/specs/epic_XX/epic.md

Analyze this epic and create requirements.json
```

## Recommended Workflow

### Sequential Agent Approach

Run each agent role in separate chat sessions:

1. **Analyst**: Create requirements
   ```
   @prompts/analyst_prompt.md
   @epic.md
   Create requirements.json
   ```

2. **Architect**: Design architecture
   ```
   @prompts/architect_prompt.md
   @requirements.json
   Create architecture.json, veto if Clean Architecture violated
   ```

3. **Tech Lead**: Plan implementation
   ```
   @prompts/tech_lead_prompt.md
   @architecture.json
   Create implementation.md with workstreams
   ```

4. **Developer**: Implement with TDD
   ```
   @prompts/developer_prompt.md
   @implementation.md
   Implement workstream 1 with tests
   ```

5. **QA**: Verify quality
   ```
   @prompts/qa_prompt.md
   Run tests, verify coverage
   ```

6. **DevOps**: Create deployment
   ```
   @prompts/devops_prompt.md
   Create deployment configuration
   ```

## Model Selection

Cursor allows selecting different models in Settings → Models.

For Consensus Workflow, consider:
- **Strategic roles** (Analyst, Architect): Use more capable models
- **Implementation roles** (Developer, QA): Faster models work well

Check Cursor's current model offerings in their settings.

## .cursorrules File

Create `.cursorrules` in project root with key protocol rules:

```
# Consensus Workflow Rules

## Language
ALL output must be in English.

## Message Format
JSON with compact keys: d, st, r, epic, sm, nx, artifacts

## Quality Gates
- No silent fallbacks (except: pass forbidden)
- Dependencies point inward (Clean Architecture)
- Test coverage ≥80% in touched areas

## Artifacts Location
- Artifacts: docs/specs/{epic}/consensus/artifacts/
- Messages: docs/specs/{epic}/consensus/messages/inbox/{agent}/
- Decision logs: docs/specs/{epic}/consensus/decision_log/
```

## Tips

1. **Keep context focused**: Only @-mention files relevant to current agent role
2. **Use quick prompts**: For routine tasks, use `prompts/quick/` (shorter, saves tokens)
3. **Verify format**: Ask agent to show JSON before saving to catch format errors
4. **Clear chat**: Start fresh chat for each agent to avoid context confusion

## Resources

- [PROTOCOL.md](../../PROTOCOL.md) - Full consensus protocol
- [RULES_COMMON.md](../../RULES_COMMON.md) - Shared agent rules
- [prompts/](../../prompts/) - All agent prompts

---

**Note**: This guide covers basic integration. Cursor features change frequently - check [cursor.com/changelog](https://cursor.com/changelog) for current capabilities.
