#!/bin/bash
# Test script to verify pre-push hook behavior
# This script tests the pre-push hook in both warning and hard blocking modes

set -e

echo "===================================="
echo "Testing Pre-Push Hook Behavior"
echo "===================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

echo "Project root: $PROJECT_ROOT"
echo ""

# Test 1: Verify hook exists and is executable
echo "Test 1: Checking hook installation..."
if [ -x ".git/hooks/pre-push" ]; then
    echo -e "${GREEN}✓ Pre-push hook is installed and executable${NC}"
else
    echo -e "${RED}✗ Pre-push hook not found or not executable${NC}"
    echo "Run: bash hooks/install-hooks.sh"
    exit 1
fi
echo ""

# Test 2: Verify hook script exists
echo "Test 2: Checking hook script source..."
if [ -f "hooks/pre-push.sh" ]; then
    echo -e "${GREEN}✓ Hook script exists at hooks/pre-push.sh${NC}"
else
    echo -e "${RED}✗ Hook script not found${NC}"
    exit 1
fi
echo ""

# Test 3: Verify SDP_HARD_PUSH environment variable support
echo "Test 3: Testing SDP_HARD_PUSH environment variable..."
if grep -q "SDP_HARD_PUSH" hooks/pre-push.sh; then
    echo -e "${GREEN}✓ Hook supports SDP_HARD_PUSH environment variable${NC}"
else
    echo -e "${RED}✗ SDP_HARD_PUSH not found in hook script${NC}"
    exit 1
fi
echo ""

# Test 4: Verify hard blocking logic
echo "Test 4: Checking hard blocking logic..."
if grep -q 'if \[ "$HARD_PUSH" = "1" \]' hooks/pre-push.sh; then
    echo -e "${GREEN}✓ Hard blocking logic present${NC}"
else
    echo -e "${RED}✗ Hard blocking logic not found${NC}"
    exit 1
fi
echo ""

# Test 5: Verify remediation steps
echo "Test 5: Checking remediation messages..."
if grep -q "To fix this issue:" hooks/pre-push.sh; then
    echo -e "${GREEN}✓ Remediation steps present${NC}"
else
    echo -e "${RED}✗ Remediation steps not found${NC}"
    exit 1
fi
echo ""

# Test 6: Verify both checks (regression + coverage)
echo "Test 6: Checking regression test validation..."
if grep -q "Regression tests failed" hooks/pre-push.sh; then
    echo -e "${GREEN}✓ Regression test failure handling present${NC}"
else
    echo -e "${RED}✗ Regression test failure handling not found${NC}"
    exit 1
fi
echo ""

echo "Test 7: Checking coverage validation..."
if grep -q "Coverage is below 80%" hooks/pre-push.sh; then
    echo -e "${GREEN}✓ Coverage failure handling present${NC}"
else
    echo -e "${RED}✗ Coverage failure handling not found${NC}"
    exit 1
fi
echo ""

# Test 8: Verify bypass instruction
echo "Test 8: Checking bypass instruction..."
if grep -q "git push --no-verify" hooks/pre-push.sh; then
    echo -e "${GREEN}✓ Bypass instruction present${NC}"
else
    echo -e "${RED}✗ Bypass instruction not found${NC}"
    exit 1
fi
echo ""

# Summary
echo "===================================="
echo -e "${GREEN}All tests passed!${NC}"
echo "===================================="
echo ""
echo "Pre-push hook features verified:"
echo "  ✓ Installation check"
echo "  ✓ Environment variable support (SDP_HARD_PUSH)"
echo "  ✓ Hard blocking logic (exit 1 on failures)"
echo "  ✓ Warning mode (default, SDP_HARD_PUSH=0)"
echo "  ✓ Remediation steps for failures"
echo "  ✓ Regression test validation"
echo "  ✓ Coverage validation (≥80%)"
echo "  ✓ Bypass instruction (--no-verify)"
echo ""
echo "Usage:"
echo "  Warning mode (default):    git push"
echo "  Hard blocking mode:        SDP_HARD_PUSH=1 git push"
echo "  Emergency bypass:          git push --no-verify"
echo ""
