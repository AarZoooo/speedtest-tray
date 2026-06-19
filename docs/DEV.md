# TO-DO for the project

## Blocking bugs

None

## Non-Blocking bugs

- [ ] Fix Windows installer sidebar bitmap image rendering (image remains blank/invisible during setup)

## Features

None

## Major changes

Here's the list of big changes to do to the application:
- [ ] Replace `speedtest-go` with a custom speedtest implementation (reduces third-party dependency and gives full control over rate calculation).
- [ ] Optimize per-run struct allocation (e.g. factory-based orchestrator injection) once a custom engine exists.
- [ ] Audit and optimize memory/CPU allocations, telemetry updates, and UI rendering performance to minimize hardware footprint
- [ ] Enable Linux builds (tray integration is complex on Linux due to desktop environment differences)
