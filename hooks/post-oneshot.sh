#!/bin/bash
# Post-oneshot hook: Integration & E2E tests after /oneshot completion

set -e

FEATURE_ID="$1"

if [[ -z "$FEATURE_ID" ]]; then
  echo "Usage: post-oneshot.sh F{XX}"
  exit 1
fi

echo "🧪 Running post-oneshot checks for $FEATURE_ID..."

# Change to project root
cd "$(git rev-parse --show-toplevel)/tools/hw_checker"

# 1. Integration Tests
echo ""
echo "=== 1. Integration Tests ==="
if poetry run pytest tests/integration/ -v --tb=short; then
  echo "✅ Integration tests passed"
else
  echo "❌ Integration tests failed"
  exit 1
fi

# 2. E2E Tests (if exist)
if [[ -d "tests/e2e" ]] && [[ $(ls tests/e2e/*.py 2>/dev/null | wc -l) -gt 0 ]]; then
  echo ""
  echo "=== 2. E2E Tests ==="
  if poetry run pytest tests/e2e/ -v --tb=short; then
    echo "✅ E2E tests passed"
  else
    echo "❌ E2E tests failed"
    exit 1
  fi
else
  echo ""
  echo "⚠️ No E2E tests found (tests/e2e/)"
fi

# 3. Full regression
echo ""
echo "=== 3. Full Regression Suite ==="
if poetry run pytest tests/unit/ -m fast -q --tb=short; then
  echo "✅ Regression suite passed"
else
  echo "❌ Regression suite failed"
  exit 1
fi

# 4. Coverage check (entire codebase)
echo ""
echo "=== 4. Overall Coverage Check ==="
COVERAGE=$(poetry run pytest tests/ --cov=hw_checker --cov-report=term-missing --cov-fail-under=80 -q | grep "TOTAL" | awk '{print $4}')
echo "Overall coverage: $COVERAGE"

echo ""
echo "✅ All post-oneshot checks passed for $FEATURE_ID"
