# Java Quick Start

Get started with SDP plugin on your Java project in 5 minutes.

## Prerequisites

- Java 17+
- Maven or Gradle
- Claude Code installed

## Installation

```bash
# 1. Copy SDP prompts to your project
cp -r sdp-plugin/prompts/* .claude/

# 2. Verify installation
ls .claude/skills/
# Should show: feature.md, design.md, build.md, review.md, etc.
```

## Your First Feature

### Step 1: Create Feature

```
@feature "Add user authentication"
```

SDP asks deep questions about requirements.

### Step 2: Plan Workstreams

```
@design feature-user-auth
```

SDP creates workstreams for implementation.

### Step 3: Execute

```
@build 00-001-01
```

Runs TDD cycle with AI validation.

### Step 4: Review

```
@review F01
```

Validates all quality gates.

## Quality Gates

- Coverage: ≥80% (JaCoCo)
- Type Safety: Complete signatures
- Error Handling: No empty catch
- File Size: <200 LOC
- Architecture: Clean package separation

## Project Structure

```
src/main/java/com/example/
├── domain/          # Entities, value objects
│   └── model/
├── application/     # Use cases, services
│   └── service/
├── infrastructure/  # DB, external APIs
│   └── persistence/
└── presentation/    # Controllers, views
    └── controller/
```

## Example: Clean Architecture

```java
// Domain Layer (no dependencies)
package com.example.domain.model;

public class User {
    private final String email;
    // Pure business logic
}

// Application Layer (imports domain)
package com.example.application.service;

import com.example.domain.model.User;

public class AuthService {
    // Use case logic
}

// Infrastructure Layer (imports domain)
package com.example.infrastructure.persistence;

import com.example.domain.model.User;

public class UserRepository {
    // Data access
}

// Presentation Layer (imports application)
package com.example.presentation.controller;

import com.example.application.service.AuthService;

@RestController
public class AuthController {
    // HTTP endpoints
}
```

See [TUTORIAL.md](../../TUTORIAL.md) for full details.
