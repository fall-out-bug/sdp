"""Extension manifest schema and validation.

Extension manifest (extension.yaml) structure:

```yaml
name: hw_checker
version: 1.0.0
description: Clean Architecture validation for hw_checker
author: SDP Team

# Optional directory overrides (defaults shown)
hooks_dir: hooks
patterns_dir: patterns
skills_dir: skills
integrations_dir: integrations
```

Example extension structure:
```
~/.sdp/extensions/hw_checker/
├── extension.yaml          # Manifest
├── hooks/                  # Validation scripts
│   ├── pre-build.sh
│   └── post-build.sh
├── patterns/               # Domain patterns
│   └── clean_architecture.md
├── skills/                 # Custom commands
│   └── build/
│       └── SKILL.md
└── integrations/           # Service configs
    └── github/
        └── config.yaml
```
"""

from pathlib import Path
from typing import Any

import yaml

from sdp.extensions.base import ExtensionManifest


class ValidationError(Exception):
    """Extension validation error."""
    
    pass


class ManifestParser:
    """Parse and validate extension.yaml manifest.
    
    Example:
        >>> parser = ManifestParser()
        >>> manifest = parser.parse_file(Path("extension.yaml"))
        >>> print(manifest.name)
        'hw_checker'
    """
    
    REQUIRED_FIELDS: tuple[str, ...] = ("name", "version", "description", "author")
    
    def parse_file(self, manifest_path: Path) -> ExtensionManifest:
        """Parse manifest from YAML file.
        
        Args:
            manifest_path: Path to extension.yaml
        
        Returns:
            Parsed manifest
        
        Raises:
            ValidationError: If manifest is invalid
        
        Example:
            >>> parser = ManifestParser()
            >>> manifest = parser.parse_file(Path("/path/to/extension.yaml"))
        """
        if not manifest_path.exists():
            raise ValidationError(f"Manifest not found: {manifest_path}")
        
        try:
            with open(manifest_path, "r", encoding="utf-8") as f:
                data = yaml.safe_load(f)
        except yaml.YAMLError as e:
            raise ValidationError(f"Invalid YAML in {manifest_path}: {e}")
        
        return self.parse_dict(data, manifest_path)
    
    def parse_dict(self, data: dict[str, Any], source: Path) -> ExtensionManifest:
        """Parse manifest from dictionary.
        
        Args:
            data: Manifest data
            source: Source file path (for error messages)
        
        Returns:
            Parsed manifest
        
        Raises:
            ValidationError: If manifest is invalid
        
        Example:
            >>> parser = ManifestParser()
            >>> manifest = parser.parse_dict({"name": "test", ...}, Path("test.yaml"))
        """
        self._validate_required_fields(data, source)
        
        return ExtensionManifest(
            name=data["name"],
            version=data["version"],
            description=data["description"],
            author=data["author"],
            hooks_dir=data.get("hooks_dir", "hooks"),
            patterns_dir=data.get("patterns_dir", "patterns"),
            skills_dir=data.get("skills_dir", "skills"),
            integrations_dir=data.get("integrations_dir", "integrations"),
        )
    
    def _validate_required_fields(self, data: dict[str, Any], source: Path) -> None:
        """Validate required fields are present.
        
        Args:
            data: Manifest data
            source: Source file path (for error messages)
        
        Raises:
            ValidationError: If required fields are missing
        """
        missing = [field for field in self.REQUIRED_FIELDS if field not in data]
        if missing:
            raise ValidationError(
                f"Missing required fields in {source}: {', '.join(missing)}"
            )
