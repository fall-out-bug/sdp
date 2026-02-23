---
name: submit-to-swarm
description: Submit a task to the SDP swarm via Intake Gateway
args:
  - name: project
    description: Project ID (from registry)
    required: true
  - name: title
    description: Task title
    required: true
---

Submit task to swarm.

Calls `POST /api/v1/intake` on the Intake Gateway with:
- project_id: {{project}}
- title: {{title}}
- source: opencode

Set INTAKE_GATEWAY_URL (default http://localhost:8081) for the gateway base URL.

Example: `/swarm sdp_dev "Add user authentication"`
