# Workstream Sizing Guide

## Overview

Proper workstream sizing ensures AI agents can execute them successfully in one shot while maintaining focus and quality.

## Size Categories

### SMALL (< 500 LOC)

**Characteristics:**
- Single file or 2-3 small files
- Simple logic, minimal branching
- Clear input/output
- Few dependencies
- Straightforward testing

**Examples:**
- Add validation function
- Create value object
- Add utility method
- Simple refactoring

**Time Estimate:** 1-2 hours

**Example:**
```
WS-001-01: Add email validation
Files:
  - src/validators.py (+50 lines)
  - tests/unit/test_validators.py (+100 lines)
Total: ~150 LOC
```

### MEDIUM (500-1500 LOC)

**Characteristics:**
- 3-5 files
- Moderate complexity
- Some branching/error handling
- Multiple dependencies
- Comprehensive testing needed

**Examples:**
- Create repository layer
- Implement service class
- Add API endpoints
- Database migration

**Time Estimate:** 2-4 hours

**Example:**
```
WS-001-02: User repository layer
Files:
  - src/infrastructure/repositories/user_repository.py (+200 lines)
  - src/application/ports/user_repository.py (+50 lines)
  - tests/unit/test_user_repository.py (+300 lines)
  - tests/integration/test_user_repository_integration.py (+200 lines)
Total: ~750 LOC
```

### LARGE (> 1500 LOC)

**Characteristics:**
- Too complex for one workstream
- Must be split into smaller workstreams

**Rule:** No workstream should be LARGE. Split into 2+ MEDIUM workstreams.

**Example of splitting:**

**Before (LARGE - 2000 LOC):**
```
WS-001-01: Complete authentication system
  - Domain entities
  - Repository layer
  - Service layer
  - API endpoints
```

**After (split into 4 MEDIUM):**
```
WS-001-01: Domain entities (400 LOC)
WS-001-02: Repository layer (700 LOC)
WS-001-03: Service layer (600 LOC)
WS-001-04: API endpoints (300 LOC)
```

## Estimation Techniques

### Bottom-Up Estimation

Start with implementation files and add:

1. **Implementation files:** Count lines needed
2. **Test files:** Multiply implementation by 1.5-2x
3. **Documentation:** Add 10-20% for docstrings

**Example:**
```
Implementation: 200 lines
Tests: 300 lines (1.5x)
Documentation: 30 lines (10%)
Total: 530 lines → MEDIUM
```

### Top-Down Estimation

Start with complexity and derive LOC:

1. **Simple AC:** ~100 LOC per AC (including tests)
2. **Moderate AC:** ~200 LOC per AC
3. **Complex AC:** ~300 LOC per AC

**Example:**
```
AC1: User can login (moderate) → 200 LOC
AC2: Session created (simple) → 100 LOC
AC3: Token generation (moderate) → 200 LOC
Total: 500 LOC → MEDIUM
```

### Historical Data

Track actual vs estimated LOC for calibration:

```
WS-001-01: Estimated 400 LOC, Actual 550 LOC (+37%)
WS-001-02: Estimated 600 LOC, Actual 650 LOC (+8%)
WS-001-03: Estimated 300 LOC, Actual 280 LOC (-7%)
```

Average error: +13% → Apply correction factor

## Complexity Factors

### Increases LOC Estimate

- **Error handling:** +20-30%
- **Validation logic:** +15-25%
- **External integrations:** +30-50%
- **Complex algorithms:** +40-60%
- **State management:** +25-35%

### Decreases LOC Estimate

- **Reusing existing code:** -20-30%
- **Simple CRUD:** -15-25%
- **Code generation:** -30-40%

## AI Agent Constraints

### Context Window Limits

AI agents have finite context windows:

- **Small WS:** Full context fits easily
- **Medium WS:** Context management needed
- **Large WS:** Context overflow likely

### Focus Degradation

As workstream size increases, agent focus degrades:

- **< 500 LOC:** High focus, consistent quality
- **500-1000 LOC:** Good focus, minor deviations
- **1000-1500 LOC:** Moderate focus, some issues
- **> 1500 LOC:** Poor focus, quality suffers

**Rule:** Keep workstreams MEDIUM or smaller for best results.

## Splitting Large Workstreams

### By Layer

Split complex feature across architectural layers:

```
Original: Authentication system (2000 LOC)

Split:
  WS-01: Domain layer (400 LOC)
  WS-02: Repository layer (600 LOC)
  WS-03: Service layer (700 LOC)
  WS-04: API layer (300 LOC)
```

### By Functionality

Split feature by functional components:

```
Original: User management (2500 LOC)

Split:
  WS-01: User CRUD (800 LOC)
  WS-02: User validation (500 LOC)
  WS-03: User authentication (700 LOC)
  WS-04: User authorization (500 LOC)
```

### By Integration

Split by integration boundaries:

```
Original: Payment processing (3000 LOC)

Split:
  WS-01: Payment domain (600 LOC)
  WS-02: Stripe integration (800 LOC)
  WS-03: Payment service (900 LOC)
  WS-04: Payment API (700 LOC)
```

## Scope Files and Sizing

Scope files impact workstream size:

### Tight Scope (Better)

```
scope_files:
  - src/domain/user.py
  - tests/unit/test_user.py
```

Focused, easy to estimate.

### Loose Scope (Worse)

```
scope_files:
  - src/domain/*.py
  - src/infrastructure/*.py
  - tests/**/*.py
```

Vague, hard to estimate, likely too large.

## Quality vs Size Tradeoff

| Size | Quality | Coverage | Completion Rate |
|------|---------|----------|-----------------|
| SMALL | Excellent | 95%+ | 98% |
| MEDIUM | Good | 85%+ | 90% |
| LARGE | Poor | <80% | <70% |

**Recommendation:** Prefer MEDIUM workstreams (500-1000 LOC) for balance of quality and progress.

## Sizing Checklist

When creating workstream, verify:

- [ ] Estimated LOC ≤ 1500
- [ ] Number of files ≤ 5
- [ ] Acceptance Criteria ≤ 5
- [ ] Dependencies ≤ 3
- [ ] Scope files explicit and minimal
- [ ] Single responsibility principle followed

## Examples

### Well-Sized Workstreams

```
✅ WS-001-01: Create User entity (300 LOC)
✅ WS-001-02: User repository (750 LOC)
✅ WS-001-03: User service (600 LOC)
```

### Poorly-Sized Workstreams

```
❌ WS-001-01: Complete user system (2500 LOC)
   → Split into 4-5 workstreams

❌ WS-001-02: Add helper function (50 LOC)
   → Too small, merge with related WS

❌ WS-001-03: Refactor everything (5000 LOC)
   → Split into 8-10 workstreams
```

## See Also

- [Design Specification](design-spec.md)
- [Acceptance Criteria Guide](acceptance-criteria.md)
- [Dependency Management](dependencies.md)
