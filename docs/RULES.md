# Developer Rules

This file outlines rules to follow during development of this application. These include code practices, directory structure to follow, test cases to write etc.

1. Do not hardcode constant values in code. Store them in `internal/config/`.
2. Write minimum to no comments in code. Code should itself be readable enough to not even require comments. 
3. Code goes in its appropriate places. Do not mix up code between files, or directories. For info on project structure read [ARCHITECTURE.md](ARCHITECTURE.md)
4. After every release, add bugs to fix, features to add, major changes to do in [DEV.md](DEV.md), and checkmark ones that are done along the way. After finishing, they move to Changelog prior to release.
5. Every commit should be logically grouped, with one-liner commit messages and no co-authored trailer part if committed by AI.
6. Every behavioral change needs relevant tests. Prefer deterministic unit tests with mocks over live network, real Wails windows, or user-disk access.
7. Before committing testable code, run `go test ./...` and `npm test` from `frontend/` when frontend files changed. Run `go test -race ./internal/speedtest_util ./internal/gui_wails` when touching runner, cancellation, or Wails event code.
8. Go config constants are the source of truth for shared frontend/backend values. After changing shared constants or phases, run `go generate ./...` and commit the generated frontend config.
