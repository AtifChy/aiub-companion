type Theme =
  | "default"
  | "claude"
  | "supabase"
  | "t3chat"
  | "github"
  | "amethyst"
  | "caffeine"
  | "catppuccin";

interface ThemeEntry {
  id: Theme;
  label: string;
}

export const THEMES: ThemeEntry[] = [
  { id: "default", label: "Default" },
  { id: "claude", label: "Claude" },
  { id: "supabase", label: "Supabase" },
  { id: "t3chat", label: "T3 Chat" },
  { id: "github", label: "GitHub" },
  { id: "amethyst", label: "Amethyst" },
  { id: "caffeine", label: "Caffeine" },
  { id: "catppuccin", label: "Catppuccin" },
] as const;
