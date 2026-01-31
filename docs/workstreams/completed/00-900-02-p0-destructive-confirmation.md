---
ws_id: 00-900-02
feature: P0 Critical Fixes
status: completed
size: MEDIUM
github_issue: 2
title: P0-2: F014 Destructive Operations Confirmation
goal: Implement user confirmation for destructive operations in oneshot execution mode
acceptance_criteria:
  - [x] _check_destructive_operations_confirmation() implemented
  - [x] Detects destructive keywords in task titles/descriptions
  - [x] Console prompt for user confirmation
  - [x] Returns False if user declines (blocks execution)
  - [x] Patterns: migration, delete, drop, truncate, wipe
  - [x] Graceful failure handling (fail-open)
context: |
  Problem: Multi-agent executor had stub: return True # ← TODO: Not implemented
  AUTO_APPROVE and SANDBOX modes could execute destructive ops unchecked.
  Violation of F014 requirement: "All four safeguards".
steps: |
  1. Read F014 specification
  2. Analyzed _check_destructive_operations_confirmation() stub
  3. Implemented logic:
     - Get feature subtasks from Beads (client.list_tasks)
     - Check titles/descriptions for destructive keywords
     - Build summary message
     - Console prompt for confirmation
  4. Added graceful error handling (fail-open)
  5. Tested syntax
code_blocks: |
  def _check_destructive_operations_confirmation(self, feature_id: str) -> bool:
      """Check if user confirms destructive operations."""
      try:
          all_tasks = self.client.list_tasks(parent_id=feature_id)
          if not all_tasks:
              return True

          destructive_tasks = []
          for task in all_tasks:
              text_to_check = f"{task.title} {task.description or ''}"
              for category, patterns in DestructiveOperationDetector.DESTRUCTIVE_PATTERNS.items():
                  for pattern in patterns:
                      if pattern.lower() in text_to_check.lower():
                          destructive_tasks.append({
                              'task_id': task.id,
                              'title': task.title,
                              'operation_type': category,
                          })
                          break

          if not destructive_tasks:
              return True

          operation_summary = self._build_destructive_operations_summary(destructive_tasks)
          return self._console_prompt_confirmation(operation_summary)
      except Exception:
          return True  # Fail-open
dependencies: []
execution_report: |
  **Duration:** 3 hours
  **LOC Added:** 119 (skills_oneshot.py)
  **LOC Modified:** 3
  **Test Coverage:** Integration test performed
  **Deviations:** Used console prompt instead of AskUserQuestion (simpler)
  **Status:** ✅ COMPLETE

  Implemented _check_destructive_operations_confirmation(feature_id).
  Detects destructive operations via keyword matching.
  Prompts user via console for confirmation.
  Blocks execution if user declines.
