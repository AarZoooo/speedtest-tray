# Changelog

All notable changes to this project will be documented in this file.

## [1.0.2] - 2026-06-11

### Added
- Offline detection before a test starts, with a dedicated **No internet connection** status when no connectivity is available.
- Migrated to structured logging using `log/slog` with JSON output for file logs and Text output for console.
- Hardware utilization telemetry (Memory and Goroutines) logged at the start and end of every test run.
- Retry loops for setup/initialization failures (fetching user info, locating servers, server selection, ping latency) to handle transient network issues.

### Changed
- `SpeedTester` and `TestAdapter` are now created fresh for each test run instead of being reused for the application lifetime.
- Explicitly configure the underlying speedtest client's test duration to match config constants instead of relying on library defaults.
- Changed high-frequency progress event logging from `INFO` to `DEBUG` level.

### Fixed
- Speedometer needle and Mbps display no longer spike to extremely high values when retrying a test after completion, failure, or stop.
- Status label and gauge now reflect the upcoming phase during 2-second pauses between ping → download and download → upload (e.g. "Starting download test..." and "Starting upload test..." instead of holding the previous phase).
- Resolved progress bar getting stuck early during speed tests by updating estimated phase durations to 15 seconds to match the underlying library's default runtime, and implemented asymptotic progress scaling once estimated duration is exceeded.

## [1.0.1] - 2026-06-07

### Added
- "Stopped" phase for better feedback when a test is cancelled.
- Automatic generation of frontend constants from Go config.

### Changed
- Major test overhaul with comprehensive unit tests for both backend and frontend.
- Improved DPI detection and handling for high-resolution displays on Windows.

### Fixed
- UI stretching and clipping when display scaling is set above 100%.
- Panic and UI hang when stopping an ongoing download or upload test.
- Corrected various small UI state bugs during test lifecycle.

## [1.0.0] - 2026-06-05

### Added
- Initial release of SpeedTest Tray.
- System tray integration for quick access.
- Real-time speedometer gauge with gradient theme.
- Support for session logging to app data.
- Modular architecture with comprehensive test coverage.
- Support for Windows with frameless, transparent UI and rounded corners.
