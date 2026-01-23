#!/usr/bin/env python3
"""
UPC Protocol v2.0 Epic Initializer

Creates the directory structure and initial files for a new epic.
"""

import argparse
import json
from datetime import datetime, timezone
from pathlib import Path


def create_status_json(epic_id: str, tier: str, mode: str) -> dict:
    """Generate initial status.json content."""
    return {
        "epic_id": epic_id,
        "tier": tier,
        "phase": "requirements",
        "mode": mode,
        "iteration": 1,
        "approvals": [],
        "blockers": [],
        "workstreams": [],
        "updated_at": datetime.now(timezone.utc).isoformat(),
        "updated_by": "init"
    }


def create_epic_md(epic_id: str, title: str) -> str:
    """Generate initial epic.md content."""
    return f"""# {title}

## Epic ID
{epic_id}

## Summary
<!-- Describe the problem or feature in 2-3 sentences -->

## Goals
<!-- What are we trying to achieve? -->
- Goal 1
- Goal 2

## Non-Goals
<!-- What are we explicitly NOT doing? -->
- Non-goal 1

## Success Criteria
<!-- How do we know we're done? -->
- Criterion 1
- Criterion 2

## Context
<!-- Background information, links to related docs -->

## Open Questions
<!-- Questions that need answers before implementation -->
- Q1: ...
"""


def init_epic(
    base_dir: Path,
    epic_id: str,
    title: str,
    tier: str = "standard",
    mode: str = "full"
) -> None:
    """Initialize a new epic with UPC v2.0 structure."""
    
    epic_dir = base_dir / epic_id
    consensus_dir = epic_dir / "consensus"
    
    # Create directories
    directories = [
        consensus_dir / "artifacts",
        consensus_dir / "messages" / "inbox" / "analyst",
        consensus_dir / "messages" / "inbox" / "architect",
        consensus_dir / "messages" / "inbox" / "tech_lead",
        consensus_dir / "messages" / "inbox" / "developer",
        consensus_dir / "messages" / "inbox" / "qa",
        consensus_dir / "messages" / "inbox" / "devops",
        consensus_dir / "decision_log",
    ]
    
    for directory in directories:
        directory.mkdir(parents=True, exist_ok=True)
        print(f"✓ Created {directory.relative_to(base_dir)}")
    
    # Create epic.md
    epic_md_path = epic_dir / "epic.md"
    epic_md_path.write_text(create_epic_md(epic_id, title), encoding="utf-8")
    print(f"✓ Created {epic_md_path.relative_to(base_dir)}")
    
    # Create status.json
    status_path = consensus_dir / "status.json"
    status_content = create_status_json(epic_id, tier, mode)
    status_path.write_text(
        json.dumps(status_content, indent=2, ensure_ascii=False),
        encoding="utf-8"
    )
    print(f"✓ Created {status_path.relative_to(base_dir)}")
    
    print()
    print(f"=== Epic '{epic_id}' initialized ===")
    print(f"Tier: {tier}")
    print(f"Mode: {mode}")
    print(f"Next: Edit {epic_md_path.relative_to(base_dir)} and run Analyst agent")


def main():
    parser = argparse.ArgumentParser(description="UPC Protocol v2.0 Epic Initializer")
    parser.add_argument(
        "epic_id",
        help="Epic identifier (e.g., EP-AUTH-001)"
    )
    parser.add_argument(
        "--title",
        default="New Epic",
        help="Epic title (default: 'New Epic')"
    )
    parser.add_argument(
        "--base-dir",
        type=Path,
        default=Path("tools/hw_checker/docs/specs"),
        help="Base directory for epics (default: tools/hw_checker/docs/specs)"
    )
    parser.add_argument(
        "--tier",
        choices=["starter", "standard", "enterprise"],
        default="standard",
        help="Protocol tier (default: standard)"
    )
    parser.add_argument(
        "--mode",
        choices=["full", "fast_track", "hotfix"],
        default="full",
        help="Execution mode (default: full)"
    )
    
    args = parser.parse_args()
    
    base_dir = args.base_dir.resolve()
    
    print(f"=== Initializing Epic: {args.epic_id} ===")
    print()
    
    init_epic(
        base_dir=base_dir,
        epic_id=args.epic_id,
        title=args.title,
        tier=args.tier,
        mode=args.mode
    )


if __name__ == "__main__":
    main()

