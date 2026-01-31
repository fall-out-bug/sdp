---
assignee: null
completed: '2026-01-30'
depends_on: []
feature: F032
github_issue: null
project_id: 0
size: SMALL
status: completed
traceability:
- ac_description: '`ACTestMapping` dataclass создан'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: '`TraceabilityReport` dataclass создан'
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: JSON schema в `docs/schema/traceability.json`
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: ''
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: Unit tests для моделей
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-20
---

## 00-032-20: AC Test Mapping Model

### 🎯 Goal

**What must WORK after completing this WS:**
- Модель данных для AC→Test mapping
- Format хранения в WS metadata
- JSON schema для валидации

**Acceptance Criteria:**
- [ ] AC1: `ACTestMapping` dataclass создан
- [ ] AC2: `TraceabilityReport` dataclass создан
- [ ] AC3: JSON schema в `docs/schema/traceability.json`
- [ ] AC4: Unit tests для моделей

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Review не проверяет соответствие AC и тестов.

**Solution**: Модель данных для трассируемости.

### Dependencies

None (independent)

### Steps

1. **Create models**

   ```python
   # src/sdp/traceability/models.py
   from dataclasses import dataclass, field
   from enum import Enum
   from typing import Literal
   
   class MappingStatus(Enum):
       MAPPED = "mapped"      # AC has test
       MISSING = "missing"    # AC has no test
       FAILED = "failed"      # Test exists but fails
   
   @dataclass
   class ACTestMapping:
       """Maps Acceptance Criterion to Test."""
       
       ac_id: str                    # "AC1", "AC2", etc.
       ac_description: str           # "User can login"
       test_file: str | None         # "tests/unit/test_auth.py"
       test_name: str | None         # "test_user_login"
       status: MappingStatus
       confidence: float = 1.0       # 0.0-1.0 for auto-detected
       
       def to_dict(self) -> dict:
           return {
               "ac_id": self.ac_id,
               "ac_description": self.ac_description,
               "test_file": self.test_file,
               "test_name": self.test_name,
               "status": self.status.value,
               "confidence": self.confidence,
           }
       
       @classmethod
       def from_dict(cls, data: dict) -> "ACTestMapping":
           return cls(
               ac_id=data["ac_id"],
               ac_description=data["ac_description"],
               test_file=data.get("test_file"),
               test_name=data.get("test_name"),
               status=MappingStatus(data["status"]),
               confidence=data.get("confidence", 1.0),
           )
   
   @dataclass
   class TraceabilityReport:
       """Report of AC→Test traceability for a workstream."""
       
       ws_id: str
       mappings: list[ACTestMapping] = field(default_factory=list)
       
       @property
       def total_acs(self) -> int:
           return len(self.mappings)
       
       @property
       def mapped_acs(self) -> int:
           return sum(1 for m in self.mappings if m.status == MappingStatus.MAPPED)
       
       @property
       def missing_acs(self) -> int:
           return sum(1 for m in self.mappings if m.status == MappingStatus.MISSING)
       
       @property
       def failed_acs(self) -> int:
           return sum(1 for m in self.mappings if m.status == MappingStatus.FAILED)
       
       @property
       def coverage_pct(self) -> float:
           if self.total_acs == 0:
               return 100.0
           return (self.mapped_acs / self.total_acs) * 100
       
       @property
       def is_complete(self) -> bool:
           return self.missing_acs == 0
       
       def to_dict(self) -> dict:
           return {
               "ws_id": self.ws_id,
               "total_acs": self.total_acs,
               "mapped_acs": self.mapped_acs,
               "missing_acs": self.missing_acs,
               "coverage_pct": self.coverage_pct,
               "mappings": [m.to_dict() for m in self.mappings],
           }
       
       def to_markdown_table(self) -> str:
           """Generate markdown table for report."""
           lines = [
               "| AC | Description | Test | Status |",
               "|----|-------------|------|--------|",
           ]
           
           for m in self.mappings:
               test = f"`{m.test_name}`" if m.test_name else "-"
               status = "✅" if m.status == MappingStatus.MAPPED else "❌"
               lines.append(f"| {m.ac_id} | {m.ac_description[:30]} | {test} | {status} |")
           
           return "\n".join(lines)
   ```

