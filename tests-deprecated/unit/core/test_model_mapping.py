"""Tests for model mapping registry."""

import pytest
from pathlib import Path

from sdp.core.model_mapping import (
    ModelProvider,
    ModelRegistry,
    ModelMappingError,
    load_model_registry,
)


class TestModelProvider:
    """Test ModelProvider dataclass."""

    def test_model_provider_creation(self) -> None:
        """Verify ModelProvider creation."""
        provider = ModelProvider(
            provider="Anthropic",
            model="Claude Sonnet",
            context="200K",
            tool_use=True,
        )

        assert provider.provider == "Anthropic"
        assert provider.model == "Claude Sonnet"
        assert provider.tool_use is True


class TestModelRegistry:
    """Test ModelRegistry."""

    def test_get_models_for_tier(self) -> None:
        """Verify get_models_for_tier returns models."""
        registry = ModelRegistry(
            tiers={
                "T0": [
                    ModelProvider(
                        provider="Anthropic",
                        model="Claude",
                        context="200K",
                        tool_use=True,
                    )
                ],
                "T1": [],
            }
        )

        models = registry.get_models_for_tier("T0")

        assert len(models) == 1
        assert models[0].provider == "Anthropic"

    def test_get_models_raises_invalid_tier(self) -> None:
        """Verify get_models_for_tier raises for invalid tier."""
        registry = ModelRegistry(tiers={"T0": []})

        with pytest.raises(ValueError, match="Invalid tier"):
            registry.get_models_for_tier("T99")


class TestLoadModelRegistry:
    """Test load_model_registry."""

    def test_load_from_file(self) -> None:
        """Verify load_model_registry loads from markdown."""
        mapping_file = (
            Path(__file__).parent.parent.parent.parent
            / "docs"
            / "internals"
            / "model-mapping.md"
        )

        registry = load_model_registry(mapping_file)

        assert registry is not None
        assert "T0" in registry.tiers
        assert "T1" in registry.tiers
        assert len(registry.tiers["T0"]) > 0

    def test_load_raises_file_not_found(self) -> None:
        """Verify load_model_registry raises when file missing."""
        with pytest.raises(ModelMappingError, match="Model mapping file not found"):
            load_model_registry(Path("/nonexistent/model-mapping.md"))
