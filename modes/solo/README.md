# Solo Mode

The simplest approach for AI-assisted development. One agent, one session, iterative dialogue.

## When to Use

- Task takes < 2 hours
- Changes < 10 files
- No major architectural decisions needed
- Bug fixes, small features, refactoring

## How It Works

```
You: "Fix the login bug"
    ↓
AI: Analyzes → Finds issue → Fixes → Tests
    ↓
You: "Looks good" or "Also check X"
    ↓
AI: Iterates until done
```

**Key insight**: A capable AI model can handle analysis, design, implementation, and testing in one conversation. You guide through dialogue, not process.

## Setup

1. Create `CLAUDE.md` in your project root (see [example](CLAUDE.md.example))
2. Start Claude Code or your preferred AI tool
3. Describe your task
4. Iterate until done

## CLAUDE.md for Solo Mode

Keep it short (50-100 lines). Include:
- Tech stack
- Key commands (test, build, run)
- Important rules
- Project-specific conventions

See [CLAUDE.md.example](CLAUDE.md.example) for a template.

## Example Prompts

See [prompts/examples.md](prompts/examples.md) for copy-paste prompts.

### Bug Fix
```
Fix bug: users get 500 error when resetting password.
Look at logs/error.log for stack trace.
Find root cause, fix it, add regression test.
```

### Small Feature
```
Add rate limiting to POST /api/login.
- Max 5 attempts per minute per IP
- Return 429 Too Many Requests when exceeded
- Add tests
```

### Refactoring
```
Refactor UserService to use repository pattern.
- Extract database calls to UserRepository
- Keep UserService focused on business logic
- All existing tests must pass
```

### Code Review
```
Review the changes in src/auth/ for:
- Security issues
- Error handling
- Test coverage
- Clean code violations

Report findings with file:line references.
```

## Tips

### Be Specific
```
# Bad
"Fix the bug"

# Good
"Fix bug in src/auth/login.py where users with special
characters in email can't login. Error: InvalidEmailFormat"
```

### Provide Context
```
# Bad
"Add caching"

# Good
"Add Redis caching to getUserById in UserService.
Cache for 5 minutes. Invalidate on user update.
We already have Redis configured in src/config/redis.py"
```

### Iterate
Don't try to get everything perfect in one prompt. Start simple, refine:

```
1. "Implement basic user registration"
2. "Add email validation"
3. "Add password strength requirements"
4. "Add rate limiting"
```

## When to Upgrade to Structured Mode

Switch to [Structured Mode](../structured/) when:
- Task grows beyond 2 hours
- You need documentation for the team
- Making architectural decisions that should be recorded
- Multiple components involved

## Anti-patterns

| Don't | Do |
|-------|-----|
| "Write the whole feature" | Break into smaller tasks |
| Skip testing | "Add tests for the changes" |
| Ignore errors | "What does this error mean?" |
| Accept first solution | "Are there other approaches?" |

## Model Selection

For Solo mode, use a capable model that can handle all aspects:

| Task Type | Recommended |
|-----------|-------------|
| Complex debugging | Most capable (opus) |
| Standard features | Balanced (sonnet) |
| Simple fixes | Fast (haiku) |

In Claude Code: `/model opus`, `/model sonnet`, `/model haiku`