2. **Create JSON schema**

   ```json
   // docs/schema/traceability.json
   {
     "$schema": "http://json-schema.org/draft-07/schema#",
     "title": "TraceabilityReport",
     "type": "object",
     "required": ["ws_id", "mappings"],
     "properties": {
       "ws_id": {
         "type": "string",
         "pattern": "^\\d{2}-\\d{3}-\\d{2}$"
       },
       "total_acs": {
         "type": "integer",
         "minimum": 0
       },
       "mapped_acs": {
         "type": "integer",
         "minimum": 0
       },
       "missing_acs": {
         "type": "integer",
         "minimum": 0
       },
       "coverage_pct": {
         "type": "number",
         "minimum": 0,
         "maximum": 100
       },
       "mappings": {
         "type": "array",
         "items": {
           "$ref": "#/definitions/ACTestMapping"
         }
       }
     },
     "definitions": {
       "ACTestMapping": {
         "type": "object",
         "required": ["ac_id", "ac_description", "status"],
         "properties": {
           "ac_id": {
             "type": "string",
             "pattern": "^AC\\d+$"
           },
           "ac_description": {
             "type": "string"
           },
           "test_file": {
             "type": ["string", "null"]
           },
           "test_name": {
             "type": ["string", "null"]
           },
           "status": {
             "type": "string",
             "enum": ["mapped", "missing", "failed"]
           },
           "confidence": {
             "type": "number",
             "minimum": 0,
             "maximum": 1
           }
         }
       }
     }
   }
   ```

3. **Write tests**

   ```python
   # tests/unit/test_traceability_models.py
   import pytest
   from sdp.traceability.models import (
       ACTestMapping,
       TraceabilityReport,
       MappingStatus,
   )
   
   class TestACTestMapping:
       def test_to_dict(self):
           mapping = ACTestMapping(
               ac_id="AC1",
               ac_description="User can login",
               test_file="tests/test_auth.py",
               test_name="test_login",
               status=MappingStatus.MAPPED,
           )
           
           d = mapping.to_dict()
           
           assert d["ac_id"] == "AC1"
           assert d["status"] == "mapped"
       
       def test_from_dict(self):
           data = {
               "ac_id": "AC1",
               "ac_description": "User can login",
               "test_file": "tests/test_auth.py",
               "test_name": "test_login",
               "status": "mapped",
           }
           
           mapping = ACTestMapping.from_dict(data)
           
           assert mapping.ac_id == "AC1"
           assert mapping.status == MappingStatus.MAPPED
   
   class TestTraceabilityReport:
       def test_coverage_pct(self):
           report = TraceabilityReport(
               ws_id="00-032-01",
               mappings=[
                   ACTestMapping("AC1", "Test", None, None, MappingStatus.MAPPED),
                   ACTestMapping("AC2", "Test", None, None, MappingStatus.MAPPED),
                   ACTestMapping("AC3", "Test", None, None, MappingStatus.MISSING),
               ]
           )
           
           assert report.coverage_pct == pytest.approx(66.67, rel=0.01)
       
       def test_is_complete(self):
           complete = TraceabilityReport(
               ws_id="00-032-01",
               mappings=[
                   ACTestMapping("AC1", "Test", "test.py", "test_1", MappingStatus.MAPPED),
               ]
           )
           
           incomplete = TraceabilityReport(
               ws_id="00-032-01",
               mappings=[
                   ACTestMapping("AC1", "Test", None, None, MappingStatus.MISSING),
               ]
           )
           
           assert complete.is_complete is True
           assert incomplete.is_complete is False
   ```

### Output Files

- `src/sdp/traceability/__init__.py`
- `src/sdp/traceability/models.py`
- `docs/schema/traceability.json`
- `tests/unit/test_traceability_models.py`

### Completion Criteria

```bash
# Module imports
python -c "from sdp.traceability.models import ACTestMapping, TraceabilityReport"

# Tests pass
pytest tests/unit/test_traceability_models.py -v

# Schema valid
python -c "import json; json.load(open('docs/schema/traceability.json'))"
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC4 — ✅

**Goal Achieved:** ______
