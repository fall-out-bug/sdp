#!/usr/bin/env python3
"""Check that Python files do not exceed size limits.

This script ensures all Python files in src/ are under 200 LOC
to maintain AI-readability and code quality.
"""

import argparse
import json
import sys
from pathlib import Path


def count_loc(file_path: Path) -> int:
    """Count lines of code in a Python file.

    Args:
        file_path: Path to Python file

    Returns:
        Number of non-blank, non-comment lines
    """
    lines = file_path.read_text(encoding='utf-8').splitlines()
    loc = 0

    for line in lines:
        stripped = line.strip()
        # Skip empty lines and comments
        if stripped and not stripped.startswith('#'):
            loc += 1

    return loc


def check_file_sizes(max_loc: int = 200) -> dict[str, list[dict[str, int | str]]]:
    """Check all Python files in src/ for size violations.

    Args:
        max_loc: Maximum allowed lines of code per file

    Returns:
        Dictionary with count and violations list
    """
    src_path = Path('src/sdp')
    violations = []

    for py_file in src_path.rglob('*.py'):
        loc = count_loc(py_file)

        if loc > max_loc:
            violations.append({
                'file': str(py_file),
                'lines': loc,
                'max': max_loc
            })

    return {
        'count': len(violations),
        'violations': sorted(violations, key=lambda x: x['lines'], reverse=True)
    }


if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        description='Check Python file sizes against LOC limits'
    )
    parser.add_argument(
        '--json',
        action='store_true',
        help='Output results as JSON'
    )
    parser.add_argument(
        '--max-loc',
        type=int,
        default=200,
        help='Maximum lines of code per file (default: 200)'
    )
    args = parser.parse_args()

    result = check_file_sizes(max_loc=args.max_loc)

    if args.json:
        print(json.dumps(result, indent=2))
    else:
        if result['count'] > 0:
            print(f"❌ Found {result['count']} file(s) exceeding {args.max_loc} LOC limit:", file=sys.stderr)
            for v in result['violations']:
                print(f"  {v['file']}: {v['lines']} LOC", file=sys.stderr)
            sys.exit(1)
        else:
            print(f"✅ All files under {args.max_loc} LOC")

    sys.exit(1 if result['count'] > 0 else 0)
