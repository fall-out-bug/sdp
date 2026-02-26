#!/bin/sh
# Install Git hooks: symlink .git/hooks/* to SDP hooks.
# Run from repo root (or from sdp/). Idempotent.
# When sdp is submodule: .git/hooks/pre-commit -> sdp/hooks/pre-commit.sh
set -e

ROOT="$(git rev-parse --show-toplevel)"
cd "$ROOT"

# Detect hooks source: sdp/hooks/ (when sdp submodule) or scripts/hooks/ (sdp_dev)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
if [ -d "$ROOT/sdp" ] && [ -f "$ROOT/sdp/hooks/pre-commit.sh" ]; then
  HOOKS_SRC="$ROOT/sdp/hooks"
  REL_HOOKS="../../sdp/hooks"
elif [ -f "$ROOT/scripts/hooks/pre-commit.sh" ]; then
  HOOKS_SRC="$ROOT/scripts/hooks"
  REL_HOOKS="../../scripts/hooks"
elif [ -f "$SCRIPT_DIR/pre-commit.sh" ]; then
  HOOKS_SRC="$SCRIPT_DIR"
  REL_HOOKS="../../$(echo "$HOOKS_SRC" | sed "s|^$ROOT/||")"
else
  echo "install-git-hooks: hooks not found (expected sdp/hooks/ or scripts/hooks/)" >&2
  exit 1
fi

HOOKS_DIR="$ROOT/.git/hooks"
mkdir -p "$HOOKS_DIR"

if [ -f "$HOOKS_SRC/pre-commit.sh" ]; then
  ln -sf "$REL_HOOKS/pre-commit.sh" "$HOOKS_DIR/pre-commit"
  chmod +x "$HOOKS_SRC/pre-commit.sh"
  echo "Installed pre-commit"
fi

if [ -f "$HOOKS_SRC/pre-push.sh" ]; then
  ln -sf "$REL_HOOKS/pre-push.sh" "$HOOKS_DIR/pre-push"
  chmod +x "$HOOKS_SRC/pre-push.sh"
  echo "Installed pre-push"
fi
