# Example: Structured Mode Feature

A complete example of using Structured mode to add a feature.

## Scenario

Add password reset functionality to an existing authentication system.

## Phase 1: Analyze

### Prompt

```
Analyze this feature request and create a specification:

"Users need to be able to reset their passwords if they forget them.
They should click 'Forgot Password', enter their email, receive a link,
and be able to set a new password. For security, the link should expire."

Save to docs/specs/password-reset.md
```

### Output: docs/specs/password-reset.md

```markdown
# Feature: Password Reset

## Overview
Allow users to reset their password via email when they forget it.

## Requirements
- REQ-1: User can request password reset by providing email
- REQ-2: System sends reset link to the email address
- REQ-3: Reset link expires after 1 hour
- REQ-4: User can set new password using valid link
- REQ-5: Old password is invalidated after reset
- REQ-6: User receives confirmation after successful reset

## User Stories
- As a user who forgot my password, I want to reset it via email
  so that I can regain access to my account.

## Acceptance Criteria
- [ ] "Forgot Password" link visible on login page
- [ ] Entering email sends reset link (no indication if email exists)
- [ ] Reset link contains unique, unguessable token
- [ ] Clicking expired link shows "Link expired" message
- [ ] New password must meet existing password requirements
- [ ] After reset, user can login with new password
- [ ] Old password no longer works

## Out of Scope
- Rate limiting (handled by existing middleware)
- Account lockout
- Password history (preventing reuse)

## Open Questions
- Q1: Should reset invalidate existing sessions? → Yes
- Q2: Email template approval needed? → Use existing template style
```

## Phase 2: Design

### Prompt

```
Read docs/specs/password-reset.md and create technical design.
Consider security implications.
Save to docs/specs/password-reset-design.md
Create ADR if making significant decisions.
```

### Output: docs/specs/password-reset-design.md

```markdown
# Design: Password Reset

## Components

### Domain Layer
- `PasswordResetToken` (entity)
  - id: UUID
  - user_id: UUID
  - token_hash: string (hashed, not plaintext)
  - created_at: datetime
  - expires_at: datetime
  - used_at: datetime | null

### Application Layer
- `PasswordResetService`
  - initiate_reset(email) → sends email
  - validate_token(token) → bool
  - complete_reset(token, new_password) → void

### Infrastructure Layer
- `PasswordResetTokenRepository` (implements port)
- `EmailService` (existing, reuse)

### Presentation Layer
- POST /api/auth/forgot-password
- POST /api/auth/reset-password

## Data Flow

### Initiate Reset
1. User submits email to /forgot-password
2. Controller calls PasswordResetService.initiate_reset(email)
3. Service checks if user exists (silently skip if not)
4. Service generates secure token, hashes it, stores in DB
5. Service sends email with unhashed token in link
6. Return 200 OK (always, to prevent email enumeration)

### Complete Reset
1. User clicks link, submits new password to /reset-password
2. Controller calls PasswordResetService.complete_reset(token, password)
3. Service finds token by hash, checks expiry
4. Service updates user password, marks token used
5. Service invalidates existing sessions
6. Return 200 OK or 400 with error

## API

### POST /api/auth/forgot-password
```json
Request:  { "email": "user@example.com" }
Response: 200 OK (always)
```

### POST /api/auth/reset-password
```json
Request:  { "token": "abc123...", "new_password": "NewPass123!" }
Response: 200 OK
          400 { "error": "Token expired or invalid" }
          400 { "error": "Password does not meet requirements" }
```

## Database

```sql
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    token_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    CONSTRAINT unique_active_token UNIQUE (user_id, used_at)
);

CREATE INDEX idx_reset_tokens_hash ON password_reset_tokens(token_hash);
```

## Security Considerations
- Tokens hashed before storage (see ADR-0007)
- Constant-time comparison to prevent timing attacks
- No indication whether email exists
- HTTPS required (existing infrastructure)
```

### Output: docs/adr/0007-hash-reset-tokens.md

