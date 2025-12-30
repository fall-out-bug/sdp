# Workstreams and Guardrails Guide

This guide explains how to break epics into workstreams and what guardrails protect quality.

## Workstreams

### What is a Workstream?

A workstream is a cohesive unit of work within an epic that can be:
- Implemented independently
- Code reviewed incrementally
- Completed in a reasonable timeframe (typically 2-4 hours)

### Breaking Epics into Workstreams

The Tech Lead creates workstreams in `implementation.md`. Good workstreams:

1. **Have clear boundaries** - Defined inputs and outputs
2. **Are independently testable** - Can verify without other workstreams
3. **Follow Clean Architecture** - Stay within one or two layers
4. **Are appropriately sized** - Not too large, not too granular

**Example**: Adding a user authentication feature

```
Workstream 1: Domain layer (User entity, AuthToken value object)
Workstream 2: Application layer (AuthService, use cases)
Workstream 3: Infrastructure layer (TokenRepository, PasswordHasher)
Workstream 4: Presentation layer (LoginController, middleware)
```

### Continuous Code Review Cycle (v1.2)

Each workstream follows this cycle:

```
Developer implements workstream
    ↓
Developer conducts incremental code review
    ↓
Developer fixes violations
    ↓
Tech Lead reviews code quality
    ↓
Tech Lead approves OR vetoes
    ↓
Next workstream (or fix violations if vetoed)
```

**Key rule**: No workstream proceeds until the previous one passes review.

### Workstream Artifacts

For each workstream, Developer should:

1. **Before implementing**: Search codebase for existing implementations
2. **During implementation**: Follow TDD (test first)
3. **After implementation**: Document in workstream completion message
4. **Code review checklist**:
   - DRY (Don't Repeat Yourself)
   - SOLID principles
   - Clean Code practices
   - Clean Architecture boundaries

## Guardrails

Guardrails are quality gates that cannot be bypassed. They exist at multiple levels.

### 1. Architecture Guardrails

**Clean Architecture boundaries**:
```
Presentation → Infrastructure → Application → Domain
```

Dependencies MUST point inward. Violations trigger automatic veto.

**Layer responsibilities**:
- **Domain**: Business entities, value objects, domain services
- **Application**: Use cases, application services, ports/interfaces
- **Infrastructure**: Database, external APIs, frameworks
- **Presentation**: Controllers, views, API endpoints

### 2. Code Quality Guardrails

**Forbidden patterns** (automatic veto):
```python
# FORBIDDEN: Silent exception handling
except: pass
except Exception: pass

# FORBIDDEN: Default values hiding errors
def get_user(id, default=None):
    try:
        return find_user(id)
    except:
        return default  # Hides the error!

# FORBIDDEN: Catch-all hiding errors
try:
    complex_operation()
except Exception as e:
    logger.debug(e)  # Debug-level logging hides errors
```

**Required patterns**:
```python
# CORRECT: Explicit error handling
try:
    return find_user(id)
except UserNotFound as e:
    logger.error(f"User not found: {id}")
    raise

# CORRECT: Proper error propagation
def get_user(id):
    user = repository.find(id)
    if user is None:
        raise UserNotFound(id)
    return user
```

### 3. Veto Protocol

Certain conditions trigger vetoes that **cannot be overridden**:

| Veto Trigger | Agent | Cannot Override |
|--------------|-------|-----------------|
| Architecture violations | Architect | Yes |
| Security issues | Security | Yes |
| Missing rollback plan | DevOps | Yes |
| Code review violations | Tech Lead, QA | Yes |
| Layer violations | Architect | Yes |

**Negotiable issues** (can be resolved through discussion):
- Scope creep (Analyst)
- Timeline concerns (Tech Lead)
- Non-critical test coverage

### 4. Testing Guardrails

**Coverage requirements**:
- Minimum 80% coverage in touched areas
- All new code must have tests
- Integration tests for cross-layer interactions

**Test quality gates**:
- Unit tests must be isolated
- No tests depending on external services
- Deterministic results (no flaky tests)

### 5. Documentation Guardrails

**All output must be in English** - This is non-negotiable.

**Required documentation**:
- `requirements.json` - Analyst output
- `architecture.json` - Architect output
- `implementation.md` - Tech Lead output
- `code_review.md` - Developer output at epic completion
- `test_results.md` - QA output

### 6. Message Format Guardrails

JSON messages use compact keys:
- `d` - date
- `st` - status
- `r` - role
- `epic` - epic ID
- `sm` - summary
- `nx` - next steps
- `artifacts` - artifact paths

**File naming**: `{YYYY-MM-DD}-{subject}.json`

## Guardrail Enforcement

### How Guardrails Are Enforced

1. **Agent prompts** include guardrail checks
2. **Code review** verifies compliance
3. **Veto protocol** blocks non-compliant work
4. **Self-verification checklist** before each completion

### Self-Verification Checklist

Before completing any work:
- [ ] Clean Architecture boundaries respected
- [ ] No forbidden patterns (except:pass, etc.)
- [ ] Engineering principles followed (DRY, SOLID)
- [ ] Documentation updated
- [ ] All messages in English
- [ ] Tests pass with adequate coverage

## Recovery from Guardrail Violations

If work is vetoed:

1. **Read the veto message** - Understand the specific violation
2. **Fix the violation** - Don't just patch around it
3. **Re-submit for review** - Same agent reviews again
4. **Document the fix** - Add to decision log

See [CONTROL_PLAYBOOK.md](../../CONTROL_PLAYBOOK.md) for complex recovery scenarios.

## Summary

| Concept | Purpose | Enforcement |
|---------|---------|-------------|
| Workstreams | Manageable units of work | Tech Lead defines, continuous review |
| Architecture guardrails | Maintain Clean Architecture | Architect veto |
| Code quality guardrails | Prevent bad patterns | Code review, Tech Lead veto |
| Testing guardrails | Ensure coverage | QA verification |
| Veto protocol | Block non-compliant work | Cannot override |

---

**Remember**: Guardrails exist to prevent problems, not slow you down. Following them from the start is faster than fixing violations later.
