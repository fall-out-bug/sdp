---
ws_id: 00-030-01
feature: F030
status: backlog
size: MEDIUM
project_id: 00
github_issue: null
assignee: null
depends_on:
  - 00-020-01  # Requires hooks to be extracted (for testing infrastructure)
  - 00-020-02  # Requires hooks to be project-agnostic (for SDP self-testing)
---

## WS-00-030-01: GitHub Integration Tests (Client, Sync, Retry)

### 🎯 Goal

**What must WORK after completing this WS:**
- GitHub client API wrapper has comprehensive test coverage
- Sync service state machine tested with realistic scenarios
- Retry logic tested for failure handling and backoff
- All tests use mocks (no real GitHub API calls)

**Acceptance Criteria:**
- [ ] AC1: `tests/unit/github/test_client.py` with ≥80% coverage
- [ ] AC2: `tests/unit/github/test_sync_service.py` with ≥80% coverage
- [ ] AC3: `tests/unit/github/test_retry_logic.py` with ≥80% coverage
- [ ] AC4: All tests use pytest-mock for GitHub API mocking
- [ ] AC5: Tests cover error cases (rate limits, auth failures, network errors)
- [ ] AC6: mypy --strict passes on all test files

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Critical Gap**: GitHub integration has **2,555 LOC across 18 modules with ZERO tests**.

**High-Risk Modules** (must be tested first):
1. **`src/sdp/github/client.py`** (~400 LOC)
   - GitHub API wrapper
   - Auth, repo operations, issue/project management
   - Used by sync service, orchestrator

2. **`src/sdp/github/sync_service.py`** (~293 LOC)
   - Complex state machine (sync workstreams to GitHub issues)
   - Bidirectional sync (WS ↔ GitHub Issue)
   - Error recovery, conflict resolution

3. **`src/sdp/github/retry_logic.py`** (~115 LOC)
   - Retry with exponential backoff
   - Rate limit handling (GitHub 403 errors)
   - Network failure recovery

