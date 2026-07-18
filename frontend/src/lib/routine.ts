import type { Course } from "@bindings/routine";

export type CourseStatus = "ongoing" | "upcoming" | "inactive";

export const getCourseStatus = (course: Course, now: Date = new Date()): CourseStatus => {
  const daysMap = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
  const currentDay = daysMap[now.getDay()];

  if (course.day !== currentDay) return "inactive";

  const currentMinutes = now.getHours() * 60 + now.getMinutes();
  const startMinutes = parseTimeToMinutes(course.startTime);
  const endMinutes = parseTimeToMinutes(course.endTime);

  if (currentMinutes >= startMinutes && currentMinutes <= endMinutes) {
    return "ongoing";
  } else if (currentMinutes < startMinutes && startMinutes - currentMinutes <= 60) {
    return "upcoming";
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
