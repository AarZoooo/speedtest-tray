# TO-DO for the project

## Blocking bugs

None

## Non-Blocking bugs

None

## Features

Here's the list of small features to add to the application:
- [x] Anchor utility window directly to the Windows system tray icon (improving window positioning accuracy over cursor location).
- [x] Implement environment-aware logging/config redirection in dev mode and add startup log truncation to 5,000 lines.
- [x] Add an option to the system tray context menu to open the log directory in an OS-agnostic way.
- [x] Add native macOS status bar menu bar status item support (accessory app hidden from the macOS Dock).
- [x] Implement macOS application activation and window key focus on show.
- [x] Implement timing-based double-toggle click prevention for tray/status bar clicks across Windows and macOS.
- [x] Port the Open Logs Directory system tray option to the macOS status item menu.
- [x] Package macOS build into a custom DMG with an Applications shortcut and custom volume icon.
- [x] Optimize GitHub Actions workflow triggers with commit message filtering to save build minutes.
- [x] Add GitHub Actions CI test automation for both backend (Go) and frontend (Vitest) on every push.
- [x] Expand GitHub Action workflow to generate builds for various system architectures like x64 and arm for windows, intel and arm for macOS etc.
- [x] Remove the click-elsewhere-to-close-ui feature (rely entirely on the close button to close the UI window).

## Major changes

Here's the list of big changes to do to the application:
- [ ] Add an installer instead of portable exe to bind into windows to also support autostart through task manager.
- [ ] Add an option to enable updates
- [x] Add headless CLI mode for the app
- [ ] Automatically register the CLI command in the system PATH (or as a wrapper symlink) on installation/setup for easy global access (e.g., `speedtest-tray`). Note: This should only be implemented after the application has moved out of a portable-only build (i.e., once the system installer is added).
- [x] Add a build for OSes other than Windows to try and test a MacOS build
- [x] Add support for successful speedtest history being saved and available to view in UI
- [ ] Revamp UI with a more sharp, minimal monochrome visual. Less animations, visual effects, more utility.
- [ ] Replace `speedtest-go` with a custom speedtest implementation (reduces third-party dependency and gives full control over rate calculation).
- [ ] Optimize per-run struct allocation (e.g. factory-based orchestrator injection) once a custom engine exists.
