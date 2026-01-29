"""Main CLI entry point for SDP package.

This module re-exports the CLI from the cli package for backward compatibility.
"""

from sdp.cli.main import main  # noqa: F401

__all__ = ["main"]


