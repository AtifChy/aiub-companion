import type { Category } from "@/hooks/use-notices";

export const CATEGORIES: Category[] = [
  "all",
  "general",
  "admission",
  "exam",
  "registration",
  "internship",
  "scholarship",
  "payment",
  "holiday",
];

export type AltCategory = Exclude<Category, "all">;

export const categoryStyles: Record<AltCategory, string> = {
  admission: "border-blue-500/20 bg-blue-500/10 text-blue-500",
  exam: "border-red-500/20 bg-red-500/10 text-red-500",
  registration: "border-violet-500/20 bg-violet-500/10 text-violet-500",
  internship: "border-cyan-500/20 bg-cyan-500/10 text-cyan-500",
  scholarship: "border-yellow-500/20 bg-yellow-500/10 text-yellow-500",
  payment: "border-green-500/20 bg-green-500/10 text-green-500",
  holiday: "border-orange-500/20 bg-orange-500/10 text-orange-500",
  general: "border-zinc-500/20 bg-zinc-500/10 text-zinc-500",
};
