#!/bin/sh
# SDP Uninstall Script
#
# Removes SDP artifacts from the current project while preserving user data.
#
# Usage:
#   sh sdp/scripts/uninstall.sh              # Remove SDP artifacts, preserve user config
#   sh sdp/scripts/uninstall.sh --purge       # Remove everything SDP-related
#   sh sdp/scripts/uninstall.sh --dry-run     # Show what would be removed
#
# Flags:
#   --dry-run   Show plan without executing
#   --purge     Remove everything: SDP checkout, backups, all IDE integration dirs
#   -y          Skip confirmation prompt
#
# Cleanup strategy:
# 1. Manifest-based: remove files listed in .sdp/manifest.jsonl (primary)
# 2. Marker/symlink-based: remove SDP blocks from .gitignore, SDP hooks by symlink target
# 3. Legacy fallback: pattern-based removal when no manifest exists (migration path)
# Directories (.claude/, .cursor/) are not removed unless empty or SDP-owned symlinks

set -eu

SDP_DIR="${SDP_DIR:-sdp}"
: "${SDP_DIR:?ERROR: SDP_DIR must not be empty}"
DRY_RUN=0
PURGE=0
SKIP_CONFIRM=0
UNINSTALL_PLAN=""

for arg in "$@"; do
    case "$arg" in
        --dry-run)
            DRY_RUN=1
            ;;
        --purge)
            PURGE=1
            ;;
        -y)
            SKIP_CONFIRM=1
            ;;
    esac
done

echo "SDP Uninstaller"
echo "================"

# Verify we're in a project root (has .git)
if ! git rev-parse --show-toplevel >/dev/null 2>&1; then
    echo "ERROR: Not in a git repository." >&2
    exit 1
fi

PROJECT_ROOT="$(git rev-parse --show-toplevel)"
cd "$PROJECT_ROOT"

# ---------------------------------------------------------------------------
# Manifest helpers
# ---------------------------------------------------------------------------

# is_sdp_symlink PATH — returns 0 if PATH is a symlink whose target contains
# "sdp" in its resolved path.  Used to guard hook removal so we never delete
# hooks that belong to another tool.
is_sdp_symlink() {
    _path="$1"
    [ -L "$_path" ] || return 1
    _target=$(readlink "$_path")
    case "$_target" in
        *sdp*) return 0 ;;
        *)     return 1 ;;
    esac
}

MANIFEST_FILE=".sdp/manifest.jsonl"
HAS_MANIFEST=0

if [ -f "$MANIFEST_FILE" ]; then
    HAS_MANIFEST=1
fi

# ---------------------------------------------------------------------------
# Plan helpers
# ---------------------------------------------------------------------------

plan_remove_dir() {
    dir="$1"
    if [ -d "$dir" ]; then
        UNINSTALL_PLAN="$UNINSTALL_PLAN
  REMOVE DIR:  $dir/"
    fi
}

plan_remove_file() {
    file="$1"
    if [ -f "$file" ]; then
        UNINSTALL_PLAN="$UNINSTALL_PLAN
  REMOVE FILE: $file"
    fi
}

plan_remove_symlink() {
    link="$1"
    if [ -L "$link" ]; then
        UNINSTALL_PLAN="$UNINSTALL_PLAN
  REMOVE LINK: $link -> $(readlink "$link")"
    fi
}

