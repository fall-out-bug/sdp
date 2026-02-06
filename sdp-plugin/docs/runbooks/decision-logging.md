# Decision Logging Runbooks

Incident response procedures for decision logging issues.

## Scenario 1: decisions.jsonl Corrupted

**Symptoms:**
```bash
sdp decisions list
# Output: Found X decision(s) ... [decisions truncated]
# Or: Parse errors in log output
```

**Diagnosis:**
```bash
# Check for malformed JSON
python3 -m json.tool docs/decisions/decisions.jsonl 2>&1 | head -20

# Count lines vs valid decisions
LINES=$(wc -l < docs/decisions/decisions.jsonl)
FOUND=$(sdp decisions list | grep "Found.*decision" | grep -o "[0-9]*" | head -1)
echo "Lines: $LINES, Valid: $FOUND"

# If LINES > FOUND + 1, corruption exists
```

**Resolution:**
```bash
# Step 1: Backup corrupted file
cp docs/decisions/decisions.jsonl docs/decisions/decisions.jsonl.corrupt

# Step 2: Extract valid decisions
python3 << 'PYTHON'
import json

valid = []
with open('docs/decisions/decisions.jsonl', 'r') as f:
    for line_num, line in enumerate(f, 1):
        try:
            valid.append(json.loads(line))
        except json.JSONDecodeError as e:
            print(f"Line {line_num}: CORRUPT - {e}")

# Write cleaned file
with open('docs/decisions/decisions.jsonl', 'w') as f:
    for d in valid:
        f.write(json.dumps(d) + '\n')

print(f"Recovered {len(valid)}/{line_num} decisions")
PYTHON

# Step 3: Verify
sdp decisions list
```

**Prevention:**
- Commit to git regularly
- Enable fsync() (already enabled)
- Monitor disk space

---

## Scenario 2: Disk Full During Log()

**Symptoms:**
```bash
sdp decisions log --question="..." --decision="..."
# Error: failed to write decision: no space left on device
```

**Diagnosis:**
```bash
# Check disk space
df -h

# Check file size
du -sh docs/decisions/decisions.jsonl

# Check available inodes
df -i
```

**Resolution:**
```bash
# Step 1: Free up space
# Remove old rotated files
rm -f docs/decisions/*.jsonl.*

# Or move to external storage
mv docs/decisions/*.jsonl.* /external/backup/

# Step 2: Retry log operation
sdp decisions log --question="..." --decision="..."
```

**Prevention:**
- Set up log rotation (already enabled at 10MB)
- Monitor disk usage: `df -h /`
- Alert at 80% disk usage

---

## Scenario 3: Permission Denied on File Write

**Symptoms:**
```bash
sdp decisions log --question="..." --decision="..."
# Error: failed to open decisions file: permission denied
```

**Diagnosis:**
```bash
# Check file permissions
ls -la docs/decisions/decisions.jsonl

# Check ownership
ls -la docs/decisions/

# Check current user
whoami
groups
```

**Resolution:**
```bash
# Fix file permissions
chmod 644 docs/decisions/decisions.jsonl

# Fix ownership (if wrong)
sudo chown $USER:$USER docs/decisions/decisions.jsonl

# Fix directory permissions
chmod 755 docs/decisions/

# Retry
sdp decisions log --question="..." --decision="..."
```

**Prevention:**
- Check umask: `umask` (should be 022)
- Don't run sdp as root
- Don't manually edit decisions.jsonl

---

## Scenario 4: LoadAll() Timeout / Slow

**Symptoms:**
```bash
sdp decisions list
# Takes >10 seconds to complete
# Or: Command hangs
```

**Diagnosis:**
```bash
# Check file size
ls -lh docs/decisions/decisions.jsonl

# Count decisions
wc -l docs/decisions/decisions.jsonl

# If >10K decisions, this is the issue
```

**Resolution:**
```bash
# Option 1: Use pagination (when available)
sdp decisions list --page 1 --per-page 100

# Option 2: Export to markdown instead
sdp decisions export
# View in editor with LSP

# Option 3: Rotate log to reduce size
rm -f docs/decisions/decisions.jsonl.old
mv docs/decisions/decisions.jsonl docs/decisions/decisions.jsonl.old
# Old decisions archived, start fresh
```

**Prevention:**
- Enable log rotation (already enabled)
- Monitor decision count weekly
- Archive old decisions periodically

---

## Scenario 5: High Decision Logging Rate

**Symptoms:**
```bash
# Log shows rapid-fire logging
[decision] Log: ... (many per second)
# File size growing quickly
```

**Diagnosis:**
```bash
# Check logging rate
watch -n 1 'wc -l docs/decisions/decisions.jsonl'

# Check who's logging
grep -r "sdp decisions log" . --include="*.py" --include="*.go"

# Check for infinite loop
# Look for repeated decisions in log
sdp decisions list | grep "Question:" | sort | uniq -c | sort -rn | head
```

**Resolution:**
```bash
# Step 1: Identify source code
# Search for sdp decisions log calls in automation

# Step 2: Add rate limiting (if automation)
# Example: Only log unique decisions
if ! sdp decisions search --query "$QUESTION"; then
    sdp decisions log --question="$QUESTION" --decision="$DECISION"
fi

# Step 3: Add throttling (if needed)
# Example: Max 10 decisions per minute
```

**Prevention:**
- Don't log in loops
- Log only at key decision points
- Review logging frequency in code reviews

---

## General Troubleshooting Steps

For any decision logging issue:

1. **Check file exists and is readable**
   ```bash
   ls -la docs/decisions/decisions.jsonl
   ```

2. **Check log output for errors**
   ```bash
   sdp decisions list 2>&1 | tee /tmp/decisions-debug.log
   grep -i error /tmp/decisions-debug.log
   ```

3. **Validate JSON format**
   ```bash
   python3 -m json.tool docs/decisions/decisions.jsonl > /dev/null
   ```

4. **Test with simple case**
   ```bash
   sdp decisions log --question="Test" --decision="Yes"
   sdp decisions list
   ```

5. **Check git history**
   ```bash
   git log -- docs/decisions/decisions.jsonl | head -20
   ```

6. **Escalate if unresolved**
   - P1: Data loss → Page on-call
   - P2: Degraded performance → Create issue
   - P3: Documentation needed → Create issue

## Contact

For issues not covered here:
- Create issue: `bd create --title="Decision logging: <issue>"`
- Escalate to: Platform team / SRE
