"""GitHub configuration.

Loads GitHub API credentials from environment variables.
Token is never logged for security.
"""
import os
from dataclasses import dataclass
from pathlib import Path

from dotenv import load_dotenv


@dataclass(frozen=True)
class GitHubConfig:
    """GitHub API configuration.

    Attributes:
        token: GitHub personal access token (never logged)
        repo: Repository in format "owner/repo"
        org: Optional organization name
    """

    token: str
    repo: str
    org: str | None = None

    @classmethod
    def from_env(cls, env_file: Path | None = None) -> "GitHubConfig":
        """Load config from .env file.

        If .env doesn't exist but env vars are set, creates .env automatically.

        Args:
            env_file: Path to .env file (default: .env in current dir)

        Returns:
            GitHubConfig instance

        Raises:
            ValueError: If required env vars missing
        """
        # Determine .env file path
        if env_file:
            env_path = env_file
        else:
            # Default: .env in sdp/ directory (where this code runs)
            # Path from sdp/src/sdp/github/config.py -> sdp/.env
            # __file__ = sdp/src/sdp/github/config.py
            # Need: sdp/.env (4 levels up from __file__)
            env_path = Path(__file__).parent.parent.parent.parent / ".env"

        # Check process env vars first (before loading .env)
        # This allows agent to auto-create .env from user's env vars
        # Treat empty strings as None (unset)
        process_token = os.getenv("GITHUB_TOKEN")
        process_repo = os.getenv("GITHUB_REPO")
        if process_token == "":
            process_token = None
        if process_repo == "":
            process_repo = None

        # Always load .env if it exists (it's the source of truth after first creation)
        if env_path.exists():
            load_dotenv(env_path, override=True)

        # Get final values: process env (if non-empty) takes precedence over .env
        token = process_token if process_token else os.getenv("GITHUB_TOKEN")
        repo = process_repo if process_repo else os.getenv("GITHUB_REPO")

        # If we have env vars in process but no .env file, create it
        # This allows agent to work: user sets GITHUB_TOKEN once, agent creates .env
        if process_token and process_repo and not env_path.exists():
            try:
                env_path.parent.mkdir(parents=True, exist_ok=True)
                with env_path.open("w") as f:
                    f.write(f"GITHUB_TOKEN={process_token}\n")
                    f.write(f"GITHUB_REPO={process_repo}\n")
                    if process_org := os.getenv("GITHUB_ORG"):
                        f.write(f"GITHUB_ORG={process_org}\n")
                # Set restrictive permissions (owner read/write only)
                env_path.chmod(0o600)
                # Reload .env we just created
                load_dotenv(env_path)
                token = process_token
                repo = process_repo
            except (OSError, PermissionError):
                # Can't create .env, continue with env vars from process
                pass

        if not token:
            raise ValueError("GITHUB_TOKEN not found in environment")
        if not repo:
            raise ValueError("GITHUB_REPO not found in environment")

        org = os.getenv("GITHUB_ORG")

        return cls(token=token, repo=repo, org=org)

    def __repr__(self) -> str:
        """Safe repr that hides token."""
        return f"GitHubConfig(token='***', repo='{self.repo}', org={self.org!r})"
