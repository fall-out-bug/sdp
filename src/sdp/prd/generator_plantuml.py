"""PlantUML diagram generator.

This module generates PlantUML format diagrams from PRD flow annotations.
"""

from .annotations import Flow


def generate_plantuml_sequence(flow: Flow) -> str:
    """Generate PlantUML sequence diagram from flow steps.

    Args:
        flow: Flow object containing steps

    Returns:
        PlantUML sequence diagram as string
    """
    lines = [
        "@startuml",
        f"title {flow.name}",
        "",
        "skinparam sequenceMessageAlign center",
        "skinparam noteBackgroundColor #FEFECE",
        "skinparam noteBorderColor #E0E000",
        "",
    ]

    # Generate sequence
    sorted_steps = sorted(flow.steps, key=lambda s: s.step_number)
    for step in sorted_steps:
        source = step.source_file.stem if step.source_file else "Unknown"
        lines.append(f"note over {source}")
        lines.append(f"  Step {step.step_number}: {step.description}")
        lines.append("end note")
        lines.append("")

    lines.append("@enduml")
    return "\n".join(lines)


def generate_plantuml_component() -> str:
    """Generate component diagram template.

    Returns:
        PlantUML component diagram as string
    """
    return """@startuml
!include <C4/C4_Component>

title Component Diagram

Container(system, "System", "Technology") {
    Container(api, "API", "FastAPI", "REST endpoints")
    Container(worker, "Worker", "Python", "Background jobs")

    ContainerDb(db, "Database", "PostgreSQL", "Primary storage")
    ContainerDb(queue, "Queue", "Redis", "Job queue")
}

Rel(api, db, "reads/writes")
Rel(api, queue, "enqueue")
Rel(worker, queue, "dequeue")
Rel(worker, db, "writes")

@enduml
"""


def generate_plantuml_deployment() -> str:
    """Generate deployment diagram template.

    Returns:
        PlantUML deployment diagram as string
    """
    return """@startuml
!include <C4/C4_Deployment>

title Deployment Diagram

Deployment_Node(docker, "Docker Compose", "docker-compose.yml") {
    Container(api, "API", "FastAPI", "REST endpoints")
    Container(worker, "Worker", "Python", "Job processing")
    ContainerDb(pg, "PostgreSQL", "Database", "Persistent storage")
    ContainerDb(redis, "Redis", "Job queue", "Message broker")
}

Rel(api, pg, "reads/writes", "TCP", "5432")
Rel(api, redis, "enqueue", "TCP", "6379")
Rel(worker, redis, "dequeue", "TCP", "6379")
Rel(worker, pg, "writes", "TCP", "5432")

@enduml
"""
