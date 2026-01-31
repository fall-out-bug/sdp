#!/usr/bin/env python3
"""Migrate workstream IDs to PP-FFF-SS format with validation and dry-run support.

This script standardizes workstream IDs across the SDP protocol by:
1. Converting WS-FFF-SS → PP-FFF-SS format
2. Adding project_id field to frontmatter
3. Updating filenames to match new format
4. Validating all references and dependencies

Usage:
    python scripts/migrate_workstream_ids.py --dry-run
    python scripts/migrate_workstream_ids.py --project-id 00
    python scripts/migrate_workstream_ids.py --project-id 02 --path ../hw_checker

Examples:
    # Dry run for SDP (project 00)
    python scripts/migrate_workstream_ids.py --dry-run

    # Migrate SDP workstreams
    python scripts/migrate_workstream_ids.py --project-id 00

    # Migrate hw_checker workstreams
    python scripts/migrate_workstream_ids.py --project-id 02 --path ../hw_checker
"""

import argparse
import re
import sys
from pathlib import Path
from typing import Dict, List, Optional, Tuple

# Workstream ID patterns
OLD_PATTERN = re.compile(r"WS-(\d+)-(\d+)")  # WS-FFF-SS
NEW_PATTERN = re.compile(r"(\d+)-(\d+)-(\d+)")  # PP-FFF-SS


class WorkstreamMigrationError(Exception):
    """Base exception for migration errors."""


class WorkstreamFile:
    """Represents a workstream file with migration capabilities."""

    def __init__(self, path: Path, project_id: str = "00") -> None:
        """Initialize workstream file.

        Args:
            path: Path to workstream markdown file
            project_id: Project ID to use (default: "00" for SDP)

        Raises:
            WorkstreamMigrationError: If file cannot be read or parsed
        """
        self.path = path
        self.project_id = project_id
        self.content = path.read_text(encoding="utf-8")
        self.old_id: Optional[str] = None
        self.new_id: Optional[str] = None
        self._parse_ids()

    def _parse_ids(self) -> None:
        """Extract workstream ID from content and filename."""
        # Try frontmatter first
        match = re.search(r"ws_id:\s*(WS-[\d-]+|[\d-]+)", self.content)
        if match:
            potential_id = match.group(1)
            if potential_id.startswith("WS-"):
                # Old format: WS-FFF-SS
                self.old_id = potential_id
            else:
                # New format: PP-FFF-SS
                self.old_id = None
                self.new_id = potential_id
                return
        else:
            # Try filename
            filename_match = OLD_PATTERN.search(self.path.stem)
            if filename_match:
                self.old_id = f"WS-{filename_match.group(1)}-{filename_match.group(2)}"
            else:
                # Check if already in new format via filename
                filename_new_match = NEW_PATTERN.search(self.path.stem)
                if filename_new_match:
                    # Already in new format
                    self.old_id = None
                    self.new_id = self.path.stem
                    return

        # Calculate new ID for old format
        if self.old_id and self.old_id.startswith("WS-"):
            parts = self.old_id[3:].split("-")  # Remove "WS-"
            if len(parts) == 2:
                feature_num = parts[0].zfill(3)
                ws_num = parts[1].zfill(2)
                self.new_id = f"{self.project_id}-{feature_num}-{ws_num}"

    def needs_migration(self) -> bool:
        """Check if file needs migration."""
        return self.old_id is not None and self.new_id is not None

    def migrate(self, dry_run: bool = False) -> Tuple[bool, str]:
        """Migrate workstream file to new format.

        Args:
            dry_run: If True, don't make actual changes

        Returns:
            Tuple of (success, message)
        """
        if not self.needs_migration():
            return True, f"Already in new format: {self.path.name}"

        if not self.old_id or not self.new_id:
            return False, f"Cannot parse ID from: {self.path.name}"

        try:
            # Update content
            new_content = self._update_content()

            # Calculate new filename
            new_filename = self._generate_filename()
            new_path = self.path.parent / new_filename

            if dry_run:
                return (
                    True,
                    f"[DRY RUN] Would rename: {self.path.name} → {new_filename}",
                )

            # Write updated content
            self.path.write_text(new_content, encoding="utf-8")

            # Rename file
            if self.path.name != new_filename:
                self.path.rename(new_path)
                self.path = new_path

            return True, f"Migrated: {self.old_id} → {self.new_id}"

        except Exception as e:
            return False, f"Failed to migrate {self.path.name}: {e}"

    def _update_content(self) -> str:
        """Update content with new workstream ID format."""
        content = self.content

        # Update ws_id in frontmatter
        if self.old_id and self.new_id:
            content = content.replace(f"ws_id: {self.old_id}", f"ws_id: {self.new_id}")

            # Add project_id if not present
            if "project_id:" not in content:
                content = re.sub(
                    r"(ws_id:\s*[^\n]+\n)",
                    rf"\1project_id: {self.project_id}\n",
                    content,
                    count=1,
                )

            # Update title headers
            content = content.replace(f"## {self.old_id}:", f"## {self.new_id}:")
            content = content.replace(f"@{self.old_id}", f"@{self.new_id}")

            # Update dependencies
            old_dep_pattern = re.compile(r"depends_on:\s*\n((?:\s*-\s*[\w-]+\n)*)")
            for match in old_dep_pattern.finditer(content):
                deps_section = match.group(1)
                new_deps = []
                for dep_line in deps_section.split("\n"):
                    if dep_line.strip():
                        dep_match = OLD_PATTERN.search(dep_line)
                        if dep_match:
                            # Convert WS-FFF-SS → 00-FFF-SS (assuming same project)
                            old_dep = f"WS-{dep_match.group(1)}-{dep_match.group(2)}"
                            new_dep = (
                                f"00-{dep_match.group(1).zfill(3)}-"
                                f"{dep_match.group(2).zfill(2)}"
                            )
                            dep_line = dep_line.replace(old_dep, new_dep)
                        new_deps.append(dep_line)
                content = content.replace(deps_section, "\n".join(new_deps))

        return content

    def _generate_filename(self) -> str:
        """Generate new filename based on new ID."""
        if self.new_id:
            return self.path.name.replace(self.path.stem.split("-")[0], self.new_id.split("-")[0])
        return self.path.name


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


