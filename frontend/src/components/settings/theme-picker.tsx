import { useTheme } from "next-themes";
import { useEffect, useState } from "react";

import { THEMES, type ThemeEntry } from "@/lib/themes";
import { cn } from "@/lib/utils";

interface PreviewColors {
  primary: string;
  secondary: string;
  accent: string;
  background: string;
}

type ThemeMap = Map<ThemeEntry["id"], PreviewColors>;

function useThemePreviewColors(themes: ThemeEntry[]): ThemeMap {
  const [colorMap, setColorMap] = useState<ThemeMap>(new Map());
  const { resolvedTheme } = useTheme();

  useEffect(() => {
    const probe = document.createElement("div");
    probe.style.position = "absolute";
    probe.style.visibility = "hidden";
    probe.style.pointerEvents = "none";
    document.body.appendChild(probe);

    const isDark = resolvedTheme === "dark";
    const map: ThemeMap = new Map();

    if (isDark) {
      probe.classList.add("dark");
    }

    for (const theme of themes) {
      probe.setAttribute("data-theme", theme.id);

      const style = getComputedStyle(probe);
      map.set(theme.id, {
        primary: style.getPropertyValue("--primary").trim(),
        secondary: style.getPropertyValue("--secondary").trim(),
        accent: style.getPropertyValue("--accent").trim(),
        background: style.getPropertyValue("--background").trim(),
      });

      // oxlint-disable-next-line react-you-might-not-need-an-effect/no-adjust-state-on-prop-change
      setColorMap(map);
    }

    return () => {
      probe.remove();
    };
  }, [resolvedTheme, themes]);

  return colorMap;
}

interface ThemePickerProps {
  selectedTheme: string;
  onSelectTheme: (themeId: string) => void;
}

export function ThemePicker({ selectedTheme, onSelectTheme }: ThemePickerProps) {
  const previewColors = useThemePreviewColors(THEMES);

  return (
    <div className="mx-auto grid max-w-2xl grid-cols-2 gap-3 sm:grid-cols-4">
      {THEMES.map((theme) => {
        const isSelected = selectedTheme === theme.id;
        const colors = previewColors.get(theme.id);

        return (
          <button
            key={theme.id}
            type="button"
            onClick={() => onSelectTheme(theme.id)}
            className={cn(
              "flex items-center gap-3 overflow-hidden rounded-xl border p-3 transition-all",
              "hover:bg-accent/30",
              isSelected
                ? "border-primary ring-2 ring-primary/40"
                : "border-border/60 bg-card/40 hover:border-border",
            )}
          >
            <div
              className="flex h-9 w-11 shrink-0 items-center justify-center gap-1 rounded-lg border"
              style={{ backgroundColor: colors?.background }}
            >
              <div
                className="h-5 w-1.5 rounded-full"
                style={{ backgroundColor: colors?.primary }}
              />
              <div className="h-5 w-1.5 rounded-full" style={{ backgroundColor: colors?.accent }} />
              <div
                className="h-5 w-1.5 rounded-full opacity-40"
                style={{ backgroundColor: colors?.secondary }}
              />
            </div>

            <span
              className={cn(
                "text-sm font-medium text-nowrap",
                isSelected ? "text-primary/90" : "text-muted-foreground",
              )}
            >
              {theme.label}
            </span>
          </button>
        );
      })}
    </div>
  );
}
