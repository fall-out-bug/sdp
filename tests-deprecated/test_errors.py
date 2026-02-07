"""Tests for SDP error framework."""

import json

import pytest
from sdp.errors import (
    ArtifactValidationError,
    BeadsNotFoundError,
    BuildValidationError,
    ConfigurationError,
    CoverageTooLowError,
    DependencyNotFoundError,
    ErrorCategory,
    HookExecutionError,
    QualityGateViolationError,
    SDPError,
    TestFailureError,
    WorkstreamValidationError,
    format_error_for_json,
    format_error_for_terminal,
)


class TestSDPError:
    """Test base SDPError class."""

    def test_create_minimal_error(self) -> None:
        """Test creating error with required fields only."""
        error = SDPError(
            category=ErrorCategory.VALIDATION,
            message="Test error",
            remediation="Fix it",
        )

        assert error.category == ErrorCategory.VALIDATION
        assert error.message == "Test error"
        assert error.remediation == "Fix it"
        assert error.docs_url is None
        assert error.context is None

    def test_create_full_error(self) -> None:
        """Test creating error with all fields."""
        error = SDPError(
            category=ErrorCategory.BUILD,
            message="Build failed",
            remediation="Check logs",
            docs_url="https://docs.sdp.dev/build",
            context={"step": "compile", "exit_code": 1},
        )

        assert error.category == ErrorCategory.BUILD
        assert error.docs_url == "https://docs.sdp.dev/build"
        assert error.context == {"step": "compile", "exit_code": 1}

    def test_str_formatting_minimal(self) -> None:
        """Test string representation without optional fields."""
        error = SDPError(
            category=ErrorCategory.TEST,
            message="Test failed",
            remediation="Fix test",
        )

        output = str(error)
        assert "❌ TEST Error" in output
        assert "Test failed" in output
        assert "💡 Remediation:" in output
        assert "Fix test" in output
        assert "📚 Documentation:" not in output
        assert "🔍 Context:" not in output

    def test_str_formatting_with_docs(self) -> None:
        """Test string representation with docs URL."""
        error = SDPError(
            category=ErrorCategory.CONFIGURATION,
            message="Config error",
            remediation="Fix config",
            docs_url="https://docs.sdp.dev/config",
        )

        output = str(error)
        assert "📚 Documentation:" in output
        assert "https://docs.sdp.dev/config" in output

    def test_str_formatting_with_context(self) -> None:
        """Test string representation with context."""
        error = SDPError(
            category=ErrorCategory.BUILD,
            message="Build error",
            remediation="Check build",
            context={"file": "test.py", "line": 42},
        )

        output = str(error)
        assert "🔍 Context:" in output
        assert "file: test.py" in output
        assert "line: 42" in output

    def test_error_is_exception(self) -> None:
        """Test SDPError inherits from Exception."""
        error = SDPError(
            category=ErrorCategory.VALIDATION,
            message="Test",
            remediation="Fix",
        )

        assert isinstance(error, Exception)
        assert isinstance(error, SDPError)

    def test_raise_and_catch_sdp_error(self) -> None:
        """Test raising and catching SDPError."""
        with pytest.raises(SDPError) as exc_info:
            raise SDPError(
                category=ErrorCategory.TEST,
                message="Test error",
                remediation="Fix test",
            )

        # Check message is in the string representation
        assert "Test error" in str(exc_info.value)


class TestBeadsNotFoundError:
    """Test BeadsNotFoundError."""

    def test_basic_error(self) -> None:
        """Test basic BeadsNotFoundError."""
        error = BeadsNotFoundError(task_id="TASK-001")

        assert error.category == ErrorCategory.BEADS
        assert "TASK-001" in error.message
        assert error.context["task_id"] == "TASK-001"

    def test_error_with_search_paths(self) -> None:
        """Test BeadsNotFoundError with search paths."""
        paths = ["/path1", "/path2"]
        error = BeadsNotFoundError(task_id="TASK-002", search_paths=paths)

        assert error.context["search_paths"] == paths

    def test_remediation_message(self) -> None:
        """Test remediation includes useful steps."""
        error = BeadsNotFoundError(task_id="TASK-003")
        output = str(error)

        assert "beads status" in output
        assert "beads tasks" in output
        assert "BEADS_HOME" in output


