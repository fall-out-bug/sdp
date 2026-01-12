# Testing Strategy: User Profile API

## Overview

This document defines the testing strategy for the User Profile API feature.

## Test Pyramid

```
        ┌─────────┐
        │  E2E    │  ← Few, slow, high confidence
        │  Tests  │
       ─┴─────────┴─
      ┌─────────────┐
      │ Integration │  ← Some, medium speed
      │   Tests     │
     ─┴─────────────┴─
    ┌─────────────────┐
    │   Unit Tests    │  ← Many, fast, focused
    │                 │
    └─────────────────┘
```

## Unit Tests

### Domain Layer
| Test Case | Input | Expected Output |
|-----------|-------|-----------------|
| Valid profile creation | Valid data | UserProfile instance |
| Name too long | 101 chars | ValidationError |
| Name empty | "" | ValidationError |
| Bio too long | 501 chars | ValidationError |
| Invalid email format | "notanemail" | ValidationError |

### Application Layer (Use Cases)
| Test Case | Setup | Expected |
|-----------|-------|----------|
| Get profile - success | User exists, authorized | Profile returned |
| Get profile - not found | User doesn't exist | NotFoundError |
| Get profile - unauthorized | Requester != target | ForbiddenError |
| Update profile - success | Valid data, authorized | Updated profile |
| Update profile - validation fail | Invalid data | ValidationError |
| Update profile - unauthorized | Requester != target | ForbiddenError |

### Coverage Target
- Domain layer: ≥90%
- Application layer: ≥85%
- Overall: ≥80%

## Integration Tests

### Repository Adapter
| Test Case | Setup | Verification |
|-----------|-------|--------------|
| Find by ID - exists | Seed DB with user | Returns correct profile |
| Find by ID - not exists | Empty DB | Returns null |
| Update profile | Existing user | DB updated correctly |
| Connection timeout | Slow DB | Timeout error thrown |
| Connection failure | DB unavailable | Connection error thrown |

### Authorization Adapter
| Test Case | Setup | Verification |
|-----------|-------|--------------|
| Can view own profile | Same user ID | Returns true |
| Can view other profile | Different user, no admin | Returns false |
| Can view as admin | Admin user | Returns true |
| Auth service timeout | Slow service | Timeout error, deny access |

## End-to-End Tests

### Happy Path
```gherkin
Feature: User Profile API

Scenario: Get own profile
  Given I am authenticated as user "123"
  When I request GET /api/users/123/profile
  Then I receive status 200
  And the response contains my profile data

Scenario: Update own profile
  Given I am authenticated as user "123"
  When I request PUT /api/users/123/profile with valid data
  Then I receive status 200
  And my profile is updated

Scenario: Cannot view other user's profile
  Given I am authenticated as user "123"
  When I request GET /api/users/456/profile
  Then I receive status 403
```

### Error Cases
```gherkin
Scenario: Profile not found
  Given I am authenticated as user "123"
  When I request GET /api/users/999/profile
  Then I receive status 404

Scenario: Invalid update data
  Given I am authenticated as user "123"
  When I request PUT /api/users/123/profile with name longer than 100 chars
  Then I receive status 400
  And the response contains validation errors
```

## Performance Tests

| Metric | Target | Method |
|--------|--------|--------|
| Response time (p50) | < 100ms | Load test with k6/artillery |
| Response time (p95) | < 200ms | Load test with k6/artillery |
| Throughput | > 100 RPS | Sustained load test |
| Error rate | < 0.1% | Error monitoring during load |

## Security Tests

| Test | Tool | Check |
|------|------|-------|
| SQL Injection | OWASP ZAP | Input sanitization |
| Authorization bypass | Manual | All endpoints require auth |
| Rate limiting | Artillery | Rate limits enforced |
| Input validation | Fuzz testing | No crashes on malformed input |

## Test Environment

### Local Development
```bash
# Run unit tests
npm test -- --coverage

# Run integration tests (requires test DB)
npm run test:integration

# Run E2E tests (requires running server)
npm run test:e2e
```

### CI/CD Pipeline
1. Unit tests run on every commit
2. Integration tests run on PR
3. E2E tests run before merge to main
4. Performance tests run nightly

## Test Data

### Fixtures
```json
{
  "validUser": {
    "id": "test-user-123",
    "email": "test@example.com",
    "name": "Test User",
    "bio": "Test bio"
  },
  "adminUser": {
    "id": "admin-user-001",
    "email": "admin@example.com",
    "name": "Admin User",
    "role": "admin"
  }
}
```

### Cleanup
- Tests must clean up created data
- Use transactions that rollback after tests
- Isolated test database for integration tests

---

**Status**: Approved by QA
**Version**: 1.0
**Last Updated**: 2025-12-27
