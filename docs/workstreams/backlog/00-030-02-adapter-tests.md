---
ws_id: 00-030-02
feature: F030
status: backlog
size: MEDIUM
project_id: 00
github_issue: null
assignee: null
depends_on: []
---

## WS-00-030-02: Adapter Tests (Claude Code, OpenCode)

### 🎯 Goal

**What must WORK after completing this WS:**
- All adapter modules have comprehensive test coverage
- Claude Code adapter file operations tested with realistic scenarios
- OpenCode adapter configuration tested
- Error handling tested for file system failures

**Acceptance Criteria:**
- [ ] AC1: `tests/unit/adapters/test_base.py` with common adapter tests
- [ ] AC2: `tests/unit/adapters/test_claude_code.py` with ≥80% coverage
- [ ] AC3: `tests/unit/adapters/test_opencode.py` with ≥80% coverage
- [ ] AC4: Tests cover file I/O operations (read, write, validate)
- [ ] AC5: Tests cover error cases (permissions, invalid formats)
- [ ] AC6: mypy --strict passes on all test files

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Critical Gap**: Adapters have **1,003 LOC across 5 modules with ZERO tests**.

**High-Risk Modules**:
1. **`src/sdp/adapters/claude_code.py`** (~232 LOC)
   - File system operations (read .claude/settings.json)
   - Configuration validation
   - Message formatting for Claude Code

2. **`src/sdp/adapters/opencode.py`** (~202 LOC)
   - OpenCode-specific configuration
   - API wrapper functionality

3. **`src/sdp/adapters/base.py`** (~180 LOC)
   - Base adapter class
   - Common interface

**Risk**: Adapters are integration points between SDP and AI IDEs. Untested adapters = unreliable integration.

### Dependencies

None (can proceed independently)

### Input Files

- `src/sdp/adapters/base.py`
- `src/sdp/adapters/claude_code.py`
- `src/sdp/adapters/opencode.py`
- `src/sdp/adapters/codex.py` (if exists)
- `src/sdp/adapters/__init__.py`

### Steps

1. **Create test directory structure**

   ```
   tests/
   └── unit/
       └── adapters/
           ├── __init__.py
           ├── test_base.py
           ├── test_claude_code.py
           ├── test_opencode.py
           └── fixtures/
               ├── claude_settings.json  # Sample Claude Code config
               └── opencode_config.yaml   # Sample OpenCode config
   ```

2. **Test base adapter class**

   ```python
   # tests/unit/adapters/test_base.py
   import pytest
   from pathlib import Path
   from sdp.adapters.base import BaseAdapter

   class TestBaseAdapter:
       """Test base adapter functionality."""

       def test_adapter_has_common_interface(self):
           """Verify all adapters implement common interface."""
           from sdp.adapters.claude_code import ClaudeCodeAdapter
           from sdp.adapters.opencode import OpenCodeAdapter

           # All adapters should have these methods
           required_methods = [
               "read_config",
               "validate_config",
               "format_message",
               "send_message",
           ]

           for adapter_class in [ClaudeCodeAdapter, OpenCodeAdapter]:
               for method in required_methods:
                   assert hasattr(adapter_class, method)

       def test_base_adapter_abstract_methods(self):
           """Verify BaseAdapter cannot be instantiated directly."""
           from sdp.adapters.base import BaseAdapter
           import abc

           # BaseAdapter should be abstract
           with pytest.raises(TypeError):
               BaseAdapter()
   ```

