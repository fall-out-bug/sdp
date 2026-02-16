package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/errors"
	"github.com/fall-out-bug/sdp/internal/recovery"
	"github.com/spf13/cobra"
)

var diagnoseCmd = &cobra.Command{
	Use:   "diagnose [error-code]",
	Short: "Generate diagnostics report for error analysis",
	Long: `Generate a comprehensive diagnostics report for error analysis.

The diagnose command helps you understand and recover from failures by:
  - Classifying the error and providing context
  - Showing environment information
  - Providing recovery playbooks with step-by-step instructions
  - Suggesting next commands to run

If an error code is provided, shows the playbook for that error.
Without arguments, shows general diagnostic information.`,
	Example: `  # Diagnose a specific error code
  sdp diagnose ENV001

  # Show all error classes
  sdp diagnose --list-classes

  # Show all error codes
  sdp diagnose --list-codes

  # Generate report for last error (from log)
  sdp diagnose --last`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDiagnose,
}

var (
	diagnoseListClasses bool
	diagnoseListCodes   bool
	diagnoseLast        bool
	diagnoseJSON        bool
	diagnoseOutput      string
)

func init() {
	diagnoseCmd.Flags().BoolVar(&diagnoseListClasses, "list-classes", false, "List all error classes")
	diagnoseCmd.Flags().BoolVar(&diagnoseListCodes, "list-codes", false, "List all error codes")
	diagnoseCmd.Flags().BoolVar(&diagnoseLast, "last", false, "Diagnose last error from log")
	diagnoseCmd.Flags().BoolVar(&diagnoseJSON, "json", false, "Output in JSON format")
	diagnoseCmd.Flags().StringVarP(&diagnoseOutput, "output", "o", "", "Save report to file")
}

func runDiagnose(cmd *cobra.Command, args []string) error {
	// Handle list modes
	if diagnoseListClasses {
		return listErrorClasses()
	}
	if diagnoseListCodes {
		return listErrorCodes()
	}

	// Handle specific error code
	if len(args) > 0 {
		return diagnoseErrorCode(args[0])
	}

	// Handle last error from log
	if diagnoseLast {
		return diagnoseLastError()
	}

	// Default: show general help
	return showDiagnoseHelp()
}

func listErrorClasses() error {
	classes := []struct {
		class errors.ErrorClass
		desc  string
	}{
		{errors.ClassEnvironment, "Environment issues (tools, permissions, filesystem)"},
		{errors.ClassProtocol, "Protocol violations (invalid IDs, malformed files)"},
		{errors.ClassDependency, "Dependency problems (blocked workstreams, cycles)"},
		{errors.ClassValidation, "Validation failures (coverage, size, tests)"},
		{errors.ClassRuntime, "Runtime errors (network, timeout, unexpected state)"},
	}

	if diagnoseJSON {
		fmt.Println(`{"classes":[`)
		for i, c := range classes {
			comma := ""
			if i < len(classes)-1 {
				comma = ","
			}
			fmt.Printf(`  {"code":"%s","description":"%s"}%s
`, c.class, c.desc, comma)
		}
		fmt.Println(`]}`)
		return nil
	}

	fmt.Println("Error Classes:")
	fmt.Println()
	for _, c := range classes {
		fmt.Printf("  %-10s %s\n", c.class, c.desc)
	}
	fmt.Println()
	fmt.Println("Use 'sdp diagnose --list-codes' to see all error codes.")
	return nil
}

