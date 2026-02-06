"""RoleLoader for loading agent role prompts from filesystem.

Provides role loading, caching, and lookup functionality for agent
prompt management.
"""

import logging
import re
from dataclasses import dataclass, field
from pathlib import Path
from typing import Optional

logger = logging.getLogger(__name__)


@dataclass
class Role:
    """Agent role with prompt and metadata.

    Attributes:
        name: Role name/identifier
        prompt: Role prompt/instructions
        description: Optional role description
        capabilities: List of role capabilities
        metadata: Additional role metadata
    """
    name: str
    prompt: str
    description: Optional[str] = None
    capabilities: list[str] = field(default_factory=list)
    metadata: dict[str, str] = field(default_factory=dict)


class RoleLoader:
    """Loader for agent role prompts from filesystem.

    Loads role definitions from markdown files in .claude/agents/,
    caches them for performance, and provides lookup functionality.
    """

    def __init__(self, agents_dir: Optional[Path] = None) -> None:
        """Initialize role loader.

        Args:
            agents_dir: Directory containing agent role files.
                Defaults to .claude/agents/
        """
        self._agents_dir = agents_dir or Path(".claude/agents")
        self._role_cache: dict[str, Role] = {}

    def load_role(self, role_name: str) -> Optional[Role]:
        """Load role from filesystem or cache.

        Args:
            role_name: Name of role to load (filename without .md)

        Returns:
            Role if found, None otherwise
        """
        # Check cache first
        if role_name in self._role_cache:
            return self._role_cache[role_name]

        # Build file path
        role_file = self._agents_dir / f"{role_name}.md"

        # Check if file exists
        if not role_file.exists() or not role_file.is_file():
            logger.warning(f"Role file not found: {role_file}")
            return None

        # Read and parse file
        try:
            content = role_file.read_text()
            role = self._parse_role(role_name, content)
            if role:
                self._role_cache[role_name] = role
            return role
        except IOError as e:
            logger.error(f"Failed to read role file {role_file}: {e}")
            return None

    def get_role(self, role_name: str) -> Optional[Role]:
        """Get role from cache.

        Args:
            role_name: Name of role to retrieve

        Returns:
            Role if in cache, None otherwise
        """
        return self._role_cache.get(role_name)

    def list_roles(self) -> list[str]:
        """List all cached role names.

        Returns:
            List of role names in cache
        """
        return list(self._role_cache.keys())

    def clear_cache(self) -> None:
        """Clear role cache."""
        self._role_cache.clear()

    def _parse_role(self, role_name: str, content: str) -> Optional[Role]:
        """Parse role from markdown content.

        Args:
            role_name: Role name/identifier
            content: Markdown file content

        Returns:
            Parsed Role or None if parsing fails
        """
        # Extract heading (title)
        heading_match = re.search(r'^#\s+(.+)$', content, re.MULTILINE)
        description = heading_match.group(1).strip() if heading_match else None

        # Extract capabilities section
        capabilities = self._extract_capabilities(content)

        # Use entire content as prompt (minus heading)
        prompt = content
        if heading_match:
            prompt = content[heading_match.end():].strip()

        return Role(
            name=role_name,
            prompt=prompt,
            description=description,
            capabilities=capabilities,
        )

    def _extract_capabilities(self, content: str) -> list[str]:
        """Extract capabilities from markdown content.

        Args:
            content: Markdown file content

        Returns:
            List of capability strings
        """
        # Look for capabilities section
        capabilities_match = re.search(
            r'\*\*Capabilities:\*\*\s*\n((?:-\s+.+\s*\n?)+)',
            content,
            re.IGNORECASE
        )

        if not capabilities_match:
            return []

        # Extract list items
        capabilities_text = capabilities_match.group(1)
        capabilities = re.findall(r'^-\s+(.+)$', capabilities_text, re.MULTILINE)

        return [cap.strip() for cap in capabilities]
