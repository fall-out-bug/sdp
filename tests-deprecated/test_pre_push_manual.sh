#!/bin/bash
# Manual test script to demonstrate pre-push hook behavior
# This script creates intentional failures to show how the hook behaves

set -e

echo "=========================================="
echo "Manual Pre-Push Hook Behavior Test"
echo "=========================================="
echo ""

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "This script will demonstrate the pre-push hook behavior"
echo ""
echo "Test scenarios:"
echo "  1. Test with SDP_HARD_PUSH=0 (warning mode, default)"
echo "  2. Test with SDP_HARD_PUSH=1 (hard blocking mode)"
echo ""
echo "Note: This requires actual test failures to demonstrate"
echo "      The hook will only fail if tests fail or coverage < 80%"
echo ""

# Check if we can run the tests
echo "Checking if test environment is ready..."
if [ ! -d "tools/hw_checker" ]; then
    echo -e "${YELLOW}⚠️  tools/hw_checker directory not found${NC}"
    echo "The pre-push hook expects tests in tools/hw_checker/tests/unit/"
    echo ""
    echo "Skipping actual test execution - showing behavior documentation instead:"
    echo ""
    exit 0
fi

echo -e "${GREEN}✓ Test environment found${NC}"
echo ""

# Test 1: Warning mode
echo "=========================================="
echo "Test 1: Warning Mode (SDP_HARD_PUSH=0)"
echo "=========================================="
echo ""
echo "Behavior:"
echo "  - Runs regression tests"
echo "  - Checks coverage >= 80%"
echo "  - Shows warnings for failures"
echo "  - Does NOT block push"
echo "  - Exit code: 0 (always)"
echo ""
echo "Expected output on test failure:"
echo "  ❌ Regression tests failed"
echo ""
echo "  To fix this issue:"
echo "    1. Run: cd tools/hw_checker && poetry run pytest tests/unit/ -m fast -v"
echo "    2. Fix failing tests"
echo "    3. Commit the fixes"
echo "    4. Push again"
echo ""
echo "  ⚠️  WARNING: Push not blocked (SDP_HARD_PUSH not set)"
echo "     Set SDP_HARD_PUSH=1 to enforce blocking in future"
echo ""

# Test 2: Hard blocking mode
echo "=========================================="
echo "Test 2: Hard Blocking Mode (SDP_HARD_PUSH=1)"
echo "=========================================="
echo ""
echo "Behavior:"
echo "  - Runs regression tests"
echo "  - Checks coverage >= 80%"
echo "  - Blocks push on failures"
echo "  - Exit code: 1 (on failures)"
echo ""
echo "Expected output on test failure:"
echo "  ❌ Regression tests failed"
echo ""
echo "  To fix this issue:"
echo "    1. Run: cd tools/hw_checker && poetry run pytest tests/unit/ -m fast -v"
echo "    2. Fix failing tests"
echo "    3. Commit the fixes"
echo "    4. Push again"
echo ""
echo "  🚫 PUSH BLOCKED (SDP_HARD_PUSH=1)"
echo "     To bypass: git push --no-verify"
echo ""

# Test 3: Coverage failure
echo "=========================================="
echo "Test 3: Coverage Failure (< 80%)"
echo "=========================================="
echo ""
echo "Behavior (both modes):"
echo "  - Checks .coverage file if exists"
echo "  - Extracts coverage percentage"
echo "  - Warns or blocks based on SDP_HARD_PUSH"
echo ""
echo "Expected output on coverage < 80%:"
echo "  ❌ Coverage is below 80% (currently: 65%)"
echo ""
echo "  To fix this issue:"
echo "    1. Run: cd tools/hw_checker && poetry run pytest --cov=. --cov-report=term-missing"
echo "    2. Add tests for uncovered lines"
echo "    3. Commit the tests"
echo "    4. Push again"
echo ""
if [ "$SDP_HARD_PUSH" = "1" ]; then
    echo "  🚫 PUSH BLOCKED (SDP_HARD_PUSH=1)"
else
    echo "  ⚠️  WARNING: Push not blocked (SDP_HARD_PUSH not set)"
fi
echo ""

# Summary
echo "=========================================="
echo "Summary"
echo "=========================================="
echo ""
echo "The pre-push hook now supports two modes:"
echo ""
echo -e "${GREEN}1. Warning Mode (default)${NC}"
echo "   - SDP_HARD_PUSH=0 or unset"
echo "   - Shows warnings but allows push"
echo "   - Good for development phase"
echo ""
echo -e "${RED}2. Hard Blocking Mode${NC}"
echo "   - SDP_HARD_PUSH=1"
echo "   - Blocks push on failures"
echo "   - Enforces quality standards"
echo ""
echo "Usage examples:"
echo ""
echo "  # Enable hard blocking for current session"
echo "  export SDP_HARD_PUSH=1"
echo ""
echo "  # Push with hard blocking"
echo "  git push"
echo ""
echo "  # Emergency bypass (use sparingly)"
echo "  git push --no-verify"
echo ""
echo "Configuration files:"
echo "  - Hook script: hooks/pre-push.sh"
echo "  - Install script: hooks/install-hooks.sh"
echo "  - Documentation: docs/runbooks/git-hooks-installation.md"
echo ""
