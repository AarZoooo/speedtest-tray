package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"speedtest-tray/internal/config"
	"speedtest-tray/internal/speedtest_util"
)

type jsonOutput struct {
	Status       string    `json:"status"`
	PingMS       float64   `json:"ping_ms,omitempty"`
	DownloadMbps float64   `json:"download_mbps,omitempty"`
	UploadMbps   float64   `json:"upload_mbps,omitempty"`
	Server       string    `json:"server,omitempty"`
	Error        string    `json:"error,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

func Run(ctx context.Context, jsonMode bool, serverID string, orchestrator speedtest_util.TestOrchestrator, stdout io.Writer) error {
	if !jsonMode {
		fmt.Fprintln(stdout, config.CLIHeader)
		fmt.Fprintln(stdout, config.CLILineSeparator)
	}

	updateCh := make(chan speedtest_util.Update, config.UpdateChannelSize)
	runner := speedtest_util.NewTestRunner(orchestrator)
	resultCh, err := runner.RunTest(ctx, updateCh)
	if err != nil {
		if jsonMode {
			printJSONError(stdout, err)
		} else {
			fmt.Fprintf(stdout, "\nError: %v\n", err)
		}
		return err
	}

	var printedInit, printedSelect, printedPing, printedDownload, printedUpload bool
	var downloadDone, uploadDone, pingDone bool
	var lastPing, lastDownload, lastUpload float64
	var selectedServer string

	for update := range updateCh {
		if update.Error != nil {
			if jsonMode {
				printJSONError(stdout, update.Error)
			} else {
				fmt.Fprintf(stdout, "\nError: %v\n", update.Error)
			}
			return update.Error
		}

		if jsonMode {
			if update.Server != "" {
				selectedServer = update.Server
			}
			if update.Ping > 0 {
				lastPing = update.Ping
			}
			if update.Download > 0 {
				lastDownload = update.Download
			}
			if update.Upload > 0 {
				lastUpload = update.Upload
			}
			continue
		}

		switch update.Phase {
		case speedtest_util.INITIALIZING:
			if !printedInit {
				fmt.Fprint(stdout, "[1/5] "+config.CLIStatusChecking+" ")
				printedInit = true
			}
		case speedtest_util.GETTING_INFO, speedtest_util.FINDING_SERVERS, speedtest_util.SELECTING_SERVER:
			if !printedInit {
				fmt.Fprint(stdout, "[1/5] "+config.CLIStatusChecking+" ")
				printedInit = true
			}
		case speedtest_util.SERVER_SELECTED:
			if printedInit && !printedSelect {
				fmt.Fprintln(stdout, config.CLIStatusCompleted)
				fmt.Fprintf(stdout, "[2/5] "+config.CLIStatusSelecting+" %s\n", update.Server)
				selectedServer = update.Server
				printedSelect = true
			}
		case speedtest_util.PING_TEST:
			if !printedPing {
				fmt.Fprint(stdout, "[3/5] "+config.CLIStatusPing+" ")
				printedPing = true
			}
			if update.Ping > 0 && !pingDone {
				fmt.Fprintf(stdout, "%s ms\n", speedtest_util.FormatNumber(update.Ping, 0))
				lastPing = update.Ping
				pingDone = true
			}
		case speedtest_util.STARTING_DOWNLOAD:
			if !pingDone && printedPing {
				fmt.Fprintln(stdout, "")
				pingDone = true
			}
		case speedtest_util.DOWNLOADING:
			if !printedDownload {
				fmt.Fprint(stdout, "[4/5] "+config.CLIStatusDownload+" ")
				printedDownload = true
			}
			phaseProgress := (update.Progress - config.ProgressDownStart) / (config.ProgressDownEnd - config.ProgressDownStart)
			if phaseProgress < 0 {
				phaseProgress = 0
			}
			if phaseProgress > 1 {
				phaseProgress = 1
			}
			bar := drawProgressBar(phaseProgress)
			lastDownload = update.Download
			fmt.Fprintf(stdout, "\r[4/5] %s %s %s Mbps", config.CLIStatusDownload, bar, speedtest_util.FormatNumber(update.Download, 2))
		case speedtest_util.STARTING_UPLOAD:
			if printedDownload && !downloadDone {
				speed := update.Download
				if speed == 0 {
					speed = lastDownload
				}
				bar := drawProgressBar(1.0)
				fmt.Fprintf(stdout, "\r[4/5] %s %s %s Mbps", config.CLIStatusDownload, bar, speedtest_util.FormatNumber(speed, 2))
				fmt.Fprintln(stdout, " "+config.CLIStatusCompleted)
				downloadDone = true
			}
		case speedtest_util.UPLOADING:
			if !printedUpload {
				fmt.Fprint(stdout, "[5/5] "+config.CLIStatusUpload+" ")
				printedUpload = true
			}
			phaseProgress := (update.Progress - config.ProgressUpStart) / (config.ProgressUpEnd - config.ProgressUpStart)
			if phaseProgress < 0 {
				phaseProgress = 0
			}
			if phaseProgress > 1 {
				phaseProgress = 1
			}
			bar := drawProgressBar(phaseProgress)
			lastUpload = update.Upload
			fmt.Fprintf(stdout, "\r[5/5] %s %s %s Mbps", config.CLIStatusUpload, bar, speedtest_util.FormatNumber(update.Upload, 2))
		case speedtest_util.COMPLETED:
			if printedUpload && !uploadDone {
				speed := update.Upload
				if speed == 0 {
					speed = lastUpload
				}
				bar := drawProgressBar(1.0)
				fmt.Fprintf(stdout, "\r[5/5] %s %s %s Mbps", config.CLIStatusUpload, bar, speedtest_util.FormatNumber(speed, 2))
				fmt.Fprintln(stdout, " "+config.CLIStatusCompleted)
				uploadDone = true
			}
		}
	}

	if !jsonMode && printedUpload && !uploadDone {
		bar := drawProgressBar(1.0)
		fmt.Fprintf(stdout, "\r[5/5] %s %s %s Mbps", config.CLIStatusUpload, bar, speedtest_util.FormatNumber(lastUpload, 2))
		fmt.Fprintln(stdout, " "+config.CLIStatusCompleted)
		uploadDone = true
	}

	var result speedtest_util.Result
	select {
	case res := <-resultCh:
		result = res
	case <-time.After(config.ResultTimeout):
		errTimeout := fmt.Errorf(config.ErrTestTimeout)
		if jsonMode {
			printJSONError(stdout, errTimeout)
		} else {
			fmt.Fprintf(stdout, "\nError: %v\n", errTimeout)
		}
		return errTimeout
	}

	if result.Error != nil {
		if jsonMode {
			printJSONError(stdout, result.Error)
		} else {
			fmt.Fprintf(stdout, "\nError: %v\n", result.Error)
		}
		return result.Error
	}

	if jsonMode {
		output := jsonOutput{
			Status:       config.JSONStatusSuccess,
			PingMS:       result.Ping,
			DownloadMbps: result.Download,
			UploadMbps:   result.Upload,
			Server:       result.Server,
			Timestamp:    time.Now().UTC(),
		}
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(output)
	}

	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, config.CLIDoubleLine)
	fmt.Fprintln(stdout, config.CLIResultHeader)
	fmt.Fprintln(stdout, config.CLIDoubleLine)
	fmt.Fprintf(stdout, "  Server:   %s\n", result.Server)
	fmt.Fprintf(stdout, "  Ping:     %s ms\n", speedtest_util.FormatNumber(result.Ping, 2))
	fmt.Fprintf(stdout, "  Download: %s Mbps\n", speedtest_util.FormatNumber(result.Download, 2))
	fmt.Fprintf(stdout, "  Upload:   %s Mbps\n", speedtest_util.FormatNumber(result.Upload, 2))
	fmt.Fprintln(stdout, config.CLIDoubleLine)

	_ = lastPing
	_ = lastDownload
	_ = lastUpload
	_ = selectedServer

	return nil
}

func drawProgressBar(progress float64) string {
	const barWidth = 20
	filled := int(progress * float64(barWidth))
	if filled < 0 {
		filled = 0
	}
	if filled > barWidth {
		filled = barWidth
	}
	bar := make([]byte, barWidth)
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar[i] = '='
		} else if i == filled && filled > 0 && filled < barWidth {
			bar[i] = '>'
		} else {
			bar[i] = ' '
		}
	}
	return fmt.Sprintf("[%s] %3.0f%%", string(bar), progress*100)
}

func printJSONError(stdout io.Writer, err error) {
	output := jsonOutput{
		Status:    config.JSONStatusFailed,
		Error:     err.Error(),
		Timestamp: time.Now().UTC(),
	}
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(output)
}