plan_remove_gitignore_block() {
    if [ -f .gitignore ]; then
        if grep -q "# >>> SDP_START >>>" .gitignore; then
            UNINSTALL_PLAN="$UNINSTALL_PLAN
  REMOVE:      SDP entries from .gitignore (between explicit markers)"
        fi
        # Also detect legacy "# SDP" markers from older installs
        if grep -q "^# SDP" .gitignore; then
            UNINSTALL_PLAN="$UNINSTALL_PLAN
  REMOVE:      Legacy SDP entries from .gitignore (# SDP marker)"
        fi
    fi
}

plan_remove_git_hooks() {
    for hook in pre-commit pre-push; do
        hook_path=".git/hooks/$hook"
        if [ -L "$hook_path" ]; then
            target=$(readlink "$hook_path")
            case "$target" in
                *sdp/hooks/*|*scripts/hooks/*)
                    UNINSTALL_PLAN="$UNINSTALL_PLAN
  REMOVE HOOK: $hook_path -> $target"
                    ;;
            esac
        fi
    done
}

# ---------------------------------------------------------------------------
# Build plan — manifest-aware
# ---------------------------------------------------------------------------

if [ "$HAS_MANIFEST" = "1" ]; then
    echo "Using install manifest: $MANIFEST_FILE"
    # Plan from manifest
    while IFS= read -r line; do
        [ -n "$line" ] || continue
        action=$(echo "$line" | sed 's/.*"action":"\([^"]*\)".*/\1/')
        path=$(echo "$line" | sed 's/.*"path":"\([^"]*\)".*/\1/')
        case "$action" in
            symlink)
                plan_remove_symlink "$path"
                ;;
            copy|merge)
                plan_remove_file "$path"
                ;;
        esac
    done < "$MANIFEST_FILE"
else
    echo "============================================================" >&2
    echo "WARNING: Legacy install detected (no .sdp/manifest.jsonl)." >&2
    echo "Removing SDP artifacts by hardcoded pattern." >&2
    echo "If this is incorrect, restore from .sdp/backup/." >&2
    echo "============================================================" >&2
    # Legacy fallback: plan by pattern
    plan_remove_symlink .claude/skills
    plan_remove_symlink .claude/agents
    plan_remove_symlink .cursor/skills
    plan_remove_symlink .cursor/agents
    plan_remove_symlink .opencode/skills
    plan_remove_symlink .opencode/agents
    plan_remove_symlink .codex/skills/sdp
    plan_remove_symlink .codex/agents
    plan_remove_file .claude/commands.json
    plan_remove_file .codex/INSTALL.md
    plan_remove_file .codex/skills/README.md
fi

# Git hooks: NOT tracked by the manifest because they live inside .git/hooks/
# and are managed by a dedicated installer script.  Cleanup is safe because we
# only remove hooks that are symlinks pointing into an SDP path (checked via
# is_sdp_symlink / readlink).
plan_remove_git_hooks

# .gitignore: NOT tracked line-by-line in the manifest because it is a shared
# config file.  Cleanup uses marker-based removal (# >>> SDP_START >>> ...
# # <<< SDP_END <<<) so only the SDP block is stripped, preserving all other
# entries.  This is safe by design — markers are unique to SDP installs.
plan_remove_gitignore_block

# .sdp/backup directory (only in purge mode)
if [ "$PURGE" = "1" ]; then
    plan_remove_dir .sdp/backup
    plan_remove_dir .sdp
fi

# SDP checkout directory
if [ "$PURGE" = "1" ]; then
    plan_remove_dir "$SDP_DIR"
fi

# Full IDE integration directories (only in purge mode)
# Only plan removal if SDP-owned symlink or empty directory — never blind rm -rf
if [ "$PURGE" = "1" ]; then
    if [ -L ".claude/hooks" ] && is_sdp_symlink ".claude/hooks"; then
        plan_remove_symlink .claude/hooks
    elif [ -d ".claude/hooks" ] && [ -z "$(ls -A .claude/hooks 2>/dev/null)" ]; then
        plan_remove_dir .claude/hooks
    fi
    if [ -L ".claude/patterns" ] && is_sdp_symlink ".claude/patterns"; then
        plan_remove_symlink .claude/patterns
    elif [ -d ".claude/patterns" ] && [ -z "$(ls -A .claude/patterns 2>/dev/null)" ]; then
        plan_remove_dir .claude/patterns
    fi
    # Note: we do NOT remove .claude/settings.json even in purge mode
    # because it likely contains user customizations beyond SDP
    UNINSTALL_PLAN="$UNINSTALL_PLAN
  PRESERVE:    .claude/settings.json (may contain user config)"
fi

# ---------------------------------------------------------------------------
# Display plan
# ---------------------------------------------------------------------------

if [ -z "$UNINSTALL_PLAN" ]; then
    echo "No SDP artifacts found. Nothing to remove."
    exit 0
fi

echo ""
if [ "$PURGE" = "1" ]; then
    echo "Mode: PURGE (remove all SDP artifacts)"
elif [ "$DRY_RUN" = "1" ]; then
    echo "Mode: DRY RUN (no changes will be made)"
