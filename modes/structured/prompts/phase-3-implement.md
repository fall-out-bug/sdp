# Phase 3: Implement

You are implementing the designed solution.

## Your Task

Write code and tests following the specification and design.

## Input

- Specification: `docs/specs/{feature}.md`
- Design: `docs/specs/{feature}-design.md`
- Any ADRs in `docs/adr/`

## Output

- Production code in `src/`
- Tests in `tests/`
- Updated documentation if needed

## Implementation Order

Follow Clean Architecture layers, inside-out:

### 1. Domain Layer
Entities, value objects, domain services:
```
src/domain/
├── entities/
│   └── password_reset_token.py
└── services/
    └── password_policy.py
```

### 2. Application Layer
Use cases, application services, ports:
```
src/application/
├── services/
│   └── password_reset_service.py
└── ports/
    └── token_repository.py  # interface
```

### 3. Infrastructure Layer
Implementations of ports, external integrations:
```
src/infrastructure/
├── repositories/
│   └── sql_token_repository.py
└── services/
    └── smtp_email_service.py
```

### 4. Presentation Layer
Controllers, API routes:
```
src/presentation/
└── controllers/
    └── password_reset_controller.py
```

## TDD Workflow

For each component:

1. **Write test first**
```python
def test_password_reset_token_expires_after_one_hour():
    token = PasswordResetToken.create(user_id="123")
    assert token.is_expired() == False

    # Simulate time passing
    token.created_at = datetime.now() - timedelta(hours=2)
    assert token.is_expired() == True
```

2. **Write minimal code to pass**
```python
class PasswordResetToken:
    def is_expired(self) -> bool:
        return datetime.now() > self.created_at + timedelta(hours=1)
```

3. **Refactor if needed**

4. **Move to next component**

## Example Prompt

```
Implement password reset feature based on:
- Spec: docs/specs/password-reset.md
- Design: docs/specs/password-reset-design.md

Start with domain layer:
1. Create PasswordResetToken entity
2. Add validation for token expiry
3. Write tests first

Then proceed to application layer.
```

## Quality Checklist

Before finishing this phase:
- [ ] All acceptance criteria from spec are implemented
- [ ] Design document was followed (or updated if deviated)
- [ ] Tests exist for all new code
- [ ] Test coverage ≥80% for new code
- [ ] No hardcoded values (use config)
- [ ] Error handling is explicit (no silent failures)
- [ ] Code follows existing patterns in codebase

## Tips

### Follow the Design
The design exists for a reason. If you need to deviate:
1. Update the design document
2. Explain why in commit message

### Small Commits
Commit after each component:
```
git commit -m "Add PasswordResetToken entity with expiry logic"
git commit -m "Add PasswordResetService application service"
git commit -m "Add SqlTokenRepository infrastructure"
git commit -m "Add password reset API endpoints"
```

### Don't Over-Engineer
Implement what's in the spec, not what might be needed later.

### Ask Questions
If spec or design is unclear, ask before guessing.
