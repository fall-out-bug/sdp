# Project Conventions

**Fill this template when integrating SDP into your project.**

This file defines project-specific conventions that AI agents must follow.

---

## Language & Communication

**Primary Language:** [e.g., English, Russian, Spanish]

**Code Comments:** [e.g., English only]

**Documentation:** [e.g., English for technical docs, bilingual for user guides]

---

## Code Style

### Python

**Formatter:** [e.g., black, ruff]

**Linter:** [e.g., flake8, ruff]

**Type Checker:** [e.g., mypy --strict]

**Line Length:** [e.g., 88, 100, 120]

**Import Style:**
```python
# ✅ DO
from project.domain.models import User
from project.application.ports import UserRepository

# ❌ DON'T
from project.domain.models import *
from project.infrastructure import UserRepositoryImpl
```

### JavaScript/TypeScript

**Formatter:** [e.g., prettier]

**Linter:** [e.g., eslint]

**Style Guide:** [e.g., Airbnb, Standard]

---

## Architecture

### Layer Boundaries

**DO:**
- ✅ Domain layer is pure business logic, no external dependencies
- ✅ Application layer defines ports (interfaces)
- ✅ Infrastructure implements ports
- ✅ Presentation uses only Application layer

**DON'T:**
- ❌ Domain importing Infrastructure
- ❌ Domain importing Presentation
- ❌ Application importing Infrastructure directly
- ❌ Business logic in controllers

### Example (fill with your project structure)

```
project/
├── domain/           # Pure business logic
├── application/      # Use cases & ports
├── infrastructure/   # DB, API, external services
└── presentation/     # API, CLI, UI
```

---

## Naming Conventions

### Files

**DO:**
```
✅ user_repository.py
✅ user-repository.ts
✅ UserRepository.java
```

**DON'T:**
```
❌ user_repo.py (unclear abbreviation)
❌ userrepository.py (no separator)
❌ User_Repository.py (inconsistent case)
```

### Variables

**DO:**
```python
✅ user_count = 10
✅ max_retry_attempts = 3
✅ is_authenticated = True
```

**DON'T:**
```python
❌ uc = 10  # unclear
❌ maxRetryAttempts = 3  # inconsistent case
❌ auth = True  # unclear abbreviation
```

### Functions/Methods

**DO:**
```python
✅ def get_user_by_id(user_id: int) -> User:
✅ def validate_email(email: str) -> bool:
✅ def calculate_total_price(items: list[Item]) -> Decimal:
```

**DON'T:**
```python
❌ def get_user(id):  # missing type hints
❌ def validate(e):  # unclear parameter
❌ def calc(items):  # unclear abbreviation
```

### Classes

**DO:**
```python
✅ class UserRepository
✅ class EmailValidator
✅ class PaymentProcessor
```

**DON'T:**
```python
❌ class UserRepo  # unclear abbreviation
❌ class user_repository  # wrong case
❌ class UserMgr  # unclear abbreviation
```

---

## Testing

### Coverage Requirements

**Minimum Coverage:** [e.g., 80%]

**Critical Paths Coverage:** [e.g., 100%]

### Test Structure

**DO:**
```python
# ✅ Clear test names
def test_user_creation_with_valid_email():
def test_authentication_fails_with_invalid_password():
def test_payment_processing_with_insufficient_funds():
```

**DON'T:**
```python
# ❌ Unclear test names
def test_user():
def test_auth():
def test_payment():
```

### Mocking

**DO:**
- ✅ Mock external services (APIs, databases)
- ✅ Use test doubles for adapters
- ✅ Test domain logic without mocks

**DON'T:**
- ❌ Mock domain entities
- ❌ Over-mock (testing mocks, not logic)
- ❌ Skip integration tests entirely

---

## Error Handling

**DO:**
```python
# ✅ Explicit error handling
try:
    user = get_user(user_id)
except UserNotFoundError as e:
    logger.error(f"User not found: {user_id}", exc_info=e)
    raise HTTPException(status_code=404, detail="User not found")
except DatabaseError as e:
    logger.error(f"Database error: {e}", exc_info=e)
    raise HTTPException(status_code=500, detail="Internal error")
```