```markdown
# ADR-0007: Hash Password Reset Tokens Before Storage

## Status
Accepted

## Context
Password reset tokens are sensitive. If database is compromised,
plaintext tokens would allow attackers to reset any user's password.

## Decision
Hash tokens using SHA-256 before storing in database.
Store only the hash; compare incoming tokens by hashing and matching.

## Consequences
### Positive
- Database leak doesn't expose usable tokens
- Follows security best practices

### Negative
- Cannot retrieve original token (not needed)
- Slight overhead for hashing (negligible)
```

## Phase 3: Implement

### Prompt

```
Implement password reset based on:
- Spec: docs/specs/password-reset.md
- Design: docs/specs/password-reset-design.md

Start with domain layer, write tests first.
```

### Output (abbreviated)

**Domain Layer:**
```python
# src/domain/entities/password_reset_token.py
from datetime import datetime, timedelta
from dataclasses import dataclass

@dataclass
class PasswordResetToken:
    id: str
    user_id: str
    token_hash: str
    created_at: datetime
    expires_at: datetime
    used_at: datetime | None = None

    @classmethod
    def create(cls, user_id: str, token_hash: str) -> "PasswordResetToken":
        now = datetime.utcnow()
        return cls(
            id=str(uuid.uuid4()),
            user_id=user_id,
            token_hash=token_hash,
            created_at=now,
            expires_at=now + timedelta(hours=1),
            used_at=None
        )

    def is_expired(self) -> bool:
        return datetime.utcnow() > self.expires_at

    def is_used(self) -> bool:
        return self.used_at is not None

    def is_valid(self) -> bool:
        return not self.is_expired() and not self.is_used()

    def mark_used(self) -> None:
        self.used_at = datetime.utcnow()
```

**Tests:**
```python
# tests/domain/test_password_reset_token.py
def test_token_expires_after_one_hour():
    token = PasswordResetToken.create("user-123", "hash")
    assert token.is_valid() == True

    # Simulate expiry
    token.expires_at = datetime.utcnow() - timedelta(minutes=1)
    assert token.is_expired() == True
    assert token.is_valid() == False

def test_token_cannot_be_reused():
    token = PasswordResetToken.create("user-123", "hash")
    token.mark_used()
    assert token.is_used() == True
    assert token.is_valid() == False
```

**Application and Infrastructure layers follow...**

## Phase 4: Review

### Prompt

```
Review password reset implementation:
1. Check acceptance criteria from docs/specs/password-reset.md
2. Verify design was followed
3. Run tests, check coverage
4. Look for security issues
```

### Output

```markdown
## Review: Password Reset

### Acceptance Criteria
- [x] "Forgot Password" link visible on login page
- [x] Entering email sends reset link
- [x] Reset link contains unique token
- [x] Clicking expired link shows error
- [x] New password meets requirements
- [x] After reset, login works with new password
- [x] Old password no longer works

### Test Results
- Tests: 18 passed, 0 failed
- Coverage: 94% for new code

### Security Check
- [x] Tokens hashed before storage
- [x] Constant-time comparison used
- [x] No email enumeration
- [x] Sessions invalidated after reset

### Issues Found
None blocking. Ready to merge.
```

## Complete File Tree

```
docs/
├── specs/
│   ├── password-reset.md
│   └── password-reset-design.md
├── adr/
│   └── 0007-hash-reset-tokens.md
src/
├── domain/entities/password_reset_token.py
├── application/services/password_reset_service.py
├── application/ports/password_reset_repository.py
├── infrastructure/repositories/sql_password_reset_repository.py
└── presentation/controllers/password_reset_controller.py
tests/
├── domain/test_password_reset_token.py
├── application/test_password_reset_service.py
└── presentation/test_password_reset_api.py
```

## Time

~4 hours total:
- Phase 1 (Analyze): 20 minutes
- Phase 2 (Design): 40 minutes
- Phase 3 (Implement): 2.5 hours
- Phase 4 (Review): 30 minutes
