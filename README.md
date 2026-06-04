# SpeedTest Tray

A small Windows tray app for running an on-demand internet speed test from a compact Wails window.

## What It Does

- Runs from the system tray.
- Opens a fixed-size floating window near the cursor.
- Measures server, ping, download speed, and upload speed.
- Streams live progress updates to the UI during the test.
- Hides the window on close or focus loss while keeping the tray app alive.

## Tech Stack

- Go
- Wails v2
- Plain HTML, CSS, and JavaScript
- `github.com/showwin/speedtest-go` for speed tests
- `github.com/energye/systray` for the tray integration

## Project Layout

```text
.
├── main.go                    # Wails and tray entry point
├── frontend/                  # Static UI assets and Wails bindings
├── internal/gui_wails/        # Wails backend bindings and window integration
├── internal/speedtest_util/   # Speed test lifecycle and phase updates
├── assets/                    # Source app assets
└── build/windows/             # Windows icon and Wails build metadata
```

## Development

Run checks:

```powershell
go test ./...
```

Run the app in development mode:

```powershell
wails dev
```

Build the app:

```powershell
wails build
```

## Notes

The app is intentionally Windows-focused. Native Windows APIs are used for tray-window positioning and rounded corners.
