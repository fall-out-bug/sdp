"""Feature orchestrator implementation.

Orchestrates @idea → @design → @oneshot with progressive menu
and approval gates for unified feature development workflow.
"""

import logging
from typing import Any, Callable, Optional

from sdp.unified.feature.models import (
    FeatureExecution,
    FeaturePhase,
    MockCheckpoint,
    StepResult,
)
from sdp.unified.feature.skills import (
    invoke_design_skill,
    invoke_idea_skill,
    invoke_oneshot_skill,
)

logger = logging.getLogger(__name__)


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

        try:
            if step == 1:
                result = invoke_idea_skill(execution.feature_id, execution.feature_name)
                if result.success:
                    execution.completed_phases.append(FeaturePhase.REQUIREMENTS)
                    return StepResult(success=True, phase=FeaturePhase.REQUIREMENTS)
                return StepResult(success=False, phase=None)

            if step == 2:
                result = invoke_design_skill(execution.feature_id, "requirements.md")
                if result.success:
                    execution.completed_phases.append(FeaturePhase.ARCHITECTURE)
                    return StepResult(success=True, phase=FeaturePhase.ARCHITECTURE)
                return StepResult(success=False, phase=None)

            if step == 3:
                result = invoke_oneshot_skill(execution.feature_id, "architecture.md")
                if result.success:
                    execution.completed_phases.append(FeaturePhase.EXECUTION)
                    return StepResult(success=True, phase=FeaturePhase.EXECUTION)
                return StepResult(success=False, phase=None)
        except Exception as e:
            logger.error(f"Error executing step {step}: {e}")
            return StepResult(success=False, phase=None)

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
