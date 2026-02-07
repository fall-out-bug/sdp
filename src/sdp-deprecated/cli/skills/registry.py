"""Skills registry for discovery and help."""
from dataclasses import dataclass
from enum import Enum


class SkillCategory(Enum):
    """Skill categories for organization."""
    WORKFLOW = "workflow"      # @feature, @idea, @design, @build
    DEBUGGING = "debugging"    # /debug, @issue
    FIXES = "fixes"           # @hotfix, @bugfix
    DEPLOYMENT = "deployment"  # @review, @deploy
    UTILITY = "utility"       # @help, /tdd


@dataclass
class SkillInfo:
    """Skill metadata."""
    name: str
    category: SkillCategory
    description: str
    usage: str
    example: str
    when_to_use: list[str]
    related: list[str]


SKILLS_REGISTRY: dict[str, SkillInfo] = {
    "feature": SkillInfo(
        name="@feature",
        category=SkillCategory.WORKFLOW,
        description="Unified entry point for feature development",
        usage="@feature \"<feature description>\"",
        example="@feature \"Add user authentication\"",
        when_to_use=[
            "Starting a new feature from scratch",
            "When you need full workflow (idea → design → build)",
        ],
        related=["@idea", "@design", "@build"],
    ),
    "idea": SkillInfo(
        name="@idea",
        category=SkillCategory.WORKFLOW,
        description="Interactive requirements gathering",
        usage="@idea \"<idea description>\"",
        example="@idea \"Add user comments\"",
        when_to_use=[
            "Exploring requirements before planning",
            "When you need deep interviewing",
        ],
        related=["@feature", "@design"],
    ),
    "design": SkillInfo(
        name="@design",
        category=SkillCategory.WORKFLOW,
        description="Decompose idea into workstreams",
        usage="@design <beads-id>",
        example="@design beads-comments",
        when_to_use=[
            "After @idea, when requirements are clear",
            "Creating workstream breakdown",
        ],
        related=["@idea", "@build"],
    ),
    "build": SkillInfo(
        name="@build",
        category=SkillCategory.WORKFLOW,
        description="Execute single workstream with TDD",
        usage="@build <WS-ID>",
        example="@build 00-034-01",
        when_to_use=[
            "Implementing a planned workstream",
            "Following TDD cycle",
        ],
        related=["@design", "@oneshot"],
    ),
    "debug": SkillInfo(
        name="/debug",
        category=SkillCategory.DEBUGGING,
        description="Systematic debugging using scientific method",
        usage="/debug \"<problem description>\"",
        example="/debug \"Test fails unexpectedly\"",
        when_to_use=[
            "Unexpected test failure",
            "Bug with unclear cause",
            "Need methodical investigation",
        ],
        related=["@issue"],
    ),
    "issue": SkillInfo(
        name="@issue",
        category=SkillCategory.DEBUGGING,
        description="Classify bug severity and route to fix",
        usage="@issue \"<bug description>\"",
        example="@issue \"Login fails on Firefox\"",
        when_to_use=[
            "Triaging a bug report",
            "Deciding between hotfix/bugfix/backlog",
        ],
        related=["/debug", "@hotfix", "@bugfix"],
    ),
    "hotfix": SkillInfo(
        name="@hotfix",
        category=SkillCategory.FIXES,
        description="Emergency P0 fix (production down)",
        usage="@hotfix \"<P0 issue>\"",
        example="@hotfix \"Critical API outage\"",
        when_to_use=[
            "Production is down",
            "Security vulnerability",
            "Data loss risk",
        ],
        related=["@bugfix", "@issue"],
    ),
    "bugfix": SkillInfo(
        name="@bugfix",
        category=SkillCategory.FIXES,
        description="Quality fix for P1/P2 bugs",
        usage="@bugfix \"<bug description>\"",
        example="@bugfix \"Incorrect totals in report\"",
        when_to_use=[
            "Bug affecting users but not critical",
            "Quality issue found in testing",
        ],
        related=["@hotfix", "@issue"],
    ),
    "review": SkillInfo(
        name="@review",
        category=SkillCategory.DEPLOYMENT,
        description="Quality review before deployment",
        usage="@review <feature-id>",
        example="@review F034",
        when_to_use=[
            "All workstreams completed",
            "Before deployment",
        ],
        related=["@deploy"],
    ),
    "deploy": SkillInfo(
        name="@deploy",
        category=SkillCategory.DEPLOYMENT,
        description="Deploy feature to production",
        usage="@deploy <feature-id>",
        example="@deploy F034",
        when_to_use=[
            "After review approval",
            "UAT completed",
        ],
        related=["@review"],
    ),
    "oneshot": SkillInfo(
        name="@oneshot",
        category=SkillCategory.WORKFLOW,
        description="Autonomous execution of all workstreams",
        usage="@oneshot <feature-id>",
        example="@oneshot F034",
        when_to_use=[
            "Feature has multiple WS to execute",
            "Want autonomous execution",
        ],
        related=["@build", "@feature"],
    ),
    "help": SkillInfo(
        name="@help",
        category=SkillCategory.UTILITY,
        description="Interactive skill discovery",
        usage="@help [query]",
        example="@help \"how to fix a bug\"",
        when_to_use=[
            "Not sure which skill to use",
            "Learning SDP workflow",
        ],
        related=[],
    ),
}
