import type { Settings } from "@bindings/settings";
import { GetSettings, SaveSettings } from "@bindings/settings/service";
import { Loader2Icon } from "lucide-react";
import { rawReturn } from "mutative";
import {
  createContext,
  useContext,
  useEffect,
  useRef,
  type ReactNode,
} from "react";
import { toast } from "sonner";
import { useMutative, type Updater } from "use-mutative";

interface SettingsContextType {
  config: Settings;
  updateConfig: Updater<Settings | undefined>;
}

const SettingsContext = createContext<SettingsContextType | null>(null);

export function SettingsProvider({ children }: { children: ReactNode }) {
  const [config, updateConfig] = useMutative<Settings | undefined>(undefined);
  const hydratedRef = useRef(false);

  useEffect(() => {
    GetSettings()
      .then((config) =>
        updateConfig(() => rawReturn(JSON.parse(JSON.stringify(config)))),
      )
      .catch((err) =>
        toast.error("Failed to load settings", {
          description: (err as Error).message,
        }),
      );
  }, [updateConfig]);

  useEffect(() => {
    if (!config) return;
    if (!hydratedRef.current) {
      hydratedRef.current = true;
      return;
    }
    const timer = setTimeout(
      () =>
        SaveSettings(config).catch((err) =>
          toast.error("Failed to save settings", {
            description: (err as Error).message,
          }),
        ),
      200,
    );
    return () => clearTimeout(timer);
  }, [config]);

  if (!config) {
    return (
      <div className="flex h-screen w-screen items-center justify-center">
        <Loader2Icon className="text-muted-foreground size-8 animate-spin" />
      </div>
    );
  }

  return (
    <SettingsContext.Provider value={{ config, updateConfig }}>
      {children}
    </SettingsContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export function useSettings() {
  const context = useContext(SettingsContext);
  if (!context)
    throw new Error("useSettings must be used within a SettingsProvider");
  return context;
}
