#!/bin/sh
# SDP One-liner Installer
# Usage:
#   curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh
#   curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only
#   curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | SDP_IDE=cursor sh
#   curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --no-overwrite-config
#
# Default: install prompts/hooks config in your project
# --binary-only: install sdp binary globally (no project files)
#
# For binary installer directly:
#   curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/scripts/install.sh | sh
#
# Installs SDP prompts/hooks in project and optionally CLI.
# Works with: Claude Code, Cursor, OpenCode, Windsurf

set -e

SDP_DIR="${SDP_DIR:-sdp}"
SDP_IDE="${SDP_IDE:-all}"
REMOTE="${SDP_REMOTE:-https://github.com/fall-out-bug/sdp.git}"
SDP_INSTALL_CLI="${SDP_INSTALL_CLI:-0}"
SDP_INSTALL_CLI_FROM_SOURCE="${SDP_INSTALL_CLI_FROM_SOURCE:-0}"
SDP_PRESERVE_CONFIG="${SDP_PRESERVE_CONFIG:-0}"
SDP_CLI_VERSION="${SDP_CLI_VERSION:-latest}"
DEFAULT_REMOTE="https://github.com/fall-out-bug/sdp.git"
DEFAULT_REPO="fall-out-bug/sdp"
SDP_REPO="${SDP_REPO:-$DEFAULT_REPO}"
SDP_INSTALL_SCRIPT_URL="${SDP_INSTALL_SCRIPT_URL:-https://raw.githubusercontent.com/${SDP_REPO}/main/scripts/install.sh}"
BINARY_ONLY=0

for arg in "$@"; do
    case "$arg" in
        --binary-only)
            BINARY_ONLY=1
            ;;
        --no-overwrite-config)
            SDP_PRESERVE_CONFIG="1"
            ;;
        --overwrite-config)
            SDP_PRESERVE_CONFIG="0"
            ;;
    esac
done

echo "🚀 SDP Installer"
echo "================"
if [ "$SDP_PRESERVE_CONFIG" = "1" ]; then
    echo "Mode: preserve existing IDE config files"
fi

if [ "$BINARY_ONLY" = "1" ]; then
    echo "📦 Installing SDP CLI binary only..."
    if ! curl -fsSL "${SDP_INSTALL_SCRIPT_URL}" | SDP_REPO="${SDP_REPO}" sh -s -- "${SDP_CLI_VERSION}"; then
        echo "❌ Binary install failed"
        exit 1
    fi

    echo ""
    echo "✅ Binary installed. Next step in your project:"
    echo "   sdp init --auto"
    echo ""
    echo "sdp init will fetch prompts automatically if local prompts are missing."
    exit 0
fi

# Check if already installed
if [ -d "$SDP_DIR" ]; then
    echo "⚠️  $SDP_DIR already exists. Updating..."
    git -C "$SDP_DIR" pull origin main
else
    echo "📦 Cloning SDP..."
    git clone --depth 1 "$REMOTE" "$SDP_DIR"
fi

cd "$SDP_DIR"

# Optional CLI install for project mode
if [ "$SDP_INSTALL_CLI" = "1" ]; then
    cli_installed=0

    if [ "$REMOTE" = "$DEFAULT_REMOTE" ]; then
        echo "🔧 Installing SDP CLI binary from latest release..."
        if sh scripts/install.sh latest; then
            cli_installed=1
        else
            echo "⚠️  Release binary installation failed."
            if [ "$SDP_INSTALL_CLI_FROM_SOURCE" = "1" ] && command -v go >/dev/null 2>&1; then
                echo "🔧 Trying source build fallback..."
                mkdir -p "${HOME}/.local/bin"
                if (
                    cd sdp-plugin && \
                    CGO_ENABLED=0 GOFLAGS=-buildvcs=false go build -o "${HOME}/.local/bin/sdp" ./cmd/sdp
                ); then
                    echo "✅ Installed CLI from source to ${HOME}/.local/bin/sdp"
                    cli_installed=1
                fi
            fi
        fi
    else
        if command -v go >/dev/null 2>&1; then
            echo "🔧 Building SDP CLI from checked-out source (custom remote)..."
            mkdir -p "${HOME}/.local/bin"
            if (
                cd sdp-plugin && \
                CGO_ENABLED=0 GOFLAGS=-buildvcs=false go build -o "${HOME}/.local/bin/sdp" ./cmd/sdp
            ); then
                echo "✅ Installed CLI from source to ${HOME}/.local/bin/sdp"
                cli_installed=1
            fi
        fi
    fi

    if [ "$cli_installed" != "1" ]; then
        echo "⚠️  CLI installation failed. Prompts are installed, but 'sdp init' may not be available yet."
        if [ "$REMOTE" = "$DEFAULT_REMOTE" ]; then
            echo "   Retry manually: sh scripts/install.sh"
        else
            echo "   Retry manually after installing Go toolchain."
        fi
    fi
