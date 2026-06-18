# Installer Plan

## Overview

This document is the authoritative technical spec for the installer, self-update,
Launch at Login, and CLI alias features. It is written for agent execution and is
structured to match the git branching strategy below.

**Read docs/RULES.md before making any changes.**

---

## Branch Structure

```
master
└── feature/installer          ← ROOT: core scaffolding (this branch)
    ├── feature/installer/cli-alias     ← PATH registration + /usr/local/bin symlink
    ├── feature/installer-autostart     ← Launch at Login toggle (tray + Wails bindings)
    ├── feature/installer/update-ui     ← Update panel UI (SVG icon, Wails events, panel)
    └── feature/installer/tests         ← unit tests + complete release.yml
```

Each sub-branch **branches from `feature/installer`** after the root work is
committed. Each sub-branch **merges back into `feature/installer`** when done.
`feature/installer` is then merged to `master` after all sub-branches are merged
and manual QA passes.

---

## Root Branch: `feature/installer` [DONE]

### Goal
Lay the foundation that all three sub-branches depend on:
- [x] Version and GitHub constants in config (kept at 1.0.2 per user feedback)
- [x] `internal/updater/` package (core logic, no tests yet)
- [x] `internal/autostart/` package (core logic, no tests yet)
- [x] NSIS installer script scaffold (Windows)
- [x] PKG installer scripts scaffold (macOS)
- [x] Skeleton `release.yml` CI workflow

### Commit message
`add installer scaffolding, updater and autostart packages [ci]`

---

### 1. `internal/config/constants.go`

Add after the existing `AppName` constant:

```go
const AppVersion = "1.2.0"

const (
    GitHubOwner = "AarZoooo"
    GitHubRepo  = "speedtest-tray"
)
```

Add to the log messages block:
```go
const (
    LogUpdateCheckStart   = "Update check started"
    LogUpdateFound        = "Update available"
    LogUpdateNoneFound    = "No update available"
    LogUpdateApplying     = "Applying update"
    LogUpdateCleanup      = "Cleaned up staged installer"
    ErrUpdateCheck        = "Failed to check for update"
    ErrUpdateApply        = "Failed to apply update"
    ErrUpdateDownload     = "Failed to download update"
    ErrUpdateBadChecksum  = "Downloaded installer size mismatch, aborting"
    ErrAutostartEnable    = "Failed to enable launch at login"
    ErrAutostartDisable   = "Failed to disable launch at login"
)
```

---

### 2. `internal/config/config.go`

Add `SkippedVersion` to `CustomConfig`:

```go
type CustomConfig struct {
    SaveLogs       bool   `json:"save_logs"`
    SkippedVersion string `json:"skipped_version,omitempty"`
}
```

Update `DefaultConfig`:
```go
var DefaultConfig = CustomConfig{
    SaveLogs:       false,
    SkippedVersion: "",
}
```

---

### 3. `internal/updater/` — New Package

Create `internal/updater/updater.go`:

```go
package updater

// UpdateInfo holds the result of a version check.
type UpdateInfo struct {
    LatestVersion      string
    ReleasePageURL     string
    AssetSizeBytes     int64
    HasUpdate          bool
}

// Check calls the GitHub Releases API and returns update info.
// Returns an UpdateInfo with HasUpdate=false if the current version
// is up to date or the skipped version matches the latest.
func Check(currentVersion, skippedVersion, owner, repo string) (UpdateInfo, error)

// Apply downloads the appropriate release asset for the current
// OS/arch, verifies the download size, then performs a platform-
// specific binary swap or installer launch.
//   - Windows: downloads installer.exe → %TEMP% → runs "/S" → os.Exit(0)
//   - macOS:   downloads .tar.gz → extracts binary → os.Rename →
//              on EPERM: osascript admin escalation → syscall.Exec relaunch
func Apply(info UpdateInfo) error

// CleanupStagedInstaller removes any leftover temp installer file
// from a previous update attempt. Safe to call even if no file exists.
func CleanupStagedInstaller()
```

**Implementation notes:**

- `Check()` calls `GET https://api.github.com/repos/{owner}/{repo}/releases/latest`
- Parse `tag_name` (strip leading `v`), compare with `currentVersion` using semver
  (simple string split on `.` is fine — versions follow `MAJOR.MINOR.PATCH`)
- If `latestVersion == skippedVersion`: return `HasUpdate=false`
- `ReleasePageURL` = `html_url` from the API response
- `AssetSizeBytes` = `size` from the matching asset in `assets[]`
- Asset selection: find asset whose `name` contains both the GOOS string
  (`windows`/`darwin`) and GOARCH string (`amd64`/`arm64`)
  - Windows: look for `.exe` or `.zip` asset
  - macOS: look for `.tar.gz` asset (the raw binary tarball for self-update)
- `Apply()` must verify `Content-Length` of download equals `AssetSizeBytes`
  before performing any swap. Return `ErrUpdateBadChecksum` error if mismatch.
- macOS quarantine: after extracting to temp, run
  `xattr -d com.apple.quarantine <tempfile>` before replacing the binary.
  Use `exec.Command` — ignore error (file may not be quarantined).
