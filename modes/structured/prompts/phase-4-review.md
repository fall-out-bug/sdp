# Phase 4: Review

You are reviewing the implementation for quality and completeness.

## Your Task

Verify the implementation meets all requirements and quality standards.

## Input

- Specification: `docs/specs/{feature}.md`
- Design: `docs/specs/{feature}-design.md`
- Implementation: `src/` and `tests/`

## Output

Either:
- Approval (all checks pass)
- Issue report: `docs/reviews/{feature}-review.md`

## Review Checklist

### 1. Acceptance Criteria (from spec)

Go through each criterion:
```markdown
## Acceptance Criteria Check

- [x] User can request password reset via email
- [x] Reset link expires after 1 hour
- [ ] User receives email within 5 minutes ← NOT TESTED
- [x] New password must meet strength requirements
```

### 2. Design Compliance

Verify implementation matches design:
```markdown
## Design Compliance

- [x] Components created as designed
- [x] Data flow follows design
- [x] API matches specification
- [ ] Database schema differs from design ← DEVIATION
  - Design: `used_at TIMESTAMP`
  - Actual: `is_used BOOLEAN`
  - Impact: Minor, acceptable
```

### 3. Test Coverage

Run tests and check coverage:
```bash
pytest --cov=src --cov-report=term-missing
```

```markdown
## Test Coverage

Overall: 84% ✓
New code coverage: 91% ✓

Uncovered areas:
- src/infrastructure/email_service.py:45-52 (error handling)
```

### 4. Code Quality

Check for issues:

```markdown
## Code Quality

### Clean Architecture
- [x] Domain has no external dependencies
- [x] Application depends only on domain
- [x] Infrastructure implements application ports
- [x] Presentation depends on application

### Error Handling
- [x] No bare `except:` clauses
- [x] Errors logged appropriately
- [x] User-facing errors are friendly

### Security
- [x] Tokens are hashed before storage
- [x] Rate limiting implemented
- [ ] No timing attack protection ← ISSUE

### Code Smells
- [ ] Long method in PasswordResetService.initiate() (45 lines)
- [ ] Duplicate validation logic in two controllers
```

## Example Prompt

```
Review the password reset implementation:

1. Check all acceptance criteria from docs/specs/password-reset.md
2. Verify design from docs/specs/password-reset-design.md was followed
3. Run tests: pytest --cov=src
4. Check for:
   - Security issues
   - Error handling
   - Code smells
   - Clean Architecture violations

Report issues with file:line references.
```

## Review Report Format

If issues found, create `docs/reviews/{feature}-review.md`:

```markdown
# Review: Password Reset Feature

Date: 2024-01-15
Reviewer: AI Assistant

## Summary
Implementation is 90% complete. 2 issues need fixing before merge.

## Blocking Issues

### 1. Missing Timing Attack Protection
**Location**: src/application/services/password_reset_service.py:34
**Issue**: Token comparison uses `==` which is vulnerable to timing attacks
**Fix**: Use `secrets.compare_digest()` for constant-time comparison

### 2. Uncovered Error Path
**Location**: src/infrastructure/email_service.py:45-52
**Issue**: SMTP errors not tested
**Fix**: Add test for email sending failure scenario

## Non-Blocking Issues

### 1. Long Method
**Location**: src/application/services/password_reset_service.py:15-60
**Issue**: Method is 45 lines, hard to read
**Suggestion**: Extract email composition to separate method

## Test Results
- Total: 42 passed, 0 failed
- Coverage: 84%
- New code coverage: 91%

## Recommendation
Fix blocking issues, then approve for merge.
```

## Quality Gates

### Blocking (must fix)
- Any acceptance criterion not met
- Security vulnerabilities
- Test coverage < 80% for new code
- Silent error handling (`except: pass`)
- Clean Architecture violations

### Non-Blocking (should fix)
- Code smells
- Minor deviations from design
- Documentation gaps
- Style inconsistencies

## Tips

- Be specific with file:line references
- Explain WHY something is an issue
- Suggest fixes, not just problems
- Distinguish blocking vs non-blocking
- Check edge cases in tests
