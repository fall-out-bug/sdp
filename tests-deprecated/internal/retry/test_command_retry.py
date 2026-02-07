"""Tests for command retry with exponential backoff."""

from unittest.mock import Mock, patch

import pytest
from sdp.internal.retry.command_retry import (
    CommandRetryError,
    RetryConfig,
    RetryResult,
    run_with_retry,
)


class TestRetryConfig:
    """Test RetryConfig dataclass."""

    def test_default_config(self) -> None:
        """Test default configuration values."""
        config = RetryConfig()

        assert config.max_retries == 3
        assert config.base_delay == 1.0
        assert config.backoff_factor == 2.0

    def test_custom_config(self) -> None:
        """Test custom configuration values."""
        config = RetryConfig(
            max_retries=5,
            base_delay=2.0,
            backoff_factor=3.0,
        )

        assert config.max_retries == 5
        assert config.base_delay == 2.0
        assert config.backoff_factor == 3.0


class TestRunWithRetry:
    """Test run_with_retry function."""

    def test_success_on_first_attempt(self) -> None:
        """Test command succeeds on first attempt."""
        with patch("subprocess.run") as mock_run:
            mock_run.return_value = Mock(
                returncode=0,
                stdout="success",
                stderr="",
            )

            result = run_with_retry(["echo", "hello"])

            assert result.success is True
            assert result.exit_code == 0
            assert result.stdout == "success"
            assert result.attempts == 1
            assert result.total_delay == 0.0
            mock_run.assert_called_once()

    def test_success_on_second_attempt(self) -> None:
        """Test command succeeds on second attempt."""
        with patch("subprocess.run") as mock_run:
            # First attempt fails, second succeeds
            mock_run.side_effect = [
                Mock(returncode=1, stdout="", stderr="error"),
                Mock(returncode=0, stdout="success", stderr=""),
            ]

            result = run_with_retry(
                ["echo", "hello"],
                config=RetryConfig(max_retries=3, base_delay=0.1),
            )

            assert result.success is True
            assert result.exit_code == 0
            assert result.attempts == 2
            assert result.total_delay == 0.1  # base_delay * 2^0
            assert mock_run.call_count == 2

    def test_success_on_third_attempt(self) -> None:
        """Test command succeeds on third attempt."""
        with patch("subprocess.run") as mock_run:
            # First two attempts fail, third succeeds
            mock_run.side_effect = [
                Mock(returncode=1, stdout="", stderr="error1"),
                Mock(returncode=1, stdout="", stderr="error2"),
                Mock(returncode=0, stdout="success", stderr=""),
            ]

            result = run_with_retry(
                ["echo", "hello"],
                config=RetryConfig(max_retries=3, base_delay=0.1),
            )

            assert result.success is True
            assert result.attempts == 3
            # 0.1 * (2^0 + 2^1) = 0.1 * 3 = 0.3
            assert result.total_delay == pytest.approx(0.3)

    def test_failure_after_max_retries(self) -> None:
        """Test command fails after max retries."""
        with patch("subprocess.run") as mock_run:
            # All attempts fail
            mock_run.return_value = Mock(
                returncode=1,
                stdout="",
                stderr="persistent error",
            )

            with pytest.raises(CommandRetryError) as exc_info:
                run_with_retry(
                    ["false"],
                    config=RetryConfig(max_retries=2, base_delay=0.1),
                )

            error = exc_info.value
            assert error.last_exit_code == 1
            assert error.attempts == 3  # initial + 2 retries
            assert error.total_delay == pytest.approx(0.3)  # 0.1 * (1 + 2)
            assert mock_run.call_count == 3

    def test_exponential_backoff_delays(self) -> None:
        """Test exponential backoff delay calculation."""
        with patch("subprocess.run") as mock_run, \
             patch("time.sleep") as mock_sleep:
            # All attempts fail
            mock_run.return_value = Mock(returncode=1, stdout="", stderr="")

            with pytest.raises(CommandRetryError):
                run_with_retry(
                    ["false"],
                    config=RetryConfig(max_retries=3, base_delay=1.0),
                )

            # Check sleep calls: 1s, 2s, 4s
            assert mock_sleep.call_count == 3
            mock_sleep.assert_any_call(1.0)
            mock_sleep.assert_any_call(2.0)
            mock_sleep.assert_any_call(4.0)

    def test_custom_backoff_factor(self) -> None:
        """Test custom backoff factor."""
        with patch("subprocess.run") as mock_run, \
             patch("time.sleep") as mock_sleep:
            # All attempts fail
            mock_run.return_value = Mock(returncode=1, stdout="", stderr="")

            with pytest.raises(CommandRetryError):
                run_with_retry(
                    ["false"],
                    config=RetryConfig(
                        max_retries=2,
                        base_delay=1.0,
                        backoff_factor=3.0,
                    ),
                )

            # Check sleep calls: 1s, 3s (factor 3)
            assert mock_sleep.call_count == 2
            mock_sleep.assert_any_call(1.0)
            mock_sleep.assert_any_call(3.0)

    def test_zero_max_retries(self) -> None:
        """Test with max_retries=0 (no retries)."""
        with patch("subprocess.run") as mock_run:
            mock_run.return_value = Mock(
                returncode=1,
                stdout="",
                stderr="error",
            )

            with pytest.raises(CommandRetryError) as exc_info:
                run_with_retry(
                    ["false"],
                    config=RetryConfig(max_retries=0),
                )

            # Should only attempt once (initial attempt, no retries)
            assert exc_info.value.attempts == 1
            assert mock_run.call_count == 1

    def test_stdout_stderr_capture(self) -> None:
        """Test stdout and stderr capture."""
        with patch("subprocess.run") as mock_run:
            mock_run.return_value = Mock(
                returncode=0,
                stdout="output line\n",
                stderr="warning message\n",
            )

            result = run_with_retry(["test"])

            assert result.stdout == "output line\n"
            assert result.stderr == "warning message\n"

    def test_subprocess_exception_propagates(self) -> None:
        """Test that subprocess exceptions are not caught."""
        with patch("subprocess.run") as mock_run:
            mock_run.side_effect = FileNotFoundError("Command not found")

            with pytest.raises(FileNotFoundError):
                run_with_retry(["nonexistent"])

    def test_telemetry_callback(self) -> None:
        """Test telemetry callback is called for each attempt."""
        telemetry_calls = []

        def telemetry_callback(result: RetryResult) -> None:
            telemetry_calls.append(result)

        with patch("subprocess.run") as mock_run:
            # First attempt fails, second succeeds
            mock_run.side_effect = [
                Mock(returncode=1, stdout="", stderr="error"),
                Mock(returncode=0, stdout="success", stderr=""),
            ]

            run_with_retry(
                ["test"],
                config=RetryConfig(max_retries=3, base_delay=0.1),
                on_retry=telemetry_callback,
            )

            # Should be called twice (first attempt + retry)
            assert len(telemetry_calls) == 2
            assert telemetry_calls[0].success is False
            assert telemetry_calls[1].success is True

    def test_working_directory(self) -> None:
        """Test command execution in custom working directory."""
        with patch("subprocess.run") as mock_run:
            mock_run.return_value = Mock(
                returncode=0,
                stdout="success",
                stderr="",
            )

            result = run_with_retry(
                ["pwd"],
                cwd="/tmp",
            )

            assert result.success is True
            mock_run.assert_called_once()
            call_kwargs = mock_run.call_args[1]
            assert call_kwargs["cwd"] == "/tmp"

    def test_timeout_passed_to_subprocess(self) -> None:
        """Test timeout is passed to subprocess."""
        with patch("subprocess.run") as mock_run:
            mock_run.return_value = Mock(
                returncode=0,
                stdout="success",
                stderr="",
            )

            result = run_with_retry(
                ["sleep", "1"],
                timeout=5.0,
            )

            assert result.success is True
            mock_run.assert_called_once()
            call_kwargs = mock_run.call_args[1]
            assert call_kwargs["timeout"] == 5.0


