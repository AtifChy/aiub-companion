import { ThemeProvider as NextThemesProvider } from "next-themes";

import { useSettings } from "@/components/providers/settings-provider";

export type Theme = "dark" | "light" | "system";

interface ThemeProviderProps {
  children: React.ReactNode;
  defaultTheme?: Theme;
}

export function ThemeProvider({ children }: ThemeProviderProps) {
  const { config } = useSettings();
  const theme = config.appearance.theme as Theme;

  return (
    <NextThemesProvider attribute="class" defaultTheme={theme} enableSystem>
      {children}
    </NextThemesProvider>
  );
}
