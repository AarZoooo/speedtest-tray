const PHASES = {
  INITIALIZING: "INITIALIZING",
  GETTING_INFO: "GETTING_INFO",
  FINDING_SERVERS: "FINDING_SERVERS",
  SELECTING_SERVER: "SELECTING_SERVER",
  SERVER_SELECTED: "SERVER_SELECTED",
  PING_TEST: "PING_TEST",
  STARTING_DOWNLOAD: "STARTING_DOWNLOAD",
  DOWNLOADING: "DOWNLOADING",
  STARTING_UPLOAD: "STARTING_UPLOAD",
  UPLOADING: "UPLOADING",
  COMPLETED: "COMPLETED",
  FAILED: "FAILED",
};

const EVENTS = {
  WINDOW_SHOWN: "window_shown",
  TEST_UPDATE: "test_update",
  TEST_COMPLETE: "test_complete",
  TEST_ERROR: "test_error",
};

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

let canHide = false;
let isTesting = false;

const elements = {
  server: document.getElementById("server"),
  ping: document.getElementById("ping"),
  download: document.getElementById("download"),
  upload: document.getElementById("upload"),
  status: document.getElementById("status"),
  runBtn: document.getElementById("run-btn"),
  speedometer: document.getElementById("speedometer"),
};

function updateGauge(speed) {
  if (elements.speedometer) {
    elements.speedometer.setValue(speed);
  }
}

window.runtime.EventsOn(EVENTS.WINDOW_SHOWN, () => {
  canHide = false;
  setTimeout(() => {
    canHide = true;
  }, 1000);
});

window.onblur = () => {
  if (canHide) {
    window.go.gui_wails.App.HideWindow();
    canHide = false;
  }
};

document.addEventListener("visibilitychange", () => {
  if (document.visibilityState === "hidden" && canHide) {
    window.go.gui_wails.App.HideWindow();
    canHide = false;
  }
});

function handleBtnClick() {
  if (isTesting) {
    stopTest();
  } else {
    startTest();
  }
}

function startTest() {
  console.log("JS: startTest called");
  if (!elements.runBtn) return;

  isTesting = true;
  elements.runBtn.innerText = TEXT.STOP;
  elements.status.innerText = TEXT.INITIALIZING;

  resetUI(TEXT.LOADER_HTML);

  console.log("JS: Invoking backend StartTest");
  window.go.gui_wails.App.StartTest()
    .then(() => console.log("JS: Backend StartTest promise resolved"))
    .catch((err) => {
      console.error("JS: Backend StartTest failed:", err);
      isTesting = false;
      elements.runBtn.innerText = TEXT.TRY_AGAIN;
    });
}

function stopTest() {
  console.log("JS: stopTest called");
  isTesting = false; // Set to false immediately to ignore incoming updates
  elements.runBtn.disabled = true; // Briefly disable while backend cleans up
  elements.status.innerText = TEXT.TEST_STOPPED;
  resetUI(TEXT.DEFAULT_VAL);

  window.go.gui_wails.App.StopTest();
}

function resetUI(val) {
  elements.server.innerHTML = val;
  elements.ping.innerHTML = val;
  elements.download.innerHTML = val;
  elements.upload.innerHTML = val;
  updateGauge(0);
}

window.runtime.EventsOn(EVENTS.TEST_UPDATE, (data) => {
  if (!isTesting) return; // Ignore updates if we've stopped

  elements.status.innerText = data.status;

  if (data.server) elements.server.innerText = data.server;

  if (parseFloat(data.ping) > 0) {
    elements.ping.innerText = data.ping + TEXT.MS_SUFFIX;
  }

  if (data.phase === PHASES.DOWNLOADING) {
    if (parseFloat(data.download) > 0) {
      elements.speedometer.setMax(1000);
      elements.download.innerText = data.download + TEXT.MBPS_SUFFIX;
      updateGauge(data.download);
    }
  } else if (data.phase === PHASES.UPLOADING) {
    if (parseFloat(data.upload) > 0) {
      elements.speedometer.setMax(100);
      elements.upload.innerText = data.upload + TEXT.MBPS_SUFFIX;
      updateGauge(data.upload);
    }
  } else {
    updateGauge(0);
  }
});

window.runtime.EventsOn(EVENTS.TEST_COMPLETE, (data) => {
  // If we stopped manually, isTesting is already false
  const wasManualStop =
    !isTesting && elements.status.innerText === TEXT.TEST_STOPPED;

  isTesting = false;
  elements.runBtn.disabled = false;
  elements.runBtn.innerText =
    data.error || wasManualStop ? TEXT.TRY_AGAIN : TEXT.START_AGAIN;

  if (wasManualStop || data.error === "Test stopped") {
    elements.status.innerText = TEXT.TEST_STOPPED;
    resetUI(TEXT.DEFAULT_VAL);
  } else {
    elements.status.innerText = data.error
      ? TEXT.TEST_FAILED
      : TEXT.TEST_COMPLETED;

    if (!data.error) {
      elements.server.innerText = data.server;
      elements.ping.innerText = data.ping + TEXT.MS_SUFFIX;
      elements.download.innerText = data.download + TEXT.MBPS_SUFFIX;
      elements.upload.innerText = data.upload + TEXT.MBPS_SUFFIX;
    } else {
      resetUI(TEXT.DEFAULT_VAL);
    }
  }

  updateGauge(0);
});

window.runtime.EventsOn(EVENTS.TEST_ERROR, (err) => {
  isTesting = false;
  elements.runBtn.innerText = TEXT.TRY_AGAIN;
  elements.status.innerText = TEXT.ERROR_PREFIX + err;
  updateGauge(0);
  resetUI(TEXT.DEFAULT_VAL);
});
