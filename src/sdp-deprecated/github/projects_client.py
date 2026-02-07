"""GitHub Projects v2 (GraphQL) API client."""

from typing import Any, cast

import requests


class ProjectsClient:
    """GitHub Projects v2 API wrapper (GraphQL)."""

    GRAPHQL_URL = "https://api.github.com/graphql"

    def __init__(self, token: str, owner: str) -> None:
        """Initialize projects client.

        Args:
            token: GitHub personal access token
            owner: Organization or user name

        """
        self._token = token
        self._owner = owner
        self._headers: dict[str, str] = {
            "Authorization": f"Bearer {token}",
            "Content-Type": "application/json",
        }

    def _query(self, query: str, variables: dict[str, Any] | None = None) -> Any:  # noqa: ANN401
        """Execute GraphQL query.

        Args:
            query: GraphQL query string
            variables: Query variables

        Returns:
            Response data

        Raises:
            RuntimeError: If query fails

        """
        response = requests.post(
            self.GRAPHQL_URL,
            json={"query": query, "variables": variables or {}},
            headers=self._headers,
            timeout=30,
        )
        response.raise_for_status()

        data: dict[str, Any] = response.json()
        if "errors" in data:
            msg = f"GraphQL errors: {data['errors']}"
            raise RuntimeError(msg)

        return data["data"]

    def get_project_by_name(self, name: str) -> dict[str, Any] | None:
        """Find project by name.

        Args:
            name: Project name (e.g., "mlsd")

        Returns:
            Project dict or None if not found

        """
        query = """
        query($owner: String!) {
          repositoryOwner(login: $owner) {
            ... on Organization {
              projectsV2(first: 20) {
                nodes {
                  id
                  number
                  title
                }
              }
            }
            ... on User {
              projectsV2(first: 20) {
                nodes {
                  id
                  number
                  title
                }
              }
            }
          }
        }
        """

        result: Any = self._query(query, {"owner": self._owner})
        owner_node = result.get("repositoryOwner") or {}
        projects: list[Any] = owner_node.get("projectsV2", {}).get("nodes", [])

        for project in projects:
            if project["title"] == name:
                return cast("dict[str, Any]", project)

        return None

    def create_project(self, name: str, description: str = "") -> dict[str, Any]:  # noqa: ARG002
        """Create new project.

        Args:
            name: Project name
            description: Project description (reserved for future use)

        Returns:
            Created project info

        """
        query = """
        mutation($ownerId: ID!, $title: String!) {
          createProjectV2(input: {ownerId: $ownerId, title: $title}) {
            projectV2 {
              id
              number
              title
            }
          }
        }
        """

        # Get owner ID
        owner_id = self._get_owner_id()

        result: Any = self._query(query, {"ownerId": owner_id, "title": name})
        return cast("dict[str, Any]", result["createProjectV2"]["projectV2"])

    def add_issue_to_project(self, project_id: str, issue_id: str) -> dict[str, Any]:
        """Add issue to project board.

        Args:
            project_id: Project global ID
            issue_id: Issue global ID

        Returns:
            Project item info

        """
        query = """
        mutation($projectId: ID!, $contentId: ID!) {
          addProjectV2ItemById(input: {projectId: $projectId, contentId: $contentId}) {
            item {
              id
            }
          }
        }
        """

        result: Any = self._query(query, {"projectId": project_id, "contentId": issue_id})
        return cast("dict[str, Any]", result["addProjectV2ItemById"]["item"])

    def update_item_field(
        self,
        project_id: str,
        item_id: str,
        field_id: str,
        value: str,
    ) -> None:
        """Update project item field (e.g., Status).

        Args:
            project_id: Project ID
            item_id: Item ID
            field_id: Field ID
            value: New value

        """
        query = """
        mutation($projectId: ID!, $itemId: ID!, $fieldId: ID!, $value: String!) {
          updateProjectV2ItemFieldValue(
            input: {
              projectId: $projectId
              itemId: $itemId
              fieldId: $fieldId
              value: {text: $value}
            }
          ) {
            projectV2Item {
              id
            }
          }
        }
        """

        self._query(
            query,
            {
                "projectId": project_id,
                "itemId": item_id,
                "fieldId": field_id,
                "value": value,
            },
        )

    def update_item_single_select(
        self,
        project_id: str,
        item_id: str,
        field_id: str,
        option_id: str,
    ) -> None:
        """Update single-select field value by option ID.

        Args:
            project_id: Project ID
            item_id: Item ID
            field_id: Field ID
            option_id: Single select option ID

        """
        query = """
        mutation(
          $projectId: ID!,
          $itemId: ID!,
          $fieldId: ID!,
          $optionId: String!
        ) {
          updateProjectV2ItemFieldValue(
            input: {
              projectId: $projectId
              itemId: $itemId
              fieldId: $fieldId
              value: {singleSelectOptionId: $optionId}
            }
          ) {
            projectV2Item {
              id
            }
          }
        }
        """
        self._query(
            query,
            {
                "projectId": project_id,
                "itemId": item_id,
                "fieldId": field_id,
                "optionId": option_id,
            },
        )

    def _get_owner_id(self) -> str:
        """Get organization/user global ID.

        Returns:
            Global ID for owner

        """
        query = """
        query($login: String!) {
          repositoryOwner(login: $login) {
            id
          }
        }
        """

        result: Any = self._query(query, {"login": self._owner})
        owner_node = result.get("repositoryOwner")
        if not owner_node:
            msg = f"Owner not found: {self._owner}"
            raise RuntimeError(msg)
        owner_id: str = owner_node["id"]
        return owner_id

    def get_or_create_project(self, name: str, description: str = "") -> dict[str, Any]:
        """Get existing project or create new one.

        Args:
            name: Project name
            description: Project description (reserved for future use)

        Returns:
            Project info

        """
        project = self.get_project_by_name(name)
        if project:
            return project

        return self.create_project(name, description)
