#!/usr/bin/env python3
"""Quick documentation drift check."""
import sys
import re
from pathlib import Path

def check_frontmatter_consistency(ws_file):
    """Check if workstream frontmatter matches scope files."""
    with open(ws_file) as f:
        content = f.read()
    
    # Extract scope files
    match = re.search(r'scope_files:\s*\n((?:[ \t]- .+\n)+)', content)
    if not match:
        return 0  # No scope files defined
    
    scope_files = re.findall(r'- (.+)', match.group(1))
    
    # Check if files exist
    missing = []
    for file in scope_files:
        if not Path(file).exists():
            missing.append(file)
    
    if missing:
        print(f"❌ Scope files not found: {', '.join(missing)}")
        return 1
    
    print("✅ All scope files exist")
    return 0

if __name__ == "__main__":
    # Check all workstreams in backlog
    backlog_dir = Path("docs/workstreams/backlog")
    exit_code = 0
    
    for ws_file in backlog_dir.glob("*.md"):
        if check_frontmatter_consistency(ws_file) != 0:
            exit_code = 1
    
    sys.exit(exit_code)
