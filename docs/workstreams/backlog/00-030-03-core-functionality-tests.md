---
ws_id: 00-030-03
feature: F030
status: backlog
size: MEDIUM
project_id: 00
github_issue: null
assignee: null
depends_on: []
---

## WS-00-030-03: Core Functionality Tests

### 🎯 Goal

**What must WORK after completing this WS:**
- Core business logic modules have comprehensive test coverage
- Workstream parser tested with various formats
- Builder router tested for model selection logic
- Model mapping registry tested for correctness

**Acceptance Criteria:**
- [ ] AC1: `tests/unit/core/test_workstream.py` with ≥80% coverage
- [ ] AC2: `tests/unit/core/test_builder_router.py` with ≥80% coverage
- [ ] AC3: `tests/unit/core/test_model_mapping.py` with ≥80% coverage
- [ ] AC4: `tests/unit/core/test_project_map_parser.py` with ≥80% coverage
- [ ] AC5: Tests cover edge cases (invalid WS IDs, circular dependencies)
- [ ] AC6: mypy --strict passes on all test files

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Critical Gap**: Core functionality has **2,288 LOC across 12 modules with minimal tests**.

**High-Risk Modules**:
1. **`src/sdp/core/workstream.py`** (~340 LOC)
   - Workstream domain entity
   - Frontmatter parsing (YAML)
   - Validation logic

2. **`src/sdp/core/builder_router.py`** (~278 LOC)
   - Model selection logic
   - Capability matching
   - Fallback mechanisms

3. **`src/sdp/core/model_mapping.py`** (~221 LOC)
   - Registry pattern for model mappings
   - Dynamic model loading
   - Fallback to default models

4. **`src/sdp/core/project_map_parser.py`** (~301 LOC)
   - Project map parsing logic
   - Dependency graph construction

**Risk**: Core business logic is the foundation. Untested core = unreliable system.

### Dependencies

None (can proceed independently)

### Input Files

- `src/sdp/core/workstream.py`
- `src/sdp/core/builder_router.py`
- `src/sdp/core/model_mapping.py`
- `src/sdp/core/project_map_parser.py`
- `src/sdp/core/feature.py` (if exists)

### Steps

1. **Create test directory structure**

   ```
   tests/
   └── unit/
       └── core/
           ├── __init__.py
           ├── test_workstream.py
           ├── test_builder_router.py
           ├── test_model_mapping.py
           ├── test_project_map_parser.py
           └── fixtures/
               ├── valid_workstreams.md      # Sample WS files
               ├── invalid_workstreams.md    # Edge cases
               └── project_maps.yaml        # Sample project maps
   ```

2. **Test workstream parsing**

   ```python
   # tests/unit/core/test_workstream.py
   import pytest
   from pathlib import Path
   from sdp.core.workstream import Workstream, WorkstreamParser
   from sdp.core.exceptions import WorkstreamParseError, ValidationError

   @pytest.fixture
   def valid_workstream_file(tmp_path):
       """Create a valid workstream file."""
       ws_file = tmp_path / "00-001-01.md"
       ws_file.write_text('''---
   ws_id: 00-001-01
   feature: F001
   status: backlog
   size: SMALL
   project_id: 00
   ---

   ## WS-00-001-01: Test Workstream

   ### 🎯 Goal

   Test goal

   ### Context

   Test context
   ''')
       return ws_file

   class TestWorkstreamParser:
       """Test workstream file parsing."""

       def test_parse_valid_workstream(self, valid_workstream_file):
           """Verify parser extracts all fields correctly."""
           parser = WorkstreamParser()
           ws = parser.parse(valid_workstream_file)

           assert ws.ws_id == "00-001-01"
           assert ws.feature == "F001"
           assert ws.status == "backlog"
           assert ws.size == "SMALL"
           assert ws.project_id == "00"
           assert "Test goal" in ws.content

       def test_parse_raises_error_for_invalid_ws_id(self, tmp_path):
           """Verify parser rejects invalid WS ID format."""
           invalid_file = tmp_path / "WS-001-01.md"  # Old format
           invalid_file.write_text("ws_id: WS-001-01")

           parser = WorkstreamParser()

           with pytest.raises(WorkstreamParseError, match="Invalid ws_id format"):
               parser.parse(invalid_file)

       def test_parse_raises_error_for_missing_required_fields(self, tmp_path):
           """Verify parser requires all required fields."""
           incomplete_file = tmp_path / "00-001-01.md"
           incomplete_file.write_text("---\nws_id: 00-001-01\n---")  # Missing feature

           parser = WorkstreamParser()

           with pytest.raises(WorkstreamParseError, match="Missing required field.*feature"):
               parser.parse(incomplete_file)

       def test_parse_handles_extra_whitespace(self, valid_workstream_file):
           """Verify parser handles extra whitespace gracefully."""
           # Add extra whitespace
           content = valid_workstream_file.read_text()
           valid_workstream_file.write_text(content.replace("\n\n", "\n\n\n"))

           parser = WorkstreamParser()
           ws = parser.parse(valid_workstream_file)

           assert ws.ws_id == "00-001-01"  # Should still parse correctly

   class TestWorkstreamValidation:
       """Test workstream validation logic."""

       def test_validate_accepts_complete_workstream(self):
           """Verify validation passes for complete workstream."""
           ws = Workstream(
               ws_id="00-001-01",
               feature="F001",
               status="backlog",
               size="SMALL",
               project_id="00",
               goal="Test goal",
               acceptance_criteria=["AC1: Test"],
               context="Test context",
               steps=["Step 1"],
           )

           # Should not raise
           ws.validate()

       def test_validate_rejects_empty_goal(self):
           """Verify validation rejects empty goal."""
           ws = Workstream(
               ws_id="00-001-01",
               feature="F001",
               status="backlog",
               size="SMALL",
               project_id="00",
               goal="",  # Empty goal
               acceptance_criteria=[],
               context="Test",
               steps=[],
           )

           with pytest.raises(ValidationError, match="Goal cannot be empty"):
               ws.validate()

       def test_validate_rejects_size_mismatch(self, tmp_path):
           """Verify validation detects size mismatch (declared SMALL but 1000 LOC)."""
           # Create file with 1000 lines
           large_file = tmp_path / "large_module.py"
           large_file.write_text("\n".join(["# line"] * 1000))

           ws = Workstream(
               ws_id="00-001-01",
               feature="F001",
               status="backlog",
               size="SMALL",  # Claims SMALL but file is LARGE
               project_id="00",
               goal="Test",
               acceptance_criteria=[],
               context="Test",
               steps=[],
               files=[large_file],
           )

           with pytest.raises(ValidationError, match="Size mismatch"):
               ws.validate()
   ```

