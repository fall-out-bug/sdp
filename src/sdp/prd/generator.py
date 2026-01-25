"""Unified diagram generation interface.

This module provides a unified interface for generating all diagram types.
"""

from pathlib import Path
from typing import TYPE_CHECKING, List

if TYPE_CHECKING:
    from .annotations import FlowStep

from .annotations import Flow
from .generator_mermaid import (
    generate_mermaid_sequence,
    generate_mermaid_component,
    generate_mermaid_deployment,
)
from .generator_plantuml import (
    generate_plantuml_sequence,
    generate_plantuml_component,
    generate_plantuml_deployment,
)


def generate_diagrams(
    flows: List[Flow],
    output_dir: Path,
    project_type: str = "service",
) -> List[Path]:
    """Generate all diagrams for flows and save to output directory.

    Args:
        flows: List of Flow objects to generate diagrams for
        output_dir: Directory to save diagram files
        project_type: Type of project (affects which templates to include)

    Returns:
        List of created file paths
    """
    output_dir.mkdir(parents=True, exist_ok=True)
    created_files = []

    for flow in flows:
        # Mermaid sequence
        mmd_path = output_dir / f"sequence-{flow.name}.mmd"
        mmd_path.write_text(generate_mermaid_sequence(flow))
        created_files.append(mmd_path)

        # PlantUML sequence
        puml_path = output_dir / f"sequence-{flow.name}.puml"
        puml_path.write_text(generate_plantuml_sequence(flow))
        created_files.append(puml_path)

    # Add architecture diagrams based on project type
    if project_type == "service":
        # Component diagram
        comp_path = output_dir / "component-overview.mmd"
        comp_path.write_text(generate_mermaid_component())
        created_files.append(comp_path)

        # Deployment diagram
        deploy_path = output_dir / "deployment-production.puml"
        deploy_path.write_text(generate_plantuml_deployment())
        created_files.append(deploy_path)

    return created_files


def generate_flow_from_steps(
    flow_name: str,
    steps: List["FlowStep"],
) -> Flow:
    """Create a Flow object from a list of FlowStep objects.

    Args:
        flow_name: Name of the flow
        steps: List of FlowStep objects

    Returns:
        Flow object
    """
    return Flow(name=flow_name, steps=steps)
