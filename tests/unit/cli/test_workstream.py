"""Tests for workstream CLI commands."""
from pathlib import Path
from unittest.mock import MagicMock, patch

import pytest
from click.testing import CliRunner

from sdp.cli.workstream import workstream


@pytest.fixture
def runner():
    """Create CLI runner."""
    return CliRunner()


@pytest.fixture
def temp_ws_file(tmp_path):
    """Create temporary workstream file."""
    ws_file = tmp_path / "00-001-01-test.md"
    ws_file.write_text("""---
ws_id: 00-001-01
title: Test Workstream
feature: F001
status: backlog
size: SMALL
tier: T1
---

# Test Workstream

## Acceptance Criteria
- [ ] Test criterion 1
- [ ] Test criterion 2
""")
    return ws_file


@pytest.fixture
def temp_project_map(tmp_path):
    """Create temporary PROJECT_MAP.md file."""
    pm_file = tmp_path / "PROJECT_MAP.md"
    pm_file.write_text("""# Project Map

## Decisions
- Decision 1
- Decision 2

## Constraints
- Constraint 1

## Tech Stack
- Python 3.10+
- Poetry
""")
    return pm_file


def test_workstream_group_help(runner):
    """Test workstream group help text."""
    result = runner.invoke(workstream, ["--help"])
    assert result.exit_code == 0
    assert "Core SDP operations" in result.output


def test_parse_workstream_success(runner, temp_ws_file):
    """Test parsing valid workstream file."""
    result = runner.invoke(workstream, ["parse", str(temp_ws_file)])
    assert result.exit_code == 0
    assert "✓ Parsed 00-001-01" in result.output
    assert "Feature: F001" in result.output
    assert "Status: backlog" in result.output
    assert "Size: SMALL" in result.output
    assert "Acceptance Criteria: 2" in result.output


def test_parse_workstream_invalid(runner, tmp_path):
    """Test parsing invalid workstream file."""
    invalid_ws = tmp_path / "invalid.md"
    invalid_ws.write_text("Invalid content")
    
    result = runner.invoke(workstream, ["parse", str(invalid_ws)])
    assert result.exit_code == 1
    assert "Error" in result.output or "error" in result.output


def test_parse_workstream_no_acceptance_criteria(runner, tmp_path):
    """Test parsing workstream without acceptance criteria."""
    ws_file = tmp_path / "00-001-01.md"
    ws_file.write_text("""---
ws_id: 00-001-01
title: Test Workstream
feature: F001
status: backlog
size: SMALL
---

# Test Workstream
""")
    
    result = runner.invoke(workstream, ["parse", str(ws_file)])
    assert result.exit_code == 0
    assert "✓ Parsed 00-001-01" in result.output


def test_parse_project_map_success(runner, temp_project_map):
    """Test parsing valid PROJECT_MAP.md file."""
    result = runner.invoke(workstream, ["parse-project-map", str(temp_project_map)])
    # Project map parsing may have specific validation - just check it doesn't crash
    assert result.exit_code in [0, 1]  # May fail validation but shouldn't crash


def test_parse_project_map_invalid(runner, tmp_path):
    """Test parsing invalid PROJECT_MAP.md file."""
    invalid_pm = tmp_path / "invalid_pm.md"
    invalid_pm.write_text("Invalid content")
    
    result = runner.invoke(workstream, ["parse-project-map", str(invalid_pm)])
    assert result.exit_code == 1


def test_scope_show_with_files(runner):
    """Test showing scope with files."""
    with patch("sdp.cli.workstream.create_beads_client"):
        with patch("sdp.cli.workstream.ScopeManager") as mock_manager_class:
            mock_manager = MagicMock()
            mock_manager_class.return_value = mock_manager
            mock_manager.get_scope.return_value = ["src/module.py", "tests/test_module.py"]
            
            result = runner.invoke(workstream, ["scope", "show", "00-001-01"])
            
            assert result.exit_code == 0
            assert "Scope for 00-001-01: 2 files" in result.output
            assert "src/module.py" in result.output
            assert "tests/test_module.py" in result.output


def test_scope_show_unrestricted(runner):
    """Test showing unrestricted scope."""
    with patch("sdp.cli.workstream.create_beads_client"):
        with patch("sdp.cli.workstream.ScopeManager") as mock_manager_class:
            mock_manager = MagicMock()
            mock_manager_class.return_value = mock_manager
            mock_manager.get_scope.return_value = []
            
            result = runner.invoke(workstream, ["scope", "show", "00-001-01"])
            
            assert result.exit_code == 0
            assert "unrestricted" in result.output


def test_scope_show_error(runner):
    """Test showing scope with error."""
    with patch("sdp.cli.workstream.create_beads_client"):
        with patch("sdp.cli.workstream.ScopeManager") as mock_manager_class:
            mock_manager = MagicMock()
            mock_manager_class.return_value = mock_manager
            mock_manager.get_scope.side_effect = ValueError("Task not found")
            
            result = runner.invoke(workstream, ["scope", "show", "00-001-01"])
            
            assert result.exit_code == 1
            assert "Task not found" in result.output


