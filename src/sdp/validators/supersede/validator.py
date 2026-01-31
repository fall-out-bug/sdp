"""Supersede validator - prevents orphan superseded workstreams."""

from pathlib import Path

from sdp.validators.supersede.chain import find_orphans as find_orphans_impl
from sdp.validators.supersede.chain import trace_chain as trace_chain_impl
from sdp.validators.supersede.models import SupersedeChain, SupersedeResult, ValidationReport
from sdp.validators.supersede.parser import find_ws_file, parse_frontmatter, update_frontmatter


class SupersedeValidator:
    """Validate supersede relationships."""

    def __init__(self, ws_dir: Path = Path("docs/workstreams")):
        """Initialize validator.

        Args:
            ws_dir: Base directory for workstreams
        """
        self.ws_dir = ws_dir

    def supersede(self, old_ws: str, new_ws: str) -> SupersedeResult:
        """Mark old_ws as superseded by new_ws.

        Validates:
        - new_ws exists
        - No cycle created
        - old_ws not already superseded

        Args:
            old_ws: Workstream to supersede
            new_ws: Replacement workstream

        Returns:
            SupersedeResult with success/failure
        """
        # Find old WS file
        old_path = self._find_ws_file(old_ws)
        if not old_path:
            return SupersedeResult(
                success=False, old_ws=old_ws, new_ws=new_ws, error=f"Workstream not found: {old_ws}"
            )

        # Find new WS file
        new_path = self._find_ws_file(new_ws)
        if not new_path:
            return SupersedeResult(
                success=False,
                old_ws=old_ws,
                new_ws=new_ws,
                error=f"Replacement not found: {new_ws}",
            )

        # Check if already superseded
        old_data = self._parse_frontmatter(old_path)
        if old_data.get("status") == "superseded":
            existing = old_data.get("superseded_by", "unknown")
            return SupersedeResult(
                success=False,
                old_ws=old_ws,
                new_ws=new_ws,
                error=f"Already superseded by {existing}",
            )

        # Check for cycle
        chain = self.trace_chain(new_ws)
        if chain.has_cycle or old_ws in chain.replacements:
            return SupersedeResult(
                success=False,
                old_ws=old_ws,
                new_ws=new_ws,
                error="Would create supersede cycle",
            )

        # Update old WS frontmatter
        self._update_frontmatter(old_path, {"status": "superseded", "superseded_by": new_ws})

        return SupersedeResult(success=True, old_ws=old_ws, new_ws=new_ws, error=None)

    def find_orphans(self, ws_dir: Path | None = None) -> list[str]:
        """Find superseded WS without valid replacement.

        Args:
            ws_dir: Directory to search (defaults to self.ws_dir)

        Returns:
            List of orphaned WS IDs
        """
        if ws_dir is None:
            ws_dir = self.ws_dir
        return find_orphans_impl(ws_dir)

    def trace_chain(self, ws_id: str) -> SupersedeChain:
        """Trace supersede chain to final WS.

        Args:
            ws_id: Starting workstream ID

        Returns:
            SupersedeChain with full path
        """
        return trace_chain_impl(ws_id, self.ws_dir)

    def validate_all(self, ws_dir: Path | None = None) -> ValidationReport:
        """Validate all supersede relationships.

        Args:
            ws_dir: Directory to search (defaults to self.ws_dir)

        Returns:
            ValidationReport with summary
        """
        if ws_dir is None:
            ws_dir = self.ws_dir

        total_superseded = 0
        orphans: list[str] = []
        cycles: list[SupersedeChain] = []
        valid_chains: list[SupersedeChain] = []

        for ws_file in ws_dir.rglob("*.md"):
            data = self._parse_frontmatter(ws_file)

            if data.get("status") == "superseded":
                total_superseded += 1
                ws_id = data.get("ws_id", ws_file.stem)

                chain = self.trace_chain(ws_id)

                if chain.has_cycle:
                    cycles.append(chain)
                elif chain.final_ws is None:
                    orphans.append(ws_id)
                else:
                    valid_chains.append(chain)

        return ValidationReport(
            total_superseded=total_superseded,
            orphans=orphans,
            cycles=cycles,
            valid_chains=valid_chains,
        )

    def _find_ws_file(self, ws_id: str) -> Path | None:
        """Find WS file by ID (backward compatibility).

        Args:
            ws_id: Workstream ID

        Returns:
            Path to WS file or None
        """
        return find_ws_file(ws_id, self.ws_dir)

    def _parse_frontmatter(self, ws_path: Path) -> dict[str, str]:
        """Parse frontmatter from WS file (backward compatibility).

        Args:
            ws_path: Path to WS file

        Returns:
            Dict of frontmatter fields
        """
        return parse_frontmatter(ws_path)

    def _update_frontmatter(self, ws_path: Path, updates: dict[str, str]) -> None:
        """Update frontmatter fields (backward compatibility).

        Args:
            ws_path: Path to WS file
            updates: Dict of fields to update
        """
        update_frontmatter(ws_path, updates)
