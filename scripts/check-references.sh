#!/bin/sh
# check-references.sh — Reference Integrity Gate for SDP
#
# Validates that all skill/command/agent references across the codebase
# resolve to actual files. Exit 1 on any broken reference.
#
# Checks:
#   1. Skills mentioned in CLAUDE.md exist in prompts/skills/
#   2. Commands in .claude/commands.json map to existing skill files
#   3. Patterns in .claude/commands.json map to existing pattern files
#   4. Agents in .claude/commands.json map to existing agent files
#   5. Harness READMEs (.cursor/README.md, .codex/INSTALL.md,
#      .opencode/README.md) reference existing skills
#   6. All symlinks resolve correctly
#
# Requirements:
#   - GNU grep (for grep -oE extended regex). Ubuntu-latest CI ships GNU grep.
#   - POSIX sh, find, sed, readlink.
#
# Usage:
#   ./scripts/check-references.sh          # from sdp/ root
#   ./scripts/check-references.sh /path    # explicit root

# --- Resolve SDP root ---
if [ -n "$1" ]; then
    SDP_ROOT="$1"
else
    SDP_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
fi

if [ ! -d "$SDP_ROOT" ]; then
    echo "ERROR: SDP root does not exist: $SDP_ROOT" >&2
    exit 1
fi

if [ ! -d "$SDP_ROOT/prompts/skills" ]; then
    echo "ERROR: $SDP_ROOT does not appear to be an SDP repo (missing prompts/skills/)" >&2
    exit 1
fi

ERRORS=0
WARNINGS=0

# Temp file for symlink collection — cleaned up on EXIT/INT/TERM
_SYMLINKS_TMP=""
_cleanup() {
    [ -n "$_SYMLINKS_TMP" ] && rm -f "$_SYMLINKS_TMP" 2>/dev/null
}
trap _cleanup EXIT INT TERM

# --- Helpers ---
log_error() {
    printf "ERROR: %s\n" "$1" >&2
    ERRORS=$((ERRORS + 1))
}

log_warn() {
    printf "WARN: %s\n" "$1" >&2
    WARNINGS=$((WARNINGS + 1))
}

log_ok() {
    printf "  ok: %s\n" "$1"
}

skill_file_exists() {
    _name="$1"
    [ -f "${SDP_ROOT}/prompts/skills/${_name}/SKILL.md" ]
}

# --- Preamble ---
printf "%s\n" "=== SDP Reference Integrity Check ==="
printf "Root: %s\n\n" "$SDP_ROOT"

# ============================================================
# 1. Skills mentioned in CLAUDE.md
# ============================================================
printf "%s\n" "--- Checking CLAUDE.md skill references ---"

CLAUDE_MD="${SDP_ROOT}/CLAUDE.md"
if [ ! -f "$CLAUDE_MD" ]; then
    log_warn "CLAUDE.md not found at ${CLAUDE_MD}"
else
    # Extract @command names from the "Commands:" line
    # Format: **Commands:** @vision @reality @feature @oneshot @build @review @deploy
    COMMANDS_LINE=$(grep -E '^\*\*Commands:\*\*' "$CLAUDE_MD" 2>/dev/null || true)

    if [ -n "$COMMANDS_LINE" ]; then
        # Parse @xxx tokens
        for token in $COMMANDS_LINE; do
            case "$token" in
                @*)
                    skill_name="${token#@}"
                    # Strip trailing punctuation
                    skill_name=$(printf '%s' "$skill_name" | sed 's/[^a-zA-Z0-9_-]//g')
                    if [ -z "$skill_name" ]; then
                        continue
                    fi
                    if skill_file_exists "$skill_name"; then
                        log_ok "CLAUDE.md @${skill_name} -> prompts/skills/${skill_name}/SKILL.md"
                    else
                        log_error "CLAUDE.md references @${skill_name} but prompts/skills/${skill_name}/SKILL.md not found"
                    fi
                    ;;
            esac
        done
    else
        log_warn "No 'Commands:' line found in CLAUDE.md"
    fi
fi