class TestCoverageTooLowError:
    """Test CoverageTooLowError."""

    def test_basic_error(self) -> None:
        """Test basic CoverageTooLowError."""
        error = CoverageTooLowError(
            coverage_pct=65.5,
            required_pct=80.0,
            module="sdp.core",
        )

        assert error.category == ErrorCategory.COVERAGE
        assert "65.5%" in error.message
        assert "80.0%" in error.message
        assert error.context["actual_coverage"] == "65.5%"
        assert error.context["required_coverage"] == "80.0%"

    def test_error_with_missing_files(self) -> None:
        """Test CoverageTooLowError with missing files."""
        missing = ["src/sdp/core/parser.py", "src/sdp/core/validator.py"]
        error = CoverageTooLowError(
            coverage_pct=50.0,
            required_pct=80.0,
            module="sdp.core",
            missing_files=missing,
        )

        assert error.context["missing_files"] == missing

    def test_remediation_includes_command(self) -> None:
        """Test remediation includes pytest command."""
        error = CoverageTooLowError(
            coverage_pct=70.0,
            required_pct=80.0,
            module="sdp.validators",
        )
        output = str(error)

        assert "pytest --cov=sdp.validators" in output


class TestQualityGateViolationError:
    """Test QualityGateViolationError."""

    def test_basic_error(self) -> None:
        """Test basic QualityGateViolationError."""
        violations = ["File too large", "Missing tests"]
        error = QualityGateViolationError(
            gate_name="file_size",
            violations=violations,
        )

        assert error.category == ErrorCategory.VALIDATION
        assert "file_size" in error.message
        assert error.context["violations"] == violations
        assert error.context["severity"] == "error"

    def test_warning_severity(self) -> None:
        """Test QualityGateViolationError with warning severity."""
        error = QualityGateViolationError(
            gate_name="complexity",
            violations=["High CC"],
            severity="warning",
        )

        assert error.context["severity"] == "warning"

    def test_critical_severity(self) -> None:
        """Test QualityGateViolationError with critical severity."""
        error = QualityGateViolationError(
            gate_name="security",
            violations=["SQL injection risk"],
            severity="critical",
        )

        assert error.context["severity"] == "critical"


class TestWorkstreamValidationError:
    """Test WorkstreamValidationError."""

    def test_basic_error(self) -> None:
        """Test basic WorkstreamValidationError."""
        errors = ["Missing Goal section", "No Acceptance Criteria"]
        error = WorkstreamValidationError(
            ws_id="WS-001-01",
            errors=errors,
        )

        assert error.category == ErrorCategory.VALIDATION
        assert "WS-001-01" in error.message
        assert error.context["ws_id"] == "WS-001-01"
        assert error.context["errors"] == errors

    def test_error_with_file_path(self) -> None:
        """Test WorkstreamValidationError with file path."""
        error = WorkstreamValidationError(
            ws_id="WS-002-01",
            errors=["Invalid format"],
            file_path="docs/workstreams/backlog/WS-002-01.md",
        )

        assert error.context["file_path"] == "docs/workstreams/backlog/WS-002-01.md"


class TestConfigurationError:
    """Test ConfigurationError."""

    def test_basic_error(self) -> None:
        """Test basic ConfigurationError."""
        errors = ["Invalid TOML", "Missing key"]
        error = ConfigurationError(
            config_file="quality-gate.toml",
            errors=errors,
        )

        assert error.category == ErrorCategory.CONFIGURATION
        assert "quality-gate.toml" in error.message
        assert error.context["errors"] == errors

    def test_error_with_missing_keys(self) -> None:
        """Test ConfigurationError with missing keys."""
        missing = ["max_file_size", "min_coverage"]
        error = ConfigurationError(
            config_file="sdp.config",
            errors=["Schema validation failed"],
            missing_keys=missing,
        )

        assert error.context["missing_keys"] == missing


class TestDependencyNotFoundError:
    """Test DependencyNotFoundError."""

    def test_basic_error(self) -> None:
        """Test basic DependencyNotFoundError."""
        error = DependencyNotFoundError(
            dependency="WS-001-01",
        )

        assert error.category == ErrorCategory.DEPENDENCY
        assert "WS-001-01" in error.message
        assert error.context["dependency"] == "WS-001-01"

    def test_error_with_ws_id(self) -> None:
        """Test DependencyNotFoundError for specific workstream."""
        error = DependencyNotFoundError(
            dependency="WS-001-01",
            ws_id="WS-001-02",
        )

        assert "WS-001-02" in error.message
        assert error.context["ws_id"] == "WS-001-02"

    def test_error_with_available_ws(self) -> None:
        """Test DependencyNotFoundError with available workstreams."""
        available = ["WS-001-01", "WS-001-03"]
        error = DependencyNotFoundError(
            dependency="WS-001-02",
            available_ws=available,
        )

        assert error.context["available_workstreams"] == available


