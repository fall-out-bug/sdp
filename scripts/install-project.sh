#!/bin/sh
# SDP Project Installer
#
# Installs prompts/hooks config into the current project.
# This script does NOT require installing the SDP CLI binary.
#
# Flags:
#   --preview             Show what would change without modifying anything.
#   --no-overwrite-config Preserve existing IDE config files (default behavior).
#   --overwrite-config    Overwrite existing config files (not recommended).

set -eu

SDP_DIR="${SDP_DIR:-sdp}"
SDP_IDE="${SDP_IDE:-auto}"
SDP_REF="${SDP_REF:-main}"
REMOTE="${SDP_REMOTE:-https://github.com/fall-out-bug/sdp.git}"
SDP_INSTALL_CLI="${SDP_INSTALL_CLI:-0}"
SDP_INSTALL_CLI_FROM_SOURCE="${SDP_INSTALL_CLI_FROM_SOURCE:-0}"
SDP_PRESERVE_CONFIG="${SDP_PRESERVE_CONFIG:-1}"
DEFAULT_REMOTE="https://github.com/fall-out-bug/sdp.git"
SDP_AUTO_FALLBACK_ALL=0
SDP_CONFIGURED_INTEGRATIONS=""
SDP_PREVIEW=0
SDP_BACKUP_DIR=""
SDP_PREVIEW_CHANGES=""
SDP_HOME=""  # resolved to project root after cd into SDP_DIR

for arg in "$@"; do
    case "$arg" in
        --no-overwrite-config)
            SDP_PRESERVE_CONFIG="1"
            ;;
        --overwrite-config)
            SDP_PRESERVE_CONFIG="0"
            ;;
        --preview)
            SDP_PREVIEW=1
            ;;
    esac
done

echo "SDP Project Installer"
echo "========================"
if [ "$SDP_PREVIEW" = "1" ]; then
    echo "Mode: PREVIEW (no changes will be made)"
elif [ "$SDP_PRESERVE_CONFIG" = "1" ]; then
    echo "Mode: preserve existing IDE config files"
fi

# ---------------------------------------------------------------------------
# Backup helpers
# ---------------------------------------------------------------------------

init_backup_dir() {
    # Lazily creates the backup directory on first use.
    # Not called from main flow — invoked by backup_file() so that
    # preview mode (which never backs up) never creates the dir.
    if [ -n "$SDP_BACKUP_DIR" ]; then
        return
    fi
    project_root="$(cd "$SDP_DIR/.." && pwd)"
    mkdir -p "$project_root/.sdp/backup"
    SDP_BACKUP_DIR=$(mktemp -d "$project_root/.sdp/backup/XXXXXX")
}

# Note: backup files are not manifest-tracked; they live under .sdp/ which
# is removed wholesale on uninstall (--purge) or left as-is in standard mode.
backup_file() {
    file="$1"
    if [ ! -e "$file" ]; then
        return
    fi
    init_backup_dir
    mkdir -p "$SDP_BACKUP_DIR"
    # Normalize path: strip leading ../ to get a clean relative path
    clean_path=$(echo "$file" | sed 's|^\.\./||')
    dest_dir="$SDP_BACKUP_DIR/$(dirname "$clean_path")"
    mkdir -p "$dest_dir"
    cp "$file" "$dest_dir/$(basename "$file")"
}

# ---------------------------------------------------------------------------
# Preview helpers
# ---------------------------------------------------------------------------

preview_note() {
    action="$1"
    target="$2"
    if [ "$SDP_PREVIEW" = "1" ]; then
        SDP_PREVIEW_CHANGES="$SDP_PREVIEW_CHANGES
  $action: $target"
        return 0
    fi
    return 1
}

preview_link() {
    target="$1"
    dest="$2"
    if preview_note "CREATE symlink" "$dest -> $target"; then
        return
    fi
    mkdir -p "$(dirname "$dest")"
    ln -sfn "$target" "$dest"
    manifest_record "symlink" "$dest" "$target"
}

preview_tree() {
    src_dir="$1"
    dest_dir="$2"
    if [ ! -d "$src_dir" ]; then
        return
    fi
    find "$src_dir" -type f | while IFS= read -r src; do
        rel="${src#$src_dir/}"
        preview_file "$src" "$dest_dir/$rel"
    done
}