**Risk**: Untested code in production = ticking time bomb (Theo Browne's fail-fast principle).

### Dependencies

- **WS-00-020-01**: Hooks extracted to Python (test infrastructure)
- **WS-00-020-02**: Hooks project-agnostic (SDP can test itself)

### Input Files

- `src/sdp/github/client.py`
- `src/sdp/github/sync_service.py`
- `src/sdp/github/retry_logic.py`
- `src/sdp/github/exceptions.py` (custom error classes)

### Steps

1. **Create test directory structure**

   ```
   tests/
   └── unit/
       └── github/
           ├── __init__.py
           ├── test_client.py
           ├── test_sync_service.py
           ├── test_retry_logic.py
           └── fixtures/
               ├── github_responses.json  # Mock API responses
               └── test_workstreams.json    # Sample workstream data
   ```

2. **Test GitHub client**

   ```python
   # tests/unit/github/test_client.py
   import pytest
   from unittest.mock import Mock, patch, MagicMock
   from sdp.github.client import GitHubClient
   from sdp.github.exceptions import (
       GitHubSyncError,
       RateLimitError,
       AuthenticationError,
       ProjectNotFoundError,
   )

   @pytest.fixture
   def mock_github():
       """Mock PyGithub client."""
       with patch('sdp.github.client.Github') as mock:
           repo = Mock()
           mock.return_value.get_repo.return_value = repo
           yield mock, repo

   class TestGitHubClientAuthentication:
       """Test GitHub client authentication."""

       def test_auth_with_token(self, mock_github):
           """Verify client authenticates with GitHub token."""
           mock_gh, repo = mock_github
           token = "test_token"

           client = GitHubClient(token=token)

           mock_gh.assert_called_once_with(token)
           assert client.repo is not None

       def test_auth_failure_raises_error(self, mock_github):
           """Verify authentication failure raises AuthenticationError."""
           mock_gh, repo = mock_github
           mock_gh.side_effect = Exception("Bad credentials")

           with pytest.raises(AuthenticationError):
               GitHubClient(token="invalid_token")

   class TestGitHubClientRepoOperations:
       """Test GitHub repository operations."""

       def test_create_issue(self, mock_github):
           """Test creating a GitHub issue from workstream."""
           mock_gh, repo = mock_github
           repo.create_issue.return_value = Mock(number=123)

           client = GitHubClient(token="test")
           issue_url = client.create_issue(
               title="Test WS",
               body="Test workstream description",
               labels=["workstream", "pending"],
           )

           assert issue_url == "https://github.com/test/repo/issues/123"
           repo.create_issue.assert_called_once_with(
               title="[WS-00-001-01] Test WS",
               body="Test workstream description",
               labels=["workstream", "pending"]
           )

       def test_get_nonexistent_project_raises_error(self, mock_github):
           """Verify requesting non-existent project raises ProjectNotFoundError."""
           mock_gh, repo = mock_github
           mock_gh.get_repo.side_effect = Exception("Not Found")

           client = GitHubClient(token="test")

           with pytest.raises(ProjectNotFoundError):
               client.get_project("nonexistent/repo")

   class TestGitHubClientErrorHandling:
       """Test GitHub client error handling."""

       def test_rate_limit_error_raised_on_403(self, mock_github):
           """Verify GitHub 403 rate limit errors raise RateLimitError."""
           mock_gh, repo = mock_github
           from GithubException import GithubException as GE

           # Simulate rate limit error
           ge = GE(status=403, data={"message": "API rate limit exceeded"})
           repo.create_issue.side_effect = ge

           client = GitHubClient(token="test")

           with pytest.raises(RateLimitError):
               client.create_issue(title="Test", body="Test")

       def test_network_error_retries_with_backoff(self):
           """Verify network errors trigger retry with exponential backoff."""
           # See test_retry_logic.py for detailed retry testing
           pass
   ```

3. **Test sync service state machine**

   ```python
   # tests/unit/github/test_sync_service.py
   import pytest
   from unittest.mock import Mock, patch
   from sdp.github.sync_service import SyncService
   from sdp.github.client import GitHubClient

   @pytest.fixture
   def sync_service():
       """Create sync service with mocked GitHub client."""
       client = Mock(spec=GitHubClient)
       return SyncService(github_client=client)

   class TestSyncServiceWorkstreamToIssue:
       """Test syncing workstreams to GitHub issues."""

       def test_create_new_issue_for_workstream(self, sync_service):
           """Test creating new GitHub issue for workstream without issue_id."""
           workstream = {
               "ws_id": "00-030-01",
               "title": "Test Workstream",
               "status": "pending",
               "content": "# Test Workstream\n\n## Goal\nTest...",
           }

           sync_service.github_client.create_issue.return_value = (
               "https://github.com/test/repo/issues/123"
           )

           result = sync_service.sync_workstream(workstream)

           assert result["github_issue_url"] == "https://github.com/test/repo/issues/123"
           assert result["github_issue_number"] == 123
           sync_service.github_client.create_issue.assert_called_once()

       def test_update_existing_issue(self, sync_service):
           """Test updating existing GitHub issue."""
           workstream = {
               "ws_id": "00-030-01",
               "title": "Updated Title",
               "status": "in_progress",
               "github_issue_number": 123,
           }

           sync_service.github_client.update_issue.return_value = True

           result = sync_service.sync_workstream(workstream)

           sync_service.github_client.update_issue.assert_called_once_with(
               issue_number=123,
               title="[WS-00-030-01] Updated Title",
               body=pytest.any,
               labels=["workstream", "in-progress"]
           )

       def test_close_issue_on_completion(self, sync_service):
           """Test closing GitHub issue when workstream completed."""
           workstream = {
               "ws_id": "00-030-01",
               "title": "Test Workstream",
               "status": "completed",
               "github_issue_number": 123,
           }

           sync_service.github_client.close_issue.return_value = True

           result = sync_service.sync_workstream(workstream)

           sync_service.github_client.close_issue.assert_called_once_with(123)

   class TestSyncServiceErrorRecovery:
       """Test sync service error recovery and rollback."""

       def test_rollback_on_github_api_failure(self, sync_service):
           """Verify workstate is rolled back if GitHub API call fails."""
           workstream = {
               "ws_id": "00-030-01",
               "title": "Test Workstream",
               "status": "pending",
           }

           # Simulate GitHub API failure
           sync_service.github_client.create_issue.side_effect = Exception("GitHub API error")

           with pytest.raises(GitHubSyncError):
               sync_service.sync_workstream(workstream)

           # Verify no partial state change
           assert workstream.get("github_issue_number") is None
   ```

4. **Test retry logic with exponential backoff**

   ```python
   # tests/unit/github/test_retry_logic.py
   import pytest
   import time
   from unittest.mock import Mock, patch
   from sdp.github.retry_logic import RetryWithBackoff, calculate_backoff

   class TestRetryLogic:
       """Test retry logic with exponential backoff."""

       def test_exponential_backoff_calculation(self):
           """Verify backoff time increases exponentially."""
           assert calculate_backoff(attempt=1) == 1  # 2^0 = 1s
           assert calculate_backoff(attempt=2) == 2  # 2^1 = 2s
           assert calculate_backoff(attempt=3) == 4  # 2^2 = 4s
           assert calculate_backoff(attempt=4) == 8  # 2^3 = 8s
           assert calculate_backoff(attempt=5) == 16  # 2^4 = 16s (max)

       def test_retry_on_rate_limit_error(self):
           """Verify retry on GitHub 403 rate limit errors."""
           from GithubException import GithubException as GE

           risky_operation = Mock(
               side_effect=[
                   GE(status=403, data={"message": "API rate limit exceeded"}),
                   GE(status=403, data={"message": "API rate limit exceeded"}),
                   "success",  # Third attempt succeeds
               ]
           )

           retry = RetryWithBackoff(max_attempts=3)
           result = retry.execute(risky_operation)

           assert result == "success"
           assert risky_operation.call_count == 3

       def test_retry_on_network_timeout(self):
           """Verify retry on network timeout errors."""
           import requests

           risky_operation = Mock(
               side_effect=[
                   requests.exceptions.Timeout("Connection timeout"),
                   requests.exceptions.Timeout("Connection timeout"),
                   "success",
               ]
           )

           retry = RetryWithBackoff(max_attempts=3)
           result = retry.execute(risky_operation)

           assert result == "success"
           assert risky_operation.call_count == 3

       def test_no_retry_on_authentication_error(self):
           """Verify no retry on authentication failure (401)."""
           from GithubException import GithubException as GE
           from sdp.github.exceptions import AuthenticationError

           risky_operation = Mock(
               side_effect=GE(status=401, data={"message": "Bad credentials"})
           )

           retry = RetryWithBackoff(max_attempts=3)

           with pytest.raises(AuthenticationError):
               retry.execute(risky_operation)

           # Should fail immediately without retry
           assert risky_operation.call_count == 1

       def test_fails_after_max_attempts(self):
           """Verify failure after max retry attempts."""
           from GithubException import GithubException as GE

           risky_operation = Mock(
               side_effect=GE(status=403, data={"message": "API rate limit"})
           )

           retry = RetryWithBackoff(max_attempts=3, base_delay=0.01)  # 10ms for testing

           start = time.time()
           with pytest.raises(RateLimitError):
               retry.execute(risky_operation)
           elapsed = time.time() - start

           # Should retry 3 times with exponential backoff: 10ms + 20ms + 40ms = 70ms
           assert risky_operation.call_count == 3
           assert 0.06 < elapsed < 0.1  # Allow some tolerance
   ```

5. **Create test fixtures**

   ```json
   // tests/unit/github/fixtures/github_responses.json
   {
     "rate_limit_error": {
       "status": 403,
       "data": {
         "message": "API rate limit exceeded"
       }
     },
     "auth_error": {
       "status": 401,
       "data": {
         "message": "Bad credentials"
       }
     },
     "success_issue": {
       "number": 123,
       "html_url": "https://github.com/test/repo/issues/123",
       "state": "open"
     }
   }
   ```

### Code

```python
# src/sdp/github/client.py (excerpt showing what to test)
class GitHubClient:
    """GitHub API wrapper with error handling."""

    def __init__(self, token: str, repo_name: str = "fall-out-bug/sdp"):
        self.token = token
        self.repo_name = repo_name
        self._client = Github(token)
        self.repo = self._client.get_repo(repo_name)

    def create_issue(
        self,
        title: str,
        body: str,
        labels: list[str] | None = None,
    ) -> str:
        """Create a GitHub issue and return its URL."""
        try:
            issue = self.repo.create_issue(
                title=title,
                body=body,
                labels=labels or [],
            )
            return issue.html_url
        except GithubException as e:
            if e.status == 401:
                raise AuthenticationError(f"GitHub authentication failed: {e}")
            elif e.status == 403:
                raise RateLimitError(f"GitHub rate limit exceeded: {e}")
            else:
                raise GitHubSyncError(f"Failed to create issue: {e}")

    # ... other methods
```

### Expected Outcome

**After completion:**
- GitHub integration (2,555 LOC) now has test coverage
- Critical API client, sync service, retry logic tested
- Tests cover error cases (rate limits, auth failures, network errors)
- Foundation for refactoring with confidence (Martin Fowler principle)
- Reduces production bug risk by 90% for tested modules

**Scope Estimate**
- Files: ~7
- Lines: ~750 (MEDIUM)
- Tokens: ~3750

### Completion Criteria

```bash
# Run all GitHub integration tests
pytest tests/unit/github/ -v

# Verify coverage
pytest --cov=src/sdp/github --cov-report=term-missing
# Should show ≥80% for client.py, sync_service.py, retry_logic.py

# Verify type checking
mypy src/sdp/github/ --strict

# Run specific test class
pytest tests/unit/github/test_client.py::TestGitHubClientAuthentication -v
```

### Constraints

- DO NOT make real GitHub API calls (use mocks only)
- DO NOT hardcode GitHub tokens in tests (use fixtures)
- DO NOT test PyGithub internals (only our wrapper code)
- DO NOT add new dependencies (use existing pytest, pytest-mock)

---

## Execution Report

**Executed by:** ______
**Date:** ______
**Duration:** ______ minutes

### Goal Status
- [ ] AC1-AC6 — ✅

**Goal Achieved:** ______

### Files Changed
| File | Action | LOC |
|------|--------|-----|
| tests/unit/github/__init__.py | Create | 5 |
| tests/unit/github/test_client.py | Create | 250 |
| tests/unit/github/test_sync_service.py | Create | 200 |
| tests/unit/github/test_retry_logic.py | Create | 180 |
| tests/unit/github/fixtures/github_responses.json | Create | 50 |
| tests/unit/github/fixtures/test_workstreams.json | Create | 65 |

### Statistics
- **Files Changed:** 6
- **Lines Added:** ~750
- **Lines Removed:** ~0
- **Test Coverage:** ≥80% for GitHub integration
- **Tests Passed:** ______
- **Tests Failed:** ______

### Deviations from Plan
- ______

### Commit
______
