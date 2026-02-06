"""Tests for builder router model selection."""

import pytest
from pathlib import Path

from sdp.core.builder_router import (
    select_model_weighted,
    select_model_for_tier,
    BuilderRouter,
    RetryPolicy,
    DEFAULT_WEIGHTS,
)
from sdp.core.model_mapping import ModelProvider, ModelRegistry, load_model_registry
from sdp.core.workstream import Workstream, WorkstreamStatus, WorkstreamSize


@pytest.fixture
def sample_models() -> list[ModelProvider]:
    """Create sample model providers for testing."""
    return [
        ModelProvider(
            provider="Anthropic",
            model="Claude Sonnet",
            context="200K",
            tool_use=True,
            cost_per_1m_tokens=3.0,
            availability_pct=0.99,
            context_window=200000,
        ),
        ModelProvider(
            provider="OpenAI",
            model="GPT-4o-mini",
            context="128K",
            tool_use=True,
            cost_per_1m_tokens=0.15,
            availability_pct=0.995,
            context_window=128000,
        ),
    ]


class TestSelectModelWeighted:
    """Test weighted model selection."""

    def test_select_model_returns_provider(self, sample_models: list[ModelProvider]) -> None:
        """Verify select_model_weighted returns a model."""
        result = select_model_weighted(sample_models, DEFAULT_WEIGHTS)

        assert result is not None
        assert result in sample_models

    def test_select_model_filters_by_context(self, sample_models: list[ModelProvider]) -> None:
        """Verify select_model_weighted filters by required_context."""
        result = select_model_weighted(
            sample_models, DEFAULT_WEIGHTS, required_context=150000
        )

        assert result.context_window >= 150000

    def test_select_model_raises_no_candidates(self, sample_models: list[ModelProvider]) -> None:
        """Verify select_model_weighted raises when no models meet context."""
        with pytest.raises(ValueError, match="No models with context"):
            select_model_weighted(
                sample_models, DEFAULT_WEIGHTS, required_context=500000
            )


class TestRetryPolicy:
    """Test retry policy."""

    def test_should_retry_within_limit(self) -> None:
        """Verify should_retry allows retries within limit."""
        policy = RetryPolicy(max_attempts=3)
        assert policy.should_retry(1) is True
        assert policy.should_retry(3) is True

    def test_should_escalate_after_failures(self) -> None:
        """Verify should_escalate when failed and exceeded attempts."""
        policy = RetryPolicy(max_attempts=3)
        assert policy.should_escalate(4, failed=True) is True
        assert policy.should_escalate(3, failed=True) is False


class TestBuilderRouter:
    """Test BuilderRouter."""

    def test_select_model_for_workstream(self) -> None:
        """Verify router selects model for workstream."""
        registry = load_model_registry(
            Path(__file__).parent.parent.parent.parent / "docs" / "internals" / "model-mapping.md"
        )
        router = BuilderRouter(registry=registry)
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
        )

        result = router.select_model(ws)

        assert result is not None
        assert result.provider
        assert result.model

    def test_get_retry_policy_t2(self) -> None:
        """Verify T2 gets 3 attempts."""
        registry = load_model_registry(
            Path(__file__).parent.parent.parent.parent / "docs" / "internals" / "model-mapping.md"
        )
        router = BuilderRouter(registry=registry)
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
        )
        setattr(ws, "capability_tier", "T2")

        policy = router.get_retry_policy(ws)

        assert policy.max_attempts == 3

    def test_should_escalate_to_human_after_max_attempts(self) -> None:
        """Verify should_escalate_to_human returns True after max attempts for T2."""
        registry = load_model_registry(
            Path(__file__).parent.parent.parent.parent / "docs" / "internals" / "model-mapping.md"
        )
        router = BuilderRouter(registry=registry)
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
        )
        setattr(ws, "capability_tier", "T2")

        # Should not escalate before max attempts
        assert router.should_escalate_to_human(ws, attempt=2, failed=True) is False
        # Should escalate after max attempts
        assert router.should_escalate_to_human(ws, attempt=4, failed=True) is True

    def test_create_escalation_error(self) -> None:
        """Verify create_escalation_error creates proper error."""
        registry = load_model_registry(
            Path(__file__).parent.parent.parent.parent / "docs" / "internals" / "model-mapping.md"
        )
        router = BuilderRouter(registry=registry)
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
        )
        setattr(ws, "capability_tier", "T2")

        error = router.create_escalation_error(ws, attempt=3, diagnostics="Test failure")

        assert error.context["ws_id"] == "00-001-01"
        assert error.context["tier"] == "T2"
        assert error.context["attempts"] == 3
        assert error.context["diagnostics"] == "Test failure"
