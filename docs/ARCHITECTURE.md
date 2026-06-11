# Architecture

## Overview

Speedtest Tray is a modular speed testing application with clean separation of concerns:

- **Config Layer** (`internal/config/`): Centralized configuration and constants
- **Business Logic** (`internal/speedtest_util/`): Speed testing orchestration and progress calculation
- **GUI Layer** (`internal/gui_wails/`): Wails framework integration
- **Frontend** (`frontend/`): Modularized JavaScript with state, handlers, and constants

## Module Structure

### internal/config

**Purpose**: Single source of truth for all configuration values

- `constants.go`: All hardcoded values (progress thresholds, window dimensions, gauge scales, test durations, UI timing)
- `config.go`: Configuration file I/O (loading/saving JSON)
- `phases.go`: Phase lifecycle constants
- `cmd/gen-frontend-config`: Generates frontend shared constants from Go config with `go generate ./...`

**Usage**: Import to access centralized values. Example:
```go
import "speedtest-tray/internal/config"
fmt.Println(config.WindowWidth)  // 320
fmt.Println(config.PhaseDownloading)  // "DOWNLOADING"
```

### internal/speedtest_util

**Purpose**: Core speed testing logic independent of GUI framework

**Key Components**:

1. **TestOrchestrator Interface** (`orchestrator_interface.go`)
   - Defines contract for speed testing
   - Enables mock implementations for testing
   - Methods: GetUserInfo, FindServers, SelectBestServer, RunPing, RunDownload, RunUpload

2. **SpeedTester Implementation** (`orchestrator_impl.go`)
   - Implements TestOrchestrator
   - Wraps speedtest-go library calls
   - Handles callbacks for progress updates (receives float64 MBPS)

3. **TestRunner** (`runner.go`)
   - Orchestrates full test workflow
   - Probes internet connectivity before any speedtest API calls
   - Retries setup/initialization phases (`GetUserInfo`, `FindServers`, `SelectBestServer`, and `RunPing`) to handle transient network errors
   - Manages context, cancellation, phase sequencing
   - Handles progress mapping (phase progress → total test progress)
   - Lifecycle: Initialize → Connectivity check → Ping → Download → Upload → Complete

4. **Connectivity Check** (`connectivity.go`)
   - HEAD request to a lightweight probe endpoint with a short timeout
   - Fails fast with `ErrNoInternet` when the machine appears offline

5. **Progress Helpers** (`progress.go`)
   - `CalculatePhaseProgress()`: Calculates progress within a single phase
   - `MapPhaseProgressToTotal()`: Maps phase progress to overall test progress
   - `FormatNumber()`: Consistent float formatting (2 decimal places)

6. **Data Types** (`types.go`)
   - `Update`: Progress notification to GUI
   - `Result`: Final test results

### internal/gui_wails

**Purpose**: Wails framework integration layer

**Key Components**:

1. **App Struct** (`app.go`)
   - Minimal Wails binding layer for window, tray, and test control
   - Holds only Wails context and the active test cancellation handle
   - On each `StartTest`, creates a fresh `SpeedTester` and `TestAdapter` for that run
   - Methods: Startup, StartTest, StopTest

2. **TestAdapter** (`adapter.go`)
   - Bridges Wails events with business logic
   - Created per test run
   - Converts TestRunner callbacks to Wails runtime events
   - Serializes results for frontend consumption
   - Zero business logic dependency

3. **Window Management** (`window_windows.go`)
   - Platform-specific window setup (Windows)
   - Uses config constants for sizing, corner radius

### frontend

**Purpose**: Modularized vanilla JavaScript UI

**Structure**:

- `index.html`: HTML structure (removed inline onclick handlers)
- `main.js`: Module orchestrator (imports all modules, initializes app)
- `src/constants.js`: Frontend event names and re-exports for generated shared constants
- `src/generated/config.js`: Generated phase and UI config constants from `internal/config`
- `src/state.js`: TestState class (centralized test state management)
- `src/ui.js`: UI update handlers (results, gauge, status, button state)
- `src/handlers.js`: Test control (start, stop, button click)
- `src/window.js`: Window events (show, blur, visibility)
- `speedometer.js`: Custom gauge component (Web Component)
- `src/speedometer-config.js`: Gauge configuration constants
- `style.css`: Global styles
- `speedometer.css`: Gauge styles

