---
ws_id: 00-031-01
feature: F031
status: backlog
size: MEDIUM
project_id: 00
github_issue: null
assignee: null
depends_on: []
---

## WS-00-031-01: Migrate Core Exceptions to SDPError

### 🎯 Goal

**What must WORK after completing this WS:**
- Core custom exception classes inherit from SDPError
- All core exceptions include category, remediation, docs_url
- Error messages are actionable with clear guidance
- Existing exception references updated to new SDPError-based classes

**Acceptance Criteria:**
- [ ] AC1: `WorkstreamParseError` inherits from SDPError with ErrorCategory.VALIDATION
- [ ] AC2: `CircularDependencyError` inherits from SDPError with ErrorCategory.DEPENDENCY
- [ ] AC3: `MissingDependencyError` inherits from SDPError with ErrorCategory.DEPENDENCY
- [ ] AC4: All exceptions include remediation steps (non-empty)
- [ ] AC5: All exception raises in core/ updated to use new format
- [ ] AC6: Coverage ≥80% for exception classes, mypy --strict passes

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Core module has **custom exception classes that don't leverage SDPError framework**.

**Current Exceptions** (in core/):
- `WorkstreamParseError` - plain Exception subclass
- `CircularDependencyError` - plain Exception subclass
- `MissingDependencyError` - plain Exception subclass
- `ContractViolationError` - plain Exception subclass
- `ModelMappingError` - plain Exception subclass
- `HumanEscalationError` - plain Exception subclass

**Impact**:
- Errors lack structured fields (category, remediation, docs_url, context)
- Users get generic error messages without actionable guidance
- Violates Theo Browne's "explicit errors with remediation" principle
- SDPError framework exists but isn't being used

**Solution**: Refactor core exceptions to inherit from SDPError with proper structure.

### Dependencies

None (can proceed independently, SDPError framework already exists)

### Input Files

- `src/sdp/core/workstream.py` (contains WorkstreamParseError)
- `src/sdp/core/feature.py` (contains dependency errors)
- `src/sdp/core/model_mapping.py` (contains ModelMappingError)
- `src/sdp/core/builder_router.py` (may contain errors)
- `src/sdp/errors.py` (SDPError base classes)

### Steps

1. **Review existing SDPError framework**

   ```python
   # src/sdp/errors.py (already exists)
   from dataclasses import dataclass
   from enum import Enum
   from typing import Optional, Dict, Any

   class ErrorCategory(Enum):
       """Error categories for filtering and routing."""
       VALIDATION = "validation"
       DEPENDENCY = "dependency"
       RUNTIME = "runtime"
       CONFIGURATION = "configuration"
       AUTHENTICATION = "authentication"
       AUTHORIZATION = "authorization"
       NOT_FOUND = "not_found"
       BUSINESS_LOGIC = "business_logic"

   @dataclass
   class SDPError(Exception):
       """Base class for all SDP exceptions with structured fields."""

       category: ErrorCategory
       message: str
       remediation: str
       docs_url: Optional[str] = None
       context: Optional[Dict[str, Any]] = None

       def format_terminal(self) -> str:
           """Format error for terminal output."""
           output = [f"❌ {self.category.value.upper()}: {self.message}"]
           output.append(f"\n💡 {self.remediation}")

           if self.docs_url:
               output.append(f"\n📖 Docs: {self.docs_url}")

           if self.context:
               output.append(f"\n📋 Context: {self.context}")

           return "\n".join(output)

       def format_json(self) -> str:
           """Format error as JSON for machine parsing."""
           import json
           return json.dumps({
               "category": self.category.value,
               "message": self.message,
               "remediation": self.remediation,
               "docs_url": self.docs_url,
               "context": self.context,
           }, indent=2)
   ```

2. **Refactor WorkstreamParseError**

   ```python
   # src/sdp/core/workstream.py
   from sdp.errors import SDPError, ErrorCategory
   from typing import Optional
   from pathlib import Path

   class WorkstreamParseError(SDPError):
       """Workstream parsing error with actionable guidance."""

       def __init__(
           self,
           message: str,
           file_path: Optional[Path] = None,
           parse_error: Optional[str] = None,
       ):
           super().__init__(
               category=ErrorCategory.VALIDATION,
               message=message,
               remediation=(
                   "1. Check WS ID format: PP-FFF-SS (e.g., 00-500-01)\n"
                   "   - PP: Project ID (00-99)\n"
                   "   - FFF: Feature ID (001-999)\n"
                   "   - SS: Sequence number (01-99)\n"
                   "2. Ensure file starts with --- frontmatter\n"
                   "3. Validate YAML syntax\n"
                   "4. See docs/workstreams/TEMPLATE.md for template"
               ),
               docs_url="https://sdp.dev/docs/workstreams#format",
               context={
                   "file_path": str(file_path) if file_path else None,
                   "parse_error": parse_error,
               } if file_path or parse_error else None,
           )

   # Update all raise sites
   # BEFORE:
   # raise WorkstreamParseError("Invalid ws_id format")

   # AFTER:
   # raise WorkstreamParseError(
   #     message=f"Invalid ws_id format: {ws_id}",
   #     file_path=file_path,
   # )
   ```

