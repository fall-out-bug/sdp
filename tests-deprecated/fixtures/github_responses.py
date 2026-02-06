"""Mock GitHub API responses for testing."""

from unittest.mock import Mock

MOCK_REPO = {
    "id": 123456,
    "name": "msu_ai_masters",
    "full_name": "fall-out-bug/msu_ai_masters",
    "owner": {"login": "fall-out-bug"},
}


MOCK_ISSUE = {
    "number": 123,
    "id": 1234567890,
    "node_id": "I_test123",
    "title": "WS-160-01: Test Workstream",
    "body": "Test issue body",
    "state": "open",
    "labels": [
        {"name": "workstream"},
        {"name": "feature/F160"},
        {"name": "size/SMALL"},
        {"name": "status/backlog"},
    ],
    "milestone": {"number": 1},
}


MOCK_MILESTONE = {
    "number": 1,
    "id": 9876543,
    "title": "Feature F160: Test Feature",
    "description": "Test milestone",
    "state": "open",
    "open_issues": 5,
    "closed_issues": 3,
}


MOCK_PROJECT = {
    "id": "PVT_test123",
    "number": 1,
    "title": "mlsd",
    "description": "Test project",
}


MOCK_PROJECT_ITEM = {
    "id": "PVTI_test456",
    "content": {
        "id": "I_test123",
    },
}


def create_mock_github_client() -> tuple[Mock, Mock]:
    """Create mock GitHub client with common responses.

    Returns:
        Tuple of (client, repo) mocks

    """
    client = Mock()
    repo = Mock()

    # Configure repo mock
    repo.get_issue.return_value = Mock(**MOCK_ISSUE)
    repo.create_issue.return_value = Mock(**MOCK_ISSUE)
    repo.get_milestone.return_value = Mock(**MOCK_MILESTONE)
    repo.create_milestone.return_value = Mock(**MOCK_MILESTONE)
    repo.get_milestones.return_value = [Mock(**MOCK_MILESTONE)]
    repo.get_labels.return_value = []

    client.get_repo.return_value = repo
    client._repo = repo

    return client, repo
