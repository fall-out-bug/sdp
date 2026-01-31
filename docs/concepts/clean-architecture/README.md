# Clean Architecture in SDP

Clean Architecture is a software design principle that separates concerns into layers with strict dependency rules.

## SDP Architecture

SDP implements Clean Architecture with the following layers:

```
                    ┌─────────────┐
                    │   domain/   │  ← Pure entities, no deps
                    └─────────────┘
                          ↑
            ┌─────────────┼─────────────┐
            │             │             │
      ┌─────────┐   ┌─────────┐   ┌─────────┐
      │  core/  │   │ beads/  │   │ unified/│
      └─────────┘   └─────────┘   └─────────┘
            ↑             ↑             ↑
            │             │             │
            └─────────────┼─────────────┘
                          │
                    ┌─────────┐
                    │   cli/  │
                    └─────────┘
```

### Layer Responsibilities

| Layer | Purpose | Allowed Dependencies |
|-------|---------|---------------------|
| **domain/** | Pure business entities (Workstream, Feature, WorkstreamID) | None (pure Python only) |
| **core/** | Application logic (parsing, validation, orchestration) | domain/ |
| **beads/** | Infrastructure layer (external task storage) | domain/, core/ |
| **unified/** | Infrastructure layer (AI agent integration) | domain/, core/ |
| **cli/** | Presentation layer (CLI commands) | domain/, core/, beads/, unified/ |

### SDP-Specific Entities

#### Domain Layer (`src/sdp/domain/`)

**Workstream Entity** (`workstream.py`):
```python
@dataclass
class Workstream:
    """Core workstream entity with business logic."""
    ws_id: str
    feature: str
    status: WorkstreamStatus
    size: WorkstreamSize
    # ... no I/O logic or external dependencies
```

**Feature Aggregate** (`feature.py`):
```python
@dataclass
class Feature:
    """Feature aggregate managing workstream dependencies."""
    feature_id: str
    workstreams: list[Workstream]
    dependency_graph: dict[str, list[str]]
    execution_order: list[str]  # Topologically sorted
```

**Value Objects** (`workstream.py`):
```python
@dataclass(frozen=True)
class WorkstreamID:
    """Immutable workstream identifier (PP-FFF-SS format)."""
    project_id: int
    feature_id: int
    sequence: int
```

**Domain Exceptions** (`exceptions.py`):
```python
class DomainError(Exception): pass
class ValidationError(DomainError): pass
class DependencyCycleError(DomainError): pass
class MissingDependencyError(DomainError): pass
```

#### Core Layer (`src/sdp/core/`)

**Application Services**:
- `workstream/parser.py` - Parse markdown files → Workstream entities
- `feature/models.py` - Feature dependency graph management
- `feature/loader.py` - Load features from directories

**Business Logic**:
- Dependency validation (cycle detection)
- Topological sorting (execution order)
- Status transitions

#### Infrastructure Layers

**Beads** (`src/sdp/beads/`):
- External task storage (JSON files)
- Skill orchestration
- Sync with domain entities

**Unified** (`src/sdp/unified/`):
- AI agent integration
- Multi-agent orchestration
- Tool abstractions

#### Presentation Layer (`src/sdp/cli/`)

**CLI Commands**:
- `sdp build WS-001-01` - Execute workstream
- `sdp review F01` - Review feature
- `sdp guard activate WS-001-01` - Activate guard

## Architecture Rules

### ✅ Allowed

```python
# domain/ → (nothing)
# core/ → domain/
from sdp.domain.workstream import Workstream, WorkstreamStatus

# beads/ → domain/, core/
from sdp.domain.workstream import Workstream
from sdp.core.workstream import parse_workstream

# cli/ → domain/, core/, beads/, unified/
from sdp.domain.workstream import WorkstreamID
from sdp.core.workstream import parse_workstream
from sdp.beads.sync_service import BeadsSync
```

### ❌ Forbidden

```python
# domain/ → ANYTHING
from sdp.core.workstream import parse_workstream  # ❌ External dependency

# beads/ → core/workstream/models
from sdp.core.workstream.models import Workstream  # ❌ Use domain/ instead

# unified/ → core/feature/models
from sdp.core.feature.models import Feature  # ❌ Use domain/ instead
```

## Checking Architecture

### Manual Check

```bash
python scripts/check_architecture.py
```

Output:
```
Checking src/sdp/domain/...
Checking src/sdp/beads/...
Checking src/sdp/unified/...

✅ All architecture checks passed!
```

### Pre-commit Hook

Add to `.git/hooks/pre-commit`:
```bash
#!/bin/bash
python scripts/check_architecture.py || exit 1
```

### CI Integration

```yaml
# .github/workflows/architecture.yml
- name: Check architecture
  run: python scripts/check_architecture.py
```

## Migration from Old Structure

### Before (Violation)

```python
# beads/sync_service.py
from sdp.core.workstream import Workstream  # ❌ Infrastructure → Application
```

**Problem**: `beads/` (infrastructure) directly imports from `core/` (application), creating tight coupling.

### After (Fixed)

```python
# beads/sync_service.py
from sdp.domain.workstream import Workstream  # ✅ Infrastructure → Domain
```

**Solution**: Both layers depend on shared domain entities.

### Backward Compatibility

Old imports still work (with deprecation warnings):

```python
# OLD: core/workstream/models.py
from sdp.core.workstream.models import Workstream
# DeprecationWarning: Use 'from sdp.domain.workstream import Workstream'

# NEW: domain/workstream.py
from sdp.domain.workstream import Workstream  # ✅ No warning
```

```
┌─────────────────────────────────────────────────────┐
│                  Presentation                        │
│              (Controllers, Views, API)               │
├─────────────────────────────────────────────────────┤
│                  Infrastructure                      │
│          (Database, External APIs, Frameworks)       │
├─────────────────────────────────────────────────────┤
│                   Application                        │
│              (Use Cases, Services)                   │
├─────────────────────────────────────────────────────┤
│                     Domain                           │
│           (Entities, Business Rules)                 │
└─────────────────────────────────────────────────────┘

        ↑ Dependencies point INWARD (toward Domain)
```

## The Layers

### Domain (Innermost)

**Contains**: Business entities, value objects, domain services, business rules.

**Knows about**: Nothing external. Pure business logic.

**Example**:
```python
# domain/entities/user.py
class User:
    def __init__(self, id: str, email: str, password_hash: str):
        self.id = id
        self.email = email
        self.password_hash = password_hash

    def can_login(self) -> bool:
        return self.is_active and not self.is_locked

# domain/value_objects/email.py
class Email:
    def __init__(self, value: str):
        if not self._is_valid(value):
            raise InvalidEmailError(value)
        self.value = value

    def _is_valid(self, value: str) -> bool:
        return "@" in value and "." in value
```

**Rules**:
- No imports from other layers
- No framework dependencies
- No database code
- Pure Python/language features only

### Application

**Contains**: Use cases, application services, ports (interfaces).

**Knows about**: Domain layer only.

**Example**:
```python
# application/services/auth_service.py
from domain.entities.user import User
from application.ports.user_repository import UserRepository  # interface

class AuthService:
    def __init__(self, user_repo: UserRepository):
        self.user_repo = user_repo

    def login(self, email: str, password: str) -> User:
        user = self.user_repo.find_by_email(email)
        if not user or not user.verify_password(password):
            raise InvalidCredentialsError()
        if not user.can_login():
            raise AccountLockedError()
        return user

# application/ports/user_repository.py
from abc import ABC, abstractmethod
from domain.entities.user import User

class UserRepository(ABC):
    @abstractmethod
    def find_by_email(self, email: str) -> User | None:
        pass

    @abstractmethod
    def save(self, user: User) -> None:
        pass
```

**Rules**:
- Imports from Domain allowed
- No imports from Infrastructure or Presentation
- Defines interfaces (ports) that Infrastructure implements

### Infrastructure

**Contains**: Database implementations, external API clients, framework integrations.

**Knows about**: Domain and Application layers.

**Example**:
```python
# infrastructure/repositories/sql_user_repository.py
from sqlalchemy.orm import Session
from domain.entities.user import User
from application.ports.user_repository import UserRepository

class SqlUserRepository(UserRepository):
    def __init__(self, session: Session):
        self.session = session

    def find_by_email(self, email: str) -> User | None:
        row = self.session.query(UserModel).filter_by(email=email).first()
        return self._to_entity(row) if row else None

    def save(self, user: User) -> None:
        model = self._to_model(user)
        self.session.add(model)
        self.session.commit()

# infrastructure/external/email_service.py
import smtplib
from application.ports.email_sender import EmailSender

class SmtpEmailService(EmailSender):
    def send(self, to: str, subject: str, body: str) -> None:
        # SMTP implementation
        pass
```

**Rules**:
- Implements interfaces from Application layer
- Can import from Domain and Application
- No imports from Presentation

### Presentation (Outermost)

**Contains**: Controllers, API endpoints, views, CLI commands.

**Knows about**: Application layer (and transitively Domain).

**Example**:
```python
# presentation/controllers/auth_controller.py
from fastapi import APIRouter, Depends, HTTPException
from application.services.auth_service import AuthService

router = APIRouter()

@router.post("/login")
def login(email: str, password: str, auth_service: AuthService = Depends()):
    try:
        user = auth_service.login(email, password)
        return {"token": create_token(user)}
    except InvalidCredentialsError:
        raise HTTPException(401, "Invalid credentials")
    except AccountLockedError:
        raise HTTPException(403, "Account locked")
```

**Rules**:
- Can import from Application
- Should not import directly from Infrastructure
- Handles HTTP, CLI, or UI concerns

## Dependency Injection

How do layers connect without violating the dependency rule?

**Dependency Injection** at the composition root (usually in main.py or a DI container):

```python
# main.py (composition root)
from infrastructure.repositories.sql_user_repository import SqlUserRepository
from application.services.auth_service import AuthService
from presentation.controllers.auth_controller import router

# Create instances with dependencies injected
db_session = create_session()
user_repo = SqlUserRepository(db_session)  # Infrastructure
auth_service = AuthService(user_repo)       # Application uses interface

# Wire up presentation
app.include_router(router)
```

## Common Violations

### Domain Importing Infrastructure

```python
# BAD: domain/entities/user.py
from sqlalchemy import Column, String  # Framework in domain!

class User(Base):  # Domain entity inherits from ORM
    __tablename__ = "users"
```

**Fix**: Keep domain entities pure, create separate ORM models in infrastructure.

### Application Importing Infrastructure

```python
# BAD: application/services/auth_service.py
from infrastructure.repositories.sql_user_repository import SqlUserRepository

class AuthService:
    def __init__(self):
        self.repo = SqlUserRepository()  # Direct infrastructure dependency
```

**Fix**: Depend on interface (port), inject implementation.

### Presentation Bypassing Application

```python
# BAD: presentation/controllers/user_controller.py
from infrastructure.repositories.sql_user_repository import SqlUserRepository

@router.get("/users/{id}")
def get_user(id: str):
    repo = SqlUserRepository()  # Skipping application layer
    return repo.find_by_id(id)
```

**Fix**: Go through application services.

## Directory Structure

```
src/
├── domain/
│   ├── entities/
│   │   ├── user.py
│   │   └── order.py
│   ├── value_objects/
│   │   ├── email.py
│   │   └── money.py
│   └── services/
│       └── pricing_service.py
│
├── application/
│   ├── services/
│   │   ├── auth_service.py
│   │   └── order_service.py
│   ├── ports/
│   │   ├── user_repository.py
│   │   └── email_sender.py
│   └── use_cases/
│       └── create_order.py
│
├── infrastructure/
│   ├── repositories/
│   │   └── sql_user_repository.py
│   ├── external/
│   │   └── stripe_payment.py
│   └── persistence/
│       └── models.py
│
└── presentation/
    ├── api/
    │   ├── routes/
    │   │   └── users.py
    │   └── schemas/
    │       └── user_schemas.py
    └── cli/
        └── commands.py
```

## Benefits

1. **Testability**: Domain and Application can be tested without database or framework
2. **Flexibility**: Swap database or framework without touching business logic
3. **Clarity**: Clear boundaries make code easier to navigate
4. **Maintainability**: Changes in one layer don't ripple through others

## Testing by Layer

```python
# Domain tests - pure unit tests, no mocks needed
def test_user_can_login_when_active():
    user = User(id="1", email="test@test.com", is_active=True)
    assert user.can_login() == True

# Application tests - mock the ports
def test_auth_service_login():
    mock_repo = Mock(spec=UserRepository)
    mock_repo.find_by_email.return_value = User(...)
    service = AuthService(mock_repo)
    user = service.login("test@test.com", "password")
    assert user is not None

# Infrastructure tests - integration tests with real DB
def test_sql_user_repository(test_db):
    repo = SqlUserRepository(test_db)
    repo.save(User(...))
    found = repo.find_by_email("test@test.com")
    assert found is not None

# Presentation tests - API tests
def test_login_endpoint(client):
    response = client.post("/login", json={"email": "...", "password": "..."})
    assert response.status_code == 200
```

## Using with AI

### Enforce in CLAUDE.md
```markdown
## Architecture Rules
- Clean Architecture: dependencies point inward
- Domain layer: no external imports
- Application layer: use ports (interfaces), not implementations
- Check layer violations before completing
```

### Ask AI to Verify
```
"Check if this implementation follows Clean Architecture.
Are there any layer violations?"
```

### Get AI to Design
```
"Design a user registration feature following Clean Architecture.
Show the components in each layer."
```
