#!/usr/bin/env bash
# F141-03: one-shot SDP installer for downstream repos.
# Usage: curl -fsSL https://raw.githubusercontent.com/fall-out-bug/sdp_lab/main/scripts/install.sh | bash
#
# Environment overrides:
#   SDP_REPO    GitHub repo slug (default: fall-out-bug/sdp_lab)
#   SDP_BRANCH  Branch/tag to clone (default: main)
#   SDP_SOURCE_DIR Local sdp_lab checkout to use instead of cloning
#   SDP_HARNESS Harness selection: auto|all|claude-code,opencode,... (default: auto)
#   SDP_TARGET  Target directory (default: .)
#   SDP_BIN_DIR Directory for the repo-local sdp binary (default: $SDP_TARGET/.sdp/bin)
set -euo pipefail

REPO="${SDP_REPO:-fall-out-bug/sdp_lab}"
BRANCH="${SDP_BRANCH:-main}"
SOURCE_DIR="${SDP_SOURCE_DIR:-}"
HARNESS="${SDP_HARNESS:-auto}"
TARGET="${SDP_TARGET:-.}"
TARGET_ABS="$(cd "$TARGET" 2>/dev/null && pwd -P || true)"
if [[ -z "$TARGET_ABS" ]]; then
  mkdir -p "$TARGET"
  TARGET_ABS="$(cd "$TARGET" && pwd -P)"
fi
BIN_DIR="${SDP_BIN_DIR:-$TARGET_ABS/.sdp/bin}"
LOCAL_SDP="$BIN_DIR/sdp"

echo "→ SDP installer: harness=$HARNESS target=$TARGET_ABS"

# Detect platform
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64)       ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
esac

echo "→ platform: $OS/$ARCH"

supports_current_init() {
  "$1" init --help 2>&1 | grep -q -- "--harness"
}

SDP_BIN=""

if [[ -n "$SOURCE_DIR" ]]; then
  SOURCE_ROOT="$(cd "$SOURCE_DIR" && pwd -P)"
  SOURCE_LABEL="$SOURCE_ROOT"
  echo "→ using local source: $SOURCE_ROOT"
else
  if ! command -v git >/dev/null 2>&1; then
    echo "error: required tool 'git' not found on PATH" >&2
    echo "       Please install git and re-run this script." >&2
    exit 1
  fi

  TMPDIR_SDP="$(mktemp -d)"
  trap 'rm -rf "$TMPDIR_SDP"' EXIT
  SOURCE_ROOT="$TMPDIR_SDP/sdp_lab"
  SOURCE_LABEL="$REPO@$BRANCH"

  echo "→ cloning $REPO@$BRANCH into $SOURCE_ROOT"
  git clone --depth=1 --branch "$BRANCH" "https://github.com/$REPO.git" "$SOURCE_ROOT" 2>&1
fi

# Strategy 1: if a compatible `sdp` binary is already on PATH, use it directly.
if command -v sdp >/dev/null 2>&1; then
  FOUND_SDP="$(command -v sdp)"
  if supports_current_init "$FOUND_SDP"; then
    echo "→ found compatible sdp binary on PATH: $FOUND_SDP"
    SDP_BIN="$FOUND_SDP"
  else
    echo "warning: found incompatible sdp binary on PATH: $FOUND_SDP" >&2
    echo "warning: it does not support 'sdp init --harness'; building $SOURCE_LABEL instead" >&2
  fi
fi

# Strategy 2: clone-and-build (offline-friendly, no GitHub Releases needed in v1).
# Requires: go (1.21+)
if [[ -z "$SDP_BIN" ]]; then
  if ! command -v go >/dev/null 2>&1; then
    echo "error: required tool 'go' not found on PATH" >&2
    echo "       Please install go and re-run this script." >&2
    exit 1
  fi

  echo "→ building sdp binary"
  mkdir -p "$BIN_DIR"
  (cd "$SOURCE_ROOT" && go build -tags "sqlite_fts5" -o "$LOCAL_SDP" ./cmd/sdp 2>&1)
  SDP_BIN="$LOCAL_SDP"
fi

if [[ ! -f "$TARGET_ABS/sdp.manifest.yaml" ]]; then
  cp "$SOURCE_ROOT/sdp.manifest.yaml" "$TARGET_ABS/sdp.manifest.yaml"
  echo "→ installed canonical sdp.manifest.yaml"
else
  echo "→ keeping existing sdp.manifest.yaml"
fi

if [[ ! -d "$TARGET_ABS/prompts" ]]; then
  cp -R "$SOURCE_ROOT/prompts" "$TARGET_ABS/prompts"
  echo "→ installed canonical prompts/"
else
  echo "→ keeping existing prompts/"
fi

echo "→ running sdp init"
"$SDP_BIN" init --harness "$HARNESS" --target "$TARGET_ABS"

mkdir -p "$BIN_DIR"
if [[ "$SDP_BIN" != "$LOCAL_SDP" ]]; then
  cp "$SDP_BIN" "$LOCAL_SDP"
  chmod 755 "$LOCAL_SDP"
fi

echo "✓ SDP installed in $TARGET_ABS"
echo "✓ repo-local binary: $LOCAL_SDP"
echo "→ for this shell: export PATH=\"$BIN_DIR:\$PATH\""
