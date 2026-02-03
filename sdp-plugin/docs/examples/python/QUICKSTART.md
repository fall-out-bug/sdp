# Python Quick Start

Get started with SDP plugin on your Python project in 5 minutes.

## Prerequisites

- Python 3.10+
- pytest installed
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

- Coverage: ≥80% (pytest)
- Type Safety: Complete hints (mypy)
- Error Handling: No bare except
- File Size: <200 LOC
- Architecture: Clean layers

See [TUTORIAL.md](../../TUTORIAL.md) for details.