3. **Refactor dependency errors**

   ```python
   # src/sdp/core/feature.py
   from sdp.errors import SDPError, ErrorCategory

   class CircularDependencyError(SDPError):
       """Circular dependency detected in workstream graph."""

       def __init__(
           self,
           ws_id: str,
           cycle: list[str],
       ):
           formatted_cycle = " → ".join(cycle + [cycle[0]])  # Show cycle
           super().__init__(
               category=ErrorCategory.DEPENDENCY,
               message=f"Circular dependency detected: {formatted_cycle}",
               remediation=(
                   f"1. Break the cycle by removing one dependency:\n"
                   f"   - {ws_id} depends on: {' → '.join(cycle)}\n"
                   f"2. Reorder workstreams to avoid circular reference\n"
                   f"3. Or split into smaller independent features\n"
                   f"4. See docs/dependency-management.md for strategies"
               ),
               docs_url="https://sdp.dev/docs/dependencies#circular",
               context={"ws_id": ws_id, "cycle": cycle},
           )

   class MissingDependencyError(SDPError):
       """Required workstream dependency not found."""

       def __init__(
           self,
           ws_id: str,
           missing_dep: str,
           available_workstreams: list[str],
       ):
           super().__init__(
               category=ErrorCategory.DEPENDENCY,
               message=f"Workstream {ws_id} depends on {missing_dep} which doesn't exist",
               remediation=(
                   f"1. Create missing workstream first: {missing_dep}\n"
                   f"2. Or remove dependency if not actually needed\n"
                   f"3. Available workstreams: {', '.join(available_workstreams[:5])}\n"
                   f"4. See docs/workflows/dependency-management.md"
               ),
               docs_url="https://sdp.dev/docs/workflows#dependencies",
               context={
                   "ws_id": ws_id,
                   "missing_dep": missing_dep,
                   "available_ws": available_workstreams,
               },
           )
   ```

4. **Update all exception raise sites**

   Search and replace in `src/sdp/core/`:

   ```python
   # BEFORE (pattern to find):
   raise WorkstreamParseError(f"Invalid ws_id: {ws_id}")

   # AFTER (refactored):
   raise WorkstreamParseError(
       message=f"Invalid ws_id: {ws_id}",
       file_path=self.file_path,
   )

   # BEFORE:
   raise CircularDependencyError(f"Cycle: {cycle}")

   # AFTER:
   raise CircularDependencyError(
       ws_id=self.ws_id,
       cycle=cycle,
   )
   ```

5. **Create tests for exception formatting**

   ```python
   # tests/unit/core/test_exceptions.py
   import pytest
   from sdp.core.workstream import WorkstreamParseError
   from sdp.core.feature import CircularDependencyError, MissingDependencyError

   class TestWorkstreamParseError:
       """Test WorkstreamParseError SDPError compliance."""

       def test_includes_error_category(self):
           """Verify error includes VALIDATION category."""
           error = WorkstreamParseError("Invalid ws_id")

           assert error.category.value == "validation"
           assert "VALIDATION" in error.format_terminal()

       def test_includes_remediation_steps(self):
           """Verify error includes actionable remediation."""
           error = WorkstreamParseError("Invalid ws_id")

           formatted = error.format_terminal()
           assert "💡" in formatted  # Remediation section
           assert "1. Check WS ID format" in formatted
           assert "PP-FFF-SS" in formatted

       def test_includes_docs_url(self):
           """Verify error includes documentation link."""
           error = WorkstreamParseError("Invalid ws_id")

           formatted = error.format_terminal()
           assert "📖 Docs:" in formatted
           assert "https://sdp.dev/docs/workstreams#format" in formatted

       def test_includes_context_when_provided(self):
           """Verify error includes context when available."""
           from pathlib import Path

           error = WorkstreamParseError(
               message="Invalid ws_id",
               file_path=Path("/path/to/ws.md"),
           )

           formatted = error.format_terminal()
           assert "📋 Context:" in formatted
           assert "/path/to/ws.md" in formatted

       def test_formats_json_for_machine_parsing(self):
           """Verify error formats as JSON for CI/CD systems."""
           import json

           error = WorkstreamParseError("Invalid ws_id")

           json_str = error.format_json()
           parsed = json.loads(json_str)

           assert parsed["category"] == "validation"
           assert parsed["message"] == "Invalid ws_id"
           assert "remediation" in parsed
           assert "docs_url" in parsed

   class TestCircularDependencyError:
       """Test CircularDependencyError SDPError compliance."""

       def test_shows_cycle_in_error_message(self):
           """Verify error message shows the dependency cycle."""
           error = CircularDependencyError(
               ws_id="00-001-01",
               cycle=["00-001-02", "00-001-03", "00-001-01"],
           )

           formatted = error.format_terminal()
           assert "00-001-02 → 00-001-03 → 00-001-01" in formatted

       def test_suggests_breaking_cycle(self):
           """Verify remediation suggests breaking the cycle."""
           error = CircularDependencyError(
               ws_id="00-001-01",
               cycle=["00-001-02", "00-001-03", "00-001-01"],
           )

           formatted = error.format_terminal()
           assert "1. Break the cycle" in formatted
           assert "remove one dependency" in formatted

   class TestMissingDependencyError:
       """Test MissingDependencyError SDPError compliance."""

       def test_shows_available_alternatives(self):
           """Verify error shows available workstreams."""
           error = MissingDependencyError(
               ws_id="00-001-01",
               missing_dep="00-001-02",
               available_workstreams=["00-001-03", "00-001-04"],
           )

           formatted = error.format_terminal()
           assert "00-001-03" in formatted
           assert "00-001-04" in formatted
   ```

