# SpeedTest Tray

A small, high-performance system tray application for running on-demand internet speed tests from a compact, modern window.

## Key Features

- **System Tray Integration**: Stays out of your way until you need it.
- **Modern Speedometer UI**: Features a custom-built, modular solid-sector gauge with a real-time synchronized kite needle.
- **Vibrant Aesthetic**: Uses a premium dual-accent gradient theme with high-quality, color-matched bloom (glow) effects and softened shadows.
- **Dynamic Scaling**: The speedometer automatically adjusts its scale for Download (1000 Mbps) and Upload (100 Mbps) phases to provide meaningful visual feedback.
- **Test Cancellation**: Immediate termination support with a dedicated "Stop" button that instantly resets the UI.
- **Persistent Logging**: Configurable file-based logging with support for standard and OneDrive-synced Documents folders.
- **Cross-Platform Readiness**: Designed with agnostic pathing and logic for both Windows and macOS support.

## Tech Stack

- **Backend**: Go
- **Frontend**: Wails v2 (HTML/CSS/JS)
- **Speed Test Engine**: `github.com/showwin/speedtest-go`
- **System Tray**: `github.com/energye/systray`
- **UI Components**: Vanilla Web Components

## Project Layout

```text
.
├── main.go                    # Wails and tray entry point
├── internal/config/           # Centralized configuration and constants
├── internal/gui_wails/        # Wails backend bindings and window integration
├── internal/speedtest_util/   # Speed test core logic and orchestration
├── frontend/                  # Modularized UI assets and Web Components
├── docs/                      # Architecture and design documentation
├── assets/                    # Source app assets
└── build/windows/             # Windows icon and Wails build metadata
```

## Development

Run the app in development mode:

```powershell
wails dev
```

Build the app:

```powershell
wails build
```

## Notes

While the app features a platform-agnostic frontend, it utilizes native Win32 APIs on Windows for precise window positioning near the taskbar and hardware-accelerated rounded corners.