else
    echo "Mode: STANDARD (remove SDP artifacts, preserve user data)"
fi

echo ""
echo "Plan:"
echo "$UNINSTALL_PLAN" | grep -v '^$'

if [ "$DRY_RUN" = "1" ]; then
    echo ""
    echo "Run without --dry-run to execute this plan."
    exit 2
fi

# ---------------------------------------------------------------------------
# Confirm
# ---------------------------------------------------------------------------

if [ "$SKIP_CONFIRM" != "1" ]; then
    echo ""
    printf "Proceed? [y/N] "
    read -r answer
    case "$answer" in
        [yY]|[yY][eE][sS])
            ;;
        *)
            echo "Aborted."
            exit 0
            ;;
    esac
fi

# ---------------------------------------------------------------------------
# Execute
# ---------------------------------------------------------------------------

echo ""
echo "Removing SDP artifacts..."

if [ "$HAS_MANIFEST" = "1" ]; then
    # Manifest-driven removal: only remove files SDP actually created
    while IFS= read -r line; do
        [ -n "$line" ] || continue
        action=$(echo "$line" | sed 's/.*"action":"\([^"]*\)".*/\1/')
        path=$(echo "$line" | sed 's/.*"path":"\([^"]*\)".*/\1/')
        case "$action" in
            symlink)
                if [ -L "$path" ]; then
                    rm -f "$path"
                    echo "  Removed symlink: $path"
                fi
                ;;
            copy|merge)
                if [ -f "$path" ]; then
                    rm -f "$path"
                    echo "  Removed file: $path"
                fi
                ;;
        esac
    done < "$MANIFEST_FILE"
else
    # Legacy fallback: remove by known patterns
    echo "  (legacy mode — pattern-based removal)" >&2
    for link in \
        .claude/skills .claude/agents \
        .cursor/skills .cursor/agents \
        .opencode/skills .opencode/agents \
        .codex/skills/sdp .codex/agents; do
        if [ -L "$link" ]; then
            rm -f "$link"
            echo "  Removed: $link"
        fi
    done

    for file in .claude/commands.json .codex/INSTALL.md .codex/skills/README.md; do
        if [ -f "$file" ]; then
            rm -f "$file"
            echo "  Removed: $file"
        fi
    done
fi

# Remove git hooks that point into an SDP checkout.
# Only removes symlinks whose target contains "sdp" — never deletes hooks
# that belong to another tool (e.g. husky, lefthook).
for hook in pre-commit pre-push; do
    hook_path=".git/hooks/$hook"
    if is_sdp_symlink "$hook_path"; then
        rm -f "$hook_path"
        echo "  Removed hook: $hook_path"
    fi
done

# Remove .gitignore SDP block (between explicit markers)
if [ -f .gitignore ] && grep -q "# >>> SDP_START >>>" .gitignore; then
    if command -v sed >/dev/null 2>&1; then
        # Remove lines between (and including) the SDP markers
        sed -i.bak '/^# >>> SDP_START >>>$/,/^# <<< SDP_END <<<$/d' .gitignore 2>/dev/null || \
        sed -i '' '/^# >>> SDP_START >>>$/,/^# <<< SDP_END <<<$/d' .gitignore 2>/dev/null || true
        rm -f .gitignore.bak
        # Clean up consecutive blank lines left behind
        sed -i.bak '/^$/{ N; /^\n$/d; }' .gitignore 2>/dev/null || \
        sed -i '' '/^$/{ N; /^\n$/d; }' .gitignore 2>/dev/null || true
        rm -f .gitignore.bak
        echo "  Cleaned: .gitignore SDP entries (between markers)"
    fi
fi