# ============================================================
# 2. Commands in .claude/commands.json -> skill files
# ============================================================
printf "\n%s\n" "--- Checking .claude/commands.json command references ---"

COMMANDS_JSON="${SDP_ROOT}/.claude/commands.json"
if [ ! -f "$COMMANDS_JSON" ]; then
    log_warn ".claude/commands.json not found"
else
    # NOTE: JSON parsing uses grep+sed for zero-dependency POSIX compat.
    # Requires commands.json to stay pretty-printed (one value per line).
    # If JSON is ever minified, add jq or python3 as a CI dependency.
    # Extract "file" values from commands section
    # Using grep+sed for POSIX compatibility (no jq dependency)
    # Pattern: "file": "skills/xxx.md"
    file_refs=$(grep '"file"' "$COMMANDS_JSON" | sed 's/.*"file"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/')

    for ref in $file_refs; do
        case "$ref" in
            skills/*)
                skill_name=$(printf '%s' "$ref" | sed 's|skills/||; s|\.md$||')
                if skill_file_exists "$skill_name"; then
                    log_ok "commands.json ${ref} -> prompts/skills/${skill_name}/SKILL.md"
                else
                    log_error "commands.json references ${ref} but prompts/skills/${skill_name}/SKILL.md not found"
                fi
                ;;
            *)
                log_warn "commands.json: unexpected file reference format: ${ref}"
                ;;
        esac
    done
fi

# ============================================================
# 3. Patterns in .claude/commands.json -> pattern files
# ============================================================
printf "\n%s\n" "--- Checking .claude/commands.json pattern references ---"

if [ -f "$COMMANDS_JSON" ]; then
    # Extract pattern value lines: "key": "patterns/xxx.md"
    pattern_refs=$(grep '"patterns/' "$COMMANDS_JSON" | sed 's/.*"\([^"]*patterns\/[^"]*\)".*/\1/')

    for ref in $pattern_refs; do
        case "$ref" in
            patterns/*)
                pattern_path="${SDP_ROOT}/.claude/${ref}"
                if [ -f "$pattern_path" ]; then
                    log_ok "commands.json ${ref} -> .claude/${ref}"
                else
                    log_error "commands.json references ${ref} but .claude/${ref} not found"
                fi
                ;;
            *)
                log_warn "commands.json: unexpected pattern reference format: ${ref}"
                ;;
        esac
    done
fi

# ============================================================
# 4. Agents in .claude/commands.json -> agent files
# ============================================================
printf "\n%s\n" "--- Checking .claude/commands.json agent references ---"

if [ -f "$COMMANDS_JSON" ]; then
    # Extract agent value lines: "key": "agents/xxx.md"
    agent_refs=$(grep '"agents/' "$COMMANDS_JSON" | sed 's/.*"\([^"]*agents\/[^"]*\)".*/\1/')

    for ref in $agent_refs; do
        case "$ref" in
            agents/*)
                agent_path="${SDP_ROOT}/prompts/${ref}"
                if [ -f "$agent_path" ]; then
                    log_ok "commands.json ${ref} -> prompts/${ref}"
                else
                    log_error "commands.json references ${ref} but prompts/${ref} not found"
                fi
                ;;
            *)
                log_warn "commands.json: unexpected agent reference format: ${ref}"
                ;;
        esac
    done
fi

# ============================================================
# 5. Harness READMEs reference existing skills
# ============================================================
printf "\n%s\n" "--- Checking harness README skill references ---"

# Known skill names — discovered dynamically from prompts/skills/*/SKILL.md
KNOWN_SKILLS=""
for _sk_dir in "${SDP_ROOT}"/prompts/skills/*/; do
    [ -f "${_sk_dir}SKILL.md" ] && KNOWN_SKILLS="${KNOWN_SKILLS} $(basename "${_sk_dir}")"
done
KNOWN_SKILLS="${KNOWN_SKILLS# }"

is_known_skill() {
    _s="$1"
    for k in $KNOWN_SKILLS; do
        [ "$_s" = "$k" ] && return 0
    done
    return 1
}

check_harness_readme() {
    _file="$1"
    _label="$2"

    if [ ! -f "$_file" ]; then
        log_warn "${_label} not found at ${_file}"
        return
    fi

    # Extract @xxx references from the file.
    # Match @skill in any context: start-of-line, after space, after [, after (, after `, after *
    for token in $(grep -oE '(^|[][ ()`*])@[a-zA-Z0-9_-]+' "$_file" 2>/dev/null | sed 's/^[^@]*//' || true); do
        skill_name="${token#@}"
        if is_known_skill "$skill_name"; then
            if skill_file_exists "$skill_name"; then
                log_ok "${_label} @${skill_name} -> prompts/skills/${skill_name}/SKILL.md"
            else
                log_error "${_label} references @${skill_name} but prompts/skills/${skill_name}/SKILL.md not found"
            fi
        fi
    done
}

