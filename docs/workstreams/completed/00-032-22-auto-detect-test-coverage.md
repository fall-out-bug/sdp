---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-21
feature: F032
github_issue: null
project_id: 0
size: MEDIUM
status: completed
traceability:
- ac_description: Parse test docstrings для AC references (`'''Tests AC1'''`)
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: Parse test names (`test_ac1_*`, `test_acceptance_criterion_1`)
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: Heuristic matching по keywords
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: '`sdp trace auto <ws_id>` запускает auto-detection'
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: Confidence score 0.0-1.0 для каждого mapping
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-22
---

## 00-032-22: Auto-Detect Test Coverage

### 🎯 Goal

**What must WORK after completing this WS:**
- Автоматическое определение какие тесты покрывают какие AC
- Parse docstrings и test names
- Confidence score для каждого mapping

**Acceptance Criteria:**
- [ ] AC1: Parse test docstrings для AC references (`'''Tests AC1'''`)
- [ ] AC2: Parse test names (`test_ac1_*`, `test_acceptance_criterion_1`)
- [ ] AC3: Heuristic matching по keywords
- [ ] AC4: Confidence score 0.0-1.0 для каждого mapping
- [ ] AC5: `sdp trace auto <ws_id>` запускает auto-detection

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Manual mapping трудоёмкий. Нужна автоматизация.

**Solution**: AST parsing + heuristics для auto-detection.

### Dependencies

- **00-032-21**: Traceability CLI

### Steps

1. **Create ACDetector**

   ```python
   # src/sdp/traceability/detector.py
   import ast
   import re
   from pathlib import Path
   from dataclasses import dataclass
   
   @dataclass
   class DetectedMapping:
       ac_id: str
       test_file: str
       test_name: str
       confidence: float
       source: str  # "docstring", "name", "keyword"
   
   class ACDetector:
       """Auto-detect AC coverage from test files."""
       
       def detect_all(
           self,
           test_dir: Path,
           ac_descriptions: dict[str, str]
       ) -> list[DetectedMapping]:
           """Detect mappings from all test files."""
           mappings = []
           
           for test_file in test_dir.rglob("test_*.py"):
               file_mappings = self.detect_from_file(test_file, ac_descriptions)
               mappings.extend(file_mappings)
           
           return mappings
       
       def detect_from_file(
           self,
           test_file: Path,
           ac_descriptions: dict[str, str]
       ) -> list[DetectedMapping]:
           """Detect mappings from single test file."""
           try:
               tree = ast.parse(test_file.read_text())
           except SyntaxError:
               return []
           
           mappings = []
           
           for node in ast.walk(tree):
               if isinstance(node, ast.FunctionDef) and node.name.startswith("test_"):
                   # Check docstring
                   doc_mappings = self._detect_from_docstring(
                       node, str(test_file)
                   )
                   mappings.extend(doc_mappings)
                   
                   # Check name
                   name_mappings = self._detect_from_name(
                       node.name, str(test_file)
                   )
                   mappings.extend(name_mappings)
                   
                   # Check keywords against AC descriptions
                   keyword_mappings = self._detect_from_keywords(
                       node, str(test_file), ac_descriptions
                   )
                   mappings.extend(keyword_mappings)
           
           return mappings
       
       def _detect_from_docstring(
           self,
           func: ast.FunctionDef,
           test_file: str
       ) -> list[DetectedMapping]:
           """Extract AC references from docstring.
           
           Patterns:
               '''Tests AC1: User can login'''
               '''Covers: AC1, AC2'''
               '''AC1'''
           """
           docstring = ast.get_docstring(func)
           if not docstring:
               return []
           
           # Pattern: AC1, AC2, etc.
           ac_refs = re.findall(r'\bAC(\d+)\b', docstring, re.IGNORECASE)
           
           return [
               DetectedMapping(
                   ac_id=f"AC{ref}",
                   test_file=test_file,
                   test_name=func.name,
                   confidence=0.95,  # High confidence from docstring
                   source="docstring"
               )
               for ref in ac_refs
           ]
       
       def _detect_from_name(
           self,
           test_name: str,
           test_file: str
       ) -> list[DetectedMapping]:
           """Extract AC from test name.
           
           Patterns:
               test_ac1_user_login -> AC1
               test_acceptance_criterion_2 -> AC2
               test_ac_1_something -> AC1
           """
           patterns = [
               (r'test_ac(\d+)', 0.90),
               (r'test_ac_(\d+)', 0.85),
               (r'test_acceptance_criterion_(\d+)', 0.90),
           ]
           
           mappings = []
           for pattern, confidence in patterns:
               match = re.search(pattern, test_name, re.IGNORECASE)
               if match:
                   mappings.append(DetectedMapping(
                       ac_id=f"AC{match.group(1)}",
                       test_file=test_file,
                       test_name=test_name,
                       confidence=confidence,
                       source="name"
                   ))
           
           return mappings
       
       def _detect_from_keywords(
           self,
           func: ast.FunctionDef,
           test_file: str,
           ac_descriptions: dict[str, str]
       ) -> list[DetectedMapping]:
           """Match test to AC by keyword similarity.
           
           Lower confidence - heuristic matching.
           """
           # Get test name words
           test_words = set(re.findall(r'\w+', func.name.lower()))
           
           # Get docstring words
           docstring = ast.get_docstring(func) or ""
           test_words.update(re.findall(r'\w+', docstring.lower()))
           
           mappings = []
           
           for ac_id, description in ac_descriptions.items():
               # Already detected by other methods?
               # Skip keyword matching for those
               
               desc_words = set(re.findall(r'\w+', description.lower()))
               
               # Remove common words
               common = {"the", "a", "an", "is", "are", "can", "should", "must", "test"}
               test_words -= common
               desc_words -= common
               
               # Calculate overlap
               if not desc_words:
                   continue
               
               overlap = len(test_words & desc_words) / len(desc_words)
               
               if overlap >= 0.5:  # At least 50% word match
                   mappings.append(DetectedMapping(
                       ac_id=ac_id,
                       test_file=test_file,
                       test_name=func.name,
                       confidence=min(0.7, overlap),  # Cap at 0.7 for keywords
                       source="keyword"
                   ))
           
           return mappings
   ```