- Staged installer temp path:
  - Windows: `filepath.Join(os.TempDir(), "speedtest-tray-update.exe")`
  - macOS: `filepath.Join(os.TempDir(), "speedtest-tray-update")`
- `CleanupStagedInstaller()` removes those paths if they exist.

Use build tags for platform-specific Apply logic:
- `updater_windows.go` — Windows apply
- `updater_darwin.go` — macOS apply
- `updater.go` — shared Check, UpdateInfo, CleanupStagedInstaller

---

### 4. `internal/autostart/` — New Package

Create `internal/autostart/autostart.go`:

```go
package autostart

import "github.com/emersion/go-autostart"

// Manager wraps go-autostart for the current executable.
type Manager struct {
    app *autostart.App
}

// New creates a Manager pointing at the current executable.
// Returns an error if os.Executable() fails.
func New() (*Manager, error)

// IsEnabled reports whether the app is registered to launch at login.
func (m *Manager) IsEnabled() bool

// SetEnabled registers or deregisters the app from OS login items.
func (m *Manager) SetEnabled(enabled bool) error
```

**Implementation notes:**
- `autostart.App` requires `Name` and `Exec` fields
- `Name` = `config.AppName` ("SpeedTest Tray")
- `Exec` = `[]string{executablePath}` from `os.Executable()`
- On macOS the plist `Name` field becomes the label:
  `dev.aarju.speedtest-tray` — set `autostart.App.Name` to this identifier
- Add `go-autostart` to go.mod: `go get github.com/emersion/go-autostart`

---

### 5. Windows NSIS Installer — `build/windows/installer/`

Run `wails build -nsis -platform windows/amd64` locally (once, not in CI) to
generate the scaffold, then commit and customize it.

Alternatively, create `build/windows/installer/project.nsi` manually. The root
branch version is the **base installer** (install + uninstall, no PATH or
autostart yet — those are added by sub-branches):

Key NSIS directives for the base script:
```nsis
!define PRODUCT_NAME "SpeedTest Tray"
!define PRODUCT_VERSION "1.2.0"
!define PRODUCT_PUBLISHER "Aarju Pal"
!define INSTALL_DIR "$LOCALAPPDATA\Programs\SpeedTest Tray"
!define UNINSTALL_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\SpeedTest Tray"
!define EXE_NAME "speedtest-tray.exe"

; Use per-user installation — no UAC required
RequestExecutionLevel user
InstallDir "${INSTALL_DIR}"

; Pages
Page directory
Page instfiles
UninstPage uninstConfirm
UninstPage instfiles

Section "Install"
    SetOutPath "$INSTDIR"
    File "${EXE_NAME}"
    ; Start Menu shortcut
    CreateDirectory "$SMPROGRAMS\SpeedTest Tray"
    CreateShortcut "$SMPROGRAMS\SpeedTest Tray\SpeedTest Tray.lnk" "$INSTDIR\${EXE_NAME}"
    ; Register uninstaller
    WriteRegStr HKCU "${UNINSTALL_KEY}" "DisplayName" "${PRODUCT_NAME}"
    WriteRegStr HKCU "${UNINSTALL_KEY}" "UninstallString" "$INSTDIR\uninstall.exe"
    WriteRegStr HKCU "${UNINSTALL_KEY}" "DisplayVersion" "${PRODUCT_VERSION}"
    WriteRegStr HKCU "${UNINSTALL_KEY}" "Publisher" "${PRODUCT_PUBLISHER}"
    WriteUninstaller "$INSTDIR\uninstall.exe"
    ; Launch app after install
    ExecShell "" "$INSTDIR\${EXE_NAME}"
SectionEnd

Section "Uninstall"
    ; TODO: PATH removal added by cli-alias branch
    ; TODO: Run key removal added by autostart branch
    ; TODO: Data cleanup dialog added by autostart branch
    Delete "$INSTDIR\${EXE_NAME}"
    Delete "$INSTDIR\uninstall.exe"
    Delete "$SMPROGRAMS\SpeedTest Tray\SpeedTest Tray.lnk"
    RMDir "$SMPROGRAMS\SpeedTest Tray"
    RMDir "$INSTDIR"
    DeleteRegKey HKCU "${UNINSTALL_KEY}"
SectionEnd
```

Also create `build/windows/installer/info.json` (used by Wails NSIS template):
```json
{
  "companyName": "Aarju Pal",
  "productName": "SpeedTest Tray",
  "productVersion": "1.2.0",
  "copyright": "Copyright © 2026 Aarju Pal"
}
```

---

### 6. macOS PKG Installer — `build/macos/pkg/`

```
build/macos/pkg/
├── scripts/
│   └── postinstall      ← executed by macOS Installer after .app is placed
└── uninstall.sh         ← bundled into .app/Contents/Resources/, exposed via tray
```

**`build/macos/pkg/scripts/postinstall`** (base — no symlink yet, added by cli-alias):
```bash
#!/bin/bash
set -e
# Placeholder: CLI symlink added by feature/installer/cli-alias
# Placeholder: LaunchAgent setup added by feature/installer/autostart
exit 0
```
Make executable: `chmod +x build/macos/pkg/scripts/postinstall`

