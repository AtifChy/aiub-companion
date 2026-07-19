import type { Schedule } from "@bindings/routine";

export type CourseStatus = "ongoing" | "upcoming" | "inactive";

export const DAYS: readonly string[] = [
  "Sunday",
  "Monday",
  "Tuesday",
  "Wednesday",
  "Thursday",
  "Friday",
  "Saturday",
];

export const getCourseStatus = (schedule: Schedule, now: Date = new Date()): CourseStatus => {
  const currentDay = DAYS[now.getDay()];
  const currentMinutes = now.getHours() * 60 + now.getMinutes();

  if (schedule.day === currentDay) {
    const startMinutes = parseTimeToMinutes(schedule.startTime);
    const endMinutes = parseTimeToMinutes(schedule.endTime);

    if (currentMinutes >= startMinutes && currentMinutes <= endMinutes) {
      return "ongoing";
    }

    if (currentMinutes < startMinutes && startMinutes - currentMinutes <= 60) {
      return "upcoming";
    }
  }

  return "inactive";
};

export const parseTimeToMinutes = (timeStr: string | undefined): number => {
  if (!timeStr) return 0;

  const match = /^(\d{1,2}):(\d{2})\s*(AM|PM)$/i.exec(timeStr.trim());
  if (!match) return 0;

  const [, hourStr, minuteStr, period] = match;

  let hours = Number(hourStr);
  const minutes = Number(minuteStr);

  if (hours < 1 || hours > 12 || minutes < 0 || minutes >= 60) return 0;

  if (period?.toUpperCase() === "PM" && hours !== 12) {
    hours += 12;
  } else if (period?.toUpperCase() === "AM" && hours === 12) {
    hours = 0;
  }

  return hours * 60 + minutes;
};
