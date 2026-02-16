# F075: Self-Healing Doctor

**Status**: Completed
**Priority**: P0
**Workstreams**: 6/6 Complete

## Overview

Implements `sdp doctor --repair --deep` with guided fix-it output, config migrations, and automatic rollback support.

## Workstreams

| ID | Title | Status |
|----|-------|--------|
| 00-075-01 | Doctor Repair Mode | Completed |
| 00-075-02 | Deep Diagnostic Mode | Completed |
| 00-075-03 | Guided Fixes UX | Completed |
| 00-075-04 | Safe Config Migrations | Completed |
| 00-075-05 | Automatic Rollback on Failure | Completed |
| 00-075-06 | Doctor Integration Tests | Completed |

## Implementation Summary

### --repair Mode
- Auto-fixes file permissions on sensitive files (0600)
- Creates missing `.claude/` subdirectories (skills, agents, validators)
- Creates missing `.sdp/log` directory
- Reports repair status (fixed/skipped/manual/failed/partial)

### --deep Mode
- Git hooks integrity check (pre-commit, pre-push)
- Skill files syntax validation (YAML frontmatter)
- Workstream circular dependency detection
- Beads database integrity check
- Config version compatibility check
- Duration tracking for each check

### --migrate Mode
- Automatic config version detection
- Timestamped backups in `.sdp/backups/`
- Migration registry for version steps
- Migration logging to `.sdp/migrations.log`
- Dry-run support with `--dry-run`

### --rollback Mode
- Restore config from any backup file
- Clear error messages for missing backups
- List available backups helper function

## Quality Metrics

| Metric | Target | Actual |
|--------|--------|--------|
| Test Coverage | >= 80% | 85.1% |
| File Size | < 200 LOC | All files pass |
| Test Status | All Pass | PASS |

## Files

### Source Files
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_checks.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_repair.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_repair_perms.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_deep.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_deep_deps.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_migrate.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_migrate_ops.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_drift.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/cmd/sdp/doctor.go`

### Test Files
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_test.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_repair_test.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_deep_test.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_migrate_test.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_config_test.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_config_validation_test.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_drift_test.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_env_test.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/internal/doctor/doctor_extra_coverage_test.go`
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-F075-self-healing/sdp-plugin/cmd/sdp/doctor_test.go`

## CLI Usage

```bash
# Standard health check
sdp doctor

# Run with drift detection
sdp doctor --drift

# Auto-fix issues
sdp doctor --repair

# Comprehensive diagnostics
sdp doctor --deep

# Migrate config to latest version
sdp doctor --migrate

# Preview migration without changes
sdp doctor --migrate --dry-run

# Rollback to previous config
sdp doctor --rollback .sdp/backups/config-20260216-120000.yml

# Combined modes
sdp doctor --repair --deep
```
