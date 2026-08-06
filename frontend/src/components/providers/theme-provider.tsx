import { ThemeProvider as NextThemesProvider } from "next-themes";
import { useEffect } from "react";

import { useSettings } from "@/components/providers/settings-provider";

export type Color = "dark" | "light" | "system";

interface ThemeProviderProps {
  children: React.ReactNode;
}

export function ThemeProvider({ children }: ThemeProviderProps) {
  const { config } = useSettings();
  const color = config.appearance.color as Color;
  const theme = config.appearance.theme || "default";

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", theme);
  }, [theme]);

  return (
    <NextThemesProvider attribute="class" defaultTheme={color} enableSystem>
      {children}
    </NextThemesProvider>
  );
}
