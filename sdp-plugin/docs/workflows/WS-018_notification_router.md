# WS-018: NotificationRouter Implementation

> **Workstream ID:** WS-018  
> **Feature:** F024: Unified Workflow  
> **Dependencies:** WS-017

## Goal

Route notifications to appropriate providers.

## Acceptance Criteria

### AC1: Event Routing
- [ ] Routes workstream events to Telegram
- [ ] Routes approval events to all providers
- [ ] Routes bug events to selected providers

### AC2: Fallback Providers
- [ ] Primary provider fails → fallback
- [ ] Multiple fallbacks supported
- [ ] Fallback logged

### AC3: Mock Provider
- [ ] Mock provider for testing
- [ ] Captures notifications in memory
- [ ] Used in unit tests

## Scope Files

**internal/orchestrator/notification_router.go** (NEW)

## Estimated Scope

- ~200 LOC
- Duration: 2 hours
