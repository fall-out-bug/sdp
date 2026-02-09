#!/usr/bin/env node

/**
 * Sync skills from .claude/skills/ to .opencode/commands/
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const PROJECT_ROOT = path.dirname(__dirname);
const SKILLS_DIR = path.join(PROJECT_ROOT, '.claude', 'skills');
const COMMANDS_DIR = path.join(PROJECT_ROOT, '.opencode', 'commands');

// Skills that should NOT have commands (internal/meta skills)
const SKIP_SKILLS = ['guard', 'think', 'tdd'];

// Map skill names to command names
const COMMAND_MAP = {
  'review': 'codereview'  // review -> codereview for consistency
};

function getSkillName(skillDir) {
  const skillFile = path.join(SKILLS_DIR, skillDir, 'SKILL.md');
  if (!fs.existsSync(skillFile)) return null;

  const content = fs.readFileSync(skillFile, 'utf-8');
  const match = content.match(/^name:\s*(.+)$/m);
  return match ? match[1].trim() : skillDir;
}

function getSkillDescription(skillDir) {
  const skillFile = path.join(SKILLS_DIR, skillDir, 'SKILL.md');
  if (!fs.existsSync(skillFile)) return null;

  const content = fs.readFileSync(skillFile, 'utf-8');
  const match = content.match(/^description:\s*(.+)$/m);
  return match ? match[1].trim() : null;
}

function getAgentForSkill(skillName) {
  // Map skills to appropriate opencode agents
  const agentMap = {
    'debug': 'planner',
    'bugfix': 'builder',
    'build': 'builder',
    'deploy': 'deployer',
    'design': 'planner',
    'hotfix': 'builder',
    'issue': 'planner',
    'oneshot': 'orchestrator',
    'test': 'builder',
    'codereview': 'reviewer'
  };
  return agentMap[skillName] || 'builder';
}

function generateCommand(skillDir) {
  const skillName = getSkillName(skillDir);
  const description = getSkillDescription(skillDir);
  const commandName = COMMAND_MAP[skillDir] || skillDir;
  const agent = getAgentForSkill(skillName);

  if (!description) {
    console.warn(`⚠️  No description found for ${skillDir}, skipping`);
    return null;
  }

  const commandPath = path.join(COMMANDS_DIR, `${commandName}.md`);

  const content = `---
description: ${description}
agent: ${agent}
---

# /${commandName} — ${skillName.charAt(0).toUpperCase() + skillName.slice(1)}

## Overview

This command implements the ${skillName} skill from the SDP workflow.

See \`/.claude/skills/${skillDir}/SKILL.md\` for complete documentation.

## Usage

\`\`\`bash
/${commandName} [arguments]
\`\`\`

## Implementation

The command delegates to the \`${skillDir}\` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: \`.claude/skills/${skillDir}/SKILL.md\`
- Agents: \`.claude/agents/${agent}.md\`
`;

  return { commandPath, content, commandName };
}

function sync() {
  console.log('🔄 Syncing skills to commands...');

  // Ensure commands directory exists
  if (!fs.existsSync(COMMANDS_DIR)) {
    fs.mkdirSync(COMMANDS_DIR, { recursive: true });
  }

  // Get all skill directories
  const skillDirs = fs.readdirSync(SKILLS_DIR, { withFileTypes: true })
    .filter(dirent => dirent.isDirectory())
    .map(dirent => dirent.name)
    .sort();

  console.log(`📁 Found ${skillDirs.length} skills`);

  let created = 0;
  let updated = 0;
  let skipped = 0;

  for (const skillDir of skillDirs) {
    if (SKIP_SKILLS.includes(skillDir)) {
      console.log(`⏭️  Skipping ${skillDir} (internal skill)`);
      skipped++;
      continue;
    }

    const result = generateCommand(skillDir);
    if (!result) continue;

    const { commandPath, content, commandName } = result;

    const exists = fs.existsSync(commandPath);
    const needsUpdate = !exists || fs.readFileSync(commandPath, 'utf-8') !== content;

    if (needsUpdate) {
      fs.writeFileSync(commandPath, content, 'utf-8');
      if (exists) {
        console.log(`✏️  Updated ${commandName}.md`);
        updated++;
      } else {
        console.log(`✅ Created ${commandName}.md`);
        created++;
      }
    } else {
      console.log(`✓ ${commandName}.md up to date`);
    }
  }

  console.log(`\n📊 Summary: ${created} created, ${updated} updated, ${skipped} skipped`);
}

// Run sync
sync();
