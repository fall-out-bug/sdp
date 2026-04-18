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
    if [ -f .gitignore ] && grep -q "# SDP" .gitignore; then
        UNINSTALL_PLAN="$UNINSTALL_PLAN
  REMOVE:      SDP entries from .gitignore"
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
# Build plan
# ---------------------------------------------------------------------------

# SDP managed symlinks in .claude
plan_remove_symlink .claude/skills
plan_remove_symlink .claude/agents

# SDP managed symlinks in .cursor
plan_remove_symlink .cursor/skills
plan_remove_symlink .cursor/agents

# SDP managed symlinks in .opencode
plan_remove_symlink .opencode/skills
plan_remove_symlink .opencode/agents

# SDP managed symlinks in .codex
plan_remove_symlink .codex/skills/sdp
plan_remove_symlink .codex/agents

# SDP managed files (only those we installed)
plan_remove_file .claude/commands.json
plan_remove_file .codex/INSTALL.md
plan_remove_file .codex/skills/README.md

# Hooks installed by SDP
plan_remove_git_hooks

# .gitignore block
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
if [ "$PURGE" = "1" ]; then
    plan_remove_dir .claude/hooks
    plan_remove_dir .claude/patterns
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

# Remove symlinks (safe: symlinks are always SDP-managed)
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

# Remove SDP-installed files
for file in .claude/commands.json .codex/INSTALL.md .codex/skills/README.md; do
    if [ -f "$file" ]; then
        rm -f "$file"
        echo "  Removed: $file"
    fi
done

# Remove git hooks that point to SDP
for hook in pre-commit pre-push; do
    hook_path=".git/hooks/$hook"
    if [ -L "$hook_path" ]; then
        target=$(readlink "$hook_path")
        case "$target" in
            *sdp/hooks/*|*scripts/hooks/*)
                rm -f "$hook_path"
                echo "  Removed hook: $hook_path"
                ;;
        esac
    fi
done

# Remove .gitignore SDP block
if [ -f .gitignore ] && grep -q "# SDP" .gitignore; then
    # Remove the SDP block: from "# SDP" line to the next blank line or EOF
    # Use a portable approach
    if command -v sed >/dev/null 2>&1; then
        # Remove lines between "# SDP" marker and next empty line (inclusive)
        sed -i.bak '/^# SDP$/,/^[[:space:]]*$/{ /^$/!d; }' .gitignore 2>/dev/null || \
        sed -i '' '/^# SDP$/,/^[[:space:]]*$/{ /^$/!d; }' .gitignore 2>/dev/null || true
        # Also remove specific SDP-managed entries (in case the block approach missed some)
        for entry in \
            "$SDP_DIR/.git" \
            ".claude/skills" \
            ".claude/agents" \
            ".cursor/skills" \
            ".cursor/agents" \
            ".opencode/skills" \
            ".opencode/agents" \
            ".codex/skills/sdp" \
            ".codex/agents" \
            ".prompts"; do
            # Remove the line (GNU or BSD sed)
            sed -i.bak "\|^${entry}$|d" .gitignore 2>/dev/null || \
            sed -i '' "\|^${entry}$|d" .gitignore 2>/dev/null || true
        done
        rm -f .gitignore.bak
        # Clean up trailing blank lines
        sed -i.bak -e :a -e '/^\n*$/{$d;N;ba' -e '}' .gitignore 2>/dev/null || \
        sed -i '' -e :a -e '/^\n*$/{$d;N;ba' -e '}' .gitignore 2>/dev/null || true
        rm -f .gitignore.bak
        echo "  Cleaned: .gitignore SDP entries"
    fi
fi

# Purge mode: remove additional directories
if [ "$PURGE" = "1" ]; then
    # Remove .claude/hooks (SDP-installed hook scripts)
    if [ -d .claude/hooks ]; then
        rm -rf .claude/hooks
        echo "  Removed: .claude/hooks/"
    fi

    # Remove .claude/patterns (SDP-installed patterns)
    if [ -d .claude/patterns ]; then
        rm -rf .claude/patterns
        echo "  Removed: .claude/patterns/"
    fi

    # Remove .sdp directory (config, backups, everything)
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
