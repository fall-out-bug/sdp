#!/bin/sh
# Regression checks for repo-install hook installation.

set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

create_project() {
    name="$1"
    hooks_path="$2"
    project_dir="$TMP_DIR/$name"

    mkdir -p "$project_dir/sdp/hooks"
    git -C "$project_dir" init >/dev/null 2>&1
    if [ -n "$hooks_path" ]; then
        git -C "$project_dir" config core.hooksPath "$hooks_path"
    fi

    for hook in install-git-hooks.sh pre-commit.sh pre-push.sh post-checkout.sh post-merge.sh; do
        cp "$ROOT_DIR/hooks/$hook" "$project_dir/sdp/hooks/$hook"
        chmod +x "$project_dir/sdp/hooks/$hook"
    done

    (cd "$project_dir" && sh sdp/hooks/install-git-hooks.sh >/dev/null)

    actual_hooks_dir=$(cd "$project_dir" && git rev-parse --git-path hooks)
    for hook in pre-commit pre-push post-checkout post-merge; do
        if [ ! -f "$project_dir/$actual_hooks_dir/$hook" ]; then
            echo "missing $hook in $actual_hooks_dir for $name" >&2
            exit 1
        fi
        if ! grep -q 'SDP-MANAGED-HOOK' "$project_dir/$actual_hooks_dir/$hook"; then
            echo "missing SDP marker in $hook for $name" >&2
            exit 1
        fi
    done
}

create_project "default-hooks" ""
create_project "custom-hooks" ".beads/hooks"

echo "install-project hook regression checks passed"
