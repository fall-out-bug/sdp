"""Model mapping package.

This module provides backward compatibility by re-exporting all public APIs.
"""

from sdp.core.model.loader import load_model_registry
from sdp.core.model.models import ModelMappingError, ModelProvider, ModelRegistry

__all__ = [
    "ModelProvider",
    "ModelRegistry",
    "ModelMappingError",
    "load_model_registry",
]
