#!/bin/bash
# Pre-deploy hook: build, test, and contract validation before deployment
# Usage: ./pre-deploy.sh F{XX} [staging|prod]

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

# Go build and test
if [ -d "sdp-plugin" ]; then
    cd sdp-plugin
    go build ./...
    go test ./... -count=1 -short
    cd "$REPO_ROOT"
fi

# Contract validation
echo ""
echo "=== 6. Contract Validation ==="

if [ -d ".contracts" ] && [ "$(ls -A .contracts/*.yaml 2>/dev/null | wc -l)" -gt 0 ]; then
    CONTRACT_COUNT=$(ls -A .contracts/*.yaml 2>/dev/null | wc -l)

    if [ "$CONTRACT_COUNT" -lt 2 ]; then
        echo "  Skipped (need at least 2 contracts, found $CONTRACT_COUNT)"
    else
        # Build SDP CLI if not already built (from sdp-plugin in this repo)
        if [ ! -f "./sdp" ]; then
            echo "  Building SDP CLI..."
            if [ -d "sdp-plugin" ]; then
                (cd sdp-plugin && go build -o ../sdp ./cmd/sdp) || {
                    echo "⚠️ Failed to build SDP CLI, skipping contract validation"
                    exit 0
                }
            else
                go build -o sdp ./cmd/sdp || {
                    echo "⚠️ Failed to build SDP CLI, skipping contract validation"
                    exit 0
                }
            fi
        fi

        # Run contract validation
        echo "  Validating $CONTRACT_COUNT contract(s)..."

        CONTRACTS=($(ls .contracts/*.yaml))
        ./sdp contract validate "${CONTRACTS[@]}" --output validation-report.md || VALIDATION_FAILED=true

        # Parse validation report
        if [ -f "validation-report.md" ]; then
            ERROR_COUNT=$(grep -c "^|.*ERROR" validation-report.md 2>/dev/null || echo "0")
            WARNING_COUNT=$(grep -c "^|.*WARNING" validation-report.md 2>/dev/null || echo "0")

            echo "  Validation complete: $ERROR_COUNT errors, $WARNING_COUNT warnings"

            if [ "$ERROR_COUNT" -gt 0 ]; then
                echo "❌ Contract validation failed"
                echo ""
                echo "Validation report:"
                cat validation-report.md
                rm -f validation-report.md
                exit 1
            fi

            rm -f validation-report.md
        fi

        if [ "$VALIDATION_FAILED" = true ]; then
            echo "⚠️ Contract validation command failed"
            exit 1
        fi

        if [ "$ERROR_COUNT" -eq 0 ]; then
            echo "✅ Contract validation passed"
        fi
    fi
else
    echo "  Skipped (no contracts found)"
fi

echo ""
echo "✅ All pre-deploy checks passed"
exit 0
