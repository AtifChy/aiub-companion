import type { AcademicEvent } from "@bindings/calendar";

export function formatEventDate(event: AcademicEvent): string {
  const date = event.date;
  const endDate = event.endDate ?? null;

  const opts: Intl.DateTimeFormatOptions = {
    month: "short",
    day: "numeric",
  };

  if (date.getFullYear() !== new Date().getFullYear()) {
    opts.year = "numeric";
  }

  if (endDate) {
    return `${date.toLocaleDateString(undefined, opts)} – ${endDate.toLocaleDateString(undefined, opts)}`;
  }

  return date.toLocaleDateString(undefined, opts);
}

export function isEventOngoing(event: AcademicEvent, now = new Date()): boolean {
  const start = new Date(event.date);
  start.setHours(0, 0, 0, 0);

  const current = new Date(now);
  current.setHours(0, 0, 0, 0);

  const end = event.endDate ? new Date(event.endDate) : new Date(start);
  end.setHours(23, 59, 59, 999);

  return current >= start && current <= end;
}
