"""DEPRECATED: Use sdp.core.model submodule instead.

This module provides backward compatibility by re-exporting from the model package.
"""

import warnings

from sdp.core.model import ModelMappingError, ModelProvider, ModelRegistry, load_model_registry

warnings.warn(
    "sdp.core.model_mapping module is deprecated. "
    "Use 'from sdp.core.model import ...' instead.",
    DeprecationWarning,
    stacklevel=2,
)

__all__ = [
    "ModelProvider",
    "ModelRegistry",
    "ModelMappingError",
    "load_model_registry",
]
