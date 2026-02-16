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

# Try to use the Go implementation if available
if command -v sdp >/dev/null 2>&1; then
    echo "Using sdp CLI to install hooks..."
    exec sdp hooks install
fi

# Fallback to manual installation if sdp is not available
echo "SDP CLI not found. Using manual installation..."

HOOKS_DIR=".git/hooks"

# Check if in git repo
if [ ! -d ".git" ]; then
    echo "Error: Not in git repository root" >&2
    exit 1
fi

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