3. **Test Claude Code adapter**

   ```python
   # tests/unit/adapters/test_claude_code.py
   import pytest
   from pathlib import Path
   from unittest.mock import Mock, patch, mock_open
   from sdp.adapters.claude_code import ClaudeCodeAdapter
   from sdp.adapters.exceptions import ConfigValidationError

   @pytest.fixture
   def claude_adapter():
       """Create Claude Code adapter instance."""
       return ClaudeCodeAdapter()

   @pytest.fixture
   def mock_settings_file(tmp_path):
       """Create mock .claude/settings.json file."""
       settings = tmp_path / ".claude" / "settings.json"
       settings.parent.mkdir(parents=True)
       settings.write_text('{"model": "claude-sonnet-4-5"}')
       return settings

   class TestClaudeCodeConfigReading:
       """Test Claude Code configuration reading."""

       def test_read_config_finds_settings_json(self, claude_adapter, mock_settings_file):
           """Verify adapter finds .claude/settings.json."""
           config = claude_adapter.read_config(mock_settings_file.parent)

           assert config["model"] == "claude-sonnet-4-5"

       def test_read_config_raises_error_for_missing_file(self, claude_adapter, tmp_path):
           """Verify error raised when .claude/settings.json missing."""
           non_existent = tmp_path / ".claude"

           with pytest.raises(FileNotFoundError):
               claude_adapter.read_config(non_existent)

       def test_read_config_validates_json_format(self, claude_adapter, tmp_path):
           """Verify config validation rejects invalid JSON."""
           invalid_file = tmp_path / ".claude" / "settings.json"
           invalid_file.parent.mkdir(parents=True)
           invalid_file.write_text('{invalid json}')

           with pytest.raises(ConfigValidationError):
               claude_adapter.read_config(tmp_path)

   class TestClaudeCodeMessageFormatting:
       """Test Claude Code message formatting."""

       def test_format_message_for_claude_code(self, claude_adapter):
           """Verify message is formatted correctly for Claude Code."""
           message = {
               "type": "message",
               "recipient": "builder",
               "content": "Execute WS-00-001-01",
           }

           formatted = claude_adapter.format_message(message)

           # Claude Code expects specific format
           assert "recipient" in formatted
           assert "Execute WS-00-001-01" in formatted

       def test_format_message_includes_context(self, claude_adapter):
           """Verify message includes context for recipient."""
           message = {
               "type": "message",
               "recipient": "builder",
               "content": "Task update",
               "context": {"ws_id": "00-001-01", "status": "in_progress"},
           }

           formatted = claude_adapter.format_message(message)

           assert "00-001-01" in formatted
           assert "in_progress" in formatted

   class TestClaudeCodeErrorHandling:
       """Test Claude Code adapter error handling."""

       def test_handle_permission_denied(self, claude_adapter, tmp_path):
           """Verify permission errors are handled gracefully."""
           read_only_file = tmp_path / "readonly.txt"
           read_only_file.write_text("read only")
           read_only_file.chmod(0o444)  # Read-only

           with pytest.raises(PermissionError):
               claude_adapter.read_config(read_only_file)

       def test_validate_config_rejects_invalid_model(self, claude_adapter):
           """Verify config validation rejects unsupported models."""
           invalid_config = {"model": "gpt-4"}  # Not a Claude model

           with pytest.raises(ConfigValidationError, match="Invalid model"):
               claude_adapter.validate_config(invalid_config)
   ```

