# Architecture: User Profile API

## Overview

This document describes the architecture for the User Profile API feature, following Clean Architecture principles.

## Layer Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    PRESENTATION LAYER                        │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │ ProfileController│  │ ProfileRoutes   │                   │
│  └────────┬────────┘  └────────┬────────┘                   │
│           │                    │                             │
└───────────┼────────────────────┼─────────────────────────────┘
            │                    │
┌───────────┼────────────────────┼─────────────────────────────┐
│           ▼                    ▼        APPLICATION LAYER    │
│  ┌─────────────────────────────────────┐                     │
│  │       ProfileService (Use Cases)    │                     │
│  │  - GetProfileUseCase                │                     │
│  │  - UpdateProfileUseCase             │                     │
│  └────────────────┬────────────────────┘                     │
│                   │                                          │
│  ┌────────────────┴────────────────┐                         │
│  │         PORTS (Interfaces)       │                        │
│  │  - ProfileRepositoryPort         │                        │
│  │  - AuthorizationPort             │                        │
│  └──────────────────────────────────┘                        │
└──────────────────────────────────────────────────────────────┘
                          │
┌─────────────────────────┼────────────────────────────────────┐
│                         ▼        INFRASTRUCTURE LAYER        │
│  ┌─────────────────────────────────────┐                     │
│  │          ADAPTERS                    │                    │
│  │  - ProfileRepositoryAdapter (DB)     │                    │
│  │  - AuthorizationAdapter              │                    │
│  └─────────────────────────────────────┘                     │
└──────────────────────────────────────────────────────────────┘
                          │
┌─────────────────────────┼────────────────────────────────────┐
│                         ▼           DOMAIN LAYER             │
│  ┌─────────────────────────────────────┐                     │
│  │           ENTITIES                   │                    │
│  │  - UserProfile                       │                    │
│  │  - ProfileValidation                 │                    │
│  └─────────────────────────────────────┘                     │
└──────────────────────────────────────────────────────────────┘
```

## Components

### Domain Layer

| Component | Description | Dependencies |
|-----------|-------------|--------------|
| UserProfile | Entity representing user profile data | None |
| ProfileValidation | Validation rules for profile fields | None |

### Application Layer

| Component | Description | Dependencies |
|-----------|-------------|--------------|
| GetProfileUseCase | Retrieves user profile | ProfileRepositoryPort, AuthorizationPort |
| UpdateProfileUseCase | Updates user profile | ProfileRepositoryPort, AuthorizationPort |
| ProfileRepositoryPort | Interface for profile storage | None (abstract) |
| AuthorizationPort | Interface for auth checks | None (abstract) |

### Infrastructure Layer

| Component | Description | Dependencies |
|-----------|-------------|--------------|
| ProfileRepositoryAdapter | Database implementation | Database client |
| AuthorizationAdapter | Auth service integration | Auth service client |

### Presentation Layer

| Component | Description | Dependencies |
|-----------|-------------|--------------|
| ProfileController | HTTP request handler | Use cases |
| ProfileRoutes | Route definitions | Controller |

## Data Flow

### GET Profile Request
```
1. HTTP Request → ProfileRoutes
2. ProfileRoutes → ProfileController.getProfile()
3. ProfileController → GetProfileUseCase.execute()
4. GetProfileUseCase → AuthorizationPort.canView()
5. GetProfileUseCase → ProfileRepositoryPort.findById()
6. Response ← ProfileController ← GetProfileUseCase
```

### PUT Profile Request
```
1. HTTP Request → ProfileRoutes
2. ProfileRoutes → ProfileController.updateProfile()
3. ProfileController → UpdateProfileUseCase.execute()
4. UpdateProfileUseCase → AuthorizationPort.canEdit()
5. UpdateProfileUseCase → ProfileValidation.validate()
6. UpdateProfileUseCase → ProfileRepositoryPort.update()
7. Response ← ProfileController ← UpdateProfileUseCase
```

## Contracts

### ProfileRepositoryPort
```typescript
interface ProfileRepositoryPort {
  findById(userId: string): Promise<UserProfile | null>;
  update(userId: string, data: ProfileUpdateData): Promise<UserProfile>;
}
```

### AuthorizationPort
```typescript
interface AuthorizationPort {
  canViewProfile(requesterId: string, targetId: string): Promise<boolean>;
  canEditProfile(requesterId: string, targetId: string): Promise<boolean>;
}
```

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Database timeout | High | Implement retry with exponential backoff, set connection timeout |
| Auth service unavailable | High | Cache authorization decisions, fail closed |
| Invalid input data | Medium | Validate at controller and domain layers |

## Security Considerations

1. **Authorization**: Every request must pass through AuthorizationPort
2. **Input validation**: Sanitize and validate all user input
3. **Error messages**: Don't leak internal details in error responses

## Decisions

- **ADR-001**: Using port/adapter pattern for repository (see decision_log)
- **ADR-002**: Authorization at use case level, not controller level

---

**Status**: Approved by Architect
**Version**: 1.0
**Last Updated**: 2025-12-27
