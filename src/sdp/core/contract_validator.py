"""Runtime contract validation for T2/T3 workstreams.

Ensures that T2/T3 builds do not modify Interface or Tests sections
of Contract-Driven workstreams.
"""

import re
from dataclasses import dataclass
from pathlib import Path
from typing import Optional

from sdp.errors import ErrorCategory, SDPError


class ContractViolationError(SDPError):
    """Raised when T2/T3 build violates contract."""

    def __init__(self, ws_id: str, tier: str, violation: str) -> None:
        """Initialize contract violation error."""
        super().__init__(
            category=ErrorCategory.BUILD,
            message=f"Contract violation in {ws_id} (tier {tier}): {violation}",
            remediation=(
                "1. Revert changes to Interface and Tests sections\n"
                "2. Only modify Implementation section in T2/T3 builds\n"
                "3. Run pre-build validation to catch contract changes\n"
                "4. See docs/reference/contract-driven-development.md"
            ),
            docs_url="https://docs.sdp.dev/contract-driven",
            context={"ws_id": ws_id, "tier": tier, "violation": violation},
        )


@dataclass
class ContractSnapshot:
    """Snapshot of contract sections (Interface + Tests).

    Args:
        interface_content: Content of Interface section
        tests_content: Content of Tests section
    """

    interface_content: str
    tests_content: str

    def equals(self, other: "ContractSnapshot") -> bool:
        """Check if two snapshots are identical.

        Args:
            other: Another snapshot to compare

        Returns:
            True if interface and tests are identical
        """
        return (
            self.interface_content == other.interface_content
            and self.tests_content == other.tests_content
        )


class ContractValidator:
    """Validate contract integrity for T2/T3 builds.

    Raises:
        ContractViolationError: If contract is modified during build
    """

    def snapshot_contract(self, ws_file: Path) -> Optional[ContractSnapshot]:
        """Extract Interface + Tests sections from WS file.

        Args:
            ws_file: Path to workstream markdown file

        Returns:
            ContractSnapshot if contract sections exist, None otherwise
        """
        if not ws_file.exists():
            return None

        content = ws_file.read_text(encoding="utf-8")

        # Find Interface section
        interface_match = re.search(
            r"#### Interface \(DO NOT MODIFY для T2/T3\)\n(.*?)(?=\n####|\n---|\Z)",
            content,
            re.DOTALL,
        )

        # Find Tests section
        tests_match = re.search(
            r"#### Tests \(DO NOT MODIFY для T2/T3\)\n(.*?)(?=\n####|\n---|\Z)",
            content,
            re.DOTALL,
        )

        if not interface_match and not tests_match:
            return None

        interface_content = interface_match.group(1).strip() if interface_match else ""
        tests_content = tests_match.group(1).strip() if tests_match else ""

        return ContractSnapshot(
            interface_content=interface_content,
            tests_content=tests_content,
        )

    def validate_contract_integrity(
        self,
        before: ContractSnapshot,
        after: ContractSnapshot,
        ws_id: str,
        tier: str,
    ) -> None:
        """Validate contract wasn't modified.

        Args:
            before: Snapshot before build
            after: Snapshot after build
            ws_id: Workstream ID
            tier: Capability tier

        Raises:
            ContractViolationError: If contract was modified
        """
        if not before.equals(after):
            # Determine what changed
            violations = []

            if before.interface_content != after.interface_content:
                violations.append("Interface section was modified")

            if before.tests_content != after.tests_content:
                violations.append("Tests section was modified")

            violation = ", ".join(violations)
            raise ContractViolationError(ws_id, tier, violation)
