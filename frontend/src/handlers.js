import { testState } from "./state.js";
import { setStatus, setButtonState, resetUI, renderHistory } from "./ui.js";
import { CONFIG } from "./constants.js";

const TEXT = {
  INITIALIZING: "Initializing...",
  STOP: "Stop",
  LOADER_HTML: '<div class="loader"></div>',
  DEFAULT_VAL: "--",
};

// Start a new test
export async function startTest() {
  console.log("JS: startTest called");

  testState.startTest();
  setButtonState(true);
  setStatus(TEXT.INITIALIZING);
  resetUI(TEXT.LOADER_HTML);

  // Play the startup sweep animation — await it so the backend only starts after the animation finishes
  const speedometer = document.getElementById("speedometer");
  if (speedometer?.playStartupSweep) {
    await speedometer.playStartupSweep();
  }

  // Only proceed if the test wasn't stopped during the animation
  if (!testState.isTesting) return;

  const header = document.querySelector("header");
  if (header) {
    header.classList.add("loading");
  }

  console.log("JS: Invoking backend StartTest");
  try {
    await window.go.gui_wails.App.StartTest();
    console.log("JS: Backend StartTest promise resolved");
  } catch (err) {
    console.error("JS: Backend StartTest failed:", err);
    testState.stopTest();
    setButtonState(false);
    if (header) {
      header.classList.remove("loading");
    }
  }
}

// Stop the running test
export function stopTest() {
  console.log("JS: stopTest called");
  testState.stopTest();

  // Stop sweep animation if running
  const speedometer = document.getElementById("speedometer");
  if (speedometer?.sweeping) speedometer.stopSweep();

  const btn = document.getElementById("run-btn");
  if (btn) {
    btn.disabled = true;
    btn.classList.remove("running");
  }
  setStatus("Test Stopped");
  resetUI(TEXT.DEFAULT_VAL);

  const header = document.querySelector("header");
  if (header) {
    header.classList.remove("loading");
  }

  window.go.gui_wails.App.StopTest();
}

// Handle button click - start or stop
export function handleButtonClick() {
  if (testState.isTesting) {
    stopTest();
  } else {
    startTest();
  }
}

// Initialize button event listener
export function initializeButtonHandler() {
  const btn = document.getElementById("run-btn");
  if (btn) {
    btn.addEventListener("click", handleButtonClick);
  }
}

const SPEED_ICON_HTML = `<svg class="icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M6.34315 17.6569C5.22433 16.538 4.4624 15.1126 4.15372 13.5607C3.84504 12.0089 4.00346 10.4003 4.60896 8.93853C5.21446 7.47672 6.23984 6.22729 7.55544 5.34824C8.87103 4.46919 10.4177 4 12 4C13.5823 4 15.129 4.46919 16.4446 5.34824C17.7602 6.22729 18.7855 7.47672 19.391 8.93853C19.9965 10.4003 20.155 12.0089 19.8463 13.5607C19.5376 15.1126 18.7757 16.538 17.6569 17.6569" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
<path d="M12 12L16 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
</svg>`;

