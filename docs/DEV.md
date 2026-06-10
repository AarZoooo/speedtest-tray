# TO-DO for the project

## Blocking bugs

None

## Non-Blocking bugs

- [x] Retrying speed tests after completion or failure still shows the speed coming down from around 3000 Mbps
- [ ] The 2-second sleep between each step is indistinguishable and the proper process labels aren't being shown in the sleep times.

## Features

Here's the list of small features to add to the application:
- [ ] Add offline indicator support for when the test is run without an internet connection
- [ ] Add retry loops for failures
- [ ] Add better logger support with hardware utilization stats (use slog)

## Major changes

Here's the list of big changes to do to the application:
- [ ] Add an installer instead of portable exe to bind into windows to also support autostart through task manager.
- [ ] Add an option to enable updates
- [ ] Add headless CLI mode for the app
- [ ] Add a build for OSes other than Windows to try and test a MacOS build
- [ ] Add support for successful speedtest history being saved and available to view in UI
- [ ] Revamp UI with a more sharp, minimal monochrome visual. Less animations, visual effects, more utility.
- [ ] Replace `speedtest-go` with a custom speedtest implementation (reduces third-party dependency and gives full control over rate calculation).
- [ ] Optimize per-run struct allocation (e.g. factory-based orchestrator injection) once a custom engine exists.
