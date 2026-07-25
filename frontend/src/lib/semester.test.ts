import { AcademicEvent } from "@bindings/calendar";
import { describe, expect, it } from "vitest";

import { formatEventDate, isEventOngoing } from "./semester";

describe("formatEventDate", () => {
  it("formats a single date correctly in the current year", () => {
    const currentYear = new Date().getFullYear();
    const event = new AcademicEvent();
    event.date = new Date(`${String(currentYear)}-05-15T12:00:00Z`);

    const result = formatEventDate(event);

    // In English locales, this is usually "May 15" or "15 May"
    expect(result).toMatch(/May 15|15 May/);
    expect(result).not.toContain(String(currentYear)); // Should exclude current year
  });

  it("formats a single date correctly in a different year", () => {
    const event = new AcademicEvent();
    event.date = new Date("2020-01-01T12:00:00Z");

    const result = formatEventDate(event);

    expect(result).toMatch(/Jan 1, 2020|1 Jan 2020/);
    expect(result).toContain("2020");
  });

  it("formats a date range correctly", () => {
    const currentYear = new Date().getFullYear();
    const event = new AcademicEvent();
    event.date = new Date(`${String(currentYear)}-08-10T12:00:00Z`);
    event.endDate = new Date(`${String(currentYear)}-08-12T12:00:00Z`);

    const result = formatEventDate(event);

    expect(result).toContain("Aug");
    expect(result).toContain("10");
    expect(result).toContain("12");
    expect(result).toContain("–"); // en-dash
  });
});

describe("isEventOngoing", () => {
  it("returns true for an ongoing event", () => {
    const event = new AcademicEvent();
    event.date = new Date("2024-06-01T00:00:00Z");
    event.endDate = new Date("2024-06-10T23:59:59Z");

    const now = new Date("2024-06-05T12:00:00Z");
    const result = isEventOngoing(event, now);

    expect(result).toBe(true);
  });

  it("returns false for a past event", () => {
    const event = new AcademicEvent();
    event.date = new Date("2024-05-01T00:00:00Z");
    event.endDate = new Date("2024-05-10T23:59:59Z");

    const now = new Date("2024-06-05T12:00:00Z");
    const result = isEventOngoing(event, now);

    expect(result).toBe(false);
  });

  it("returns false for a future event", () => {
    const event = new AcademicEvent();
    event.date = new Date("2024-07-01T00:00:00Z");
    event.endDate = new Date("2024-07-10T23:59:59Z");

    const now = new Date("2024-06-05T12:00:00Z");
    const result = isEventOngoing(event, now);

    expect(result).toBe(false);
  });

  it("returns true for an event that starts today", () => {
    const event = new AcademicEvent();
    event.date = new Date("2024-06-05T00:00:00Z");
    event.endDate = new Date("2024-06-10T23:59:59Z");

    const now = new Date("2024-06-05T12:00:00Z");
    const result = isEventOngoing(event, now);

    expect(result).toBe(true);
  });

  it("returns true for an event that ends today", () => {
    const event = new AcademicEvent();
    event.date = new Date("2024-06-01T00:00:00Z");
    event.endDate = new Date("2024-06-05T23:59:59Z");

    const now = new Date("2024-06-05T12:00:00Z");
    const result = isEventOngoing(event, now);

    expect(result).toBe(true);
  });
});
