export interface WailsError {
  message: string;
  cause: unknown;
  kind: string;
}

function isWailsError(value: unknown): value is WailsError {
  return (
    typeof value === "object" &&
    value !== null &&
    "message" in value &&
    "kind" in value
  );
}

export function parseWailsError(raw: unknown): WailsError | null {
  const value = raw instanceof Error ? raw.message : raw;

  if (typeof value === "string") {
    try {
      const parsed: unknown = JSON.parse(value);
      if (isWailsError(parsed)) return parsed;
    } catch {
      return { message: value, cause: null, kind: "Unknown" };
    }
  }

  return null;
}
