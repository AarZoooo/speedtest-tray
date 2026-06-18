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
import { initializeButtonHandler, initializeHistoryHandlers, initializeCloseHandler } from "./src/handlers.js";

// Initialize on DOM ready
document.addEventListener("DOMContentLoaded", function () {
  // Setup UI elements
  initializeElements();

  // Setup event listeners
  initializeWindowEvents();
  initializeButtonHandler();
  initializeHistoryHandlers();
  initializeCloseHandler();

  // Setup Wails runtime events
  window.runtime.EventsOn(EVENTS.TEST_UPDATE, handleTestUpdate);
  window.runtime.EventsOn(EVENTS.TEST_COMPLETE, handleTestComplete);
  window.runtime.EventsOn(EVENTS.TEST_ERROR, handleTestError);
});
