# SDP Demo: Real OpenCode CLI workflow
# Run from sdp repo root: vhs demo/demo.tape
# Requires: opencode, go, curl. Output: demo/demo.gif

Output demo/demo.gif

Set FontSize 12
Set Width 900
Set Height 560
Set Padding 16

Require opencode
Require go
Require curl

Set Shell "bash"

# 1. Create project, install SDP (per README)
Type "mkdir -p demo-beads-viz && cd demo-beads-viz && git init -q"
Enter
Sleep 500ms

Type "curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | SDP_IDE=opencode sh"
Enter
Sleep 8s

# 2. Init SDP + Beads (bd init for @feature)
Type "sdp init --auto"
Enter
Sleep 3s
Type "bd init 2>/dev/null || true"
Enter
Sleep 2s

# 3. @feature via OpenCode CLI
Type "opencode run '@feature beads viz in opencode' --dir ."
Enter
Sleep 90s

# 4. @oneshot via OpenCode CLI
Type "opencode run '@oneshot F001' --dir . --continue"
Enter
Sleep 120s

# 5. @review via OpenCode CLI
Type "opencode run '@review F001' --dir . --continue"
Enter
Sleep 60s

# 6. Validation
Type "sdp verify 00-001-01"
Enter
Sleep 5s

# 7. Show result (beads viz or bd list)
Type "bd list 2>/dev/null || ls -la docs/workstreams/backlog/"
Enter
Sleep 3s
