#!/usr/bin/env python3
"""Quality gate validation script for git hooks.

This is a thin wrapper that imports the modular quality checking system.

The actual implementation is split into:
- scripts/quality/checker.py - Main orchestrator
- scripts/quality/security.py - Security checks
- scripts/quality/documentation.py - Documentation checks
- scripts/quality/performance.py - Performance checks
- scripts/quality/models.py - Data models

This wrapper maintains backward compatibility with existing hooks.
"""

import sys
from pathlib import Path

# Add repo root to path for imports
repo_root = Path(__file__).parent.parent
sys.path.insert(0, str(repo_root))

from scripts.quality.main import main

if __name__ == "__main__":
    sys.exit(main())
