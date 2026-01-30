---
ws_id: 00-901-02
feature: P1 High Priority Tasks
status: completed
size: MEDIUM
github_issue: 5
title: P1-05: Make Pre-push Hard Blocking
goal: Implement hard blocking mode for pre-push hook to prevent bad code from reaching remote
acceptance_criteria:
  - [x] SDP_HARD_PUSH environment variable support
  - [x] Default: warning mode (backward compatible)
  - [x] Hard blocking: exit 1 on failures when SDP_HARD_PUSH=1
  - [x] Coverage < 80% blocks push
  - [x] Regression failures block push
  - [x] Actionable error messages with remediation steps
  - [x] Fixed for SDP repository (not hw_checker)
  - [x] Tests created (test_pre_push_hook.sh)
  - [x] Documentation created
context: |
  Problem: pre-push.sh had "Don't block push, just warn" for coverage and regression.
  This allowed bad code to reach remote repository.

  Solution: Add SDP_HARD_PUSH=1 environment variable for hard blocking mode.
  Default behavior unchanged (backward compatible).
steps: |
  1. Read hooks/pre-push.sh
  2. Added SDP_HARD_PUSH variable check
  3. Implemented hard blocking logic:
     - exit 1 on test failures when SDP_HARD_PUSH=1
     - exit 1 on coverage < 80% when SDP_HARD_PUSH=1
  4. Fixed for SDP (was hw_checker-specific):
     - Remove cd tools/hw_checker
     - Use pytest tests/ directly
     - Add command -v checks
  5. Created test_pre_push_hook.sh (8 tests, all passing)
  6. Created documentation
code_blocks: |
  # Check if strict mode is enabled
  HARD_PUSH=${SDP_HARD_PUSH:-0}

  if [ "$HARD_PUSH" = "1" ]; then
      echo "🔒 HARD PUSH MODE enabled"
      # Will exit 1 on failures
  else
      echo "⚠️ Soft mode (set SDP_HARD_PUSH=1)"
      # Will warn but not block
  fi

  # Later...
  if [ "$HARD_PUSH" = "1" ] && tests_fail; then
      echo "🚫 PUSH BLOCKED (SDP_HARD_PUSH=1)"
      echo "To bypass: git push --no-verify"
      exit 1
  fi
dependencies: []
execution_report: |
  **Duration:** 3 hours
  **LOC Added:** 171 (pre-push.sh + tests + docs)
  **LOC Modified:** 48 (pre-push.sh improvements)
  **Test Coverage:** 8/8 tests passing
  **Deviations:** None
  **Status:** ✅ COMPLETE

  Implemented SDP_HARD_PUSH=1 for hard blocking mode.
  Default: warning mode (backward compatible).
  Fixed pre-push.sh for SDP repository.
  Created comprehensive tests and documentation.

  Usage:
  # Warning mode (default)
  git push

  # Hard blocking mode
  export SDP_HARD_PUSH=1
  git push

  # Bypass if needed
  git push --no-verify
