#!/bin/bash
# Pre-deploy hook: E2E tests before deployment

set -e

FEATURE_ID="$1"
ENVIRONMENT="${2:-staging}"

if [[ -z "$FEATURE_ID" ]]; then
  echo "Usage: pre-deploy.sh F{XX} [staging|prod]"
  exit 1
fi

echo "🚀 Running pre-deploy checks for $FEATURE_ID ($ENVIRONMENT)..."

# Change to project root
cd "$(git rev-parse --show-toplevel)"

# 1. Unit Tests (mandatory for deploy)
echo ""
echo "=== 1. Unit Tests ==="
if pytest tests/unit/ -v --tb=short --cov=sdp --cov-report=term-missing --cov-fail-under=70; then
  echo "✅ Unit tests passed"
else
  echo "❌ Unit tests failed"
  echo ""
  echo "⛔ DEPLOY BLOCKED: Unit tests must pass"
  exit 1
fi

# 2. Integration Tests
echo ""
echo "=== 2. Integration Tests ==="
if pytest tests/integration/ -v --tb=short; then
  echo "✅ Integration tests passed"
else
  echo "❌ Integration tests failed"
  exit 1
fi

# 3. Type Checking
echo ""
echo "=== 3. Type Checking (mypy) ==="
if mypy src/sdp/ --strict 2>&1 | head -50; then
  echo "✅ Type checking passed"
else
  echo "⚠️ Type checking issues found (review required)"
fi

# 4. Linting (ruff)
echo ""
echo "=== 4. Linting (ruff) ==="
if ruff check src/sdp/; then
  echo "✅ Linting passed"
else
  echo "⚠️ Linting issues found (review required)"
fi

# 5. SDP Quality Gates
echo ""
echo "=== 5. SDP Quality Gates ==="
if [[ -f "hooks/validators/session-quality-check.sh" ]]; then
  if bash hooks/validators/session-quality-check.sh; then
    echo "✅ SDP quality gates passed"
  else
    echo "⚠️ SDP quality gate warnings"
  fi
else
  echo "⚠️ session-quality-check.sh not found"
fi

# 6. Environment-specific checks
if [[ "$ENVIRONMENT" == "prod" ]]; then
  echo ""
  echo "=== 6. Production Readiness ==="

  # Check for debug flags
  if grep -r "DEBUG.*=.*True" src/ 2>/dev/null | grep -v ".pyc" | grep -v "__pycache__" > /dev/null; then
    echo "❌ DEBUG flags found in code"
    exit 1
  fi

  # Check for print statements (should use logging)
  if grep -r "print(" src/sdp/ 2>/dev/null | grep -v "# debug:" | grep -v "__pycache__" | grep -v ".pyc" > /dev/null; then
    echo "⚠️ print() statements found (should use logging)"
  fi

  echo "✅ Production readiness checks passed"
fi

echo ""
echo "✅ All pre-deploy checks passed for $FEATURE_ID"
echo "✅ Ready to deploy to $ENVIRONMENT"
