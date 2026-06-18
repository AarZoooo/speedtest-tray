// @vitest-environment jsdom
import { describe, expect, it } from "vitest";
import { initializeWindowEvents } from "./window.js";

describe("window events", () => {
  it("initializes without throwing", () => {
    expect(() => initializeWindowEvents()).not.toThrow();
  });
});
