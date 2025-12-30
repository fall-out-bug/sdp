# Example: Multi-Agent Epic

A complete example of using Multi-Agent mode for a larger feature.

## Scenario

Add a complete user notification system: users can set preferences, receive notifications via email and in-app, and view notification history.

## Setup

```
docs/specs/epic-notifications/
├── epic.md
└── consensus/
    ├── artifacts/
    ├── messages/inbox/{agent}/
    └── decision_log/
```

## Epic Definition

### docs/specs/epic-notifications/epic.md

```markdown
# Epic: User Notification System

## Business Need
Users need to receive notifications about important events (order updates,
security alerts, marketing) through their preferred channels.

## High-Level Requirements
1. Users can set notification preferences (email, in-app, push)
2. System sends notifications through selected channels
3. Users can view notification history
4. Users can mark notifications as read
5. Admin can send broadcast notifications

## Success Metrics
- 80% of users set preferences within first week
- <5 minute delivery time for critical notifications
- <1% unsubscribe rate

## Timeline
Target: 3 weeks
```

## Agent 1: Analyst

### Session

```
Read prompts/analyst_prompt.md
Analyze docs/specs/epic-notifications/epic.md
Create detailed requirements
```

### Output: consensus/artifacts/requirements.md

```markdown
# Requirements: User Notification System

## Functional Requirements

### FR-1: Notification Preferences
- FR-1.1: User can enable/disable email notifications
- FR-1.2: User can enable/disable in-app notifications
- FR-1.3: User can enable/disable push notifications
- FR-1.4: User can set quiet hours (no notifications during period)
- FR-1.5: Preferences persist across sessions

### FR-2: Notification Delivery
- FR-2.1: System routes notifications based on user preferences
- FR-2.2: Critical notifications bypass quiet hours
- FR-2.3: Failed deliveries are retried (3 attempts)
- FR-2.4: Delivery status is tracked

### FR-3: Notification History
- FR-3.1: User can view last 100 notifications
- FR-3.2: User can filter by type (order, security, marketing)
- FR-3.3: User can mark individual notifications as read
- FR-3.4: User can mark all as read

### FR-4: Admin Broadcast
- FR-4.1: Admin can send notification to all users
- FR-4.2: Admin can send to segment (by preference)
- FR-4.3: Broadcast has scheduling option

## Non-Functional Requirements
- NFR-1: Delivery within 5 minutes for normal, 1 minute for critical
- NFR-2: Support 10,000 concurrent notifications
- NFR-3: 99.9% delivery success rate
- NFR-4: Notification history loads in <500ms

## User Stories
1. As a user, I want to control how I receive notifications so I'm not overwhelmed
2. As a user, I want to see my notification history so I don't miss anything
3. As an admin, I want to broadcast announcements to all users
```

### Message to Architect

```json
// consensus/messages/inbox/architect/2024-01-15-requirements-ready.json
{
  "d": "2024-01-15",
  "r": "analyst",
  "st": "completed",
  "epic": "epic-notifications",
  "sm": "Requirements analysis complete. 4 feature areas identified.",
  "nx": ["architect"],
  "artifacts": ["consensus/artifacts/requirements.md"]
}
```

## Agent 2: Architect

### Session

```
Read prompts/architect_prompt.md
Read consensus/messages/inbox/architect/
Review requirements and create architecture
VETO if Clean Architecture would be violated
```

### Output: consensus/artifacts/architecture.md

```markdown
# Architecture: User Notification System

## Component Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Presentation Layer                       │
│  NotificationController │ PreferencesController │ AdminAPI   │
├─────────────────────────────────────────────────────────────┤
│                     Application Layer                        │
│  NotificationService │ PreferencesService │ DeliveryService │
├─────────────────────────────────────────────────────────────┤
│                     Domain Layer                             │
│  Notification │ UserPreferences │ DeliveryChannel           │
├─────────────────────────────────────────────────────────────┤
│                     Infrastructure Layer                     │
│  NotificationRepo │ EmailSender │ PushSender │ QueueWorker  │
└─────────────────────────────────────────────────────────────┘
```

## Domain Model

### Notification
- id, user_id, type, title, body, priority
- created_at, read_at, delivered_at

### UserPreferences
- user_id, email_enabled, inapp_enabled, push_enabled
- quiet_hours_start, quiet_hours_end

### DeliveryChannel (enum)
- EMAIL, IN_APP, PUSH

## Key Decisions

### Async Processing
Notifications queued for async delivery to handle scale.
See ADR-0008.

### Preference Caching
User preferences cached in Redis (5 min TTL) to reduce DB load.
See ADR-0009.

## API Endpoints

### User Endpoints
- GET /api/notifications - list user's notifications
- PUT /api/notifications/:id/read - mark as read
- GET /api/preferences - get notification preferences
- PUT /api/preferences - update preferences

### Admin Endpoints
- POST /api/admin/notifications/broadcast - send broadcast

## Database Schema

### notifications
- id, user_id, type, title, body, priority
- created_at, read_at
- INDEX on (user_id, created_at)

### notification_preferences
- user_id (PK), email_enabled, inapp_enabled, push_enabled
- quiet_hours_start, quiet_hours_end

