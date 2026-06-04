package main

import (
	"log"
	"speedtest-tray/internal/gui"
	"speedtest-tray/internal/speedtest_util"
)

func main() {
	// Initialize logger
	// if err := logger.Init("logs/speedtest.log"); err != nil {
	// 	log.Fatalf("Failed to initialize logger: %v", err)
	// }
	log.Println("Starting Speedtest Utility...")

	// 1. Initialize the backend speedtester instance
	st := speedtest_util.New()

	// 2. Initialize the GUI layout framework with the backend engine passed in
	ui := gui.New(st)

	// 3. Launch the window panel and block execution on the main OS thread layout pump
	ui.Window.ShowAndRun()
}
