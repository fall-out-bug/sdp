#!/bin/bash
# Validate SDP tutorial quality and completeness

set -e

TUTORIAL_FILE="docs/beginner/00-quick-start.md"
PRACTICE_FILE="docs/beginner/tutorial-practice.py"
TESTS_FILE="docs/beginner/tutorial-tests.py"

echo "🔍 SDP Tutorial Validation"
echo "============================"
echo ""

# Check 1: Files exist
echo "Check 1: Files exist"
if [ -f "$TUTORIAL_FILE" ]; then
    echo "✅ Tutorial file exists"
else
    echo "❌ Tutorial file missing"
    exit 1
fi

if [ -f "$PRACTICE_FILE" ]; then
    echo "✅ Practice file exists"
else
    echo "❌ Practice file missing"
    exit 1
fi

if [ -f "$TESTS_FILE" ]; then
    echo "✅ Tests file exists"
else
    echo "❌ Tests file missing"
    exit 1
fi

echo ""

# Check 2: Content completeness
echo "Check 2: Content completeness"

# Count steps in tutorial
STEP_COUNT=$(grep -c "^## Step" "$TUTORIAL_FILE" || echo "0")
if [ "$STEP_COUNT" -ge 6 ]; then
    echo "✅ Tutorial has $STEP_COUNT steps (≥6)"
else
    echo "❌ Tutorial has only $STEP_COUNT steps (need ≥6)"
    exit 1
fi

# Check for prerequisites section
if grep -q "## Prerequisites" "$TUTORIAL_FILE"; then
    echo "✅ Prerequisites section present"
else
    echo "❌ Prerequisites section missing"
    exit 1
fi

# Check for troubleshooting section
if grep -q "## Troubleshooting" "$TUTORIAL_FILE"; then
    echo "✅ Troubleshooting section present"
else
    echo "❌ Troubleshooting section missing"
    exit 1
fi

echo ""

# Check 3: Code examples
echo "Check 3: Code examples quality"

# Count code blocks
CODE_BLOCKS=$(grep -c '```' "$TUTORIAL_FILE" || echo "0")
CODE_BLOCK_COUNT=$((CODE_BLOCKS / 2))  # Each block has opening and closing
if [ "$CODE_BLOCK_COUNT" -ge 20 ]; then
    echo "✅ $CODE_BLOCK_COUNT code blocks (≥20)"
else
    echo "⚠️  Only $CODE_BLOCK_COUNT code blocks (recommend ≥20)"
fi

# Check for expected output examples
EXPECTED_OUTPUT=$(grep -c "Expected output:" "$TUTORIAL_FILE" || echo "0")
if [ "$EXPECTED_OUTPUT" -ge 6 ]; then
    echo "✅ $EXPECTED_OUTPUT expected output examples (≥6)"
else
    echo "⚠️  Only $EXPECTED_OUTPUT expected output examples (recommend ≥6)"
fi

echo ""

# Check 4: Time estimates
echo "Check 4: Time estimates"

TIME_ESTIMATES=$(grep -c "minutes" "$TUTORIAL_FILE" || echo "0")
if [ "$TIME_ESTIMATES" -ge 6 ]; then
    echo "✅ $TIME_ESTIMATES time estimates found (≥6)"
else
    echo "⚠️  Only $TIME_ESTIMATES time estimates (recommend ≥6)"
fi

# Check total time adds up to ~15 minutes
if grep -q "15 Minutes" "$TUTORIAL_FILE"; then
    echo "✅ Total time advertised as 15 minutes"
else
    echo "⚠️  Total time not clearly advertised"
fi

echo ""

# Check 5: Checkpoint validation
echo "Check 5: Checkpoint markers"

CHECKPOINTS=$(grep -c "✅ Checkpoint:" "$TUTORIAL_FILE" || echo "0")
if [ "$CHECKPOINTS" -ge 6 ]; then
    echo "✅ $CHECKPOINTS checkpoint markers (≥6)"
else
    echo "⚠️  Only $CHECKPOINTS checkpoint markers (recommend ≥6)"
fi

echo ""

# Check 6: Practice file quality
echo "Check 6: Practice file quality"

