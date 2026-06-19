# TO-DO for the project

## Blocking bugs

None

## Non-Blocking bugs

- [ ] Installer looks hideous and feels like from early 2000 Windows XP times
- [ ] Installer exe doesn't have the app icon, instead shows the old generic .exe icon
- [ ] History shows dummy data on fresh install (history.json is prefilled with placeholder data)

## Features

Here's the list of small features to add to the application:
None - all features completed in v1.1.0

## Major changes

Here's the list of big changes to do to the application:
- [ ] Replace `speedtest-go` with a custom speedtest implementation (reduces third-party dependency and gives full control over rate calculation).
- [ ] Optimize per-run struct allocation (e.g. factory-based orchestrator injection) once a custom engine exists.
