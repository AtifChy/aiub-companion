import { ThemeProvider as NextThemesProvider } from "next-themes";
import { useEffect } from "react";

import { useSettings } from "@/components/providers/settings-provider";

export type Theme = "dark" | "light" | "system";

interface ThemeProviderProps {
  children: React.ReactNode;
  defaultTheme?: Theme;
}

export function ThemeProvider({ children }: ThemeProviderProps) {
  const { config } = useSettings();
  const theme = config.appearance.theme as Theme;
  const colorscheme = "default";

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", colorscheme);
  }, []);

  return (
    <NextThemesProvider attribute="class" defaultTheme={theme} enableSystem>
      {children}
    </NextThemesProvider>
  );
}
