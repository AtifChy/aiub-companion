import { THEMES } from "@/lib/themes";
import { cn } from "@/lib/utils";

interface ThemePickerProps {
  selectedTheme: string;
  onSelectTheme: (themeId: string) => void;
}

export function ThemePicker({ selectedTheme, onSelectTheme }: ThemePickerProps) {
  return (
    <div className="mx-auto grid max-w-2xl grid-cols-2 gap-3 sm:grid-cols-4">
      {THEMES.map((theme) => {
        const isSelected = selectedTheme === theme.id;

        return (
          <button
            key={theme.id}
            type="button"
            onClick={() => onSelectTheme(theme.id)}
            className={cn(
              "flex items-center gap-3 overflow-hidden rounded-xl border p-3 transition-all",
              "outline-none focus-visible:border focus-visible:border-primary focus-visible:ring-2 focus-visible:ring-ring/50",
              "hover:bg-accent/30",
              isSelected
                ? "border-primary/60 ring-1 ring-ring/50"
                : "border-border bg-background/40 hover:border-border",
            )}
          >
            <div
              data-theme={theme.id}
              className="flex h-9 w-11 shrink-0 items-center justify-center gap-1 rounded-lg border bg-background shadow-md"
            >
              <div className="h-5 w-1.5 rounded-full bg-primary" />
              <div className="h-5 w-1.5 rounded-full bg-accent" />
              <div className="h-5 w-1.5 rounded-full bg-secondary opacity-40" />
            </div>

            <span
              className={cn(
                "text-sm font-medium",
                isSelected ? "text-primary" : "text-muted-foreground",
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
