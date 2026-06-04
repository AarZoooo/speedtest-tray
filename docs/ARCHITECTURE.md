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
- `config.go`: Configuration file I/O (loading/saving YAML)
- `phases.go`: Phase lifecycle constants (used by both Go and JavaScript)

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
   - Manages context, cancellation, phase sequencing
   - Handles progress mapping (phase progress → total test progress)
   - Lifecycle: Initialize → Ping → Download → Upload → Complete

4. **Progress Helpers** (`progress.go`)
   - `CalculatePhaseProgress()`: Calculates progress within a single phase
   - `MapPhaseProgressToTotal()`: Maps phase progress to overall test progress
   - `FormatNumber()`: Consistent float formatting (2 decimal places)

5. **Data Types** (`types.go`)
   - `Update`: Progress notification to GUI
   - `Result`: Final test results

### internal/gui_wails

**Purpose**: Wails framework integration layer

**Key Components**:

1. **App Struct** (`app.go`)
   - Minimal Wails binding layer
   - Delegates to TestAdapter
   - Methods: Startup, StartTest, StopTest

2. **TestAdapter** (`adapter.go`)
   - Bridges Wails events with business logic
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
- `src/constants.js`: Phase constants, event names, configuration (must match Go values)
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
Go: App.StartTest() → TestAdapter.Start()
  ↓
TestRunner executes phases
  ↓
Callbacks emit progress → window.runtime.EventsEmit()
  ↓
Frontend receives test_update event
  ↓
ui.js: handleTestUpdate() updates gauge/results
```

## Design Patterns

### 1. Dependency Injection

TestRunner accepts TestOrchestrator interface, enabling:
- Production use: real SpeedTester implementation
- Testing: mock orchestrator for unit tests

### 2. Callback-Based Progress

TestRunner uses callbacks instead of shared state:
- Ping callback: `func(float64)` (milliseconds)
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

## Configuration

All configurable values in `internal/config/constants.go`:

```go
// Progress thresholds (0.0-1.0)
const ProgressDownStart = 0.30  // Download starts at 30% total
const ProgressDownEnd = 0.70    // Download ends at 70% total

// UI timing
const UIHideDelayMs = 1000      // Hide window after 1 second of showing

// Window properties
const WindowWidth = 320
const WindowHeight = 560

// Gauge scales (Mbps)
const GaugeMaxDownload = 1000
const GaugeMaxUpload = 100
```

Edit here to tune behavior without touching business logic.

## Testing

**Unit Tests**:

- `internal/config/config_test.go`: Validates constants
- `internal/speedtest_util/progress_test.go`: Tests progress calculations
- `internal/speedtest_util/mock_orchestrator_test.go`: Mock implementation for testing

**Running Tests**:

```bash
go test ./internal/...
```

**Mock Orchestrator Usage**:

```go
mock := &speedtest_util.MockOrchestrator{
    UserInfoResult: &UserInfo{IP: "1.2.3.4"},
    PingResult: 10.5,
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

## Future Improvements

Listed in ascending order of implementation difficulty:

1. **Frontend Bundler**: Currently vanilla JS modules in browser. Migrating to Webpack or Vite would enable tree-shaking, minification, and more efficient CSS bundling.
2. **CLI Interface**: Add a `--cli` flag to allow running headless speed tests with JSON output, making the tool useful for automation and scripting.
3. **Error Recovery**: Improve resilience by implementing retry logic for transient network failures and supporting partial results if a specific phase (e.g., upload) fails.
4. **Enhanced Logging**: Transition to structured logging (using Go's `slog`) and incorporate real-time hardware usage data (CPU, RAM) to correlate system performance with test results.
5. **Memory Footprint Reduction**: Conduct deep profiling of the Wails/Webview lifecycle to minimize memory usage, ensuring the tray app remains as lightweight as possible during idle periods.
6. **Historical Results**: Implement a local persistence layer (SQLite or JSON) to save test history and create a new UI view to visualize performance trends over time.
7. **MacOS Native Builds**: Extend platform support to macOS, requiring native window positioning logic, tray handling adjustments, and a full notarization/packaging workflow.
