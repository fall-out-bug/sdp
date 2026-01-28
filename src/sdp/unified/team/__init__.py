"""TeamManager role registry for coordinated multi-agent workflows."""

from sdp.unified.team.errors import TeamManagerError
from sdp.unified.team.manager import TeamManager
from sdp.unified.team.models import Role, RoleState

__all__ = [
    "TeamManager",
    "Role",
    "RoleState",
    "TeamManagerError",
]
