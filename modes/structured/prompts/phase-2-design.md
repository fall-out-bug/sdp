# Phase 2: Design

You are designing the technical solution for a feature.

## Your Task

Read the specification and create a technical design document.

## Input

- Specification from Phase 1: `docs/specs/{feature}.md`

## Output

### Primary: Design Document

Create `docs/specs/{feature}-design.md` with:

#### 1. Components
List all components (classes, modules, services) involved:
```
- UserService (application) - orchestrates password reset flow
- EmailService (infrastructure) - sends reset emails
- TokenRepository (infrastructure) - stores reset tokens
```

#### 2. Data Flow
Describe how data moves through the system:
```
1. User submits email → PasswordResetController
2. Controller calls UserService.initiateReset(email)
3. UserService generates token, saves via TokenRepository
4. UserService calls EmailService.sendResetLink(email, token)
5. User clicks link → PasswordResetController.validateToken
6. ...
```

#### 3. API Changes
New or modified endpoints:
```
POST /api/auth/forgot-password
  Request: { email: string }
  Response: 200 OK (always, for security)

POST /api/auth/reset-password
  Request: { token: string, newPassword: string }
  Response: 200 OK | 400 Bad Request (invalid/expired token)
```

#### 4. Database Changes
New tables, columns, or indexes:
```sql
CREATE TABLE password_reset_tokens (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  token_hash VARCHAR(255) NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  used_at TIMESTAMP
);
```

#### 5. Dependencies
New libraries or services:
```
- None required (using existing email infrastructure)
```

### Secondary: ADRs

Create `docs/adr/NNNN-{title}.md` for any architectural decisions:

**Create ADR when:**
- Choosing between alternatives (library A vs B)
- Making trade-offs (security vs convenience)
- Deviating from existing patterns
- Introducing new technology

**Don't create ADR for:**
- Obvious choices with no alternatives
- Following existing patterns

## Example Prompt

```
Read docs/specs/password-reset.md and create technical design.

Consider:
- Where to store reset tokens (DB vs Redis)
- How to generate secure tokens
- Rate limiting to prevent abuse

If making decisions, create ADRs in docs/adr/
Save design to docs/specs/password-reset-design.md
```

## Quality Checklist

Before finishing this phase:
- [ ] All requirements from spec are addressed
- [ ] Components follow Clean Architecture
- [ ] Data flow covers success and error paths
- [ ] API design is RESTful and consistent
- [ ] Database changes are migration-ready
- [ ] ADRs created for significant decisions

## Tips

- Reference the specification frequently
- Consider security implications
- Think about scalability (what if 10x users?)
- Don't over-design (YAGNI - You Aren't Gonna Need It)
- Mark assumptions clearly
