#!/bin/bash
# sdp/scripts/test_hook_enforcement.sh
# Test script to verify hooks enforce quality gates

set -e

echo "🧪 Testing Hook Enforcement"
echo "============================"

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

# Create temporary test directory
TEST_DIR=$(mktemp -d)
echo "Test directory: $TEST_DIR"

cleanup() {
    rm -rf "$TEST_DIR"
}
trap cleanup EXIT

# Test 1: except: pass should fail
echo ""
echo "Test 1: except: pass (should FAIL)"
cat > "$TEST_DIR/test_except_pass.py" << 'EOF'
def bad_function():
    try:
        risky_operation()
    except Exception:
        pass  # This should fail the hook
EOF

echo "Created test file with except: pass"
echo "Expected: Hook should fail and block commit"
echo ""

# Test 2: Large file (>200 LOC) should fail
echo ""
echo "Test 2: Large file >200 LOC (should FAIL)"
cat > "$TEST_DIR/test_large_file.py" << 'EOF'
# This file is intentionally large (>200 lines)

def function_1():
    pass

def function_2():
    pass

# ... (repeat to exceed 200 lines)
EOF

# Add 200 more lines
for i in {3..203}; do
    echo "def function_$i():" >> "$TEST_DIR/test_large_file.py"
    echo "    pass" >> "$TEST_DIR/test_large_file.py"
    echo "" >> "$TEST_DIR/test_large_file.py"
done

LINES=$(wc -l < "$TEST_DIR/test_large_file.py")
echo "Created test file with $LINES lines"
echo "Expected: Hook should fail and block commit"
echo ""

# Test 3: Missing type hints should fail
echo ""
echo "Test 3: Missing type hints (should FAIL)"
cat > "$TEST_DIR/test_no_hints.py" << 'EOF'
def add(a, b):  # Missing type hints
    return a + b
EOF

echo "Created test file without type hints"
echo "Expected: Hook should fail with mypy errors"
echo ""

# Test 4: Direct main commit should fail
echo ""
echo "Test 4: Commit to main branch (should FAIL)"
echo "Expected: Hook should fail with message about feature branches"
echo ""

# Summary
echo ""
echo "============================"
echo "Test Setup Complete"
echo ""
echo "To manually test hooks:"
echo "1. Copy test files to your working directory"
echo "2. Stage them: git add <test-files>"
echo "3. Try to commit: git commit -m 'test: violating code'"
echo "4. Hook should fail with actionable error messages"
echo ""
echo "Expected behavior:"
echo "  ✓ except: pass → FAILS with message to handle explicitly"
echo "  ✓ Large files → FAILS with message to split modules"
echo "  ✓ No type hints → FAILS with mypy errors"
echo "  ✓ Main commit → FAILS with feature branch suggestion"
echo ""
