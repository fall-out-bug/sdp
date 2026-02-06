"""Skip flag parser for @oneshot workflow."""

import logging
from typing import Optional

from sdp.unified.gates.models import GateType

logger = logging.getLogger(__name__)


class SkipFlagParser:
    """Parser for --skip-* command line flags."""

    # Flag mappings
    FLAG_MAP = {
        "--skip-requirements": GateType.REQUIREMENTS,
        "--skip-architecture": GateType.ARCHITECTURE,
        "--skip-uat": GateType.UAT,
    }

    def __init__(self, args: list[str]) -> None:
        """Initialize skip flag parser.

        Args:
            args: Command line arguments list
        """
        self.args = args
        self._skips: Optional[set[GateType]] = None

    def parse(self) -> set[GateType]:
        """Parse skip flags from arguments.

        Returns:
            Set of gate types to skip
        """
        if self._skips is not None:
            return self._skips

        skips = set()
        for arg in self.args:
            if arg in self.FLAG_MAP:
                gate_type = self.FLAG_MAP[arg]
                skips.add(gate_type)
                logger.debug(f"Skip flag detected: {arg}")

        self._skips = skips
        return skips

    def is_skip_required(self, gate_type: GateType) -> bool:
        """Check if a gate should be skipped.

        Args:
            gate_type: Type of gate to check

        Returns:
            True if gate should be skipped, False otherwise
        """
        skips = self.parse()
        return gate_type in skips
