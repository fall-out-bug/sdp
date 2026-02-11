# WS-019: TelegramNotifier + Mock Provider

> **Workstream ID:** WS-019  
> **Feature:** F024: Unified Workflow  
> **Dependencies:** WS-018

## Goal

Telegram bot integration for notifications.

## Acceptance Criteria

### AC1: Telegram Integration
- [ ] Sends formatted messages to channel/group
- [ ] Uses Telegram Bot API
- [ ] Handles rate limits

### AC2: Message Format
- [ ] Messages formatted with Markdown
- [ ] Includes: feature, workstream, status, duration
- [ ] Escape special characters

### AC3: Error Handling
- [ ] Send errors logged
- [ ] Doesn't crash on API failures
- [ ] Retry on transient errors

### AC4: Mock Provider
- [ ] Mock provider for unit tests
- [ ] Captures messages in memory
- [ ] Returns instantly

## Scope Files

**internal/orchestrator/telegram_notifier.go** (NEW)
**internal/orchestrator/mock_notifier.go** (NEW)

## Implementation Details

```go
type TelegramNotifier struct {
    botToken string
    chatID   string
    client   *http.Client
}

func (tn *TelegramNotifier) Send(msg Message) error {
    // Send to Telegram Bot API
    url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tn.botToken)
    // POST request with message
}
```

## Estimated Scope

- ~250 LOC
- Duration: 2.5 hours
