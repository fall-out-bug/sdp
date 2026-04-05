package main

import "os"

var defaultSkillsDirCandidates = []string{
	".claude/skills",
	".cursor/skills",
	".opencode/skills",
	".codex/skills/sdp",
}

func resolveDefaultSkillsDir() string {
	for _, candidate := range defaultSkillsDirCandidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate
		}
	}

	return ".claude/skills"
}
