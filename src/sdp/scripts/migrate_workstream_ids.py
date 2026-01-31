"""Workstream migration orchestrator."""

from pathlib import Path
from typing import Dict, List, Tuple

from sdp.scripts.migrate_models import WorkstreamFile, WorkstreamMigrationError


class WorkstreamMigrator:
    """Orchestrates workstream migration."""

    def __init__(
        self,
        root_path: Path,
        project_id: str = "00",
        dry_run: bool = False,
    ) -> None:
        """Initialize migrator.

        Args:
            root_path: Root directory of project
            project_id: Project ID (default: "00" for SDP)
            dry_run: Enable dry-run mode
        """
        self.root_path = root_path
        self.project_id = project_id
        self.dry_run = dry_run
        self.ws_dir = root_path / "docs" / "workstreams"
        self.results: List[Tuple[bool, str]] = []

    def migrate(self) -> Dict[str, int]:
        """Execute migration.

        Returns:
            Dictionary with migration statistics
        """
        if not self.ws_dir.exists():
            raise WorkstreamMigrationError(
                f"Workstreams directory not found: {self.ws_dir}"
            )

        print(f"{'=' * 70}")
        print("Workstream ID Migration")
        print(f"{'=' * 70}")
        print(f"Project ID: {self.project_id}")
        print(f"Path: {self.ws_dir}")
        print(f"Mode: {'DRY RUN' if self.dry_run else 'LIVE'}")
        print(f"{'=' * 70}\n")

        # Find all workstream files
        ws_files = self._find_workstream_files()

        if not ws_files:
            print("⚠️  No workstream files found to migrate")
            return {"total": 0, "migrated": 0, "skipped": 0, "failed": 0}

        print(f"Found {len(ws_files)} workstream files\n")

        # Migrate each file
        for ws_file in ws_files:
            try:
                ws = WorkstreamFile(ws_file, self.project_id)
                success, message = ws.migrate(self.dry_run)
                self.results.append((success, message))

                status = "✓" if success else "✗"
                print(f"{status} {message}")

            except Exception as e:
                self.results.append((False, f"Error: {ws_file.name}: {e}"))
                print(f"✗ Error processing {ws_file.name}: {e}")

        # Print summary
        return self._print_summary()

    def _find_workstream_files(self) -> List[Path]:
        """Find all workstream markdown files."""
        files: List[Path] = []

        # Find old format (WS-FFF-SS)
        files.extend(self.ws_dir.rglob("WS-*.md"))

        # Find files with old format ws_id in frontmatter
        for md_file in self.ws_dir.rglob("*.md"):
            content = md_file.read_text(encoding="utf-8")
            if "ws_id: WS-" in content:
                files.append(md_file)

        # Remove duplicates and sort
        return sorted(set(files))

    def _print_summary(self) -> Dict[str, int]:
        """Print migration summary."""
        print(f"\n{'=' * 70}")
        print("Migration Summary")
        print(f"{'=' * 70}")

        stats = {
            "total": len(self.results),
            "migrated": sum(1 for s, _ in self.results if s and "Migrated" in _),
            "skipped": sum(1 for s, _ in self.results if s and "Already" in _),
            "failed": sum(1 for s, _ in self.results if not s),
        }

        print(f"Total files: {stats['total']}")
        print(f"✓ Migrated: {stats['migrated']}")
        print(f"⊘ Skipped: {stats['skipped']}")
        print(f"✗ Failed: {stats['failed']}")

        if stats['failed'] > 0:
            print("\nFailed files:")
            for success, msg in self.results:
                if not success:
                    print(f"  - {msg}")

        # Verification commands
        print(f"\n{'=' * 70}")
        print("Verification")
        print(f"{'=' * 70}")
        print("\nRun these commands to verify migration:\n")

        print("  # Check for remaining old format")
        print(f"  grep -r 'ws_id: WS-' {self.ws_dir}")
        print("  # Should return empty\n")

        print("  # Verify new format")
        print(f"  grep -r 'project_id:' {self.ws_dir}")
        print("  # Should show all files with project_id\n")

        if not self.dry_run:
            print("  # Count files by format")
            print(f"  find {self.ws_dir} -name 'WS-*.md' | wc -l  # Old format (should be 0)")
            print(f"  find {self.ws_dir} -name '{self.project_id}-*.md' | wc -l  # New format")

        return stats