6. **Update imports in files that use these exceptions**

   Search for `from sdp.core.workstream import WorkstreamParseError` etc. and verify compatibility.

   The new exception classes inherit from SDPError which inherits from Exception, so all `except WorkstreamParseError` blocks will still work.

### Code

```python
# src/sdp/core/feature.py (example of updated usage)
from sdp.errors import SDPError, ErrorCategory

class FeatureDecomposer:
    """Decomposes features into workstreams."""

    def validate_dependencies(self, workstreams: list[Workstream]) -> None:
        """Validate workstream dependencies for cycles."""
        # Build dependency graph
        graph = self._build_graph(workstreams)

        # Detect cycles
        for ws_id, deps in graph.items():
            if self._has_cycle(ws_id, deps, graph):
                # Use new SDPError-based exception
                raise CircularDependencyError(
                    ws_id=ws_id,
                    cycle=self._extract_cycle(ws_id, graph),
                )

    def _build_graph(self, workstreams: list[Workstream]) -> dict:
        """Build dependency graph from workstreams."""
        graph = {}
        for ws in workstreams:
            graph[ws.ws_id] = ws.dependencies or []
        return graph
```

### Expected Outcome

**After completion:**
- Core custom exceptions now leverage SDPError framework
- All errors include category, remediation, docs_url, context
- Error messages are actionable (Theo Browne principle satisfied)
- Foundation for consistent error handling across codebase
- Terminal output shows structured errors with emoji indicators
- JSON output available for CI/CD parsing

**Scope Estimate**
- Files: ~12
- Lines: ~600 (MEDIUM)
- Tokens: ~3000

### Completion Criteria

```bash
# Verify exceptions inherit from SDPError
python -c "from sdp.core.workstream import WorkstreamParseError; from sdp.errors import SDPError; assert issubclass(WorkstreamParseError, SDPError)"

# Run exception tests
pytest tests/unit/core/test_exceptions.py -v

# Verify coverage
pytest --cov=src/sdp/core --cov-report=term-missing
# Should show ≥80% for exception classes

# Verify type checking
mypy src/sdp/core/ --strict

# Test error formatting
python -c "
from sdp.core.workstream import WorkstreamParseError
error = WorkstreamParseError('Invalid ws_id')
print(error.format_terminal())
"
# Should show formatted error with remediation
```

### Constraints

- DO NOT break existing exception catching (inheritance preserves this)
- DO NOT change exception names (keep class names for compatibility)
- DO NOT remove fields from context (add more, don't remove)
- DO NOT create new exception types (only refactor existing)

---

## Execution Report

**Executed by:** ______
**Date:** ______
**Duration:** ______ minutes

### Goal Status
- [ ] AC1-AC6 — ✅

**Goal Achieved:** ______

### Files Changed
| File | Action | LOC |
|------|--------|-----|
| src/sdp/core/workstream.py | Modify | +100 |
| src/sdp/core/feature.py | Modify | +150 |
| src/sdp/core/model_mapping.py | Modify | +50 |
| tests/unit/core/test_exceptions.py | Create | 250 |
| docs/reference/error-handling.md | Modify | +50 |

### Statistics
- **Files Changed:** 5
- **Lines Added:** ~600
- **Lines Removed:** ~0
- **Test Coverage:** ≥80% for core exceptions
- **Tests Passed:** ______
- **Tests Failed:** ______

### Deviations from Plan
- ______

### Commit
______
