# WS-014: RoleLoader and Prompt Management

> **Workstream ID:** WS-014  
> **Feature:** Unified Workflow  
> **Dependencies:** WS-012

## Goal

Load role prompts from prompts/roles/ directory.

## Acceptance Criteria

### AC1: Load Roles from Filesystem
- [ ] Loads 100+ roles from .md files
- [ ] Parses role name, description, permissions
- [ ] Validates role schema

### AC2: Role Validation
- [ ] Required fields present
- [ ] Permissions valid
- [ ] No duplicate roles

### AC3: Role Caching
- [ ] Roles cached for performance
- [ ] Cache invalidated on file change
- [ ] <100ms load time for 100 roles

## Scope Files

**internal/orchestrator/role_loader.go** (NEW)
**internal/orchestrator/role_loader_test.go** (NEW)

## Estimated Scope

- ~200 LOC
- Duration: 2 hours
