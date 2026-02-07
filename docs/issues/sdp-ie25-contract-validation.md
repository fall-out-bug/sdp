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
- Automated code analysis to extract existing contracts
- Integration point discovery from codebase
- Endpoint mismatch detection
- Automated contract generation/validation

**Why this happens:**
1. @design creates NEW code but doesn't analyze EXISTING integration points
2. No automated scanning of backend routes/frontend calls/SDK methods
3. No cross-component validation during design phase
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

### Workstreams

**WS 00-053-01:** Code analysis agent - extract endpoints from existing code
**WS 00-053-02:** Contract generation agent - create OpenAPI/protobuf from code
**WS 00-053-03:** Contract validation agent - cross-reference components
**WS 00-053-04:** Update @design skill - add contract analysis phase
**WS 00-053-05:** Implement `sdp contract extract` CLI command
**WS 00-053-06:** Implement `sdp contract validate` CLI command
**WS 00-053-07:** Update @review skill - add contract validation

### Dependencies
- 00-053-01 → 00-053-02 (extraction before generation)
- 00-053-02 → 00-053-03 (generation before validation)
- 00-053-03 → 00-053-04 (validation logic used in @design)
- 00-053-01 → 00-053-05 (extraction agent used by CLI)
- 00-053-03 → 00-053-06 (validation agent used by CLI)

## Success Criteria

**Functional:**
1. Agents automatically extract contracts from existing code
2. @design analyzes integration points without human input
3. `sdp contract validate` detects endpoint mismatches
4. @review enforces contract compliance

**Quality:**
- Contract extraction accuracy ≥ 95%
- Integration mismatch detection 100%
- Rework due to contract issues reduced by 80%
- Zero manual contract maintenance required

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