**`build/macos/pkg/uninstall.sh`** (base):
```bash
#!/bin/bash
set -e

APP="/Applications/SpeedTest Tray.app"

if [ ! -d "$APP" ]; then
    echo "SpeedTest Tray is not installed in /Applications."
    exit 1
fi

# Placeholder: symlink removal added by feature/installer/cli-alias
# Placeholder: LaunchAgent removal added by feature/installer/autostart

# Remove app
rm -rf "$APP"

# Placeholder: data cleanup dialog added by feature/installer/autostart

echo "SpeedTest Tray has been uninstalled."
```
Make executable: `chmod +x build/macos/pkg/uninstall.sh`

---

### 7. `wails.json` — Version Bump

Change `"version"` and `"productVersion"` from `"1.0.2"` to `"1.2.0"`.

---

### 8. `.github/workflows/release.yml` — Skeleton

Create the file with the version-sync gate and job stubs. The complete build
steps are filled in by the `feature/installer/tests` sub-branch.

```yaml
name: Release

on:
  push:
    tags:
      - 'v*.*.*'

jobs:
  verify-version:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Verify AppVersion matches tag
        run: |
          TAG_VERSION="${GITHUB_REF_NAME#v}"
          GO_VERSION=$(grep 'AppVersion' internal/config/constants.go | grep -oP '"[^"]+"' | tr -d '"')
          if [ "$TAG_VERSION" != "$GO_VERSION" ]; then
            echo "ERROR: Tag $GITHUB_REF_NAME does not match AppVersion=$GO_VERSION in constants.go"
            exit 1
          fi
          echo "Version sync OK: $GO_VERSION"

  build-windows-amd64:
    needs: verify-version
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4
      # TODO: full build steps added by feature/installer/tests

  build-windows-arm64:
    needs: verify-version
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4
      # TODO: full build steps added by feature/installer/tests

  build-macos-intel:
    needs: verify-version
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      # TODO: full build steps added by feature/installer/tests

  build-macos-arm:
    needs: verify-version
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      # TODO: full build steps added by feature/installer/tests

  create-release:
    needs: [build-windows-amd64, build-windows-arm64, build-macos-intel, build-macos-arm]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      # TODO: asset collection and release creation added by feature/installer/tests
```

---

## Sub-Branch 1: `feature/installer/cli-alias` [DONE]

### Branches from: `feature/installer`
### Merges back to: `feature/installer`
### Commit message: `add CLI alias: PATH registration (Windows) and /usr/local/bin symlink (macOS)`

### Goal
Register `speedtest-tray` as a usable terminal command after installation.

### 1. `build/windows/installer/project.nsi`

Add to the **Install** section (after Start Menu shortcut creation):
```nsis
; Add install dir to user PATH
ReadRegStr $0 HKCU "Environment" "Path"
StrCpy $0 "$0;$INSTDIR"
WriteRegExpandStr HKCU "Environment" "Path" "$0"
; Broadcast WM_SETTINGCHANGE so open terminals pick up new PATH
SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000
```

Add to the **Uninstall** section (remove the PATH entry):
```nsis
; Remove install dir from user PATH
ReadRegStr $0 HKCU "Environment" "Path"
; Strip "$INSTDIR;" from PATH string
${StrRep} $0 $0 "$INSTDIR;" ""
${StrRep} $0 $0 ";$INSTDIR" ""
WriteRegExpandStr HKCU "Environment" "Path" "$0"
SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000
```

Requires the `StrRep` NSIS macro — include at top of script:
```nsis
!include "WordFunc.nsh"
```

### 2. `build/macos/pkg/scripts/postinstall`

Replace the CLI symlink placeholder:
```bash
# Create CLI symlink
CLI_LINK="/usr/local/bin/speedtest-tray"
APP_BINARY="/Applications/SpeedTest Tray.app/Contents/MacOS/SpeedTest Tray"

if [ -L "$CLI_LINK" ]; then
    rm "$CLI_LINK"
fi
ln -sf "$APP_BINARY" "$CLI_LINK"
```

### 3. `build/macos/pkg/uninstall.sh`

Replace the symlink removal placeholder:
```bash
# Remove CLI symlink
CLI_LINK="/usr/local/bin/speedtest-tray"
if [ -L "$CLI_LINK" ]; then
    rm "$CLI_LINK"
    echo "Removed CLI symlink."
fi
```

---

## Sub-Branch 2: `feature/installer-autostart` [DONE]

### Branches from: `feature/installer`
### Merges back to: `feature/installer`
### Commit message: `add Launch at Login toggle, NSIS autostart option, LaunchAgent, and uninstall data cleanup`

### Goal
- [x] In-app Launch at Login toggle in both tray menus
- [x] NSIS: optional autostart checkbox during install
- [x] macOS: LaunchAgent plist install in postinstall
- [x] Uninstaller: data cleanup dialog on both platforms

### 1. `internal/gui_wails/app.go`

Add new Wails-bound methods and update startup sequence.

**New imports needed:** `speedtest-tray/internal/autostart`, `speedtest-tray/internal/updater`

**New fields on `App` struct:**
```go
autostartMgr *autostart.Manager
updateInfo   updater.UpdateInfo
```

