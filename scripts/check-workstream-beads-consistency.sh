#!/usr/bin/env bash
set -euo pipefail

ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
WORKSTREAM_DIR="${ROOT}/docs/workstreams/completed"
ISSUES_FILE="${ROOT}/.beads/issues.jsonl"

if [ ! -d "${WORKSTREAM_DIR}" ]; then
  echo "No completed workstreams directory found; skipping consistency check"
  exit 0
fi

if [ ! -f "${ISSUES_FILE}" ]; then
  echo "No .beads/issues.jsonl found; skipping consistency check"
  exit 0
fi

python3 - "${WORKSTREAM_DIR}" "${ISSUES_FILE}" <<'PY'
import json
import pathlib
import re
import sys

workstream_dir = pathlib.Path(sys.argv[1])
issues_file = pathlib.Path(sys.argv[2])

issues = []
for line in issues_file.read_text(encoding="utf-8").splitlines():
    line = line.strip()
    if not line:
        continue
    try:
        issues.append(json.loads(line))
    except json.JSONDecodeError as exc:
        print(f"ERROR: malformed JSON in {issues_file}: {exc}")
        sys.exit(1)

issues_by_ref = {}
for issue in issues:
    ext = issue.get("external_ref")
    if ext:
        issues_by_ref.setdefault(ext, []).append(issue)

mismatches = []
pattern = re.compile(r'^ws_id:\s*"?([0-9]{2}-[0-9]{3}-[0-9]{2})"?\s*$')

for ws_file in sorted(workstream_dir.glob("*.md")):
    ws_id = ws_file.stem
    for line in ws_file.read_text(encoding="utf-8", errors="ignore").splitlines():
        match = pattern.match(line.strip())
        if match:
            ws_id = match.group(1)
            break

    related = issues_by_ref.get(ws_id, [])
    if not related:
        continue

    not_closed = [i for i in related if i.get("status") != "closed"]
    if not_closed:
        mismatches.append((ws_id, ws_file.name, not_closed))

if mismatches:
    print("ERROR: completed workstream docs with non-closed beads issues:")
    for ws_id, filename, records in mismatches:
        states = ", ".join(f"{r.get('id')}:{r.get('status')}" for r in records)
        print(f"- {ws_id} ({filename}) -> {states}")
    print("Fix by closing the listed issues or moving the workstream doc out of completed/.")
    sys.exit(1)

print("OK: completed workstreams are consistent with beads issue status")
PY
