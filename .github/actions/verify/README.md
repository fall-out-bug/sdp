# SDP Verify Action

Run SDP verification gates on pull requests with evidence tracking.

## Usage

### Basic Example

```yaml
name: SDP Verify

on:
  pull_request:
    branches: [main]

jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: ./.github/actions/verify
        with:
          gates: 'types,tests,coverage'
          evidence-required: 'true'
```

### Inputs

| Input | Required | Default | Description |
|-------|----------|---------|-------------|
| `gates` | No | `types,tests,coverage` | Comma-separated list of gates to run |
| `evidence-required` | No | `true` | Whether to check evidence chain integrity |
| `comment` | No | `true` | Whether to post PR comment with results |
| `version` | No | `latest` | SDP CLI version to install |
| `working-directory` | No | `.` | Working directory for SDP checks |

### Available Gates

- **types** - Type checking (go vet, mypy, etc.)
- **tests** - Run all tests
- **coverage** - Test coverage analysis (≥80% required)
- **evidence** - Evidence chain integrity check

### Outputs

| Output | Description |
|--------|-------------|
| `result` | Verification result (`pass`/`fail`) |
| `gates_passed` | Number of gates that passed |
| `gates_failed` | Number of gates that failed |

### Example with Go Project

```yaml
jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - uses: ./.github/actions/verify
        with:
          gates: 'types,tests,coverage'
          working-directory: './my-go-app'
```

## Exit Codes

- `0` - All gates passed
- `1` - One or more gates failed

## Performance

Execution time: < 2 minutes for average Go project

## License

MIT
