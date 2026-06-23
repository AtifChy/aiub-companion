import { Service as LogService } from "@bindings/log";

export const logger = {
  debug: (message: string, ...args: unknown[]) => {
    console.debug(message, ...args);
    void LogService.Debug(format(message, args));
  },
  info: (message: string, ...args: unknown[]) => {
    console.info(message, ...args);
    void LogService.Info(format(message, args));
  },
  warn: (message: string, ...args: unknown[]) => {
    console.warn(message, ...args);
    void LogService.Warn(format(message, args));
  },
  error: (message: string, ...args: unknown[]) => {
    console.error(message, ...args);
    void LogService.Error(format(message, args));
  },
};

function format(message: string, args: unknown[]): string {
  if (!args.length) return message;
  return `${message}${args.length ? ":" : ""} ${args.map((arg) => String((arg as Error).message)).join(" ")}`;
}
