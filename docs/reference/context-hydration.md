# Context Hydration Guarantees

`internal/orchestrate/Hydrate` builds `.sdp/context-packet.json` before `@build` and `@review` execution.

## Guarantees

- `Hydrate` fails fast when workstream file or `AGENTS.md` cannot be read.
- Dependency lookup failures are captured per dependency in the packet, not dropped silently.
- Drift collection failures are captured as explicit `ERROR:` text in `drift_status`.
- Packet writing stays atomic (`.tmp` + rename).

## Why This Matters

- Keeps prompt context deterministic and auditable.
- Prevents hidden context loss that can produce non-reproducible agent behavior.
- Preserves partial diagnostics for dependency/drift sources without bypassing required quality-gate input.

## Testability Contract

`RunBuildPhase` and `RunReviewPhase` accept an `LLMInvoker` interface. Passing `nil` uses the default opencode invoker.

- Production path: `DefaultLLMInvoker`.
- Tests: fake invokers can simulate output and exit codes without spawning subprocesses.