**New Wails-bound methods:**
```go
func (a *App) GetLaunchAtLogin() bool {
    if a.autostartMgr == nil {
        return false
    }
    return a.autostartMgr.IsEnabled()
}

func (a *App) SetLaunchAtLogin(enabled bool) {
    if a.autostartMgr == nil {
        return
    }
    if err := a.autostartMgr.SetEnabled(enabled); err != nil {
        slog.Error(config.ErrAutostartEnable, config.KeyError, err)
    }
}

func (a *App) GetUpdateInfo() updater.UpdateInfo {
    return a.updateInfo
}

func (a *App) ApplyUpdate() {
    if err := updater.Apply(a.updateInfo); err != nil {
        slog.Error(config.ErrUpdateApply, config.KeyError, err)
    }
}

func (a *App) SkipUpdate(version string) {
    cfg := config.LoadConfigOrDefault()
    cfg.SkippedVersion = version
    config.SaveConfig(cfg)
}
```

**Update startup (in `startup` method or wherever the app initialises):**
```go
updater.CleanupStagedInstaller()

mgr, err := autostart.New()
if err != nil {
    slog.Error("Failed to init autostart manager", config.KeyError, err)
} else {
    a.autostartMgr = mgr
}

go func() {
    cfg := config.LoadConfigOrDefault()
    info, err := updater.Check(
        config.AppVersion, cfg.SkippedVersion,
        config.GitHubOwner, config.GitHubRepo,
    )
    if err != nil {
        slog.Error(config.ErrUpdateCheck, config.KeyError, err)
        return
    }
    a.updateInfo = info
    if info.HasUpdate {
        slog.Info(config.LogUpdateFound, "version", info.LatestVersion)
        runtime.EventsEmit(a.ctx, "update:available", info)
    } else {
        slog.Info(config.LogUpdateNoneFound)
    }
}()
```

### 2. `internal/gui_wails/tray_windows.go`

Add Launch at Login checkbox. Update `StartTray` signature:
```go
func StartTray(app *App, iconBytes []byte, macIconBytes []byte,
    appConfig *config.CustomConfig,
    toggleLogging func(bool),
    toggleLaunchAtLogin func(bool))
```

Add menu item after the `show` item's separator:
```go
launchAtLogin := systray.AddMenuItemCheckbox(
    "Launch at Login",
    "Start SpeedTest Tray automatically on login",
    app.GetLaunchAtLogin(),
)
launchAtLogin.Click(func() {
    enabled := !launchAtLogin.Checked()
    if enabled {
        launchAtLogin.Check()
    } else {
        launchAtLogin.Uncheck()
    }
    toggleLaunchAtLogin(enabled)
})

systray.AddSeparator()
```

### 3. `internal/gui_wails/tray_darwin.go` + `tray_darwin.m`

Add `toggleLaunchAtLoginCallback func(bool)` package-level var.

Add new exported Go callback:
```go
//export onLaunchAtLoginClick
func onLaunchAtLoginClick(enabled C.int) {
    slog.Info("onLaunchAtLoginClick received from Objective-C", "enabled", enabled != 0)
    if toggleLaunchAtLoginCallback != nil {
        go toggleLaunchAtLoginCallback(enabled != 0)
    }
}
```

Update `StartTray` signature to match Windows (gains `toggleLaunchAtLogin func(bool)`).
Set `toggleLaunchAtLoginCallback = toggleLaunchAtLogin` in `StartTray`.

Update `initStatusItem` C declaration and call to pass `initialLaunchAtLoginState int`.

In `tray_darwin.m`: add a checkmark NSMenuItem for "Launch at Login" with the
`onLaunchAtLoginClick` target, positioned between Show and the logging item.

### 4. `build/windows/installer/project.nsi`

Add optional autostart checkbox during install. At the top add:
```nsis
!define RUN_KEY "Software\Microsoft\Windows\CurrentVersion\Run"
Var LaunchAtLoginCheckbox

Section "Install"
    ; ... existing install code ...

    ; Launch at Login checkbox
    ${If} $LaunchAtLoginCheckbox == ${BST_CHECKED}
        WriteRegStr HKCU "${RUN_KEY}" "${PRODUCT_NAME}" "$INSTDIR\${EXE_NAME}"
    ${EndIf}
SectionEnd
```

Add to installer UI page (custom page or using nsDialogs for a checkbox).

Add to Uninstall section:
```nsis
; Remove autostart Run key if present
DeleteRegValue HKCU "${RUN_KEY}" "${PRODUCT_NAME}"
```

Add data cleanup dialog to Uninstall section:
```nsis
MessageBox MB_YESNO "Remove configuration, logs, and history? \
    This includes config.json, app.log, and history.json." \
    IDNO SkipDataRemoval
    RMDir /r "$APPDATA\SpeedTest Tray"
SkipDataRemoval:
```

### 5. `build/macos/pkg/scripts/postinstall`

Replace LaunchAgent placeholder:
```bash
# Install LaunchAgent for current user
PLIST_DIR="$HOME/Library/LaunchAgents"
PLIST_PATH="$PLIST_DIR/dev.aarju.speedtest-tray.plist"

mkdir -p "$PLIST_DIR"
cat > "$PLIST_PATH" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
    "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>dev.aarju.speedtest-tray</string>
    <key>ProgramArguments</key>
    <array>
        <string>/Applications/SpeedTest Tray.app/Contents/MacOS/SpeedTest Tray</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
EOF
launchctl load "$PLIST_PATH" 2>/dev/null || true
```

