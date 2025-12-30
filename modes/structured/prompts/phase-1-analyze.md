# Phase 1: Analyze

You are analyzing requirements for a new feature.

## Your Task

Read the feature request and create a clear specification document.

## Input

A feature request, bug report, user story, or business requirement.

## Output

Create `docs/specs/{feature-name}.md` with:

### 1. Overview
One or two sentences explaining what this feature does and why it matters.

### 2. Requirements
Numbered list of functional requirements:
- REQ-1: The system must...
- REQ-2: The system must...

### 3. User Stories
Who benefits and how:
- As a [type of user], I want [action] so that [benefit].

### 4. Acceptance Criteria
Testable conditions that define "done":
- [ ] When X happens, Y should result
- [ ] User can see Z after doing W

### 5. Out of Scope
Explicitly state what this feature does NOT include (prevents scope creep).

### 6. Open Questions
List anything that needs clarification before proceeding.

## Example Prompt

```
Analyze this feature request and create a specification:

"We need users to be able to reset their passwords via email.
They should click 'Forgot Password', enter email, receive a link,
and set a new password. The link should expire after 1 hour."

Save to docs/specs/password-reset.md
```

## Quality Checklist

Before finishing this phase:
- [ ] All requirements are testable (not vague)
- [ ] User stories cover all user types
- [ ] Acceptance criteria are specific
- [ ] Out of scope is clearly defined
- [ ] Open questions are listed (even if none)

## Tips

- Ask clarifying questions if requirements are ambiguous
- Consider edge cases (what if email doesn't exist?)
- Think about error scenarios (what if link expired?)
- Don't over-specify implementation details (that's Phase 2)
