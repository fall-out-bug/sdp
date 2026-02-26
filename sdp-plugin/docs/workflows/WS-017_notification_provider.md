# WS-017: NotificationProvider Interface

> **Workstream ID:** WS-017  
> **Feature:** Unified Workflow  
> **Dependencies:** WS-013

## Goal

Define notification provider interface for multiple notification backends.

## Acceptance Criteria

### AC1: Interface Defined
- [ ] NotificationProvider interface with 3 methods
- [ ] Send(message) async
- [ ] Notify(event) async

### AC2: Event Types
- [ ] workstream_complete
- [ ] approval_needed
- [ ] bug_found
- [ ] feature_complete

### AC3: Provider Registration
- [ ] RegisterProvider(name, provider)
- [ ] UnregisterProvider(name)
- [ ] GetProvider(name)

## Scope Files

**internal/orchestrator/notification_provider.go** (NEW)

## Estimated Scope

- ~80 LOC (interface only)
- Duration: 1 hour