const HISTORY_ICON_HTML = `<svg class="icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
<path fill-rule="evenodd" clip-rule="evenodd" d="M9.09958 2.39754C9.24874 2.78396 9.05641 3.21814 8.66999 3.36731C8.52855 3.42191 8.38879 3.47988 8.2508 3.54114C7.87221 3.70921 7.42906 3.53856 7.261 3.15997C7.09293 2.78139 7.26358 2.33824 7.64217 2.17017C7.80267 2.09892 7.96526 2.03147 8.12981 1.96795C8.51623 1.81878 8.95041 2.01112 9.09958 2.39754ZM5.6477 4.24026C5.93337 4.54021 5.92178 5.01495 5.62183 5.30061C5.51216 5.40506 5.40505 5.51216 5.30061 5.62183C5.01495 5.92178 4.54021 5.93337 4.24026 5.6477C3.94031 5.36204 3.92873 4.88731 4.21439 4.58736C4.33566 4.46003 4.46002 4.33566 4.58736 4.21439C4.88731 3.92873 5.36204 3.94031 5.6477 4.24026ZM3.15997 7.261C3.53856 7.42907 3.70921 7.87221 3.54114 8.2508C3.47988 8.38879 3.42191 8.52855 3.36731 8.66999C3.21814 9.05641 2.78396 9.24874 2.39754 9.09958C2.01112 8.95041 1.81878 8.51623 1.96795 8.12981C2.03147 7.96526 2.09892 7.80267 2.17017 7.64217C2.33824 7.26358 2.78139 7.09293 3.15997 7.261ZM2.02109 11.004C2.43518 11.0141 2.76276 11.3579 2.75275 11.7719C2.75092 11.8477 2.75 11.9237 2.75 12C2.75 12.0763 2.75092 12.1523 2.75275 12.2281C2.76276 12.6421 2.43518 12.9859 2.02109 12.996C1.60699 13.006 1.26319 12.6784 1.25319 12.2643C1.25107 12.1764 1.25 12.0883 1.25 12C1.25 11.9117 1.25107 11.8236 1.25319 11.7357C1.26319 11.3216 1.60699 10.994 2.02109 11.004ZM21.6025 14.9004C21.9889 15.0496 22.1812 15.4838 22.032 15.8702C21.9685 16.0347 21.9011 16.1973 21.8298 16.3578C21.6618 16.7364 21.2186 16.9071 20.84 16.739C20.4614 16.5709 20.2908 16.1278 20.4589 15.7492C20.5201 15.6112 20.5781 15.4714 20.6327 15.33C20.7819 14.9436 21.216 14.7513 21.6025 14.9004ZM2.39754 14.9004C2.78396 14.7513 3.21814 14.9436 3.36731 15.33C3.42191 15.4714 3.47988 15.6112 3.54114 15.7492C3.70921 16.1278 3.53856 16.5709 3.15997 16.739C2.78139 16.9071 2.33824 16.7364 2.17017 16.3578C2.09892 16.1973 2.03147 16.0347 1.96795 15.8702C1.81878 15.4838 2.01112 15.0496 2.39754 14.9004ZM19.7597 18.3523C20.0597 18.638 20.0713 19.1127 19.7856 19.4126C19.6643 19.54 19.54 19.6643 19.4126 19.7856C19.1127 20.0713 18.638 20.0597 18.3523 19.7597C18.0666 19.4598 18.0782 18.9851 18.3782 18.6994C18.4878 18.5949 18.5949 18.4878 18.6994 18.3782C18.9851 18.0782 19.4598 18.0666 19.7597 18.3523ZM4.24026 18.3523C4.54021 18.0666 5.01495 18.0782 5.30061 18.3782C5.40506 18.4878 5.51216 18.5949 5.62183 18.6994C5.92178 18.9851 5.93337 19.4598 5.6477 19.7597C5.36204 20.0597 4.88731 20.0713 4.58736 19.7856C4.46003 19.6643 4.33566 19.54 4.21439 19.4126C3.92873 19.1127 3.94031 18.638 4.24026 18.3523ZM7.261 20.84C7.42907 20.4614 7.87221 20.2908 8.2508 20.4589C8.38879 20.5201 8.52855 20.5781 8.66999 20.6327C9.05641 20.7819 9.24874 21.216 9.09958 21.6025C8.95041 21.9889 8.51623 22.1812 8.12981 22.032C7.96526 21.9685 7.80267 21.9011 7.64217 21.8298C7.26358 21.6618 7.09293 21.2186 7.261 20.84ZM16.739 20.84C16.9071 21.2186 16.7364 21.6618 16.3578 21.8298C16.1973 21.9011 16.0347 21.9685 15.8702 22.032C15.4838 22.1812 15.0496 21.9889 14.9004 21.6025C14.7513 21.216 14.9436 20.7819 15.33 20.6327C15.4714 20.5781 15.6112 20.5201 15.7492 20.4589C16.1278 20.2908 16.5709 20.4614 16.739 20.84ZM11.004 21.9789C11.0141 21.5648 11.3579 21.2372 11.7719 21.2472C11.8477 21.2491 11.9237 21.25 12 21.25C12.0763 21.25 12.1523 21.2491 12.2281 21.2472C12.6421 21.2372 12.9859 21.5648 12.996 21.9789C13.006 22.393 12.6784 22.7368 12.2643 22.7468C12.1764 22.7489 12.0883 22.75 12 22.75C11.9117 22.75 11.8236 22.7489 11.7357 22.7468C11.3216 22.7468 10.994 22.393 11.004 21.9789ZM12 2.75C17.1086 2.75 21.25 6.89137 21.25 12C21.25 12.4142 21.5858 12.75 22 12.75C22.4142 12.75 22.75 12.4142 22.75 12C22.75 6.06294 17.9371 1.25 12 1.25C11.5858 1.25 11.25 1.58579 11.25 2C11.25 2.41421 11.5858 2.75 12 2.75ZM12 8.25C12.4142 8.25 12.75 8.58579 12.75 9V12.25H16C16.4142 12.25 16.75 12.5858 16.75 13C16.75 13.4142 16.4142 13.75 16 13.75H12C11.5858 13.75 11.25 13.4142 11.25 13V9C11.25 8.58579 11.5858 8.25 12 8.25Z" fill="currentColor"/>
</svg>`;

