#!/bin/sh
# SDP Project Installer
#
# Installs prompts/hooks config into the current project.
# This script does NOT require installing the SDP CLI binary.

set -e

SDP_DIR="${SDP_DIR:-sdp}"
SDP_IDE="${SDP_IDE:-auto}"
SDP_REF="${SDP_REF:-main}"
REMOTE="${SDP_REMOTE:-https://github.com/fall-out-bug/sdp.git}"
SDP_INSTALL_CLI="${SDP_INSTALL_CLI:-0}"
SDP_INSTALL_CLI_FROM_SOURCE="${SDP_INSTALL_CLI_FROM_SOURCE:-0}"
SDP_PRESERVE_CONFIG="${SDP_PRESERVE_CONFIG:-0}"
DEFAULT_REMOTE="https://github.com/fall-out-bug/sdp.git"

for arg in "$@"; do
    case "$arg" in
        --no-overwrite-config)
            SDP_PRESERVE_CONFIG="1"
            ;;
        --overwrite-config)
            SDP_PRESERVE_CONFIG="0"
            ;;
    esac
done

echo "🚀 SDP Project Installer"
echo "========================"
if [ "$SDP_PRESERVE_CONFIG" = "1" ]; then
    echo "Mode: preserve existing IDE config files"
fi

detect_auto_ide() {
    detected=""

    if [ -d "../.claude" ] || command -v claude >/dev/null 2>&1; then
        detected="$detected claude"
    fi
    if [ -d "../.cursor" ] || command -v cursor >/dev/null 2>&1; then
        detected="$detected cursor"
    fi
    if [ -d "../.opencode" ] || command -v opencode >/dev/null 2>&1 || command -v windsurf >/dev/null 2>&1; then
        detected="$detected opencode"
    fi

    if [ -z "$detected" ]; then
        echo "No IDE detected from PATH/project; falling back to all integrations." >&2
        echo "claude cursor opencode"
        return
    fi

    echo "$detected"
}

# Check if already installed
if [ -d "$SDP_DIR" ]; then
    echo "⚠️  $SDP_DIR already exists. Updating..."
    git -C "$SDP_DIR" fetch origin "$SDP_REF" && git -C "$SDP_DIR" checkout "$SDP_REF" 2>/dev/null || git -C "$SDP_DIR" pull origin main
else
    echo "📦 Cloning SDP (ref: $SDP_REF)..."
    git clone --depth 1 -b "$SDP_REF" "$REMOTE" "$SDP_DIR" 2>/dev/null || git clone --depth 1 "$REMOTE" "$SDP_DIR"
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
            echo ""
            echo "   Retry CLI install:"
            echo "   curl -sSL https://raw.githubusercontent.com/${SDP_REPO:-fall-out-bug/sdp}/main/install.sh | sh -s -- --binary-only"
        else
            echo "   Retry: install Go, then run 'cd sdp/sdp-plugin && go build -o \${HOME}/.local/bin/sdp ./cmd/sdp'"
        fi
    fi
fi

# Setup for selected IDE
if [ "$SDP_IDE" = "auto" ]; then
    SDP_IDE_LIST=$(detect_auto_ide)
    echo "🔗 Setting up for auto-detected IDEs: $SDP_IDE_LIST"
else
    SDP_IDE_LIST="$SDP_IDE"
    echo "🔗 Setting up for: $SDP_IDE"
fi

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

for ide in $SDP_IDE_LIST; do
    case "$ide" in
        claude|claude-code)
            setup_claude
            ;;
        cursor)
            setup_cursor
            ;;
        opencode|windsurf)
            setup_opencode
            ;;
        all)
            setup_claude
            setup_cursor
            setup_opencode
            ;;
        *)
            echo "⚠️  Unknown SDP_IDE value '$ide', skipping"
            ;;
    esac
done

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
echo "✅ SDP project assets installed successfully!"
echo ""
if [ -x "${HOME}/.local/bin/sdp" ]; then
    echo "CLI: ${HOME}/.local/bin/sdp"
    echo "Try: ${HOME}/.local/bin/sdp init --auto"
    echo "     (update CLI: curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only)"
elif command -v sdp >/dev/null 2>&1; then
    cli_path=$(command -v sdp)
    echo "CLI: ${cli_path}"
    echo "Try: sdp init --auto"
    echo "     (update CLI: curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only)"
else
    echo "CLI not found in PATH. Install binary with:"
    echo "  curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only"
fi
