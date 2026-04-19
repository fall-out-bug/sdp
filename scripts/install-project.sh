#!/bin/sh
# SDP Project Installer
#
# Installs prompts/hooks config into the current project.
# This script does NOT require installing the SDP CLI binary.

set -eu

SDP_DIR="${SDP_DIR:-sdp}"
SDP_IDE="${SDP_IDE:-auto}"
SDP_REF="${SDP_REF:-main}"
REMOTE="${SDP_REMOTE:-https://github.com/fall-out-bug/sdp.git}"
SDP_INSTALL_CLI="${SDP_INSTALL_CLI:-0}"
SDP_INSTALL_CLI_FROM_SOURCE="${SDP_INSTALL_CLI_FROM_SOURCE:-0}"
SDP_PRESERVE_CONFIG="${SDP_PRESERVE_CONFIG:-0}"
DEFAULT_REMOTE="https://github.com/fall-out-bug/sdp.git"
SDP_AUTO_FALLBACK_ALL=0
SDP_CONFIGURED_INTEGRATIONS=""

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
    if [ -d "../.codex" ] || command -v codex >/dev/null 2>&1; then
        detected="$detected codex"
    fi

    if [ -z "$detected" ]; then
        echo "No supported IDE detected from PATH/project; installing all supported integrations." >&2
        echo "claude cursor opencode codex"
        return 10
    fi

    echo "$detected"
    return 0
}

register_integration() {
    label="$1"
    path="$2"
    entry="$label ($path)"

    case "
$SDP_CONFIGURED_INTEGRATIONS
" in
        *"
$entry
"*)
            return
            ;;
    esac

    if [ -n "$SDP_CONFIGURED_INTEGRATIONS" ]; then
        SDP_CONFIGURED_INTEGRATIONS="$SDP_CONFIGURED_INTEGRATIONS
$entry"
    else
        SDP_CONFIGURED_INTEGRATIONS="$entry"
    fi
}

print_configured_integrations() {
    if [ -z "$SDP_CONFIGURED_INTEGRATIONS" ]; then
        return
    fi

    echo "Configured integrations:"
    printf '%s\n' "$SDP_CONFIGURED_INTEGRATIONS" | while IFS= read -r entry; do
        [ -n "$entry" ] || continue
        echo "  - $entry"
    done
}

sync_file() {
    src="$1"
    dest="$2"

    if [ "$SDP_PRESERVE_CONFIG" = "1" ] && [ -e "$dest" ]; then
        return
    fi

    mkdir -p "$(dirname "$dest")"
    cp "$src" "$dest"
}

sync_tree_files() {
    src_dir="$1"
    dest_dir="$2"

    if [ ! -d "$src_dir" ]; then
        return
    fi

    find "$src_dir" -type f | while IFS= read -r src; do
        rel="${src#$src_dir/}"
        sync_file "$src" "$dest_dir/$rel"
    done
}

sync_link() {
    target="$1"
    dest="$2"

    if [ "$SDP_PRESERVE_CONFIG" = "1" ] && [ -e "$dest" ]; then
        return
    fi

    mkdir -p "$(dirname "$dest")"
    ln -sfn "$target" "$dest"
}

ensure_managed_checkout() {
    if ! git -C "$SDP_DIR" rev-parse --git-dir >/dev/null 2>&1; then
        echo "ERROR: $SDP_DIR exists but is not a git checkout." >&2
        echo "Move or remove it, then rerun the installer." >&2
        exit 1
    fi
}

ensure_clean_checkout() {
    if [ -n "$(git -C "$SDP_DIR" status --porcelain)" ]; then
        echo "ERROR: $SDP_DIR has local changes." >&2
        echo "Commit, stash, or remove them before rerunning the installer." >&2
        git -C "$SDP_DIR" status --short >&2 || true
        exit 1
    fi
}

update_existing_checkout() {
    echo "⚠️  $SDP_DIR already exists. Updating..."
    ensure_managed_checkout
    ensure_clean_checkout

    if git -C "$SDP_DIR" remote get-url origin >/dev/null 2>&1; then
        git -C "$SDP_DIR" remote set-url origin "$REMOTE"
    else
        git -C "$SDP_DIR" remote add origin "$REMOTE"
    fi

    git -C "$SDP_DIR" fetch --depth 1 origin "$SDP_REF"
    git -C "$SDP_DIR" checkout -B "$SDP_REF" FETCH_HEAD >/dev/null 2>&1
}

