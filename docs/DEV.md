# TO-DO for the project

## Blocking bugs

- [x] On display scaling set to a value other than 100%, the UI gets stretched and clipped.

## Non-Blocking bugs

- [ ] Retrying speed tests after completion or failure still shows the speed coming down from around 3000 Mbps
- [ ] The 2-second sleep between each step is indistinguishable and the proper process labels aren't being shown in the sleep times.
- [x] The tests are buggy, incomplete, and don't cover everything.

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