### notification_deliveries
- id, notification_id, channel, status, attempted_at, delivered_at
```

### ADRs Created

- `docs/adr/0008-async-notification-queue.md`
- `docs/adr/0009-preference-caching.md`

### Message to Tech Lead

```json
// consensus/messages/inbox/tech_lead/2024-01-16-architecture-approved.json
{
  "d": "2024-01-16",
  "r": "architect",
  "st": "approved",
  "epic": "epic-notifications",
  "sm": "Architecture approved. 3 workstreams recommended.",
  "nx": ["tech_lead"],
  "artifacts": [
    "consensus/artifacts/architecture.md",
    "docs/adr/0008-async-notification-queue.md",
    "docs/adr/0009-preference-caching.md"
  ]
}
```

## Agent 3: Tech Lead

### Output: consensus/artifacts/implementation.md

```markdown
# Implementation Plan: User Notification System

## Workstreams

### Workstream 1: Core Domain & Preferences (Day 1-3)
- Notification entity and repository
- UserPreferences entity and repository
- PreferencesService
- Preferences API endpoints
- Tests

### Workstream 2: Notification Delivery (Day 4-7)
- NotificationService
- DeliveryService with channel routing
- Email sender implementation
- Push sender implementation
- Queue worker for async processing
- Tests

### Workstream 3: History & Admin (Day 8-10)
- Notification history API
- Mark as read functionality
- Admin broadcast endpoint
- Admin UI integration
- Tests

## Dependencies
- Workstream 2 depends on Workstream 1 (needs preferences)
- Workstream 3 depends on Workstream 2 (needs delivery)

## Review Points
- After each workstream: code review with Tech Lead
- After Workstream 2: Architecture audit
- After Workstream 3: Full QA cycle
```

## Agent 4: Developer (Workstream 1)

### Session

```
Read prompts/developer_prompt.md
Implement Workstream 1 from implementation.md
TDD approach, domain layer first
```

### After Implementation

```json
// consensus/messages/inbox/tech_lead/2024-01-18-workstream1-complete.json
{
  "d": "2024-01-18",
  "r": "developer",
  "st": "review_requested",
  "epic": "epic-notifications",
  "sm": "Workstream 1 complete. 94% test coverage. Incremental review done.",
  "nx": ["tech_lead"],
  "artifacts": [
    "src/domain/entities/notification.py",
    "src/domain/entities/user_preferences.py",
    "src/application/services/preferences_service.py",
    "src/presentation/controllers/preferences_controller.py"
  ]
}
```

## Agent 5: Tech Lead Review

### Session

```
Read consensus/messages/inbox/tech_lead/
Review Workstream 1 code quality
Approve or veto
```

### Output

```json
// consensus/messages/inbox/developer/2024-01-18-workstream1-approved.json
{
  "d": "2024-01-18",
  "r": "tech_lead",
  "st": "approved",
  "epic": "epic-notifications",
  "sm": "Workstream 1 approved. Minor suggestion: add validation for quiet hours overlap.",
  "nx": ["developer"],
  "artifacts": []
}
```

*Workstreams 2 and 3 follow same pattern...*

## Agent 6: QA

### Session

```
Read prompts/qa_prompt.md
Verify all acceptance criteria
Run full test suite
Check coverage
```

### Output: consensus/artifacts/test-report.md

```markdown
# Test Report: User Notification System

## Summary
- Status: APPROVED
- Tests: 156 passed, 0 failed
- Coverage: 89%

## Acceptance Criteria Verification

### Preferences
- [x] User can enable/disable email notifications
- [x] User can enable/disable in-app notifications
- [x] User can enable/disable push notifications
- [x] User can set quiet hours
- [x] Preferences persist

### Delivery
- [x] Notifications route based on preferences
- [x] Critical notifications bypass quiet hours
- [x] Failed deliveries retry 3 times
- [x] Delivery status tracked

### History
- [x] User can view last 100 notifications
- [x] Filter by type works
- [x] Mark as read works
- [x] Mark all as read works

### Performance
- [x] History loads in <500ms (measured: 180ms avg)
- [x] 10,000 concurrent notifications handled

## Issues Found
None blocking.

## Recommendation
Ready for deployment.
```

## Agent 7: DevOps

### Output: consensus/artifacts/deployment.md

```markdown
# Deployment: User Notification System

## Changes
- New database migrations (2)
- New Redis keys for preference caching
- New queue workers (notification_worker x 3)
- New API endpoints (5)

## Deployment Steps
1. Run database migrations
2. Deploy queue workers
3. Deploy API changes
4. Verify health checks
5. Enable feature flag

## Rollback Plan
1. Disable feature flag
2. Rollback API deployment
3. Stop queue workers
4. Rollback migrations (if needed)

## Monitoring
- Alert: notification_delivery_failures > 1%
- Alert: notification_queue_depth > 10000
- Dashboard: notification system health
```

## Decision Log

```markdown
# Decision Log: Epic Notifications

## 2024-01-15 - Analyst
Completed requirements analysis. 4 feature areas: preferences,
delivery, history, admin broadcast.

## 2024-01-16 - Architect
Approved architecture. Async queue for delivery, Redis caching
for preferences. Created ADR-0008, ADR-0009.

## 2024-01-17 - Tech Lead
Created 3 workstreams with 10-day timeline.

## 2024-01-18 - Developer
Completed Workstream 1 (preferences). 94% coverage.

## 2024-01-18 - Tech Lead
Approved Workstream 1. Minor feedback incorporated.

...

## 2024-01-28 - QA
Full verification passed. 156 tests, 89% coverage.

## 2024-01-29 - DevOps
Deployment plan created with rollback procedure.
```

## Time

~10 days elapsed:
- Day 1: Analyst, Architect
- Day 2: Tech Lead planning
- Days 3-5: Workstream 1
- Days 6-8: Workstream 2
- Days 9-10: Workstream 3
- Day 10: QA, DevOps
