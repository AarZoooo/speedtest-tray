import { PHASES, CONFIG } from "./constants.js";

// Test state manager
class TestState {
  constructor() {
    this.isTesting = false;
    this.canHide = false;
    this.currentPhase = null;
    this.results = {
      server: "--",
      ping: "--",
      download: "--",
      upload: "--",
    };
  }

  startTest() {
    this.isTesting = true;
    this.currentPhase = PHASES.INITIALIZING;
  }

  stopTest() {
    this.isTesting = false;
    this.currentPhase = null;
  }

  setPhase(phase) {
    this.currentPhase = phase;
  }

  updateResults(server, ping, download, upload) {
    this.results = { server, ping, download, upload };
  }

  resetResults() {
    this.results = {
      server: "--",
      ping: "--",
      download: "--",
      upload: "--",
    };
  }

  setCanHide(value) {
    this.canHide = value;
  }

  getState() {
    return {
      isTesting: this.isTesting,
      canHide: this.canHide,
      currentPhase: this.currentPhase,
      results: { ...this.results },
    };
  }
}

export const testState = new TestState();