class TestRetryResult:
    """Test RetryResult dataclass."""

    def test_result_fields(self) -> None:
        """Test all result fields are populated."""
        result = RetryResult(
            success=True,
            exit_code=0,
            stdout="output",
            stderr="error",
            attempts=2,
            total_delay=1.5,
        )

        assert result.success is True
        assert result.exit_code == 0
        assert result.stdout == "output"
        assert result.stderr == "error"
        assert result.attempts == 2
        assert result.total_delay == 1.5


class TestCommandRetryError:
    """Test CommandRetryError exception."""

    def test_error_attributes(self) -> None:
        """Test error exception attributes."""
        error = CommandRetryError(
            command=["test"],
            last_exit_code=1,
            attempts=3,
            total_delay=3.0,
            last_output="error output",
        )

        assert str(error) == "Command ['test'] failed after 3 attempts"
        assert error.command == ["test"]
        assert error.last_exit_code == 1
        assert error.attempts == 3
        assert error.total_delay == 3.0
        assert error.last_output == "error output"

    def test_error_raise_and_catch(self) -> None:
        """Test error can be raised and caught."""
        with pytest.raises(CommandRetryError) as exc_info:
            raise CommandRetryError(
                command=["false"],
                last_exit_code=1,
                attempts=1,
                total_delay=0.0,
            )

        assert exc_info.value.attempts == 1
