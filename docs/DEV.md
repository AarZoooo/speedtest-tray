# TO-DO for the project

## Blocking bugs

- [ ] Fix Windows system tray idle hang: App freezes/becomes unresponsive after remaining minimized in the system tray for extended periods.

## Non-Blocking bugs

- [ ] Address macOS-specific UI polish and layout inconsistencies (e.g., menu bar styling, window sizing)

## Features

- [ ] Build a dedicated Settings panel/tab in the Wails GUI.
- [ ] Add Speedometer Range configuration (e.g., select max bounds like 100, 500, 1000, 2500 Mbps).
- [ ] Move "Enable Session Logging", "Start Minimized", and "Launch at Login" toggles into the Settings GUI.
- [ ] Implement robust trace logging from application startup to close.

## Major changes

Here's the list of big changes to do to the application:
- [ ] Replace `speedtest-go` with a custom speedtest implementation (reduces third-party dependency and gives full control over rate calculation).
- [ ] Optimize per-run struct allocation (e.g. factory-based orchestrator injection) once a custom engine exists.
- [ ] Audit and optimize memory/CPU allocations, telemetry updates, and UI rendering performance to minimize hardware footprint.
- [ ] Rigorously test and debug on physical MacBook hardware using local development builds (macOS release polish).
- [ ] Enable Linux builds (tray integration is complex on Linux due to desktop environment differences).