3. **Test builder router**

   ```python
   # tests/unit/core/test_builder_router.py
   import pytest
   from unittest.mock import Mock, patch
   from sdp.core.builder_router import BuilderRouter, ModelCapabilities
   from sdp.core.exceptions import ModelSelectionError

   @pytest.fixture
   def router():
       """Create builder router with default models."""
       return BuilderRouter()

   class TestModelCapabilities:
       """Test model capability detection."""

       def test_claude_sonnet_supports_system_prompt(self, router):
           """Verify Claude Sonnet is detected as supporting system prompts."""
           capabilities = router.get_capabilities("claude-sonnet-4-5")

           assert capabilities.supports_system_prompt
           assert capabilities.supports_tools
           assert capabilities.supports_extended_context

       def test_gpt4_supports_function_calling(self, router):
           """Verify GPT-4 is detected as supporting function calling."""
           capabilities = router.get_capabilities("gpt-4")

           assert capabilities.supports_function_calling
           assert capabilities.supports_tools

   class TestModelSelection:
       """Test model selection logic."""

       def test_select_model_for_capability_requirement(self, router):
           """Verify router selects correct model for capability."""
           requirement = ModelCapabilities(
               supports_system_prompt=True,
               supports_tools=True,
           )

           selected = router.select_model(requirement)

           # Should prefer Claude Sonnet (best match)
           assert "claude" in selected.lower()

       def test_select_model_raises_error_when_no_match(self, router):
           """Verify error raised when no model supports requirement."""
           impossible_requirement = ModelCapabilities(
               supports_brain_computer_interface=True,  # Nobody supports this
           )

           with pytest.raises(ModelSelectionError, match="No model supports required capability"):
               router.select_model(impossible_requirement)

       def test_select_uses_fallback_when_primary_unavailable(self, router):
           """Verify fallback model used when primary unavailable."""
           # Mock primary model as unavailable
           router.models["claude-sonnet-4-5"]["available"] = False

           requirement = ModelCapabilities(supports_tools=True)

           selected = router.select_model(requirement)

           # Should fall back to gpt-4
           assert "gpt" in selected.lower()

   class TestBuilderRouter:
       """Test builder router orchestration."""

       def test_route_task_to_best_available_model(self, router):
           """Verify router routes task to best available model."""
           task = {
               "type": "build",
               "requires_tools": True,
               "requires_system_prompt": False,
           }

           model = router.route(task)

           assert model is not None
           assert router.is_model_available(model)

       def test_route_returns_none_if_no_model_available(self, router):
           """Verify router returns None when all models unavailable."""
           # Mock all models as unavailable
           for model_name in router.models:
               router.models[model_name]["available"] = False

           task = {"type": "build", "requires_tools": True}

           model = router.route(task)

           assert model is None
   ```

