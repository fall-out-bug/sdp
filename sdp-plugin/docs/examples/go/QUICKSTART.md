# Go Quick Start

Get started with SDP plugin on your Go project in 5 minutes.

## Prerequisites

- Go 1.21+
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

- Coverage: ≥80% (go test -cover)
- Type Safety: Complete signatures
- Error Handling: No ignored errors
- File Size: <200 LOC
- Architecture: Clean import paths

## Project Structure

```
github.com/user/project/
├── domain/          # Entities, business logic
│   └── entities/
├── application/     # Use cases, services
│   └── service/
├── infrastructure/  # DB, external APIs
│   └── persistence/
└── presentation/    # HTTP handlers
    └── api/
```

## Example: Clean Architecture

```go
// Domain Layer (no dependencies)
package entities

type User struct {
    Email string
    // Pure business logic
}

// Application Layer (imports domain)
package service

import "github.com/user/project/domain/entities"

type AuthService struct {
    // Use case logic
}

// Infrastructure Layer (imports domain)
package persistence

import "github.com/user/project/domain/entities"

type UserRepository struct {
    // Data access
}

// Presentation Layer (imports application)
package api

import "github.com/user/project/service"

type AuthHandler struct {
    // HTTP handlers
}
```

## Error Handling

**❌ BAD: Ignore errors**
```go
data, _ := fetchData()  // Error lost
```

**✅ GOOD: Check errors**
```go
data, err := fetchData()
if err != nil {
    log.Printf("Fetch failed: %v", err)
    return err
}
```

See [TUTORIAL.md](../../TUTORIAL.md) for full details.
