#!/bin/bash
# SDP Workstream Artifact Validation
# Validates workstream files for required fields and format

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track errors and warnings
errors=0
warnings=0

# Function to validate a single workstream file
validate_workstream() {
    local ws_file="$1"

    # Check file exists
    if [ ! -f "$ws_file" ]; then
        echo -e "${RED}✗${NC} File not found: $ws_file"
        ((errors++))
        return 1
    fi

    echo -e "\n🔍 Validating: $ws_file"

    # Extract frontmatter (between --- lines)
    local frontmatter=$(sed -n '/^---$/,/^---$/{ /^---$/d; p; }' "$ws_file")

    # Required fields
    local required_fields=("ws_id" "feature" "status" "size" "goal" "AC" "context" "steps")
    local missing_fields=()

    for field in "${required_fields[@]}"; do
        if ! echo "$frontmatter" | grep -q "^${field}:"; then
            missing_fields+=("$field")
        fi
    done

    if [ ${#missing_fields[@]} -gt 0 ]; then
        echo -e "${RED}✗${NC} Missing required fields: ${missing_fields[*]}"
        ((errors++))
    else
        echo -e "${GREEN}✓${NC} All required fields present"
    fi

    # Optional fields (warn if missing)
    local optional_fields=("depends_on" "github_issue" "assignee")

    for field in "${optional_fields[@]}"; do
        if ! echo "$frontmatter" | grep -q "^${field}:"; then
            echo -e "${YELLOW}⚠${NC} Optional field missing: $field"
            ((warnings++))
        fi
    done

    # Validate AC format (should start with "- [ ]" or "- [x")
    local ac_section=$(grep -A 20 "Acceptance Criteria" "$ws_file" || true)

    if [ -z "$ac_section" ]; then
        echo -e "${YELLOW}⚠${NC} No Acceptance Criteria section found"
        ((warnings++))
    else
        # Check if AC items are properly formatted
        if ! echo "$ac_section" | grep -q "^- \[.\]"; then
            echo -e "${RED}✗${NC} AC items must start with '- [ ]' or '- [x]'"
            ((errors++))
        else
            echo -e "${GREEN}✓${NC} AC format valid"
        fi
    fi

    # Check for Execution Report if status is completed
    local status=$(echo "$frontmatter" | grep "^status:" | cut -d':' -f2 | tr -d ' ')

    if [ "$status" = "completed" ]; then
        if grep -q "## Execution Report" "$ws_file"; then
            echo -e "${GREEN}✓${NC} Execution Report present"
        else
            echo -e "${RED}✗${NC} Completed WS must have Execution Report"
            ((errors++))
        fi
    fi

    # Validate status values (skip if it looks like a template with |)
    local valid_statuses=("backlog" "active" "completed" "blocked")
    local status_valid=false

    # Skip validation if it contains | (template placeholder)
    if [[ ! "$status" == *"|"* ]]; then
        for valid_status in "${valid_statuses[@]}"; do
            if [ "$status" = "$valid_status" ]; then
                status_valid=true
                break
            fi
        done

        if [ "$status_valid" = false ] && [ -n "$status" ]; then
            echo -e "${RED}✗${NC} Invalid status: $status (must be: ${valid_statuses[*]})"
            ((errors++))
        fi
    fi

    # Validate size values (skip if it looks like a template with |)
    local size=$(echo "$frontmatter" | grep "^size:" | cut -d':' -f2 | tr -d ' ')
    local valid_sizes=("SMALL" "MEDIUM" "LARGE")
    local size_valid=false

    # Skip validation if it contains | (template placeholder)
    if [[ ! "$size" == *"|"* ]]; then
        for valid_size in "${valid_sizes[@]}"; do
            if [ "$size" = "$valid_size" ]; then
                size_valid=true
                break
            fi
        done

        if [ "$size_valid" = false ] && [ -n "$size" ]; then
            echo -e "${RED}✗${NC} Invalid size: $size (must be: ${valid_sizes[*]})"
            ((errors++))
        fi
    fi
}

# Main validation logic
main() {
    echo "🔍 SDP Workstream Artifact Validation"
    echo "======================================"

    # If arguments provided, validate those files
    if [ $# -gt 0 ]; then
        for ws_file in "$@"; do
            validate_workstream "$ws_file"
        done
    else
        # No arguments: validate all workstream files in common locations
        local ws_dirs=(
            "docs/workstreams/backlog"
            "docs/workstreams/in_progress"
            "docs/workstreams/completed"
        )

        for ws_dir in "${ws_dirs[@]}"; do
            if [ -d "$ws_dir" ]; then
                for ws_file in "$ws_dir"/*.md; do
                    if [ -f "$ws_file" ]; then
                        validate_workstream "$ws_file"
                    fi
                done
            fi
        done
    fi

    # Summary
    echo ""
    echo "======================================"
    echo "Validation Summary:"
    echo -e "  Errors: ${RED}${errors}${NC}"
    echo -e "  Warnings: ${YELLOW}${warnings}${NC}"

    if [ $errors -eq 0 ]; then
        echo -e "\n${GREEN}✓ Artifact validation passed${NC}"
        exit 0
    else
        echo -e "\n${RED}✗ Artifact validation failed with $errors error(s)${NC}"
        exit 1
    fi
}

# Run main function
main "$@"