**Data Flow**:

```
User clicks "Start"
  ↓
handlers.js: startTest()
  ↓
Calls window.go.gui_wails.App.StartTest()
  ↓
Go: App.StartTest() → fresh SpeedTester + TestAdapter → TestAdapter.RunTest()
  ↓
TestRunner executes phases (new instance per run)
  ↓
Callbacks emit progress → window.runtime.EventsEmit()
  ↓
Frontend receives test_update event
  ↓
ui.js: handleTestUpdate() updates gauge/results
```

### 7. Logging and Telemetry

**Purpose**: Structured logging and process performance tracking

**Key Components**:

1. **Structured Logging** (`log/slog`)
   - Uses `slog.JSONHandler` for persistent file logging (`app.log`)
   - Uses `slog.TextHandler` for standard output
   - Global logger set via `slog.SetDefault()` in `main.go`
   - Log messages and keys centralized in `internal/config/constants.go`

2. **Process Telemetry** (`internal/speedtest_util/telemetry.go`)
   - Captures process-specific metrics: `AllocMB`, `SysMB`, `NumGoroutine`
   - Integrated into `TestRunner.RunTest()` to log hardware utilization at test start and completion

## Design Patterns

### 1. Dependency Injection

`TestRunner` and `TestAdapter` accept a `TestOrchestrator` interface, enabling:
- Production use: a fresh `SpeedTester` created in `App.StartTest()`
- Testing: mock orchestrator for unit tests

`App` constructs the production orchestrator on demand rather than receiving a long-lived instance from `main`.

### 2. Callback-Based Progress

TestRunner uses callbacks instead of shared state:
- Download callback: `func(float64)` (Mbps)
- Upload callback: `func(float64)` (Mbps)

This keeps phases decoupled and testable.

### 3. Event-Driven Frontend

Wails runtime events push updates to frontend:
- `test_update`: Periodic progress notifications
- `test_complete`: Final results
- `test_error`: Error handling

Frontend is reactive—no polling, no tight coupling.

### 4. Phase Orchestration

Progress tracked as phases progress:
- Each phase has start/end thresholds (e.g., download: 0.30-0.70)
- Phase progress (0-1) mapped to total progress (start + (end-start)*phaseProgress)
- Frontend displays percentage complete (total progress * 100)
- Between major steps, `TestRunner` sleeps for `PhaseSleepDuration` (2 seconds). The next-phase update (e.g. `STARTING_DOWNLOAD`, `STARTING_UPLOAD`) is emitted **before** the sleep so the UI status label and gauge reset reflect the upcoming step during the pause, not the phase that just finished

## CI/CD and Build Pipeline

The project uses GitHub Actions (defined in `.github/workflows/build.yml`) to build and package releases for Windows and macOS.

### Job Trigger Restrictions
To save GitHub Actions runner minutes (which are limited to 2,000 free minutes per month for private repositories), the build jobs only run under two conditions:
1. The commit message contains the string `[release commit]`.
2. The workflow is manually dispatched from the GitHub Actions UI.

All standard commits will automatically skip the compilation and packaging jobs.

### Platform Builds
- **Windows Build:** Compiles the application using Wails for the `windows/amd64` platform and outputs a portable `.exe` executable.
- **macOS Build:** 
  - Compiles the application as a universal binary (`darwin/universal`).
  - Packages the resulting `SpeedTest Tray.app` bundle into a `.dmg` (Disk Image) file using macOS `hdiutil`.
  - The DMG mounts with a custom volume icon (`.VolumeIcon.icns` generated from the app icon) and includes a symlink to `/Applications` for a native drag-to-install experience.

## Configuration

All configurable values in `internal/config/constants.go`:

```go
// Progress thresholds (0.0-1.0)
const ProgressDownStart = 0.30  // Download starts at 30% total
const ProgressDownEnd = 0.70    // Download ends at 70% total

// UI timing
const UIHideDelayMs = 2000      // Hide window after 2 seconds of showing

// Window properties
const WindowWidth = 320
const WindowHeight = 560

// Gauge scales (Mbps)
const GaugeMaxDownload = 1000
const GaugeMaxUpload = 100
```

