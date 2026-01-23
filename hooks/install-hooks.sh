#!/bin/bash
# Install SDP git hooks

set -e

HOOKS_DIR=".git/hooks"
SDP_HOOKS_DIR="sdp/hooks"

# Check if in git repo
if [ ! -d ".git" ]; then
    echo "Error: Not in git repository root" >&2
    exit 1
fi

# Install post-commit hook
echo "Installing post-commit hook..."
ln -sf "../../${SDP_HOOKS_DIR}/post-commit.sh" "${HOOKS_DIR}/post-commit"
chmod +x "${HOOKS_DIR}/post-commit"

# Install other existing hooks (if any)
if [ -f "${SDP_HOOKS_DIR}/pre-commit.sh" ]; then
    echo "Installing pre-commit hook..."
    ln -sf "../../${SDP_HOOKS_DIR}/pre-commit.sh" "${HOOKS_DIR}/pre-commit"
    chmod +x "${HOOKS_DIR}/pre-commit"
fi

if [ -f "${SDP_HOOKS_DIR}/commit-msg.sh" ]; then
    echo "Installing commit-msg hook..."
    ln -sf "../../${SDP_HOOKS_DIR}/commit-msg.sh" "${HOOKS_DIR}/commit-msg"
    chmod +x "${HOOKS_DIR}/commit-msg"
fi

if [ -f "${SDP_HOOKS_DIR}/pre-push.sh" ]; then
    echo "Installing pre-push hook..."
    ln -sf "../../${SDP_HOOKS_DIR}/pre-push.sh" "${HOOKS_DIR}/pre-push"
    chmod +x "${HOOKS_DIR}/pre-push"
fi

echo "✅ Git hooks installed successfully"
echo ""
echo "Hooks:"
echo "  - pre-commit:  quality checks (time, tech debt, Python, Clean Arch, WS format)"
echo "  - post-commit: GitHub issue sync (if GITHUB_TOKEN set)"
echo "  - pre-push:    regression tests"
echo ""
echo "Required environment variables for GitHub integration:"
echo "  GITHUB_TOKEN  - GitHub personal access token"
echo "  GITHUB_REPO   - Repository in format 'owner/repo'"
echo ""
echo "Add to .env file or export in shell"
