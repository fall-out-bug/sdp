# Telegram Notifications

Automated Telegram notifications for critical workflow events.

## Setup

1. **Create Telegram Bot**

```bash
# Talk to @BotFather on Telegram
/newbot
# Follow instructions, save bot token
```

2. **Get Chat ID**

```bash
# Start chat with your bot
# Send any message
# Then run:
curl https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates

# Find "chat":{"id":12345678} in response
```

3. **Set Environment Variables**

```bash
export TELEGRAM_BOT_TOKEN="1234567890:ABCdefGHIjklMNOpqrsTUVwxyz"
export TELEGRAM_CHAT_ID="12345678"
```

Or add to `~/.bashrc` / `~/.zshrc`:

```bash
echo 'export TELEGRAM_BOT_TOKEN="..."' >> ~/.bashrc
echo 'export TELEGRAM_CHAT_ID="..."' >> ~/.bashrc
source ~/.bashrc
```

## Events

### Oneshot Lifecycle

| Event | When | Severity |
|-------|------|----------|
| `oneshot_started` | `/oneshot` begins | INFO |
| `oneshot_completed` | All WS done | SUCCESS |
| `oneshot_blocked` | Critical failure | CRITICAL |

### Workstream Events

| Event | When | Severity |
|-------|------|----------|
| `ws_failed` | WS execution fails | WARNING |

### Review Events

| Event | When | Severity |
|-------|------|----------|
| `review_failed` | Review verdict: CHANGES_REQUESTED | WARNING |

### Deployment Events

| Event | When | Severity |
|-------|------|----------|
| `breaking_changes` | Breaking changes detected | WARNING |
| `e2e_failed` | E2E tests fail before deploy | CRITICAL |
| `deploy_success` | Deployment complete | SUCCESS |

### Hotfix Events

| Event | When | Severity |
|-------|------|----------|
| `hotfix_deployed` | Hotfix deployed to production | CRITICAL |

## Usage

### Manual Testing

```bash
# Test notification
bash sdp/notifications/telegram.sh oneshot_started "F60" "4"

# Test blocking notification
bash sdp/notifications/telegram.sh oneshot_blocked "F60" "WS-060-02" "Coverage < 80%"
```

### In Commands

Notifications are automatically sent from:
- `/oneshot` (start, complete, blocked)
- `/codereview` (failed)
- `/deploy` (success)
- `/hotfix` (deployed)

### In Hooks

Notifications are sent from:
- `pre-commit.sh` (breaking changes)
- `pre-deploy.sh` (e2e failed)

## Notification Format

### oneshot_started
```
🚀 *Oneshot Started*

Feature: `F60`
Workstreams: 4
Status: Executing autonomously

Branch: `feature/lms-integration`
```

### oneshot_completed
```
✅ *Oneshot Completed*

Feature: `F60`
Elapsed (telemetry): 2h 15m
Status: All WS executed successfully

Ready for review: `/codereview F60`
```

### oneshot_blocked (CRITICAL)
```
🚨 *CRITICAL: Oneshot Blocked*

Feature: `F60`
Blocked at: `WS-060-02`
Reason: Coverage < 80%

⚠️ *HUMAN INTERVENTION REQUIRED*

To resume: `/oneshot F60 --resume`
```

### review_failed
```
⚠️ *Review Failed*

Feature: `F60`
Issues: 3

Status: CHANGES_REQUESTED

Check: `tools/hw_checker/docs/workstreams/reports/F60-review.md`
```

### breaking_changes
```
⚠️ *Breaking Changes Detected*

Changes: 3

Action required:
1. Review `BREAKING_CHANGES.md`
2. Complete `MIGRATION_GUIDE.md`
3. Add both to commit
```

### e2e_failed (CRITICAL)
```
🚨 *E2E Tests Failed*

Feature: `F60`
Failed tests: 2

⛔ Deployment blocked until fixed
```

### deploy_success
```
🎉 *Deployment Successful*

Feature: `F60`
Environment: production
Version: v1.5.0

Status: Live
```

### hotfix_deployed (CRITICAL)
```
🔥 *Hotfix Deployed*

Issue: `001`
Elapsed (telemetry): 30m
Target: Production

Status: Emergency fix live
```

## Disable Notifications

If `TELEGRAM_BOT_TOKEN` or `TELEGRAM_CHAT_ID` are not set, notifications are silently skipped.

```bash
# Temporarily disable
unset TELEGRAM_BOT_TOKEN
unset TELEGRAM_CHAT_ID
```

## Troubleshooting

### No notifications received

```bash
# Check env vars
echo $TELEGRAM_BOT_TOKEN
echo $TELEGRAM_CHAT_ID

# Test bot manually
curl -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
  -d "chat_id=${TELEGRAM_CHAT_ID}" \
  -d "text=Test message"
```

### Wrong chat

```bash
# Get all chats bot has access to
curl https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/getUpdates

# Update TELEGRAM_CHAT_ID
export TELEGRAM_CHAT_ID="new_id"
```

### Rate limits

Telegram limits: 30 messages/second per bot.

If you hit limits:
- Reduce notification frequency
- Batch multiple events into one message
- Use Telegram Bot API business plan

## Multiple Recipients

Send to multiple chats:

```bash
# telegram.sh
for CHAT_ID in "$TELEGRAM_CHAT_ID_1" "$TELEGRAM_CHAT_ID_2"; do
  curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
    -d "chat_id=${CHAT_ID}" \
    -d "text=${message}" \
    > /dev/null 2>&1
done
```

## Security

- **Never commit** bot tokens to Git
- Use environment variables or secrets management
- Restrict bot permissions (only send messages)
- Regularly rotate bot token
- Monitor bot activity in @BotFather

## Alternative: Telegram Channels

For team notifications, use channels:

1. Create channel in Telegram
2. Add bot as admin
3. Get channel ID (starts with `-100`)
4. Use channel ID as `TELEGRAM_CHAT_ID`

## Integration with Other Tools

### Slack (alternative)

Replace telegram.sh with slack.sh:

```bash
curl -X POST https://hooks.slack.com/services/YOUR/WEBHOOK/URL \
  -H 'Content-Type: application/json' \
  -d "{\"text\":\"${message}\"}"
```

### Discord (alternative)

```bash
curl -X POST https://discord.com/api/webhooks/YOUR/WEBHOOK \
  -H 'Content-Type: application/json' \
  -d "{\"content\":\"${message}\"}"
```
