package gui

import (
	"runtime"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procGetCursorPos = user32.NewProc("GetCursorPos")
	procFindWindowW  = user32.NewProc("FindWindowW")
	procSetWindowPos = user32.NewProc("SetWindowPos")
)

type POINT struct {
	X, Y int32
}

const (
	HWND_TOPMOST   = ^uintptr(0) // -1
	SWP_NOSIZE     = 0x0001
	SWP_NOZORDER   = 0x0004
	SWP_SHOWWINDOW = 0x0040
)

func (g *GUI) setupTray() {
	if desk, ok := g.App.(desktop.App); ok {
		// Set the main window to be controlled by the tray.
		desk.SetSystemTrayWindow(g.Window)

		// Set the tray icon explicitly
		// This icon is what users hover over.
		desk.SetSystemTrayIcon(theme.InfoIcon())

		// Create the tray menu (for right-click)
		showItem := fyne.NewMenuItem("Show", func() {
			g.ShowNearTray()
		})

		quitItem := fyne.NewMenuItem("Quit", func() {
			g.App.Quit()
		})

		menu := fyne.NewMenu("Speedtest Tray", showItem, quitItem)
		desk.SetSystemTrayMenu(menu)
	}
}

func (g *GUI) ShowNearTray() {
	if runtime.GOOS == "windows" {
		g.positionWindows()
	}
	g.Window.Show()
	g.Window.RequestFocus()
}

func (g *GUI) positionWindows() {
	var pt POINT
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))

	// Find the window by its title
	titlePtr, _ := syscall.UTF16PtrFromString("Speedtest Tray")
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(titlePtr)))

	if hwnd != 0 {
		// Window size: 320x460 (updated with padding)
		w := int32(320)
		h := int32(460)

		x := pt.X - (w / 2)
		y := pt.Y - h - 10

		// Ensure it doesn't go off screen (simple check)
		if x < 0 { x = 0 }

		procSetWindowPos.Call(hwnd, HWND_TOPMOST, uintptr(x), uintptr(y), 0, 0, SWP_NOSIZE|SWP_SHOWWINDOW)
	}
}
