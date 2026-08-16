export type Color = "dark" | "light" | "system";
export type Theme = "default" | "claude" | "supabase" | "t3chat";

export interface ThemeEntry {
  id: Theme;
  label: string;
}

export const THEMES: ThemeEntry[] = [
  { id: "default", label: "Default" },
  { id: "claude", label: "Claude" },
  { id: "supabase", label: "Supabase" },
  { id: "t3chat", label: "T3 Chat" },
] as const;
