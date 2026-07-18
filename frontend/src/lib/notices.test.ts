import { describe, expect, it } from "vitest";

import { formatDate } from "./notices";

describe("formatDate", () => {
  it("formats a standard ISO date correctly", () => {
    // Note: ToLocaleDateString relies on system timezone, but en-GB fixes the format structure
    const result = formatDate("2026-07-16T12:00:00Z");
    // Depending on the local timezone where test is run, it will be 16 Jul 2026
    expect(result).toMatch(/\d{1,2} [A-Z][a-z]{2} \d{4}/);
    expect(result).toContain("2026");
    expect(result).toContain("Jul");
  });

  it("returns empty string if raw date is empty", () => {
    expect(formatDate("")).toBe("");
  });

  it("returns the raw string if date is invalid", () => {
    expect(formatDate("invalid-date-format")).toBe("invalid-date-format");
  });

  it("handles standard YYYY-MM-DD date strings", () => {
    const result = formatDate("2026-01-05");
    expect(result).toContain("2026");
    expect(result).toContain("Jan");
  });
});
