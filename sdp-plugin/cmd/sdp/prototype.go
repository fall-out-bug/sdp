package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/spf13/cobra"
)

func prototypeCmd() *cobra.Command {
	var featureID string
	var workstreamCount int
	var skipInterview bool
	var noDebt bool
	var immediate bool

	cmd := &cobra.Command{
		Use:   "prototype <feature-description>",
		Short: "Rapid prototyping shortcut for experienced vibecoders",
		Long: `Ultra-fast feature planning and execution for experienced developers.

Workflow (15 minutes):
1. Ultra-fast interview (5 questions)
2. Auto-generate 1-3 monolithic workstreams
3. Launch @oneshot with relaxed quality gates
4. Track tech debt for later cleanup

Quality Gates (Relaxed):
- TDD: Optional (tests after code)
- Coverage: No requirement
- File size: No limit
- Architecture: Monolithic OK
- Documentation: Comments only

This is PROTOTYPE MODE - working prototype over clean architecture.
All violations tracked as tech debt issues.

Examples:
  sdp prototype "Add user authentication"
  sdp prototype "Payment processing" --feature=F060
  sdp prototype "Dashboard widgets" --skip-interview
  sdp prototype "API refactor" --workstreams=2 --immediate`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			featureDesc := args[0]

			// Get feature ID or auto-generate
			if featureID == "" {
				// Simple auto-generation: use next available feature ID
				featureID = "AUTO"
			}

			fmt.Printf("🚀 PROTOTYPE MODE\n")
			fmt.Printf("Feature: %s\n", featureDesc)
			fmt.Printf("Feature ID: %s\n", featureID)
			fmt.Printf("\n")

			// Step 1: Ultra-fast interview (5 questions)
			if !skipInterview {
				fmt.Printf("📋 Step 1: Ultra-Fast Interview (5 questions)\n")
				fmt.Printf("⏱️  Estimated time: 5-10 minutes\n\n")

				// In real implementation, this would invoke @idea skill
				// with prototype-mode questions
				fmt.Printf("⚠️  Interview mode: @prototype skill not yet integrated\n")
				fmt.Printf("   Please run: @prototype \"%s\"\n", featureDesc)
				fmt.Printf("   Then return to continue with workstream generation\n")
				return nil
			}

			// Step 2: Auto-generate workstreams
			fmt.Printf("📊 Step 2: Generate Monolithic Workstreams\n")

			numWS := workstreamCount
			if numWS == 0 {
				// Auto-detect from interview or default to 1
				numWS = 1
			}

			fmt.Printf("   Generating %d workstream(s)...\n", numWS)

			for i := 1; i <= numWS; i++ {
				wsID := fmt.Sprintf("%s-%02d", featureID, i)
				fmt.Printf("   - %s: WS-%s\n", wsID, wsID)
			}
			// F056-03: emit generation event (non-blocking)
			evidence.Emit(evidence.SkillEvent("prototype", "generation", "00-000-00", map[string]interface{}{
				"workstream_count": numWS,
				"feature_id":       featureID,
			}))
			fmt.Printf("\n")

			// Step 3: Launch @oneshot
			fmt.Printf("🚀 Step 3: Launch @oneshot (Prototype Mode)\n")
			fmt.Printf("   Quality gates: RELAXED\n")
			fmt.Printf("   Tech debt tracking: %s\n", map[bool]string{true: "ENABLED", false: "DISABLED"}[!noDebt])
			fmt.Printf("\n")

			if !immediate {
				// Ask for confirmation
				fmt.Printf("Ready to launch agents? (Y/n): ")
				var response string
				_, err := fmt.Scanln(&response)
				if err != nil {
					return fmt.Errorf("failed to read response: %w", err)
				}
				if response != "Y" && response != "y" && response != "yes" {
					fmt.Printf("❌ Cancelled\n")
					return nil
				}
			}

			fmt.Printf("⚠️  PROTOTYPE MODE - Relaxed quality gates active\n")
			fmt.Printf("⚠️  Tech debt will be tracked but NOT blocking\n\n")

			// Create workstream files and invoke opencode/sdp-orchestrate
			projectRoot, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}
			resolvedFeatureID := resolveFeatureID(featureID)
			if err := createPrototypeWorkstreams(projectRoot, resolvedFeatureID, featureDesc, numWS); err != nil {
				return fmt.Errorf("create workstreams: %w", err)
			}
			if err := launchOneshot(projectRoot, resolvedFeatureID); err != nil {
				return err
			}

			// Step 4: Show what would happen
			fmt.Printf("📋 What Happens Next:\n\n")
			fmt.Printf("1. Agents execute workstreams with relaxed gates:\n")
			fmt.Printf("   ✅ TDD: Optional (tests after code)\n")
			fmt.Printf("   ✅ Coverage: No requirement\n")
			fmt.Printf("   ✅ File size: No limit\n")
			fmt.Printf("   ✅ Architecture: Monolithic OK\n")
			fmt.Printf("\n")
			fmt.Printf("2. Tech debt automatically tracked:\n")
			fmt.Printf("   📝 Architecture violations\n")
			fmt.Printf("   📝 Missing test coverage\n")
			fmt.Printf("   📝 Files >200 LOC\n")
			fmt.Printf("\n")
			fmt.Printf("3. Post-prototype:\n")
			fmt.Printf("   ⚠️  NOT PRODUCTION-READY\n")
			fmt.Printf("   🔧 Fix tech debt or refactor for production\n")
			fmt.Printf("\n")

			return nil
		},
	}

	cmd.Flags().StringVar(&featureID, "feature", "", "Explicit feature ID (default: auto-generate)")
	cmd.Flags().IntVar(&workstreamCount, "workstreams", 0, "Number of workstreams (1-3, default: auto-detect)")
	cmd.Flags().BoolVar(&skipInterview, "skip-interview", false, "Skip questions, use defaults")
	cmd.Flags().BoolVar(&noDebt, "no-debt", false, "Don't create tech debt issues")
	cmd.Flags().BoolVar(&immediate, "immediate", false, "Launch @oneshot without confirmation")

	return cmd
}