# Remove legacy "# SDP" marker blocks from older installs.
if [ -f .gitignore ] && grep -q "^# SDP" .gitignore; then
    if command -v sed >/dev/null 2>&1; then
        sed -i.bak '/^# SDP/d' .gitignore 2>/dev/null || \
        sed -i '' '/^# SDP/d' .gitignore 2>/dev/null || true
        rm -f .gitignore.bak
        for entry in "$SDP_DIR/.git" ".claude/skills" ".claude/agents" \
                     ".cursor/skills" ".cursor/agents" \
                     ".opencode/skills" ".opencode/agents" \
                     ".codex/skills/sdp" ".codex/agents" ".prompts"; do
            sed -i.bak "\|^${entry}\$|d" .gitignore 2>/dev/null || \
            sed -i '' "\|^${entry}\$|d" .gitignore 2>/dev/null || true
            rm -f .gitignore.bak
        done
        sed -i.bak '/^$/{ N; /^\n$/d; }' .gitignore 2>/dev/null || \
        sed -i '' '/^$/{ N; /^\n$/d; }' .gitignore 2>/dev/null || true
        rm -f .gitignore.bak
        echo "  Cleaned: .gitignore legacy SDP entries (# SDP marker)"
    fi
fi

# Purge mode: remove additional directories.
# .claude/hooks and .claude/patterns are only removed in purge mode because
# they may contain user-authored scripts.  Removal is guarded:
# - SDP-owned symlink: always safe to remove
# - Real directory: only removed if empty (rmdir, not rm -rf)
if [ "$PURGE" = "1" ]; then
    # Remove .claude/hooks — only if it is an SDP-owned symlink or empty dir
    if [ -L ".claude/hooks" ] && is_sdp_symlink ".claude/hooks"; then
        rm -f ".claude/hooks"
        echo "  Removed symlink: .claude/hooks/"
    elif [ -d ".claude/hooks" ] && [ -z "$(ls -A .claude/hooks 2>/dev/null)" ]; then
        rmdir ".claude/hooks"
        echo "  Removed empty dir: .claude/hooks/"
    fi

    # Remove .claude/patterns — only if it is an SDP-owned symlink or empty dir
    if [ -L ".claude/patterns" ] && is_sdp_symlink ".claude/patterns"; then
        rm -f ".claude/patterns"
        echo "  Removed symlink: .claude/patterns/"
    elif [ -d ".claude/patterns" ] && [ -z "$(ls -A .claude/patterns 2>/dev/null)" ]; then
        rmdir ".claude/patterns"
        echo "  Removed empty dir: .claude/patterns/"
    fi

    # Remove .sdp directory (config, backups, manifest, everything)
    if [ -d .sdp ]; then
        rm -rf .sdp
        echo "  Removed: .sdp/"
    fi

    # Remove SDP checkout directory
    if [ -n "${SDP_DIR:-}" ] && [ -d "$SDP_DIR" ]; then
        rm -rf "$SDP_DIR"
        echo "  Removed: $SDP_DIR/"
    fi
fi

# Clean up empty directories left behind
for dir in .codex/skills .codex .claude .cursor .opencode; do
    if [ -d "$dir" ] && [ -z "$(ls -A "$dir" 2>/dev/null)" ]; then
        rmdir "$dir" 2>/dev/null && echo "  Cleaned empty dir: $dir/" || true
    fi
done

# Remove manifest itself in standard mode (if it still exists)
if [ "$HAS_MANIFEST" = "1" ] && [ "$PURGE" != "1" ] && [ -f "$MANIFEST_FILE" ]; then
    rm -f "$MANIFEST_FILE"
    echo "  Removed: $MANIFEST_FILE"
    # Clean up .sdp dir if empty
    if [ -d .sdp ] && [ -z "$(ls -A .sdp 2>/dev/null)" ]; then
        rmdir .sdp 2>/dev/null && echo "  Cleaned empty dir: .sdp/" || true
    fi
fi

echo ""
if [ "$PURGE" = "1" ]; then
    echo "  SDP fully purged from this project."
    echo ""
    echo "  Preserved (may contain your data):"
    echo "    .claude/settings.json"
    echo "    .cursor/ (non-SDP files)"
    echo "    .opencode/ (non-SDP files)"
else
    echo "  SDP artifacts removed. User data preserved."
    echo ""
    echo "  Preserved:"
    echo "    .claude/settings.json (your IDE settings)"
    echo "    .sdp/ (project config and backups)"
    echo "    $SDP_DIR/ (SDP checkout)"
    echo ""
    echo "  To remove everything, run: sh $SDP_DIR/scripts/uninstall.sh --purge"
    echo "  To reinstall: sh $SDP_DIR/scripts/install-project.sh"
fi
