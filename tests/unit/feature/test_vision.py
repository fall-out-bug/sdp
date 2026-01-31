"""Tests for product vision management."""

from datetime import datetime
from pathlib import Path

from sdp.feature.vision import ProductVision, VisionManager


def test_product_vision_to_markdown():
    """Test converting ProductVision to markdown."""
    vision = ProductVision(
        mission="Make development faster",
        users=["Developers", "DevOps engineers"],
        success_metrics=["<1 hour to running code", "30 min onboarding"],
        strategic_tradeoffs={"DX vs Control": "Prioritize DX"},
        non_goals=["Real-time collaboration"],
        updated=datetime(2026, 1, 26, 12, 0, 0),
    )

    md = vision.to_markdown()

    assert "Make development faster" in md
    assert "**Developers**" in md
    assert "- [ ] <1 hour to running code" in md
    assert "**DX vs Control**: Prioritize DX" in md
    assert "- Real-time collaboration" in md
    assert "2026-01-26" in md


def test_vision_manager_save_and_load(tmp_path):
    """Test saving and loading product vision."""
    vision = ProductVision(
        mission="Test mission",
        users=["Test user"],
        success_metrics=["Test metric"],
        strategic_tradeoffs={},
        non_goals=[],
        updated=datetime(2026, 1, 26),
    )

    manager = VisionManager(tmp_path)
    manager.save(vision)

    loaded = manager.load()
    assert loaded is not None
    assert loaded.mission == "Test mission"
    assert "Test user" in loaded.users


def test_vision_manager_load_returns_none_when_missing(tmp_path):
    """Test loading when file doesn't exist."""
    manager = VisionManager(tmp_path)
    assert manager.load() is None
