# Pipeline Hooks Security

This guide explains how SDP executes pipeline hooks safely and how to write hook commands that pass security validation.

## Why This Exists

Pipeline hooks run inside CI and local automation. If hook commands are interpreted by a shell, malformed or malicious input can execute unintended commands.

SDP uses a fail-closed model:

- invalid hook command -> hook execution is rejected
- `on_fail: halt` -> phase fails immediately
- `on_fail: warn` or `ignore` -> behavior remains explicit and auditable

## Execution Model

Hooks are loaded from `.sdp/pipeline-hooks.yaml`.

Each hook includes:

- `phase`: `build`, `review`, or `ci`
- `when`: `pre` or `post`
- `command`: executable plus arguments
- `on_fail`: `halt`, `warn`, or `ignore`
- `timeout`: optional seconds (default `60`)

Runtime environment variables:

- `WS_ID`
- `FEATURE_ID`
- `PHASE`
- `CHECKPOINT_PATH`

## Command Validation Rules

### 1) No shell metacharacters

The following characters are rejected in `command`:

`; | & \` $ < > ( )` and newline/carriage-return.

This blocks command chaining and shell injection patterns.

### 2) Allowlisted executables only

Allowed direct executables:

- `bd`, `echo`, `false`, `git`, `go`, `make`
- `notify`, `slack-notify`
- `sdp`, `sdp-doc-sync`, `sdp-evidence`, `sdp-protocol-check`
- `trivy`, `true`

### 3) Local script paths are allowed, but sandboxed

Commands using `./` or `/` are allowed only if the resolved path stays inside the project root.

Examples:

- allowed: `hooks/pre-build.sh`
- allowed: `./scripts/check.sh`
- rejected: `../../outside.sh`

## Safe Configuration Examples

```yaml
hooks:
  - phase: build
    when: pre
    command: hooks/pre-build.sh
    on_fail: halt
    timeout: 30

  - phase: review
    when: post
    command: sdp-evidence validate .sdp/evidence/run.json
    on_fail: warn
```

## Migration From Shell-Style Commands

Replace shell forms with explicit commands:

- bad: `sh -c "echo ok && make test"`
- good: split into two hooks, one command per hook

- bad: `echo ok; rm -rf /tmp/x`
- good: `echo ok` (and perform cleanup in a vetted script under repo)

If you need complex logic, put it in a repository script (for example `hooks/pre-build.sh`) and call that script directly.

## Troubleshooting

### "command contains disallowed shell metacharacters"

Remove shell operators and split logic into multiple hooks or a script file.

### "command \"...\" is not in allowlist"

Use an allowlisted executable or call a script inside the repository.

### Hook times out

Increase `timeout` or optimize the called command.
