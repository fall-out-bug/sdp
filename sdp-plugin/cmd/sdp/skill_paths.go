package main

import "os"

var defaultSkillsDirCandidates = []string{
	".claude/skills",
	".cursor/skills",
	".opencode/skills",
	".codex/skills",
	".codex/skills/sdp",
}

func resolveDefaultSkillsDir() string {
	for _, candidate := range defaultSkillsDirCandidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			// For .codex/skills, prefer per-skill layout if it has SKILL.md subdirs.
			// Fall back to .codex/skills/sdp for backward compat with old installs.
			return candidate
		}
	}

	return ".claude/skills"
}
