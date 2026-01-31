"""Product vision management for feature development."""

from dataclasses import dataclass
from datetime import datetime
from pathlib import Path


@dataclass
class ProductVision:
    """Product vision manifesto."""

    mission: str
    users: list[str]
    success_metrics: list[str]
    strategic_tradeoffs: dict[str, str]
    non_goals: list[str]
    updated: datetime

    def to_markdown(self) -> str:
        """Convert to markdown format."""
        users_md = "\n".join(f"{i+1}. **{u}**" for i, u in enumerate(self.users))
        metrics_md = "\n".join(f"- [ ] {m}" for m in self.success_metrics)
        tradeoffs_md = "\n".join(
            f"- **{k}**: {v}" for k, v in self.strategic_tradeoffs.items()
        )
        nongoals_md = "\n".join(f"- {ng}" for ng in self.non_goals)

        return f"""# PRODUCT_VISION.md

> **Last Updated:** {self.updated.strftime("%Y-%m-%d")}
> **Version:** 1.0

## Mission

{self.mission}

## Users

{users_md}

## Success Metrics

{metrics_md}

## Strategic Tradeoffs

{tradeoffs_md}

## Non-Goals

{nongoals_md}
"""


class VisionManager:
    """Manages PRODUCT_VISION.md file."""

    def __init__(self, project_root: str | Path = "."):
        self._root = Path(project_root)
        self._vision_file = self._root / "PRODUCT_VISION.md"

    def save(self, vision: ProductVision) -> None:
        """Save vision to file.

        Args:
            vision: ProductVision to save
        """
        self._vision_file.write_text(vision.to_markdown())

    def load(self) -> ProductVision | None:
        """Load vision from file.

        Returns:
            ProductVision if file exists, None otherwise
        """
        if not self._vision_file.exists():
            return None

        content = self._vision_file.read_text()

        # Parse mission
        if "## Mission" in content:
            mission_start = content.index("## Mission") + len("## Mission")
            mission_end = content.index("## Users", mission_start)
            mission = content[mission_start:mission_end].strip()
        else:
            mission = "Unknown mission"

        # Parse users
        users = []
        if "## Users" in content:
            users_start = content.index("## Users") + len("## Users")
            users_end = content.index("## Success Metrics", users_start)
            users_section = content[users_start:users_end].strip()
            for line in users_section.split("\n"):
                # Format: "1. **User**" or "**User**"
                if "**" in line:
                    parts = line.strip().split("**")
                    if len(parts) >= 2:
                        user = parts[1].strip()
                        if user:
                            users.append(user)

        return ProductVision(
            mission=mission,
            users=users,
            success_metrics=[],
            strategic_tradeoffs={},
            non_goals=[],
            updated=datetime.now(),
        )
