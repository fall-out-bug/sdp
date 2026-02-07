# Notification System Guide

**Team notifications via multiple providers**

## Overview

The notification system enables agents to send notifications to team members via various providers (Telegram, Slack, Email, etc.).

## Architecture

```
┌─────────────┐
│   Agent     │
└──────┬──────┘
       │
       ▼
┌─────────────┐    ┌──────────────┐
│ Notification│───→│ Provider     │
│   Router    │    │ (Telegram)   │
└─────────────┘    └──────────────┘
                          │
                          ▼
                   ┌──────────────┐
                   │  Team Chat   │
                   └──────────────┘
```

## Components

### 1. Notification Provider

Interface for notification providers:

```go
type NotificationProvider interface {
    Send(ctx context.Context, msg Notification) error
    Name() string
    Validate(msg Notification) error
}
```

### 2. Notification Router

Routes notifications to appropriate providers:

```go
type NotificationRouter struct {
    providers map[string]NotificationProvider
    rules     []RoutingRule
}

func (nr *NotificationRouter) Send(ctx context.Context, msg Notification) error
func (nr *NotificationRouter) RegisterProvider(provider NotificationProvider)
func (nr *NotificationRouter) AddRule(rule RoutingRule)
```

### 3. Notification Message

```go
type Notification struct {
    Title     string
    Body      string
    Severity  string // info, warning, error, critical
    Tags      []string
    Metadata  map[string]string
    Timestamp time.Time
}
```

## Providers

### Telegram Provider

```go
type TelegramProvider struct {
    botToken string
    chatID   string
}

func (tp *TelegramProvider) Send(ctx context.Context, msg Notification) error {
    // Send via Telegram Bot API
    url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tp.botToken)
    payload := map[string]interface{}{
        "chat_id": tp.chatID,
        "text":    formatMessage(msg),
        "parse_mode": "Markdown",
    }
    // HTTP POST to Telegram API
    return nil
}
```

### Slack Provider

```go
type SlackProvider struct {
    webhookURL string
    channel    string
}

func (sp *SlackProvider) Send(ctx context.Context, msg Notification) error {
    // Send via Slack Incoming Webhook
    payload := map[string]interface{}{
        "channel":  sp.channel,
        "username": "SDP Bot",
        "text":     formatMessage(msg),
        "attachments": []Attachment{...},
    }
    // HTTP POST to Slack webhook
    return nil
}
```

### Email Provider

```go
type EmailProvider struct {
    smtpHost string
    smtpPort int
    username string
    password string
    from     string
    to       []string
}

func (ep *EmailProvider) Send(ctx context.Context, msg Notification) error {
    // Send via SMTP
    return nil
}
```

## Routing Rules

```go
type RoutingRule struct {
    Match    func(Notification) bool
    Provider string
}

// Example rules
rules := []RoutingRule{
    {
        Match: func(n Notification) bool {
            return n.Severity == "critical"
        },
        Provider: "telegram", // Critical → Telegram (immediate)
    },
    {
        Match: func(n Notification) bool {
            return n.Severity == "info"
        },
        Provider: "slack", // Info → Slack (batched)
    },
}
```

## Usage Examples

### Send Notification

```bash
# Agent sends notification
@notify --title="WS-001 complete" \
  --body="Domain models implemented with 92% coverage" \
  --severity="info" \
  --tags=["feature-complete", "ws-001"]
```

### Multi-Provider Routing

```yaml
notification_rules:
  - match:
      severity: critical
    providers:
      - telegram
      - slack
      - email

  - match:
      severity: info
    providers:
      - slack

  - match:
      tags: ["deploy"]
    providers:
      - slack
      - email
```

## Notification Templates

### Feature Complete

```markdown
✅ **Feature Complete**

**Feature:** {{ .FeatureID }}
**Workstreams:** {{ .CompletedWS }}/{{ .TotalWS }}
**Duration:** {{ .Duration }}
**Coverage:** {{ .Coverage }}%

**Summary:** {{ .Summary }}

**Next Steps:** {{ .NextSteps }}
```

### Error Notification

```markdown
🚨 **Error Detected**

**Workstream:** {{ .WorkstreamID }}
**Error:** {{ .Error }}
**Severity:** {{ .Severity }}

**Context:**
```
{{ .StackTrace }}
```

**Action Required:** {{ .Action }}
```

### Deployment Notification

```markdown
🚀 **Deployment Complete**

**Environment:** {{ .Env }}
**Version:** {{ .Version }}
**Commit:** {{ .Commit }}

**Changes:**
{{ range .Changes }}
- {{ . }}
{{ end }}

**Monitor:** {{ .MonitorURL }}
```

## Configuration

```yaml
notifications:
  enabled: true

  providers:
    telegram:
      bot_token: "${TELEGRAM_BOT_TOKEN}"
      chat_id: "${TELEGRAM_CHAT_ID}"
      enabled: true

    slack:
      webhook_url: "${SLACK_WEBHOOK_URL}"
      channel: "#dev-notifications"
      enabled: true

    email:
      smtp_host: "smtp.gmail.com"
      smtp_port: 587
      from: "sdp-bot@example.com"
      to: ["team@example.com"]
      enabled: false

  routing:
    critical:
      providers: ["telegram", "slack"]
      immediate: true

    error:
      providers: ["slack"]
      immediate: true

    warning:
      providers: ["slack"]
      immediate: false

    info:
      providers: ["slack"]
      immediate: false
```

## Best Practices

1. **Don't spam**: Batch non-critical notifications
2. **Use severity correctly**: Critical = immediate, Info = batched
3. **Provide context**: Include relevant information (WS ID, error, stack trace)
4. **Actionable messages**: Include what needs to be done
5. **Tag notifications**: Enable filtering and routing

## Testing

```go
func TestNotificationProvider(t *testing.T) {
    provider := NewMockNotificationProvider()

    msg := Notification{
        Title: "Test",
        Body:  "Test notification",
        Severity: "info",
    }

    err := provider.Send(context.Background(), msg)

    if err != nil {
        t.Fatalf("Failed to send: %v", err)
    }
}
```

## Security Considerations

1. **Encrypt credentials**: Store tokens in encrypted vault
2. **Use environment variables**: Don't hardcode tokens
3. **Validate recipients**: Prevent notification spam
4. **Rate limiting**: Prevent abuse
5. **Audit logging**: Log all notifications sent
