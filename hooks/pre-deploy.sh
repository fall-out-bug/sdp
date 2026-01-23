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
cd "$(git rev-parse --show-toplevel)/tools/hw_checker"

# 1. E2E Tests (mandatory for deploy)
echo ""
echo "=== 1. E2E Tests (Critical Path) ==="
if [[ -d "tests/e2e" ]] && [[ $(ls tests/e2e/*.py 2>/dev/null | wc -l) -gt 0 ]]; then
  if poetry run pytest tests/e2e/ -v --tb=short -m "not slow"; then
    echo "✅ E2E tests passed"
  else
    echo "❌ E2E tests failed"
    echo ""
    echo "⛔ DEPLOY BLOCKED: E2E tests must pass"
    
    # Send notification
    FAILED_COUNT=$(poetry run pytest tests/e2e/ -v --tb=short -m "not slow" 2>&1 | grep -c "FAILED" || echo "unknown")
    bash ../../sdp/notifications/telegram.sh e2e_failed "$FEATURE_ID" "$FAILED_COUNT"
    
    exit 1
  fi
else
  echo "⚠️ No E2E tests found"
  echo "⚠️ Deploying without E2E coverage (risky!)"
fi

# 2. Smoke Tests
echo ""
echo "=== 2. Smoke Tests ==="
if [[ -f "tests/smoke/test_critical_path.py" ]]; then
  if poetry run pytest tests/smoke/ -v --tb=short; then
    echo "✅ Smoke tests passed"
  else
    echo "❌ Smoke tests failed"
    exit 1
  fi
else
  echo "⚠️ No smoke tests found (tests/smoke/)"
fi

# 3. Docker Build Test
echo ""
echo "=== 3. Docker Build Test ==="
if docker build -t hw-checker:pre-deploy-test . > /dev/null 2>&1; then
  echo "✅ Docker build successful"
  docker rmi hw-checker:pre-deploy-test > /dev/null 2>&1
else
  echo "❌ Docker build failed"
  exit 1
fi

# 4. Security Scan (if available)
echo ""
echo "=== 4. Security Scan ==="
if command -v bandit &> /dev/null; then
  if bandit -r src/hw_checker/ -ll -q; then
    echo "✅ No security issues found"
  else
    echo "⚠️ Security issues detected (review required)"
  fi
else
  echo "⚠️ bandit not installed, skipping security scan"
fi

# 5. Environment-specific checks
if [[ "$ENVIRONMENT" == "prod" ]]; then
  echo ""
  echo "=== 5. Production Readiness ==="
  
  # Check for debug flags
  if grep -r "DEBUG.*=.*True" src/ > /dev/null 2>&1; then
    echo "❌ DEBUG flags found in code"
    exit 1
  fi
  
  # Check for print statements (should use logging)
  if grep -r "print(" src/hw_checker/ | grep -v "# debug:" > /dev/null 2>&1; then
    echo "⚠️ print() statements found (should use logging)"
  fi
  
  echo "✅ Production readiness checks passed"
fi

echo ""
echo "✅ All pre-deploy checks passed for $FEATURE_ID"
echo "✅ Ready to deploy to $ENVIRONMENT"
