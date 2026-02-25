# SDP Demo: Creating OpenCode Beads Viz Plugin
# Run: vhs demo.tape  (requires: brew install vhs)
# Output: docs/demo.gif

Output docs/demo.gif

Set FontSize 14
Set Width 900
Set Height 500
Set Padding 20

# Title
Type "echo '╔══════════════════════════════════════════════════════════════╗'"
Enter
Sleep 200ms
Type "echo '║  SDP: Creating OpenCode Plugin for Beads Visualization   ║'"
Enter
Sleep 200ms
Type "echo '╚══════════════════════════════════════════════════════════════╝'"
Enter
Sleep 500ms

# Step 1: Install
Type "echo ''"
Enter
Type "echo '\$ curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh'"
Enter
Sleep 300ms
Type "echo '🔗 Installing SDP project assets (prompts/hooks/config)...'"
Enter
Type "echo '✅ SDP installed to ./sdp'"
Enter
Sleep 800ms

# Step 2: Init
Type "echo ''"
Enter
Type "echo '\$ sdp init --auto'"
Enter
Sleep 200ms
Type "echo '✅ .sdp/config.yml created'"
Enter
Type "echo '✅ Guard rules ready'"
Enter
Sleep 600ms

# Step 3: Feature discovery
Type "echo ''"
Enter
Type "echo '\$ @feature \"OpenCode plugin for beads visualization\"'"
Enter
Sleep 400ms
Type "echo ''"
Enter
Type "echo '→ @idea: Requirements gathered'"
Enter
Type "echo '→ @design: Workstreams created'"
Enter
Sleep 300ms
Type "echo ''"
Enter
Type "echo '  docs/workstreams/backlog/00-001-01.md  # Plugin scaffold'"
Enter
Type "echo '  docs/workstreams/backlog/00-001-02.md  # Beads API client'"
Enter
Type "echo '  docs/workstreams/backlog/00-001-03.md  # Terminal UI'"
Enter
Sleep 800ms

# Step 4: Build
Type "echo ''"
Enter
Type "echo '\$ @oneshot F001'"
Enter
Sleep 300ms
Type "echo '→ Executing workstreams...'"
Enter
Type "echo '  ✅ 00-001-01  TDD pass'"
Enter
Type "echo '  ✅ 00-001-02  TDD pass'"
Enter
Type "echo '  ✅ 00-001-03  TDD pass'"
Enter
Sleep 600ms

# Step 5: Review & Deploy
Type "echo ''"
Enter
Type "echo '\$ @review F001'"
Enter
Sleep 200ms
Type "echo '→ QA | Security | DevOps | SRE | TechLead | Docs'"
Enter
Type "echo '  Verdict: APPROVED'"
Enter
Sleep 400ms
Type "echo ''"
Enter
Type "echo '\$ @deploy F001'"
Enter
Sleep 200ms
Type "echo '→ PR merged to main'"
Enter
Sleep 500ms

# Beads viz output
Type "echo ''"
Enter
Type "echo '\$ bd viz'"
Enter
Sleep 300ms
Type "echo ''"
Enter
Type "echo '  ┌─────────────────────────────────────────┐'"
Enter
Type "echo '  │ BEADS                    ready: 3       │'"
Enter
Type "echo '  ├─────────────────────────────────────────┤'"
Enter
Type "echo '  │ sdp-abc1  Add auth          P1  in_prog │'"
Enter
Type "echo '  │ sdp-def2  Fix login bug     P2  ready   │'"
Enter
Type "echo '  │ sdp-ghi3  Beads viz plugin   P1  done   │'"
Enter
Type "echo '  └─────────────────────────────────────────┘'"
Enter
Sleep 1s