> Note: The LaunchAgent is installed for every PKG install. The in-app toggle
> (`SetLaunchAtLogin`) manages it at runtime. If the user toggles it off in-app,
> the plist is removed by `go-autostart`.

### 6. `build/macos/pkg/uninstall.sh`

Replace LaunchAgent removal and add data cleanup placeholder:
```bash
# Remove LaunchAgent
PLIST="$HOME/Library/LaunchAgents/dev.aarju.speedtest-tray.plist"
if [ -f "$PLIST" ]; then
    launchctl unload "$PLIST" 2>/dev/null || true
    rm "$PLIST"
    echo "Removed LaunchAgent."
fi

# Optional: remove app data
read -r -p "Remove configuration, logs, and history? [y/N] " response
if [[ "$response" =~ ^[Yy]$ ]]; then
    DATA_DIR="$HOME/Library/Application Support/SpeedTest Tray"
    if [ -d "$DATA_DIR" ]; then
        rm -rf "$DATA_DIR"
        echo "Removed app data."
    fi
fi
```

---

## Sub-Branch 3: `feature/installer/tests` [DONE]

### Branches from: `feature/installer`
### Merges back to: `feature/installer`
### Commit message: `add updater and autostart unit tests, complete release.yml`

### Goal
Full unit test coverage for the updater and autostart packages, plus the
complete `release.yml` workflow with all 8 release assets.

### 1. `internal/updater/updater_test.go`

All tests use `httptest.NewServer` for mock GitHub API. No real network calls.
No disk writes outside `t.TempDir()`.

Test matrix (implement all of these):

| Test name | Setup | Assert |
|-----------|-------|--------|
| `TestCheck_NewerVersion` | Mock returns tag `v9.9.9` | `HasUpdate=true`, `LatestVersion="9.9.9"` |
| `TestCheck_SameVersion` | Mock returns tag matching `AppVersion` | `HasUpdate=false` |
| `TestCheck_OlderVersion` | Mock returns tag `v0.0.1` | `HasUpdate=false` |
| `TestCheck_SkippedVersion` | Latest = `v9.9.9`, skipped = `9.9.9` | `HasUpdate=false` |
| `TestCheck_NetworkError` | Mock server closes conn immediately | Returns error, no panic |
| `TestCheck_MalformedJSON` | Mock returns `{invalid` | Returns error |
| `TestCheck_NoMatchingAsset` | Mock has no asset matching GOOS/GOARCH | Returns error |
| `TestAssetSelection_WindowsAmd64` | Assets list with various names | Selects `*windows*amd64*` |
| `TestAssetSelection_WindowsArm64` | Assets list | Selects `*windows*arm64*` |
| `TestAssetSelection_DarwinAmd64` | Assets list | Selects `*darwin*amd64*` |
| `TestAssetSelection_DarwinArm64` | Assets list | Selects `*darwin*arm64*` |
| `TestAssetSizeBytes` | Mock returns `Content-Length: 12345` | `AssetSizeBytes == 12345` |
| `TestCleanupStagedInstaller_FileExists` | Create file at staged path | File removed after call |
| `TestCleanupStagedInstaller_NoFile` | No file at staged path | No error |

For download size verification (simulated without actually calling `Apply` end-to-end):
```go
func TestApply_ContentLengthMismatch(t *testing.T) {
    // Mock server returns Content-Length: 100 but only sends 10 bytes
    // Apply() should return ErrUpdateBadChecksum
}
```

### 2. `internal/autostart/autostart_test.go`

Use build tag to guard OS writes:
```go
//go:build !ci
```

Tests:
```go
func TestAutostartRoundTrip(t *testing.T) {
    mgr, err := New()
    require.NoError(t, err)

    // Ensure disabled at start
    _ = mgr.SetEnabled(false)
    assert.False(t, mgr.IsEnabled())

    // Enable
    require.NoError(t, mgr.SetEnabled(true))
    assert.True(t, mgr.IsEnabled())

    // Disable again (cleanup)
    require.NoError(t, mgr.SetEnabled(false))
    assert.False(t, mgr.IsEnabled())
}
```

### 3. Complete `.github/workflows/release.yml`

Fill in the `# TODO` stubs from the root branch. Final workflow:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*.*.*'

