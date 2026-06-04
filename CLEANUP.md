# Cleanup Progress

## Phase 1: Config System ✅ DONE
- Created `internal/config` package with constants.go and config.go
- Centralized all hardcoded values (timeouts, progress thresholds, window properties, gauge scales)
- Removed config logic from main.go
- Updated all imports to use config package
- Build verified: ✓
- **Commit**: Phase 1: Centralize config and hardcoded constants

## Phase 2: Eliminate Redundancy ✅ DONE
- Unified PHASES constants: defined once in config/phases.go, imported everywhere
- Extracted progress calculation helpers: CalculatePhaseProgress() and MapPhaseProgressToTotal()
- Converted inline helpers to methods: checkContextCancelled() and sleepWithInterrupt()
- Removed duplicate logic from download/upload test callbacks
- Build verified: ✓
- **Commit**: Phase 2: Unify phase constants and extract progress helpers

## Phase 3: Modularize SpeedTester ✅ DONE
- Created TestOrchestrator interface: defines contract for test execution
- Implemented interface in SpeedTester: GetUserInfo, FindServers, SelectBestServer, RunPing, RunDownload, RunUpload
- Simplified RunTest(): cleaner orchestration, easier to understand phase flow
- Removed old phase methods: phases_init.go and phases_tests.go deleted
- All test logic now in orchestrator methods (orchestrator_impl.go)
- Build verified: ✓
- **Commit**: Phase 3: Create TestOrchestrator interface and simplify RunTest

## Phase 4: Separate Wails/Business Logic [PENDING]
- [ ] Create adapter layer
- [ ] Move speedtest logic to separate package
- [ ] Simplify App struct

## Phase 5: Frontend Modularization [PENDING]
- [ ] Extract state management
- [ ] Separate event handlers
- [ ] Extract constants to modules
- [ ] Move speedometer config

## Phase 6: Testing [PENDING]
- [ ] Add Go unit tests
- [ ] Add mock TestOrchestrator
- [ ] Frontend test setup

## Phase 7: Cleanup & Docs [PENDING]
- [ ] Remove unnecessary comments
- [ ] Add ARCHITECTURE.md
- [ ] Update README

---

Files Modified in Phase 1:
- internal/config/constants.go (NEW)
- internal/config/config.go (NEW)
- main.go
- internal/gui_wails/app.go
- internal/gui_wails/window_windows.go
- internal/speedtest_util/types.go
- internal/speedtest_util/speedtester.go
- internal/speedtest_util/phases_init.go
- internal/speedtest_util/phases_tests.go
- frontend/main.js
- config.yaml.template (NEW)