const UPDATE_ICON_HTML = `<svg class="icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M4 12A8 8 0 0 1 18.93 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
<path d="M20 12A8 8 0 0 1 5.07 16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
<polyline points="14 8 19 8 19 3" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
<polyline points="10 16 5 16 5 21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
</svg>`;

let isConfirmingClear = false;
let clearConfirmTimeout = null;

function handleDocumentClick(event) {
  const clearBtn = document.getElementById("clear-history-btn");
  if (clearBtn && !clearBtn.contains(event.target)) {
    resetClearConfirmState();
  }
}

function resetClearConfirmState() {
  isConfirmingClear = false;
  if (clearConfirmTimeout) {
    clearTimeout(clearConfirmTimeout);
    clearConfirmTimeout = null;
  }
  document.removeEventListener("click", handleDocumentClick);
  const clearBtn = document.getElementById("clear-history-btn");
  if (clearBtn) {
    clearBtn.classList.remove("confirming");
  }
}

// Switch between panels cleanly (only one active, toggling off goes back to test-view)
export function showView(viewId) {
  const views = ["test-view", "history-view", "update-view"];
  views.forEach(id => {
    const el = document.getElementById(id);
    if (el) {
      if (id === viewId) {
        el.classList.remove("view-hidden");
        el.classList.add("view-active");
      } else {
        el.classList.remove("view-active");
        el.classList.add("view-hidden");
      }
    }
  });

  // Toggle selected class
  const historyToggleBtn = document.getElementById("history-toggle-btn");
  if (historyToggleBtn) {
    if (viewId === "history-view") {
      historyToggleBtn.classList.add("selected");
    } else {
      historyToggleBtn.classList.remove("selected");
    }
  }

  const updateToggleBtn = document.getElementById("update-toggle-btn");
  if (updateToggleBtn) {
    if (viewId === "update-view") {
      updateToggleBtn.classList.add("selected");
    } else {
      updateToggleBtn.classList.remove("selected");
    }
  }
}

