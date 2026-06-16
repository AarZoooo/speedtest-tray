//go:build windows

package gui_wails

import (
	"syscall"
	"unsafe"

	"speedtest-tray/internal/config"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	user32                     = syscall.NewLazyDLL("user32.dll")
	gdi32                      = syscall.NewLazyDLL("gdi32.dll")
	shcore                     = syscall.NewLazyDLL("shcore.dll")
	shell32                    = syscall.NewLazyDLL("shell32.dll")
	procGetCursorPos           = user32.NewProc("GetCursorPos")
	procFindWindowW            = user32.NewProc("FindWindowW")
	procCreateRoundRectRgn     = gdi32.NewProc("CreateRoundRectRgn")
	procSetWindowRgn           = user32.NewProc("SetWindowRgn")
	procGetDpiForWindow        = user32.NewProc("GetDpiForWindow")
	procGetDpiForSystem        = user32.NewProc("GetDpiForSystem")
	procMonitorFromPoint       = user32.NewProc("MonitorFromPoint")
	procGetDpiForMonitor       = shcore.NewProc("GetDpiForMonitor")
	procShellNotifyIconGetRect = shell32.NewProc("Shell_NotifyIconGetRect")
	procGetMonitorInfoW        = user32.NewProc("GetMonitorInfoW")
)

type point struct {
	x int32
	y int32
}

type rect struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

type notifyIconIdentifier struct {
	cbSize   uint32
	hWnd     uintptr
	uID      uint32
	guidItem [16]byte
}

type monitorInfo struct {
	cbSize    uint32
	rcMonitor rect
	rcWork    rect
	dwFlags   uint32
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
	className, _ := syscall.UTF16PtrFromString("SystrayClass")
	systrayHwnd, _, _ := procFindWindowW.Call(uintptr(unsafe.Pointer(className)), 0)

	var trayRect rect
	var success bool

	if systrayHwnd != 0 {
		var identifier notifyIconIdentifier
		identifier.cbSize = uint32(unsafe.Sizeof(identifier))
		identifier.hWnd = systrayHwnd
		identifier.uID = config.SystrayIconID

		ret, _, _ := procShellNotifyIconGetRect.Call(
			uintptr(unsafe.Pointer(&identifier)),
			uintptr(unsafe.Pointer(&trayRect)),
		)
		if ret == 0 {
			success = true
		}
	}

	var x, y int
	var scaleFactor float64

	if success {
		iconCenterX := (trayRect.left + trayRect.right) / 2
		iconCenterY := (trayRect.top + trayRect.bottom) / 2

		scaleFactor = a.getScaleFactorAtPoint(iconCenterX, iconCenterY)

		scaledWidth := int(float64(config.WindowWidth) * scaleFactor)
		scaledHeight := int(float64(config.WindowHeight) * scaleFactor)
		gap := int(float64(config.WindowTrayGap) * scaleFactor)

		packedPoint := uint64(uint32(iconCenterX)) | (uint64(uint32(iconCenterY)) << 32)
		monitor, _, _ := procMonitorFromPoint.Call(uintptr(packedPoint), uintptr(config.MonitorDefaultToNearest))

		var mi monitorInfo
		mi.cbSize = uint32(unsafe.Sizeof(mi))
		getMonInfoRet, _, _ := procGetMonitorInfoW.Call(monitor, uintptr(unsafe.Pointer(&mi)))

		if getMonInfoRet != 0 {
			var targetX, targetY int32

			if mi.rcWork.bottom < mi.rcMonitor.bottom {
				targetX = iconCenterX - int32(scaledWidth)/2
				targetY = mi.rcWork.bottom - int32(scaledHeight) - int32(gap)
			} else if mi.rcWork.top > mi.rcMonitor.top {
				targetX = iconCenterX - int32(scaledWidth)/2
				targetY = mi.rcWork.top + int32(gap)
			} else if mi.rcWork.right < mi.rcMonitor.right {
				targetX = mi.rcWork.right - int32(scaledWidth) - int32(gap)
				targetY = iconCenterY - int32(scaledHeight)/2
			} else if mi.rcWork.left > mi.rcMonitor.left {
				targetX = mi.rcWork.left + int32(gap)
				targetY = iconCenterY - int32(scaledHeight)/2
			} else {
				targetX = mi.rcWork.right - int32(scaledWidth) - int32(gap)
				targetY = mi.rcWork.bottom - int32(scaledHeight) - int32(gap)
			}

			if targetX < mi.rcWork.left {
				targetX = mi.rcWork.left + int32(gap)
			}
			if targetX+int32(scaledWidth) > mi.rcWork.right {
				targetX = mi.rcWork.right - int32(scaledWidth) - int32(gap)
			}
			if targetY < mi.rcWork.top {
				targetY = mi.rcWork.top + int32(gap)
			}
			if targetY+int32(scaledHeight) > mi.rcWork.bottom {
				targetY = mi.rcWork.bottom - int32(scaledHeight) - int32(gap)
			}

			x = int(targetX)
			y = int(targetY)
		} else {
			success = false
		}
	}

	if !success {
		var cursor point
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&cursor)))

		scaleFactor = a.getScaleFactorAtPoint(cursor.x, cursor.y)

		scaledWidth := int(float64(config.WindowWidth) * scaleFactor)
		scaledHeight := int(float64(config.WindowHeight) * scaleFactor)
		scaledOffsetY := int(float64(config.WindowOffsetYPixels) * scaleFactor)

		x = int(cursor.x) - scaledWidth/2
		y = int(cursor.y) - scaledHeight + scaledOffsetY

		if x < 0 {
			x = 0
		}
	}

	scaledWidth := int(float64(config.WindowWidth) * scaleFactor)
	scaledHeight := int(float64(config.WindowHeight) * scaleFactor)

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
