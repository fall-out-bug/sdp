#!/bin/bash
# Watch mode for pytest using pytest-watcher
# Usage: bash scripts/watch-tests.sh
#
# This script runs pytest in watch mode, automatically re-running tests
# whenever Python files change. Perfect for TDD workflow.
#
# Features:
# - Ignores integration tests by default (fast feedback)
# - Clears screen between runs for clean output
# - Passes additional args to pytest
# - Uses 0.2s delay to wait for file saves/formatters
#
# Examples:
#   scripts/watch-tests.sh                    # Watch all unit tests
#   scripts/watch-tests.sh --ignore=tests/unit/adapters  # Exclude specific path
#   scripts/watch-tests.sh -k test_specific    # Run specific tests only
#   scripts/watch-tests.sh --cov              # Run with coverage

set -e

# Default arguments (pytest-watcher options)
WATCHER_ARGS="--now --clear --delay 0.2"
PYTEST_ARGS="--ignore=tests/integration -v"

# Parse additional arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --with-integration)
      PYTEST_ARGS="${PYTEST_ARGS/--ignore=tests\/integration/}"
      shift
      ;;
    --cov|--coverage)
      PYTEST_ARGS="$PYTEST_ARGS --cov=src/sdp --cov-report=term-missing"
      shift
      ;;
    -k)
      PYTEST_ARGS="$PYTEST_ARGS -k $2"
      shift 2
      ;;
    *)
      PYTEST_ARGS="$PYTEST_ARGS $1"
      shift
      ;;
  esac
done

echo "🔍 Starting pytest-watcher..."
echo "   Watching: . (recursive)"
echo "   Ignoring: tests/integration"
echo "   Args: $PYTEST_ARGS"
echo ""
echo "Tip: Press Ctrl+C to stop"
echo ""

# Run pytest-watcher
# Note: pytest-watcher passes all args after path to pytest
poetry run pytest-watcher . $WATCHER_ARGS -- $PYTEST_ARGS