**DON'T:**
```python
# ❌ Silent failures
try:
    user = get_user(user_id)
except:
    pass

# ❌ Bare exceptions
try:
    user = get_user(user_id)
except Exception:
    return None

# ❌ Default values hiding errors
def get_user(user_id):
    try:
        return db.query(User).get(user_id)
    except:
        return User(id=0, name="Unknown")  # Hiding error!
```

---

## Git Workflow

### Branch Naming

**DO:**
```
✅ feature/user-authentication
✅ fix/login-redirect
✅ hotfix/critical-api-failure
```

**DON'T:**
```
❌ feature-user-auth (missing category prefix)
❌ fix_login (inconsistent separator)
❌ johns-branch (unclear purpose)
```

### Commit Messages

**Format:** [e.g., Conventional Commits]

**DO:**
```
✅ feat(auth): WS-001-01 - add user registration
✅ fix(api): WS-002-03 - handle timeout in payment API
✅ test(domain): WS-001-02 - add validation tests
```

**DON'T:**
```
❌ added stuff
❌ fixed bug
❌ WIP
```

---

## Documentation

### Docstrings

**Style:** [e.g., Google, NumPy, Sphinx]

**DO:**
```python
def process_payment(amount: Decimal, user_id: int) -> PaymentResult:
    """Process payment for user.
    
    Args:
        amount: Payment amount in USD
        user_id: User identifier
        
    Returns:
        PaymentResult with status and transaction_id
        
    Raises:
        InsufficientFundsError: If user balance < amount
        PaymentGatewayError: If external API fails
        
    Example:
        result = process_payment(Decimal("19.99"), user_id=123)
    """
```

**DON'T:**
```python
def process_payment(amount, user_id):
    # Process payment
    pass
```

### README Sections

**Required:**
- [ ] Project description
- [ ] Installation instructions
- [ ] Quick start example
- [ ] Development setup
- [ ] Testing instructions
- [ ] Deployment guide

---

## Dependencies

### Adding Dependencies

**DO:**
- ✅ Use latest stable versions
- ✅ Pin major versions in requirements.txt / package.json
- ✅ Document why dependency is needed
- ✅ Check license compatibility

**DON'T:**
- ❌ Use outdated packages with vulnerabilities
- ❌ Add dependencies without review
- ❌ Use exact pins without updates (security risk)

---

## Security

### Secrets Management

**DO:**
- ✅ Use environment variables
- ✅ Use secret management tools (Vault, AWS Secrets Manager)
- ✅ Never commit .env files
- ✅ Rotate secrets regularly

**DON'T:**
- ❌ Hardcode API keys
- ❌ Commit secrets to git
- ❌ Use default passwords
- ❌ Share secrets in Slack/email

---

## Performance

### Database Queries

**DO:**
- ✅ Use indexes on frequently queried columns
- ✅ Add timeouts to all queries
- ✅ Use connection pooling
- ✅ Paginate large result sets

**DON'T:**
- ❌ N+1 queries
- ❌ SELECT * without filters
- ❌ Missing indexes on foreign keys
- ❌ Long-running queries without timeouts

---

## AI-Specific Rules

### File Size Limits

**Max LOC per file:** [e.g., 200]

**Action if exceeded:** [e.g., Split into multiple modules]

### Complexity Limits

**Max Cyclomatic Complexity:** [e.g., 10]

**Action if exceeded:** [e.g., Refactor into smaller functions]

---

## Project-Specific DO/DON'T

### [Your Domain Area 1]

**DO:**
- ✅ [Add your domain-specific best practice]
- ✅ [Add another best practice]

**DON'T:**
- ❌ [Add your domain-specific anti-pattern]
- ❌ [Add another anti-pattern]

### [Your Domain Area 2]

**DO:**
- ✅ [Add best practice]

**DON'T:**
- ❌ [Add anti-pattern]

---

## Review Checklist

Before marking any workstream as DONE, verify:

- [ ] All DO/DON'T rules followed
- [ ] Tests pass (coverage ≥ [X]%)
- [ ] Type hints complete
- [ ] Documentation updated
- [ ] No TODO/FIXME in code
- [ ] Clean Architecture followed
- [ ] No secrets in code
- [ ] Performance targets met
- [ ] Security review passed

---

**Last Updated:** [YYYY-MM-DD]  
**Version:** 1.0  
**Maintained By:** [Team/Person]
