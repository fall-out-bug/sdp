"""Supersede validator - prevents orphan superseded workstreams."""

from dataclasses import dataclass
from pathlib import Path


@dataclass
class SupersedeChain:
    """Supersede chain representation."""

    original_ws: str
    replacements: list[str]
    has_cycle: bool
    final_ws: str | None  # None if cycle


@dataclass
class SupersedeResult:
    """Result of supersede operation."""

    success: bool
    old_ws: str
    new_ws: str
    error: str | None


@dataclass
class ValidationReport:
    """Report of all supersede validations."""

    total_superseded: int
    orphans: list[str]
    cycles: list[SupersedeChain]
    valid_chains: list[SupersedeChain]


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

        orphans: list[str] = []

        for ws_file in ws_dir.rglob("*.md"):
            data = self._parse_frontmatter(ws_file)

            if data.get("status") == "superseded":
                superseded_by = data.get("superseded_by")

                if not superseded_by:
                    # No replacement specified
                    orphans.append(data.get("ws_id", ws_file.stem))
                else:
                    # Check if replacement exists
                    replacement_path = self._find_ws_file(superseded_by)
                    if not replacement_path:
                        orphans.append(data.get("ws_id", ws_file.stem))

        return orphans

    def trace_chain(self, ws_id: str) -> SupersedeChain:
        """Trace supersede chain to final WS.

        Args:
            ws_id: Starting workstream ID

        Returns:
            SupersedeChain with full path
        """
        visited: set[str] = set()
        chain: list[str] = [ws_id]
        current = ws_id

        while True:
            if current in visited:
                # Cycle detected
                return SupersedeChain(
                    original_ws=ws_id, replacements=chain, has_cycle=True, final_ws=None
                )

            visited.add(current)

            # Find current WS
            ws_path = self._find_ws_file(current)
            if not ws_path:
                # Dead end
                return SupersedeChain(
                    original_ws=ws_id, replacements=chain, has_cycle=False, final_ws=None
                )

            # Check if superseded
            data = self._parse_frontmatter(ws_path)
            if data.get("status") != "superseded":
                # Found final WS
                return SupersedeChain(
                    original_ws=ws_id, replacements=chain, has_cycle=False, final_ws=current
                )

            # Continue to replacement
            next_ws = data.get("superseded_by")
            if not next_ws:
                # Orphan
                return SupersedeChain(
                    original_ws=ws_id, replacements=chain, has_cycle=False, final_ws=None
                )

            chain.append(next_ws)
            current = next_ws

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
        """Find WS file by ID.

        Args:
            ws_id: Workstream ID

        Returns:
            Path to WS file or None
        """
        search_dirs = [
            self.ws_dir / "backlog",
            self.ws_dir / "in_progress",
            self.ws_dir / "completed",
        ]

        for search_dir in search_dirs:
            if not search_dir.exists():
                continue

            for ws_file in search_dir.glob(f"{ws_id}*.md"):
                return ws_file

        return None

    def _parse_frontmatter(self, ws_path: Path) -> dict:
        """Parse frontmatter from WS file.

        Args:
            ws_path: Path to WS file

        Returns:
            Dict of frontmatter fields
        """
        content = ws_path.read_text(encoding="utf-8")
        data: dict = {}

        in_frontmatter = False
        for line in content.splitlines():
            if line.strip() == "---":
                if not in_frontmatter:
                    in_frontmatter = True
                else:
                    break
                continue

            if in_frontmatter and ":" in line:
                key, value = line.split(":", 1)
                data[key.strip()] = value.strip()

        return data

    def _update_frontmatter(self, ws_path: Path, updates: dict) -> None:
        """Update frontmatter fields.

        Args:
            ws_path: Path to WS file
            updates: Dict of fields to update
        """
        content = ws_path.read_text(encoding="utf-8")
        lines = content.splitlines()

        # Find frontmatter bounds
        fm_start = -1
        fm_end = -1
        for i, line in enumerate(lines):
            if line.strip() == "---":
                if fm_start == -1:
                    fm_start = i
                else:
                    fm_end = i
                    break

        if fm_start == -1 or fm_end == -1:
            raise ValueError("Frontmatter not found")

        # Update frontmatter
        new_fm_lines: list[str] = []
        updated_keys: set[str] = set()

        for i in range(fm_start + 1, fm_end):
            line = lines[i]
            if ":" in line:
                key, value = line.split(":", 1)
                key = key.strip()

                if key in updates:
                    new_fm_lines.append(f"{key}: {updates[key]}")
                    updated_keys.add(key)
                else:
                    new_fm_lines.append(line)
            else:
                new_fm_lines.append(line)

        # Add new fields
        for key, value in updates.items():
            if key not in updated_keys:
                new_fm_lines.append(f"{key}: {value}")

        # Reconstruct file
        new_lines = (
            lines[:fm_start] + ["---"] + new_fm_lines + ["---"] + lines[fm_end + 1 :]
        )

        ws_path.write_text("\n".join(new_lines) + "\n", encoding="utf-8")