def main() -> int:
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="Migrate workstream IDs to PP-FFF-SS format",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Dry run for SDP
  python scripts/migrate_workstream_ids.py --dry-run

  # Migrate SDP workstreams
  python scripts/migrate_workstream_ids.py --project-id 00

  # Migrate hw_checker workstreams
  python scripts/migrate_workstream_ids.py --project-id 02 --path ../hw_checker
        """,
    )

    parser.add_argument(
        "--project-id",
        default="00",
        help="Project ID (default: 00 for SDP)",
    )
    parser.add_argument(
        "--path",
        type=Path,
        default=None,
        help="Path to project root (default: script parent directory)",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Show what would be changed without making changes",
    )

    args = parser.parse_args()

    # Determine root path
    if args.path:
        root_path = args.path
    else:
        root_path = Path(__file__).parent.parent

    # Validate project_id format
    if not re.match(r"^\d{2}$", args.project_id):
        print(f"✗ Invalid project_id: {args.project_id}")
        print("  Project ID must be 2 digits (e.g., 00, 01, 02)")
        return 1

    try:
        migrator = WorkstreamMigrator(
            root_path=root_path,
            project_id=args.project_id,
            dry_run=args.dry_run,
        )
        stats = migrator.migrate()

        # Return error code if any failures
        return 1 if stats["failed"] > 0 else 0

    except WorkstreamMigrationError as e:
        print(f"✗ Migration error: {e}")
        return 1
    except KeyboardInterrupt:
        print("\n✗ Migration cancelled by user")
        return 130
    except Exception as e:
        print(f"✗ Unexpected error: {e}")
        import traceback

        traceback.print_exc()
        return 1


if __name__ == "__main__":
    sys.exit(main())
