import { describe, it, expect } from "vitest";
import { bgColor, formatClock } from "./logic.js";

describe("bgColor", () => {
  it("green when safe", () => expect(bgColor(true)).toBe("#16a34a"));
  it("red when not safe", () => expect(bgColor(false)).toBe("#dc2626"));
  it("slate when unknown", () => expect(bgColor(null)).toBe("#334155"));
});

describe("formatClock", () => {
  // 2026-08-12T03:00:00Z = 10:00 in Jakarta (UTC+7)
  const d = new Date("2026-08-12T03:00:00Z");

  it("formats in the given timezone", () => {
    const s = formatClock(d, "Asia/Jakarta");
    expect(s).toContain("10:00:00");
    expect(s).toContain("2026");
  });

  it("shifts with timezone", () => {
    // same instant, Tokyo (UTC+9) = 12:00
    expect(formatClock(d, "Asia/Tokyo")).toContain("12:00:00");
  });

  it("returns empty string when timezone missing", () => {
    expect(formatClock(d, undefined)).toBe("");
    expect(formatClock(d, "")).toBe("");
  });

  it("returns empty string on a bad timezone", () => {
    expect(formatClock(d, "Not/AZone")).toBe("");
  });
});
