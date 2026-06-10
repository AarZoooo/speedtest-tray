import { describe, expect, it } from "vitest";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { CONFIG, ERRORS, MESSAGES, PHASES } from "./constants.js";

describe("constants", () => {
  it("keeps frontend phases aligned with Go phases", () => {
    const phasesGo = readFileSync(resolve(process.cwd(), "../internal/config/phases.go"), "utf8");
    const matches = [...phasesGo.matchAll(/Phase\w+\s+Phase = "([^"]+)"/g)].map((match) => match[1]);

    expect(Object.values(PHASES)).toEqual(matches);
  });

  it("keeps UI config aligned with Go constants", () => {
    const constantsGo = readFileSync(resolve(process.cwd(), "../internal/config/constants.go"), "utf8");

    expect(CONFIG.UI_HIDE_DELAY_MS).toBe(goIntConstant(constantsGo, "UIHideDelayMs"));
    expect(CONFIG.GAUGE_MAX_DOWNLOAD).toBe(goIntConstant(constantsGo, "GaugeMaxDownload"));
    expect(CONFIG.GAUGE_MAX_UPLOAD).toBe(goIntConstant(constantsGo, "GaugeMaxUpload"));
  });

  it("keeps error messages aligned with Go constants", () => {
    const constantsGo = readFileSync(resolve(process.cwd(), "../internal/config/constants.go"), "utf8");

    expect(ERRORS.NO_INTERNET).toBe(goStringConstant(constantsGo, "ErrNoInternet"));
    expect(MESSAGES.NO_INTERNET).toBe(goStringConstant(constantsGo, "MsgNoInternet"));
  });
});

function goStringConstant(source, name) {
  const match = source.match(new RegExp(`${name}\\s*=\\s*"([^"]+)"`));
  if (!match) {
    throw new Error(`missing Go constant ${name}`);
  }
  return match[1];
}

function goIntConstant(source, name) {
  const match = source.match(new RegExp(`${name}\\s*=\\s*(\\d+)`));
  if (!match) {
    throw new Error(`missing Go constant ${name}`);
  }
  return Number(match[1]);
}
