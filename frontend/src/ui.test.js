// @vitest-environment jsdom
import { beforeEach, describe, expect, it, vi } from "vitest";
import { CONFIG, ERRORS, MESSAGES, PHASES } from "./constants.js";
import { testState } from "./state.js";
import {
  handleTestComplete,
  handleTestError,
  handleTestUpdate,
  initializeElements,
  resetUI,
  setButtonState,
  setStatus,
  updateGauge,
  updateResults,
  renderHistory,
  updateHistoryToggleState,
} from "./ui.js";

describe("ui", () => {
  beforeEach(() => {
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
    document.getElementById("speedometer").setMax = vi.fn();
    testState.stopTest();
    testState.resetResults();
    initializeElements();
  });

  it("updates result text and gauge", () => {
    updateResults("Server", "20 ms", "90 Mbps", "18 Mbps");
    updateGauge(90);

    expect(text("server")).toBe("Server");
    expect(text("ping")).toBe("20 ms");
    expect(text("download")).toBe("90 Mbps");
    expect(text("upload")).toBe("18 Mbps");
    expect(document.getElementById("speedometer").setValue).toHaveBeenCalledWith(90);
  });

  it("resets UI and button state", () => {
    resetUI();
    setStatus("Running");
    setButtonState(true);

    expect(text("server")).toBe("--");
    expect(text("status")).toBe("Running");
    expect(text("run-btn")).toBe("Stop");
    expect(document.getElementById("run-btn").disabled).toBe(false);
  });

  it("ignores updates when no test is running", () => {
    handleTestUpdate({ status: "Ignored", phase: PHASES.DOWNLOADING, download: "90" });

    expect(text("status")).toBe("Ready");
  });

  it("applies download updates while testing", () => {
    testState.startTest();

    handleTestUpdate({
      status: "Running download test...",
      phase: PHASES.DOWNLOADING,
      server: "Nearest",
      ping: "20",
      download: "90.5",
      upload: "0",
    });

    expect(text("server")).toBe("Nearest");
    expect(text("ping")).toBe("20 ms");
    expect(text("download")).toBe("90.5 Mbps");
    expect(document.getElementById("speedometer").setMax).toHaveBeenCalledWith(CONFIG.GAUGE_MAX_DOWNLOAD);
    expect(document.getElementById("speedometer").setValue).toHaveBeenCalledWith("90.5");
  });

  it("applies upload updates while testing", () => {
    testState.startTest();

    handleTestUpdate({
      status: "Running upload test...",
      phase: PHASES.UPLOADING,
      ping: "20",
      download: "90.5",
      upload: "18.2",
    });

    expect(text("upload")).toBe("18.2 Mbps");
    expect(document.getElementById("speedometer").setMax).toHaveBeenCalledWith(CONFIG.GAUGE_MAX_UPLOAD);
    expect(document.getElementById("speedometer").setValue).toHaveBeenCalledWith("18.2");
  });

  it("handles successful completion", () => {
    testState.startTest();

    handleTestComplete({
      server: "Nearest",
      ping: "20",
      download: "90.5",
      upload: "18.2",
    });

    expect(testState.isTesting).toBe(false);
    expect(text("status")).toBe("Test Completed");
    expect(text("run-btn")).toBe("Start Again");
    expect(text("server")).toBe("Nearest");
  });

  it("handles stopped completion", () => {
    testState.startTest();

    handleTestComplete({ error: "test stopped" });

    expect(text("status")).toBe("Test Stopped");
    expect(text("run-btn")).toBe("Try Again");
    expect(text("server")).toBe("--");
  });

  it("handles offline completion", () => {
    testState.startTest();

    handleTestComplete({ error: ERRORS.NO_INTERNET });

    expect(testState.isTesting).toBe(false);
    expect(text("status")).toBe(MESSAGES.NO_INTERNET);
    expect(text("run-btn")).toBe("Try Again");
    expect(text("server")).toBe("--");
  });

  it("handles test errors", () => {
    testState.startTest();

    handleTestError("network failed");

    expect(testState.isTesting).toBe(false);
    expect(text("status")).toBe("Error: network failed");
    expect(text("run-btn")).toBe("Try Again");
    expect(text("server")).toBe("--");
  });

  it("updates history toggle state", () => {
    const toggleBtn = document.getElementById("history-toggle-btn");

    setButtonState(true);
    expect(toggleBtn.disabled).toBe(true);
    expect(toggleBtn.style.opacity).toBe("0.5");

    setButtonState(false);
    expect(toggleBtn.disabled).toBe(false);
    expect(toggleBtn.style.opacity).toBe("1");
  });

  it("renders history entries", () => {
    const historyList = document.getElementById("history-list");
    const clearBtn = document.getElementById("clear-history-btn");

    renderHistory([]);
    expect(historyList.innerHTML).toContain("No test history yet");
    expect(clearBtn.disabled).toBe(true);

    const mockHistory = [
      {
        timestamp: "2026-06-16T15:53:05Z",
        server: "Test Server",
        ping: 15.5,
        download: 120.4,
        upload: 45.2,
      }
    ];
    renderHistory(mockHistory);
    expect(historyList.innerHTML).not.toContain("No test history yet");
    expect(historyList.innerHTML).toContain("Test Server");
    expect(historyList.innerHTML).toContain("120.4");
    expect(historyList.innerHTML).toContain("45.2");
    expect(historyList.innerHTML).toContain("16");
    expect(clearBtn.disabled).toBe(false);
  });
});

function text(id) {
  const element = document.getElementById(id);
  return element.innerText ?? element.textContent;
}
