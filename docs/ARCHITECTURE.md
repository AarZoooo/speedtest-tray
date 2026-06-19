# Architecture

## Overview

Speedtest Tray is a modular speed testing application with clean separation of concerns:

- **Config Layer** (`internal/config/`): Centralized configuration and constants
- **Business Logic** (`internal/speedtest_util/`): Speed testing orchestration, progress calculation, and data persistence (history management)
- **Autostart Layer** (`internal/autostart/`): Platform-specific launch-at-login settings (registry manipulation for Windows, LaunchAgent plists for macOS)
- **Updater Layer** (`internal/updater/`): In-app self-update checks, downloads, progress tracking, and binary swapping/installer execution
- **GUI Layer** (`internal/gui_wails/`): Wails framework integration
- **CLI Layer** (`internal/cli/`): Headless CLI framework integration and progress output
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

7. **History Management** (`history.go`)
   - `HistoryEntry`: Data struct representing a speedtest result
   - `SaveToHistory()`: Prepends successful speedtest results up to a maximum cap (50 entries)
   - `LoadHistory()`: Reads and unmarshals JSON records, filtering out invalid/corrupted ones
   - `ClearHistory()`: Resets the history file by writing an empty array `[]`

### internal/autostart

**Purpose**: Platform-specific launch-at-login configuration and autostart manager

- `autostart.go`: Defines the core `Manager` struct referencing the current running executable.
- `autostart_windows.go`: Registry manager accessing `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` to enable/disable startup.
- `autostart_darwin.go`: Integrates with Objective-C AppKit menu items and manages user LaunchAgents plists to register the app bundle for login launch.
- `autostart_other.go`: No-op placeholder for other operating systems.

### internal/updater

**Purpose**: GitHub-based self-update engine

- `updater.go`: Performs update checks by querying the GitHub Releases API. Compares semantic versioning tags to check for updates. Holds path constants for stashing temp downloaded installers.
- `updater_windows.go`: Handles downloading the update binary, verifying checksum size, executing the new installer silently with `/S`, and terminating the active instance.
- `updater_darwin.go`: Downloads the macOS update asset, swaps the bundle, clears Gatekeeper quarantine attributes, and manages permission recovery.
- `updater_other.go`: No-op updater implementation for unsupported systems.

### internal/gui_wails

**Purpose**: Wails framework integration layer

**Key Components**:

1. **App Struct** (`app.go`)
   - Minimal Wails binding layer for window, tray, test control, autostart configuration, and updates
   - Holds Wails context, the active test cancellation handle, and cached update information
   - On each `StartTest`, creates a fresh `SpeedTester` and `TestAdapter` for that run
   - Methods: Startup, StartTest, StopTest, GetLaunchAtLogin, SetLaunchAtLogin, CheckForUpdate, ApplyUpdate, SkipUpdate

2. **TestAdapter** (`adapter.go`)
   - Bridges Wails events with business logic
   - Created per test run
   - Converts TestRunner callbacks to Wails runtime events
   - Serializes results for frontend consumption
   - Zero business logic dependency

3. **Window Management** (`window_windows.go`)
   - Platform-specific window setup (Windows)
   - Uses config constants for sizing, corner radius

### internal/cli

**Purpose**: Headless CLI framework integration layer

**Key Components**:

1. **CLI Engine** (`cli.go`)
   - Handles headless execution of speed testing using a custom terminal renderer
   - Listens to the `Update` channel emitted by `TestRunner`
   - Supports two output formats: interactive pretty-print with carriage returns (`\r`), and parseable JSON output
   - Methods: Run

### frontend

**Purpose**: Modularized vanilla JavaScript UI

**Structure**:

