#!/bin/bash
# Run SDP demo in a clean environment.
# Creates a temp dir, clones sdp, runs vhs, saves demo.gif.
set -e

REPO="${SDP_REPO:-https://github.com/fall-out-bug/sdp.git}"
REF="${SDP_REF:-main}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SDP_SOURCE="$(cd "$SCRIPT_DIR/.." && pwd)"

DEMO_ROOT=$(mktemp -d -t sdp-demo-XXXXXX)
echo "Clean environment: $DEMO_ROOT"
cd "$DEMO_ROOT"

git clone --depth 1 -b "$REF" "$REPO" sdp
cd sdp

# Use local demo/ if present (for branches not yet on main)
if [ -f "$SDP_SOURCE/demo/demo.tape" ]; then
  echo "Using local demo/ from $SDP_SOURCE"
  rm -rf demo
  cp -r "$SDP_SOURCE/demo" demo
fi

echo "Running vhs demo/demo.tape..."
vhs demo/demo.tape

# vhs Output demo/demo.gif is relative to cwd ($DEMO_ROOT/sdp)
OUT="$DEMO_ROOT/sdp/demo/demo.gif"
DEST="$SCRIPT_DIR/demo.gif"
if [ -f "$OUT" ]; then
  cp "$OUT" "$DEST"
  echo ""
  echo "Demo saved: $DEST"
else
  echo "Demo GIF not found (vhs may have failed)"
fi

echo ""
echo "Temp dir kept: $DEMO_ROOT"
echo "To remove: rm -rf $DEMO_ROOT"
