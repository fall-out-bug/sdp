"""Tests for sdp.hooks.ws_complete.

AC1: VerificationResult formatting via _handle_verification_passed/failed
AC2: PostWSCompleteHook (WSCompleteChecker) class methods
AC3: Integration with WSCompletionVerifier (verify_output_files, verify_commands,
     verify_coverage) exercised via verifier.verify() in run()
"""

from pathlib import Path
from unittest.mock import MagicMock

import pytest
from click.testing import CliRunner
from sdp.hooks.ws_complete import HookResult, PostWSCompleteHook, main
from sdp.validators.ws_completion.models import CheckResult, VerificationResult
from sdp.validators.ws_completion.verifier import WSCompletionVerifier


def test_hook_result_dataclass() -> None:
    """HookResult stores passed, ws_id, messages, bypass_used, bypass_reason."""
    result = HookResult(
        passed=True,
        ws_id="00-020-01",
        messages=["OK"],
        bypass_used=False,
        bypass_reason=None,
    )
    assert result.passed is True
    assert result.ws_id == "00-020-01"
    assert result.messages == ["OK"]
    assert result.bypass_used is False
    assert result.bypass_reason is None


def test_post_ws_complete_hook_init() -> None:
    """PostWSCompleteHook stores verifier on init."""
    verifier = MagicMock(spec=WSCompletionVerifier)
    hook = PostWSCompleteHook(verifier)
    assert hook.verifier is verifier


def test_handle_bypass_without_reason_returns_failure() -> None:
    """_handle_bypass returns failure when reason is empty."""
    verifier = MagicMock(spec=WSCompletionVerifier)
    hook = PostWSCompleteHook(verifier)
    result = hook._handle_bypass("00-020-01", "")
    assert result.passed is False
    assert "Bypass requires --reason" in result.messages[0]
    assert result.bypass_used is False
    assert result.bypass_reason is None


def test_handle_bypass_with_reason_returns_success(tmp_path: Path) -> None:
    """_handle_bypass returns success and logs when reason provided."""
    verifier = MagicMock(spec=WSCompletionVerifier)
    hook = PostWSCompleteHook(verifier)
    sdp_dir = tmp_path / ".sdp"
    sdp_dir.mkdir()
    log_file = sdp_dir / "bypass_log.txt"
    with pytest.MonkeyPatch.context() as m:
        m.chdir(tmp_path)
        result = hook._handle_bypass("00-020-01", "Manual override")
    assert result.passed is True
    assert "Bypassing verification" in result.messages[0]
    assert "Manual override" in result.messages[1]
    assert result.bypass_used is True
    assert result.bypass_reason == "Manual override"
    assert log_file.exists()
    assert "00-020-01" in log_file.read_text()
    assert "BYPASS" in log_file.read_text()
    assert "Manual override" in log_file.read_text()


def test_handle_verification_passed_without_coverage() -> None:
    """_handle_verification_passed formats result without coverage."""
    verifier = MagicMock(spec=WSCompletionVerifier)
    hook = PostWSCompleteHook(verifier)
    vr = VerificationResult(
        ws_id="00-020-01",
        passed=True,
        checks=[CheckResult("A", True, "ok", None)],
        coverage_actual=None,
        missing_files=[],
        failed_commands=[],
    )
    result = hook._handle_verification_passed("00-020-01", vr)
    assert result.passed is True
    assert "Verifying" in result.messages[0]
    assert "PASSED" in result.messages[1]
    assert "1/1 passed" in result.messages[2]
    assert len(result.messages) == 3


def test_handle_verification_passed_with_coverage() -> None:
    """_handle_verification_passed includes coverage when present."""
    verifier = MagicMock(spec=WSCompletionVerifier)
    hook = PostWSCompleteHook(verifier)
    vr = VerificationResult(
        ws_id="00-020-01",
        passed=True,
        checks=[CheckResult("A", True, "ok", None)],
        coverage_actual=85.5,
        missing_files=[],
        failed_commands=[],
    )
    result = hook._handle_verification_passed("00-020-01", vr)
    assert result.passed is True
    assert "85.5%" in result.messages[3]


def test_handle_verification_failed_basic() -> None:
    """_handle_verification_failed formats failed checks."""
    verifier = MagicMock(spec=WSCompletionVerifier)
    hook = PostWSCompleteHook(verifier)
    vr = VerificationResult(
        ws_id="00-020-01",
        passed=False,
        checks=[
            CheckResult("File: foo.py", False, "Missing: foo.py", None),
        ],
        coverage_actual=None,
        missing_files=["foo.py"],
        failed_commands=[],
    )
    result = hook._handle_verification_failed("00-020-01", vr)
    assert result.passed is False
    assert "FAILED" in result.messages[1]
    assert "1/1" in result.messages[2]
    assert "foo.py" in " ".join(result.messages)
    assert "Remediation steps" in " ".join(result.messages)
    assert "Create missing files" in " ".join(result.messages)


def test_handle_verification_failed_with_failed_commands() -> None:
    """_handle_verification_failed includes failed commands section."""
    verifier = MagicMock(spec=WSCompletionVerifier)
    hook = PostWSCompleteHook(verifier)
    vr = VerificationResult(
        ws_id="00-020-01",
        passed=False,
        checks=[CheckResult("Command: pytest", False, "Exit code: 1", None)],
        coverage_actual=None,
        missing_files=[],
        failed_commands=["pytest"],
    )
    result = hook._handle_verification_failed("00-020-01", vr)
    assert "Fix failing commands" in " ".join(result.messages)
    assert "pytest" in " ".join(result.messages)


