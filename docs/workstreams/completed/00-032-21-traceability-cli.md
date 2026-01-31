---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-20
feature: F032
github_issue: null
project_id: 0
size: MEDIUM
status: completed
traceability:
- ac_description: '`sdp trace check <ws_id>` показывает mapping table'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: '`sdp trace add <ws_id> --ac AC1 --test test_func` добавляет mapping'
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: Exit code 1 если есть unmapped ACs
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_ws_tracker.py
  test_name: test_cannot_activate_second_ws
- ac_description: Unit tests (≥80%)
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: '`--json` flag для machine-readable output'
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_ws_tracker.py
  test_name: test_get_active_returns_ws_id
ws_id: 00-032-21
---

## 00-032-21: Traceability CLI

### 🎯 Goal

**What must WORK after completing this WS:**
- CLI для проверки и создания AC→Test mappings
- Exit code 1 если есть unmapped ACs
- JSON output для CI integration

**Acceptance Criteria:**
- [ ] AC1: `sdp trace check <ws_id>` показывает mapping table
- [ ] AC2: `sdp trace add <ws_id> --ac AC1 --test test_func` добавляет mapping
- [ ] AC3: Exit code 1 если есть unmapped ACs
- [ ] AC4: `--json` flag для machine-readable output
- [ ] AC5: Unit tests (≥80%)

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Нет инструмента для проверки трассируемости.

**Solution**: CLI с human и JSON output.

### Dependencies

- **00-032-20**: AC Test Mapping Model

### Steps

1. **Create TraceabilityService**

   ```python
   # src/sdp/traceability/service.py
   from pathlib import Path
   import re
   
   from sdp.traceability.models import (
       ACTestMapping,
       TraceabilityReport,
       MappingStatus,
   )
   from sdp.beads import BeadsClient
   
   class TraceabilityService:
       """Check and manage AC→Test traceability."""
       
       def __init__(self, client: BeadsClient):
           self._client = client
       
       def check_traceability(self, ws_id: str) -> TraceabilityReport:
           """Check traceability for workstream."""
           # Get WS from Beads
           task = self._client.get_task(ws_id)
           if not task:
               raise ValueError(f"WS not found: {ws_id}")
           
           # Extract ACs from description
           acs = self._extract_acs(task.description)
           
           # Get existing mappings
           stored_mappings = task.sdp_metadata.get("traceability", [])
           
           # Build report
           mappings = []
           for ac_id, ac_desc in acs:
               # Find stored mapping
               stored = next(
                   (m for m in stored_mappings if m["ac_id"] == ac_id),
                   None
               )
               
               if stored:
                   mappings.append(ACTestMapping.from_dict(stored))
               else:
                   mappings.append(ACTestMapping(
                       ac_id=ac_id,
                       ac_description=ac_desc,
                       test_file=None,
                       test_name=None,
                       status=MappingStatus.MISSING,
                   ))
           
           return TraceabilityReport(ws_id=ws_id, mappings=mappings)
       
       def add_mapping(
           self,
           ws_id: str,
           ac_id: str,
           test_file: str,
           test_name: str
       ) -> None:
           """Add AC→Test mapping."""
           task = self._client.get_task(ws_id)
           if not task:
               raise ValueError(f"WS not found: {ws_id}")
           
           # Get current mappings
           metadata = task.sdp_metadata.copy()
           mappings = metadata.get("traceability", [])
           
           # Update or add
           existing = next((m for m in mappings if m["ac_id"] == ac_id), None)
           if existing:
               existing["test_file"] = test_file
               existing["test_name"] = test_name
               existing["status"] = "mapped"
           else:
               mappings.append({
                   "ac_id": ac_id,
                   "ac_description": "",  # Will be filled from WS
                   "test_file": test_file,
                   "test_name": test_name,
                   "status": "mapped",
               })
           
           metadata["traceability"] = mappings
           self._client.update_metadata(ws_id, metadata)
       
       def _extract_acs(self, description: str) -> list[tuple[str, str]]:
           """Extract ACs from WS description.
           
           Looks for patterns like:
           - [ ] AC1: Description
           - AC1: Description
           """
           pattern = r'(?:- \[[ x]\] )?(AC\d+):\s*(.+?)(?:\n|$)'
           matches = re.findall(pattern, description, re.IGNORECASE)
           return [(m[0].upper(), m[1].strip()) for m in matches]
   ```

