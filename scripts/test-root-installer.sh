#!/bin/sh
# Regression checks for root install.sh dispatcher behavior.

set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

MOCK_REMOTE_DIR="$TMP_DIR/mock-remote"
MOCK_BIN_DIR="$TMP_DIR/mock-bin"
mkdir -p "$MOCK_REMOTE_DIR" "$MOCK_BIN_DIR"

cat > "$MOCK_REMOTE_DIR/install-project.sh" <<'EOF'
#!/bin/sh
set -eu
printf 'project %s\n' "$*" > "$MOCK_OUT"
EOF

cat > "$MOCK_REMOTE_DIR/install.sh" <<'EOF'
#!/bin/sh
set -eu
printf 'binary %s\n' "$*" > "$MOCK_OUT"
EOF

cat > "$MOCK_BIN_DIR/curl" <<'EOF'
#!/bin/sh
set -eu
url=""
for arg in "$@"; do
    case "$arg" in
        http://*|https://*)
            url="$arg"
            ;;
    esac
done

case "$url" in
    */scripts/install-project.sh)
        cat "$MOCK_REMOTE_DIR/install-project.sh"
        ;;
    */scripts/install.sh)
        cat "$MOCK_REMOTE_DIR/install.sh"
        ;;
    *)
        echo "unexpected URL for mock curl: $url" >&2
        exit 1
        ;;
esac
EOF

chmod +x "$MOCK_REMOTE_DIR/install-project.sh" "$MOCK_REMOTE_DIR/install.sh" "$MOCK_BIN_DIR/curl"

run_case() {
    mode="$1"
    args="$2"

    case_dir="$TMP_DIR/$mode"
    mkdir -p "$case_dir"
    cp "$ROOT_DIR/install.sh" "$case_dir/install.sh"

    out_file="$TMP_DIR/${mode}.out"
    : > "$out_file"

    # shellcheck disable=SC2086
    (cd "$case_dir" && PATH="$MOCK_BIN_DIR:$PATH" MOCK_REMOTE_DIR="$MOCK_REMOTE_DIR" MOCK_OUT="$out_file" SDP_REPO="fall-out-bug/sdp" SDP_REF="main" sh ./install.sh $args)

    if [ ! -s "$out_file" ]; then
        echo "root installer did not execute remote script in $mode mode" >&2
        exit 1
    fi

    case "$mode" in
        project)
            grep -q '^project ' "$out_file"
            ;;
        binary)
            grep -q '^binary latest$' "$out_file"
            ;;
    esac
}

run_case "project" ""
run_case "binary" "--binary-only"

echo "root installer regression checks passed"
