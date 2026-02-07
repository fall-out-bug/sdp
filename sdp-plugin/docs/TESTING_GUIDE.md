# Testing Guide for F024

**Comprehensive testing strategy**

## Test Pyramid

```
        ┌─────────┐
        │   E2E   │  (WS-022, WS-023)
        │  10%   │
        ├─────────┤
        │Integration│ (WS-021)
        │  30%   │
        ├─────────┤
        │  Unit   │  (WS-020)
        │  60%   │
        └─────────┘
```

## Unit Tests (WS-020)

### Coverage Requirements

- **Minimum:** 80% per package
- **Target:** 90% per package
- **Measured:** `go test -coverprofile=coverage.out`

### Test Structure

```go
// Internal orchestrator tests
func TestFeatureCoordinator_ExecuteFeature_Success(t *testing.T) {
    // Arrange
    loader := &MockWorkstreamLoader{...}
    executor := &MockWorkstreamExecutor{...}
    coordinator := NewFeatureCoordinator(loader, executor, nil, 2)

    // Act
    err := coordinator.ExecuteFeature("F001")

    // Assert
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}
```

### Test Categories

1. **Happy Path:** Normal execution
2. **Error Cases:** Invalid inputs, failures
3. **Edge Cases:** Empty data, boundaries
4. **Concurrency:** Race conditions, deadlocks

### Running Tests

```bash
# Run all unit tests
go test ./internal/... -v

# Run with coverage
go test ./internal/... -coverprofile=coverage.out

# Check coverage threshold
go tool cover -func=coverage.out | grep total
```

## Integration Tests (WS-021)

### Components Tested

1. **Orchestrator + Beads:** Workstream execution with task tracking
2. **Skill Invocation:** @idea, @design, @oneshot calls
3. **Agent Spawning:** Multi-agent coordination
4. **Message Routing:** Agent-orchestrator communication
5. **Role Loading:** Role file parsing and caching

### Example Test

```go
func TestIntegration_OrchestratorWithBeads(t *testing.T) {
    // Setup: Initialize Beads
    bd := setupTestBeads(t)
    defer cleanupTestBeads(bd)

    // Execute feature
    orchestrator := NewOrchestrator(...)
    err := orchestrator.ExecuteFeature("F001")

    // Verify: Beads tasks created/closed
    tasks := bd.ListTasks("F001")
    if len(tasks) == 0 {
        t.Error("Expected Beads tasks to be created")
    }
}
```

### Running Integration Tests

```bash
# Run integration tests
go test ./tests/integration/... -v

# With test database
TEST_DATABASE_URL="postgres://localhost/test" go test ./tests/integration/...
```

## E2E Tests with Beads (WS-022)

### Test Scenarios

1. **Full Feature Workflow:**
   - @feature → @design → @oneshot → review
   - Verify all workstreams execute
   - Verify Beads integration

2. **Checkpoint/Resume:**
   - Start feature execution
   - Interrupt mid-execution
   - Resume from checkpoint
   - Verify continuation

3. **Quality Gates:**
   - Execute workstream with failing tests
   - Verify rejection
   - Fix and retry
   - Verify acceptance

### Example Test

```go
func TestE2E_FullWorkflow(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }

    // Setup: Clean workspace, init Beads
    workspace := setupWorkspace(t)
    defer cleanupWorkspace(workspace)

    // Execute: Run full workflow
    result := runSDPCommand(t, "@feature", "Add user login --execute")

    // Verify: Feature complete
    if result.ExitCode != 0 {
        t.Errorf("Expected success, got exit code %d", result.ExitCode)
    }

    // Verify: All files created
    assertFileExists(t, "PRODUCT_VISION.md")
    assertFileExists(t, "docs/drafts/idea-user-login.md")
    assertWorkstreamsComplete(t, "F001")
}
```

### Running E2E Tests

```bash
# Run all E2E tests
go test ./tests/e2e/... -v -timeout=30m

# Run specific E2E test
go test ./tests/e2e/... -v -run TestE2E_FullWorkflow

# With Beads
BEADS_ENABLED=true go test ./tests/e2e/... -v
```

