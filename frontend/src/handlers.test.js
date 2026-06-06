// @vitest-environment jsdom
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { testState } from "./state.js";
import { handleButtonClick, initializeButtonHandler, startTest, stopTest } from "./handlers.js";
import { initializeElements } from "./ui.js";

describe("handlers", () => {
  beforeEach(() => {
    vi.spyOn(console, "log").mockImplementation(() => {});
    vi.spyOn(console, "error").mockImplementation(() => {});
    document.body.innerHTML = `
      <span id="server">--</span>
      <span id="ping">--</span>
      <span id="download">--</span>
      <span id="upload">--</span>
      <div id="status">Ready</div>
      <button id="run-btn">Start</button>
      <speedometer-gauge id="speedometer"></speedometer-gauge>
    `;
    document.getElementById("speedometer").setValue = vi.fn();
    testState.stopTest();
    initializeElements();
    window.go = {
      gui_wails: {
        App: {
          StartTest: vi.fn().mockResolvedValue(undefined),
          StopTest: vi.fn(),
        },
      },
    };
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("starts a test through the backend", async () => {
    await startTest();

    expect(testState.isTesting).toBe(true);
    expect(document.getElementById("run-btn").innerText).toBe("Stop");
    expect(window.go.gui_wails.App.StartTest).toHaveBeenCalledOnce();
  });

  it("rolls back state when backend start fails", async () => {
    window.go.gui_wails.App.StartTest.mockRejectedValue(new Error("failed"));

    await startTest();

    expect(testState.isTesting).toBe(false);
    expect(document.getElementById("run-btn").innerText).toBe("Try Again");
  });

  it("stops a running test through the backend", () => {
    testState.startTest();

    stopTest();

    expect(testState.isTesting).toBe(false);
    expect(document.getElementById("run-btn").disabled).toBe(true);
    expect(window.go.gui_wails.App.StopTest).toHaveBeenCalledOnce();
  });

  it("routes button clicks by current state", () => {
    handleButtonClick();
    expect(window.go.gui_wails.App.StartTest).toHaveBeenCalledOnce();

    testState.startTest();
    handleButtonClick();
    expect(window.go.gui_wails.App.StopTest).toHaveBeenCalledOnce();
  });

  it("registers the button click listener", () => {
    initializeButtonHandler();

    document.getElementById("run-btn").click();

    expect(window.go.gui_wails.App.StartTest).toHaveBeenCalledOnce();
  });
});
