"""PRD annotation data structures.

This module defines the data structures for representing PRD flow
annotations extracted from code.
"""

from dataclasses import dataclass
from pathlib import Path


@dataclass
class FlowStep:
    """A single step in a PRD flow.

    Attributes:
        flow_name: Name of the flow this step belongs to
        step_number: Sequential number of this step in the flow
        description: Human-readable description of what this step does
        source_file: Path to the source file containing this step
        line_number: Line number in the source file
        participant: Optional participant name for sequence diagrams
    """
    flow_name: str
    step_number: int
    description: str
    source_file: Path
    line_number: int
    participant: str | None = None


@dataclass
class Flow:
    """A PRD flow containing multiple steps.

    Attributes:
        name: Name of the flow
        steps: List of steps in this flow
    """
    name: str
    steps: list[FlowStep]

    def add_step(self, step: FlowStep) -> None:
        """Add a step to this flow.

        Args:
            step: The step to add
        """
        self.steps.append(step)

    def get_sorted_steps(self) -> list[FlowStep]:
        """Get steps sorted by step number.

        Returns:
            Sorted list of steps
        """
        return sorted(self.steps, key=lambda s: s.step_number)

    def __len__(self) -> int:
        """Return the number of steps in this flow."""
        return len(self.steps)
