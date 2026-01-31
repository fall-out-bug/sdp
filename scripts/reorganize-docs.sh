#!/bin/bash
# Reorganize docs/ by role-based structure
# Creates beginner/, reference/, and internals/ directories

set -e

DOCS_ROOT="$(cd "$(dirname "$0")/../docs" && pwd)"

echo "📁 Reorganizing documentation in $DOCS_ROOT"
echo "=========================================="

# Create directories
mkdir -p "$DOCS_ROOT/beginner"
mkdir -p "$DOCS_ROOT/reference"
mkdir -p "$DOCS_ROOT/internals"

# Move files to appropriate directories

echo ""
echo "📚 Moving files to beginner/..."
# Beginner docs - progressive learning
mv "$DOCS_ROOT/TUTORIAL.md" "$DOCS_ROOT/beginner/" 2>/dev/null || true
# Note: SITEMAP.md, PROJECT_MAP.md stay in root for navigation

echo ""
echo "📖 Moving files to reference/..."
# Reference docs - lookup materials
mv "$DOCS_ROOT/GLOSSARY.md" "$DOCS_ROOT/reference/" 2>/dev/null || true
mv "$DOCS_ROOT/PRINCIPLES.md" "$DOCS_ROOT/reference/" 2>/dev/null || true
mv "$DOCS_ROOT/quality-gate-quick-reference.md" "$DOCS_ROOT/reference/" 2>/dev/null || true
mv "$DOCS_ROOT/quality-gate-schema.md" "$DOCS_ROOT/reference/" 2>/dev/null || true
mv "$DOCS_ROOT/quality-gates.md" "$DOCS_ROOT/reference/" 2>/dev/null || true

echo ""
echo "🔧 Moving files to internals/..."
# Internal docs - maintainer/advanced
mv "$DOCS_ROOT/architecture-checking.md" "$DOCS_ROOT/internals/" 2>/dev/null || true
mv "$DOCS_ROOT/completion-protocol.md" "$DOCS_ROOT/internals/" 2>/dev/null || true
mv "$DOCS_ROOT/error_patterns.md" "$DOCS_ROOT/internals/" 2>/dev/null || true
mv "$DOCS_ROOT/model-mapping.md" "$DOCS_ROOT/internals/" 2>/dev/null || true
mv "$DOCS_ROOT/multi-ide-parity.md" "$DOCS_ROOT/internals/" 2>/dev/null || true
mv "$DOCS_ROOT/overview-for-leads.md" "$DOCS_ROOT/internals/" 2>/dev/null || true
mv "$DOCS_ROOT/troubleshooting.md" "$DOCS_ROOT/internals/" 2>/dev/null || true
mv "$DOCS_ROOT/two-stage-review.md" "$DOCS_ROOT/internals/" 2>/dev/null || true
mv "$DOCS_ROOT/verification-protocol.md" "$DOCS_ROOT/internals/" 2>/dev/null || true

# Workstream completion docs move to workstreams/
mv "$DOCS_ROOT/WS-DX-03-completion.md" "$DOCS_ROOT/workstreams/" 2>/dev/null || true
mv "$DOCS_ROOT/WS-REF-01-progress.md" "$DOCS_ROOT/workstreams/" 2>/dev/null || true

echo ""
echo "✅ Reorganization complete!"
echo ""
echo "New structure:"
echo "  docs/beginner/     - Progressive learning (TUTORIAL.md)"
echo "  docs/reference/    - Lookup docs (GLOSSARY.md, PRINCIPLES.md, quality gates)"
echo "  docs/internals/    - Maintainer docs (architecture, patterns, protocols)"
echo "  docs/runbooks/     - Preserved (unchanged)"
echo "  docs/workstreams/  - Preserved (unchanged)"
