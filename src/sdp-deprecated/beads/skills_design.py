"""
Feature decomposition into Beads workstreams.

Implements @design skill logic: decompose a feature task into
sub-tasks with sequential dependencies.
"""

from dataclasses import dataclass, field
from typing import List, Optional

from .client import BeadsClient
from .models import BeadsPriority


@dataclass
class WorkstreamSpec:
    """Specification for a workstream to create."""

    title: str
    sequence: int = 1
    size: str = "MEDIUM"
    dependencies: List[str] = field(default_factory=list)


class FeatureDecomposer:
    """Decompose feature tasks into workstreams with dependencies."""

    # Default workstream template
    DEFAULT_WORKSTREAMS = [
        WorkstreamSpec(title="Domain entities", sequence=1, size="MEDIUM"),
        WorkstreamSpec(title="Repository layer", sequence=2, size="MEDIUM"),
        WorkstreamSpec(title="Service layer", sequence=3, size="MEDIUM"),
    ]

    def __init__(self, client: BeadsClient):
        """Initialize decomposer.

        Args:
            client: BeadsClient instance (mock or real)
        """
        self.client = client

    def decompose(
        self,
        feature_id: str,
        workstreams: Optional[List[WorkstreamSpec]] = None,
    ) -> List[str]:
        """Decompose feature into workstreams.

        Creates sub-tasks in Beads with sequential dependencies.
        Each workstream blocks the next one.

        Args:
            feature_id: Parent feature task ID
            workstreams: Custom workstream specs (uses default if None)

        Returns:
            List of created workstream task IDs

        Example:
            decomposer = FeatureDecomposer(client)
            ws_ids = decomposer.decompose("bd-0001")

            # Creates:
            # bd-0001.1: Domain entities (ready)
            # bd-0001.2: Repository (blocked by bd-0001.1)
            # bd-0001.3: Service (blocked by bd-0001.2)
        """
        from . import BeadsDependency, BeadsDependencyType, BeadsTaskCreate

        # Use default workstreams if none provided
        if workstreams is None:
            workstreams = self.DEFAULT_WORKSTREAMS

        # Create workstreams with sequential dependencies
        ws_ids = []
        prev_ws_id = None

        for ws_spec in workstreams:
            # Build dependencies
            dependencies: list[BeadsDependency] = []
            if prev_ws_id is not None:
                dependencies.append(
                    BeadsDependency(
                        task_id=prev_ws_id, type=BeadsDependencyType.BLOCKS
                    )
                )

            # Create workstream task
            ws_task = self.client.create_task(
                BeadsTaskCreate(
                    title=ws_spec.title,
                    parent_id=feature_id,
                    priority=BeadsPriority.MEDIUM,
                    dependencies=dependencies,
                    sdp_metadata={
                        "sequence": ws_spec.sequence,
                        "size": ws_spec.size,
                    },
                )
            )

            ws_ids.append(ws_task.id)
            prev_ws_id = ws_task.id

        return ws_ids
