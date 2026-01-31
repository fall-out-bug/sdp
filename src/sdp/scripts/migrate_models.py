"""Workstream migration data models."""

import re
from pathlib import Path
from typing import Optional, Tuple

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
