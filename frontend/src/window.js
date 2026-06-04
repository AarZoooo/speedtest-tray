import { CONFIG, EVENTS } from "./constants.js";
import { testState } from "./state.js";
import { resetUI } from "./ui.js";

const TEXT = {
  DEFAULT_VAL: "--",
};

// Handle window shown event
export function onWindowShown() {
  testState.setCanHide(false);
  setTimeout(() => {
    testState.setCanHide(true);
  }, CONFIG.UI_HIDE_DELAY_MS);
}

// Handle window blur
export function onWindowBlur() {
  if (testState.canHide) {
    window.go.gui_wails.App.HideWindow();
    testState.setCanHide(false);
  }
}

// Handle visibility change
export function onVisibilityChange() {
  if (document.visibilityState === "hidden" && testState.canHide) {
    window.go.gui_wails.App.HideWindow();
    testState.setCanHide(false);
  }
}

// Initialize window event listeners
export function initializeWindowEvents() {
  window.runtime.EventsOn(EVENTS.WINDOW_SHOWN, onWindowShown);
  window.onblur = onWindowBlur;
  document.addEventListener("visibilitychange", onVisibilityChange);
}
