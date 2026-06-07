# Changelog

All notable changes to this project will be documented in this file.

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
