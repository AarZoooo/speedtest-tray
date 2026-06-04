import { PHASES, CONFIG } from "./constants.js";
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

  if (wasManualStop || data.error === "Test stopped") {
    setStatus(TEXT.TEST_STOPPED);
    resetUI(TEXT.DEFAULT_VAL);
  } else {
    setStatus(data.error ? TEXT.TEST_FAILED : TEXT.TEST_COMPLETED);

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

// Handle test error
export function handleTestError(err) {
  testState.stopTest();
  if (elements.runBtn) elements.runBtn.innerText = TEXT.TRY_AGAIN;
  setStatus(TEXT.ERROR_PREFIX + err);
  updateGauge(0);
  resetUI(TEXT.DEFAULT_VAL);
}
