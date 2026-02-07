"""Tests for GitHub retry logic."""

from unittest.mock import Mock

import pytest

from sdp.github.retry_logic import retry_on_rate_limit, with_error_handling
from sdp.github.exceptions import GitHubSyncError, RateLimitError, AuthenticationError


class TestRetryOnRateLimit:
    """Test retry_on_rate_limit decorator."""

    def test_success_on_first_call(self) -> None:
        """Verify function succeeds without retry when no error."""
        @retry_on_rate_limit(max_retries=2, base_delay=0.01)  # type: ignore[untyped-decorator]
        def succeed() -> str:
            return "ok"

        result = succeed()

        assert result == "ok"

    def test_retry_on_rate_limit_then_succeed(self) -> None:
        """Verify retry on 403 rate limit then success."""
        from github import GithubException

        call_count = 0

        @retry_on_rate_limit(max_retries=2, base_delay=0.001)  # type: ignore[untyped-decorator]
        def flaky() -> str:
            nonlocal call_count
            call_count += 1
            if call_count < 2:
                raise GithubException(403, {"message": "API rate limit exceeded"})
            return "ok"

        result = flaky()

        assert result == "ok"
        assert call_count == 2

    def test_raises_rate_limit_after_max_retries(self) -> None:
        """Verify RateLimitError raised after max retries."""
        from github import GithubException

        @retry_on_rate_limit(max_retries=2, base_delay=0.001)  # type: ignore[untyped-decorator]
        def always_fail() -> str:
            raise GithubException(403, {"message": "API rate limit exceeded"})

        with pytest.raises(RateLimitError):
            always_fail()

    def test_no_retry_on_non_rate_limit_error(self) -> None:
        """Verify non-403 errors are not retried."""
        from github import GithubException

        call_count = 0

        @retry_on_rate_limit(max_retries=2, base_delay=0.01)  # type: ignore[untyped-decorator]
        def auth_fail() -> str:
            nonlocal call_count
            call_count += 1
            raise GithubException(401, {"message": "Bad credentials"})

        with pytest.raises(GithubException):
            auth_fail()

        assert call_count == 1


class TestWithErrorHandling:
    """Test with_error_handling decorator."""

    def test_success_passes_through(self) -> None:
        """Verify success passes through unchanged."""
        @with_error_handling  # type: ignore[untyped-decorator]
        def succeed() -> str:
            return "ok"

        assert succeed() == "ok"

    def test_401_raises_authentication_error(self) -> None:
        """Verify 401 raises AuthenticationError."""
        from github import GithubException

        @with_error_handling  # type: ignore[untyped-decorator]
        def auth_fail() -> str:
            raise GithubException(401, {"message": "Bad credentials"})

        with pytest.raises(AuthenticationError):
            auth_fail()

    def test_other_github_error_raises_sync_error(self) -> None:
        """Verify other GitHub errors raise GitHubSyncError."""
        from github import GithubException

        @with_error_handling  # type: ignore[untyped-decorator]
        def other_fail() -> str:
            raise GithubException(404, {"message": "Not Found"})

        with pytest.raises(GitHubSyncError, match="GitHub API error"):
            other_fail()

    def test_generic_exception_raises_sync_error(self) -> None:
        """Verify generic exceptions raise GitHubSyncError."""

        @with_error_handling  # type: ignore[untyped-decorator]
        def generic_fail() -> str:
            raise ValueError("Something went wrong")

        with pytest.raises(GitHubSyncError, match="Unexpected error"):
            generic_fail()
