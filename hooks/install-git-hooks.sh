#!/bin/sh
# Install Git hooks: symlink .git/hooks/* to scripts/hooks/*.sh.
# Run from repo root. Idempotent.
# F128: hooks are native in scripts/hooks/. The sdp/ local clone is a downstream
# mirror (gitignored) and must not provide hook sources to avoid split-brain.
set -e

ROOT="$(git rev-parse --show-toplevel)"
cd "$ROOT"

HOOKS_DIR=".git/hooks"
SCRIPTS_DIR="scripts/hooks"
mkdir -p "$HOOKS_DIR"

if [ -f "$SCRIPTS_DIR/pre-commit.sh" ]; then
  ln -sf ../../scripts/hooks/pre-commit.sh "$HOOKS_DIR/pre-commit"
  chmod +x "$SCRIPTS_DIR/pre-commit.sh"
  echo "Installed pre-commit"
fi

if [ -f "$SCRIPTS_DIR/pre-push.sh" ]; then
  ln -sf ../../scripts/hooks/pre-push.sh "$HOOKS_DIR/pre-push"
  chmod +x "$SCRIPTS_DIR/pre-push.sh"
  echo "Installed pre-push"
fi
