"""Extension system for SDP.

Allows project-specific customization via extensions:
- Custom hooks (validation, pre/post checks)
- Domain patterns documentation
- Custom skills (slash commands)
- Integrations (GitHub, GitLab, Telegram)
"""

from sdp.extensions.base import BaseExtension, Extension, ExtensionManifest
from sdp.extensions.loader import ExtensionLoader
from sdp.extensions.manifest import ManifestParser, ValidationError
from sdp.extensions.validator import ExtensionValidator

__all__ = [
    "BaseExtension",
    "Extension",
    "ExtensionManifest",
    "ExtensionLoader",
    "ExtensionValidator",
    "ManifestParser",
    "ValidationError",
]
