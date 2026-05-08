#!/bin/sh
# Install SDP git hooks
#
# DEPRECATED: This script is deprecated in favor of the Go implementation.
# Use: sdp hooks install
# This script is kept for backward compatibility only.

set -e

echo "WARNING: This install method is deprecated."
echo "Please use: sdp hooks install"
echo ""

# Try to use the Go implementation if available.
# Resolve the CLI relative to the script/project, never from PATH
# (on macOS /usr/bin/sdp is an Xcode tool, not this project's CLI).
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

sdp_cli=""
# 1. Local build inside optional sdp/ checkout
if [ -x "${PROJECT_ROOT}/sdp/sdp-plugin/sdp" ]; then
    sdp_cli="${PROJECT_ROOT}/sdp/sdp-plugin/sdp"
fi
# 2. Project's own binary (if built at project root)
if [ -z "${sdp_cli}" ] && [ -x "${PROJECT_ROOT}/sdp" ]; then
    case "$(file "${PROJECT_ROOT}/sdp" 2>/dev/null)" in
        *"ELF"*|*"Mach-O"*|*"PE32"*)
            sdp_cli="${PROJECT_ROOT}/sdp"
            ;;
    esac
fi

if [ -n "${sdp_cli}" ]; then
    echo "Using sdp CLI (${sdp_cli}) to install hooks..."
    exec "${sdp_cli}" hooks install
fi

# Fallback to manual installation if sdp is not available
echo "SDP CLI not found. Using manual installation..."

# Resolve hooks dir via git (works in both regular repos and worktrees
# where .git is a file, not a directory).
if ! GIT_DIR="$(git rev-parse --git-dir 2>/dev/null)"; then
    echo "Error: Not in a git repository" >&2
    exit 1
fi
HOOKS_DIR="${GIT_DIR}/hooks"

# Create hooks directory if needed
mkdir -p "${HOOKS_DIR}"

# SDP marker for ownership tracking
SDP_MARKER="# SDP-MANAGED-HOOK"

# Install canonical hooks
for hook in pre-commit pre-push post-merge post-checkout; do
    hook_file="${HOOKS_DIR}/${hook}"

    # Check if hook already exists and is SDP-managed
    if [ -f "${hook_file}" ]; then
        if grep -q "${SDP_MARKER}" "${hook_file}" 2>/dev/null; then
            echo "Updating ${hook}..."
        else
            echo "Warning: ${hook} exists but not SDP-managed, skipping"
            continue
        fi
    else
        echo "Installing ${hook}..."
    fi

    # Write hook with SDP marker
    cat > "${hook_file}" << EOF
#!/bin/sh
${SDP_MARKER}
# SDP Git Hook
# This hook is managed by SDP. Do not edit manually.

# Check if sdp binary exists
if ! command -v sdp >/dev/null 2>&1; then
    echo "Warning: SDP CLI (sdp) not found in PATH"
    echo "Install SDP to enable quality checks: https://github.com/fall-out-bug/sdp"
    exit 0
fi

# Add your custom checks here
EOF

    chmod +x "${hook_file}"
done

echo ""
echo "Git hooks installed!"
echo "Note: Install the SDP CLI for full functionality: sdp hooks install"