## E2E Tests with Telegram (WS-023)

### Prerequisites

- Telegram bot token (test environment)
- Test chat ID
- Test bot configured

### Test Scenarios

1. **Notification Delivery:**
   - Trigger critical notification
   - Verify Telegram message received
   - Verify message format

2. **Interactive Commands:**
   - Send command via Telegram
   - Verify agent execution
   - Verify response

### Example Test

```go
func TestE2E_TelegramNotification(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping Telegram test in short mode")
    }

    if telegramToken := os.Getenv("TELEGRAM_TEST_TOKEN"); telegramToken == "" {
        t.Skip("TELEGRAM_TEST_TOKEN not set")
    }

    // Setup: Create Telegram client
    bot := setupTestBot(t)

    // Execute: Trigger notification
    router := NewNotificationRouter()
    router.RegisterProvider(NewTelegramProvider(telegramToken, testChatID))

    msg := Notification{
        Title: "Test Notification",
        Body: "E2E test message",
        Severity: "info",
    }

    err := router.Send(context.Background(), msg)

    // Verify: No error
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    // Verify: Message received (manual check or Telegram API)
}
```

### Running Telegram E2E Tests

```bash
# Set test credentials
export TELEGRAM_TEST_TOKEN="your-test-bot-token"
export TELEGRAM_TEST_CHAT_ID="your-test-chat-id"

# Run Telegram E2E tests
go test ./tests/e2e/telegram/... -v -timeout=10m

# Dry run (no actual messages sent)
TELEGRAM_DRY_RUN=true go test ./tests/e2e/telegram/... -v
```

## Test Quality Gates

### Pre-Commit

```bash
# Run in pre-commit hook
go test ./... -short
go vet ./...
golint ./...
```

### Pre-Merge

```bash
# Full test suite
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | awk '/total/ && ($3+0) < 80 {exit 1}'
```

### Pre-Deploy

```bash
# E2E tests
go test ./tests/e2e/... -v -timeout=30m
```

## Test Data Management

### Fixtures

```go
// tests/fixtures/feature_spec.go
package fixtures

var TestFeature = Feature{
    ID: "F999",
    Name: "Test Feature",
    Workstreams: []Workstream{
        {ID: "WS-999-01", Name: "Test WS 1"},
        {ID: "WS-999-02", Name: "Test WS 2"},
    },
}
```

### Test Database

```bash
# Setup test database
docker run -d --name test-postgres \
  -e POSTGRES_PASSWORD=test \
  -e POSTGRES_DB=test \
  -p 5433:5432 postgres:15

export TEST_DATABASE_URL="postgres://postgres:test@localhost:5433/test?sslmode=disable"
```

## Continuous Integration

### GitHub Actions

```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run unit tests
        run: go test ./internal/... -coverprofile=coverage.out

      - name: Check coverage
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$coverage < 80" | bc -l) )); then
            echo "Coverage ${coverage}% is below 80%"
            exit 1
          fi

      - name: Run integration tests
        run: go test ./tests/integration/... -v

      - name: Run E2E tests
        if: github.ref == 'refs/heads/main'
        run: go test ./tests/e2e/... -v -timeout=30m
```

## Best Practices

1. **Test isolation:** Each test should be independent
2. **Use table-driven tests:** Test multiple cases efficiently
3. **Mock external dependencies:** Don't rely on external services
4. **Clean up resources:** Use defer for cleanup
5. **Use t.Parallel():** Run independent tests in parallel
6. **Descriptive names:** Test names should describe what they test
7. **Assert messages:** Include context in assertions
8. **Subtests:** Use t.Run() for related test cases

## Debugging Tests

### Verbose Output

```bash
go test -v ./internal/orchestrator/...
```

### Race Detection

```bash
go test -race ./internal/...
```

### Benchmark

```bash
go test -bench=. -benchmem ./internal/...
```

### Coverage HTML Report

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```
