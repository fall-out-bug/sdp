#!/bin/bash
# Quality metrics checker for SDP
# Checks: 1) coverage >= 80%  2) test/code ratio >= 1.5

set -e

echo "=== Quality Metrics Check ==="
echo ""

# Coverage check
echo "1. Coverage Check (target: >= 80%)"
echo "--------------------------------"
FAILED_COVERAGE=""

for pkg in $(go list ./src/sdp/... ./sdp-plugin/internal/... 2>/dev/null | grep -v "/test$" | sort); do
    coverage=$(go test -cover "$pkg" 2>/dev/null | grep -oP 'coverage:\s*\K[0-9.]+')
    if [ -n "$coverage" ]; then
        if (( $(echo "$coverage < 80.0" | bc -l) )); then
            echo "  ❌ $pkg: ${coverage}%"
            FAILED_COVERAGE="$FAILED_COVERAGE $pkg"
        else
            echo "  ✅ $pkg: ${coverage}%"
        fi
    fi
done

echo ""
echo "2. Test/Code Ratio Check (target: 1.5 - 2.0)"
echo "---------------------------------------------"

check_ratio() {
    local pkg_path=$1
    local pkg_name=$2

    # Count production code lines (exclude test files)
    prod_lines=$(find "$pkg_path" -name "*.go" ! -name "*_test.go" -exec cat {} \; 2>/dev/null | grep -v "^$" | grep -v "^//" | wc -l | tr -d ' ')

    # Count test code lines
    test_lines=$(find "$pkg_path" -name "*_test.go" -exec cat {} \; 2>/dev/null | grep -v "^$" | grep -v "^//" | wc -l | tr -d ' ')

    if [ "$prod_lines" -gt 0 ]; then
        ratio=$(echo "scale=2; $test_lines / $prod_lines" | bc)
        if (( $(echo "$ratio < 1.5" | bc -l) )); then
            echo "  ❌ $pkg_name: ${ratio} (${test_lines}/${prod_lines} lines) - BELOW MINIMUM"
            return 1
        elif (( $(echo "$ratio > 2.0" | bc -l) )); then
            echo "  ⚠️  $pkg_name: ${ratio} (${test_lines}/${prod_lines} lines) - ABOVE MAXIMUM"
            return 0  # Warning only, not a failure
        else
            echo "  ✅ $pkg_name: ${ratio} (${test_lines}/${prod_lines} lines)"
            return 0
        fi
    fi
    return 0
}

FAILED_RATIO=""
for dir in src/sdp/*/ sdp-plugin/internal/*/; do
    if [ -d "$dir" ]; then
        pkg_name=$(basename "$dir")
        if ! check_ratio "$dir" "$pkg_name"; then
            FAILED_RATIO="$FAILED_RATIO $pkg_name"
        fi
    fi
done

echo ""
echo "=== Summary ==="
if [ -z "$FAILED_COVERAGE" ] && [ -z "$FAILED_RATIO" ]; then
    echo "✅ All quality metrics passed!"
    exit 0
else
    [ -n "$FAILED_COVERAGE" ] && echo "❌ Coverage failures:$FAILED_COVERAGE"
    [ -n "$FAILED_RATIO" ] && echo "❌ Test/code ratio failures:$FAILED_RATIO"
    exit 1
fi
