# SDP CLI Reference

Run `sdp --help` for the full list. Key commands:

| Command | Purpose |
|---------|---------|
| `sdp init` | Initialize SDP in project |
| `sdp doctor` | Health check (Git, .claude/) |
| `sdp guard activate/deactivate` | Scope enforcement for @build |
| `sdp drift detect <ws-id>` | Detect doc-code drift |
| `sdp build <ws-id>` | Execute single workstream |
| `sdp apply` | Execute ready workstreams |
| `sdp plan` | Generate workstreams |
| `sdp log show` | Evidence log |
| `sdp status` | Project state |
| `sdp checkpoint` | Long-running feature checkpoints |
| `sdp completion` | Shell completion |

See [reference/skills.md](reference/skills.md) for skill catalog and [PROTOCOL.md](PROTOCOL.md) for full spec.
