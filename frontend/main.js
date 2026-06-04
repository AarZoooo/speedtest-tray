const spinnerFrames = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"];
let spinnerInterval = null;
let canHide = false;

const elements = {
  server: document.getElementById("server"),
  ping: document.getElementById("ping"),
  download: document.getElementById("download"),
  upload: document.getElementById("upload"),
  dlProgress: document.getElementById("dl-progress"),
  ulProgress: document.getElementById("ul-progress"),
  status: document.getElementById("status"),
  runBtn: document.getElementById("run-btn"),
};

window.runtime.EventsOn("window_shown", () => {
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

function startTest() {
  elements.runBtn.disabled = true;
  elements.status.innerText = "Initializing...";

  resetUI(spinnerFrames[0]);
  startSpinner();

  window.go.gui_wails.App.StartTest();
}

function resetUI(val) {
  elements.server.innerText = val;
  elements.ping.innerText = val;
  elements.download.innerText = val;
  elements.upload.innerText = val;
  elements.dlProgress.style.width = "0%";
  elements.ulProgress.style.width = "0%";
}

function startSpinner() {
  if (spinnerInterval) clearInterval(spinnerInterval);
  let i = 0;
  spinnerInterval = setInterval(() => {
    const frame = spinnerFrames[i];
    if (isSpinner(elements.server.innerText)) elements.server.innerText = frame;
    if (isSpinner(elements.ping.innerText)) elements.ping.innerText = frame;
    if (isSpinner(elements.download.innerText))
      elements.download.innerText = frame;
    if (isSpinner(elements.upload.innerText)) elements.upload.innerText = frame;
    i = (i + 1) % spinnerFrames.length;
  }, 100);
}

function stopSpinner() {
  if (spinnerInterval) {
    clearInterval(spinnerInterval);
    spinnerInterval = null;
  }
}

function isSpinner(text) {
  return spinnerFrames.includes(text) || text === "--";
}

window.runtime.EventsOn("test_update", (data) => {
  elements.status.innerText = data.status;

  if (data.server) elements.server.innerText = data.server;

  if (parseFloat(data.ping) > 0) {
    elements.ping.innerText = data.ping + " ms";
  }

  if (parseFloat(data.download) > 0) {
    elements.download.innerText = data.download + " Mbps";
    let w = (parseFloat(data.download) / 1000) * 100;
    elements.dlProgress.style.width = Math.min(w, 100) + "%";
  }

  if (parseFloat(data.upload) > 0) {
    elements.upload.innerText = data.upload + " Mbps";
    let w = (parseFloat(data.upload) / 1000) * 100;
    elements.ulProgress.style.width = Math.min(w, 100) + "%";
  }
});

window.runtime.EventsOn("test_complete", (data) => {
  stopSpinner();
  elements.runBtn.disabled = false;
  elements.runBtn.innerText = data.error ? "Try Again" : "Start Again";
  elements.status.innerText = data.error ? "Test Failed" : "Test Completed";

  if (!data.error) {
    elements.server.innerText = data.server;
    elements.ping.innerText = data.ping + " ms";
    elements.download.innerText = data.download + " Mbps";
    elements.upload.innerText = data.upload + " Mbps";
  }
});

window.runtime.EventsOn("test_error", (err) => {
  stopSpinner();
  elements.runBtn.disabled = false;
  elements.runBtn.innerText = "Try Again";
  elements.status.innerText = "Error: " + err;
});
