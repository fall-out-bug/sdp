"""Examples of using the SDP error framework."""

from sdp import errors


def example_basic_error() -> None:
    """Example: Create a basic SDP error."""
    error = errors.SDPError(
        category=errors.ErrorCategory.VALIDATION,
        message="Workstream file is missing required sections",
        remediation="Add Goal and Acceptance Criteria sections to the workstream file",
        docs_url="https://docs.sdp.dev/troubleshooting#ws-validation",
    )

    print("Basic Error:")
    print(str(error))


def example_coverage_error() -> None:
    """Example: Coverage too low error."""
    error = errors.CoverageTooLowError(
        coverage_pct=65.5,
        required_pct=80.0,
        module="sdp.core",
        missing_files=["src/sdp/core/parser.py", "src/sdp/core/validator.py"],
    )

    print("\nCoverage Error:")
    print(str(error))


def example_dependency_error() -> None:
    """Example: Dependency not found error."""
    error = errors.DependencyNotFoundError(
        dependency="WS-001-01",
        ws_id="WS-001-02",
        available_ws=["WS-001-01", "WS-001-03"],
    )

    print("\nDependency Error:")
    print(str(error))


def example_hook_error() -> None:
    """Example: Hook execution error."""
    error = errors.HookExecutionError(
        hook_name="pre-commit",
        stage="pre-commit",
        output="Time estimates found in WS files",
        exit_code=1,
    )

    print("\nHook Error:")
    print(str(error))


def example_quality_gate_error() -> None:
    """Example: Quality gate violation error."""
    error = errors.QualityGateViolationError(
        gate_name="file_size",
        violations=["module.py: 250 lines (max: 200)", "utils.py: 180 lines"],
        severity="warning",
    )

    print("\nQuality Gate Error:")
    print(str(error))


def example_json_formatting() -> None:
    """Example: Format error as JSON."""
    error = errors.BeadsNotFoundError(task_id="TASK-001")
    json_data = errors.format_error_for_json(error)

    print("\nJSON Format:")
    import json
    print(json.dumps(json_data, indent=2))


if __name__ == "__main__":
    print("SDP Error Framework Examples")
    print("=" * 50)

    example_basic_error()
    example_coverage_error()
    example_dependency_error()
    example_hook_error()
    example_quality_gate_error()
    example_json_formatting()

    print("\n" + "=" * 50)
    print("For more information, see: docs/troubleshooting.md")
