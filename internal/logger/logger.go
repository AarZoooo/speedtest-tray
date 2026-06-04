package logger

import (
	"log"
	"os"
	"path/filepath"
)

var logFile *os.File

func Init(logFilePath string) error {
	// Ensure the directory exists
	dir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Open the file with O_TRUNC to overwrite it every run
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	// Set the output of the standard logger to the file
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Logger initialized")
	return nil
}
