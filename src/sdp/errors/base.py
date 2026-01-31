"""Base error classes for SDP.

Defines ErrorCategory enum, SDPError base exception, and formatting functions.
"""

from dataclasses import dataclass
from enum import Enum
from typing import Any


class ErrorCategory(str, Enum):
    """Category of error for filtering and routing.

    Categories:
    - VALIDATION: Artifact validation, schema violations
    - BUILD: Compilation, quality gate violations
    - TEST: Test failures
    - CONFIGURATION: Configuration issues
    - DEPENDENCY: Missing dependencies
    - HOOK: Git hook failures
    - ARTIFACT: Artifact issues
    - BEADS: Beads integration errors
    - COVERAGE: Coverage violations
    """

    VALIDATION = "validation"
    BUILD = "build"
    TEST = "test"
    CONFIGURATION = "configuration"
    DEPENDENCY = "dependency"
    HOOK = "hook"
    ARTIFACT = "artifact"
    BEADS = "beads"
    COVERAGE = "coverage"


@dataclass
class SDPError(Exception):
    """Base exception for all SDP errors.

    Provides structured error information with:
    - Category for filtering/routing
    - Message: Human-readable error description
    - Remediation steps for user guidance
    - Docs URL: Link to documentation (optional)
    - Context for debugging

    Example:
        raise SDPError(
            category=ErrorCategory.VALIDATION,
            message="Workstream file not found",
            remediation="Check the WS-ID and ensure file exists",
            docs_url="https://docs.sdp.dev/troubleshooting",
            context={"ws_id": "WS-001-01"},
        )
    """

    category: ErrorCategory
    message: str
    remediation: str
    docs_url: str | None = None
    context: dict[str, Any] | None = None

    def __str__(self) -> str:
        """Format error for display."""
        lines = [
            f"❌ {self.category.value.upper()} Error",
            f"   {self.message}",
            "",
            "   💡 Remediation:",
            f"   {self.remediation}",
        ]

        if self.docs_url:
            lines.extend([
                "",
                "   📚 Documentation:",
                f"   {self.docs_url}",
            ])

        if self.context:
            lines.extend([
                "",
                "   🔍 Context:",
            ])
            for key, value in self.context.items():
                lines.append(f"   {key}: {value}")

        return "\n".join(lines)


def format_error_for_terminal(error: Exception) -> str:
    """Format any exception for terminal display.

    Args:
        error: Exception to format

    Returns:
        Formatted error message

    Example:
        >>> try:
        ...     raise SDPError(
        ...         category=ErrorCategory.VALIDATION,
        ...         message="Test error",
        ...         remediation="Fix it",
        ...     )
        ... except SDPError as e:
        ...     print(format_error_for_terminal(e))
    """
    if isinstance(error, SDPError):
        return str(error)

    # Standard Python exception
    return f"❌ Error: {type(error).__name__}\n   {error}"


def format_error_for_json(error: Exception) -> dict[str, Any]:
    """Format exception as JSON-serializable dict.

    Args:
        error: Exception to format

    Returns:
        Dictionary with error details

    Example:
        >>> try:
        ...     raise SDPError(
        ...         category=ErrorCategory.VALIDATION,
        ...         message="Test error",
        ...         remediation="Fix it",
        ...     )
        ... except SDPError as e:
        ...     import json
        ...     print(json.dumps(format_error_for_json(e)))
    """
    if isinstance(error, SDPError):
        return {
            "type": type(error).__name__,
            "category": error.category.value,
            "message": error.message,
            "remediation": error.remediation,
            "docs_url": error.docs_url,
            "context": error.context,
        }

    return {
        "type": type(error).__name__,
        "category": "unknown",
        "message": str(error),
    }