export async function handleHistoryToggleClick() {
  if (testState.isTesting) return;

  const historyView = document.getElementById("history-view");
  if (!historyView) return;

  resetClearConfirmState();

  if (historyView.classList.contains("view-active")) {
    showView("test-view");
  } else {
    showView("history-view");
    try {
      const history = await window.go.gui_wails.App.GetHistory();
      renderHistory(history);
    } catch (err) {
      console.error("Failed to load history:", err);
    }
  }
  window.focus();
}

export async function handleClearHistoryClick(event) {
  if (event) {
    event.stopPropagation();
  }
  const clearBtn = document.getElementById("clear-history-btn");
  if (!clearBtn) return;

  if (!isConfirmingClear) {
    isConfirmingClear = true;
    clearBtn.classList.add("confirming");
    document.addEventListener("click", handleDocumentClick);
    clearConfirmTimeout = setTimeout(() => {
      resetClearConfirmState();
    }, 3000);
  } else {
    resetClearConfirmState();
    try {
      await window.go.gui_wails.App.ClearHistory();
      renderHistory([]);
    } catch (err) {
      console.error("Failed to clear history:", err);
    }
  }
}

export async function handleOpenJsonClick() {
  try {
    await window.go.gui_wails.App.OpenHistoryJSON();
  } catch (err) {
    console.error("Failed to open history JSON:", err);
  }
}

export function initializeHistoryHandlers() {
  const toggleBtn = document.getElementById("history-toggle-btn");
  if (toggleBtn) {
    toggleBtn.addEventListener("click", handleHistoryToggleClick);
  }

  const clearBtn = document.getElementById("clear-history-btn");
  if (clearBtn) {
    clearBtn.addEventListener("click", handleClearHistoryClick);
  }

  const openJsonBtn = document.getElementById("open-json-btn");
  if (openJsonBtn) {
    openJsonBtn.addEventListener("click", handleOpenJsonClick);
  }
}

export function handleCloseClick() {
  window.go.gui_wails.App.HideWindow();
}

export function initializeCloseHandler() {
  const closeBtn = document.getElementById("close-btn");
  if (closeBtn) {
    closeBtn.addEventListener("click", handleCloseClick);
  }
}

export function handleBannerClick() {
  if (testState.isTesting) return;
  showView("test-view");
}

export function initializeBannerHandler() {
  const banner = document.getElementById("header-banner");
  if (banner) {
    banner.addEventListener("click", handleBannerClick);
  }
}

export function showUpdateState(stateName) {
  const states = ["update-checking-state", "update-uptodate-state", "update-available-state", "update-error-state"];
  states.forEach(id => {
    const el = document.getElementById(id);
    if (el) {
      if (id === stateName) {
        el.classList.remove("view-hidden");
      } else {
        el.classList.add("view-hidden");
      }
    }
  });
}

let updateData = null;

export async function performUpdateCheck() {
  showUpdateState("update-checking-state");

  try {
    const info = await window.go.gui_wails.App.CheckForUpdate();
    window.localStorage.setItem("last_update_check_time", Date.now().toString());

    if (info.HasUpdate) {
      updateData = info;
      const versionVal = document.getElementById("update-version-val");
      const sizeVal = document.getElementById("update-size-val");
      if (versionVal) versionVal.innerText = "v" + info.LatestVersion;
      if (sizeVal) {
        const mb = (info.AssetSizeBytes / (1024 * 1024)).toFixed(2);
        sizeVal.innerText = mb + " MB";
      }

      const toggleBtn = document.getElementById("update-toggle-btn");
      if (toggleBtn) {
        toggleBtn.classList.add("has-badge");
        toggleBtn.title = "Update Available";
      }

      showUpdateState("update-available-state");
    } else {
      updateData = null;
      const versionVal = document.getElementById("current-version-display");
      if (versionVal) {
        versionVal.innerText = "v" + CONFIG.APP_VERSION;
      }
      const toggleBtn = document.getElementById("update-toggle-btn");
      if (toggleBtn) {
        toggleBtn.title = "Check for Updates";
      }
      showUpdateState("update-uptodate-state");
    }
  } catch (err) {
    console.error("Failed manual update check:", err);
    const errorMsgEl = document.getElementById("update-error-msg");
    if (errorMsgEl) {
      errorMsgEl.innerText = (err && err.message) ? err.message : (err || "Failed to check for updates");
    }
    showUpdateState("update-error-state");
  }
}

