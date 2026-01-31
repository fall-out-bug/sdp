"""Models for @feature skill orchestrator.

Provides data models for feature execution state, phases, and skip flags.
"""

from dataclasses import dataclass, field
from enum import Enum
from typing import Optional


class FeaturePhase(str, Enum):
    """Feature development phases."""

    REQUIREMENTS = "requirements"
    ARCHITECTURE = "architecture"
    EXECUTION = "execution"


@dataclass
class SkipFlags:
    """Skip flags for bypassing approval gates."""

    skip_requirements: bool = False
    skip_architecture: bool = False
    skip_uat: bool = False


@dataclass
class FeatureExecution:
    """Feature execution state.

    Tracks progress through @feature workflow phases including
    requirements, architecture, and execution with skip flag support.
    """

    feature_id: str
    feature_name: str
    skip_flags: SkipFlags = field(default_factory=SkipFlags)
    completed_phases: list[FeaturePhase] = field(default_factory=list)
    status: str = "in_progress"


@dataclass
class SkillResult:
    """Result from skill invocation."""

    success: bool
    artifacts: dict[str, str] = field(default_factory=dict)
    error: Optional[str] = None


@dataclass
class MockCheckpoint:
    """Mock checkpoint for testing."""

    feature: str
    completed_phases: list[str] = field(default_factory=list)
    metrics: dict[str, object] = field(default_factory=dict)


class StepResult:
    """Result from executing a step."""

    def __init__(self, success: bool, phase: Optional[FeaturePhase] = None) -> None:
        self.success = success
        self.phase = phase
