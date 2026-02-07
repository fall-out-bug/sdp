# Issue: Component Contract Validation Gap

**ID:** sdp-ie25  
**Type:** Feature  
**Priority:** P1 (High)  
**Status:** Open  
**Created:** 2026-02-07

## Problem Statement

During SDP-assisted feature implementation, agents ignored component integration points. Frontend, backend, and SDK ended up with incompatible API endpoints, rendering the feature non-functional despite all components being "complete."

**Symptom:**
- Telemetry feature implementation failed
- Backend exposed endpoints that frontend couldn't consume
- SDK methods didn't match backend API
- Components were developed in isolation without integration validation

## Root Cause Analysis

**Current @design workflow focuses on:**
- Domain entity modeling
- Workstream decomposition
- File organization
- Clean architecture layers

**Missing critical aspects:**
- PRE-CODE contract agreement (before parallel implementation)
- Multi-agent contract synthesis
- Conflict resolution for integration points
- Contract-first parallel execution workflow

**Why this happens:**
1. @design creates workstreams but doesn't specify contracts
2. Agents implement components in PARALLEL without coordination
3. No contract agreement phase BEFORE writing code
4. Each agent "implements as they see it"
5. Integration validation only happens AFTER code is written (too late!)

**Real-world example (user's project):**
```
Time T0: Requirements gathered via @idea
Time T1: 4 agents start implementation in parallel:
  - Frontend agent:  → implements POST /api/v1/telemetry/submit
  - Backend agent:   → implements POST /api/v1/telemetry/events
  - SDK agent:       → implements telemetryClient.submit(data)
  - Infra agent:     → implements collector setup
Time T2: All components report "complete"
Time T3: Integration test → 404 Not Found

Problem: No contract agreement at T0-T1, agents diverged
Solution: Contract synthesis MUST happen BEFORE T1
```

## Impact Assessment

**Severity:** P1 (Feature-breaking but not production-down)

**Impact:**
- Features ship but don't integrate properly
- Rework required post-implementation
- Wasted development time
- Reduced trust in SDP workflow

**Example from user's project:**
```
Frontend: POST /api/v1/telemetry/submit
Backend:  POST /api/v1/telemetry/events
SDK:      telemetryClient.submit(data)

Result: 404 Not Found
```

## Proposed Solution

### Phase 0: Contract Synthesis (BEFORE Implementation)

**Critical insight:** Contract agreement MUST happen BEFORE agents write code.

**Workflow:**
```
@feature "Add telemetry"
  ↓
@idea (gather requirements)
  ↓
@design creates workstreams
  ↓
**NEW: Contract Synthesis Phase**
  ↓
Architect Agent proposes initial contract
  ↓
Multi-Agent Review (parallel):
  - Frontend Agent: "Need /batch endpoint"
  - Backend Agent:  "Works for us"
  - SDK Agent:      "Matches our method naming"
  ↓
Synthesizer Agent resolves conflicts
  ↓
All agents agree → FIXED CONTRACT
  ↓
Now parallel implementation can begin
```

**0.1 Contract Synthesis Agent** (NEW)
```bash
sdp contract synthesize --feature=telemetry

# Analyzes requirements
# Proposes openapi.yaml
# Sends to frontend/backend/sdk agents
# Collects feedback
# Resolves conflicts via synthesis rules
# Outputs: agreed_contract.yaml
```

**0.2 Contract Agreement Protocol**
```yaml
phase: contract_synthesis
participants:
  - architect: proposes initial contract
  - frontend: reviews from client perspective
  - backend: reviews from server perspective
  - sdk: reviews from library perspective
  - synthesizer: resolves conflicts

synthesis_rules:
  - domain_expertise: frontend/backend/sdk have veto
  - quality_gate: all must agree
  - merge: combine suggestions
  - escalate: human if unresolved
```

**0.3 Contract Lock**
```bash
# Once agreed, contract is locked
sdp contract lock --feature=telemetry --sha=abc123

# Implementation must match locked contract
# @build checks compliance automatically
```

### Phase 1: Contract Extraction & Validation (POST-Implementation)

**1.1 Update @design Skill**
Add code analysis workstreams (NO manual questions to user):
- Scan existing backend code → extract routes/endpoints
- Scan frontend code → extract API calls
- Scan SDK code → extract public methods
- Compare and find mismatches
- Generate integration report

**1.2 Contract Extraction Agent**
Implement automated contract discovery:
```bash
sdp contract extract --component=backend
# Analyzes Go/Python/Node code
# Extracts: routes, handlers, request/response types
# Generates: openapi.yaml
```

**1.3 Contract Validation Agent**
Implement automated contract validation:
```bash
sdp contract validate
# Cross-references:
# - Frontend fetch() calls vs backend routes
# - SDK methods vs backend endpoints
# - Request/response schema consistency
# Reports mismatches with file locations
```

**1.4 Update @review Checklist**
Add automated contract compliance checks:
- [ ] Contract extraction ran successfully
- [ ] No endpoint mismatches found
- [ ] All integration points documented
- [ ] Frontend-backend alignment verified

### Phase 2: Tooling Support (00-053-05 to 00-053-07)

**2.1 Contract Generator**
```bash
sdp contract generate --feature=telemetry
# Generates:
# - openapi.yaml
# - integration.md
# - sequence-diagram.mmd
```

**2.2 Contract Linter**
```bash
sdp contract lint
# Validates:
# - Frontend API calls match backend routes
# - SDK methods align with backend endpoints
# - Request/response schemas consistent
```

**2.3 Mock Server Generation**
```bash
sdp contract mock --port=8080
# Starts mock server from OpenAPI spec
# Allows frontend development before backend ready
```

## Implementation Plan

### Phase 0: Contract Synthesis (PRE-Implementation)

**WS 00-053-00:** Contract synthesis agent - multi-agent agreement
- Analyzes requirements from @idea
- Architect agent proposes initial contract
- Frontend/backend/sdk agents review in parallel
- Synthesizer resolves conflicts
- Outputs locked contract (YAML/protobuf)
- MUST complete before parallel implementation

### Phase 1: Contract Extraction & Validation (POST-Implementation)

**WS 00-053-01:** Code analysis agent - extract endpoints from existing code
**WS 00-053-02:** Contract generation agent - create OpenAPI/protobuf from code
**WS 00-053-03:** Contract validation agent - cross-reference components
**WS 00-053-04:** Update @design skill - add contract synthesis phase
**WS 00-053-05:** Implement `sdp contract synthesize` CLI command
**WS 00-053-06:** Implement `sdp contract lock` CLI command
**WS 00-053-07:** Update @review skill - add contract compliance checks

### Dependencies
**Phase 0 (Critical Path):**
- 00-053-00 MUST complete before any component implementation starts
- Contract lock prevents agent divergence

**Phase 1:**
- 00-053-01 → 00-053-02 (extraction before generation)
- 00-053-02 → 00-053-03 (generation before validation)
- 00-053-03 → 00-053-04 (validation logic used in @design)
- 00-053-00 → 00-053-05 (synthesis agent used by CLI)
- 00-053-00 → 00-053-06 (synthesis outputs lockable contract)

## Success Criteria

**Functional:**
1. **PRE-CODE:** Contract synthesis happens BEFORE implementation
2. Multi-agent agreement via synthesis rules (domain expertise veto)
3. Locked contract guides parallel implementation
4. POST-CODE: Validation detects implementation drift
5. Zero integration mismatches in shipped features

**Quality:**
- Contract synthesis success rate ≥ 95% (auto-resolution without human)
- Contract compliance 100% (implementation matches contract)
- Rework due to contract issues reduced by 90%
- Parallel agent implementation synchronized from day 1

## Migration Guide

For existing features:
1. Run `sdp contract extract` to reverse-engineer contracts from code
2. Validate generated contracts
3. Fix mismatches
4. Add to documentation

## Alternatives Considered

**Alternative 1:** Manual contract documentation
- Pros: Simple, no tooling
- Cons: Easy to forget, not enforced, inconsistent

**Alternative 2:** Post-hoc contract generation
- Pros: Can add after feature complete
- Cons: Too late, contracts don't drive design

**Alternative 3:** Type-safe contracts (TypeScript/Go interfaces shared)
- Pros: Compile-time safety
- Cons: Language-specific, doesn't cover HTTP boundaries

**Selected:** Contract-first design with validation ✅
- Drives implementation from contracts
- Language-agnostic
- Enforced in review gate

## Related Issues

- User feedback: "Что толку в длинных опросах, если в итоге продукт разваливается?"
- Similar to: Interface design pattern enforcement
- Related to: API-first development best practices

## References

- [API Design Guide](https://google.github.io/styleguide/restapi-api-guide.html)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Contract-Driven Development](https://martinfowler.com/articles/contract-test.html)

---

**Created by:** /issue command  
**Assigned:** TBD  
**Sprint:** TBD  
**Story Points:** TBD
