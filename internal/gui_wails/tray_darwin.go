//go:build darwin

package gui_wails

/*
#cgo LDFLAGS: -framework AppKit -framework Foundation
#include <stdlib.h>
void initStatusItem(const char* title, const void* iconData, int iconLength, int initialLoggingState, int initialLaunchAtLoginState);
void getStatusItemPosition(double *x, double *y, double *width, double *height, double *screenWidth);
void removeStatusItem(void);
*/
import "C"
import (
	"log/slog"
	"speedtest-tray/internal/config"
	"unsafe"
)

var (
	globalApp                   *App
	toggleLoggingCallback       func(bool)
	toggleLaunchAtLoginCallback func(bool)
	initialLoggingVal           bool
	initialLaunchAtLoginVal     bool
)

//export onStatusItemClick
func onStatusItemClick() {
	slog.Info("onStatusItemClick received from Objective-C")
	if globalApp != nil {
		go globalApp.ToggleWindow()
	}
}

//export onQuitClick
func onQuitClick() {
	slog.Info("onQuitClick received from Objective-C")
	if globalApp != nil {
		globalApp.Quit()
	}
}

//export onToggleLoggingClick
func onToggleLoggingClick(enabled C.int) {
	slog.Info("onToggleLoggingClick received from Objective-C", "enabled", enabled != 0)
	if toggleLoggingCallback != nil {
		go toggleLoggingCallback(enabled != 0)
	}
}

//export onLaunchAtLoginClick
func onLaunchAtLoginClick(enabled C.int) {
	slog.Info("onLaunchAtLoginClick received from Objective-C", "enabled", enabled != 0)
	if toggleLaunchAtLoginCallback != nil {
		go toggleLaunchAtLoginCallback(enabled != 0)
	}
}

//export onOpenLogsClick
func onOpenLogsClick() {
	slog.Info("onOpenLogsClick received from Objective-C")
	config.OpenDirectory(config.GetConfigDir())
}

func StartTray(app *App, iconBytes []byte, macIconBytes []byte, appConfig *config.CustomConfig, toggleLogging func(bool), toggleLaunchAtLogin func(bool)) {
	slog.Info("StartTray called on macOS")
	app.MacIcon = macIconBytes
	toggleLoggingCallback = toggleLogging
	toggleLaunchAtLoginCallback = toggleLaunchAtLogin
	initialLoggingVal = appConfig.SaveLogs
	initialLaunchAtLoginVal = app.GetLaunchAtLogin()
}

func (a *App) initMacStatusItem() {
	slog.Info("initMacStatusItem starting")
	globalApp = a
	title := C.CString(config.AppName)
	defer C.free(unsafe.Pointer(title))

	var iconPtr unsafe.Pointer
	var iconLen C.int
	if len(a.MacIcon) > 0 {
		iconPtr = unsafe.Pointer(&a.MacIcon[0])
		iconLen = C.int(len(a.MacIcon))
		slog.Info("initMacStatusItem passing icon bytes", "len", len(a.MacIcon))
	} else {
		slog.Info("initMacStatusItem has no icon bytes")
	}

	loggingVal := C.int(0)
	if initialLoggingVal {
		loggingVal = C.int(1)
	}

	launchAtLoginVal := C.int(0)
	if initialLaunchAtLoginVal {
		launchAtLoginVal = C.int(1)
	}

	C.initStatusItem(title, iconPtr, iconLen, loggingVal, launchAtLoginVal)
}

func (a *App) cleanupMacStatusItem() {
	slog.Info("cleanupMacStatusItem called")
	C.removeStatusItem()
}
