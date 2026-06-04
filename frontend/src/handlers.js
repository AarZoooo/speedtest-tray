import { testState } from "./state.js";
import { setStatus, setButtonState, resetUI } from "./ui.js";

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

  console.log("JS: Invoking backend StartTest");
  try {
    await window.go.gui_wails.App.StartTest();
    console.log("JS: Backend StartTest promise resolved");
  } catch (err) {
    console.error("JS: Backend StartTest failed:", err);
    testState.stopTest();
    setButtonState(false);
  }
}

// Stop the running test
export function stopTest() {
  console.log("JS: stopTest called");
  testState.stopTest();
  if (document.getElementById("run-btn")) {
    document.getElementById("run-btn").disabled = true;
  }
  setStatus("Test Stopped");
  resetUI(TEXT.DEFAULT_VAL);

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
