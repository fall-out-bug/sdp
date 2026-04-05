#!/bin/sh
# Install Git hooks into the path Git actually uses (.git/hooks or core.hooksPath).
set -e

ROOT="$(git rev-parse --show-toplevel)"
cd "$ROOT"

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname "$0")" && pwd)"
if [ -d "$ROOT/sdp" ] && [ -f "$ROOT/sdp/hooks/pre-commit.sh" ]; then
  HOOKS_SRC="$ROOT/sdp/hooks"
elif [ -f "$ROOT/scripts/hooks/pre-commit.sh" ]; then
  HOOKS_SRC="$ROOT/scripts/hooks"
elif [ -f "$SCRIPT_DIR/pre-commit.sh" ]; then
  HOOKS_SRC="$SCRIPT_DIR"
else
  echo "install-git-hooks: hooks not found (expected sdp/hooks/ or scripts/hooks/)" >&2
  exit 1
fi

HOOKS_DIR="$(git rev-parse --git-path hooks)"
mkdir -p "$HOOKS_DIR"

install_hook() {
  hook_name="$1"
  source_path="$HOOKS_SRC/$hook_name.sh"
  target_path="$HOOKS_DIR/$hook_name"

  if [ ! -f "$source_path" ]; then
    return 0
  fi

  if [ -f "$target_path" ] && ! grep -q 'SDP-MANAGED-HOOK' "$target_path"; then
    echo "Skipped $hook_name (existing hook is not SDP-managed)"
    return 0
  fi

  tmp_path="$target_path.tmp.$$"
  if grep -q 'SDP-MANAGED-HOOK' "$source_path"; then
    cp "$source_path" "$tmp_path"
  else
    awk '
      NR == 1 && $0 ~ /^#!/ {
        print
        print "# SDP-MANAGED-HOOK"
        next
      }
      NR == 1 {
        print "#!/bin/sh"
        print "# SDP-MANAGED-HOOK"
      }
      { print }
    ' "$source_path" > "$tmp_path"
  fi

  chmod +x "$tmp_path"
  mv "$tmp_path" "$target_path"
  echo "Installed $hook_name"
}

for hook_name in pre-commit pre-push post-checkout post-merge; do
  install_hook "$hook_name"
done
