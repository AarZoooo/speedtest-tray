# Changelog

All notable changes to this project will be documented in this file.

## [1.1.2] - 2026-06-19

### Added
- Added installer options for launch at login, CLI alias installation, and desktop shortcut creation (pre-checked by default).
- Added a finish page to the installer with an option to launch the application on close.
- Added custom sidebar bitmap image support to the installer screens.
- Added launch-minimized and launch-at-login settings to the config file and system tray context menu.
- Added a check-failure state view in the update panel to display error details and a retry option when update checks fail.

### Changed
- Configured the application to read and apply launch preferences from the configuration file at startup.

### Fixed
- Fixed the Windows installer branding to display "SpeedTest Tray" instead of "Name".
- Fixed macOS PKG packages failing to open with a "damaged or incomplete" error by defining the correct installation path.
- Fixed installer welcome and finish pages rendering blank on the left side by correcting the bitmap format.
- Fixed uninstaller to prompt the user to close the app if it is currently running.
- Fixed update view turning blank on manual update check error by correcting a nested HTML container tag.
- Fixed running Start button styling to correctly display the idle secondary style (card background, no border).

## [1.1.1] - 2026-06-19

### Fixed
- **Installer overhaul**: Migrated Windows installer to Modern UI 2 with customized branding and app icon integration.
- **Corner radius gap**: Set HTML container corner radius to `8px` (and manual Win10 rounding fallback to `8px`) to match Windows 11 and macOS native window container standards, eliminating visible gaps.
- **Test isolation**: Prevented development/placeholder history entries from polluting the user's real speedtest history during testing or build phases.

## [1.1.0] - 2026-06-19

### Added
- **macOS support**: Native status bar menu bar integration with accessory app hidden from Dock
- **Windows installer**: NSIS-based installer with autostart support through Task Manager
- **CLI mode**: Headless CLI mode with interactive progress, JSON output, and history support
- **Updater system**: Auto-update functionality with platform-specific implementations (Windows/macOS)
- **History persistence**: Speedtest results saved to disk and available to view in UI
- **Multi-architecture builds**: GitHub Actions CI/CD for Windows x64/arm and macOS Intel/ARM
- **CI/CD automation**: Test automation for both backend (Go) and frontend (Vitest) on every push
- **Startup sweep animation**: Speedometer needle animation (0→max→0, ~1.5s) before backend test begins
- **Inter Variable font**: Embedded locally for crisp, non-pixelated text rendering in WebView
- **CSS design tokens**: Centralized font-size values as CSS custom properties (`--fs-*`) in `:root`
- **History pill labels**: Permanent pill labels on history action buttons with instant text swap and smooth width transitions
- **Native OS translucency**: Windows Acrylic backdrop and native OS translucency support
- **Window anchoring**: Utility window anchored directly to Windows system tray icon for improved positioning accuracy
- **Open logs directory**: System tray/context menu option to open log directory in OS-agnostic way
- **macOS DMG packaging**: Custom DMG with Applications shortcut and custom volume icon
- **macOS app activation**: Application activation and window key focus on show
- **Double-toggle prevention**: Timing-based double-toggle click prevention for tray/status bar clicks across Windows and macOS
- **Environment-aware logging**: Dev mode config redirection and startup log truncation to 5,000 lines
- **GitHub Actions optimization**: Commit message filtering to save build minutes

### Changed
- **UI revamp**: Comprehensive UI redesign with better layout, spacing, and visual polish
- **Close behavior**: Removed click-elsewhere-to-close feature; rely entirely on close button
- **Gradient theme**: Updated default accent gradient to bottom-right for true diagonal flow
- **Dark theme**: Updated dark theme colors and removed danger-gradient
- **Button styling**: Running button hover border drawn internally to prevent overflow and layout shifting
- **Corner radius**: Matched CSS corner radius to Windows and macOS native window defaults

### Fixed
- **Window rounding**: Set fallback corner radius and cleanup window rounding comment
- **Running button hover**: Fixed hover border radius clipping using background origin/clip
- **Layout spacing**: Unified layout spacing and button bottom positioning globally
- **Update view**: Vertically centered checking loader and text in update view

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
