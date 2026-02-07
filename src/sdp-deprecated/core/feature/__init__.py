"""Feature decomposition and dependency management package.

This module provides backward compatibility by re-exporting all public APIs.
"""

from sdp.core.feature.errors import CircularDependencyError, MissingDependencyError
from sdp.core.feature.loader import load_feature_from_directory, load_feature_from_spec
from sdp.core.feature.models import Feature

__all__ = [
    "CircularDependencyError",
    "MissingDependencyError",
    "Feature",
    "load_feature_from_directory",
    "load_feature_from_spec",
]
