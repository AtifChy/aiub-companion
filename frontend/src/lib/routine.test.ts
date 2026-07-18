import { Schedule } from "@bindings/routine";
import { describe, expect, it } from "vitest";

import { getCourseStatus, parseTimeToMinutes } from "./routine";

describe("Routine Utils", () => {
  describe("parseTimeToMinutes", () => {
    it("parses AM times correctly", () => {
      expect(parseTimeToMinutes("08:00 AM")).toBe(8 * 60);
      expect(parseTimeToMinutes("11:30 AM")).toBe(11 * 60 + 30);
    });

    it("parses PM times correctly", () => {
      expect(parseTimeToMinutes("02:00 PM")).toBe(14 * 60);
      expect(parseTimeToMinutes("04:45 PM")).toBe(16 * 60 + 45);
    });

    it("handles 12 PM (noon) correctly", () => {
      expect(parseTimeToMinutes("12:00 PM")).toBe(12 * 60);
      expect(parseTimeToMinutes("12:30 PM")).toBe(12 * 60 + 30);
    });

    it("handles 12 AM (midnight) correctly", () => {
      expect(parseTimeToMinutes("12:00 AM")).toBe(0);
      expect(parseTimeToMinutes("12:15 AM")).toBe(15);
    });

    it("returns 0 for invalid inputs", () => {
      expect(parseTimeToMinutes(undefined)).toBe(0);
      expect(parseTimeToMinutes("")).toBe(0);
      expect(parseTimeToMinutes("Invalid Time")).toBe(0);
    });
  });

  describe("getCourseStatus", () => {
    // Create a fixed reference time for testing: Tuesday at 10:00 AM
    // Using UTC Date constructor can be flaky, so we construct components
    // 2026-07-21 is a Tuesday
    const testNow = new Date("2026-07-21T10:00:00.000+06:00");

    const createSchedule = (day: string, start: string, end: string): Schedule => {
      const schedule = new Schedule();
      schedule.day = day;
      schedule.startTime = start;
      schedule.endTime = end;
      return schedule;
    };

    it("returns 'inactive' if the course is not today", () => {
      const course = createSchedule("Monday", "10:00 AM", "12:00 PM");
      expect(getCourseStatus(course, testNow)).toBe("inactive");
    });

    it("returns 'ongoing' if current time is within course bounds", () => {
      // Tuesday 09:30 AM to 11:30 AM -> 10:00 AM is inside
      const course = createSchedule("Tuesday", "09:30 AM", "11:30 AM");
      expect(getCourseStatus(course, testNow)).toBe("ongoing");
    });

    it("returns 'upcoming' if current time is exactly 60 minutes before start", () => {
      // Tuesday 11:00 AM -> 10:00 AM is exactly 60m before
      const course = createSchedule("Tuesday", "11:00 AM", "12:00 PM");
      expect(getCourseStatus(course, testNow)).toBe("upcoming");
    });

    it("returns 'upcoming' if current time is 30 minutes before start", () => {
      // Tuesday 10:30 AM -> 10:00 AM is 30m before
      const course = createSchedule("Tuesday", "10:30 AM", "12:00 PM");
      expect(getCourseStatus(course, testNow)).toBe("upcoming");
    });

    it("returns 'inactive' if current time is more than 60 minutes before start", () => {
      // Tuesday 11:05 AM -> 10:00 AM is 65m before
      const course = createSchedule("Tuesday", "11:05 AM", "12:00 PM");
      expect(getCourseStatus(course, testNow)).toBe("inactive");
    });

    it("returns 'inactive' if current time is after course ends", () => {
      // Tuesday 08:00 AM to 09:50 AM -> 10:00 AM is after
      const course = createSchedule("Tuesday", "08:00 AM", "09:50 AM");
      expect(getCourseStatus(course, testNow)).toBe("inactive");
    });
  });
});
