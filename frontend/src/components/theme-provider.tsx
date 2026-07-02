import { useSettings } from "@/components/settings-provider";
import {
  ThemeProvider as NextThemesProvider,
  useTheme as useNextTheme,
} from "next-themes";
import { createContext, useEffect } from "react";

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

  return (
    <ThemeProviderContext.Provider {...props} value={theme}>
      <NextThemesProvider
        attribute="class"
        defaultTheme={defaultTheme}
        enableSystem
      >
        <ThemeSync theme={theme} />
        {children}
      </NextThemesProvider>
    </ThemeProviderContext.Provider>
  );
}

function ThemeSync({ theme }: { theme: Theme }) {
  const { setTheme } = useNextTheme();

  useEffect(() => {
    setTheme(theme);
  }, [theme, setTheme]);

  return null;
}
