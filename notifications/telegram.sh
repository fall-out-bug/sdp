#!/bin/bash
# sdp/notifications/telegram.sh
# Send Telegram notifications for workflow events

set -e

# Configuration
TELEGRAM_BOT_TOKEN="${TELEGRAM_BOT_TOKEN:-}"
TELEGRAM_CHAT_ID="${TELEGRAM_CHAT_ID:-}"

# Check if configured
if [[ -z "$TELEGRAM_BOT_TOKEN" ]] || [[ -z "$TELEGRAM_CHAT_ID" ]]; then
    # Silently skip if not configured
    exit 0
fi

# Send notification
send_telegram() {
    local message="$1"
    local parse_mode="${2:-Markdown}"
    
    curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
        -d "chat_id=${TELEGRAM_CHAT_ID}" \
        -d "text=${message}" \
        -d "parse_mode=${parse_mode}" \
        -d "disable_web_page_preview=true" \
        > /dev/null 2>&1
}

# Main function
main() {
    local event_type="$1"
    shift
    
    case "$event_type" in
        "oneshot_started")
            local feature_id="$1"
            local ws_count="$2"
            send_telegram "🚀 *Oneshot Started*

Feature: \`$feature_id\`
Workstreams: $ws_count
Status: Executing autonomously

Branch: \`$(git branch --show-current)\`"
            ;;
        
        "oneshot_completed")
            local feature_id="$1"
            local duration="$2"
            send_telegram "✅ *Oneshot Completed*

Feature: \`$feature_id\`
Elapsed (telemetry): $duration
Status: All WS executed successfully

Ready for review: \`/codereview $feature_id\`"
            ;;
        
        "oneshot_blocked")
            local feature_id="$1"
            local ws_id="$2"
            local reason="$3"
            send_telegram "🚨 *CRITICAL: Oneshot Blocked*

Feature: \`$feature_id\`
Blocked at: \`$ws_id\`
Reason: $reason

⚠️ *HUMAN INTERVENTION REQUIRED*

To resume: \`/oneshot $feature_id --resume\`"
            ;;
        
        "ws_failed")
            local ws_id="$1"
            local error="$2"
            send_telegram "❌ *Workstream Failed*

WS: \`$ws_id\`
Error: $error

Retry: \`/build $ws_id\`"
            ;;
        
        "review_failed")
            local feature_id="$1"
            local issues_count="$2"
            send_telegram "⚠️ *Review Failed*

Feature: \`$feature_id\`
Issues: $issues_count

Status: CHANGES_REQUESTED

Check: \`tools/hw_checker/docs/workstreams/reports/$feature_id-review.md\`"
            ;;
        
        "breaking_changes")
            local count="$1"
            send_telegram "⚠️ *Breaking Changes Detected*

Changes: $count

Action required:
1. Review \`BREAKING_CHANGES.md\`
2. Complete \`MIGRATION_GUIDE.md\`
3. Add both to commit"
            ;;
        
        "e2e_failed")
            local feature_id="$1"
            local failed_tests="$2"
            send_telegram "🚨 *E2E Tests Failed*

Feature: \`$feature_id\`
Failed tests: $failed_tests

⛔ Deployment blocked until fixed"
            ;;
        
        "deploy_success")
            local feature_id="$1"
            local environment="$2"
            local version="$3"
            send_telegram "🎉 *Deployment Successful*

Feature: \`$feature_id\`
Environment: $environment
Version: $version

Status: Live"
            ;;
        
        "hotfix_deployed")
            local issue_id="$1"
            local duration="$2"
            send_telegram "🔥 *Hotfix Deployed*

Issue: \`$issue_id\`
Elapsed (telemetry): $duration
Target: Production

Status: Emergency fix live"
            ;;
        
        *)
            echo "Unknown event type: $event_type"
            exit 1
            ;;
    esac
}

# Run
main "$@"
