"""Mermaid diagram generator.

This module generates Mermaid format diagrams from PRD flow annotations.
"""

from .annotations import Flow


def generate_mermaid_sequence(flow: Flow) -> str:
    """Generate Mermaid sequence diagram from flow steps.

    Args:
        flow: Flow object containing steps

    Returns:
        Mermaid sequence diagram as string
    """
    lines = [
        "sequenceDiagram",
    ]

    # Collect unique participants
    participants = set()
    for step in flow.steps:
        if step.participant:
            participants.add(step.participant)
        # Also add source file as participant
        if step.source_file:
            participants.add(step.source_file.stem)

    # Add participants
    for p in sorted(participants):
        lines.append(f"    participant {p}")

    lines.append("")

    # Generate sequence
    sorted_steps = sorted(flow.steps, key=lambda s: s.step_number)
    for step in sorted_steps:
        source = step.source_file.stem if step.source_file else "Unknown"
        lines.append(f"    Note over {source}: Step {step.step_number}")
        lines.append(f"    Note over {source}: {step.description}")
        lines.append("")

    return "\n".join(lines)


def generate_mermaid_component() -> str:
    """Generate component diagram template.

    Returns:
        Mermaid component diagram as string
    """
    return """flowchart TB
    subgraph Presentation["Presentation Layer"]
        API[FastAPI]
        CLI[Click CLI]
    end

    subgraph Application["Application Layer"]
        UseCase[Use Cases]
        Ports[Ports/Interfaces]
    end

    subgraph Domain["Domain Layer"]
        Entities[Entities]
        Services[Domain Services]
    end

    subgraph Infrastructure["Infrastructure Layer"]
        DB[(PostgreSQL)]
        Queue[(Redis)]
        External[External APIs]
    end

    API --> UseCase
    CLI --> UseCase
    UseCase --> Entities
    UseCase --> Ports
    Ports --> DB
    Ports --> Queue
    Ports --> External
"""


def generate_mermaid_deployment() -> str:
    """Generate deployment diagram template.

    Returns:
        Mermaid deployment diagram as string
    """
    return """graph TB
    subgraph "Docker Environment"
        Client[Clients]

        subgraph "Application Server"
            API[API Service]
            Worker[Background Worker]
        end

        subgraph "Data Layer"
            PG[(PostgreSQL)]
            Redis[(Redis Queue)]
        end

        subgraph "External Services"
            S3[S3 Storage]
            Email[Email Service]
        end
    end

    Client --> API
    API --> PG
    API --> Redis
    Worker --> Redis
    Worker --> PG
    Worker --> S3
    Worker --> Email
"""