jobs:
  verify-version:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Verify AppVersion matches tag
        run: |
          TAG_VERSION="${GITHUB_REF_NAME#v}"
          GO_VERSION=$(grep 'AppVersion' internal/config/constants.go | grep -oP '"[^"]+"' | tr -d '"')
          if [ "$TAG_VERSION" != "$GO_VERSION" ]; then
            echo "ERROR: Tag $GITHUB_REF_NAME does not match AppVersion=$GO_VERSION"
            exit 1
          fi

  build-windows-amd64:
    needs: verify-version
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.26' }
      - name: Install Wails
        run: go install github.com/wailsapp/wails/v2/cmd/wails@latest
      - name: Install NSIS
        run: choco install nsis -y
      - name: Generate frontend config
        run: go generate ./...
      - name: Build Windows amd64 installer
        run: wails build -platform windows/amd64 -nsis -o speedtest-tray.exe -clean
      - name: Rename installer
        run: mv build/bin/speedtest-tray-amd64-installer.exe SpeedTest-Tray-windows-amd64-installer.exe
      - uses: actions/upload-artifact@v4
        with:
          name: windows-amd64
          path: SpeedTest-Tray-windows-amd64-installer.exe

  build-windows-arm64:
    needs: verify-version
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.26' }
      - name: Install Wails
        run: go install github.com/wailsapp/wails/v2/cmd/wails@latest
      - name: Install NSIS
        run: choco install nsis -y
      - name: Generate frontend config
        run: go generate ./...
      - name: Build Windows arm64 installer
        run: wails build -platform windows/arm64 -nsis -o speedtest-tray.exe -clean
      - name: Rename installer
        run: mv build/bin/speedtest-tray-arm64-installer.exe SpeedTest-Tray-windows-arm64-installer.exe
      - uses: actions/upload-artifact@v4
        with:
          name: windows-arm64
          path: SpeedTest-Tray-windows-arm64-installer.exe

  build-macos-intel:
    needs: verify-version
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.26' }
      - name: Install Wails and tools
        run: |
          go install github.com/wailsapp/wails/v2/cmd/wails@latest
          brew install create-dmg fileicon
      - name: Generate frontend config
        run: go generate ./...
      - name: Build macOS Intel
        run: wails build -platform darwin/amd64 -clean
      - name: Build PKG
        run: |
          pkgbuild --root "build/bin" \
                   --identifier dev.aarju.speedtest-tray \
                   --version ${{ github.ref_name }} \
                   --scripts build/macos/pkg/scripts \
                   SpeedTest-Tray-macOS-Intel.pkg
      - name: Build DMG
        run: |
          create-dmg \
            --volname "SpeedTest Tray" \
            --window-pos 200 120 \
            --window-size 600 400 \
            --icon-size 100 \
            --icon "SpeedTest Tray.app" 175 190 \
            --hide-extension "SpeedTest Tray.app" \
            --app-drop-link 425 190 \
            "SpeedTest-Tray-macOS-Intel.dmg" \
            "build/bin/"
      - name: Package raw binary for self-update
        run: |
          tar -czf speedtest-tray-darwin-amd64.tar.gz \
            -C "build/bin/SpeedTest Tray.app/Contents/MacOS" "SpeedTest Tray"
      - uses: actions/upload-artifact@v4
        with:
          name: macos-intel
          path: |
            SpeedTest-Tray-macOS-Intel.pkg
            SpeedTest-Tray-macOS-Intel.dmg
            speedtest-tray-darwin-amd64.tar.gz

  build-macos-arm:
    needs: verify-version
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.26' }
      - name: Install Wails and tools
        run: |
          go install github.com/wailsapp/wails/v2/cmd/wails@latest
          brew install create-dmg fileicon
      - name: Generate frontend config
        run: go generate ./...
      - name: Build macOS ARM
        run: wails build -platform darwin/arm64 -clean
      - name: Build PKG
        run: |
          pkgbuild --root "build/bin" \
                   --identifier dev.aarju.speedtest-tray \
                   --version ${{ github.ref_name }} \
                   --scripts build/macos/pkg/scripts \
                   SpeedTest-Tray-macOS-ARM.pkg
      - name: Build DMG
        run: |
          create-dmg \
            --volname "SpeedTest Tray" \
            --window-pos 200 120 \
            --window-size 600 400 \
            --icon-size 100 \
            --icon "SpeedTest Tray.app" 175 190 \
            --hide-extension "SpeedTest Tray.app" \
            --app-drop-link 425 190 \
            "SpeedTest-Tray-macOS-ARM.dmg" \
            "build/bin/"
      - name: Package raw binary for self-update
        run: |
          tar -czf speedtest-tray-darwin-arm64.tar.gz \
            -C "build/bin/SpeedTest Tray.app/Contents/MacOS" "SpeedTest Tray"
      - uses: actions/upload-artifact@v4
        with:
          name: macos-arm
          path: |
            SpeedTest-Tray-macOS-ARM.pkg
            SpeedTest-Tray-macOS-ARM.dmg
            speedtest-tray-darwin-arm64.tar.gz

  create-release:
    needs: [build-windows-amd64, build-windows-arm64, build-macos-intel, build-macos-arm]
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/download-artifact@v4
        with: { path: artifacts }
      - uses: softprops/action-gh-release@v2
        with:
          files: |
            artifacts/windows-amd64/SpeedTest-Tray-windows-amd64-installer.exe
            artifacts/windows-arm64/SpeedTest-Tray-windows-arm64-installer.exe
            artifacts/macos-intel/SpeedTest-Tray-macOS-Intel.pkg
            artifacts/macos-intel/SpeedTest-Tray-macOS-Intel.dmg
            artifacts/macos-intel/speedtest-tray-darwin-amd64.tar.gz
            artifacts/macos-arm/SpeedTest-Tray-macOS-ARM.pkg
            artifacts/macos-arm/SpeedTest-Tray-macOS-ARM.dmg
            artifacts/macos-arm/speedtest-tray-darwin-arm64.tar.gz
          generate_release_notes: true
