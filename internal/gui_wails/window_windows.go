//go:build windows

package gui_wails

import (
	"syscall"
	"unsafe"

	"speedtest-tray/internal/config"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	user32                 = syscall.NewLazyDLL("user32.dll")
	gdi32                  = syscall.NewLazyDLL("gdi32.dll")
	shcore                 = syscall.NewLazyDLL("shcore.dll")
	procGetCursorPos       = user32.NewProc("GetCursorPos")
	procFindWindowW        = user32.NewProc("FindWindowW")
	procCreateRoundRectRgn = gdi32.NewProc("CreateRoundRectRgn")
	procSetWindowRgn       = user32.NewProc("SetWindowRgn")
	procGetDpiForWindow    = user32.NewProc("GetDpiForWindow")
	procGetDpiForSystem    = user32.NewProc("GetDpiForSystem")
	procMonitorFromPoint   = user32.NewProc("MonitorFromPoint")
	procGetDpiForMonitor   = shcore.NewProc("GetDpiForMonitor")
)

type point struct {
	x int32
	y int32
}

func (a *App) getScaleFactorAtPoint(x, y int32) float64 {
	packedPoint := uint64(uint32(x)) | (uint64(uint32(y)) << 32)
	monitor, _, _ := procMonitorFromPoint.Call(uintptr(packedPoint), uintptr(config.MonitorDefaultToNearest))

	if monitor != 0 {
		var dpiX, dpiY uint32
		procGetDpiForMonitor.Call(
			monitor,
			uintptr(config.MdtEffectiveDpi),
			uintptr(unsafe.Pointer(&dpiX)),
			uintptr(unsafe.Pointer(&dpiY)),
		)
		if dpiX != 0 {
			return float64(dpiX) / config.StandardDPI
		}
	}

	if procGetDpiForSystem.Find() == nil {
		dpi, _, _ := procGetDpiForSystem.Call()
		if dpi != 0 {
			return float64(dpi) / config.StandardDPI
		}
	}

	return 1.0
}

func (a *App) getScaleFactorForWindow(hwnd uintptr) float64 {
	if hwnd != 0 && procGetDpiForWindow.Find() == nil {
		dpi, _, _ := procGetDpiForWindow.Call(hwnd)
		if dpi != 0 {
			return float64(dpi) / config.StandardDPI
		}
	}

	if procGetDpiForSystem.Find() == nil {
		dpi, _, _ := procGetDpiForSystem.Call()
		if dpi != 0 {
			return float64(dpi) / config.StandardDPI
		}
	}

	return 1.0
}

func (a *App) positionWindow() {
	var cursor point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&cursor)))

	scaleFactor := a.getScaleFactorAtPoint(cursor.x, cursor.y)

	scaledWidth := int(float64(config.WindowWidth) * scaleFactor)
	scaledHeight := int(float64(config.WindowHeight) * scaleFactor)
	scaledOffsetY := int(float64(config.WindowOffsetYPixels) * scaleFactor)

	x := int(cursor.x) - scaledWidth/2
	y := int(cursor.y) - scaledHeight + scaledOffsetY

	if x < 0 {
		x = 0
	}

	wailsRuntime.WindowSetSize(a.ctx, scaledWidth, scaledHeight)
	wailsRuntime.WindowSetPosition(a.ctx, x, y)
}

func (a *App) ApplyRoundedCorners() {
	title, _ := syscall.UTF16PtrFromString(config.AppName)
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(title)))
	if hwnd == 0 {
		return
	}

	scaleFactor := a.getScaleFactorForWindow(hwnd)

	scaledWidth := int(float64(config.WindowWidth) * scaleFactor)
	scaledHeight := int(float64(config.WindowHeight) * scaleFactor)
	scaledRadius := int(float64(config.WindowCornerRadius) * scaleFactor)

	region, _, _ := procCreateRoundRectRgn.Call(
		0,
		0,
		uintptr(scaledWidth),
		uintptr(scaledHeight),
		uintptr(scaledRadius),
		uintptr(scaledRadius),
	)
	procSetWindowRgn.Call(hwnd, region, 1)
}

func (a *App) focusApp() {}
