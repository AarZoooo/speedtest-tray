import { PHASES, CONFIG, ERRORS, MESSAGES } from "./constants.js";
import { testState } from "./state.js";

const TEXT = {
  INITIALIZING: "Initializing...",
  TRY_AGAIN: "Try Again",
  START_AGAIN: "Start Again",
  STOP: "Stop",
  TEST_FAILED: "Test Failed",
  TEST_STOPPED: "Test Stopped",
  TEST_COMPLETED: "Test Completed",
  ERROR_PREFIX: "Error: ",
  MS_SUFFIX: " ms",
  MBPS_SUFFIX: " Mbps",
  DEFAULT_VAL: "--",
  LOADER_HTML: '<div class="loader"></div>',
};

// UI element cache
const elements = {
  server: null,
  ping: null,
  download: null,
  upload: null,
  status: null,
  runBtn: null,
  speedometer: null,
  historyToggleBtn: null,
  clearHistoryBtn: null,
  openJsonBtn: null,
  testView: null,
  historyView: null,
  historyList: null,
};

// Initialize element references
export function initializeElements() {
  elements.server = document.getElementById("server");
  elements.ping = document.getElementById("ping");
  elements.download = document.getElementById("download");
  elements.upload = document.getElementById("upload");
  elements.status = document.getElementById("status");
  elements.runBtn = document.getElementById("run-btn");
  elements.speedometer = document.getElementById("speedometer");
  elements.historyToggleBtn = document.getElementById("history-toggle-btn");
  elements.clearHistoryBtn = document.getElementById("clear-history-btn");
  elements.openJsonBtn = document.getElementById("open-json-btn");
  elements.testView = document.getElementById("test-view");
  elements.historyView = document.getElementById("history-view");
  elements.historyList = document.getElementById("history-list");
  return elements;
}

// Update UI with test result
export function updateResults(server, ping, download, upload) {
  if (elements.server) elements.server.innerText = server;
  if (elements.ping) elements.ping.innerText = ping;
  if (elements.download) elements.download.innerText = download;
  if (elements.upload) elements.upload.innerText = upload;
}

// Update gauge display
export function updateGauge(speed) {
  if (elements.speedometer) {
    elements.speedometer.setValue(speed);
  }
}

// Reset UI to default values
export function resetUI(val = TEXT.DEFAULT_VAL) {
  if (elements.server) elements.server.innerHTML = val;
  if (elements.ping) elements.ping.innerHTML = val;
  if (elements.download) elements.download.innerHTML = val;
  if (elements.upload) elements.upload.innerHTML = val;
  updateGauge(0);
}

// Update status text
export function setStatus(text) {
  if (elements.status) elements.status.innerText = text;
}

// Update button state
export function setButtonState(isTesting) {
  if (elements.runBtn) {
    elements.runBtn.innerText = isTesting ? TEXT.STOP : TEXT.TRY_AGAIN;
    elements.runBtn.disabled = false;
  }
  updateHistoryToggleState(isTesting);
}

// Handle test update event
export function handleTestUpdate(data) {
  if (!testState.isTesting) return; // Ignore if we stopped

  setStatus(data.status);

  if (data.server) {
    if (elements.server) elements.server.innerText = data.server;
  }

  if (parseFloat(data.ping) > 0) {
    if (elements.ping) elements.ping.innerText = data.ping + TEXT.MS_SUFFIX;
  }

  if (data.phase === PHASES.DOWNLOADING) {
    if (parseFloat(data.download) > 0) {
      elements.speedometer.setMax(CONFIG.GAUGE_MAX_DOWNLOAD);
      if (elements.download) elements.download.innerText = data.download + TEXT.MBPS_SUFFIX;
      updateGauge(data.download);
    }
  } else if (data.phase === PHASES.UPLOADING) {
    if (parseFloat(data.upload) > 0) {
      elements.speedometer.setMax(CONFIG.GAUGE_MAX_UPLOAD);
      if (elements.upload) elements.upload.innerText = data.upload + TEXT.MBPS_SUFFIX;
      updateGauge(data.upload);
    }
  } else {
    updateGauge(0);
  }
}

