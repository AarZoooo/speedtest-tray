// @vitest-environment jsdom
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { testState } from "./state.js";
import { handleButtonClick, initializeButtonHandler, startTest, stopTest, handleHistoryToggleClick, handleClearHistoryClick, handleOpenJsonClick, initializeHistoryHandlers } from "./handlers.js";
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
      <button id="history-toggle-btn">🕒</button>
      <div id="test-view" class="view-active"></div>
      <div id="history-view" class="view-hidden">
        <div class="history-header">
          <button id="clear-history-btn">Clear history</button>
          <button id="open-json-btn">Open json</button>
        </div>
        <div id="history-list"></div>
      </div>
    `;
    document.getElementById("speedometer").setValue = vi.fn();
    testState.stopTest();
    initializeElements();
    window.go = {
      gui_wails: {
        App: {
          StartTest: vi.fn().mockResolvedValue(undefined),
          StopTest: vi.fn(),
          GetHistory: vi.fn().mockResolvedValue([]),
          ClearHistory: vi.fn().mockResolvedValue(undefined),
          OpenHistoryJSON: vi.fn().mockResolvedValue(undefined),
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

  it("toggles to history view and back", async () => {
    const testView = document.getElementById("test-view");
    const historyView = document.getElementById("history-view");
    const toggleBtn = document.getElementById("history-toggle-btn");

    expect(testView.classList.contains("view-active")).toBe(true);
    expect(historyView.classList.contains("view-hidden")).toBe(true);

    await handleHistoryToggleClick();

    expect(testView.classList.contains("view-hidden")).toBe(true);
    expect(historyView.classList.contains("view-active")).toBe(true);
    expect(toggleBtn.innerHTML).toContain("stroke-width=\"2\"");
    expect(window.go.gui_wails.App.GetHistory).toHaveBeenCalledOnce();

    await handleHistoryToggleClick();

    expect(testView.classList.contains("view-active")).toBe(true);
    expect(historyView.classList.contains("view-hidden")).toBe(true);
    expect(toggleBtn.innerHTML).toContain("fill-rule=\"evenodd\"");
  });

  it("ignores history toggle when testing is active", async () => {
    testState.startTest();
    await handleHistoryToggleClick();
    expect(document.getElementById("test-view").classList.contains("view-active")).toBe(true);
    expect(window.go.gui_wails.App.GetHistory).not.toHaveBeenCalled();
  });

  it("uses a 2-click flow to clear history", async () => {
    const clearBtn = document.getElementById("clear-history-btn");

    await handleClearHistoryClick();
    expect(clearBtn.innerText).toBe("Clear history (Sure?)");
    expect(clearBtn.classList.contains("danger")).toBe(true);
    expect(window.go.gui_wails.App.ClearHistory).not.toHaveBeenCalled();

    await handleClearHistoryClick();
    expect(clearBtn.innerText).toBe("Clear history");
    expect(clearBtn.classList.contains("danger")).toBe(false);
    expect(window.go.gui_wails.App.ClearHistory).toHaveBeenCalledOnce();
  });

  it("resets clear confirmation state after a timeout", () => {
    vi.useFakeTimers();
    const clearBtn = document.getElementById("clear-history-btn");

    handleClearHistoryClick();
    expect(clearBtn.innerText).toBe("Clear history (Sure?)");

    vi.advanceTimersByTime(3000);
    expect(clearBtn.innerText).toBe("Clear history");
    expect(clearBtn.classList.contains("danger")).toBe(false);
  });

  it("opens history json file", async () => {
    await handleOpenJsonClick();

    expect(window.go.gui_wails.App.OpenHistoryJSON).toHaveBeenCalledOnce();
  });

  it("registers history listeners", () => {
    initializeHistoryHandlers();

    document.getElementById("history-toggle-btn").click();
    expect(window.go.gui_wails.App.GetHistory).toHaveBeenCalledOnce();

    document.getElementById("open-json-btn").click();
    expect(window.go.gui_wails.App.OpenHistoryJSON).toHaveBeenCalledOnce();
  });
});