```

## Sub-Branch 4: `feature/installer/update-ui` [DONE]

### Branches from: `feature/installer`
### Merges back to: `feature/installer`
### Commit message: `add update panel UI, update-toggle-btn, badge, skip and install update interaction`

### Goal
Implement the frontend update UI:
- Add a new header action button: `update-toggle-btn` to the left of history button.
- Toggle-able update panel/view that matches the styling and layout of `history-view`.
- Handle the `update:available` event from Wails and update the button/view dynamically.
- Implement click handlers for Update and Skip buttons.

### 1. `frontend/index.html`

In the header actions container:
```html
<div class="header-actions">
    <!-- Add before history-toggle-btn -->
    <button
        id="update-toggle-btn"
        class="icon-btn update-btn-hidden"
        title="Update Available"
    >
        <svg class="icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M4 12A8 8 0 0 1 18.93 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M20 12A8 8 0 0 1 5.07 16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <polyline points="14 8 19 8 19 3" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <polyline points="10 16 5 16 5 21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
    </button>
    <button
        id="history-toggle-btn"
...
```

Add the `update-view` inside the `<main class="main-content">` tag, below the `history-view`:
```html
<!-- Update View -->
<div id="update-view" class="view-hidden">
    <div class="update-card">
        <h2 class="update-title">New Version Available!</h2>
        <div class="update-info-group">
            <div class="update-info-row">
                <span class="update-info-label">Version:</span>
                <span id="update-version-val" class="update-info-val">v1.2.0</span>
            </div>
            <div class="update-info-row">
                <span class="update-info-label">File Size:</span>
                <span id="update-size-val" class="update-info-val">--</span>
            </div>
        </div>
        <div class="update-actions">
            <button id="update-now-btn" class="update-action-btn primary-btn">Update Now</button>
            <button id="update-skip-btn" class="update-action-btn secondary-btn">Skip Version</button>
        </div>
        <div class="update-notes-link">
            <a href="#" id="update-notes-btn">View Release Notes</a>
        </div>
    </div>
</div>
```

### 2. `frontend/style.css`

Add update UI styles:
```css
/* Update Button styling with badge */
#update-toggle-btn {
    position: relative;
    transition: color 0.2s ease;
}

.update-btn-hidden {
    display: none !important;
}

.update-btn-visible {
    display: block !important;
}

#update-toggle-btn.has-badge::after {
    content: '';
    position: absolute;
    top: 2px;
    right: 2px;
    width: 8px;
    height: 8px;
    background-color: var(--danger);
    border-radius: 50%;
    border: 1.5px solid var(--bg-color);
}

/* Update Panel */
.update-card {
    background-color: var(--card-bg);
    border-radius: 12px;
    padding: 24px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
    border: 1px solid var(--border-dim);
    margin: auto 0;
}

.update-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-main);
    text-align: center;
}

.update-info-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
    width: 100%;
    max-width: 200px;
    margin: 8px 0;
}

.update-info-row {
    display: flex;
    justify-content: space-between;
    font-size: 0.9rem;
}

.update-info-label {
    color: var(--text-dim);
}

.update-info-val {
    color: var(--text-main);
    font-weight: 500;
}

.update-actions {
    display: flex;
    gap: 12px;
    width: 100%;
}

.update-action-btn {
    flex: 1;
    padding: 10px;
    border-radius: 8px;
    border: none;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s, transform 0.1s;
}

.update-action-btn:active {
    transform: scale(0.98);
}

.primary-btn {
    background: linear-gradient(135deg, var(--accent-start), var(--accent-end));
    color: white;
    box-shadow: var(--shadow-btn);
}

.primary-btn:hover {
    filter: brightness(1.1);
    box-shadow: var(--shadow-btn-hover);
}

.secondary-btn {
    background-color: transparent;
    border: 1px solid var(--border-dim);
    color: var(--text-dim);
}

.secondary-btn:hover {
    background-color: rgba(255, 255, 255, 0.05);
    color: var(--text-main);
}

.update-notes-link {
    margin-top: 8px;
}

.update-notes-link a {
    color: var(--text-dim);
    font-size: 0.85rem;
    text-decoration: underline;
    cursor: pointer;
}

.update-notes-link a:hover {
    color: var(--text-main);
}

.update-status-msg {
    font-size: 0.9rem;
    color: var(--text-dim);
    text-align: center;
    margin: 12px 0;
}
```

### 3. `frontend/src/handlers.js`

Add and export update view handling functions:

```javascript
let updateData = null;

export function handleUpdateAvailable(info) {
    updateData = info;
    const toggleBtn = document.getElementById("update-toggle-btn");
    const versionVal = document.getElementById("update-version-val");
    const sizeVal = document.getElementById("update-size-val");

    if (toggleBtn) {
        toggleBtn.classList.replace("update-btn-hidden", "update-btn-visible");
        toggleBtn.classList.add("has-badge");
    }

    if (versionVal) {
        versionVal.innerText = "v" + info.LatestVersion;
    }

    if (sizeVal) {
        // Format size: e.g. 4.52 MB
        const mb = (info.AssetSizeBytes / (1024 * 1024)).toFixed(2);
        sizeVal.innerText = mb + " MB";
    }
}

