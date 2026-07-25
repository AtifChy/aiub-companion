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