func listErrorCodes() error {
	codes := []struct {
		code errors.ErrorCode
		desc string
	}{
		// Environment
		{errors.ErrGitNotFound, "Git not found"},
		{errors.ErrGoNotFound, "Go not found"},
		{errors.ErrClaudeNotFound, "Claude Code CLI not found"},
		{errors.ErrBeadsNotFound, "Beads CLI not found"},
		{errors.ErrPermissionDenied, "Permission denied"},
		{errors.ErrWorktreeNotFound, "Git worktree not found"},
		{errors.ErrConfigNotFound, "Config file not found"},
		// Protocol
		{errors.ErrInvalidWorkstreamID, "Invalid workstream ID"},
		{errors.ErrInvalidFeatureID, "Invalid feature ID"},
		{errors.ErrMalformedYAML, "YAML parsing error"},
		{errors.ErrHashChainBroken, "Evidence hash chain broken"},
		{errors.ErrSessionCorrupted, "Session corrupted"},
		// Dependency
		{errors.ErrBlockedWorkstream, "Workstream blocked"},
		{errors.ErrCircularDependency, "Circular dependency"},
		{errors.ErrFeatureNotFound, "Feature not found"},
		{errors.ErrWorkstreamNotFound, "Workstream not found"},
		{errors.ErrCollisionDetected, "Scope collision detected"},
		// Validation
		{errors.ErrCoverageLow, "Coverage below threshold"},
		{errors.ErrFileTooLarge, "File too large"},
		{errors.ErrTestFailed, "Tests failed"},
		{errors.ErrLintFailed, "Linting failed"},
		{errors.ErrQualityGateFailed, "Quality gate failed"},
		{errors.ErrDriftDetected, "Drift detected"},
		// Runtime
		{errors.ErrCommandFailed, "Command failed"},
		{errors.ErrTimeoutExceeded, "Timeout exceeded"},
		{errors.ErrInternalError, "Internal error"},
	}

	if diagnoseJSON {
		fmt.Println(`{"codes":[`)
		for i, c := range codes {
			comma := ""
			if i < len(codes)-1 {
				comma = ","
			}
			fmt.Printf(`  {"code":"%s","class":"%s","description":"%s"}%s
`, c.code, c.code.Class(), c.desc, comma)
		}
		fmt.Println(`]}`)
		return nil
	}

	fmt.Println("Error Codes:")
	fmt.Println()

	var currentClass errors.ErrorClass
	for _, c := range codes {
		class := c.code.Class()
		if class != currentClass {
			currentClass = class
			fmt.Printf("\n[%s]\n", class)
		}
		fmt.Printf("  %-12s %s\n", c.code, c.desc)
	}

	fmt.Println()
	fmt.Println("Use 'sdp diagnose <code>' for detailed recovery steps.")
	return nil
}

func diagnoseErrorCode(codeStr string) error {
	code := errors.ErrorCode(codeStr)
	if !code.IsValid() {
		return fmt.Errorf("unknown error code: %s", codeStr)
	}

	pb := recovery.GetPlaybook(code)
	if pb == nil {
		return fmt.Errorf("no playbook found for error code: %s", codeStr)
	}

	if diagnoseJSON {
		data, err := json.MarshalIndent(pb, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal playbook: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Printf("== %s [%s] ==\n", pb.Title, pb.Severity)
	fmt.Printf("Error Code: %s\n", code)
	fmt.Printf("Class:      %s\n", code.Class())
	fmt.Printf("Message:    %s\n", code.Message())
	fmt.Println()
	fmt.Println(recovery.FormatPlaybook(pb))

	if diagnoseOutput != "" {
		data := fmt.Sprintf("# %s [%s]\n\n", pb.Title, pb.Severity)
		data += recovery.FormatPlaybook(pb)
		return os.WriteFile(diagnoseOutput, []byte(data), 0644)
	}

	return nil
}

func diagnoseLastError() error {
	// Find project root
	_, err := config.FindProjectRoot()
	if err != nil {
		return fmt.Errorf("find project root: %w", err)
	}

	// Generate a sample report for the most common error
	// In a real implementation, this would read from the evidence log
	fmt.Println("Note: --last requires an active error context.")
	fmt.Println("Use 'sdp diagnose <error-code>' to see playbooks for specific errors.")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  sdp diagnose ENV001  # Diagnose Git not found error")
	fmt.Println()

	return nil
}

func showDiagnoseHelp() error {
	fmt.Println("SDP Diagnostics Tool")
	fmt.Println()
	fmt.Println("Quick diagnosis and recovery for SDP failures.")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  sdp diagnose <code>       Show recovery playbook for error code")
	fmt.Println("  sdp diagnose --list-classes  List all error classes")
	fmt.Println("  sdp diagnose --list-codes    List all error codes")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  sdp diagnose ENV001       Diagnose Git not found error")
	fmt.Println("  sdp diagnose VAL001       Diagnose low coverage error")
	fmt.Println("  sdp diagnose --list-codes Show all error codes")
	fmt.Println()
	fmt.Println("Error classes:")
	fmt.Println("  ENV     - Environment issues (tools, permissions)")
	fmt.Println("  PROTO   - Protocol violations (invalid IDs, bad YAML)")
	fmt.Println("  DEP     - Dependency problems (blocked, cycles)")
	fmt.Println("  VAL     - Validation failures (coverage, tests)")
	fmt.Println("  RUNTIME - Runtime errors (network, timeout)")
	return nil
}