// resolveFeatureID converts AUTO to F099; otherwise normalizes (e.g. F060 -> F060).
func resolveFeatureID(id string) string {
	id = strings.ToUpper(strings.TrimSpace(id))
	if id != "" && id != "AUTO" {
		if !strings.HasPrefix(id, "F") {
			id = "F" + id
		}
		return id
	}
	return "F099" // F099 for prototype mode when AUTO
}

// createPrototypeWorkstreams writes minimal workstream files for the feature.
func createPrototypeWorkstreams(projectRoot, featureID, featureDesc string, numWS int) error {
	backlogDir := filepath.Join(projectRoot, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(backlogDir, 0755); err != nil {
		return err
	}
	rawNum := strings.TrimPrefix(featureID, "F")
	n, _ := strconv.Atoi(rawNum)
	if n <= 0 {
		n = 99
	}
	for i := 1; i <= numWS; i++ {
		wsID := fmt.Sprintf("00-%03d-%02d", n, i)
		content := fmt.Sprintf(`---
ws_id: %s
feature_id: %s
status: pending
scope_files: []
---

# %s: %s (Prototype)

## Goal

%s

## Scope Files

- (to be defined during build)

## Acceptance Criteria

- [ ] Working prototype
`, wsID, featureID, wsID, featureDesc, featureDesc)
		path := filepath.Join(backlogDir, wsID+".md")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}
	// Append to .beads-sdp-mapping.jsonl if it exists
	mappingPath := filepath.Join(projectRoot, ".beads-sdp-mapping.jsonl")
	f, err := os.OpenFile(mappingPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		for i := 1; i <= numWS; i++ {
			wsID := fmt.Sprintf("00-%03d-%02d", n, i)
			fmt.Fprintf(f, "%s\tprototype-%s\n", wsID, wsID)
		}
		f.Close()
	}
	return nil
}

// launchOneshot invokes sdp-orchestrate --runtime opencode or opencode run '@oneshot'.
func launchOneshot(projectRoot, featureID string) error {
	if !strings.HasPrefix(featureID, "F") {
		featureID = "F" + featureID
	}
	if os.Getenv("GLM_API_KEY") == "" {
		fmt.Printf("⚠️  GLM_API_KEY not set — LLM execution skipped\n")
		fmt.Printf("   Set GLM_API_KEY and run: @oneshot %s --mode=prototype\n", featureID)
		return nil
	}
	// Prefer sdp-orchestrate (from sdp_dev) if on PATH
	if _, err := exec.LookPath("sdp-orchestrate"); err == nil {
		cmd := exec.Command("sdp-orchestrate", "--feature", featureID, "--runtime", "opencode")
		cmd.Dir = projectRoot
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("sdp-orchestrate: %w", err)
		}
		return nil
	}
	// Fallback: opencode run '@oneshot FXXX --mode=prototype'
	if _, err := exec.LookPath("opencode"); err == nil {
		cmd := exec.Command("opencode", "run", fmt.Sprintf("@oneshot %s --mode=prototype", featureID), "--dir", projectRoot, "--continue")
		cmd.Dir = projectRoot
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("opencode run: %w", err)
		}
		return nil
	}
	fmt.Printf("⚠️  Neither sdp-orchestrate nor opencode on PATH\n")
	fmt.Printf("   Install sdp-orchestrate (sdp_dev) or opencode, set GLM_API_KEY, then run:\n")
	fmt.Printf("   @oneshot %s --mode=prototype\n", featureID)
	return nil
}
