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
	procGetCursorPos       = user32.NewProc("GetCursorPos")
	procFindWindowW        = user32.NewProc("FindWindowW")
	procCreateRoundRectRgn = gdi32.NewProc("CreateRoundRectRgn")
	procSetWindowRgn       = user32.NewProc("SetWindowRgn")
)

type point struct {
	x int32
	y int32
}

func (a *App) positionWindow() {
	var cursor point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&cursor)))

	x := int(cursor.x) - config.WindowWidth/2
	y := int(cursor.y) - config.WindowHeight + config.WindowOffsetYPixels
	if x < 0 {
		x = 0
	}

	wailsRuntime.WindowSetPosition(a.ctx, x, y)
}

func (a *App) ApplyRoundedCorners() {
	title, _ := syscall.UTF16PtrFromString(config.AppName)
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(title)))
	if hwnd == 0 {
		return
	}

	region, _, _ := procCreateRoundRectRgn.Call(
		0,
		0,
		uintptr(config.WindowWidth),
		uintptr(config.WindowHeight),
		uintptr(config.WindowCornerRadius),
		uintptr(config.WindowCornerRadius),
	)
	procSetWindowRgn.Call(hwnd, region, 1)
}
