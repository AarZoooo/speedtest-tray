// @vitest-environment jsdom
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { testState } from "./state.js";
import { 
  handleButtonClick, 
  initializeButtonHandler, 
  startTest, 
  stopTest, 
  handleHistoryToggleClick, 
  handleClearHistoryClick, 
  handleOpenJsonClick, 
  initializeHistoryHandlers, 
  handleCloseClick, 
  initializeCloseHandler,
  handleUpdateAvailable,
  handleUpdateToggleClick,
  handleUpdateNowClick,
  handleUpdateProgress,
  handleUpdateError,
  handleUpdateSkipClick,
  handleReleaseNotesClick,
  initializeUpdateHandlers,
  handleBannerClick,
  initializeBannerHandler
} from "./handlers.js";
import { initializeElements } from "./ui.js";

describe("handlers", () => {
  beforeEach(() => {
    vi.spyOn(console, "log").mockImplementation(() => {});
    vi.spyOn(console, "error").mockImplementation(() => {});
    document.body.innerHTML = `
      <img id="header-banner" src="assets/banner.png" />
      <span id="server">--</span>
      <span id="ping">--</span>
      <span id="download">--</span>
      <span id="upload">--</span>
      <div id="status">Ready</div>
      <button id="run-btn">Start</button>
      <speedometer-gauge id="speedometer"></speedometer-gauge>
      <button id="history-toggle-btn">🕒</button>
      <button id="close-btn">✖</button>
      <button id="update-toggle-btn">Update</button>
      <div id="test-view" class="view-active"></div>
      <div id="history-view" class="view-hidden">
        <div class="history-header">
          <button id="clear-history-btn">Clear history</button>
          <button id="open-json-btn">Open json</button>
        </div>
        <div id="history-list"></div>
      </div>
      <div id="update-view" class="view-hidden">
        <span id="update-version-val"></span>
        <span id="update-size-val"></span>
        <button id="update-now-btn">Update Now</button>
        <button id="update-skip-btn">Skip Version</button>
        <a id="update-notes-btn">Release Notes</a>
      </div>
    `;
    document.getElementById("speedometer").setValue = vi.fn();
    testState.stopTest();
    initializeElements();
    const store = {};
    window.localStorage = {
      getItem: vi.fn(key => store[key] || null),
      setItem: vi.fn((key, value) => { store[key] = value.toString(); }),
      removeItem: vi.fn(key => { delete store[key]; }),
      clear: vi.fn(() => { for (const key in store) delete store[key]; })
    };
    window.runtime = {
      BrowserOpenURL: vi.fn(),
    };
    window.go = {
      gui_wails: {
        App: {
          StartTest: vi.fn().mockResolvedValue(undefined),
          StopTest: vi.fn(),
          GetHistory: vi.fn().mockResolvedValue([]),
          ClearHistory: vi.fn().mockResolvedValue(undefined),
          OpenHistoryJSON: vi.fn().mockResolvedValue(undefined),
          HideWindow: vi.fn().mockResolvedValue(undefined),
          ApplyUpdate: vi.fn().mockResolvedValue(undefined),
          SkipUpdate: vi.fn().mockResolvedValue(undefined),
          CheckForUpdate: vi.fn().mockResolvedValue({
            HasUpdate: false,
            LatestVersion: "1.0.1",
            AssetSizeBytes: 0,
            ReleasePageURL: ""
          }),
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
    expect(toggleBtn.classList.contains("selected")).toBe(true);
    expect(window.go.gui_wails.App.GetHistory).toHaveBeenCalledOnce();

    await handleHistoryToggleClick();

    expect(testView.classList.contains("view-active")).toBe(true);
    expect(historyView.classList.contains("view-hidden")).toBe(true);
    expect(toggleBtn.classList.contains("selected")).toBe(false);
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

  it("hides the window on close click", () => {
    handleCloseClick();
    expect(window.go.gui_wails.App.HideWindow).toHaveBeenCalledOnce();
  });

  it("registers the close button click listener", () => {
    initializeCloseHandler();
    document.getElementById("close-btn").click();
    expect(window.go.gui_wails.App.HideWindow).toHaveBeenCalledOnce();
  });

  it("handles update available events by populating UI and showing badge", () => {
    const toggleBtn = document.getElementById("update-toggle-btn");

    handleUpdateAvailable({
      LatestVersion: "1.2.0",
      AssetSizeBytes: 1048576 * 4.5,
      ReleasePageURL: "https://github.com/foo/bar/releases",
    });

    expect(toggleBtn.classList.contains("has-badge")).toBe(true);
    expect(document.getElementById("update-version-val").innerText).toBe("v1.2.0");
    expect(document.getElementById("update-size-val").innerText).toBe("4.50 MB");
  });

  it("toggles update view and handles badge removal", () => {
    const testView = document.getElementById("test-view");
    const updateView = document.getElementById("update-view");
    const toggleBtn = document.getElementById("update-toggle-btn");

    handleUpdateAvailable({
      LatestVersion: "1.2.0",
      AssetSizeBytes: 1000000,
      ReleasePageURL: "https://github.com/foo/bar",
    });

    expect(toggleBtn.classList.contains("has-badge")).toBe(true);
    expect(testView.classList.contains("view-active")).toBe(true);
    expect(updateView.classList.contains("view-hidden")).toBe(true);

    handleUpdateToggleClick();

    expect(toggleBtn.classList.contains("has-badge")).toBe(false);
    expect(testView.classList.contains("view-hidden")).toBe(true);
    expect(updateView.classList.contains("view-active")).toBe(true);

    handleUpdateToggleClick();

    expect(testView.classList.contains("view-active")).toBe(true);
    expect(updateView.classList.contains("view-hidden")).toBe(true);
  });

  it("triggers update action, changes UI card to show progress bar", () => {
    handleUpdateNowClick();
    expect(window.go.gui_wails.App.ApplyUpdate).toHaveBeenCalledOnce();
    const view = document.getElementById("update-view");
    expect(view.innerHTML).toContain("update-progress-fill");
    expect(view.innerHTML).toContain("Downloading update...");
  });

  it("updates progress bar and percentage text", () => {
    handleUpdateNowClick();

    handleUpdateProgress(45);

    const fill = document.getElementById("update-progress-fill");
    const percentText = document.getElementById("update-progress-percent");
    const statusText = document.getElementById("update-install-status");

    expect(fill.style.width).toBe("45%");
    expect(percentText.textContent).toBe("45%");
    expect(statusText.textContent).toBe("Downloading update...");

    handleUpdateProgress(100);
    expect(statusText.textContent).toBe("Installing update...");
  });

  it("handles update error by displaying it and restoring view after timeout", () => {
    vi.useFakeTimers();

    handleUpdateAvailable({
      LatestVersion: "1.2.0",
      AssetSizeBytes: 1000000,
      ReleasePageURL: "https://github.com/foo/bar",
    });

    handleUpdateNowClick();
    handleUpdateError("mock error");

    const statusText = document.getElementById("update-install-status");
    expect(statusText.textContent).toBe("Error: mock error");
    expect(statusText.style.color).toBe("var(--danger)");

    // Advance time to trigger recovery
    vi.advanceTimersByTime(3000);

    // It should restore the original HTML view
    expect(document.getElementById("update-now-btn")).not.toBeNull();
  });

  it("skips the update version and removes badge", () => {
    handleUpdateAvailable({
      LatestVersion: "1.2.0",
      AssetSizeBytes: 1000000,
      ReleasePageURL: "https://github.com/foo/bar",
    });

    // Make update view active first
    handleUpdateToggleClick();
    expect(document.getElementById("update-view").classList.contains("view-active")).toBe(true);

    handleUpdateSkipClick();

    expect(window.go.gui_wails.App.SkipUpdate).toHaveBeenCalledWith("1.2.0");
    expect(document.getElementById("update-toggle-btn").classList.contains("has-badge")).toBe(false);
    expect(document.getElementById("update-view").classList.contains("view-hidden")).toBe(true);
    expect(document.getElementById("test-view").classList.contains("view-active")).toBe(true);
  });

  it("opens release notes url in browser", () => {
    handleUpdateAvailable({
      LatestVersion: "1.2.0",
      AssetSizeBytes: 1000000,
      ReleasePageURL: "https://github.com/foo/bar",
    });

    const event = { preventDefault: vi.fn() };
    handleReleaseNotesClick(event);

    expect(event.preventDefault).toHaveBeenCalledOnce();
    expect(window.runtime.BrowserOpenURL).toHaveBeenCalledWith("https://github.com/foo/bar");
  });

  it("registers update handler listeners correctly", () => {
    initializeUpdateHandlers();

    document.getElementById("update-toggle-btn").click();
    // Toggle switches views
    expect(document.getElementById("update-view").classList.contains("view-active")).toBe(true);

    document.getElementById("update-now-btn").click();
    expect(window.go.gui_wails.App.ApplyUpdate).toHaveBeenCalledOnce();
  });

  it("returns to test view on banner click", () => {
    // Navigate to history first
    document.getElementById("test-view").classList.replace("view-active", "view-hidden");
    document.getElementById("history-view").classList.replace("view-hidden", "view-active");

    handleBannerClick();

    expect(document.getElementById("test-view").classList.contains("view-active")).toBe(true);
    expect(document.getElementById("history-view").classList.contains("view-hidden")).toBe(true);
  });

  it("registers banner click handler listener", () => {
    initializeBannerHandler();

    // Navigate to history first
    document.getElementById("test-view").classList.replace("view-active", "view-hidden");
    document.getElementById("history-view").classList.replace("view-hidden", "view-active");

    document.getElementById("header-banner").click();

    expect(document.getElementById("test-view").classList.contains("view-active")).toBe(true);
    expect(document.getElementById("history-view").classList.contains("view-hidden")).toBe(true);
  });
});

