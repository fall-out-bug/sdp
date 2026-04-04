#!/bin/sh
# Regression checks for install-project.sh end-to-end behavior.

set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
TMP_DIR=$(mktemp -d)
trap 'chmod -R u+w "$TMP_DIR" 2>/dev/null || true; rm -rf "$TMP_DIR"' EXIT

SOURCE_BARE="$TMP_DIR/source.git"
ADMIN_DIR="$TMP_DIR/admin"
HOME_DIR="$TMP_DIR/home"
PROJECT_DIR="$TMP_DIR/project"
FULL_PROJECT_DIR="$TMP_DIR/project-full"
CODEX_PROJECT_DIR="$TMP_DIR/project-codex"

mkdir -p "$HOME_DIR"

git clone --bare "$ROOT_DIR" "$SOURCE_BARE" >/dev/null
git clone "$SOURCE_BARE" "$ADMIN_DIR" >/dev/null
git -C "$ADMIN_DIR" config user.name "SDP Installer Test"
git -C "$ADMIN_DIR" config user.email "installer-test@example.com"
git -C "$ADMIN_DIR" checkout -B main >/dev/null
git -C "$ADMIN_DIR" push origin HEAD:refs/heads/main >/dev/null

run_install() {
    project_dir="$1"
    log_file="$2"
    shift 2

    mkdir -p "$project_dir"
    if [ ! -d "$project_dir/.git" ]; then
        git -C "$project_dir" init -q
    fi

    (
        cd "$project_dir"
        HOME="$HOME_DIR" \
        SDP_REMOTE="$SOURCE_BARE" \
        SDP_IDE=claude \
        "$@" \
        sh "$ROOT_DIR/scripts/install-project.sh"
    ) >"$log_file" 2>&1
}

hash_file() {
    file="$1"
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$file" | awk '{print $1}'
    else
        shasum -a 256 "$file" | awk '{print $1}'
    fi
}

assert_contains() {
    needle="$1"
    file="$2"
    if ! grep -Fq "$needle" "$file"; then
        echo "expected '$needle' in $file" >&2
        exit 1
    fi
}

update_admin_file() {
    file="$1"
    from="$2"
    to="$3"
    message="$4"

    perl -0pi -e "s/\Q$from\E/$to/" "$ADMIN_DIR/$file"
    git -C "$ADMIN_DIR" commit -am "$message" >/dev/null
    git -C "$ADMIN_DIR" push origin HEAD:refs/heads/main >/dev/null
}

# Cold start / clean install
run_install "$PROJECT_DIR" "$TMP_DIR/clean-install.log" env
test -d "$PROJECT_DIR/sdp/.git"
test -L "$PROJECT_DIR/.claude/skills"
test -f "$PROJECT_DIR/.claude/commands.json"
assert_contains '"version": "1.1.0"' "$PROJECT_DIR/.claude/commands.json"

# Clean reinstall / update should refresh vendored checkout and managed files.
update_admin_file ".claude/commands.json" '"version": "1.1.0"' '"version": "9.9.9"' "test: update commands manifest"
run_install "$PROJECT_DIR" "$TMP_DIR/update-install.log" env
assert_contains '"version": "9.9.9"' "$PROJECT_DIR/.claude/commands.json"
assert_contains '"version": "9.9.9"' "$PROJECT_DIR/sdp/.claude/commands.json"

# Dirty reinstall should fail clearly before git noise if managed checkout changed.
perl -0pi -e 's/"version": "9\.9\.9"/"version": "LOCAL-DIRTY"/' "$PROJECT_DIR/sdp/.claude/commands.json"
update_admin_file ".claude/commands.json" '"version": "9.9.9"' '"version": "10.0.0"' "test: update commands manifest again"
if run_install "$PROJECT_DIR" "$TMP_DIR/dirty-install.log" env; then
    echo "dirty installer rerun unexpectedly succeeded" >&2
    exit 1
fi
assert_contains "ERROR: sdp has local changes." "$TMP_DIR/dirty-install.log"
assert_contains ".claude/commands.json" "$TMP_DIR/dirty-install.log"

# Full install path should build/update CLI when requested.
run_install "$FULL_PROJECT_DIR" "$TMP_DIR/full-install.log" env SDP_INSTALL_CLI=1
test -x "$HOME_DIR/.local/bin/sdp"
before_hash=$(hash_file "$HOME_DIR/.local/bin/sdp")
update_admin_file "sdp-plugin/cmd/sdp/main.go" 'var version = "dev"' 'var version = "dev-full-install-update"' "test: update cli source"
run_install "$FULL_PROJECT_DIR" "$TMP_DIR/full-update.log" env SDP_INSTALL_CLI=1
after_hash=$(hash_file "$HOME_DIR/.local/bin/sdp")
if [ "$before_hash" = "$after_hash" ]; then
    echo "full installer did not refresh CLI binary" >&2
    exit 1
fi

# Codex install/update should provision project-level Codex surface and refresh managed links.
mkdir -p "$CODEX_PROJECT_DIR"
: > "$CODEX_PROJECT_DIR/.gitignore"
run_install "$CODEX_PROJECT_DIR" "$TMP_DIR/codex-install.log" env SDP_IDE=codex
test -d "$CODEX_PROJECT_DIR/sdp/.git"
test -f "$CODEX_PROJECT_DIR/.codex/INSTALL.md"
test -f "$CODEX_PROJECT_DIR/.codex/skills/README.md"
test -L "$CODEX_PROJECT_DIR/.codex/skills/sdp"
test -L "$CODEX_PROJECT_DIR/.codex/agents"
test ! -e "$CODEX_PROJECT_DIR/.claude"
assert_contains ".codex/skills/sdp" "$CODEX_PROJECT_DIR/.gitignore"

printf '\n<!-- codex update marker -->\n' >> "$ADMIN_DIR/prompts/skills/build/SKILL.md"
git -C "$ADMIN_DIR" commit -am "test: update codex skill source" >/dev/null
git -C "$ADMIN_DIR" push origin HEAD:refs/heads/main >/dev/null
run_install "$CODEX_PROJECT_DIR" "$TMP_DIR/codex-update.log" env SDP_IDE=codex
assert_contains "codex update marker" "$CODEX_PROJECT_DIR/.codex/skills/sdp/build/SKILL.md"

echo "install-project regression checks passed"
