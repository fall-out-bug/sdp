# F060: Shared Contracts for Parallel Features

> Beads: sdp-2z3f | Priority: P1

---

## Problem

F054 scope collision detection tells you "two workstreams touch the same file." But it doesn't help you prevent the conflict — it just signals after the fact.

When 10 agents build 5 features in parallel, they need to agree on shared interfaces BEFORE building. Not after.

## Solution

Extend `@design` to detect shared boundaries across parallel features and generate interface contracts.

### How It Works

1. **Boundary Detection**: When `@design` runs for Feature B while Feature A is in progress, scan Feature A's workstreams for overlapping scope
2. **Contract Generation**: For shared surfaces (API endpoints, data models, function signatures), generate explicit contract files
3. **Contract-First Build**: Parallel workstreams build against shared contracts, not assumptions
4. **Contract Synthesis**: Use existing synthesis engine to resolve conflicts between features' needs

### Example

```
Feature A (@design): needs User model with email, name
Feature B (@design): needs User model with email, role, org_id

→ Contract: User model must have email, name, role, org_id
→ Both features implement against this contract
→ No surprise at merge time
```

### Contract File Format

```yaml
# .contracts/user-model.yaml
contract:
  type: data_model
  name: User
  features: [F054, F055]
  fields:
    - name: email
      type: string
      required_by: [F054, F055]
    - name: name
      type: string
      required_by: [F054]
    - name: role
      type: string
      required_by: [F055]
  status: locked
  locked_at: 2026-02-08T14:00:00Z
```

## Constraints

- Builds on F054 scope collision detection (detection → prevention)
- Uses existing synthesis engine for conflict resolution
- Contracts are suggestions, not enforcement (in P1)
- P2 adds enforcement (contract validation in CI)

## Users

- Teams running parallel features (3+ simultaneous)
- Architecture teams wanting interface stability
- SDP orchestrator (`@oneshot`) for automatic coordination

## Success Metrics

- Shared boundary detected in at least 1 real parallel development case
- Contract prevents at least 1 merge conflict that scope collision would have only warned about
- Synthesis engine resolves at least 1 cross-feature conflict

## Dependencies

- F054 (scope collision detection — foundation)
- Existing synthesis engine, contract validation system

## Notes

- This is the P1 step in the North Star path: P0 collision → **P1 contracts** → P2 cross-branch → beyond
- The contract validation system already exists in SDP (OpenAPI contracts)
- Extension: from single-feature contracts to cross-feature contracts
