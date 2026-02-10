#!/usr/bin/env bash
# Unit tests for GitHub Action
# Run with: bash .github/actions/verify/test-comment.sh

echo "🧪 Running unit tests for GitHub Action..."
echo "==========================================="
echo ""

TESTS_PASSED=0
TESTS_FAILED=0

# Test 1: Verify action.yml exists
test_action_yml_exists() {
    [ -f ".github/actions/verify/action.yml" ]
}

# Test 2: Verify comment.sh exists
test_comment_sh_exists() {
    [ -f ".github/actions/verify/comment.sh" ]
}

# Test 3: Verify comment.sh is executable
test_comment_sh_executable() {
    [ -x ".github/actions/verify/comment.sh" ]
}

# Test 4: Verify action has no sudo (except comments)
test_no_sudo() {
    # Check for sudo in actual code (not comments)
    ! grep -v "^[[:space:]]*#" .github/actions/verify/action.yml | grep -q "sudo"
}

# Test 5: Verify RUNNER_TEMP is used
test_uses_runner_temp() {
    grep -q "RUNNER_TEMP" .github/actions/verify/action.yml
}

# Test 6: Verify GITHUB_PATH is used
test_uses_github_path() {
    grep -q "GITHUB_PATH" .github/actions/verify/action.yml
}

# Test 7: Verify action has required inputs
test_required_inputs() {
    local content
    content=$(cat .github/actions/verify/action.yml)
    echo "$content" | grep -q "gates:" && \
    echo "$content" | grep -q "evidence-required:" && \
    echo "$content" | grep -q "comment:" && \
    echo "$content" | grep -q "version:"
}

# Test 8: Verify action has required outputs
test_required_outputs() {
    local content
    content=$(cat .github/actions/verify/action.yml)
    echo "$content" | grep -q "result:" && \
    echo "$content" | grep -q "gates-passed:" && \
    echo "$content" | grep -q "gates-failed:"
}

# Test 9: Verify action uses composite
test_uses_composite() {
    grep -q "using: 'composite'" .github/actions/verify/action.yml
}

# Test 10: Verify comment.sh has shebang
test_comment_shebang() {
    head -n 1 .github/actions/verify/comment.sh | grep -q "#!/usr/bin/env bash"
}

# Test 11: Verify comment.sh has set -euo pipefail
test_comment_error_handling() {
    grep -q "set -euo pipefail" .github/actions/verify/comment.sh
}

# Test 12: Verify action.yml has install step
test_install_step() {
    grep -q "Install SDP CLI" .github/actions/verify/action.yml
}

# Test 13: Verify action.yml has verification step
test_verify_step() {
    grep -q "Run Verification Gates" .github/actions/verify/action.yml
}

# Test 14: Verify action.yml has PR comment step
test_comment_step() {
    grep -q "Post PR Comment" .github/actions/verify/action.yml
}

# Test 15: Verify action uses if: always() for comment
test_comment_always() {
    grep -A2 "Post PR Comment" .github/actions/verify/action.yml | grep -q "if: always()"
}

# Test 16: Verify test workflow exists
test_workflow_exists() {
    [ -f ".github/workflows/test-verify-action.yml" ]
}

# Test 17: Verify README exists
test_readme_exists() {
    [ -f ".github/actions/verify/README.md" ]
}

# Test 18: Verify action has curl download
test_has_curl_download() {
    grep -q "curl -fsSL" .github/actions/verify/action.yml
}

# Test 19: Verify action has chmod +x
test_has_chmod() {
    grep -q "chmod +x" .github/actions/verify/action.yml
}

# Test 20: Verify comment step checks PR context
test_pr_context_check() {
    grep -q "pull_request" .github/actions/verify/action.yml
}

# Run all tests
run_test() {
    local test_name="$1"
    local test_func="$2"

    echo -n "Testing: $test_name ... "

    if $test_func; then
        echo "✅ PASS"
        ((TESTS_PASSED++))
        return 0
    else
        echo "❌ FAIL"
        ((TESTS_FAILED++))
        return 1
    fi
}

run_test "action.yml exists" test_action_yml_exists
run_test "comment.sh exists" test_comment_sh_exists
run_test "comment.sh is executable" test_comment_sh_executable
run_test "action has no sudo" test_no_sudo
run_test "action uses RUNNER_TEMP" test_uses_runner_temp
run_test "action uses GITHUB_PATH" test_uses_github_path
run_test "action has required inputs" test_required_inputs
run_test "action has required outputs" test_required_outputs
run_test "action uses composite" test_uses_composite
run_test "comment.sh has shebang" test_comment_shebang
run_test "comment.sh has error handling" test_comment_error_handling
run_test "action has install step" test_install_step
run_test "action has verification step" test_verify_step
run_test "action has comment step" test_comment_step
run_test "comment uses if: always()" test_comment_always
run_test "test workflow exists" test_workflow_exists
run_test "README exists" test_readme_exists
run_test "action has curl download" test_has_curl_download
run_test "action has chmod +x" test_has_chmod
run_test "comment checks PR context" test_pr_context_check

echo ""
echo "==========================================="
echo "Test Results:"
echo "  Total:  $((TESTS_PASSED + TESTS_FAILED))"
echo "  Passed: $TESTS_PASSED"
echo "  Failed: $TESTS_FAILED"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo "✅ All tests passed!"
    exit 0
else
    echo "❌ Some tests failed"
    exit 1
fi
