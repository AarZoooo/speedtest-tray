//go:build darwin

package gui_wails

/*
#cgo LDFLAGS: -framework AppKit -framework Foundation
#include <stdlib.h>
void initStatusItem(const char* title, const void* iconData, int iconLength, int initialLoggingState);
void getStatusItemPosition(double *x, double *y, double *width, double *height, double *screenWidth);
void removeStatusItem(void);
*/
import "C"
import (
	"log/slog"
	"unsafe"
	"speedtest-tray/internal/config"
)

var (
	globalApp             *App
	toggleLoggingCallback func(bool)
	initialLoggingVal     bool
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

func StartTray(app *App, iconBytes []byte, macIconBytes []byte, appConfig *config.CustomConfig, toggleLogging func(bool)) {
	slog.Info("StartTray called on macOS")
	app.MacIcon = macIconBytes
	toggleLoggingCallback = toggleLogging
	initialLoggingVal = appConfig.SaveLogs
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

	C.initStatusItem(title, iconPtr, iconLen, loggingVal)
}

func (a *App) cleanupMacStatusItem() {
	slog.Info("cleanupMacStatusItem called")
	C.removeStatusItem()
}