preview_file() {
    src="$1"
    dest="$2"

    if preview_note "COPY" "$dest"; then
        return
    fi

    # Backup before overwrite (after preview check to avoid side effects)
    if [ -e "$dest" ]; then
        backup_file "$dest"
    fi

    mkdir -p "$(dirname "$dest")"
    cp "$src" "$dest"
    manifest_record "copy" "$dest"
}

# ---------------------------------------------------------------------------
# Manifest helpers — record every file/symlink we create so uninstall is safe
# ---------------------------------------------------------------------------

# Resolve the project root (parent of SDP_DIR).  Called once after cd into
# $SDP_DIR so that relative paths like ../.claude/... are stable.
sdp_resolve_home() {
    SDP_HOME="$(cd "$SDP_DIR/.." && pwd)"
}

manifest_record() {
    local action="$1" path="$2" extra="${3:-}"
    [ "$SDP_PREVIEW" = "1" ] && return 0
    [ -z "$SDP_HOME" ] && sdp_resolve_home
    mkdir -p "$SDP_HOME/.sdp"
    local entry="{\"action\":\"$action\",\"path\":\"$path\""
    [ -n "$extra" ] && entry="$entry,\"target\":\"$extra\""
    entry="$entry,\"timestamp\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}"
    echo "$entry" >> "$SDP_HOME/.sdp/manifest.jsonl"
}

# Note: directories are not manifest-tracked; uninstall handles empty-dir
# cleanup separately (rmdir) and never force-removes populated dirs.
safe_mkdir() {
    [ "$SDP_PREVIEW" = "1" ] && return 0
    mkdir -p "$1"
}

# ---------------------------------------------------------------------------
# JSON merge for settings.json (POSIX-compatible, no jq dependency)
# Merges SDP hook entries into existing settings.json arrays.
# Preserves all user-existing keys unchanged.
# ---------------------------------------------------------------------------