Edit here to tune behavior without touching business logic. If a value is shared with the frontend, run:

```bash
go generate ./...
```

This regenerates `frontend/src/generated/config.js`, keeping Go as the source of truth.

## Testing

The project uses deterministic tests instead of live network tests or a running Wails window.

**Backend Tests**:

- `internal/config/config_test.go`: Validates constants, phase strings, and config file persistence with temp directories.
- `internal/speedtest_util/progress_test.go`: Tests progress calculation, clamping, and formatting.
- `internal/speedtest_util/mock_orchestrator_test.go`: Mock implementation of `TestOrchestrator`.
- `internal/speedtest_util/connectivity_test.go`: Tests connectivity probing and offline error detection.
- `internal/speedtest_util/runner_test.go`: Tests phase sequencing, progress mapping, cancellation, offline failure, failures, channel closure, and final results.
- `internal/gui_wails/adapter_test.go`: Tests serialization and Wails event routing through an injected emitter.

**Frontend Tests**:

- Vitest with JSDOM runs from `frontend/`.
- `frontend/src/constants.test.js`: Verifies generated frontend constants stay aligned with Go.
- `frontend/src/state.test.js`: Tests isolated `TestState` behavior.
- `frontend/src/ui.test.js`: Tests DOM updates, gauge calls, completion, stop, and error states.
- `frontend/src/handlers.test.js`: Tests start/stop button behavior with mocked Wails bindings.
- `frontend/src/window.test.js`: Tests window event behavior with mocked runtime APIs.

**Running Tests**:

```bash
go generate ./...
go test ./...
npm test --prefix frontend
go test -race ./internal/speedtest_util ./internal/gui_wails
```

**Mock Orchestrator Usage**:

```go
mock := &speedtest_util.MockOrchestrator{
    PingResult: 20 * time.Millisecond,
    DownloadResult: 90.5,
    UploadResult: 18.2,
}
runner := speedtest_util.NewTestRunner(mock)
```

## Key Decisions

### 1. Central Config Package

**Why**: 28+ hardcoded values scattered across files. Central config reduces:
- Maintenance burden (one edit point)
- Sync bugs (phase names, thresholds)
- Duplication

### 2. TestOrchestrator Interface

**Why**: SpeedTester had mixed responsibilities. Interface abstracts:
- Library calls (implemented in SpeedTester)
- Lifecycle/sequencing (implemented in TestRunner)
- Enables testing without real speedtest library

### 3. Wails/Business Logic Separation

**Why**: Originally business logic tangled in Wails App struct. Separation enables:
- Reusable TestRunner (swap orchestrator for testing)
- Cleaner App struct (minimal binding layer)
- Unit testable business logic

### 4. Modularized Frontend

**Why**: 200+ lines in single main.js with global state. Modules provide:
- Maintainability (split by concern)
- Reusability (export handlers, state)
- Clear data flow (state → UI updates)

### 5. Generated Shared Constants

**Why**: Phase and UI config values are needed in both Go and browser JavaScript. Go remains the source of truth, while `go generate ./...` writes `frontend/src/generated/config.js`.

This avoids drift between backend phases, frontend state handling, and UI timing/gauge values.

### 6. Deterministic Test Boundaries

**Why**: The real speedtest library, Wails runtime, browser DOM, and user config directory are external systems. Tests use mocks, temp directories, injected emitters, and JSDOM so they remain fast and reliable.

### 7. Per-Run Test Engine Lifecycle

**Why**: Reusing a single `speedtest-go` client across runs leaked EWMA rate state into the first download/upload callbacks of the next test, causing the speedometer to briefly show thousands of Mbps. Creating a fresh `SpeedTester` and `TestAdapter` per run guarantees each test starts from zero without manual reset logic against library internals.

**Trade-off**: A few extra allocations per test, which is negligible compared to network I/O and test buffers. Struct reuse or a custom speedtest engine can be revisited later if needed.
