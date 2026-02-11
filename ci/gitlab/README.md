# GitLab CI Templates for SDP

This directory contains GitLab CI templates for running SDP verification gates.

## .gitlab-ci-sdp-verify.yml

Runs SDP verification gates on merge requests.

### Usage

#### Option 1: Include from SDP Repository

```yaml
# .gitlab-ci.yml
include:
  - project: 'fall-out-bug/sdp'
    ref: main
    file: 'ci/gitlab/.gitlab-ci-sdp-verify.yml'

# Customize variables
variables:
  SDP_GATES: "types,tests,coverage"
  SDP_WORKING_DIR: "./my-app"
```

#### Option 2: Copy Template

Copy `.gitlab-ci-sdp-verify.yml` to your repository and customize:

```yaml
# .gitlab-ci.yml
include:
  - local: '.gitlab-ci-sdp-verify.yml'

variables:
  SDP_VERSION: "latest"
  SDP_GATES: "types,tests,coverage"
  SDP_EVIDENCE_REQUIRED: "true"
  SDP_WORKING_DIR: "."
```

### Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SDP_VERSION` | `latest` | SDP CLI version to install |
| `SDP_GATES` | `types,tests,coverage` | Comma-separated list of gates |
| `SDP_EVIDENCE_REQUIRED` | `true` | Check evidence chain integrity |
| `SDP_WORKING_DIR` | `.` | Working directory for checks |

### MR Comments

To enable MR comments, set `GITLAB_TOKEN` in your project CI/CD variables:

1. Go to **Settings > CI/CD > Variables**
2. Add variable `GITLAB_TOKEN` with `api` scope
3. Token can be a personal access token or project access token

### Available Gates

- **types** - Type checking (go vet, mypy, etc.)
- **tests** - Run all tests
- **coverage** - Test coverage analysis (≥80% required)
- **evidence** - Evidence chain integrity check

### Example: Monorepo

```yaml
variables:
  SDP_WORKING_DIR: "./services/backend"

sdp_verify_backend:
  extends: .sdp_verify
  only:
    variables:
      - $CI_MERGE_REQUEST_TARGET_BRANCH_NAME == "main"
```

### Example: Multiple Services

```yaml
sdp_verify_api:
  extends: .sdp_verify
  variables:
    SDP_WORKING_DIR: "./services/api"

sdp_verify_worker:
  extends: .sdp_verify
  variables:
    SDP_WORKING_DIR: "./services/worker"
```

## Requirements

- Runner with shell executor
- Go 1.23+ (for `go test` gate)
- Internet access (to download SDP CLI)

## Performance

Execution time: < 2 minutes for average Go project

## Troubleshooting

### Permission Denied

```bash
chmod +x /usr/local/bin/sdp
```

### Unsupported Architecture

The template supports `amd64` and `arm64`. For other architectures, modify the download URL.

### Go Module Cache

To enable Go module caching:

```yaml
cache:
  paths:
    - .go/pkg/mod/
```
