"""Tests for TDD runner."""

from unittest.mock import Mock, patch

import pytest

from sdp.tdd.runner import TDDRunner, Phase


def test_red_phase_expects_failure(tmp_path):
    """RED phase: test should fail for TDD cycle to start."""
    # Create a test that fails
    test_file = tmp_path / "test_fail.py"
    test_file.write_text("def test_not_impl(): assert False")

    runner = TDDRunner(tmp_path)
    result = runner.red_phase(str(test_file))

    assert result.phase == Phase.RED
    assert result.success is True  # Failure is expected in RED phase!
    assert "FAILED" in result.output
    assert result.next_phase == Phase.GREEN


def test_green_phase_expects_success(tmp_path):
    """GREEN phase: test should pass after implementation."""
    test_file = tmp_path / "test_pass.py"
    test_file.write_text("def test_pass(): assert True")

    runner = TDDRunner(tmp_path)
    result = runner.green_phase(str(test_file))

    assert result.phase == Phase.GREEN
    assert result.success is True
    assert "PASSED" in result.output
    assert result.next_phase == Phase.REFACTOR


def test_refactor_phase_maintain_passing(tmp_path):
    """REFACTOR phase: tests should still pass after refactoring."""
    test_file = tmp_path / "test_refactor.py"
    test_file.write_text("def test_still_passes(): assert True")

    runner = TDDRunner(tmp_path)
    result = runner.refactor_phase(str(test_file))

    assert result.phase == Phase.REFACTOR
    assert result.success is True
    assert "PASSED" in result.output
    assert result.next_phase is None  # Cycle complete


def test_green_phase_failure_stays_green(tmp_path):
    """GREEN phase: if test fails, stay in GREEN phase."""
    test_file = tmp_path / "test_still_fails.py"
    test_file.write_text("def test_not_fixed_yet(): assert False")

    runner = TDDRunner(tmp_path)
    result = runner.green_phase(str(test_file))

    assert result.phase == Phase.GREEN
    assert result.success is False  # Test didn't pass
    assert result.next_phase == Phase.GREEN  # Stay in GREEN


def test_runner_with_retry_disabled(tmp_path):
    """Test runner can disable retry functionality."""
    test_file = tmp_path / "test_pass.py"
    test_file.write_text("def test_pass(): assert True")

    runner = TDDRunner(tmp_path, enable_retry=False)
    result = runner.green_phase(str(test_file))

    assert result.phase == Phase.GREEN
    assert result.success is True


def test_runner_custom_max_retries(tmp_path):
    """Test runner can configure custom max retries."""
    test_file = tmp_path / "test_pass.py"
    test_file.write_text("def test_pass(): assert True")

    runner = TDDRunner(tmp_path, max_retries=5)
    result = runner.green_phase(str(test_file))

    assert result.phase == Phase.GREEN
    assert result.success is True


def test_retry_integration_with_tdd(tmp_path):
    """Test that retry logic integrates properly with TDD phases."""
    from unittest.mock import patch

    test_file = tmp_path / "test_flaky.py"
    test_file.write_text("def test_flaky(): assert True")

    # Mock run_with_retry to verify it's called
    with patch("sdp.tdd.runner.run_with_retry") as mock_retry:
        mock_retry.return_value = Mock(
            success=True,
            exit_code=0,
            stdout="PASSED",
            stderr="",
            attempts=1,
            total_delay=0.0,
        )

        runner = TDDRunner(tmp_path, enable_retry=True)
        result = runner.green_phase(str(test_file))

        assert result.success is True
        mock_retry.assert_called_once()


def test_direct_execution_without_retry(tmp_path):
    """Test that direct execution works when retry is disabled."""
    from unittest.mock import patch

    test_file = tmp_path / "test_direct.py"
    test_file.write_text("def test_direct(): assert True")

    # Mock subprocess.run to verify it's called directly
    with patch("subprocess.run") as mock_subprocess:
        mock_subprocess.return_value = Mock(
            returncode=0,
            stdout="PASSED",
            stderr="",
        )

        runner = TDDRunner(tmp_path, enable_retry=False)
        result = runner.green_phase(str(test_file))

        assert result.success is True
        mock_subprocess.assert_called_once()