# Check if already installed
if [ -d "$SDP_DIR" ]; then
    update_existing_checkout
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
    if SDP_IDE_LIST=$(detect_auto_ide); then
        echo "🔗 Setting up for auto-detected IDEs: $SDP_IDE_LIST"
    else
        status=$?
        if [ "$status" -ne 10 ]; then
            exit "$status"
        fi
        SDP_AUTO_FALLBACK_ALL=1
        echo "🔗 Setting up for all supported IDEs: $SDP_IDE_LIST"
    fi
else
    SDP_IDE_LIST="$SDP_IDE"
    echo "🔗 Setting up for: $SDP_IDE"
fi

setup_claude() {
    mkdir -p ../.claude
    sync_link "../$SDP_DIR/prompts/skills" "../.claude/skills"
    sync_link "../$SDP_DIR/prompts/agents" "../.claude/agents"
    sync_file .claude/commands.json ../.claude/commands.json
    sync_tree_files .claude/hooks ../.claude/hooks
    sync_tree_files .claude/patterns ../.claude/patterns
    sync_file .claude/settings.json ../.claude/settings.json
    register_integration "Claude" ".claude/"
}

setup_cursor() {
    mkdir -p ../.cursor
    sync_link "../$SDP_DIR/prompts/skills" "../.cursor/skills"
    sync_link "../$SDP_DIR/prompts/agents" "../.cursor/agents"
    mkdir -p ../.cursor/commands
    sync_tree_files .cursor/commands ../.cursor/commands
    register_integration "Cursor" ".cursor/"
}

setup_opencode() {
    mkdir -p ../.opencode
    sync_link "../$SDP_DIR/prompts/skills" "../.opencode/skills"
    sync_link "../$SDP_DIR/prompts/agents" "../.opencode/agents"
    mkdir -p ../.opencode/commands
    sync_tree_files .opencode/commands ../.opencode/commands
    register_integration "OpenCode" ".opencode/"
}

setup_codex() {
    mkdir -p ../.codex/skills
    # Individual skill symlinks (language-agnostic, per-skill granularity)
    for _skill_dir in "$SDP_DIR"/prompts/skills/*/; do
        _skill_name="$(basename "$_skill_dir")"
        [ -f "${_skill_dir}SKILL.md" ] && sync_link "../../$SDP_DIR/prompts/skills/$_skill_name" "../.codex/skills/$_skill_name"
    done
    sync_link "../$SDP_DIR/prompts/agents" "../.codex/agents"
    sync_file .codex/INSTALL.md ../.codex/INSTALL.md
    sync_file .codex/skills/README.md ../.codex/skills/README.md
    register_integration "Codex" ".codex/"
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
        codex)
            setup_codex
            ;;
        all)
            setup_claude
            setup_cursor
            setup_opencode
            setup_codex
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
        echo ".opencode/skills" >> ../.gitignore
        echo ".opencode/agents" >> ../.gitignore
        echo ".codex/skills" >> ../.gitignore
        echo ".codex/agents" >> ../.gitignore
        echo ".prompts" >> ../.gitignore
        echo "✅ Added entries to .gitignore"
    fi
fi

# Install Git hooks (pre-commit, pre-push)
if [ -f hooks/install-git-hooks.sh ]; then
    if (cd .. && sh "$SDP_DIR/hooks/install-git-hooks.sh" 2>/dev/null); then
        echo "✅ Git hooks installed (pre-commit, pre-push)"
    fi
fi

echo ""
echo "✅ SDP project assets installed successfully!"
echo ""
print_configured_integrations

if [ "$SDP_AUTO_FALLBACK_ALL" = "1" ]; then
    echo ""
    echo "Note: no supported IDE was detected, so SDP installed all supported integrations."
    echo "To install only one surface, rerun with SDP_IDE=claude|cursor|opencode|codex."
fi

cli_path=""
init_cmd=""
update_cmd="curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only"

if [ -x "${HOME}/.local/bin/sdp" ]; then
    cli_path="${HOME}/.local/bin/sdp"
    init_cmd="${HOME}/.local/bin/sdp init --auto"
elif command -v sdp >/dev/null 2>&1; then
    cli_path=$(command -v sdp)
    init_cmd="sdp init --auto"
fi

echo ""
if [ -n "$cli_path" ]; then
    echo "CLI: ${cli_path}"
    echo ""
    echo "Next:"
    echo "  1. Run ${init_cmd}"
    echo "  2. After init, review .sdp/config.yml"
    echo "  3. Open this repo in your IDE"
    echo ""
    echo "Update CLI:"
    echo "  ${update_cmd}"
else
    echo "CLI not found in PATH. Install binary with:"
    echo "  ${update_cmd}"
    echo ""
    echo "Next:"
    echo "  1. Install the CLI command above"
    echo "  2. Run sdp init --auto"
    echo "  3. After init, review .sdp/config.yml"
fi
