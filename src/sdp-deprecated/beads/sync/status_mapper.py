"""
Status and priority mapping between SDP and Beads.

Handles conversion of status values and size/priority mappings.
"""

from ..models import BeadsPriority, BeadsStatus

# SDP workstream status values
SDP_STATUS_BACKLOG = "backlog"
SDP_STATUS_ACTIVE = "active"
SDP_STATUS_COMPLETED = "completed"
SDP_STATUS_BLOCKED = "blocked"


def map_sdp_status_to_beads(sdp_status: str | object) -> BeadsStatus:
    """Map SDP status to Beads status.

    Args:
        sdp_status: SDP status value (string or enum)

    Returns:
        Corresponding BeadsStatus
    """
    status_str = getattr(sdp_status, "value", str(sdp_status))
    mapping = {
        SDP_STATUS_BACKLOG: BeadsStatus.OPEN,
        SDP_STATUS_ACTIVE: BeadsStatus.IN_PROGRESS,
        SDP_STATUS_COMPLETED: BeadsStatus.CLOSED,
        SDP_STATUS_BLOCKED: BeadsStatus.BLOCKED,
    }
    return mapping.get(status_str, BeadsStatus.OPEN)


def map_beads_status_to_sdp(beads_status: BeadsStatus) -> str:
    """Map Beads status to SDP status.

    Args:
        beads_status: Beads status value

    Returns:
        Corresponding SDP status string
    """
    mapping = {
        BeadsStatus.OPEN: SDP_STATUS_BACKLOG,
        BeadsStatus.IN_PROGRESS: SDP_STATUS_ACTIVE,
        BeadsStatus.CLOSED: SDP_STATUS_COMPLETED,
        BeadsStatus.BLOCKED: SDP_STATUS_BLOCKED,
    }
    return mapping.get(beads_status, SDP_STATUS_BACKLOG)


def map_sdp_size_to_beads_priority(size: str | object) -> BeadsPriority:
    """Map SDP size to Beads priority.

    Args:
        size: SDP size value (string or enum)

    Returns:
        Corresponding BeadsPriority

    Note:
        Larger workstream = higher priority (lower number = more critical)
    """
    size_str = getattr(size, "value", str(size)) if size else "MEDIUM"
    mapping = {
        "SMALL": BeadsPriority(2),
        "MEDIUM": BeadsPriority(1),
        "LARGE": BeadsPriority(0),
    }
    return mapping.get(str(size_str).upper(), BeadsPriority(2))
