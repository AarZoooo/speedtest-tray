//go:build darwin

package gui_wails

/*
#include <stdlib.h>
void getStatusItemPosition(double *x, double *y, double *width, double *height, double *screenWidth);
*/
import "C"
import (
	"speedtest-tray/internal/config"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) positionWindow() {
	var x, y, width, height, screenWidth C.double
	C.getStatusItemPosition(&x, &y, &width, &height, &screenWidth)

	if width > 0 {
		// Center the window horizontally under the status item
		windowX := int(float64(x) + float64(width)/2 - float64(config.WindowWidth)/2)
		windowY := int(float64(y))

		// Ensure it doesn't go off the right edge of the screen
		margin := 10
		maxRight := int(float64(screenWidth)) - config.WindowWidth - margin
		if windowX > maxRight {
			windowX = maxRight
		}
		if windowX < margin {
			windowX = margin
		}

		wailsRuntime.WindowSetPosition(a.ctx, windowX, windowY)
	}
}

func (a *App) ApplyRoundedCorners() {}