// Handle test completion
export function handleTestComplete(data) {
  const wasManualStop = !testState.isTesting && elements.status.innerText === TEXT.TEST_STOPPED;

  testState.stopTest();
  if (elements.runBtn) {
    elements.runBtn.disabled = false;
    elements.runBtn.innerText = data.error || wasManualStop ? TEXT.TRY_AGAIN : TEXT.START_AGAIN;
  }

  updateHistoryToggleState(false);

  if (wasManualStop || data.error?.toLowerCase() === "test stopped") {
    setStatus(TEXT.TEST_STOPPED);
    resetUI(TEXT.DEFAULT_VAL);
  } else {
    setStatus(data.error ? failureStatus(data.error) : TEXT.TEST_COMPLETED);

    if (!data.error) {
      if (elements.server) elements.server.innerText = data.server;
      if (elements.ping) elements.ping.innerText = data.ping + TEXT.MS_SUFFIX;
      if (elements.download) elements.download.innerText = data.download + TEXT.MBPS_SUFFIX;
      if (elements.upload) elements.upload.innerText = data.upload + TEXT.MBPS_SUFFIX;
    } else {
      resetUI(TEXT.DEFAULT_VAL);
    }
  }

  updateGauge(0);
}

function failureStatus(error) {
  if (error?.toLowerCase() === ERRORS.NO_INTERNET) {
    return MESSAGES.NO_INTERNET;
  }
  return TEXT.TEST_FAILED;
}

// Handle test error
export function handleTestError(err) {
  testState.stopTest();
  if (elements.runBtn) elements.runBtn.innerText = TEXT.TRY_AGAIN;
  setStatus(TEXT.ERROR_PREFIX + err);
  updateGauge(0);
  resetUI(TEXT.DEFAULT_VAL);
  updateHistoryToggleState(false);
}

export function updateHistoryToggleState(isTesting) {
  if (elements.historyToggleBtn) {
    elements.historyToggleBtn.disabled = isTesting;
    if (isTesting) {
      elements.historyToggleBtn.style.opacity = "0.5";
      elements.historyToggleBtn.style.cursor = "not-allowed";
    } else {
      elements.historyToggleBtn.style.opacity = "1";
      elements.historyToggleBtn.style.cursor = "pointer";
    }
  }
}

export function renderHistory(history) {
  if (!elements.historyList) return;

  elements.historyList.innerHTML = "";

  if (!history || history.length === 0) {
    elements.historyList.innerHTML = `
      <div class="empty-history">
        <div class="empty-icon">🕒</div>
        <div>No test history yet</div>
      </div>
    `;
    if (elements.clearHistoryBtn) {
      elements.clearHistoryBtn.disabled = true;
    }
    return;
  }

  if (elements.clearHistoryBtn) {
    elements.clearHistoryBtn.disabled = false;
  }

  history.forEach(entry => {
    const card = document.createElement("div");
    card.className = "history-card";

    let dateStr = "";
    try {
      const date = new Date(entry.timestamp);
      dateStr = date.toLocaleString(undefined, {
        month: "short",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit"
      });
    } catch (e) {
      dateStr = entry.timestamp;
    }

    card.innerHTML = `
      <div class="history-card-header">
        <span class="history-server">${entry.server || "Unknown Server"}</span>
        <span class="history-date">${dateStr}</span>
      </div>
      <div class="history-metrics">
        <div class="history-metric">
          <span class="metric-icon">↓</span>
          <span class="metric-value">${parseFloat(entry.download).toFixed(1)}</span>
          <span class="metric-unit">Mbps</span>
        </div>
        <div class="history-metric">
          <span class="metric-icon">↑</span>
          <span class="metric-value">${parseFloat(entry.upload).toFixed(1)}</span>
          <span class="metric-unit">Mbps</span>
        </div>
        <div class="history-metric">
          <span class="metric-icon">⇄</span>
          <span class="metric-value">${Math.round(entry.ping)}</span>
          <span class="metric-unit">ms</span>
        </div>
      </div>
    `;
    elements.historyList.appendChild(card);
  });
}