- `index.html`: HTML structure (removed inline onclick handlers)
- `main.js`: Module orchestrator (imports all modules, initializes app)
- `src/constants.js`: Frontend event names and re-exports for generated shared constants
- `src/generated/config.js`: Generated phase and UI config constants from `internal/config`
- `src/state.js`: TestState class (centralized test state management)
- `src/ui.js`: UI update handlers (results, gauge, status, button state, rendering history cards)
- `src/handlers.js`: App controls (start/stop speedtests, window close, toggle history view, clear history with 2-click confirm, open json natively)
- `src/window.js`: Window events initialization
- `speedometer.js`: Custom gauge Web Component — includes `playStartupSweep()` / `stopSweep()` for the pre-test needle animation
- `src/speedometer-config.js`: Gauge configuration constants
- `style.css`: Global styles — design tokens (colors, spacing, radii, shadows, font-size scale) defined as CSS custom properties in `:root`
- `speedometer.css`: Gauge styles (SVG layout, needle, fill, bloom, easing and sweep transition classes)
- `assets/fonts/InterVariable.woff2`: Embedded Inter Variable font (weight 100–900, single file) for crisp rendering in the WebView

**Data Flow**:

```
User clicks "Start"
  ↓
handlers.js: startTest()
  ↓
Speedometer plays startup sweep animation (needle 0→max→0, ~1.5s) — awaited before backend call
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
   - Automatically truncates `app.log` to the last 5,000 lines on startup to prevent unbounded file growth.
   - Environment-aware routing: redirects `app.log` and `config.json` to the project root directory during development (when running via `wails dev`, built with the `dev` tag) and ignores them in Git.

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

The project uses GitHub Actions (defined in `.github/workflows/build.yml`) to test the codebase on every push and build release packages for Windows and macOS.

### Job Structure

The workflow defines jobs that run in order:

```
test ──► build-windows (amd64)
     ├──► build-windows (arm64)
     ├──► build-macos   (amd64 / Intel)
     └──► build-macos   (arm64 / Apple Silicon)
```

#### `test` — Runs on every push
- Executes on a fast `ubuntu-latest` runner (no runner-minute cost concern).
- Runs `go generate ./...` to ensure frontend config is up to date.
- Runs the full Go test suite: `go test -v ./...`
- Runs Go race detector tests on concurrency-sensitive packages: `go test -race ./internal/speedtest_util ./internal/gui_wails`
- Installs frontend dependencies and runs `npm test` (Vitest).
- **Both build jobs declare `needs: test`**, so they are skipped entirely if any test fails.

#### `build-windows` and `build-macos` — Runs only on release commits
To save GitHub Actions runner minutes (macOS and Windows runners consume minutes faster), these jobs only run when:
1. The commit message contains the string `[ci]`.
2. The workflow is manually dispatched from the GitHub Actions UI.

Standard commits will skip compilation and packaging but always run the `test` job.

### Platform Builds

Each build job runs a **matrix strategy**, producing specific packaging formats per architecture. No fat/universal binaries are used — each artifact contains only the code that runs natively on the target CPU.

| Job             | Arch    | Runner           | Output Formats / Assets |
|-----------------|---------|------------------|-------------------------|
| `build-windows` | `amd64` | `windows-latest` | `SpeedTest-Tray-windows-amd64-installer.exe` |
| `build-windows` | `arm64` | `windows-latest` | `SpeedTest-Tray-windows-arm64-installer.exe` |
| `build-macos`   | `amd64` | `macos-latest`   | `SpeedTest-Tray-macOS-Intel.dmg`<br>`SpeedTest-Tray-macOS-Intel.pkg`<br>`speedtest-tray-darwin-amd64.tar.gz` |
| `build-macos`   | `arm64` | `macos-latest`   | `SpeedTest-Tray-macOS-ARM.dmg`<br>`SpeedTest-Tray-macOS-ARM.pkg`<br>`speedtest-tray-darwin-arm64.tar.gz` |

- **Windows builds:** Compiled using Wails and packaged into per-user NSIS installers (`SpeedTest-Tray-windows-*-installer.exe`). These register the CLI `speedtest-tray` in the user's `PATH`. The installer is also used directly by the in-app self-updater for silent updates.
- **macOS builds:** Compiled natively (for `arm64`) or cross-compiled (for `amd64`/Intel) using the macOS SDK. Each architecture produces three formats:
  * `.dmg` (Disk Image): Standard drag-and-drop consumer installer.
  * `.pkg` (Installer Package): Standard installer wizard that automates autostart plist registration and CLI symlinking (`/usr/local/bin/speedtest-tray`).
  * `.tar.gz` (Archive): Contains the raw executable app bundle, used specifically by the in-app hot-swap auto-updater.

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