export function handleUpdateAvailable(info) {
  updateData = info;
  const toggleBtn = document.getElementById("update-toggle-btn");
  const versionVal = document.getElementById("update-version-val");
  const sizeVal = document.getElementById("update-size-val");

  if (toggleBtn) {
    toggleBtn.classList.add("has-badge");
    toggleBtn.title = "Update Available";
  }

  if (versionVal) {
    versionVal.innerText = "v" + info.LatestVersion;
  }

  if (sizeVal) {
    const mb = (info.AssetSizeBytes / (1024 * 1024)).toFixed(2);
    sizeVal.innerText = mb + " MB";
  }

  const updateView = document.getElementById("update-view");
  if (updateView && updateView.classList.contains("view-active")) {
    showUpdateState("update-available-state");
  }
}

export function handleUpdateToggleClick() {
  if (testState.isTesting) return;

  const updateView = document.getElementById("update-view");
  if (!updateView) return;

  const toggleBtn = document.getElementById("update-toggle-btn");
  if (toggleBtn) {
    toggleBtn.classList.remove("has-badge");
    toggleBtn.title = "Check for Updates";
  }

  if (updateView.classList.contains("view-active")) {
    showView("test-view");
  } else {
    showView("update-view");

    const lastCheck = window.localStorage.getItem("last_update_check_time");
    const sixHours = 6 * 60 * 60 * 1000;
    const now = Date.now();

    if (updateData) {
      showUpdateState("update-available-state");
    } else if (lastCheck && now - parseInt(lastCheck) < sixHours) {
      showUpdateState("update-uptodate-state");
    } else {
      performUpdateCheck();
    }
  }
}

export function handleUpdateNowClick() {
  const view = document.getElementById("update-view");
  if (!view) return;

  view.innerHTML = `
      <div class="update-state-container">
          <div class="update-status-msg" id="update-install-status">Downloading update...</div>
          <div class="update-progress-bar-container">
              <div class="update-progress-fill" id="update-progress-fill"></div>
          </div>
          <div class="update-progress-text" id="update-progress-percent">0%</div>
      </div>
  `;

  window.go.gui_wails.App.ApplyUpdate();
}

export function handleUpdateProgress(percent) {
  const fill = document.getElementById("update-progress-fill");
  const percentText = document.getElementById("update-progress-percent");
  const statusText = document.getElementById("update-install-status");

  if (fill) fill.style.width = percent + "%";
  if (percentText) percentText.textContent = percent + "%";

  if (percent >= 100 && statusText) {
    statusText.textContent = "Installing update...";
  }
}