2. **Create CLI commands**

   ```python
   # src/sdp/cli/trace.py
   import typer
   import json
   
   from sdp.beads import create_beads_client
   from sdp.traceability.service import TraceabilityService
   
   app = typer.Typer(help="Traceability commands")
   
   @app.command("check")
   def check_traceability(
       ws_id: str,
       json_output: bool = typer.Option(False, "--json", help="JSON output")
   ) -> None:
       """Check AC→Test traceability for workstream."""
       service = TraceabilityService(create_beads_client())
       
       try:
           report = service.check_traceability(ws_id)
       except ValueError as e:
           typer.echo(f"❌ {e}")
           raise typer.Exit(1)
       
       if json_output:
           typer.echo(json.dumps(report.to_dict(), indent=2))
       else:
           typer.echo(f"\nTraceability Report: {ws_id}")
           typer.echo("=" * 50)
           typer.echo(report.to_markdown_table())
           typer.echo("")
           typer.echo(f"Coverage: {report.coverage_pct:.0f}% ({report.mapped_acs}/{report.total_acs} ACs mapped)")
           
           if report.is_complete:
               typer.echo("Status: ✅ COMPLETE")
           else:
               typer.echo(f"Status: ❌ INCOMPLETE ({report.missing_acs} unmapped)")
       
       # Exit 1 if incomplete
       if not report.is_complete:
           raise typer.Exit(1)
   
   @app.command("add")
   def add_mapping(
       ws_id: str,
       ac: str = typer.Option(..., "--ac", help="AC ID (e.g., AC1)"),
       test: str = typer.Option(..., "--test", help="Test function name"),
       file: str = typer.Option(None, "--file", help="Test file path")
   ) -> None:
       """Add AC→Test mapping."""
       service = TraceabilityService(create_beads_client())
       
       try:
           service.add_mapping(ws_id, ac.upper(), file or "", test)
           typer.echo(f"✅ Mapped {ac} → {test}")
       except ValueError as e:
           typer.echo(f"❌ {e}")
           raise typer.Exit(1)
   
   @app.command("auto")
   def auto_detect(ws_id: str) -> None:
       """Auto-detect AC→Test mappings (see WS-22)."""
       typer.echo("Auto-detection not implemented yet. See WS-00-032-22.")
       raise typer.Exit(1)
   ```

3. **Write tests**

   ```python
   # tests/unit/test_traceability_cli.py
   import pytest
   from typer.testing import CliRunner
   from sdp.cli.main import app
   
   runner = CliRunner()
   
   class TestTraceCLI:
       def test_check_shows_table(self, mock_beads_with_ws):
           """AC1: check shows mapping table."""
           result = runner.invoke(app, ["trace", "check", "00-032-01"])
           
           assert "Traceability Report" in result.output
           assert "AC" in result.output
           assert "Status" in result.output
       
       def test_check_exits_1_if_incomplete(self, mock_beads_incomplete):
           """AC3: exit code 1 if unmapped ACs."""
           result = runner.invoke(app, ["trace", "check", "00-032-01"])
           
           assert result.exit_code == 1
           assert "INCOMPLETE" in result.output
       
       def test_json_output(self, mock_beads_with_ws):
           """AC4: --json flag returns JSON."""
           result = runner.invoke(app, ["trace", "check", "00-032-01", "--json"])
           
           import json
           data = json.loads(result.output)
           assert "ws_id" in data
           assert "mappings" in data
   ```

### Output Files

- `src/sdp/traceability/service.py`
- `src/sdp/cli/trace.py`
- `tests/unit/test_traceability_cli.py`

### Completion Criteria

```bash
# CLI works
sdp trace --help
sdp trace check --help

# Tests pass
pytest tests/unit/test_traceability_cli.py -v

# Coverage
pytest tests/unit/test_traceability_cli.py --cov=src/sdp/cli/trace --cov-fail-under=80
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC5 — ✅

**Goal Achieved:** ______
