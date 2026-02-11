# Bug Report Integration Guide

**Report bugs via specialist roles**

## Overview

Bug report integration enables agents to file bug reports directly to issue tracking systems (GitHub Issues, GitLab Issues, etc.).

## BugReporter Interface

```go
type BugReporter struct {
    issueTracker IssueTracker
    template     BugTemplate
}

func (br *BugReporter) FileBug(report BugReport) (string, error)
func (br *BugReporter) UpdateBug(bugID string, update string) error
func (br *BugReporter) CloseBug(bugID string, reason string) error
```

## BugReport Structure

```go
type BugReport struct {
    Title       string
    Description string
    Severity    string // P0, P1, P2, P3
    Type        string // bug, feature, enhancement
    Environment string
    Reproducible bool
    Steps       []string
    Logs        string
    Screenshots []string
}
```

## Workflow

```
┌──────────┐    ┌──────────┐    ┌──────────┐
│  Agent   │───→│ Bug      │───→│ GitHub   │
│  Detects │    │ Reporter  │    │ Issues   │
│  Issue   │    │          │    │          │
└──────────┘    └──────────┘    └──────────┘
                      │
                      ▼
               ┌──────────┐
               │ Bug ID   │
               │ Returned │
               └──────────┘
```

## Example Usage

```bash
# Agent detects bug during execution
@bugreport --title="Memory leak in user service" \
  --severity="P0" \
  --type="bug" \
  --steps="1. Run user service, 2. Observe memory growth"

# Bug filed automatically
# Output: Bug filed: https://github.com/org/repo/issues/123
```

## Integration Points

### GitHub Issues
```go
type GitHubIssueTracker struct {
    token   string
    repo    string // owner/repo
    baseURL string
}

func (gh *GitHubIssueTracker) CreateIssue(report BugReport) (string, error)
```

### GitLab Issues
```go
type GitLabIssueTracker struct {
    token   string
    project string // project ID or path
    baseURL string
}

func (gl *GitLabIssueTracker) CreateIssue(report BugReport) (string, error)
```

## Bug Template

```markdown
---
title: "{{ .Title }}"
severity: "{{ .Severity }}"
type: "{{ .Type }}"
---

## Description

{{ .Description }}

## Environment

- OS: {{ .Environment.OS }}
- Version: {{ .Environment.Version }}
- Architecture: {{ .Environment.Arch }}

## Steps to Reproduce

{{ range $i, $step := .Steps }}
{{ $i }}. {{ $step }}
{{ end }}

## Expected Behavior

What should have happened?

## Actual Behavior

What actually happened?

## Logs

```
{{ .Logs }}
```

## Screenshots

{{ range .Screenshots }}
![Screenshot]({{ . }})
{{ end }}
```

## Best Practices

1. **File bugs with context**: Include logs, environment info, steps to reproduce
2. **Use severity levels**: P0 (critical) → P3 (minor)
3. **Categorize correctly**: bug, feature, enhancement, documentation
4. **Provide minimal reproduction**: Smallest case that demonstrates the issue
5. **Monitor filed bugs**: Track status and updates

## Automation

Agents can automatically file bugs when:
- Tests fail consistently
- Quality gates reject code
- Performance degrades
- Security vulnerabilities detected
- Integration failures occur

## Configuration

```yaml
bug_tracking:
  enabled: true
  provider: github
  repo: owner/repo
  token_env: GITHUB_TOKEN
  auto_file: true
  severity_threshold: P1
```
