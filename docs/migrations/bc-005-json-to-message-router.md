# Breaking Change: ```

> **Migration Guide**: See [Breaking Changes Summary](breaking-changes.md) for all migrations

---

## Overview

**Versions**: v0.5.0

```

#### Timeline

- **Deprecated:** 2025-12-01 (v1.2)
- **Removed:** 2026-01-01 (v0.3.0)
- **Migration Support:** Ended 2026-03-01

---

### 5. JSON → Message Router

#### What Changed

The **JSON-based messaging** was replaced with **Message Router** system.

**OLD (JSON Messaging):**
```json
// consensus/messages/inbox/developer/message-001.json
{
  "from": "architect",
  "to": "developer",
  "subject": "API design clarification",
  "body": "Please use REST for now...",
  "timestamp": "2025-12-31T10:00:00Z"
}
```

**NEW (Message Router):**
```python
from sdp.unified.agent.router import SendMessageRouter, Message

router = SendMessageRouter()
router.send_message(Message(
    sender="architect",
    content="Please use REST for now...",
    recipient="developer",
))
```

#### Why It Changed

| Problem | Solution |
|---------|----------|
| Writing JSON files manually is error-prone | Python API with type hints |
| No message validation | Message schema enforced by code |
| Hard to send messages between agents | Router handles delivery |
| No notification system | Integrated with Telegram |

#### Migration Steps

**Step 1: Install Message Router**

```bash
# Already included in SDP v0.5.0
poetry install
```

**Step 2: Convert JSON Messages to Python**

```python
# OLD: consensus/messages/inbox/developer/message-001.json
# {
#   "from": "architect",
#   "to": "developer",
#   "subject": "API design",
#   "body": "Use REST...",
#   "timestamp": "2025-12-31T10:00:00Z"
# }

# NEW: Use Message Router
from sdp.unified.agent.router import SendMessageRouter, Message

router = SendMessageRouter()
router.send_message(Message(
    sender="architect",
    content="Use REST for API endpoints. GraphQL will be considered later.",
    recipient="developer",
))
```

**Step 3: Update Agent Prompts**

Old agents read JSON messages:

```markdown
# OLD: consensus/prompts/developer.md
## Context Files
Read consensus/messages/inbox/developer/ for messages
```

New agents use SendMessageRouter:

```markdown
# NEW: .claude/skills/build/SKILL.md
## Teammate Communication
Use SendMessage tool to communicate with other agents.
```

**Step 4: Remove Message JSON Files**

```bash
# After migrating to Message Router
rm -rf consensus/messages/

# Messages now sent via Python API
```

**Step 5: Enable Telegram Notifications (Optional)**

```bash
# .env
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id
```

```python
from sdp.unified.notifications.telegram import TelegramNotifier

notifier = TelegramNotifier(config=config)
notifier.send(Notification(
    type=NotificationType.SUCCESS,
    message="Workstream completed successfully",
))
```

#### Before/After Comparison

**OLD (JSON Messaging):**
```bash
# Architect sends message to developer
cat > consensus/messages/inbox/developer/msg-001.json << EOF
{
  "from": "architect",
  "to": "developer",
  "body": "Use REST, not GraphQL",
  "timestamp": "$(date -Iseconds)"
}
EOF

# Developer reads messages
ls consensus/messages/inbox/developer/
# msg-001.json
```

**NEW (Message Router):**
```python
# Architect sends message
router.send_message(Message(
    sender="architect",
    content="Use REST, not GraphQL",
