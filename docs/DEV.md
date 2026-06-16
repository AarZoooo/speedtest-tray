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

## Major changes

Here's the list of big changes to do to the application:
- [ ] Add an installer instead of portable exe to bind into windows to also support autostart through task manager.
- [ ] Add an option to enable updates
- [ ] Add headless CLI mode for the app
- [x] Add a build for OSes other than Windows to try and test a MacOS build
- [ ] Add support for successful speedtest history being saved and available to view in UI
- [ ] Revamp UI with a more sharp, minimal monochrome visual. Less animations, visual effects, more utility.
- [ ] Replace `speedtest-go` with a custom speedtest implementation (reduces third-party dependency and gives full control over rate calculation).
- [ ] Optimize per-run struct allocation (e.g. factory-based orchestrator injection) once a custom engine exists.
