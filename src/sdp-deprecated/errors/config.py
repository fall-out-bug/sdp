"""Configuration and setup errors.

Covers: Configuration issues, missing dependencies, hook failures.
"""

from .base import ErrorCategory, SDPError


class ConfigurationError(SDPError):
    """SDP configuration error."""

    def __init__(
        self,
        config_file: str,
        errors: list[str],
        missing_keys: list[str] | None = None,
    ) -> None:
        """Initialize ConfigurationError.

        Args:
            config_file: Configuration file path
            errors: List of configuration errors
            missing_keys: Missing required configuration keys
        """
        super().__init__(
            category=ErrorCategory.CONFIGURATION,
            message=f"Configuration error in '{config_file}': {len(errors)} issue(s)",
            remediation=(
                "1. Review configuration file format\n"
                "2. Add missing required keys\n"
                "3. Check syntax (TOML/YAML/JSON)\n"
                "4. See docs for config schema"
            ),
            docs_url="https://docs.sdp.dev/configuration",
            context={
                "config_file": config_file,
                "errors": errors,
                "missing_keys": missing_keys or [],
            },
        )


class DependencyNotFoundError(SDPError):
    """Required dependency not found."""

    def __init__(
        self,
        dependency: str,
        ws_id: str | None = None,
        available_ws: list[str] | None = None,
    ) -> None:
        """Initialize DependencyNotFoundError.

        Args:
            dependency: Missing dependency identifier
            ws_id: Workstream that depends on this
            available_ws: List of available workstreams
        """
        ws_context = f" for workstream '{ws_id}'" if ws_id else ""
        super().__init__(
            category=ErrorCategory.DEPENDENCY,
            message=f"Dependency '{dependency}' not found{ws_context}",
            remediation=(
                f"1. Complete dependency workstream first: {dependency}\n"
                f"2. Check INDEX.md for workstream status\n"
                f"3. Verify dependency ID is correct\n"
                f"4. Update workstream frontmatter if needed"
            ),
            docs_url="https://docs.sdp.dev/workstreams#dependencies",
            context={
                "dependency": dependency,
                "ws_id": ws_id,
                "available_workstreams": available_ws or [],
            },
        )


class HookExecutionError(SDPError):
    """Git hook or build hook execution failed."""

    def __init__(
        self,
        hook_name: str,
        stage: str,
        output: str,
        exit_code: int,
    ) -> None:
        """Initialize HookExecutionError.

        Args:
            hook_name: Name of the hook that failed
            stage: pre-commit, post-build, etc.
            output: Hook output/stderr
            exit_code: Hook exit code
        """
        super().__init__(
            category=ErrorCategory.HOOK,
            message=f"Hook '{hook_name}' failed during {stage} (exit code: {exit_code})",
            remediation=(
                "1. Review hook output above for specific errors\n"
                "2. Fix the issue that caused hook to fail\n"
                "3. Test hook manually: hooks/{hook_name}.sh\n"
                "4. Bypass with SKIP_CHECK=1 if needed (not recommended)"
            ),
            docs_url="https://docs.sdp.dev/hooks#troubleshooting",
            context={
                "hook": hook_name,
                "stage": stage,
                "exit_code": exit_code,
                "output": output[:500],  # Truncate long output
            },
        )
