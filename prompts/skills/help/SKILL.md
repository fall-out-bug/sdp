---
name: help
description: Interactive skill discovery and guidance
tools:
  - Read
---

# @help - Skill Discovery

Help users find the right skill for their task.

## Quick Reference

```
@help                        # List all skills
@help "how to fix a bug"     # Get skill recommendation
@help workflow               # Show workflow skills
```

## Trigger

When user types:
- `@help`
- `@help "query"`
- `@help what skill to use`

## Workflow

### 1. Understand Intent

What is the user trying to do?
- New feature → Workflow skills
- Bug fix → Fixes/Debugging skills
- Deployment → Deployment skills
- Not sure → Ask clarifying question

### 2. Match to Skill

Map intent to skill category:

```
User intent:
├── New feature → @feature (or @idea → @design → @build)
├── Fix bug
│   ├── Production down (P0) → @hotfix
│   ├── Quality issue (P1/P2) → @bugfix
│   └── Needs triage → @issue
├── Debug problem → /debug
├── Execute workstream → @build
├── Deploy feature → @review → @deploy
└── Not sure → Ask clarifying question
```

### 3. Provide Guidance

- Show usage example
- Explain when to use
- Mention related skills

## Example Interactions

### Bug Fix Routing

**User:** @help "I found a bug in login"

**Response:**
> Based on your description, here are your options:
>
> 1. **@issue** - If you need to classify severity first
>    ```
>    @issue "Bug in login"
>    ```
>
> 2. **@bugfix** - If it's a quality issue (not critical)
>    ```
>    @bugfix "Login bug description"
>    ```
>
> 3. **/debug** - If you need to investigate the root cause
>    ```
>    /debug "Login fails under condition X"
>    ```
>
> Which approach fits your situation?

### New Feature

**User:** @help how to add a new feature

**Response:**
> For new features, use **@feature** for the full workflow:
>
> ```
> @feature "Add user comments"
> ```
>
> This will:
> 1. Interview you about requirements (@idea)
> 2. Create workstream breakdown (@design)
> 3. Guide you through execution (@build)
>
> Alternatively, you can use individual skills:
> - `@idea` → requirements gathering
> - `@design` → workstream planning
> - `@build` → execute one workstream

## Output

Always provide:
1. Recommended skill with command
2. Brief explanation of why
3. Example usage
4. Related alternatives

## See Also

- [.claude/skills/](../../.claude/skills/) - All skill definitions
- [CLAUDE.md](../../CLAUDE.md) - Claude Code integration guide
