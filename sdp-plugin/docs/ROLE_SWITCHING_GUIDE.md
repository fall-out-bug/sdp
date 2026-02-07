# Role Switching Guide

**Dynamic role switching for SDP agents**

## Overview

Role switching allows agents to dynamically change their persona and capabilities during execution, enabling flexible multi-agent workflows.

## Role Switching Interface

### RoleSwitcher

```go
type RoleSwitcher struct {
    currentRole string
    roleLoader  *RoleLoader
    messageRouter *MessageRouter
}

func (rs *RoleSwitcher) SwitchTo(agentID string, newRole string) error
func (rs *RoleSwitcher) GetCurrentRole(agentID string) (string, error)
func (rs *RoleSwitcher) ListAvailableRoles() []string
```

## Switching Workflow

```
┌──────────┐    ┌──────────┐    ┌──────────┐
│  Agent   │───→│ Switcher │───→│New Role  │
│  (old)   │    │          │    │ Prompt   │
└──────────┘    └──────────┘    └──────────┘
                      │
                      ▼
               ┌──────────┐
               │ Context  │
               │Preserved│
               └──────────┘
```

## Usage Example

```bash
# Start as builder
@role builder
@build WS-001

# Switch to security reviewer
@switch security-reviewer
@review WS-001

# Switch back to builder
@switch builder
@build WS-002
```

## Context Preservation

When switching roles, the following context is preserved:
- Workstream ID
- Feature context
- File changes
- Test results
- Messages exchanged

## Role States

1. **Active**: Currently executing
2. **Paused**: Temporarily suspended
3. **Switching**: In transition
4. **Terminated**: Execution complete

## Best Practices

- Switch roles at logical boundaries (workstream completion)
- Preserve context between switches
- Document role transitions in logs
- Validate role compatibility before switching
