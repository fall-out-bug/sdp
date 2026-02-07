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
- API contract definition (OpenAPI/Swagger/protobuf)
- Integration point specification
- Endpoint alignment validation
- Contract-first design methodology

**Why this happens:**
1. @idea asks about requirements but NOT integration contracts
2. @design creates workstreams but NOT API specifications
3. @build implements in isolation without contract validation
4. @review checks code quality but NOT integration compatibility

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

### Phase 1: Contract-First Design (00-053-01 to 00-053-04)

**1.1 Update @idea Skill**
Add contract-focused questions:
- "What are the API contracts between components?"
- "Which components need to communicate?"
- "Are there existing contracts to extend?"
- "Will this use REST, gRPC, WebSocket, or message queue?"

**1.2 Update @design Skill**
Add contract specification workstreams:
- Generate OpenAPI/Swagger specs
- Define protobuf messages if using gRPC
- Document integration points
- Create sequence diagrams for component interaction

**1.3 Add Contract Validation**
Create `sdp contract validate` command:
- Check OpenAPI spec completeness
- Validate endpoint consistency across components
- Verify request/response schemas match
- Ensure all integration points are documented

**1.4 Update @review Checklist**
Add contract compliance checks:
- [ ] API contract exists (OpenAPI/protobuf)
- [ ] All endpoints match contract
- [ ] Request/response schemas documented
- [ ] Integration points validated

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

### Workstreams

**WS 00-053-01:** Update @idea skill - add contract questions  
**WS 00-053-02:** Update @design skill - add contract spec generation  
**WS 00-053-03:** Implement `sdp contract validate` command  
**WS 00-053-04:** Update @review - add contract compliance checks  
**WS 00-053-05:** Implement contract generator  
**WS 00-053-06:** Implement contract linter  
**WS 00-053-07:** Implement mock server generator  

### Dependencies
- 00-053-01 → 00-053-02 (design depends on idea)
- 00-053-02 → 00-053-03 (validate needs specs)
- 00-053-03 → 00-053-04 (review uses validation)
- 00-053-03 → 00-053-05, 00-053-06, 00-053-07 (parallel tooling)

## Success Criteria

**Functional:**
1. @idea asks about integration contracts
2. @design generates OpenAPI/protobuf specs
3. @review validates contract compliance
4. No more integration mismatches

**Quality:**
- Contract completeness ≥ 90%
- Integration validation coverage 100%
- Rework due to contract issues reduced by 80%

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
