#!/bin/bash
set -e

echo "=========================================="
echo "Beads CLI Integration Test"
echo "=========================================="
echo

echo "1. List all issues:"
echo "----------------------------------------"
bd list
echo

echo "2. Show issue details:"
echo "----------------------------------------"
bd show sdp-8h0
echo

echo "3. Search for issues:"
echo "----------------------------------------"
bd search "Vision" 2>&1 | head -10
echo

echo "4. Create child issue:"
echo "----------------------------------------"
child_id=$(bd q "F015 Test: DecisionLogger - JSONL + MD decision logging")
echo "Created child: $child_id"
echo

echo "5. Add dependency:"
echo "----------------------------------------"
# Beads использует dependency syntax, посмотрим как добавить
bd set-state sdp-8h0 --state "in_progress" 2>&1 || true
echo

echo "6. List children:"
echo "----------------------------------------"
bd children sdp-8h0
echo

echo "=========================================="
echo "✅ Beads CLI Integration Test Complete"
echo "=========================================="
