# WS-011: @idea/@design/@oneshot Invocation

> **Workstream ID:** WS-011  
> **Feature:** F024: Unified Workflow  
> **Dependencies:** WS-010

## Goal

Enable orchestrator to invoke sub-skills (@idea, @design, @oneshot) and capture their results.

## Acceptance Criteria

### AC1: Skill Tool Integration
- [ ] Orchestrator can call Skill tool
- [ ] Sub-skill executed with correct parameters
- [ ] Results captured and saved

### AC2: @idea Invocation
- [ ] Orchestrator calls @idea for requirements gathering
- [ ] AskUserQuestion answers captured
- [ ] Results saved to feature context

### AC3: @design Invocation
- [ ] Orchestrator calls @design for workstream planning
- [ ] Workstream breakdown captured
- [ ] Creates docs/drafts/idea-{slug}.md

### AC4: @oneshot Invocation
- [ ] Orchestrator calls @oneshot for autonomous execution
- [ ] Agent spawned with checkpoint support
- [ ] Results aggregated

## Scope Files

### Files to Create
**internal/orchestrator/skill_invoker.go** (NEW)
```go
package orchestrator

type SkillInvoker struct {
    skillTool SkillTool
}

func (si *SkillInvoker) InvokeIdea(featureID string) (IdeaResult, error)
func (si *SkillInvoker) InvokeDesign(specPath string) (DesignResult, error)
func (si *SkillInvoker) InvokeOneshot(featureID string) error)
```

## Implementation Steps

1. Create SkillInvoker struct
2. Implement InvokeIdea() - calls Skill tool with @idea
3. Implement InvokeDesign() - calls Skill tool with @design
4. Implement InvokeOneshot() - calls Skill tool with @oneshot
5. Add error handling and result capture

## Definition of Done

- [ ] All 4 AC met
- [ ] Skill invocations working
- [ ] Results captured
- [ ] Tests ≥80% coverage

## Estimated Scope

- ~150 LOC implementation
- ~80 LOC tests
- Duration: 1.5 hours
- Size: SMALL
