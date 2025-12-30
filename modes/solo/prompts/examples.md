# Solo Mode Prompt Examples

Copy-paste prompts for common tasks. Customize for your project.

## Bug Fixes

### Simple Bug
```
Fix bug: [describe the bug]
Location: [file path if known]
Expected: [what should happen]
Actual: [what happens now]
```

### Bug with Error Log
```
Fix this error:

[paste error/stack trace]

The error occurs when [describe trigger].
Find root cause and fix.
```

### Regression Bug
```
Fix regression: [feature] stopped working after [recent change].
It worked before commit [hash or date].
Compare the changes and restore correct behavior.
```

## Features

### Add Endpoint
```
Add API endpoint: [METHOD] /api/[path]

Request:
- [param1]: [type] - [description]
- [param2]: [type] - [description]

Response:
- Success (200): [describe response]
- Error (400): [describe error cases]

Add validation, error handling, and tests.
```

### Add Validation
```
Add validation to [function/endpoint]:
- [rule 1]
- [rule 2]
- [rule 3]

Return appropriate error messages.
Add tests for each validation rule.
```

### Add Configuration
```
Make [value] configurable:
- Environment variable: [VAR_NAME]
- Default value: [default]
- Add to config documentation

Update all usages of hardcoded value.
```

## Refactoring

### Extract Function
```
Extract the [logic description] from [file:function]
into a separate function.
- Keep it in same file or move to [target file]
- Update all callers
- Ensure tests pass
```

### Apply Pattern
```
Refactor [class/module] to use [pattern name] pattern.
- Current: [describe current structure]
- Target: [describe desired structure]
- Keep existing functionality working
```

### Improve Error Handling
```
Improve error handling in [file/module]:
- Replace bare except with specific exceptions
- Add proper logging for errors
- Ensure errors propagate correctly
- Add tests for error cases
```

## Testing

### Add Unit Tests
```
Add unit tests for [class/function]:
- Test happy path
- Test edge cases: [list edge cases]
- Test error cases: [list error cases]
- Mock external dependencies
```

### Add Integration Test
```
Add integration test for [feature]:
- Setup: [describe test data needed]
- Action: [what to test]
- Verify: [expected outcomes]
```

### Fix Flaky Test
```
Fix flaky test: [test name]
It fails intermittently with: [error]
Likely cause: [your hypothesis]
Make it deterministic.
```

## Code Review

### Review Changes
```
Review changes in [files/directory]:
1. Security issues
2. Error handling
3. Test coverage
4. Performance concerns
5. Clean code violations

Report each issue with file:line reference.
```

### Review for Production
```
Review [feature] for production readiness:
- Error handling complete?
- Logging sufficient?
- Configuration externalized?
- Tests adequate?
- Documentation updated?

List any blockers.
```

## Documentation

### Add Docstrings
```
Add docstrings to public functions in [file]:
- Description of purpose
- Args with types and descriptions
- Returns description
- Raises (if applicable)

Use [Google/NumPy/Sphinx] style.
```

### Update README
```
Update README.md to include:
- [new feature] usage instructions
- Any new configuration options
- Updated examples
```

## Performance

### Optimize Query
```
Optimize slow query in [file:function]:
Current: [describe current approach]
Problem: [what makes it slow]
Target: [performance goal if known]

Show before/after and explain the improvement.
```

### Add Caching
```
Add caching to [function]:
- Cache key: [what to use as key]
- TTL: [duration]
- Invalidation: [when to clear cache]

We use [Redis/in-memory/etc].
```

## Tips for Better Prompts

### 1. Be Specific
```
# Vague
"Make it faster"

# Specific
"Optimize getUserById to use a single query instead of N+1.
Currently takes 500ms for user with 100 orders."
```

### 2. Provide Context
```
# Missing context
"Add authentication"

# With context
"Add JWT authentication to /api/admin/* endpoints.
We already have User model with password_hash field.
Use existing src/utils/jwt.py for token generation."
```

### 3. Define Success
```
# Unclear success
"Fix the performance"

# Clear success
"Reduce /api/search response time from 2s to <200ms.
Measure with: time curl localhost:8000/api/search?q=test"
```
