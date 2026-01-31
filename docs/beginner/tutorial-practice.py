"""Tutorial practice file - Example feature implementation.

This file demonstrates the SDP workflow with a simple "greeting" feature.
Follow along with the tutorial at docs/beginner/00-quick-start.md
"""

from typing import Literal


def greet(name: str, style: Literal["formal", "casual", "friendly"] = "friendly") -> str:
    """Generate a greeting message.

    Args:
        name: The person's name
        style: Greeting style (formal, casual, or friendly)

    Returns:
        A personalized greeting message

    Examples:
        >>> greet("Alice")
        'Hello, Alice! 👋'
        >>> greet("Bob", style="formal")
        'Good day, Mr. Bob.'
        >>> greet("Charlie", style="casual")
        'Yo Charlie!'
    """
    if style == "formal":
        return f"Good day, Mr. {name}."
    elif style == "casual":
        return f"Yo {name}!"
    else:  # friendly
        return f"Hello, {name}! 👋"


def greet_multiple(names: list[str], style: Literal["formal", "casual", "friendly"] = "friendly") -> list[str]:
    """Generate greetings for multiple people.

    Args:
        names: List of names
        style: Greeting style

    Returns:
        List of greeting messages
    """
    return [greet(name, style) for name in names]


def greet_with_time(name: str, hour: int) -> str:
    """Generate a time-appropriate greeting.

    Args:
        name: The person's name
        hour: Current hour (0-23)

    Returns:
        A time-appropriate greeting

    Examples:
        >>> greet_with_time("Alice", 9)
        'Good morning, Alice!'
        >>> greet_with_time("Bob", 14)
        'Good afternoon, Bob!'
        >>> greet_with_time("Charlie", 20)
        'Good evening, Charlie!'
    """
    if 5 <= hour < 12:
        period = "morning"
    elif 12 <= hour < 17:
        period = "afternoon"
    else:
        period = "evening"

    return f"Good {period}, {name}!"


# Example usage (delete this in your actual implementation)
if __name__ == "__main__":
    print(greet("Alice"))
    print(greet("Bob", style="formal"))
    print(greet_with_time("Charlie", 14))