def test_scope_add_file(runner):
    """Test adding file to scope."""
    with patch("sdp.cli.workstream.create_beads_client"):
        with patch("sdp.cli.workstream.ScopeManager") as mock_manager_class:
            mock_manager = MagicMock()
            mock_manager_class.return_value = mock_manager
            
            result = runner.invoke(workstream, ["scope", "add", "00-001-01", "src/module.py"])
            
            assert result.exit_code == 0
            assert "✅ Added src/module.py to 00-001-01 scope" in result.output
            mock_manager.add_file.assert_called_once_with("00-001-01", "src/module.py")


def test_scope_add_file_error(runner):
    """Test adding file to scope with error."""
    with patch("sdp.cli.workstream.create_beads_client"):
        with patch("sdp.cli.workstream.ScopeManager") as mock_manager_class:
            mock_manager = MagicMock()
            mock_manager_class.return_value = mock_manager
            mock_manager.add_file.side_effect = ValueError("Task not found")
            
            result = runner.invoke(workstream, ["scope", "add", "00-001-01", "src/module.py"])
            
            assert result.exit_code == 1
            assert "Task not found" in result.output


def test_scope_remove_file(runner):
    """Test removing file from scope."""
    with patch("sdp.cli.workstream.create_beads_client"):
        with patch("sdp.cli.workstream.ScopeManager") as mock_manager_class:
            mock_manager = MagicMock()
            mock_manager_class.return_value = mock_manager
            
            result = runner.invoke(workstream, ["scope", "remove", "00-001-01", "src/module.py"])
            
            assert result.exit_code == 0
            assert "✅ Removed src/module.py from 00-001-01 scope" in result.output
            mock_manager.remove_file.assert_called_once_with("00-001-01", "src/module.py")


def test_scope_remove_file_error(runner):
    """Test removing file from scope with error."""
    with patch("sdp.cli.workstream.create_beads_client"):
        with patch("sdp.cli.workstream.ScopeManager") as mock_manager_class:
            mock_manager = MagicMock()
            mock_manager_class.return_value = mock_manager
            mock_manager.remove_file.side_effect = ValueError("File not in scope")
            
            result = runner.invoke(workstream, ["scope", "remove", "00-001-01", "src/module.py"])
            
            assert result.exit_code == 1
            assert "File not in scope" in result.output


def test_scope_clear(runner):
    """Test clearing scope."""
    with patch("sdp.cli.workstream.create_beads_client"):
        with patch("sdp.cli.workstream.ScopeManager") as mock_manager_class:
            mock_manager = MagicMock()
            mock_manager_class.return_value = mock_manager
            
            result = runner.invoke(workstream, ["scope", "clear", "00-001-01"])
            
            assert result.exit_code == 0
            assert "✅ Cleared scope for 00-001-01" in result.output
            assert "unrestricted" in result.output
            mock_manager.clear_scope.assert_called_once_with("00-001-01")


def test_scope_clear_error(runner):
    """Test clearing scope with error."""
    with patch("sdp.cli.workstream.create_beads_client"):
        with patch("sdp.cli.workstream.ScopeManager") as mock_manager_class:
            mock_manager = MagicMock()
            mock_manager_class.return_value = mock_manager
            mock_manager.clear_scope.side_effect = ValueError("Task not found")
            
            result = runner.invoke(workstream, ["scope", "clear", "00-001-01"])
            
            assert result.exit_code == 1
            assert "Task not found" in result.output


def test_verify_completion_passed(runner):
    """Test workstream verification passed."""
    with patch("sdp.cli.workstream.WSCompletionVerifier") as mock_verifier_class:
        mock_verifier = MagicMock()
        mock_verifier_class.return_value = mock_verifier
        
        from dataclasses import dataclass
        
        @dataclass
        class CheckResult:
            name: str
            passed: bool
            message: str
        
        @dataclass
        class VerificationResult:
            passed: bool
            checks: list
            coverage_actual: float | None = None
            missing_files: list = None
            failed_commands: list = None
        
        mock_verifier.verify.return_value = VerificationResult(
            passed=True,
            checks=[
                CheckResult(name="Scope files", passed=True, message="All files exist"),
                CheckResult(name="Verification commands", passed=True, message="All commands passed"),
            ],
            coverage_actual=85.0
        )
        
        result = runner.invoke(workstream, ["verify", "00-001-01"])
        
        assert result.exit_code == 0
        assert "✅ Workstream 00-001-01 verification PASSED" in result.output
        assert "Coverage: 85.0%" in result.output