class TestHookExecutionError:
    """Test HookExecutionError."""

    def test_basic_error(self) -> None:
        """Test basic HookExecutionError."""
        error = HookExecutionError(
            hook_name="pre-commit",
            stage="pre-commit",
            output="Command failed",
            exit_code=1,
        )

        assert error.category == ErrorCategory.HOOK
        assert "pre-commit" in error.message
        assert error.context["exit_code"] == 1
        assert error.context["hook"] == "pre-commit"

    def test_output_truncation(self) -> None:
        """Test long output is truncated."""
        long_output = "x" * 1000
        error = HookExecutionError(
            hook_name="post-build",
            stage="post-build",
            output=long_output,
            exit_code=2,
        )

        assert len(error.context["output"]) <= 500


class TestFailedTestsError:
    """Test TestFailureError."""

    def test_basic_error(self) -> None:
        """Test basic TestFailureError."""
        failed = ["test_foo", "test_bar"]
        error = TestFailureError(
            test_command="pytest",
            failed_tests=failed,
            total_tests=10,
            passed_tests=8,
        )

        assert error.category == ErrorCategory.TEST
        assert "8/10" in error.message
        assert error.context["passed"] == 8
        assert error.context["total"] == 10
        assert error.context["failed"] == 2

    def test_failed_tests_limit(self) -> None:
        """Test failed tests are limited to first 10."""
        failed = [f"test_{i}" for i in range(20)]
        error = TestFailureError(
            test_command="pytest",
            failed_tests=failed,
            total_tests=20,
            passed_tests=0,
        )

        assert len(error.context["failed_tests"]) == 10


class TestBuildValidationError:
    """Test BuildValidationError."""

    def test_basic_error(self) -> None:
        """Test basic BuildValidationError."""
        error = BuildValidationError(
            ws_id="WS-001-01",
            stage="pre-build",
            check_name="Goal section",
            details="Goal section is missing",
        )

        assert error.category == ErrorCategory.BUILD
        assert "Goal section" in error.message
        assert error.context["ws_id"] == "WS-001-01"
        assert error.context["stage"] == "pre-build"


class TestArtifactValidationError:
    """Test ArtifactValidationError."""

    def test_basic_error(self) -> None:
        """Test basic ArtifactValidationError."""
        errors = ["Too long", "Missing type hints"]
        error = ArtifactValidationError(
            artifact_type="code",
            artifact_path="src/sdp/module.py",
            errors=errors,
        )

        assert error.category == ErrorCategory.ARTIFACT
        assert "code" in error.message
        assert error.context["path"] == "src/sdp/module.py"
        assert error.context["errors"] == errors


class TestFormatErrorForTerminal:
    """Test format_error_for_terminal utility."""

    def test_format_sdp_error(self) -> None:
        """Test formatting SDPError for terminal."""
        error = SDPError(
            category=ErrorCategory.VALIDATION,
            message="Test error",
            remediation="Fix it",
        )

        output = format_error_for_terminal(error)
        assert "❌ VALIDATION Error" in output
        assert "Test error" in output
        assert "Fix it" in output

    def test_format_standard_exception(self) -> None:
        """Test formatting standard exception."""
        error = ValueError("Invalid value")

        output = format_error_for_terminal(error)
        assert "❌ Error: ValueError" in output
        assert "Invalid value" in output


class TestFormatErrorForJson:
    """Test format_error_for_json utility."""

    def test_format_sdp_error_json(self) -> None:
        """Test formatting SDPError as JSON."""
        error = SDPError(
            category=ErrorCategory.BUILD,
            message="Build failed",
            remediation="Check logs",
            docs_url="https://docs.sdp.dev/build",
            context={"step": "compile"},
        )

        data = format_error_for_json(error)
        assert data["category"] == "build"
        assert data["message"] == "Build failed"
        assert data["remediation"] == "Check logs"
        assert data["docs_url"] == "https://docs.sdp.dev/build"
        assert data["context"] == {"step": "compile"}
        assert data["type"] == "SDPError"

    def test_json_is_serializable(self) -> None:
        """Test JSON output is serializable."""
        error = BeadsNotFoundError(task_id="TASK-001")
        data = format_error_for_json(error)

        # Should not raise
        json.dumps(data)

    def test_format_standard_exception_json(self) -> None:
        """Test formatting standard exception as JSON."""
        error = ValueError("Invalid value")
        data = format_error_for_json(error)

        assert data["type"] == "ValueError"
        assert data["message"] == "Invalid value"
        assert data["category"] == "unknown"


class TestErrorCategoryEnum:
    """Test ErrorCategory enum."""

    def test_all_categories_exist(self) -> None:
        """Test all expected categories are defined."""
        expected = {
            "validation",
            "build",
            "test",
            "configuration",
            "dependency",
            "hook",
            "artifact",
            "beads",
            "coverage",
        }

        actual = {category.value for category in ErrorCategory}
        assert actual == expected