2. **Update CLI**

   ```python
   # src/sdp/cli/trace.py (update auto command)
   @app.command("auto")
   def auto_detect(
       ws_id: str,
       test_dir: str = typer.Option("tests/", help="Test directory"),
       apply: bool = typer.Option(False, help="Apply detected mappings")
   ) -> None:
       """Auto-detect AC→Test mappings."""
       from sdp.traceability.detector import ACDetector
       
       service = TraceabilityService(create_beads_client())
       detector = ACDetector()
       
       # Get WS and extract ACs
       report = service.check_traceability(ws_id)
       ac_descriptions = {
           m.ac_id: m.ac_description
           for m in report.mappings
       }
       
       # Detect mappings
       detected = detector.detect_all(Path(test_dir), ac_descriptions)
       
       if not detected:
           typer.echo("No mappings detected")
           return
       
       typer.echo(f"\nDetected {len(detected)} potential mappings:\n")
       
       for d in sorted(detected, key=lambda x: (-x.confidence, x.ac_id)):
           conf_bar = "█" * int(d.confidence * 10) + "░" * (10 - int(d.confidence * 10))
           typer.echo(f"  {d.ac_id} → {d.test_name}")
           typer.echo(f"    File: {d.test_file}")
           typer.echo(f"    Confidence: [{conf_bar}] {d.confidence:.0%} ({d.source})")
           typer.echo("")
       
       if apply:
           for d in detected:
               if d.confidence >= 0.8:  # Only apply high-confidence
                   service.add_mapping(ws_id, d.ac_id, d.test_file, d.test_name)
                   typer.echo(f"✅ Applied: {d.ac_id} → {d.test_name}")
       else:
           typer.echo("Run with --apply to save high-confidence mappings")
   ```

### Output Files

- `src/sdp/traceability/detector.py`
- `src/sdp/cli/trace.py` (updated)
- `tests/unit/test_ac_detector.py`

### Completion Criteria

```bash
# Detector works
python -c "from sdp.traceability.detector import ACDetector"

# CLI works
sdp trace auto 00-032-01 --help

# Tests pass
pytest tests/unit/test_ac_detector.py -v
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC5 — ✅

**Goal Achieved:** ______