4. **Test model mapping registry**

   ```python
   # tests/unit/core/test_model_mapping.py
   import pytest
   from sdp.core.model_mapping import ModelMappingRegistry, ModelMapping
   from sdp.core.exceptions import ModelMappingError

   @pytest.fixture
   def registry():
       """Create model mapping registry."""
       return ModelMappingRegistry()

   class TestModelMappingRegistry:
       """Test model mapping registry pattern."""

       def test_register_mapping(self, registry):
           """Verify registering a new model mapping."""
           mapping = ModelMapping(
               name="custom-model",
               capabilities=["tools", "system_prompt"],
               endpoint="custom-api.com",
           )

           registry.register(mapping)

           assert "custom-model" in registry.mappings

       def test_get_mapping_returns_correct_mapping(self, registry):
           """Verify get_mapping returns registered mapping."""
           mapping = ModelMapping(
               name="claude-sonnet-4-5",
               capabilities=["tools"],
               endpoint="anthropic.com",
           )
           registry.register(mapping)

           retrieved = registry.get_mapping("claude-sonnet-4-5")

           assert retrieved.name == "claude-sonnet-4-5"
           assert retrieved.endpoint == "anthropic.com"

       def test_get_mapping_raises_error_for_unknown_model(self, registry):
           """Verify get_mapping raises error for unknown model."""
           with pytest.raises(ModelMappingError, match="Model not found"):
               registry.get_mapping("unknown-model")

       def test_list_mappings_filters_by_capability(self, registry):
           """Verify list_mappings filters models by capability."""
           tool_models = [
               ModelMapping(name="model-a", capabilities=["tools"], endpoint="a.com"),
               ModelMapping(name="model-b", capabilities=["tools", "system_prompt"], endpoint="b.com"),
               ModelMapping(name="model-c", capabilities=["system_prompt"], endpoint="c.com"),
           ]
           for mapping in tool_models:
               registry.register(mapping)

           tools_only = registry.list_mappings(capabilities=["tools"])

           assert len(tools_only) == 2
           assert all("tools" in m.capabilities for m in tools_only)
   ```

### Code

```python
# src/sdp/core/workstream.py (excerpt showing what to test)
from dataclasses import dataclass
from typing import Optional, List

@dataclass
class Workstream:
    """Workstream domain entity."""

    ws_id: str
    feature: str
    status: str  # backlog|in_progress|completed
    size: str  # SMALL|MEDIUM|LARGE
    project_id: str
    goal: str
    acceptance_criteria: List[str]
    context: str
    steps: List[str]
    files: Optional[List[Path]] = None
    dependencies: Optional[List[str]] = None
    completion_criteria: Optional[List[str]] = None

    def validate(self) -> None:
        """Validate workstream completeness and consistency.

        Raises:
            ValidationError: If workstream is invalid
        """
        errors = []

        # Check required fields
        if not self.goal or self.goal.strip() == "":
            errors.append("Goal cannot be empty")

        if not self.acceptance_criteria:
            errors.append("At least one acceptance criterion required")

        if not self.steps:
            errors.append("At least one step required")

        # Check size consistency
        if self.files:
            total_loc = sum(len(f.read_text().split('\n')) for f in self.files)
            if self.size == "SMALL" and total_loc > 500:
                errors.append(f"Size mismatch: Claims SMALL but {total_loc} LOC")
            elif self.size == "MEDIUM" and (total_loc <= 500 or total_loc > 1500):
                errors.append(f"Size mismatch: Claims MEDIUM but {total_loc} LOC")
            elif self.size == "LARGE" and total_loc <= 1500:
                errors.append(f"Size mismatch: Claims LARGE but {total_loc} LOC")

        if errors:
            raise ValidationError("\n".join(errors))
```

### Expected Outcome

**After completion:**
- Core business logic (2,288 LOC) now has test coverage
- Workstream parser tested with various formats
- Model selection logic verified with edge cases
- Foundation for refactoring core domain with confidence

**Scope Estimate**
- Files: ~9
- Lines: ~900 (MEDIUM)
- Tokens: ~4500

### Completion Criteria

```bash
# Run all core tests
pytest tests/unit/core/ -v

# Verify coverage
pytest --cov=src/sdp/core --cov-report=term-missing
# Should show ≥80% for workstream.py, builder_router.py, model_mapping.py

# Verify type checking
mypy src/sdp/core/ --strict

# Run specific test class
pytest tests/unit/core/test_workstream.py::TestWorkstreamParser -v
```

### Constraints

- DO NOT test external model APIs (only our routing logic)
- DO NOT make real model API calls (use mocks)
- DO NOT hardcode model lists (use fixtures)
- DO NOT add new dependencies

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
| tests/unit/core/__init__.py | Create | 5 |
| tests/unit/core/test_workstream.py | Create | 250 |
| tests/unit/core/test_builder_router.py | Create | 200 |
| tests/unit/core/test_model_mapping.py | Create | 180 |
| tests/unit/core/test_project_map_parser.py | Create | 180 |
| tests/unit/core/fixtures/valid_workstreams.md | Create | 40 |
| tests/unit/core/fixtures/invalid_workstreams.md | Create | 30 |
| tests/unit/core/fixtures/project_maps.yaml | Create | 15 |

### Statistics
- **Files Changed:** 8
- **Lines Added:** ~900
- **Lines Removed:** ~0
- **Test Coverage:** ≥80% for core modules
- **Tests Passed:** ______
- **Tests Failed:** ______

### Deviations from Plan
- ______

### Commit
______
