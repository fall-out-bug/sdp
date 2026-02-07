"""TeamManager role registry for coordinated multi-agent workflows."""

from sdp.unified.team.errors import TeamManagerError
from sdp.unified.team.lifecycle import create_team, delete_team, get_team
from sdp.unified.team.manager import TeamManager
from sdp.unified.team.models import Role, RoleState

__all__ = [
    "TeamManager",
    "Role",
    "RoleState",
    "TeamManagerError",
    "create_team",
    "delete_team",
    "get_team",
]