export function handleUpdateError(err) {
  const statusText = document.getElementById("update-install-status");
  if (statusText) {
    statusText.textContent = "Error: " + err;
    statusText.style.color = "var(--danger)";
  }

  setTimeout(() => {
    const view = document.getElementById("update-view");
    if (view && updateData) {
      // Restore update available state view structure
      view.innerHTML = `
        <div id="update-checking-state" class="update-state-container view-hidden">
            <div class="loader"></div>
            <div class="update-status-msg">Checking for updates...</div>
        </div>
        <div id="update-uptodate-state" class="update-state-container view-hidden">
            <h2 class="update-title">SpeedTest Tray</h2>
            <div class="update-center-content">
                <svg class="update-dim-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
                    <polyline points="22 4 12 14.01 9 11.01" />
                </svg>
                <span class="update-version-text" id="current-version-display">--</span>
                <div class="update-size-container">
                    <span class="update-size-label">You're up to date!</span>
                </div>
            </div>
            <div class="update-actions">
                <button id="manual-check-btn" class="update-action-btn primary-btn">Check again</button>
            </div>
        </div>
        <div id="update-available-state" class="update-state-container view-hidden">
            <h2 class="update-title">New Version Available!</h2>
            <div class="update-center-content">
                <svg class="update-dim-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M4 12A8 8 0 0 1 18.93 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M20 12A8 8 0 0 1 5.07 16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <polyline points="14 8 19 8 19 3" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <polyline points="10 16 5 16 5 21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                <span id="update-version-val" class="update-version-text">--</span>
                <div class="update-size-container">
                    <span class="update-size-label">Update Size:</span>
                    <span id="update-size-val" class="update-size-val">--</span>
                </div>
            </div>
            <div class="update-notes-link">
                <a href="#" id="update-notes-btn">View Release Notes</a>
            </div>
            <div class="update-actions">
                <button id="update-now-btn" class="update-action-btn primary-btn">Update Now</button>
                <button id="update-skip-btn" class="update-action-btn secondary-btn">Skip Version</button>
            </div>
        </div>
        <div id="update-error-state" class="update-state-container view-hidden">
            <div class="update-center-content">
                <svg class="update-dim-icon" viewBox="0 0 24 24" fill="none" stroke="var(--danger)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="12" y1="8" x2="12" y2="12"></line>
                    <line x1="12" y1="16" x2="12.01" y2="16"></line>
                </svg>
                <span class="update-version-text" style="color: var(--danger);">Check Failed</span>
                <div class="update-size-container">
                    <span class="update-size-label" id="update-error-msg">Failed to check for updates</span>
                </div>
            </div>
            <div class="update-actions">
                <button id="error-check-btn" class="update-action-btn primary-btn">Try again</button>
            </div>
        </div>
      `;
      // Re-initialize update view button listeners and state
      initializeUpdateHandlers();
      handleUpdateAvailable(updateData);
    }
  }, 3000);
}

export function handleUpdateSkipClick() {
  if (!updateData) return;

  window.go.gui_wails.App.SkipUpdate(updateData.LatestVersion);

  updateData = null;
  const toggleBtn = document.getElementById("update-toggle-btn");
  if (toggleBtn) {
    toggleBtn.classList.remove("has-badge");
    toggleBtn.title = "Check for Updates";
  }

  showView("test-view");
}

export function handleReleaseNotesClick(e) {
  e.preventDefault();
  if (updateData && updateData.ReleasePageURL) {
    window.runtime.BrowserOpenURL(updateData.ReleasePageURL);
  }
}

export function initializeUpdateHandlers() {
  const versionDisplay = document.getElementById("current-version-display");
  if (versionDisplay) {
    versionDisplay.innerText = "v" + CONFIG.APP_VERSION;
  }

  const toggleBtn = document.getElementById("update-toggle-btn");
  if (toggleBtn) {
    toggleBtn.addEventListener("click", handleUpdateToggleClick);
  }

  const updateNowBtn = document.getElementById("update-now-btn");
  if (updateNowBtn) {
    updateNowBtn.addEventListener("click", handleUpdateNowClick);
  }

  const updateSkipBtn = document.getElementById("update-skip-btn");
  if (updateSkipBtn) {
    updateSkipBtn.addEventListener("click", handleUpdateSkipClick);
  }

  const notesBtn = document.getElementById("update-notes-btn");
  if (notesBtn) {
    notesBtn.addEventListener("click", handleReleaseNotesClick);
  }

  const manualCheckBtn = document.getElementById("manual-check-btn");
  if (manualCheckBtn) {
    manualCheckBtn.addEventListener("click", () => performUpdateCheck(true));
  }

  const errorCheckBtn = document.getElementById("error-check-btn");
  if (errorCheckBtn) {
    errorCheckBtn.addEventListener("click", () => performUpdateCheck());
  }
}