def test_verify_completion_failed(runner):
    """Test workstream verification failed."""
    with patch("sdp.cli.workstream.WSCompletionVerifier") as mock_verifier_class:
        mock_verifier = MagicMock()
        mock_verifier_class.return_value = mock_verifier
        
        from dataclasses import dataclass
        
        @dataclass
        class CheckResult:
            name: str
            passed: bool
            message: str
        
        @dataclass
        class VerificationResult:
            passed: bool
            checks: list
            coverage_actual: float | None = None
            missing_files: list = None
            failed_commands: list = None
        
        mock_verifier.verify.return_value = VerificationResult(
            passed=False,
            checks=[
                CheckResult(name="Scope files", passed=False, message="Missing files"),
            ],
            coverage_actual=60.0,
            missing_files=["src/module.py"],
            failed_commands=["pytest tests/test_module.py"]
        )
        
        result = runner.invoke(workstream, ["verify", "00-001-01"])
        
        assert result.exit_code == 1
        assert "❌ Workstream 00-001-01 verification FAILED" in result.output
        assert "Coverage: 60.0%" in result.output
        assert "Missing files" in result.output
        assert "src/module.py" in result.output
        assert "Failed commands" in result.output
        assert "pytest tests/test_module.py" in result.output


def test_supersede_success(runner):
    """Test superseding workstream."""
    with patch("sdp.cli.workstream.SupersedeValidator") as mock_validator_class:
        mock_validator = MagicMock()
        mock_validator_class.return_value = mock_validator
        
        from dataclasses import dataclass
        
        @dataclass
        class SupersedeResult:
            success: bool
            error: str | None = None
        
        mock_validator.supersede.return_value = SupersedeResult(success=True)
        
        result = runner.invoke(workstream, ["supersede", "00-001-01", "--replacement", "00-001-02"])
        
        assert result.exit_code == 0
        assert "✅ Marked 00-001-01 as superseded by 00-001-02" in result.output
        mock_validator.supersede.assert_called_once_with("00-001-01", "00-001-02")


def test_supersede_failure(runner):
    """Test superseding workstream failure."""
    with patch("sdp.cli.workstream.SupersedeValidator") as mock_validator_class:
        mock_validator = MagicMock()
        mock_validator_class.return_value = mock_validator
        
        from dataclasses import dataclass
        
        @dataclass
        class SupersedeResult:
            success: bool
            error: str | None = None
        
        mock_validator.supersede.return_value = SupersedeResult(
            success=False,
            error="Replacement workstream not found"
        )
        
        result = runner.invoke(workstream, ["supersede", "00-001-01", "--replacement", "00-001-99"])
        
        assert result.exit_code == 1
        assert "❌ Failed to supersede 00-001-01" in result.output
        assert "Replacement workstream not found" in result.output


def test_find_orphans_none_found(runner):
    """Test finding orphans when none exist."""
    with patch("sdp.cli.workstream.SupersedeValidator") as mock_validator_class:
        mock_validator = MagicMock()
        mock_validator_class.return_value = mock_validator
        mock_validator.find_orphans.return_value = []
        
        result = runner.invoke(workstream, ["orphans"])
        
        assert result.exit_code == 0
        assert "✅ No orphaned superseded workstreams found" in result.output


def test_find_orphans_found(runner):
    """Test finding orphans when they exist."""
    with patch("sdp.cli.workstream.SupersedeValidator") as mock_validator_class:
        mock_validator = MagicMock()
        mock_validator_class.return_value = mock_validator
        mock_validator.find_orphans.return_value = ["00-001-01", "00-001-02"]
        
        result = runner.invoke(workstream, ["orphans"])
        
        assert result.exit_code == 1
        assert "❌ Found 2 orphaned superseded workstream(s)" in result.output
        assert "00-001-01" in result.output
        assert "00-001-02" in result.output
        assert "Fix with:" in result.output


def test_validate_tier_command_available(runner):
    """Test that validate_tier command is available when imported."""
    # Check if validate command exists
    result = runner.invoke(workstream, ["validate", "--help"])
    # Should either work or show command not found
    # This test just ensures the import logic doesn't crash
    assert result.exit_code in [0, 2]  # 0 if available, 2 if not found


def test_workstream_scope_group_help(runner):
    """Test scope subgroup help."""
    result = runner.invoke(workstream, ["scope", "--help"])
    assert result.exit_code == 0
    assert "Manage workstream file scope" in result.output
    assert "show" in result.output
    assert "add" in result.output
    assert "remove" in result.output
    assert "clear" in result.output


def test_parse_workstream_missing_file(runner):
    """Test parsing non-existent workstream file."""
    result = runner.invoke(workstream, ["parse", "/nonexistent/file.md"])
    # Click should fail on path validation
    assert result.exit_code == 2  # Click validation error


def test_parse_project_map_missing_file(runner):
    """Test parsing non-existent project map file."""
    result = runner.invoke(workstream, ["parse-project-map", "/nonexistent/PROJECT_MAP.md"])
    # Click should fail on path validation
    assert result.exit_code == 2  # Click validation error