fi

# Setup for selected IDE
echo "🔗 Setting up for: $SDP_IDE"

setup_claude() {
    mkdir -p ../.claude
    if [ "$SDP_PRESERVE_CONFIG" = "1" ] && [ -e "../.claude/skills" ]; then
        :
    else
        ln -sf "../$SDP_DIR/prompts/skills" "../.claude/skills" 2>/dev/null || true
    fi
    if [ "$SDP_PRESERVE_CONFIG" = "1" ] && [ -e "../.claude/agents" ]; then
        :
    else
        ln -sf "../$SDP_DIR/prompts/agents" "../.claude/agents" 2>/dev/null || true
    fi
    cp -n .claude/commands.json ../.claude/ 2>/dev/null || true
    cp -rn .claude/hooks ../.claude/ 2>/dev/null || true
    cp -rn .claude/patterns ../.claude/ 2>/dev/null || true
    cp -n .claude/settings.json ../.claude/ 2>/dev/null || true
}

setup_cursor() {
    mkdir -p ../.cursor
    if [ "$SDP_PRESERVE_CONFIG" = "1" ] && [ -e "../.cursor/skills" ]; then
        :
    else
        ln -sf "../$SDP_DIR/prompts/skills" "../.cursor/skills" 2>/dev/null || true
    fi
    if [ "$SDP_PRESERVE_CONFIG" = "1" ] && [ -e "../.cursor/agents" ]; then
        :
    else
        ln -sf "../$SDP_DIR/prompts/agents" "../.cursor/agents" 2>/dev/null || true
    fi
    mkdir -p ../.cursor/commands
    for cmd in .cursor/commands/*.md; do
        [ -f "$cmd" ] && cp -n "$cmd" ../.cursor/commands/ 2>/dev/null || true
    done
}

setup_opencode() {
    mkdir -p ../.opencode
    if [ "$SDP_PRESERVE_CONFIG" = "1" ] && [ -e "../.opencode/skills" ]; then
        :
    else
        ln -sf "../$SDP_DIR/prompts/skills" "../.opencode/skills" 2>/dev/null || true
    fi
    if [ "$SDP_PRESERVE_CONFIG" = "1" ] && [ -e "../.opencode/agents" ]; then
        :
    else
        ln -sf "../$SDP_DIR/prompts/agents" "../.opencode/agents" 2>/dev/null || true
    fi
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
if [ -x "${HOME}/.local/bin/sdp" ]; then
    echo "CLI: ${HOME}/.local/bin/sdp"
    if "${HOME}/.local/bin/sdp" init --help 2>/dev/null | grep -q -- "--guided"; then
        echo "Try: ${HOME}/.local/bin/sdp init --guided"
    else
        echo "Try: ${HOME}/.local/bin/sdp init --auto"
    fi
elif command -v sdp >/dev/null 2>&1; then
    cli_path=$(command -v sdp)
    echo "CLI: ${cli_path}"
    if "$cli_path" init --help 2>/dev/null | grep -q -- "--guided"; then
        echo "Try: sdp init --guided"
    else
        echo "Try: sdp init --auto"
    fi
else
    echo "CLI not found in PATH. Restart shell or add ~/.local/bin to PATH."
fi
echo ""
echo "Usage:"
echo "  @vision \"your product\"    # Strategic planning"
echo "  @feature \"add feature\"    # Plan feature"
echo "  @build 00-001-01           # Execute workstream"
echo "  @review F01                # Quality check"
echo ""
echo "Docs: $SDP_DIR/README.md"
