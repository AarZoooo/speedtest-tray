package main

import (
	"log"
	"speedtest-tray/internal/gui"
	"speedtest-tray/internal/speedtest_util"
)

func main() {
	log.Println("Starting Speedtest Utility...")

	// 1. Initialize the backend speedtester instance
	st := speedtest_util.New()

	// 2. Initialize the GUI layout framework with the backend engine passed in
	ui := gui.New(st)

	// 3. Start the application lifecycle.
	// The window stays hidden until the tray icon is clicked.
	ui.App.Run()
}
