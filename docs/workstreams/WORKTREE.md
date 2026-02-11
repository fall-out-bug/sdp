# F059 Worktree

**Location:** `/Users/fall_out_bug/projects/vibe_coding/sdp-F059`

**Branch:** `feature/F059`

**Purpose:** Isolated workspace for F059 Observability Bridge Design

**Feature:** sdp-pom6 | Workstreams: 00-059-01, 00-059-02

**Sync from dev:**
```bash
cd /Users/fall_out_bug/projects/vibe_coding/sdp-F059
git fetch origin dev
git merge origin/dev
```

**Push:**
```bash
git push -u origin feature/F059
```

**Cleanup (after merge):**
```bash
cd /Users/fall_out_bug/projects/vibe_coding/sdp
git worktree remove ../sdp-F059
git branch -d feature/F059
```
