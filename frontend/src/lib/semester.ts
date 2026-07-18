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
