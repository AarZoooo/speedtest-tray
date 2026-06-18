import { PHASES } from "./constants.js";

// Test state manager
export class TestState {
  constructor() {
    this.isTesting = false;
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

  getState() {
    return {
      isTesting: this.isTesting,
      currentPhase: this.currentPhase,
      results: { ...this.results },
    };
  }
}

export const testState = new TestState();
