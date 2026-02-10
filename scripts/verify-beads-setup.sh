#!/bin/bash
# Script to verify Beads installation and setup

set -e

echo "🔍 Verifying Beads installation..."
echo ""

# Check Go
echo "1. Checking Go installation..."
if command -v go &> /dev/null; then
    GO_VERSION=$(go version)
    echo "   ✅ Go installed: $GO_VERSION"
else
    echo "   ❌ Go not found"
    echo ""
    echo "   To install Go:"
    echo "   brew install go"
    echo ""
    exit 1
fi

# Check Go version
echo ""
echo "2. Checking Go version..."
GO_MAJOR=$(go version | awk '{print $3}' | cut -d. -f1 | sed 's/go//')
if [ "$GO_MAJOR" -ge 1 ] && [ "$GO_MAJOR" -lt 2 ]; then
    echo "   ✅ Go version compatible (1.x)"
else
    echo "   ⚠️  Go version: $(go version)"
    echo "   Required: Go 1.24+"
fi

# Check Beads CLI
echo ""
echo "3. Checking Beads CLI..."
if command -v bd &> /dev/null; then
    echo "   ✅ Beads CLI installed"
    BD_VERSION=$(bd --version 2>/dev/null || echo "unknown")
    echo "   Version: $BD_VERSION"
else
    echo "   ❌ Beads CLI not found"
    echo ""
    echo "   To install Beads:"
    echo "   go install github.com/steveyegge/beads/cmd/bd@latest"
    echo ""
    exit 1
fi

# Check if .beads directory exists
echo ""
echo "4. Checking Beads initialization..."
if [ -d ".beads" ]; then
    echo "   ✅ Beads initialized (.beads/ exists)"

    # Check database
    if [ -f ".beads/beads.db" ]; then
        echo "   ✅ Beads database exists"
    else
        echo "   ⚠️  Beads database not found"
        echo "   Run: bd init"
    fi

    # Check JSONL
    if [ -f ".beads/issues.jsonl" ]; then
        ISSUES=$(wc -l < .beads/issues.jsonl)
        echo "   ✅ Beads issues.jsonl ($ISSUES issues)"
    fi
else
    echo "   ⚠️  Beads not initialized"
    echo ""
    echo "   To initialize Beads:"
    echo "   bd init"
    echo ""
fi

echo ""
echo "✅ Beads verification complete!"
echo ""
echo "Next steps:"
echo "  1. Set BEADS_USE_MOCK=false in environment"
echo "  2. Test with real Beads:"
echo "     bd ready   # list available work"
