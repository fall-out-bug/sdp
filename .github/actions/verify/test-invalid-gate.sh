#!/usr/bin/env bash
#
# Test script to validate AC1-AC4 for 00-058-05
# Tests that invalid gate handling works correctly

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Test helper functions
run_test() {
  local test_name="$1"
  local test_command="$2"

  TESTS_RUN=$((TESTS_RUN + 1))
  echo ""
  echo "Test $TESTS_RUN: $test_name"

  if eval "$test_command"; then
    echo -e "${GREEN}✅ PASS${NC}: $test_name"
    TESTS_PASSED=$((TESTS_PASSED + 1))
    return 0
  else
    echo -e "${RED}❌ FAIL${NC}: $test_name"
    TESTS_FAILED=$((TESTS_FAILED + 1))
    return 1
  fi
}

echo "==================================="
echo "Testing WS 00-058-05"
echo "Add continue-on-error to invalid gate test"
echo "==================================="

# AC1: Verify continue-on-error is present
run_test "AC1: continue-on-error is set on invalid gate test" \
  "grep -q 'continue-on-error: true' .github/workflows/test-verify-action.yml && grep -A5 'Test error handling - invalid gate' .github/workflows/test-verify-action.yml | grep -q 'continue-on-error: true'"

# AC2: Verify verification step exists
run_test "AC2: Verification step checks outputs.result" \
  "grep -A30 'Test error handling - invalid gate' .github/workflows/test-verify-action.yml | grep -q 'outputs.result'"

# AC2 continued: Verify it checks outcome
run_test "AC2: Verification step checks step outcome" \
  "grep -A30 'Test error handling - invalid gate' .github/workflows/test-verify-action.yml | grep -q 'steps.test-invalid.outcome'"

# AC4: Verify it checks gates_failed and gates_passed
run_test "AC4: Verification checks gates_failed" \
  "grep -A30 'Test error handling - invalid gate' .github/workflows/test-verify-action.yml | grep -q 'gates_failed'"

run_test "AC4: Verification checks gates_passed" \
  "grep -A30 'Test error handling - invalid gate' .github/workflows/test-verify-action.yml | grep -q 'gates_passed'"

# Verify structure: The invalid gate test should be at line 96-102
run_test "Invalid gate test has correct structure" \
  "sed -n '96,102p' .github/workflows/test-verify-action.yml | grep -q 'id: test-invalid'"

# Verify verification step structure (should be at line 104-137)
run_test "Verification step has correct structure" \
  "sed -n '104,137p' .github/workflows/test-verify-action.yml | grep -q 'Verify action fails on invalid gate'"

# Check that the verification validates the failure
run_test "Verification validates failure outcome" \
  "sed -n '104,137p' .github/workflows/test-verify-action.yml | grep -q 'outcome.*failure'"

# Check that outputs are validated
run_test "Verification validates outputs.result equals 'fail'" \
  "sed -n '104,137p' .github/workflows/test-verify-action.yml | grep -q \"outputs.result.*!= ''fail''\" || sed -n '104,137p' .github/workflows/test-verify-action.yml | grep -q 'result.*fail'"

# Check for proper error messages
run_test "Verification has proper error messages" \
  "sed -n '104,137p' .github/workflows/test-verify-action.yml | grep -q '❌ FAIL:'"

# Summary
echo ""
echo "==================================="
echo "Test Summary"
echo "==================================="
echo "Tests run: $TESTS_RUN"
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Failed: $TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
  echo -e "${GREEN}✅ All tests passed!${NC}"
  exit 0
else
  echo -e "${RED}❌ Some tests failed${NC}"
  exit 1
fi
