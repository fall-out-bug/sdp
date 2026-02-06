"""Sync workstreams to GitHub project boards."""

from typing import Any, cast

from github import Issue

from sdp.github.projects_client import ProjectsClient

# Status column mapping
STATUS_TO_COLUMN: dict[str, str] = {
    "backlog": "Backlog",
    "active": "In Progress",
    "completed": "Done",
    "blocked": "Blocked",
}

STATUS_OPTION_ALIASES: dict[str, list[str]] = {
    "Backlog": ["Backlog", "Todo", "To Do", "Planned"],
    "In Progress": ["In Progress", "In progress", "Doing"],
    "Done": ["Done", "Completed"],
    "Blocked": ["Blocked", "On Hold"],
}


class ProjectBoardSync:
    """Manage project board automation."""

    def __init__(self, projects_client: ProjectsClient, project_name: str) -> None:
        """Initialize board sync.

        Args:
            projects_client: GraphQL projects client
            project_name: Project name (e.g., "mlsd")

        """
        self._client = projects_client
        self._project_name = project_name
        self._project: dict[str, Any] | None = None
        self._status_field_id: str | None = None
        self._status_option_ids: dict[str, str] | None = None

    def _ensure_project(self) -> dict[str, Any]:
        """Ensure project exists, create if missing.

        Returns:
            Project info

        """
        if self._project is None:
            self._project = self._client.get_or_create_project(
                name=self._project_name,
                description=f"Workstreams for {self._project_name} course",
            )
        return self._project

    def add_issue_to_board(self, issue: Issue.Issue) -> str:
        """Add issue to project board.

        Args:
            issue: PyGithub Issue instance

        Returns:
            Project item ID

        """
        project = self._ensure_project()

        # Add issue to project
        item = self._client.add_issue_to_project(
            project_id=project["id"],
            issue_id=issue.node_id,  # Global ID
        )

        # Set initial status to Backlog
        self._update_item_status(cast(str, item["id"]), "backlog")

        return cast(str, item["id"])

    def update_issue_status(self, issue: Issue.Issue, ws_status: str) -> None:
        """Update issue status on board.

        Args:
            issue: PyGithub Issue instance
            ws_status: WS status (backlog, active, completed, blocked)

        """
        project = self._ensure_project()

        # Find item for issue
        item_id = self._find_item_for_issue(project["id"], issue.node_id)
        if not item_id:
            # Issue not on board yet, add it
            item_id = self.add_issue_to_board(issue)

        # Update status
        self._update_item_status(item_id, ws_status)

    def _find_item_for_issue(self, project_id: str, issue_id: str) -> str | None:
        """Find project item ID for issue.

        Args:
            project_id: Project ID
            issue_id: Issue global ID

        Returns:
            Item ID or None if not found

        """
        # Query project items to find matching issue
        query = """
        query($projectId: ID!) {
          node(id: $projectId) {
            ... on ProjectV2 {
              items(first: 100) {
                nodes {
                  id
                  content {
                    ... on Issue {
                      id
                    }
                  }
                }
              }
            }
          }
        }
        """

        result: Any = self._client._query(query, {"projectId": project_id})
        items: list[Any] = result["node"]["items"]["nodes"]

        for item in items:
            if item["content"]["id"] == issue_id:
                return cast(str, item["id"])

        return None

    def _update_item_status(self, item_id: str, ws_status: str) -> None:
        """Update item status field.

        Args:
            item_id: Project item ID
            ws_status: WS status

        """
        project = self._ensure_project()

        # Get status field ID and options (cached)
        if self._status_field_id is None or self._status_option_ids is None:
            field_id, option_ids = self._get_status_field_info(project["id"])
            self._status_field_id = field_id
            self._status_option_ids = option_ids

        # Map WS status to column
        column_name = STATUS_TO_COLUMN.get(ws_status, "Backlog")

        if self._status_option_ids is None:
            msg = "Status field options not loaded"
            raise ValueError(msg)

        option_id = self._status_option_ids.get(column_name)
        if not option_id:
            for alias in STATUS_OPTION_ALIASES.get(column_name, []):
                option_id = self._status_option_ids.get(alias)
                if option_id:
                    break
        if not option_id:
            msg = f"Status option not found: {column_name}"
            raise ValueError(msg)

        # Update field with single select option ID
        self._client.update_item_single_select(
            project_id=project["id"],
            item_id=item_id,
            field_id=self._status_field_id,
            option_id=option_id,
        )

    def _get_status_field_info(self, project_id: str) -> tuple[str, dict[str, str]]:
        """Get Status field ID and options from project.

        Args:
            project_id: Project ID

        Returns:
            Tuple of Status field ID and option name → option ID map

        Raises:
            ValueError: If Status field not found

        """
        query = """
        query($projectId: ID!) {
          node(id: $projectId) {
            ... on ProjectV2 {
              fields(first: 20) {
                nodes {
                  ... on ProjectV2SingleSelectField {
                    id
                    name
                    options {
                      id
                      name
                    }
                  }
                }
              }
            }
          }
        }
        """

        result: Any = self._client._query(query, {"projectId": project_id})
        fields: list[Any] = result["node"]["fields"]["nodes"]

        # Find "Status" field
        for field in fields:
            if field.get("name") == "Status":
                options = {
                    option["name"]: option["id"]
                    for option in field.get("options", [])
                }
                return cast(str, field["id"]), options

        msg = "Status field not found in project"
        raise ValueError(msg)
