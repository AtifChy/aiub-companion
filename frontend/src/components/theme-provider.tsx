import { useSettings } from "@/components/settings-provider";
import {
  ThemeProvider as NextThemesProvider,
  useTheme as useNextTheme,
} from "next-themes";
import { useEffect } from "react";

export type Theme = "dark" | "light" | "system";

type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: Theme;
};

export function ThemeProvider({
  children,
  defaultTheme = "system",
}: ThemeProviderProps) {
  const { config } = useSettings();
  const theme = (config.theme as Theme) || defaultTheme;

  return (
    <NextThemesProvider
      attribute="class"
      defaultTheme={defaultTheme}
      enableSystem
    >
      <ThemeSync theme={theme} />
      {children}
    </NextThemesProvider>
  );
}

function ThemeSync({ theme }: { theme: Theme }) {
  const { setTheme } = useNextTheme();

  useEffect(() => {
    setTheme(theme);
  }, [theme, setTheme]);

  return null;
}
