"""Retry logic for GitHub API calls."""

import functools
import time
from typing import Any, Callable, TypeVar

from github import GithubException

from sdp.github.exceptions import GitHubSyncError, RateLimitError

T = TypeVar("T")


def retry_on_rate_limit(
    max_retries: int = 3,
    base_delay: float = 1.0,
) -> Callable[[Callable[..., T]], Callable[..., T]]:
    """Decorator to retry on rate limit with exponential backoff.

    Args:
        max_retries: Maximum retry attempts
        base_delay: Base delay in seconds (exponentially increased)

    Returns:
        Decorated function

    """

    def decorator(func: Callable[..., T]) -> Callable[..., T]:
        @functools.wraps(func)
        def wrapper(*args: Any, **kwargs: Any) -> T:
            last_exception: GithubException | None = None

            for attempt in range(max_retries + 1):
                try:
                    return func(*args, **kwargs)
                except GithubException as e:
                    # Check if rate limit error
                    if (
                        e.status == 403
                        and "rate limit" in str(e).lower()
                    ):
                        if attempt < max_retries:
                            delay = base_delay * (2**attempt)
                            print(
                                f"Rate limit hit, retrying in {delay}s "
                                f"(attempt {attempt + 1}/{max_retries})..."
                            )
                            time.sleep(delay)
                            last_exception = e
                            continue
                        else:
                            # Max retries exceeded
                            reset_time = getattr(
                                e,
                                "reset_time",
                                int(time.time()) + 3600,
                            )
                            raise RateLimitError(reset_time) from e
                    else:
                        # Not rate limit, re-raise immediately
                        raise

            # Should not reach here, but satisfy type checker
            if last_exception:
                raise last_exception
            raise RuntimeError("Retry logic failed unexpectedly")

        return wrapper

    return decorator


def with_error_handling(
    func: Callable[..., T]
) -> Callable[..., T]:
    """Decorator to convert GitHub exceptions to custom exceptions.

    Args:
        func: Function to wrap

    Returns:
        Wrapped function with error handling

    """

    @functools.wraps(func)
    def wrapper(*args: Any, **kwargs: Any) -> T:
        from sdp.github.exceptions import (
            AuthenticationError,
        )

        try:
            return func(*args, **kwargs)
        except GithubException as e:
            # Authentication error
            if e.status == 401:
                raise AuthenticationError() from e
            # Other GitHub errors
            else:
                data = getattr(e, "data", {})
                message = data.get("message", str(e))
                raise GitHubSyncError(
                    f"GitHub API error: {message}",
                    action="Check GitHub status: "
                    "https://www.githubstatus.com/",
                ) from e
        except Exception as e:
            # Unexpected error
            raise GitHubSyncError(
                f"Unexpected error: {type(e).__name__}: {e}",
                action="Report this error to maintainers",
            ) from e

    return wrapper