export function handleUpdateToggleClick() {
    if (testState.isTesting) return;

    const testView = document.getElementById("test-view");
    const updateView = document.getElementById("update-view");
    const historyView = document.getElementById("history-view");
    const toggleBtn = document.getElementById("update-toggle-btn");

    if (!testView || !updateView || !toggleBtn) return;

    // Remove badge on click
    toggleBtn.classList.remove("has-badge");

    if (testView.classList.contains("view-active")) {
        testView.classList.replace("view-active", "view-hidden");
        updateView.classList.replace("view-hidden", "view-active");
        toggleBtn.innerHTML = SPEED_ICON_HTML;
        toggleBtn.title = "Run Test";
    } else {
        updateView.classList.replace("view-active", "view-hidden");
        testView.classList.replace("view-hidden", "view-active");
        toggleBtn.innerHTML = `<svg class="icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M4 12A8 8 0 0 1 18.93 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M20 12A8 8 0 0 1 5.07 16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <polyline points="14 8 19 8 19 3" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <polyline points="10 16 5 16 5 21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>`;
        toggleBtn.title = "Update Available";
    }
}

export function handleUpdateNowClick() {
    const card = document.querySelector(".update-card");
    if (!card) return;

    // Swap content with spinner and progress message
    card.innerHTML = `
        <div class="loader"></div>
        <div class="update-status-msg">Downloading and installing update...<br>App will restart automatically.</div>
    `;

    // Trigger backend apply
    window.go.gui_wails.App.ApplyUpdate();
}

export function handleUpdateSkipClick() {
    if (!updateData) return;

    // Notify backend to save skipped version
    window.go.gui_wails.App.SkipUpdate(updateData.LatestVersion);

    // Hide update toggle button and return to test view
    const toggleBtn = document.getElementById("update-toggle-btn");
    if (toggleBtn) {
        toggleBtn.classList.replace("update-btn-visible", "update-btn-hidden");
    }

    const testView = document.getElementById("test-view");
    const updateView = document.getElementById("update-view");
    if (testView && updateView && updateView.classList.contains("view-active")) {
        updateView.classList.replace("view-active", "view-hidden");
        testView.classList.replace("view-hidden", "view-active");
    }
}

export function handleReleaseNotesClick(e) {
    e.preventDefault();
    if (updateData && updateData.ReleasePageURL) {
        window.runtime.BrowserOpenURL(updateData.ReleasePageURL);
    }
}

export function initializeUpdateHandlers() {
    const toggleBtn = document.getElementById("update-toggle-btn");
    if (toggleBtn) {
        toggleBtn.addEventListener("click", handleUpdateToggleClick);
    }

    const updateNowBtn = document.getElementById("update-now-btn");
    if (updateNowBtn) {
        updateNowBtn.addEventListener("click", handleUpdateNowClick);
    }

    const updateSkipBtn = document.getElementById("update-skip-btn");
    if (updateSkipBtn) {
        updateSkipBtn.addEventListener("click", handleUpdateSkipClick);
    }

    const notesBtn = document.getElementById("update-notes-btn");
    if (notesBtn) {
        notesBtn.addEventListener("click", handleReleaseNotesClick);
    }
}
```

### 4. `frontend/main.js`

Import and initialize update handlers and register Wails runtime event:

```javascript
import { initializeUpdateHandlers, handleUpdateAvailable } from "./src/handlers.js";

// Inside DOMContentLoaded listener:
initializeUpdateHandlers();

// Register the update event:
window.runtime.EventsOn("update:available", handleUpdateAvailable);
```

---

## Merge Order

```
feature/installer/cli-alias  ──┐
feature/installer-autostart  ──┼──► feature/installer ──► master
feature/installer/update-ui  ──┤
feature/installer/tests      ──┘
```

No merge order dependency between the four sub-branches. Merge conflicts
are expected to be minimal and easily resolvable.

After all four are merged into `feature/installer`, run the full manual QA
checklist from the implementation_plan.md artifact before merging to master.


---

## Release Assets Summary

| Asset | Platform | Purpose |
|-------|----------|---------|
| `SpeedTest-Tray-windows-amd64-installer.exe` | Windows amd64 | Fresh install + self-update |
| `SpeedTest-Tray-windows-arm64-installer.exe` | Windows arm64 | Fresh install + self-update |
| `SpeedTest-Tray-macOS-Intel.pkg` | macOS Intel | Fresh install |
| `SpeedTest-Tray-macOS-ARM.pkg` | macOS ARM | Fresh install |
| `speedtest-tray-darwin-amd64.tar.gz` | macOS Intel | In-app self-update only |
| `speedtest-tray-darwin-arm64.tar.gz` | macOS ARM | In-app self-update only |
| `SpeedTest-Tray-macOS-Intel.dmg` | macOS Intel | Optional drag-install |
| `SpeedTest-Tray-macOS-ARM.dmg` | macOS ARM | Optional drag-install |