def test_handle_verification_failed_with_low_coverage() -> None:
    """_handle_verification_failed includes coverage remediation when < 80%."""
    verifier = MagicMock(spec=WSCompletionVerifier)
    hook = PostWSCompleteHook(verifier)
    vr = VerificationResult(
        ws_id="00-020-01",
        passed=False,
        checks=[CheckResult("Coverage", False, "Coverage: 65.0%", None)],
        coverage_actual=65.0,
        missing_files=[],
        failed_commands=[],
    )
    result = hook._handle_verification_failed("00-020-01", vr)
    assert "Increase test coverage" in " ".join(result.messages)
    assert "65.0%" in " ".join(result.messages)
    assert "Target: 80%" in " ".join(result.messages)


def test_handle_verification_failed_truncates_more_than_five_checks() -> None:
    """_handle_verification_failed truncates to 5 checks with 'and N more'."""
    verifier = MagicMock(spec=WSCompletionVerifier)
    hook = PostWSCompleteHook(verifier)
    vr = VerificationResult(
        ws_id="00-020-01",
        passed=False,
        checks=[
            CheckResult(f"C{i}", False, f"msg{i}", None) for i in range(7)
        ],
        coverage_actual=None,
        missing_files=[],
        failed_commands=[],
    )
    result = hook._handle_verification_failed("00-020-01", vr)
    assert any("... and 2 more" in m for m in result.messages)


def test_run_with_bypass_no_reason() -> None:
    """run() returns failure when bypass=True and reason empty."""
    verifier = MagicMock(spec=WSCompletionVerifier)
    hook = PostWSCompleteHook(verifier)
    result = hook.run("00-020-01", bypass=True, reason="")
    assert result.passed is False
    verifier.verify.assert_not_called()


def test_run_with_bypass_and_reason(tmp_path: Path) -> None:
    """run() returns success when bypass=True and reason provided."""
    verifier = MagicMock(spec=WSCompletionVerifier)
    hook = PostWSCompleteHook(verifier)
    (tmp_path / ".sdp").mkdir()
    with pytest.MonkeyPatch.context() as m:
        m.chdir(tmp_path)
        result = hook.run("00-020-01", bypass=True, reason="Skip")
    assert result.passed is True
    verifier.verify.assert_not_called()


def test_run_verification_passed() -> None:
    """run() returns success when verifier.verify passes."""
    verifier = MagicMock(spec=WSCompletionVerifier)
    verifier.verify.return_value = VerificationResult(
        ws_id="00-020-01",
        passed=True,
        checks=[CheckResult("A", True, "ok", None)],
        coverage_actual=None,
        missing_files=[],
        failed_commands=[],
    )
    hook = PostWSCompleteHook(verifier)
    result = hook.run("00-020-01", bypass=False)
    assert result.passed is True
    verifier.verify.assert_called_once_with("00-020-01")


def test_run_verification_failed() -> None:
    """run() returns failure when verifier.verify fails."""
    verifier = MagicMock(spec=WSCompletionVerifier)
    verifier.verify.return_value = VerificationResult(
        ws_id="00-020-01",
        passed=False,
        checks=[CheckResult("A", False, "fail", None)],
        coverage_actual=None,
        missing_files=[],
        failed_commands=[],
    )
    hook = PostWSCompleteHook(verifier)
    result = hook.run("00-020-01", bypass=False)
    assert result.passed is False


def test_main_exits_zero_on_pass() -> None:
    """main() exits 0 when verification passes."""
    runner = CliRunner()
    with pytest.MonkeyPatch.context() as m:
        mock_verify = MagicMock(
            return_value=VerificationResult(
                ws_id="00-020-01",
                passed=True,
                checks=[CheckResult("A", True, "ok", None)],
                coverage_actual=None,
                missing_files=[],
                failed_commands=[],
            )
        )
        mock_verifier = MagicMock()
        mock_verifier.verify = mock_verify
        m.setattr("sdp.hooks.ws_complete.WSCompletionVerifier", lambda: mock_verifier)
        result = runner.invoke(main, ["00-020-01"])
    assert result.exit_code == 0
    assert "PASSED" in result.output


def test_main_exits_one_on_fail() -> None:
    """main() exits 1 when verification fails."""
    runner = CliRunner()
    with pytest.MonkeyPatch.context() as m:
        mock_verify = MagicMock(
            return_value=VerificationResult(
                ws_id="00-020-01",
                passed=False,
                checks=[CheckResult("A", False, "fail", None)],
                coverage_actual=None,
                missing_files=[],
                failed_commands=[],
            )
        )
        mock_verifier = MagicMock()
        mock_verifier.verify = mock_verify
        m.setattr("sdp.hooks.ws_complete.WSCompletionVerifier", lambda: mock_verifier)
        result = runner.invoke(main, ["00-020-01"])
    assert result.exit_code == 1
    assert "FAILED" in result.output


def test_main_bypass_with_reason_exits_zero(tmp_path: Path) -> None:
    """main() with --bypass --reason exits 0."""
    runner = CliRunner()
    (tmp_path / ".sdp").mkdir()
    with pytest.MonkeyPatch.context() as m:
        m.chdir(tmp_path)
        result = runner.invoke(main, ["00-020-01", "--bypass", "--reason", "Skip"])
    assert result.exit_code == 0
    assert "Bypassing" in result.output


def test_main_bypass_without_reason_exits_one() -> None:
    """main() with --bypass but no --reason exits 1."""
    runner = CliRunner()
    result = runner.invoke(main, ["00-020-01", "--bypass"])
    assert result.exit_code == 1
    assert "Bypass requires" in result.output
