"""
Destructive operations detection and user confirmation.

Scans workstreams for potentially destructive operations and prompts for confirmation.
"""

from typing import TYPE_CHECKING, Any

from ..execution_mode import DestructiveOperationDetector

if TYPE_CHECKING:
    from ..client import BeadsClient


def check_destructive_operations_confirmation(
    client: "BeadsClient", feature_id: str
) -> bool:
    """Check if user confirms destructive operations.

    Args:
        client: BeadsClient instance
        feature_id: Parent feature task ID

    Returns:
        True if user confirms or no destructive operations, False otherwise

    Implementation:
    1. Get all subtasks for the feature
    2. Check task titles/descriptions for destructive keywords
    3. Prompt user for confirmation if found
    4. Return False if user declines
    """
    try:
        # Get all subtasks for this feature
        all_tasks = client.list_tasks(parent_id=feature_id)

        if not all_tasks:
            # No tasks found, auto-confirm
            return True

        # Check each task for destructive patterns
        destructive_tasks = []

        for task in all_tasks:
            # Combine title and description for checking
            text_to_check = f"{task.title} {task.description or ''}"

            # Check against destructive patterns
            for (
                category,
                patterns,
            ) in DestructiveOperationDetector.DESTRUCTIVE_PATTERNS.items():
                for pattern in patterns:
                    if pattern.lower() in text_to_check.lower():
                        destructive_tasks.append(
                            {
                                "task_id": task.id,
                                "title": task.title,
                                "operation_type": category,
                                "pattern": pattern,
                            }
                        )
                        break

        if not destructive_tasks:
            # No destructive operations found, auto-confirm
            return True

        # Destructive operations detected - need user confirmation
        operation_summary = _build_destructive_operations_summary(destructive_tasks)

        # Prompt user for confirmation
        return _console_prompt_confirmation(operation_summary)

    except Exception:
        # If detection fails, fail-open (allow execution)
        # This prevents blocking all executions due to bugs in detection logic
        return True


def _build_destructive_operations_summary(
    destructive_tasks: list[dict[str, Any]]
) -> str:
    """Build a summary of destructive operations for user confirmation.

    Args:
        destructive_tasks: List of dicts with operation details

    Returns:
        Formatted summary string
    """
    lines = [
        "⚠️  DESTRUCTIVE OPERATIONS DETECTED",
        "",
        f"Found {len(destructive_tasks)} task(s) with destructive operations:",
        "",
    ]

    for task_info in destructive_tasks:
        lines.append(f"  • {task_info['task_id']}: {task_info['title'][:60]}")
        lines.append(f"    Operation: {task_info['operation_type']}")
        lines.append(f"    Pattern matched: {task_info['pattern']}")
        lines.append("")

    lines.extend(
        [
            "These operations may:",
            "  - Delete files or data",
            "  - Modify database schema",
            "  - Cause irreversible changes",
            "",
            "Do you want to proceed?",
        ]
    )

    return "\n".join(lines)


def _console_prompt_confirmation(summary: str) -> bool:
    """Prompt user via console for confirmation.

    Args:
        summary: Summary message to show user

    Returns:
        True if user confirms, False otherwise
    """
    print("\n" + "=" * 70)
    print(summary)
    print("=" * 70 + "\n")

    while True:
        response = input("Type 'yes' to proceed, 'no' to cancel: ").strip().lower()

        if response in ("yes", "y"):
            print("✓ User confirmed - proceeding with destructive operations\n")
            return True
        elif response in ("no", "n"):
            print("✗ User declined - cancelling execution\n")
            return False
        else:
            print("Please type 'yes' or 'no'\n")
