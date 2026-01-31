---
ws_id: 00-900-03
feature: P0 Critical Fixes
status: completed
size: SMALL
github_issue: 3
title: P0-3: Resolve Markdown vs Beads Workflow Confusion
goal: Clarify which workflow to use and when to use each approach
acceptance_criteria:
  - [x] Decision document created (docs/workflow-decision.md)
  - [x] Recommendation: Beads-first (@feature) is PRIMARY
  - [x] Alternative: Traditional markdown (@idea) documented
  - [x] Decision matrix: when to use each approach
  - [x] Migration path from markdown to Beads-first
  - [x] FAQ included
  - [x] CLAUDE.md updated with reference
context: |
  Problem: Both Beads-first (@feature) and markdown (@idea) workflows documented,
  users confused about which approach to use.
  No clear recommendation or decision matrix.
steps: |
  1. Analyzed both workflows
  2. Created docs/workflow-decision.md (223 LOC) with:
     - Recommendation: Beads-first for new features
     - Detailed comparison table
     - Decision matrix (when to use which)
     - Migration path
     - FAQ
     - Future evolution timeline
  3. Updated CLAUDE.md with reference to decision doc
  4. Marked P0-3 as complete in deep analysis report
code_blocks: |
  # Beads-First Workflow (Recommended)

  **Entry Point:** `@feature "Feature description"`

  **When to Use:**
  - ✅ Default choice for all new features
  - ✅ Multi-person projects (shared task tracking)
  - ✅ Complex features (10+ workstreams)
  - ✅ Projects needing dependency tracking

  # Traditional Markdown Workflow (Alternative)

  **Entry Point:** `@idea "Feature description"`

  **When to Use:**
  - ⚠️ Solo development (single person)
  - ⚠️ Quick prototypes (1-3 workstreams)
  - ⚠️ Beads CLI not available
  - ⚠️ Exploratory work (may abandon)
dependencies: []
execution_report: |
  **Duration:** 2 hours
  **LOC Added:** 238 (workflow-decision.md + CLAUDE.md update)
  **Test Coverage:** N/A (documentation)
  **Deviations:** None
  **Status:** ✅ COMPLETE

  Created comprehensive workflow decision guide.
  Beads-first workflow recommended as PRIMARY approach.
  Traditional markdown documented as fallback option.
  Clear decision matrix provided for users.
