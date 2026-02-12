// Package git provides git operation wrappers with session validation.
package git

import "strings"

// CommandCategory represents the type of git command.
type CommandCategory int

const (
	// CategorySafe commands are read-only (status, log, diff, show)
	CategorySafe CommandCategory = iota
	// CategoryWrite commands modify the repository (add, commit, reset)
	CategoryWrite
	// CategoryRemote commands interact with remotes (push, fetch, pull)
	CategoryRemote
	// CategoryBranch commands change branches (checkout, branch, merge)
	CategoryBranch
)

// String returns the string representation of a command category.
func (c CommandCategory) String() string {
	switch c {
	case CategorySafe:
		return "safe"
	case CategoryWrite:
		return "write"
	case CategoryRemote:
		return "remote"
	case CategoryBranch:
		return "branch"
	default:
		return "unknown"
	}
}

// safeCommands lists commands that are read-only.
var safeCommands = map[string]bool{
	"status":    true,
	"log":       true,
	"diff":      true,
	"show":      true,
	"ls-files":  true,
	"rev-parse": true,
	"branch":    true,
	"remote":    true,
	"tag":       true,
}

// writeCommands lists commands that modify the repository.
var writeCommands = map[string]bool{
	"add":    true,
	"commit": true,
	"reset":  true,
	"rm":     true,
	"mv":     true,
	"stash":  true,
}

// remoteCommands lists commands that interact with remotes.
var remoteCommands = map[string]bool{
	"push":  true,
	"fetch": true,
	"pull":  true,
	"clone": true,
}

// branchChangeCommands lists commands that change branches.
var branchChangeCommands = map[string]bool{
	"checkout":     true,
	"switch":       true,
	"merge":        true,
	"rebase":       true,
	"cherry-pick":  true,
}

// CategorizeCommand determines the category of a git command.
func CategorizeCommand(cmd string) CommandCategory {
	cmd = strings.ToLower(cmd)

	if safeCommands[cmd] {
		return CategorySafe
	}
	if writeCommands[cmd] {
		return CategoryWrite
	}
	if remoteCommands[cmd] {
		return CategoryRemote
	}
	if branchChangeCommands[cmd] {
		return CategoryBranch
	}

	return CategorySafe
}

// NeedsSessionCheck returns true if the command requires session validation.
func NeedsSessionCheck(cmd string) bool {
	return true
}

// NeedsPostCheck returns true if the command needs post-execution validation.
func NeedsPostCheck(cmd string) bool {
	cmd = strings.ToLower(cmd)
	return writeCommands[cmd] || branchChangeCommands[cmd] || cmd == "push" || cmd == "pull"
}
