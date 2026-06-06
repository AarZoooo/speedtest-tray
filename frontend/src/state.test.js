import { describe, expect, it } from "vitest";
import { PHASES } from "./constants.js";
import { TestState } from "./state.js";

describe("TestState", () => {
  it("starts with idle defaults", () => {
    const state = new TestState();

    expect(state.getState()).toEqual({
      isTesting: false,
      canHide: false,
      currentPhase: null,
      results: {
        server: "--",
        ping: "--",
        download: "--",
        upload: "--",
      },
    });
  });

  it("tracks test lifecycle and phase changes", () => {
    const state = new TestState();

    state.startTest();
    expect(state.isTesting).toBe(true);
    expect(state.currentPhase).toBe(PHASES.INITIALIZING);

    state.setPhase(PHASES.DOWNLOADING);
    expect(state.currentPhase).toBe(PHASES.DOWNLOADING);

    state.stopTest();
    expect(state.isTesting).toBe(false);
    expect(state.currentPhase).toBeNull();
  });

  it("updates, resets, and protects result state", () => {
    const state = new TestState();

    state.updateResults("Server", "20", "90", "18");
    const snapshot = state.getState();
    snapshot.results.server = "Mutated";

    expect(state.results.server).toBe("Server");

    state.resetResults();
    expect(state.results).toEqual({
      server: "--",
      ping: "--",
      download: "--",
      upload: "--",
    });
  });

  it("tracks hide eligibility", () => {
    const state = new TestState();

    state.setCanHide(true);
    expect(state.canHide).toBe(true);
  });
});
