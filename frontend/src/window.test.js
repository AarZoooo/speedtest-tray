// @vitest-environment jsdom
import { beforeEach, describe, expect, it, vi } from "vitest";
import { EVENTS } from "./constants.js";
import { testState } from "./state.js";
import { initializeWindowEvents, onVisibilityChange, onWindowBlur, onWindowShown } from "./window.js";

describe("window events", () => {
  beforeEach(() => {
    vi.useFakeTimers();
    testState.stopTest();
    testState.setCanHide(false);
    window.go = {
      gui_wails: {
        App: {
          HideWindow: vi.fn(),
        },
      },
    };
    window.runtime = {
      EventsOn: vi.fn(),
    };
  });

  it("allows hiding after the show delay", () => {
    onWindowShown();
    expect(testState.canHide).toBe(false);

    vi.runAllTimers();

    expect(testState.canHide).toBe(true);
  });

  it("hides on blur when allowed", () => {
    testState.setCanHide(true);

    onWindowBlur();

    expect(window.go.gui_wails.App.HideWindow).toHaveBeenCalledOnce();
    expect(testState.canHide).toBe(false);
  });

  it("hides on visibility change when allowed", () => {
    Object.defineProperty(document, "visibilityState", {
      configurable: true,
      value: "hidden",
    });
    testState.setCanHide(true);

    onVisibilityChange();

    expect(window.go.gui_wails.App.HideWindow).toHaveBeenCalledOnce();
    expect(testState.canHide).toBe(false);
  });

  it("registers Wails and browser event handlers", () => {
    const addEventListener = vi.spyOn(document, "addEventListener");

    initializeWindowEvents();

    expect(window.runtime.EventsOn).toHaveBeenCalledWith(EVENTS.WINDOW_SHOWN, onWindowShown);
    expect(window.onblur).toBe(onWindowBlur);
    expect(addEventListener).toHaveBeenCalledWith("visibilitychange", onVisibilityChange);
  });
});
