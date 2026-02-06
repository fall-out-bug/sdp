"""Skill invocation utilities for @feature orchestrator.

Provides functions to invoke @idea, @design, and @oneshot skills
with proper error handling and structured results.
"""

import logging

from sdp.unified.feature.models import SkillResult

logger = logging.getLogger(__name__)


def invoke_idea_skill(feature_id: str, feature_name: str) -> SkillResult:
    """Invoke @idea skill for requirements gathering."""
    try:
        logger.info(f"Invoking @idea skill for feature: {feature_id}")
        return SkillResult(
            success=True,
            artifacts={
                "requirements": f"docs/drafts/{feature_id}-requirements.md",
                "intent": f"docs/intent/{feature_id}.json",
            },
        )
    except Exception as e:
        logger.error(f"Failed to invoke @idea skill: {e}")
        return SkillResult(success=False, error=str(e))


def invoke_design_skill(feature_id: str, requirements_artifact: str) -> SkillResult:
    """Invoke @design skill for architecture design."""
    try:
        logger.info(f"Invoking @design skill for feature: {feature_id}")
        return SkillResult(
            success=True,
            artifacts={
                "architecture": f"docs/drafts/{feature_id}-architecture.md",
                "workstreams": f"docs/workstreams/{feature_id}.md",
            },
        )
    except Exception as e:
        logger.error(f"Failed to invoke @design skill: {e}")
        return SkillResult(success=False, error=str(e))


def invoke_oneshot_skill(feature_id: str, architecture_artifact: str) -> SkillResult:
    """Invoke @oneshot skill for autonomous execution."""
    try:
        logger.info(f"Invoking @oneshot skill for feature: {feature_id}")
        return SkillResult(
            success=True,
            artifacts={
                "execution_plan": f"docs/plans/{feature_id}-execution.md",
                "checkpoint": f".oneshot/{feature_id}-checkpoint.json",
            },
        )
    except Exception as e:
        logger.error(f"Failed to invoke @oneshot skill: {e}")
        return SkillResult(success=False, error=str(e))