4. **Test OpenCode adapter**

   ```python
   # tests/unit/adapters/test_opencode.py
   import pytest
   from pathlib import Path
   from sdp.adapters.opencode import OpenCodeAdapter

   @pytest.fixture
   def opencode_adapter():
       """Create OpenCode adapter instance."""
       return OpenCodeAdapter()

   class TestOpenCodeConfig:
       """Test OpenCode configuration."""

       def test_read_opencode_yaml_config(self, opencode_adapter, tmp_path):
           """Verify adapter reads OpenCode YAML config."""
           config_file = tmp_path / ".opencode" / "config.yaml"
           config_file.parent.mkdir(parents=True)
           config_file.write_text("model: opencode-beta\napi_key: test123")

           config = opencode_adapter.read_config(tmp_path)

           assert config["model"] == "opencode-beta"
           assert config["api_key"] == "test123"

       def test_validate_config_requires_api_key(self, opencode_adapter):
           """Verify config validation requires api_key."""
           config_without_key = {"model": "opencode-beta"}

           with pytest.raises(ConfigValidationError, match="api_key.*required"):
               opencode_adapter.validate_config(config_without_key)

   class TestOpenCodeIntegration:
       """Test OpenCode integration scenarios."""

       def test_send_message_to_opencode(self, opencode_adapter):
           """Verify message sending to OpenCode API."""
           # Mock OpenCode API call
           with patch('sdp.adapters.opencode.requests.post') as mock_post:
               mock_post.return_value = Mock(status_code=200)

               result = opencode_adapter.send_message(
                   recipient="builder",
                   content="Execute WS-00-001-01",
               )

               assert result is True
               mock_post.assert_called_once()

       def test_send_message_handles_api_failure(self, opencode_adapter):
           """Verify API failure is handled gracefully."""
           import requests

           with patch('sdp.adapters.opencode.requests.post') as mock_post:
               mock_post.side_effect = requests.exceptions.ConnectionError("API down")

               with pytest.raises(ConnectionError):
                   opencode_adapter.send_message(
                       recipient="builder",
                       content="Test message",
                   )
   ```

### Code

```python
# src/sdp/adapters/base.py (example of what to test)
from abc import ABC, abstractmethod
from pathlib import Path
from typing import Any, Dict

class BaseAdapter(ABC):
    """Base class for AI-IDE adapters."""

    @abstractmethod
    def read_config(self, project_root: Path) -> Dict[str, Any]:
        """Read adapter configuration from project.

        Args:
            project_root: Path to project root directory

        Returns:
            Configuration dictionary

        Raises:
            FileNotFoundError: If config file missing
            ConfigValidationError: If config format invalid
        """
        pass

    @abstractmethod
    def validate_config(self, config: Dict[str, Any]) -> None:
        """Validate adapter configuration.

        Args:
            config: Configuration dictionary to validate

        Raises:
            ConfigValidationError: If config invalid
        """
        pass

    @abstractmethod
    def format_message(self, message: Dict[str, Any]) -> str:
        """Format message for adapter's specific format.

        Args:
            message: Message dict with type, recipient, content

        Returns:
            Formatted message string
        """
        pass
```

### Expected Outcome

**After completion:**
- Adapters (1,003 LOC) now have test coverage
- Integration points between SDP and AI IDEs tested
- Error handling verified for file system failures
- Foundation for adding new adapters with confidence

**Scope Estimate**
- Files: ~7
- Lines: ~650 (MEDIUM)
- Tokens: ~3250

### Completion Criteria

```bash
# Run all adapter tests
pytest tests/unit/adapters/ -v

# Verify coverage
pytest --cov=src/sdp/adapters --cov-report=term-missing
# Should show ≥80% for claude_code.py, opencode.py, base.py

# Verify type checking
mypy src/sdp/adapters/ --strict

# Run specific test class
pytest tests/unit/adapters/test_claude_code.py::TestClaudeCodeConfigReading -v
```

### Constraints

- DO NOT test external adapter code (only our wrappers)
- DO NOT make real API calls to OpenCode/Claude Code
- DO NOT hardcode file paths (use tmp_path fixtures)
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
| tests/unit/adapters/__init__.py | Create | 5 |
| tests/unit/adapters/test_base.py | Create | 120 |
| tests/unit/adapters/test_claude_code.py | Create | 250 |
| tests/unit/adapters/test_opencode.py | Create | 200 |
| tests/unit/adapters/fixtures/claude_settings.json | Create | 10 |
| tests/unit/adapters/fixtures/opencode_config.yaml | Create | 15 |

### Statistics
- **Files Changed:** 6
- **Lines Added:** ~600
- **Lines Removed:** ~0
- **Test Coverage:** ≥80% for adapters
- **Tests Passed:** ______
- **Tests Failed:** ______

### Deviations from Plan
- ______

### Commit
______
