package gui

import (
	"context"
	"fmt"
	"log"
	"time"

	"speedtest-tray/internal/speedtest_util"

	"fyne.io/fyne/v2"
)

func (g *GUI) startTest() {
	log.Println("Run Test button clicked")
	g.runButton.Disable()
	g.statusLabel.SetText("Initializing...")

	// Reset UI values to loading state immediately
	loadingIcon := spinnerFrames[0]
	g.serverLabel.SetText(loadingIcon)
	g.pingLabel.SetText(loadingIcon)
	g.dlLabel.SetText(loadingIcon)
	g.ulLabel.SetText(loadingIcon)
	g.dlProgress.SetValue(0)
	g.ulProgress.SetValue(0)

	// Start spinner animation
	if g.spinnerStop != nil {
		close(g.spinnerStop)
	}
	g.spinnerStop = make(chan struct{})
	go g.animateSpinner(g.spinnerStop)

	updateCh := make(chan speedtest_util.Update)

	// Fire the speedtest in a background goroutine to keep UI responsive
	resultCh, err := g.tester.RunTest(context.Background(), updateCh)
	if err != nil {
		log.Printf("Failed to start test: %v\n", err)
		g.statusLabel.SetText(fmt.Sprintf("Error: %v", err))
		g.runButton.Enable()
		return
	}

	// Goroutine to consume live metric updates from the channel
	go func() {
		for update := range updateCh {
			// Update individual UI elements on the main lifecycle thread
			u := update // capture for closure
			fyne.Do(func() {
				// Update Status Label using centralized Phase logic
				g.statusLabel.SetText(u.Phase.String())

				// Server info resolution
				if u.Server != "" {
					g.serverLabel.SetText(u.Server)
				}

				if u.Ping > 0 {
					g.pingLabel.SetText(fmt.Sprintf("%.0f ms", u.Ping))
				}
				if u.Download > 0 {
					g.dlLabel.SetText(fmt.Sprintf("%.1f Mbps", u.Download))
					g.dlProgress.SetValue(u.Download)
				}
				if u.Upload > 0 {
					g.ulLabel.SetText(fmt.Sprintf("%.1f Mbps", u.Upload))
					g.ulProgress.SetValue(u.Upload)
				}
			})
		}

		log.Println("Update channel closed, waiting for final result...")
		// Handle final completion state from the result channel
		result, ok := <-resultCh
		if !ok {
			log.Println("Result channel closed unexpectedly")
			fyne.Do(func() {
				if g.spinnerStop != nil {
					close(g.spinnerStop)
					g.spinnerStop = nil
				}
				g.runButton.Enable()
				g.runButton.SetText("Start Again")
				g.statusLabel.SetText("Error")
			})
			return
		}

		log.Printf("Final result received: Error=%v, Server=%s\n", result.Error, result.Server)
		fyne.Do(func() {
			if g.spinnerStop != nil {
				close(g.spinnerStop)
				g.spinnerStop = nil
			}

			if result.Error != nil {
				g.statusLabel.SetText(speedtest_util.FAILED.String())
				g.runButton.SetText("Try Again")
			} else {
				g.statusLabel.SetText(speedtest_util.COMPLETED.String())
				if result.Server != "" {
					g.serverLabel.SetText(result.Server)
				}
				g.runButton.SetText("Start Again")
			}
			g.runButton.Enable()
		})
	}()
}

func (g *GUI) animateSpinner(stop chan struct{}) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	i := 0
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			i = (i + 1) % len(spinnerFrames)
			frame := spinnerFrames[i]
			fyne.Do(func() {
				// Only update if they are still in "loading" state
				if isSpinner(g.serverLabel.Text) {
					g.serverLabel.SetText(frame)
				}
				if isSpinner(g.pingLabel.Text) {
					g.pingLabel.SetText(frame)
				}
				if isSpinner(g.dlLabel.Text) {
					g.dlLabel.SetText(frame)
				}
				if isSpinner(g.ulLabel.Text) {
					g.ulLabel.SetText(frame)
				}
			})
		}
	}
}
