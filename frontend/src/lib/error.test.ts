import { describe, expect, it } from "vitest";

import { parseWailsError } from "./error";

describe("parseWailsError", () => {
  it("parses a valid WailsError JSON string", () => {
    const errorJson = JSON.stringify({
      message: "Network failure",
      cause: "timeout",
      kind: "NetworkError",
    });

    const result = parseWailsError(errorJson);

    expect(result).toEqual({
      message: "Network failure",
      cause: "timeout",
      kind: "NetworkError",
    });
  });

  it("handles a plain string error gracefully", () => {
    const rawError = "Some generic error occurred";

    const result = parseWailsError(rawError);

    expect(result).toEqual({
      message: "Some generic error occurred",
      cause: null,
      kind: "Unknown",
    });
  });

  it("extracts the message if passed a standard Error object", () => {
    const errorObj = new Error("Something went wrong");

    const result = parseWailsError(errorObj);

    expect(result).toEqual({
      message: "Something went wrong",
      cause: null,
      kind: "Unknown",
    });
  });

  it("extracts the message from an Error object that contains JSON", () => {
    const errorJson = JSON.stringify({
      message: "Database connection lost",
      cause: "EOF",
      kind: "DBError",
    });
    const errorObj = new Error(errorJson);

    const result = parseWailsError(errorObj);

    expect(result).toEqual({
      message: "Database connection lost",
      cause: "EOF",
      kind: "DBError",
    });
  });

  it("returns null for completely invalid types (like numbers or arrays)", () => {
    expect(parseWailsError(12345)).toBeNull();
    expect(parseWailsError([1, 2, 3])).toBeNull();
    expect(parseWailsError(null)).toBeNull();
  });
});
