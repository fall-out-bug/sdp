# Go Quick Start

SDP workflow for Go projects with go test, go tool cover, go vet, and golint.

## Prerequisites

```bash
# Go 1.21+
go version

# Install linting tools
go install golang.org/x/lint/golint@latest
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
```

## Project Structure

```
my-project/
├── go.mod
├── service/
│   └── user.go
├── service/
│   └── user_test.go
└── .claude/
    └── skills/
```

## Workflow

### 1. Initialize Project

```bash
# Create feature
@feature "Add user authentication"
```

Claude will interview you about:
- Mission and users
- Technical approach
- Success criteria
- Tradeoffs

### 2. Plan Workstreams

```bash
@design feature-auth
```

Claude will:
- Explore Go project structure
- Design workstream decomposition
- Define dependencies
- Request approval

### 3. Execute Workstream

```bash
@build 00-001-01
```

SDP will:
1. Detect Go project (go.mod found)
2. Run tests: `go test ./...`
3. Run coverage: `go test -coverprofile=coverage.out ./...`
4. Run vetting: `go vet ./...`
5. Run linting: `golint ./...`
6. AI validators check architecture, errors, complexity

### 4. Review Quality

```bash
@review feature-auth
```

SDP will run AI validators:
- `/coverage-validator` - Analyzes go tool cover output
- `/architecture-validator` - Checks import paths
- `/error-validator` - Finds ignored errors
- `/complexity-validator` - Identifies complex functions

## Quality Gates

### Test Coverage ≥80%

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
```

**Output:**
```
ok      github.com/user/my-project/service    0.123s   coverage: 85.4% of statements
ok      github.com/user/my-project/model     0.045s   coverage: 92.1% of statements

github.com/user/my-project/    87.3%  ✅ PASS (≥80%)
```

### Type Checking

```bash
go vet ./...
```

**Output:**
```
✅ PASS (no warnings)
```

### Linting

```bash
golint ./...
```

**Output:**
```
✅ PASS (no warnings)
```

### Complexity

```bash
gocyclo -over 10 .
```

**Output:**
```
✅ PASS (no functions >10 complexity)
```

### File Size

```bash
find . -name "*.go" -not -path "./vendor/*" -exec wc -l {} + | awk '$1 > 200'
```

**Output:**
```
(no output = all files <200 LOC) ✅ PASS
```

## Example Workflow

```bash
# Start feature
@feature "Add user login"

# Plan workstreams
@design feature-login

# Execute first workstream
@build 00-001-01

# Expected output:
# ✓ Project type detected: Go (go.mod)
# ✓ Running tests: go test ./...
# ✓ Coverage: 87% (≥80% required)
# ✓ Type checking: go vet ./...
# ✓ Linting: golint ./...
# ✓ AI validators: PASS
#
# Workstream 00-001-01 complete!

# Execute next workstream
@build 00-001-02

# Review all workstreams
@review feature-login

# Deploy
@deploy feature-login
```

## Test Example

```go
// service/user.go
package service

type User struct {
    Name string
}

func NewUser(name string) (*User, error) {
    if name == "" {
        return nil, fmt.Errorf("name cannot be empty")
    }
    return &User{Name: name}, nil
}

func (u *User) Greet() string {
    return "Hello, " + u.Name
}
```

```go
// service/user_test.go
package service

import "testing"

func TestNewUser(t *testing.T) {
    user, err := NewUser("Alice")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "Alice" {
        t.Errorf("expected Alice, got %s", user.Name)
    }
}

func TestNewUserEmpty(t *testing.T) {
    _, err := NewUser("")
    if err == nil {
        t.Fatal("expected error for empty name")
    }
}

func TestUserGreet(t *testing.T) {
    user := &User{Name: "Alice"}
    got := user.Greet()
    expected := "Hello, Alice"
    if got != expected {
        t.Errorf("expected %s, got %s", expected, got)
    }
}
```

## Common Issues

### Issue: Coverage <80%

**Solution:** Add more table-driven tests
```go
func TestNewUser(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid", "Alice", false},
        {"empty", "", true},           // Add this
        {"spaces", "   ", true},       // Add this
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewUser(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Issue: Ignored errors

**Solution:** Always check errors
```go
// Before (FAIL)
data, _ := fetchData()  // Error lost

// After (PASS)
data, err := fetchData()
if err != nil {
    return fmt.Errorf("fetch failed: %w", err)
}
```

### Issue: go vet warnings

**Solution:** Fix vetting issues
```bash
go vet ./...
# Review warnings and fix
```

## Tips

1. **Run tests with race detection:**
   ```bash
   go test -race ./...
   ```

2. **Run specific test:**
   ```bash
   go test -v -run TestNewUser ./service
   ```

3. **Generate coverage HTML:**
   ```bash
   go tool cover -html=coverage.out
   open coverage.html
   ```

4. **Benchmark tests:**
   ```go
   func BenchmarkUserGreet(b *testing.B) {
       user := &User{Name: "Alice"}
       for i := 0; i < b.N; i++ {
           user.Greet()
       }
   }
   ```

   ```bash
   go test -bench=. -benchmem
   ```

5. **Format code:**
   ```bash
   go fmt ./...
   ```

6. **Pre-commit hooks:**
   ```bash
   # .git/hooks/pre-commit
   go test ./...
   go vet ./...
   go fmt ./...
   ```

## Module Example

**go.mod:**
```go
module github.com/user/my-project

go 1.21

require (
    github.com/stretchr/testify v1.8.4 // indirect
)
```

**Import in tests:**
```go
package service_test

import (
    "testing"

    "github.com/user/my-project/service"
    "github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
    user, err := service.NewUser("Alice")
    assert.NoError(t, err)
    assert.Equal(t, "Alice", user.Name)
}
```

## Next Steps

- [Python Quick Start](../python/QUICKSTART.md)
- [Java Quick Start](../java/QUICKSTART.md)
- [Full Tutorial](../../TUTORIAL.md)
