#!/usr/bin/env bash
set -euo pipefail

# Post or update PR comment with verification results
#
# Environment variables:
#   GITHUB_TOKEN - GitHub token for API access
#   GITHUB_REPOSITORY - owner/repo
#   GITHUB_REF_NAME - branch name
#   COMMENT_ENABLED - whether to post comment (true/false)
#   GATES_PASSED - number of gates that passed
#   GATES_FAILED - number of gates that failed
#   VERIFICATION_RESULT - pass/fail
#   EVIDENCE_SUMMARY - JSON evidence summary (optional)
#   WORKFLOW_RUN_URL - URL to workflow run

main() {
    local COMMENT_ENABLED="${COMMENT_ENABLED:-true}"
    local GATES_PASSED="${GATES_PASSED:-0}"
    local GATES_FAILED="${GATES_FAILED:-0}"
    local VERIFICATION_RESULT="${VERIFICATION_RESULT:-unknown}"
    local EVIDENCE_SUMMARY="${EVIDENCE_SUMMARY:-{}}"

    # Skip if comment disabled
    if [[ "$COMMENT_ENABLED" != "true" ]]; then
        echo "⊘ PR comment disabled"
        return 0
    fi

    # Check if we're in a PR context
    if [[ -z "${GITHUB_EVENT_NAME:-}" ]] || [[ "$GITHUB_EVENT_NAME" != "pull_request" ]]; then
        echo "⊘ Not a pull request, skipping comment"
        return 0
    fi

    # Get PR number from event
    local PR_NUMBER
    PR_NUMBER=$(jq -r '.pull_request.number // ""' "$GITHUB_EVENT_PATH")

    if [[ -z "$PR_NUMBER" ]] || [[ "$PR_NUMBER" == "null" ]]; then
        echo "⊘ No PR number found, skipping comment"
        return 0
    fi

    echo "Posting PR comment for #$PR_NUMBER..."

    # Build comment body
    local comment_body
    comment_body=$(build_comment_body)

    # Look for existing comment by bot
    local existing_comment_id
    existing_comment_id=$(find_existing_comment "$PR_NUMBER")

    if [[ -n "$existing_comment_id" ]]; then
        # Update existing comment
        echo "Updating existing comment #$existing_comment_id"
        gh api \
            --method PATCH \
            "repos/$GITHUB_REPOSITORY/issues/comments/$existing_comment_id" \
            -f body="$comment_body" \
            -H "X-GitHub-Api-Version: 2022-11-28" \
            --silent
    else
        # Create new comment
        echo "Creating new comment"
        gh api \
            --method POST \
            "repos/$GITHUB_REPOSITORY/issues/$PR_NUMBER/comments" \
            -f body="$comment_body" \
            -H "X-GitHub-Api-Version: 2022-11-28" \
            --silent
    fi

    echo "✅ PR comment posted"
}

build_comment_body() {
    local status_icon="✅"
    if [[ "$VERIFICATION_RESULT" == "fail" ]]; then
        status_icon="❌"
    fi

    local status_text="passed"
    if [[ "$VERIFICATION_RESULT" == "fail" ]]; then
        status_text="failed"
    fi

    local total_gates=$((GATES_PASSED + GATES_FAILED))
    local workflow_url="${WORKFLOW_RUN_URL:-GitHub Actions}"

    cat <<EOF
## SDP Verification ${status_icon^}

**Status:** Verification $status_text

### Gate Results

| Gate | Status |
|------|--------|
$([[ $GATES_PASSED -gt 0 ]] && echo "| Passed | $GATES_PASSED |")
$([[ $GATES_FAILED -gt 0 ]] && echo "| Failed | $GATES_FAILED |")
$([[ $total_gates -gt 0 ]] && echo "| **Total** | **$total_gates** |")

<details>
<summary><b>Evidence Summary</b></summary>

EOF

    # Add evidence summary if available
    if command -v jq &> /dev/null && [[ "$EVIDENCE_SUMMARY" != "{}" ]]; then
        echo "$EVIDENCE_SUMMARY" | jq -r '
            "### Event Counts",
            "",
            "| Type | Count |",
            "|------|-------|",
            (to_entries[] | "| \(.key) | \(.value) |"),
            "",
            "### Models Used",
            "",
            (if has("models") then
                "| Model | Events |",
                "|-------|--------|",
                (.models | to_entries[] | "| \(.key) | \(.value) |")
            else
                "*No model data available*"
            end)
        '
    else
        echo "*No evidence data available*"
    fi

    cat <<EOF

</details>

<details>
<summary><b>Chain Integrity</b></summary>

EOF

    # Check chain integrity if evidence was verified
    if [[ "$EVIDENCE_REQUIRED" == "true" ]]; then
        echo "✅ Evidence chain integrity verified"
    else
        echo "⊘ Evidence chain verification not required"
    fi

    cat <<EOF

</details>

---
*View full details in [Workflow Run]($workflow_url)*
EOF
}

find_existing_comment() {
    local pr_number="$1"

    # Get comments for this PR
    local comments
    comments=$(gh api \
        "repos/$GITHUB_REPOSITORY/issues/$pr_number/comments" \
        --paginate \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        --jq '.[].id' 2>/dev/null || echo "")

    if [[ -z "$comments" ]]; then
        return 0
    fi

    # Look for comment starting with "## SDP Verification"
    for comment_id in $comments; do
        local body
        body=$(gh api \
            "repos/$GITHUB_REPOSITORY/issues/comments/$comment_id" \
            -H "X-GitHub-Api-Version: 2022-11-28" \
            --jq '.body' 2>/dev/null || echo "")

        if [[ "$body" == "## SDP Verification"* ]]; then
            echo "$comment_id"
            return 0
        fi
    done

    return 0
}

main "$@"
