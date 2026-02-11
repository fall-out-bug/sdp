# Contract-First Development Workflow

This guide explains the contract-first development workflow for coordinating parallel features that share boundaries.

## Overview

When multiple features work on the same types/interfaces, contract-first development ensures they agree on the interface before implementing.

## Workflow

### 1. Detect Shared Boundaries

When planning parallel features, detect shared boundaries:

```bash
sdp collision detect
```

Output:
```
🔗 Shared boundaries detected:

  File: internal/model/user.go
  Type: User
  Features: [F054 F055]
  Fields:
    - Email: string
    - Name: string
  Recommendation: Create shared interface contract

  1 shared boundary(ies)
```

### 2. Generate Contracts

From detected boundaries, generate contract files:

```bash
sdp contract generate --features F054,F055
```

This creates `.contracts/User.yaml`:

```yaml
typeName: User
fields:
  - name: Email
    type: string
  - name: Name
    type: string
requiredBy:
  - F054
  - F055
status: draft
fileName: internal/model/user.go
```

### 3. Lock Contract

Lock the contract to establish it as the source of truth:

```bash
sdp contract lock --contract .contracts/User.yaml
```

This creates `.contracts/User.lock` with SHA256 checksum.

### 4. Implement Against Contract

Both features implement against the locked contract:

```go
// Feature F054 implementation
type User struct {
    Email string
    Name  string
    // F054-specific fields allowed (warning in P1)
}
```

### 5. Validate Implementation

Validate implementation against contract:

```bash
sdp contract validate --impl-dir internal/model --contracts-dir .contracts
```

Output:
```
✓ No contract violations found
  All implementations match their contracts
```

## Violation Types

| Type | Severity | Description |
|------|----------|-------------|
| `missing_field` | Error | Required field not in implementation |
| `type_mismatch` | Warning | Field type doesn't match contract |
| `extra_field` | Warning | Field not in contract (P1: warning, P2: error) |

## Integration with @build

When `@build` executes workstreams with shared contracts:

1. Reads contract from `.contracts/`
2. Validates implementation against contract
3. Reports violations as warnings (P1) or errors (P2)

## Example: Two Features with Shared User Type

### Feature F054: Authentication
- Needs: `User.Email`, `User.Name`

### Feature F055: Profile
- Needs: `User.Email`, `User.Name`

### Workflow

```bash
# 1. Detect shared boundary
sdp collision detect

# 2. Generate contract
sdp contract generate --features F054,F055

# 3. Review and edit contract if needed
vim .contracts/User.yaml

# 4. Lock contract
sdp contract lock --contract .contracts/User.yaml

# 5. Features implement against contract
# F054 and F055 teams work in parallel

# 6. Validate before merge
sdp contract validate --impl-dir internal/model
```

## Best Practices

1. **Detect Early**: Run `sdp collision detect` before starting implementation
2. **Lock Agreed Contract**: Only lock after all features agree on interface
3. **Validate Often**: Run validation during development, not just at end
4. **Handle Conflicts**: Use synthesis engine to resolve field conflicts
5. **Update Contracts**: If contract needs change, agree then re-lock

## Migration Guide

### From Ad-Hoc Coordination

**Before (ad-hoc):**
- Email threads to agree on types
- Manual code review to check consistency
- Runtime errors from mismatches

**After (contract-first):**
- `sdp collision detect` finds boundaries
- `sdp contract generate` creates agreement
- `sdp contract validate` ensures consistency
- Compile-time validation

### From Single-Feature Contracts

**Before (single-feature):**
- One contract per feature
- No cross-feature coordination
- Duplicate type definitions

**After (cross-feature):**
- One contract per shared type
- Multiple features reference same contract
- Single source of truth

## Reference

- `sdp collision detect` - Find shared boundaries
- `sdp contract generate` - Create contracts from boundaries
- `sdp contract lock` - Lock contract as source of truth
- `sdp contract validate` - Validate implementation against contract
- `sdp contract verify` - Verify contract matches lock

## Status: P1 Complete

- ✅ Boundary detection (00-060-01)
- ✅ Contract generation (00-060-02)
- ✅ Contract validation (00-060-03)
- ⏳ Contract enforcement in @build (P2)
- ⏳ Cross-branch contract sync (P2)
