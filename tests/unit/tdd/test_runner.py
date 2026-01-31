"""Tests for TDD runner."""

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
