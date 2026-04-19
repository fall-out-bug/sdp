package main

import "os"

var defaultSkillsDirCandidates = []string{
	".claude/skills",
	".cursor/skills",
	".opencode/skills",
	".codex/skills",
	".codex/skills/sdp",
}

// hasSkillFiles checks whether a directory contains at least one
// subdirectory with a SKILL.md file (the per-skill layout).
func hasSkillFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(dir + "/" + e.Name() + "/SKILL.md"); err == nil {
			return true
		}
	}
	return false
}

func resolveDefaultSkillsDir() string {
	for _, candidate := range defaultSkillsDirCandidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() && hasSkillFiles(candidate) {
			return candidate
		}
	}

	return ".claude/skills"
}
