#!/bin/bash
# SDP One-liner Installer
# Usage: curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | bash
#
# Installs SDP prompts and commands into your project.
# Works with: Claude Code, Cursor, OpenCode, Windsurf

set -e

SDP_DIR="${SDP_DIR:-sdp}"
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

# Create symlinks for all platforms
echo "🔗 Setting up symlinks..."

# Claude Code
mkdir -p ../.claude
ln -sf "../$SDP_DIR/prompts/skills" "../.claude/skills" 2>/dev/null || true
ln -sf "../$SDP_DIR/prompts/agents" "../.claude/agents" 2>/dev/null || true
cp -n .claude/commands.json ../.claude/ 2>/dev/null || true
cp -rn .claude/hooks ../.claude/ 2>/dev/null || true
cp -rn .claude/patterns ../.claude/ 2>/dev/null || true
cp -n .claude/settings.json ../.claude/ 2>/dev/null || true

# Cursor
mkdir -p ../.cursor
ln -sf "../$SDP_DIR/prompts/skills" "../.cursor/skills" 2>/dev/null || true
ln -sf "../$SDP_DIR/prompts/agents" "../.cursor/agents" 2>/dev/null || true
mkdir -p ../.cursor/commands
for cmd in .cursor/commands/*.md; do
    cp -n "$cmd" ../.cursor/commands/ 2>/dev/null || true
done

# OpenCode / Windsurf (generic .prompts)
mkdir -p ../.prompts
ln -sf "../$SDP_DIR/prompts/skills" "../.prompts/skills" 2>/dev/null || true
ln -sf "../$SDP_DIR/prompts/agents" "../.prompts/agents" 2>/dev/null || true

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
