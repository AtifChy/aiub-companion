import { ThemeProvider as NextThemesProvider } from "next-themes";
import { useEffect } from "react";

import { useSettings } from "@/components/providers/settings-provider";

const THEME_STORAGE_KEY = "data-theme";

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const { config } = useSettings();
  const color = config.appearance.color;
  const theme = config.appearance.theme || "default";

  useEffect(() => {
    document.documentElement.setAttribute(THEME_STORAGE_KEY, theme);
    localStorage.setItem(THEME_STORAGE_KEY, theme);
  }, [theme]);

  useEffect(() => {
    const handleStoreage = (event: StorageEvent) => {
      if (event.key === THEME_STORAGE_KEY && event.newValue) {
        document.documentElement.setAttribute(THEME_STORAGE_KEY, event.newValue);
      }
    };

    window.addEventListener("storage", handleStoreage);
    return () => window.removeEventListener("storage", handleStoreage);
  }, []);

  return (
    <NextThemesProvider attribute="class" defaultTheme={color} enableSystem>
      {children}
    </NextThemesProvider>
  );
}
