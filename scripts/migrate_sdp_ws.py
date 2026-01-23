#!/usr/bin/env python3
"""Migrate SDP workstreams to PP-FFF-SS format.

Usage:
    python migrate_sdp_ws.py
"""
import re
from pathlib import Path


def migrate_ws_file(ws_file: Path) -> Path:
    """Rename and update WS file from old to new format.

    Args:
        ws_file: Path to WS markdown file

    Returns:
        New path of migrated file
    """
    content = ws_file.read_text(encoding="utf-8")

    # Extract ws_id from frontmatter
    match = re.search(r'ws_id:\s*(WS-[\d-]+)', content)
    if not match:
        print(f"⚠️  No ws_id found in {ws_file}")
        return ws_file

    old_id = match.group(1)
    print(f"Processing {ws_file.name}: {old_id}")

    # Convert WS-XXX-YY to 00-XXX-YY
    new_id = f"00-{old_id[3:]}"  # WS-410-06 → 00-410-06

    # Update frontmatter ws_id
    content = content.replace(f"ws_id: {old_id}", f"ws_id: {new_id}")

    # Add project_id: 00 after ws_id if not present
    if "project_id:" not in content:
        content = re.sub(
            r"(ws_id:\s*[^\n]+\n)",
            r"\1project_id: 00\n",
            content,
            count=1
        )

    # Update title from ## WS-XXX-YY to ## 00-XXX-YY
    content = content.replace(f"## {old_id}:", f"## {new_id}:")

    # Write updated content
    ws_file.write_text(content, encoding="utf-8")

    # Rename file
    old_filename = ws_file.name
    new_filename = old_filename.replace("WS-", "00-")
    new_path = ws_file.parent / new_filename

    if old_filename != new_filename:
        ws_file.rename(new_path)
        print(f"  ✓ Renamed: {old_filename} → {new_filename}")
    else:
        print(f"  ✓ Updated: {old_filename}")

    return new_path


def main() -> None:
    """Migrate all SDP workstreams."""
    sdp_root = Path(__file__).parent.parent
    ws_dir = sdp_root / "docs" / "workstreams"

    if not ws_dir.exists():
        print(f"⚠️  Workstreams directory not found: {ws_dir}")
        return

    print(f"Migrating SDP workstreams in: {ws_dir}")
    print("-" * 60)

    # Find all WS markdown files
    ws_files = list(ws_dir.rglob("WS-*.md"))

    if not ws_files:
        print("No WS-*.md files found (already migrated?)")
        return

    print(f"Found {len(ws_files)} workstream files to migrate\n")

    for ws_file in ws_files:
        migrate_ws_file(ws_file)

    print("-" * 60)
    print("✅ Migration complete")
    print("\nVerification:")
    print(f"  Old format remaining: {len(list(ws_dir.rglob('WS-*.md')))} files")
    print(f"  New format files: {len(list(ws_dir.rglob('00-*.md')))} files")


if __name__ == "__main__":
    main()
