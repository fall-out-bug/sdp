# Telemetry: How It Works

## Overview

SDP telemetry is **LOCAL ONLY** — no data is sent to external servers.

## Data Collection

### Automatic Collection

Every command execution is automatically tracked:

```go
// In main.go PersistentPreRunE
telemetry.TrackCommandStart(cmd.Name(), args)

// In main.go PersistentPostRunE
telemetry.TrackCommandComplete(true, "")
```

### What Gets Recorded

**Command Start Event:**
```json
{
  "type": "command_start",
  "timestamp": "2026-02-06T12:00:00Z",
  "data": {
    "command": "doctor",
    "args": []
  }
}
```

**Command Complete Event:**
```json
{
  "type": "command_complete",
  "timestamp": "2026-02-06T12:00:01Z",
  "data": {
    "command": "doctor",
    "duration": 123400000000, // nanoseconds
    "success": true,
    "error": ""
  }
}
```

### Event Types

- `command_start` — Command execution started
- `command_complete` — Command finished (with success/failure)
- `ws_start` — Workstream execution started
- `ws_complete` — Workstream finished
- `quality_gate_result` — Quality gate passed/failed

## Storage

### Location

```
~/.sdp/telemetry.jsonl
```

**Format:** JSONL (one JSON object per line)

### Example File

```jsonl
{"type":"command_start","timestamp":"2026-02-06T10:30:00Z","data":{"command":"doctor","args":[]}}
{"type":"command_complete","timestamp":"2026-02-06T10:30:01Z","data":{"command":"doctor","duration":123400000000,"success":true,"error":""}}
{"type":"command_start","timestamp":"2026-02-06T10:31:00Z","data":{"command":"init","args":["."]}}
{"type":"command_complete","timestamp":"2026-02-06T10:31:05Z","data":{"command":"init","duration":5000000000,"success":true,"error":""}}
{"type":"command_start","timestamp":"2026-02-06T10:32:00Z","data":{"command":"build","args":["00-001-01"]}}
{"type":"command_complete","timestamp":"2026-02-06T10:32:45Z","data":{"command":"build","duration":45000000000,"success":false,"error":"test failed"}}
```

### File Permissions

```
-rw------- 1 user user 2.5K Feb  6 10:32 telemetry.jsonl
```

**Permissions:** `0600` (owner read/write only)

## Data Analysis

### 1. Check Status

```bash
$ sdp telemetry status

Telemetry Status:
  Enabled: No
  Events: 5
  File: /home/user/.sdp/telemetry.jsonl
```

### 2. Analyze Data

```bash
$ sdp telemetry analyze

📊 Telemetry Analysis Report
==========================
Total Events: 5

📈 Command Statistics:
  doctor:
    Total Runs: 1
    Success Rate: 100.0%
    Avg Duration: 123ms

  init:
    Total Runs: 1
    Success Rate: 100.0%
    Avg Duration: 5s

  build:
    Total Runs: 1
    Success Rate: 0.0%
    Avg Duration: 45s

❌ Top Errors:
  1. test failed (1 occurrences)
```

### 3. Upload for Sharing

```bash
# Package as JSON (structured format)
$ sdp telemetry upload --format json
✓ Collected 15 events
✓ Packaged into: telemetry_upload_2026-02-06.json
  Size: 2.5KB

🔒 Privacy Reminder:
  Review the file before sharing to ensure no sensitive data.

  You can now:
  - Attach to GitHub Issue
  - Send via email
  - Share for debugging

# Package as archive (tar.gz)
$ sdp telemetry upload --format archive
✓ Collected 15 events
✓ Packaged into: telemetry_upload_2026-02-06.tar.gz
  Size: 1.2KB
```

**Upload Package Structure:**
```json
{
  "metadata": {
    "version": "1.0",
    "generated_at": "2026-02-06T12:30:28Z",
    "event_count": 15,
    "format": "json"
  },
  "events": [
    {
      "type": "command_start",
      "timestamp": "2026-02-06T10:30:00Z",
      "data": {
        "command": "doctor",
        "args": []
      }
    }
  ]
}
```

### 4. Export Data

```bash
# Export as JSON
$ sdp telemetry export json
Exported telemetry to telemetry_export.json

# Export as CSV
$ sdp telemetry export csv
Exported telemetry to telemetry_export.csv
```

### 4. View Raw Data

```bash
# View last 10 events
$ tail -10 ~/.sdp/telemetry.jsonl

# Count events
$ wc -l ~/.sdp/telemetry.jsonl
```

## Usage Examples

### Example 1: Check What's Collected

```bash
# Run some commands
$ sdp doctor
$ sdp telemetry enable
$ sdp init .

# Check what was recorded
$ cat ~/.sdp/telemetry.jsonl | jq .

# Output:
{"type":"command_start",...}
{"type":"command_complete",...}
```

### Example 2: Analyze Your Usage Patterns

```bash
$ sdp telemetry analyze

📊 Command Usage:
  doctor: 15 times (30%)
  build: 20 times (40%)
  review: 10 times (20%)
  init: 5 times (10%)

📈 Success Rate:
  Overall: 85%
  doctor: 100%
  build: 75%
  review: 90%
```

### Example 3: Export and Share

```bash
# Export for analysis
$ sdp telemetry export json
$ cat telemetry_export.json

# Or review locally before sharing
$ less ~/.sdp/telemetry.jsonl
```

## Privacy Checklist

✅ **No PII collected** — Only command names, durations, success/failure
✅ **Local only** — File never leaves your machine
✅ **You control** — Enable/disable anytime
✅ **Transparent** — View your data anytime
✅ **Secure** — 0600 permissions (owner only)
✅ **Opt-in** — Disabled by default

## No Remote Transmission

**Important:** SDP does NOT have any code to:

- ❌ Upload data to servers
- ❌ Send analytics via HTTP
- ❌ Phone home to external services
- ❌ Integrate with third-party analytics

**All data stays on your machine.**

## Future: Voluntary Sharing

If you want to share telemetry (e.g., for debugging):

1. **Export your data:**
   ```bash
   sdp telemetry export my-telemetry.jsonl
   ```

2. **Review the file:**
   ```bash
   # Check for sensitive information
   cat my-telemetry.jsonl | jq '.data | keys'
   ```

3. **Share voluntarily:**
   - Attach to GitHub issue
   - Send to maintainer via email
   - Paste in chat for debugging

**We will ONLY use data you explicitly share.**

## Troubleshooting

### No Events Recorded

```bash
$ sdp telemetry status
Events: 0

# Check if enabled:
$ cat ~/.sdp/telemetry.json
# Should see: {"enabled":true}

# If disabled:
$ sdp telemetry enable
```

### File Too Large

```bash
# Check file size
$ ls -lh ~/.sdp/telemetry.jsonl
-rw------- 1 user user 2.5M Feb  6 10:32 telemetry.jsonl

# Clear old data (keeps last 90 days per policy)
$ sdp telemetry clear
```

### First Run No Prompt

```bash
# First run should prompt, but didn't?
$ rm ~/.sdp/telemetry.json
$ sdp doctor
# Should show consent prompt now
```

## Summary

**What:**
- Command usage tracking
- Execution duration
- Success/failure rates

**Where:**
- `~/.sdp/telemetry.jsonl` (local file)

**How:**
- Automatic on every command
- JSONL format (one line per event)
- 0600 permissions (secure)

**Control:**
- Opt-in (disabled by default)
- Enable/disable anytime
- View/export your data
- Clear when you want

**NOT:**
- ❌ No PII
- ❌ No remote transmission
- ❌ No third-party analytics
- ❌ No cloud services
