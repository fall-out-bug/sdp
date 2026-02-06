"""Feature-related error classes."""

from sdp.errors import ErrorCategory, SDPError


class CircularDependencyError(SDPError):
    """Circular dependency detected in workstream graph."""

    def __init__(self, ws_id: str, cycle: list[str]) -> None:
        formatted_cycle = " → ".join(cycle + [cycle[0]]) if cycle else ws_id
        super().__init__(
            category=ErrorCategory.DEPENDENCY,
            message=f"Circular dependency detected: {formatted_cycle}",
            remediation=(
                f"1. Break the cycle by removing one dependency:\n"
                f"   - {ws_id} depends on: {' → '.join(cycle)}\n"
                "2. Reorder workstreams to avoid circular reference\n"
                "3. Or split into smaller independent features\n"
                "4. See docs/dependency-management.md for strategies"
            ),
            docs_url="https://sdp.dev/docs/dependencies#circular",
            context={"ws_id": ws_id, "cycle": cycle},
        )


class MissingDependencyError(SDPError):
    """Required workstream dependency not found."""

    def __init__(
        self,
        ws_id: str,
        missing_dep: str,
        available_workstreams: list[str],
    ) -> None:
        super().__init__(
            category=ErrorCategory.DEPENDENCY,
            message=f"Workstream {ws_id} depends on {missing_dep} which doesn't exist",
            remediation=(
                f"1. Create missing workstream first: {missing_dep}\n"
                "2. Or remove dependency if not actually needed\n"
                f"3. Available workstreams: {', '.join(available_workstreams[:5])}\n"
                "4. See docs/workflows/dependency-management.md"
            ),
            docs_url="https://docs.sdp.dev/workflows#dependencies",
            context={
                "ws_id": ws_id,
                "missing_dep": missing_dep,
                "available_ws": available_workstreams,
            },
        )
