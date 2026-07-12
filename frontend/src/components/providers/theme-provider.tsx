import { ThemeProvider as NextThemesProvider } from "next-themes";

import { useSettings } from "@/components/providers/settings-provider";

export type Theme = "dark" | "light" | "system";

interface ThemeProviderProps {
  children: React.ReactNode;
  defaultTheme?: Theme;
}

export function ThemeProvider({ children, defaultTheme = "system" }: ThemeProviderProps) {
  const { config } = useSettings();
  const theme = config.theme as Theme;

  return (
    <NextThemesProvider
      attribute="class"
      defaultTheme={defaultTheme}
      enableSystem
      forcedTheme={theme}
    >
      {children}
    </NextThemesProvider>
  );
}
