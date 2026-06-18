
<p align="center">
  <img src="frontend/assets/banner.png" alt="SpeedTest Tray Banner" width="360">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26.4-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/Wails-v2.12.0-blue" alt="Wails Version">
  <img src="https://img.shields.io/badge/platform-windows%20%7C%20macOS-lightgrey" alt="Platform Support">
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
</p>

---

**SpeedTest Tray** is a lightweight, high-performance system tray application designed for running on-demand internet speed tests from a compact, modern window. Built with Go and Wails, it combines native performance with a vibrant, hardware-accelerated UI.

<p align="center">
  <img src="assets/preview.gif" alt="SpeedTest Tray Demo" width="400">
</p>

## 📑 Table of Contents

- [Key Features](#-key-features)
- [Installation](#-installation)
- [Usage](#-usage)
- [Tech Stack](#-tech-stack)
- [Project Layout](#-project-layout)
- [Development](#-development)
- [License](#-license)

## ✨ Key Features

- **System Tray & Menu Bar Integration**: Stays docked in your taskbar (Windows) or status bar (macOS), launching a focused window only when needed.
- **Modern Speedometer UI**: Features a custom-built, modular solid-sector gauge with a real-time synchronized kite needle.
- **Vibrant Aesthetic**: Premium dual-accent gradient theme with color-matched bloom effects and softened shadows.
- **Dynamic Scaling**: Automatically adjusts its scale for Download (1000 Mbps) and Upload (100 Mbps) phases for optimal visual feedback.
- **Reliable Retests**: Each run uses a fresh test engine, so the speedometer starts at 0 Mbps on every retry.
- **Clear Step Labels**: Status text updates to the upcoming phase during pauses between ping, download, and upload.
- **Offline Detection**: Shows **No internet connection** and offers **Try Again** when a test is started without connectivity.
- **Immediate Termination**: Dedicated "Stop" button for instant test cancellation and UI reset.
- **Persistent Logging**: Configurable file-based logging stored in your system's application data folder.
- **Speedtest History**: Persists and displays recent speedtest run records (Download, Upload, Ping, Server, and Timestamp) directly in the UI. Features a 2-click clear confirmation flow and a button to open the raw JSON file natively.
- **Headless CLI Mode**: Run speed tests directly from command line interfaces with interactive progress bars (`-c`/`--cli`) or parseable structured JSON (`-j`/`--json`).

## 🚀 Installation

### Windows (Recommended)
1. Download the latest `speedtest-tray-portable.exe` from the [**Releases**](https://github.com/AarZoooo/speedtest-tray/releases) page.        
2. Run the executable. It will automatically initialize in your system tray.

### macOS
1. Download the latest `SpeedTest Tray.dmg` from the [**Releases**](https://github.com/AarZoooo/speedtest-tray/releases) page.
2. Open the `.dmg` file, and drag `SpeedTest Tray` to your `Applications` folder.
3. Open it from your Applications folder. It will initialize as a native menu bar accessory app.

### From Source
If you prefer to build it yourself, ensure you have [Go](https://go.dev/) and [Wails](https://wails.io/) installed.

```powershell
# Clone the repository
git clone https://github.com/AarZoooo/speedtest-tray.git
cd speedtest-tray

# Build the application
wails build
```

## 💡 Usage

1. **Launch**: Open the app; it will appear as a speedometer icon in your system tray or macOS menu bar.
2. **Start Test**: Click the tray icon (or right-click → Open / Show) to bring up the window, then hit **Start**.
3. **Monitor**: Watch real-time progress as the needle synchronizes with your network throughput.
4. **Stop**: Click **Stop** at any time to abort the test.
5. **Logs**: View your test history in your `%APPDATA%/SpeedTest Tray` folder (on Windows) or `~/Library/Application Support/SpeedTest Tray` folder (on macOS).

### 🖥️ Headless CLI Mode

Run speed tests directly in the terminal without spawning the GUI:

```bash
# Run test with interactive terminal progress bar
speedtest-tray-portable.exe -c

# Output results as structured JSON (ideal for scripts/cron jobs)
speedtest-tray-portable.exe -j

# Select a specific speedtest server by ID
speedtest-tray-portable.exe -s <server_id>
```

> [!NOTE]
> **Windows Terminal Execution**: Because the compiled portable app is a GUI subsystem binary, Windows shells (PowerShell/CMD) launch it in the background asynchronously by default. To output to the terminal synchronously, pipe it to a host stream:
> * **PowerShell**: `.\speedtest-tray-portable.exe -c | Out-Host`
> * **CMD**: `speedtest-tray-portable.exe -c | cat`

## 🛠 Tech Stack

- **Backend**: [Go](https://go.dev/)
- **Frontend**: [Wails v2](https://wails.io/) (HTML/CSS/JS)
- **Speed Test Engine**: `speedtest-go`
- **System Tray**: `energye/systray` (Windows) / Native AppKit Objective-C status item via CGO (macOS)
- **UI Components**: Vanilla Web Components

## 📁 Project Layout

```text
.
├── main.go                    # Wails and tray entry point
├── internal/cli/              # Headless CLI loop and progress rendering
├── internal/config/           # Centralized configuration and constants
├── internal/gui_wails/        # Wails backend bindings and window integration
├── internal/speedtest_util/   # Speed test core logic and orchestration
├── frontend/                  # Modularized UI assets and Web Components
├── docs/                      # Architecture and design documentation
├── build/windows/             # Windows icon and Wails build metadata
└── build/darwin/              # macOS icon, Info.plist, and plist templates
```

## 📖 Documentation

For more detailed information about the project's internals and guidelines, please refer to the following:

- [**Architecture Overview**](docs/ARCHITECTURE.md): Deep dive into the system design, module structure, and data flow.
- [**Development Roadmap**](docs/DEV.md): Current status, known bugs, and planned features.
- [**Engineering Rules**](docs/RULES.md): Coding standards, testing requirements, and development workflows.

## 🛠 Development

Run in development mode with hot-reload:
```powershell
wails dev
```

Run all tests (Go and Frontend):
```powershell
# Backend
go generate ./...
go test ./...

# Frontend
npm test --prefix frontend
```

For more details on our development workflow and testing standards, see [**docs/RULES.md**](docs/RULES.md).

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
