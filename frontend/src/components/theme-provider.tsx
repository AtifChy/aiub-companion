import { useSettings } from "@/components/settings-provider";
import { createContext, useContext, useEffect } from "react";

export type Theme = "dark" | "light" | "system";

const ThemeProviderContext = createContext<Theme>("system");

type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: Theme;
};

export function ThemeProvider({
  children,
  defaultTheme = "system",
  ...props
}: ThemeProviderProps) {
  const { config } = useSettings();
  const theme = (config.theme as Theme) || defaultTheme;

  useEffect(() => {
    const root = window.document.documentElement;

    const setTheme = (theme: Exclude<Theme, "system">) => {
      root.classList.remove("light", "dark");
      root.classList.add(theme);
    };

    if (theme === "system") {
      const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");

      const applySystemTheme = () => {
        setTheme(mediaQuery.matches ? "dark" : "light");
      };
      applySystemTheme();

      mediaQuery.addEventListener("change", applySystemTheme);
      return () => mediaQuery.removeEventListener("change", applySystemTheme);
    }

    setTheme(theme);

    return undefined;
  }, [theme]);

  return (
    <ThemeProviderContext.Provider {...props} value={theme}>
      {children}
    </ThemeProviderContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export const useTheme = () => {
  const context = useContext(ThemeProviderContext);

  if (context === undefined)
    throw new Error("useTheme must be used within a ThemeProvider");

  return context;
};
