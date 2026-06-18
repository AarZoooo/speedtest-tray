// Import modules
import { PHASES, EVENTS, CONFIG } from "./src/constants.js";
import { testState } from "./src/state.js";
import { 
  initializeElements, 
  handleTestUpdate, 
  handleTestComplete, 
  handleTestError 
} from "./src/ui.js";
import { initializeWindowEvents } from "./src/window.js";
import { initializeButtonHandler, initializeHistoryHandlers, initializeCloseHandler, initializeUpdateHandlers, handleUpdateAvailable, handleUpdateProgress, handleUpdateError, initializeBannerHandler } from "./src/handlers.js";

// Initialize on DOM ready
document.addEventListener("DOMContentLoaded", function () {
  // Setup UI elements
  initializeElements();

  // Setup event listeners
  initializeWindowEvents();
  initializeButtonHandler();
  initializeHistoryHandlers();
  initializeCloseHandler();
  initializeUpdateHandlers();
  initializeBannerHandler();

  // Setup Wails runtime events
  window.runtime.EventsOn(EVENTS.TEST_UPDATE, handleTestUpdate);
  window.runtime.EventsOn(EVENTS.TEST_COMPLETE, handleTestComplete);
  window.runtime.EventsOn(EVENTS.TEST_ERROR, handleTestError);
  window.runtime.EventsOn("update:available", handleUpdateAvailable);
  window.runtime.EventsOn("update:progress", handleUpdateProgress);
  window.runtime.EventsOn("update:error", handleUpdateError);
});

