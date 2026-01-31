"""Beads client exceptions."""


class BeadsClientError(Exception):
    """Beads client operation error.

    Raised when Beads CLI operations fail.
    """

    pass


class BeadsNotInstalledError(BeadsClientError):
    """Beads CLI is not installed or not found in PATH.

    Raised when attempting to use real Beads client without installation.
    """

    pass
