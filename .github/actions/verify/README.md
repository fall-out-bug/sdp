# SDP Verify Action

![SDP Verify](https://img.shields.io/badge/SDP-Verify-blue)
![Version](https://img.shields.io/badge/version-v1.0.0-green)

Run SDP verification gates on pull requests with evidence tracking.

## Features

- **Quality Gates**: Run types, tests, and coverage checks
- **Evidence Tracking**: Verify evidence chain integrity
- **PR Comments**: Post verification results directly to pull requests
- **Fast Execution**: < 2 minutes for average Go project
- **Flexible Configuration**: Customize gates and working directory

## Badges

![Verification Status](https://img.shields.io/badge/verification-passed-success)

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

**Detailed Metrics:**
- p50 latency: < 2 minutes
- p95 latency: < 3 minutes
- p99 latency: < 5 minutes
- Success rate: ≥ 99%

See [SLOS.md](SLOS.md) for complete Service Level Objectives and monitoring.

## Marketplace

This action is published to GitHub Marketplace as `fall-out-bug/sdp-verify-action`.

To use from marketplace:

```yaml
- uses: fall-out-bug/sdp-verify-action@v1
  with:
    gates: 'types,tests,coverage'
```

## GitLab CI

Equivalent GitLab CI template available at [`ci/gitlab/.gitlab-ci-sdp-verify.yml`](../../ci/gitlab/.gitlab-ci-sdp-verify.yml).

See [GitLab CI documentation](../../ci/gitlab/README.md) for usage.

## Branding

- **Name**: SDP Verify
- **Description**: Run SDP verification gates with evidence tracking
- **Category**: CI

## License

MIT