check_harness_readme "${SDP_ROOT}/.cursor/README.md" ".cursor/README.md"
check_harness_readme "${SDP_ROOT}/.codex/INSTALL.md" ".codex/INSTALL.md"
check_harness_readme "${SDP_ROOT}/.opencode/README.md" ".opencode/README.md"

# Also check .codex/skills/README.md if it exists
check_harness_readme "${SDP_ROOT}/.codex/skills/README.md" ".codex/skills/README.md"

# ============================================================
# 6. All symlinks resolve correctly
# ============================================================
printf "\n%s\n" "--- Checking symlink integrity ---"

# Use a temp file to collect symlinks (avoids subshell variable scope issues)
_SYMLINKS_TMP="${TMPDIR:-/tmp}/sdp-check-refs-$$"
find "${SDP_ROOT}" -type l ! -path '*/.git/*' 2>/dev/null > "$_SYMLINKS_TMP" || true

while read -r link; do
    [ -z "$link" ] && continue
    if [ ! -e "$link" ]; then
        target=$(readlink "$link")
        log_error "Broken symlink: ${link##${SDP_ROOT}/} -> ${target}"
    else
        target=$(readlink "$link")
        log_ok "symlink: ${link##${SDP_ROOT}/} -> ${target}"
    fi
done < "$_SYMLINKS_TMP"

# ============================================================
# 7. llm_subagents in commands.json — logical names, NOT file refs
# ============================================================
printf "\n%s\n" "--- Checking .claude/commands.json llm_subagents references ---"

if [ -f "$COMMANDS_JSON" ]; then
    # llm_subagents are logical role names used at runtime by the LLM
    # (e.g. "analyst", "product-manager", "quality-reviewer", "documentation").
    # They do NOT map to files in prompts/agents/ and are resolved by the
    # orchestrator at spawn time.  This is intentional — do not validate them
    # as file paths.
    llm_names=$(grep -oE '"llm_subagents"[[:space:]]*:[[:space:]]*\[[^]]*\]' "$COMMANDS_JSON" | grep -oE '"[a-z-]+"' | tr -d '"' | sort -u)
    for name in $llm_names; do
        printf "  info: llm_subagent '%s' — logical name, not a file reference (intentional)\n" "$name"
    done
fi

# ============================================================
# 8. Harness symlink directories resolve to prompts/skills
# ============================================================
printf "\n%s\n" "--- Checking harness skill/agent symlinks ---"

for harness in .cursor .codex .opencode .claude; do
    for sub in skills agents; do
        link_path="${SDP_ROOT}/${harness}/${sub}"
        if [ -L "$link_path" ]; then
            target=$(readlink "$link_path")
            resolved="${SDP_ROOT}/${harness}/${target}"
            if [ -d "$resolved" ]; then
                log_ok "${harness}/${sub} -> ${target} (resolves)"
            else
                log_error "${harness}/${sub} -> ${target} (DOES NOT RESOLVE)"
            fi
        fi
    done
done

# ============================================================
# Summary
# ============================================================
printf "\n%s\n" "=== Summary ==="
printf "Errors:   %d\n" "$ERRORS"
printf "Warnings: %d\n" "$WARNINGS"

if [ "$ERRORS" -gt 0 ]; then
    printf "\nFAIL: %d broken reference(s) found.\n" "$ERRORS"
    exit 1
fi

printf "\nPASS: All references are intact.\n"
exit 0
