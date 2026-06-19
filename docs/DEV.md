# TO-DO for the project

## Blocking bugs

None

## Non-Blocking bugs

- [ ] Show error message in UI if update check fails (currently silently falls back to up-to-date)

## Features

Here's the list of small features to add to the application:
- [ ] Add context menu option (checkbox) to configure whether the application starts minimized to tray or shows the UI on startup
  - [ ] Save this launch-minimized preference in `config.json`
- [ ] Add the launch at login (autostart) preference to the `config.json` settings file
- [ ] Read and apply `config.json` settings at application startup (e.g., launch minimized, launch at login)

## Major changes

Here's the list of big changes to do to the application:
- [ ] Replace `speedtest-go` with a custom speedtest implementation (reduces third-party dependency and gives full control over rate calculation).
- [ ] Optimize per-run struct allocation (e.g. factory-based orchestrator injection) once a custom engine exists.
- [ ] Audit and optimize memory/CPU allocations, telemetry updates, and UI rendering performance to minimize hardware footprint
