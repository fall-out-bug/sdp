# F071: Team UX and Collaboration Loop — Workstream Summary

> Feature: F071 | Priority: P1 | 5 workstreams
> Theme: predictable handoff and multi-actor coordination

## Goal

Make SDP comfortable for team usage by formalizing ownership, handoff, collision signals, and shared operational rhythm.

## Dependency Graph

```text
00-071-01 Team Operating Model + Ownership
   ├─► 00-071-02 Handoff Package Format
   ├─► 00-071-03 Cross-Feature Scope Collision UX
   └─► 00-071-04 Team Dashboard/Notification Integration

00-071-05 Team Adoption Playbook
    depends on: 02, 03, 04
```

## Workstream Registry

| WS ID | Title | Priority | Size | Depends On |
|-------|-------|----------|------|------------|
| 00-071-01 | Team Operating Model and Role Boundaries | P1 | SMALL | — |
| 00-071-02 | Handoff Package and Session Continuity Contract | P1 | MEDIUM | 01 |
| 00-071-03 | Scope Collision Collaboration UX | P1 | MEDIUM | 01 |
| 00-071-04 | Team Dashboard and Notification Hooks | P2 | MEDIUM | 01 |
| 00-071-05 | Team Adoption Playbook and Guardrails | P1 | SMALL | 02,03,04 |

## Execution Phases

### Phase 1: Operating Model
- WS-01 defines team boundaries and expectations.

### Phase 2: Collaboration Surfaces
- WS-02 handoff package
- WS-03 collision UX
- WS-04 dashboard/notifications

### Phase 3: Rollout
- WS-05 adoption playbook and team default policies.

## Feature-Level Deliverables

- Team roles and ownership model aligned with protocol
- Repeatable handoff format for agent/human transitions
- Non-blocking scope collision alerts across features
- Operational playbook for adoption in multi-dev teams
