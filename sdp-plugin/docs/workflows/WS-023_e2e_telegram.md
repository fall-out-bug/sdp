# WS-023: E2E Tests with Real Telegram

> **Workstream ID:** WS-023  
> **Feature:** F024: Unified Workflow  
> **Dependencies:** WS-019

## Goal

Test Telegram notifications end-to-end.

## Acceptance Criteria

### AC1: Test Message Received
- [ ] Test message sent to Telegram
- [ ] Message appears in channel/group

### AC2: Message Format
- [ ] Markdown formatted correctly
- [ ] Special characters escaped

### AC3: Error Handling
- [ ] Invalid bot token logged
- [ ] Network error doesn't crash

## Scope

E2E tests with real Telegram bot (use test channel).

## Estimated Scope

- ~150 LOC (test code)
- Duration: 1.5 hours
