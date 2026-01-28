"""Feature orchestrator implementation.

Orchestrates @idea → @design → @oneshot with progressive menu
and approval gates for unified feature development workflow.
"""

import logging
from dataclasses import dataclass, field
from typing import Any, Callable, Optional

from sdp.unified.feature.models import FeatureExecution, FeaturePhase

logger = logging.getLogger(__name__)


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


def call_skill(skill_name: str) -> None:
    """Call a skill by name."""
    logger.info(f"Calling @{skill_name} skill")


class FeatureOrchestrator:
    """Orchestrates unified @feature workflow.

    Manages progressive menu, skill invocation, approval gates,
    and checkpoint integration throughout feature development.
    """

    def __init__(
        self, repository: Optional[Callable[[], Any]] = None
    ) -> None:
        """Initialize orchestrator."""
        self.current_phase = FeaturePhase.REQUIREMENTS
        self._repository = repository

    def generate_progressive_menu(
        self, execution: FeatureExecution, step: int
    ) -> str:
        """Generate progressive menu for current step."""
        total_steps = 5
        completed = len(execution.completed_phases)
        progress_pct = completed / total_steps * 100

        # Step emoji mapping
        step_emojis = {
            1: "📋",
            2: "🏗️",
            3: "⚙️",
            4: "✅",
            5: "🚀",
        }

        # Phase names
        phase_names = {
            1: "Requirements Gathering",
            2: "Architecture Design",
            3: "Implementation",
            4: "Quality Assurance",
            5: "User Acceptance",
        }

        # Build progress bar
        filled = int(completed / total_steps * 10)
        bar = "█" * filled + "░" * (10 - filled)

        # Build menu sections
        lines = []
        lines.append("=" * 50)
        lines.append(f"Feature: {execution.feature_name}")
        lines.append("=" * 50)
        lines.append("")
        emoji = step_emojis.get(step, "")
        name = phase_names.get(step, "Unknown")
        lines.append(f"Step {step} of {total_steps}: {emoji} {name}")
        lines.append("")
        lines.append(f"Progress: [{bar}] {progress_pct:.0f}%")
        lines.append("-" * 50)

        # Skip flags
        if step == 1 and execution.skip_flags.skip_requirements:
            lines.append("")
            lines.append("⚠️  [SKIP] Requirements phase will be skipped")
        elif step == 2 and execution.skip_flags.skip_architecture:
            lines.append("")
            lines.append("⚠️  [SKIP] Architecture phase will be skipped")

        lines.append("")
        lines.append("Remaining steps:")

        # Show remaining steps
        for i in range(step, total_steps + 1):
            emoji = step_emojis.get(i, "")
            name = phase_names.get(i, "")
            if i == step:
                lines.append(f"  → {emoji} {name} (current)")
            else:
                lines.append(f"    {emoji} {name}")

        lines.append("")
        lines.append("=" * 50)

        return "\n".join(lines)

    def execute_step(
        self, execution: FeatureExecution, step: int
    ) -> StepResult:
        """Execute a single workflow step."""
        if step == 1 and execution.skip_flags.skip_requirements:
            return StepResult(success=True, phase=FeaturePhase.ARCHITECTURE)
        step_handlers = {
            1: ("idea", FeaturePhase.REQUIREMENTS),
            2: ("design", FeaturePhase.ARCHITECTURE),
            3: ("oneshot", FeaturePhase.EXECUTION),
        }
        if step in step_handlers:
            skill_name, phase = step_handlers[step]
            call_skill(skill_name)
            execution.completed_phases.append(phase)
            return StepResult(success=True, phase=phase)
        return StepResult(success=False, phase=None)

    def execute_feature(self, execution: FeatureExecution) -> StepResult:
        """Execute complete feature workflow."""
        if all([
            execution.skip_flags.skip_requirements,
            execution.skip_flags.skip_architecture,
            execution.skip_flags.skip_uat,
        ]):
            return StepResult(success=True, phase=FeaturePhase.EXECUTION)
        for step in [1, 2, 3]:
            self.execute_step(execution, step)
        return StepResult(success=True, phase=FeaturePhase.EXECUTION)

    def request_approval(
        self,
        phase: FeaturePhase,
        artifacts_path: str,
        execution: Optional[FeatureExecution] = None,
    ) -> str:
        """Request approval for a phase."""
        if execution:
            if phase == FeaturePhase.REQUIREMENTS and execution.skip_flags.skip_requirements:
                return "skipped"
            if phase == FeaturePhase.ARCHITECTURE and execution.skip_flags.skip_architecture:
                return "skipped"
        try:
            response = input(f"Approve {phase.value}? (y/n): ")
            return "approved" if response.lower() == "y" else "rejected"
        except EOFError:
            return "approved"

    def after_phase(self, execution: FeatureExecution, phase: FeaturePhase) -> None:
        """Handle post-phase checkpoint save."""
        if self._repository:
            phases_list = [p.value for p in execution.completed_phases]
            checkpoint = MockCheckpoint(
                feature=execution.feature_id,
                completed_phases=phases_list,
                metrics={"completed_phases": phases_list},
            )
            self._repository().save_checkpoint(checkpoint)