if [ -f "$PRACTICE_FILE" ]; then
    # Check for type hints
    TYPE_HINTS=$(grep -c "def.*->" "$PRACTICE_FILE" || echo "0")
    if [ "$TYPE_HINTS" -ge 3 ]; then
        echo "✅ Practice file has $TYPE_HINTS functions with type hints"
    else
        echo "⚠️  Practice file has only $TYPE_HINTS type hints (recommend ≥3)"
    fi

    # Check for docstrings
    DOCSTRINGS=$(grep -c '"""' "$PRACTICE_FILE" || echo "0")
    DOCSTRING_COUNT=$((DOCSTRINGS / 2))
    if [ "$DOCSTRING_COUNT" -ge 3 ]; then
        echo "✅ Practice file has $DOCSTRING_COUNT docstrings (≥3)"
    else
        echo "⚠️  Practice file has only $DOCSTRING_COUNT docstrings (recommend ≥3)"
    fi

    # Check for examples in docstrings
    if grep -q ">>>" "$PRACTICE_FILE"; then
        echo "✅ Practice file has doctest examples"
    else
        echo "⚠️  Practice file missing doctest examples"
    fi
fi

echo ""

# Check 7: Tests file quality
echo "Check 7: Tests file quality"

if [ -f "$TESTS_FILE" ]; then
    # Count test functions
    TEST_FUNCTIONS=$(grep -c "def test_" "$TESTS_FILE" || echo "0")
    if [ "$TEST_FUNCTIONS" -ge 10 ]; then
        echo "✅ Tests file has $TEST_FUNCTIONS test functions (≥10)"
    else
        echo "⚠️  Tests file has only $TEST_FUNCTIONS test functions (recommend ≥10)"
    fi

    # Check for test classes
    TEST_CLASSES=$(grep -c "^class Test" "$TESTS_FILE" || echo "0")
    if [ "$TEST_CLASSES" -ge 3 ]; then
        echo "✅ Tests file has $TEST_CLASSES test classes (≥3)"
    else
        echo "⚠️  Tests file has only $TEST_CLASSES test classes (recommend ≥3)"
    fi
fi

echo ""

# Check 8: Readability
echo "Check 8: Readability metrics"

# Count words in tutorial
WORD_COUNT=$(wc -w < "$TUTORIAL_FILE" | awk '{print $1}')
if [ "$WORD_COUNT" -ge 1500 ] && [ "$WORD_COUNT" -le 3000 ]; then
    echo "✅ Tutorial word count: $WORD_COUNT (1500-3000 range)"
elif [ "$WORD_COUNT" -lt 1500 ]; then
    echo "⚠️  Tutorial word count: $WORD_COUNT (<1500, might be too brief)"
else
    echo "⚠️  Tutorial word count: $WORD_COUNT (>3000, might be too long)"
fi

# Check for clear section headers
HEADER_COUNT=$(grep -c "^#" "$TUTORIAL_FILE" || echo "0")
echo "✅ Tutorial has $HEADER_COUNT section headers"

echo ""

# Check 9: Interactive elements
echo "Check 9: Interactive elements"

# Check for questions/prompts
QUESTIONS=$(grep -c "❓" "$TUTORIAL_FILE" || echo "0")
if [ "$QUESTIONS" -gt 0 ]; then
    echo "✅ Tutorial has $QUESTIONS interactive questions"
else
    echo "⚠️  Tutorial missing interactive questions"
fi

# Check for emojis (visual markers)
EMOJIS=$(grep -oE "🎯|📋|🔨|✅|❌|⚠️|🎉|📖|💬|🐛|📧|🚀" "$TUTORIAL_FILE" | wc -l)
if [ "$EMOJIS" -ge 20 ]; then
    echo "✅ Tutorial has $EMOJIS visual markers (≥20)"
else
    echo "⚠️  Tutorial has only $EMOJIS visual markers (recommend ≥20)"
fi

echo ""

# Check 10: Troubleshooting coverage
echo "Check 10: Troubleshooting coverage"

TROUBLESHOOTING_ISSUES=$(grep -c "^### Issue" "$TUTORIAL_FILE" || echo "0")
if [ "$TROUBLESHOOTING_ISSUES" -ge 5 ]; then
    echo "✅ Troubleshooting covers $TROUBLESHOOTING_ISSUES issues (≥5)"
else
    echo "⚠️  Troubleshooting covers only $TROUBLESHOOTING_ISSUES issues (recommend ≥5)"
fi

echo ""
echo "============================"
echo "✅ Tutorial validation complete!"
echo ""
echo "Summary:"
echo "  - Steps: $STEP_COUNT (≥6)"
echo "  - Code blocks: $CODE_BLOCK_COUNT (≥20)"
echo "  - Time estimates: $TIME_ESTIMATES (≥6)"
echo "  - Checkpoints: $CHECKPOINTS (≥6)"
echo "  - Troubleshooting: $TROUBLESHOOTING_ISSUES (≥5)"
echo "  - Word count: $WORD_COUNT (target: 1500-3000)"
echo ""
echo "🎉 Tutorial meets quality standards!"
