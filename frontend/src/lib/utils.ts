import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function titleCase(str: string) {
  return str
    .toLowerCase()
    .split(/\s+/)
    .map((word, index, arr) => {
      if (index !== 0 && index !== arr.length - 1 && word.length <= 3) {
        return word;
      }
      return word.charAt(0).toUpperCase() + word.slice(1);
    })
    .join(" ");
}

export function formatBytes(bytes: number): string {
  if (bytes <= 0 || !isFinite(bytes)) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB"];
  const exponent = Math.floor(Math.log(bytes) / Math.log(1024));
  const value = bytes / Math.pow(1024, exponent);
  return `${value.toFixed(2)} ${units[exponent] ?? ""}`;
}

export function formatSpeed(rate: number): string {
  if (rate <= 0 || !isFinite(rate)) return "0 B/s";
  return `${formatBytes(rate)}/s`;
}
