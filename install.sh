#!/bin/bash
# SDP One-liner Installer
# Usage:
#   curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | bash
#   SDP_IDE=cursor curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | bash
#
# Installs SDP prompts and commands into your project.
# Works with: Claude Code, Cursor, OpenCode, Windsurf

set -e

SDP_DIR="${SDP_DIR:-sdp}"
SDP_IDE="${SDP_IDE:-all}"
REMOTE="${SDP_REMOTE:-https://github.com/fall-out-bug/sdp.git}"

echo "🚀 SDP Installer"
echo "================"

# Check if already installed
if [ -d "$SDP_DIR" ]; then
    echo "⚠️  $SDP_DIR already exists. Updating..."
    cd "$SDP_DIR" && git pull origin main
else
    echo "📦 Cloning SDP..."
    git clone --depth 1 "$REMOTE" "$SDP_DIR"
fi

cd "$SDP_DIR"

# Setup for selected IDE
echo "🔗 Setting up for: $SDP_IDE"

setup_claude() {
    mkdir -p ../.claude
    ln -sf "../$SDP_DIR/prompts/skills" "../.claude/skills" 2>/dev/null || true
    ln -sf "../$SDP_DIR/prompts/agents" "../.claude/agents" 2>/dev/null || true
    cp -n .claude/commands.json ../.claude/ 2>/dev/null || true
    cp -rn .claude/hooks ../.claude/ 2>/dev/null || true
    cp -rn .claude/patterns ../.claude/ 2>/dev/null || true
    cp -n .claude/settings.json ../.claude/ 2>/dev/null || true
}

setup_cursor() {
    mkdir -p ../.cursor
    ln -sf "../$SDP_DIR/prompts/skills" "../.cursor/skills" 2>/dev/null || true
    ln -sf "../$SDP_DIR/prompts/agents" "../.cursor/agents" 2>/dev/null || true
    mkdir -p ../.cursor/commands
    for cmd in .cursor/commands/*.md; do
        [ -f "$cmd" ] && cp -n "$cmd" ../.cursor/commands/ 2>/dev/null || true
    done
}

setup_opencode() {
    mkdir -p ../.opencode
    ln -sf "../$SDP_DIR/prompts/skills" "../.opencode/skills" 2>/dev/null || true
    ln -sf "../$SDP_DIR/prompts/agents" "../.opencode/agents" 2>/dev/null || true
    mkdir -p ../.opencode/commands
    for cmd in .opencode/commands/*.md; do
        [ -f "$cmd" ] && cp -n "$cmd" ../.opencode/commands/ 2>/dev/null || true
    done
}

case "$SDP_IDE" in
    claude|claude-code)
        setup_claude
        ;;
    cursor)
        setup_cursor
        ;;
    opencode|windsurf)
        setup_opencode
        ;;
    all|*)
        setup_claude
        setup_cursor
        setup_opencode
        ;;
esac

# Add to .gitignore
if [ -f ../.gitignore ]; then
    if ! grep -q "$SDP_DIR/.git" ../.gitignore; then
        echo "" >> ../.gitignore
        echo "# SDP" >> ../.gitignore
        echo "$SDP_DIR/.git" >> ../.gitignore
        echo ".claude/skills" >> ../.gitignore
        echo ".claude/agents" >> ../.gitignore
        echo ".cursor/skills" >> ../.gitignore
        echo ".cursor/agents" >> ../.gitignore
        echo ".prompts" >> ../.gitignore
        echo "✅ Added entries to .gitignore"
    fi
fi

echo ""
echo "✅ SDP installed successfully!"
echo ""
echo "Usage:"
echo "  @vision \"your product\"    # Strategic planning"
echo "  @feature \"add feature\"    # Plan feature"
echo "  @build 00-001-01           # Execute workstream"
echo "  @review F01                # Quality check"
echo ""
echo "Docs: $SDP_DIR/README.md"