merge_settings_json() {
    src="$1"
    dest="$2"
    backup_path="$SDP_BACKUP_DIR"

    # If dest does not exist, simple copy
    if [ ! -f "$dest" ]; then
        if preview_note "CREATE" "$dest"; then
            return
        fi
        mkdir -p "$(dirname "$dest")"
        cp "$src" "$dest"
        manifest_record "copy" "$dest"
        return
    fi

    # Backup existing file
    if [ "$SDP_PREVIEW" = "0" ]; then
        backup_file "$dest"
    fi

    # Parse the source SDP settings to extract keys
    # We use a simple approach: for each top-level key in src,
    # if dest has it and it's an object/array, merge; otherwise set from src.

    # Check if jq is available BEFORE attempting any JSON validation.
    # If jq is absent and dest exists, we cannot safely merge or validate.
    has_jq=0
    if command -v jq >/dev/null 2>&1; then
        has_jq=1
    fi

    # Extract existing dest content
    dest_content=$(cat "$dest")

    # Guard: skip merge if dest is empty or not valid JSON.
    # Only runs when jq IS available (needed to validate JSON).
    if [ "$has_jq" = "1" ]; then
        if [ -z "$dest_content" ] || ! printf '%s' "$dest_content" | jq -e . >/dev/null 2>&1; then
            echo "  Note: $dest is empty or invalid JSON — replacing with SDP defaults."
            if [ "$SDP_PREVIEW" = "1" ]; then
                preview_note "REPLACE (invalid JSON)" "$dest"
                return
            fi
            backup_file "$dest"
            mkdir -p "$(dirname "$dest")"
            cp "$src" "$dest"
            manifest_record "copy" "$dest"
            return
        fi
    fi

    src_content=$(cat "$src")

    # Check if jq is available for proper merge
    if [ "$has_jq" = "1" ]; then
        # Proper deep merge with jq:
        # - For arrays under "hooks.*", append new entries (avoiding duplicates)
        # - For objects, merge recursively
        # - For scalars, take the source value only if missing from dest

        merged=$(jq -n \
            --argjson user "$dest_content" \
            --argjson sdp "$src_content" \
        '
            # Deep merge where USER values win on scalar conflicts.
            # SDP only ADDS keys that the user does not have yet.
            # Arrays: concatenate and deduplicate (by value equality).
            # Objects: recurse. Scalars: user wins.
            $user as $user |
            $sdp as $sdp |

            def deepmerge(u; s):
              if (u | type) == "object" and (s | type) == "object" then
                (s + u) | to_entries | map(
                  .key as $k |
                  if ($k | in(s)) == false then .
                  elif ($k | in(u)) == false then .key as $kk | {key: $kk, value: s[$kk]}
                  elif (u[$k] | type) == "object" and (s[$k] | type) == "object" then
                    {key: $k, value: (deepmerge(u[$k]; s[$k]))}
                  elif (u[$k] | type) == "array" and (s[$k] | type) == "array" then
                    {key: $k, value: (u[$k] + s[$k] | unique)}
                  else
                    {key: $k, value: u[$k]}
                  end
                ) | from_entries
              elif (u | type) == "array" and (s | type) == "array" then
                u + s | unique
              else
                u
              end;

            deepmerge($user; $sdp)
        ')

        if [ "$SDP_PREVIEW" = "1" ]; then
            preview_note "MERGE (jq)" "$dest"
            SDP_PREVIEW_CHANGES="$SDP_PREVIEW_CHANGES
    Merged content preview:
$(echo "$merged" | head -20)"
        else
            echo "$merged" > "$dest"
            manifest_record "merge" "$dest"
        fi
    else
        # No jq: fall back to preserving existing file entirely.
        # Only copy if dest doesn't already have SDP markers.
        if printf '%s' "$dest_content" | grep -q '"hooks"'; then
            # Dest already has hooks — preserve it, log a warning
            if [ "$SDP_PREVIEW" = "1" ]; then
                preview_note "SKIP (existing hooks preserved, install jq for proper merge)" "$dest"
            else
                echo "  Note: $dest already has hooks config. Install jq for automatic merge."
                echo "  Existing file preserved. SDP settings available at: $src"
            fi
        else
            # Dest exists but has no "hooks" key.  Without jq we cannot safely
            # merge — refuse to overwrite and tell the user what to do.
            if [ "$SDP_PREVIEW" = "1" ]; then
                preview_note "SKIP (jq required to merge safely)" "$dest"
            else
                echo "  WARNING: $dest exists but jq is not installed." >&2
                echo "  jq is required to merge settings.json safely." >&2
                echo "  Install jq or merge manually. SDP defaults available at: $src" >&2
            fi
        fi
    fi
}

# ---------------------------------------------------------------------------
# IDE detection
# ---------------------------------------------------------------------------

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

# ---------------------------------------------------------------------------
# File sync (with backup and merge awareness)
# ---------------------------------------------------------------------------

sync_file() {
    src="$1"
    dest="$2"

    if [ "$SDP_PRESERVE_CONFIG" = "1" ] && [ -e "$dest" ]; then
        return
    fi

    preview_file "$src" "$dest"
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

    preview_link "$target" "$dest"
}

# ---------------------------------------------------------------------------
# Managed checkout helpers
# ---------------------------------------------------------------------------

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
    echo "Warning: $SDP_DIR already exists. Updating..."
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

# ---------------------------------------------------------------------------
# Clone or update SDP checkout
# ---------------------------------------------------------------------------

if [ "$SDP_PREVIEW" = "1" ]; then
    if [ -d "$SDP_DIR" ]; then
        SDP_PREVIEW_CHANGES="$SDP_PREVIEW_CHANGES
  UPDATE: $SDP_DIR (existing checkout)"
    else
        SDP_PREVIEW_CHANGES="$SDP_PREVIEW_CHANGES
  CLONE: $SDP_DIR from $REMOTE ($SDP_REF)"
    fi
else
    if [ -d "$SDP_DIR" ]; then
        update_existing_checkout
    else
        echo "Cloning SDP (ref: $SDP_REF)..."
        git clone --depth 1 -b "$SDP_REF" "$REMOTE" "$SDP_DIR" 2>/dev/null || git clone --depth 1 "$REMOTE" "$SDP_DIR"
    fi
fi

# In preview mode, if SDP checkout doesn't exist we cannot inspect its
# contents.  Print the plan gathered so far and exit early.
if [ "$SDP_PREVIEW" = "1" ] && [ ! -d "$SDP_DIR" ]; then
    echo ""
    echo "=== PREVIEW: Changes that would be made ==="
    if [ -n "$SDP_PREVIEW_CHANGES" ]; then
        echo "$SDP_PREVIEW_CHANGES" | grep -v '^$' | sed 's/^/  /'
    else
        echo "  (no changes)"
    fi
    echo ""
    echo "Note: $SDP_DIR/ does not exist locally. A full preview requires the"
    echo "checkout to be present. Run without --preview to clone and install."
    echo "Existing configs will be backed up to .sdp/backup/ before modification."
    exit 0
fi

cd "$SDP_DIR"
sdp_resolve_home

# ---------------------------------------------------------------------------
# Optional CLI install for project mode
# ---------------------------------------------------------------------------

if [ "$SDP_INSTALL_CLI" = "1" ] && [ "$SDP_PREVIEW" = "0" ]; then
    cli_installed=0

    if [ "$REMOTE" = "$DEFAULT_REMOTE" ]; then
        echo "Installing SDP CLI binary from latest release..."
        if sh scripts/install.sh latest; then
            cli_installed=1
        else
            echo "Warning: Release binary installation failed."
            if [ "$SDP_INSTALL_CLI_FROM_SOURCE" = "1" ] && command -v go >/dev/null 2>&1; then
                echo "Trying source build fallback..."
                mkdir -p "${HOME}/.local/bin"
                if (
                    cd sdp-plugin && \
                    CGO_ENABLED=0 GOFLAGS=-buildvcs=false go build -o "${HOME}/.local/bin/sdp" ./cmd/sdp
                ); then
                    echo "  Installed CLI from source to ${HOME}/.local/bin/sdp"
                    cli_installed=1
                fi
            fi
        fi
    else
        if command -v go >/dev/null 2>&1; then
            echo "Building SDP CLI from checked-out source (custom remote)..."
            mkdir -p "${HOME}/.local/bin"
            if (
                cd sdp-plugin && \
                CGO_ENABLED=0 GOFLAGS=-buildvcs=false go build -o "${HOME}/.local/bin/sdp" ./cmd/sdp
            ); then
                echo "  Installed CLI from source to ${HOME}/.local/bin/sdp"
                cli_installed=1
            fi
        fi
    fi

    if [ "$cli_installed" != "1" ]; then
        echo "Warning: CLI installation failed. Prompts are installed, but 'sdp init' may not be available yet."
        if [ "$REMOTE" = "$DEFAULT_REMOTE" ]; then
            echo ""
            echo "   Retry CLI install:"
            echo "   curl -sSL https://raw.githubusercontent.com/${SDP_REPO:-fall-out-bug/sdp}/main/install.sh | sh -s -- --binary-only"
        else
            echo "   Retry: install Go, then run 'cd sdp/sdp-plugin && go build -o \${HOME}/.local/bin/sdp ./cmd/sdp'"
        fi
    fi
elif [ "$SDP_INSTALL_CLI" = "1" ] && [ "$SDP_PREVIEW" = "1" ]; then
    SDP_PREVIEW_CHANGES="$SDP_PREVIEW_CHANGES
  INSTALL CLI: sdp binary to ${HOME}/.local/bin/"
fi

# ---------------------------------------------------------------------------
# Setup for selected IDE
# ---------------------------------------------------------------------------

if [ "$SDP_IDE" = "auto" ]; then
    if SDP_IDE_LIST=$(detect_auto_ide); then
        echo "Setting up for auto-detected IDEs: $SDP_IDE_LIST"
    else
        status=$?
        if [ "$status" -ne 10 ]; then
            exit "$status"
        fi
        SDP_AUTO_FALLBACK_ALL=1
        echo "Setting up for all supported IDEs: $SDP_IDE_LIST"
    fi
else
    SDP_IDE_LIST="$SDP_IDE"
    echo "Setting up for: $SDP_IDE"
fi

setup_claude() {
    safe_mkdir ../.claude
    sync_link "../$SDP_DIR/prompts/skills" "../.claude/skills"
    sync_link "../$SDP_DIR/prompts/agents" "../.claude/agents"
    sync_file .claude/commands.json ../.claude/commands.json
    sync_tree_files .claude/hooks ../.claude/hooks
    sync_tree_files .claude/patterns ../.claude/patterns
    # Merge settings.json instead of overwriting
    if [ "$SDP_PRESERVE_CONFIG" = "1" ] && [ -f ../.claude/settings.json ]; then
        if [ "$SDP_PREVIEW" = "1" ]; then
            preview_note "SKIP (preserving existing)" "../.claude/settings.json"
        else
            echo "  Preserving existing .claude/settings.json (use --overwrite-config to replace)"
        fi
    else
        merge_settings_json .claude/settings.json ../.claude/settings.json
    fi
    register_integration "Claude" ".claude/"
}

setup_cursor() {
    safe_mkdir ../.cursor
    sync_link "../$SDP_DIR/prompts/skills" "../.cursor/skills"
    sync_link "../$SDP_DIR/prompts/agents" "../.cursor/agents"
    safe_mkdir ../.cursor/commands
    sync_tree_files .cursor/commands ../.cursor/commands
    register_integration "Cursor" ".cursor/"
}

setup_opencode() {
    safe_mkdir ../.opencode
    sync_link "../$SDP_DIR/prompts/skills" "../.opencode/skills"
    sync_link "../$SDP_DIR/prompts/agents" "../.opencode/agents"
    safe_mkdir ../.opencode/commands
    sync_tree_files .opencode/commands ../.opencode/commands
    register_integration "OpenCode" ".opencode/"
}

setup_codex() {
    safe_mkdir ../.codex/skills
    if [ -L ../.codex/skills/sdp ]; then
        if [ "$SDP_PREVIEW" = "1" ]; then
            SDP_PREVIEW_CHANGES="$SDP_PREVIEW_CHANGES
  REMOVE legacy symlink: ../.codex/skills/sdp"
        else
            rm -f ../.codex/skills/sdp
        fi
    fi
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
            echo "  Warning: Unknown SDP_IDE value '$ide', skipping"
            ;;
    esac
done

# ---------------------------------------------------------------------------
# .gitignore entries
# ---------------------------------------------------------------------------

if [ -f ../.gitignore ]; then
    if ! grep -q "# >>> SDP_START >>>" ../.gitignore; then
        if [ "$SDP_PREVIEW" = "1" ]; then
            SDP_PREVIEW_CHANGES="$SDP_PREVIEW_CHANGES
  APPEND: .gitignore (SDP entries with explicit markers)"
        else
            echo "" >> ../.gitignore
            echo "# >>> SDP_START >>>" >> ../.gitignore
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
            echo "# <<< SDP_END <<<" >> ../.gitignore
            echo "  Added entries to .gitignore (with explicit markers)"
            manifest_record "append" "../.gitignore"
        fi
    fi
fi

# ---------------------------------------------------------------------------
# Git hooks
# ---------------------------------------------------------------------------

# Hooks are installed by a dedicated script that manages its own manifest
# entries for individual hook files. We record the hook directory here so
# the uninstaller knows hooks were set up by SDP.
if [ -f hooks/install-git-hooks.sh ]; then
    if [ "$SDP_PREVIEW" = "1" ]; then
        SDP_PREVIEW_CHANGES="$SDP_PREVIEW_CHANGES
  INSTALL: git hooks (pre-commit, pre-push)"
    else
        if (cd .. && sh "$SDP_DIR/hooks/install-git-hooks.sh" 2>/dev/null); then
            echo "  Git hooks installed (pre-commit, pre-push)"
            manifest_record "hooks" ".git/hooks"
        fi
    fi
fi

# ---------------------------------------------------------------------------
# Preview summary
# ---------------------------------------------------------------------------

if [ "$SDP_PREVIEW" = "1" ]; then
    echo ""
    echo "=== PREVIEW: Changes that would be made ==="
    if [ -n "$SDP_PREVIEW_CHANGES" ]; then
        echo "$SDP_PREVIEW_CHANGES" | grep -v '^$' | sed 's/^/  /'
    else
        echo "  (no changes)"
    fi
    echo ""
    echo "Run without --preview to apply these changes."
    echo "Existing configs will be backed up to .sdp/backup/ before modification."
    exit 0
fi

# ---------------------------------------------------------------------------
# Final summary
# ---------------------------------------------------------------------------

echo ""
echo "  SDP project assets installed successfully!"
echo ""

# Report backup location
if [ -n "$SDP_BACKUP_DIR" ] && [ -d "$SDP_BACKUP_DIR" ]; then
    echo "  Backup: $(cd .. && pwd)/.sdp/backup/"
fi

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
