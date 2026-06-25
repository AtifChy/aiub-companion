import { logger } from "@/lib/logger";
import { Service as SettingsService, type Settings } from "@bindings/settings";
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
  setConfig: Updater<Settings | undefined>;
}

const SettingsContext = createContext<SettingsContextType | null>(null);

export function SettingsProvider({ children }: { children: ReactNode }) {
  const [config, setConfig] = useMutative<Settings | undefined>(undefined);
  const skipNextSave = useRef(false);

  useEffect(() => {
    SettingsService.GetSettings()
      .then((config) => {
        skipNextSave.current = true;
        setConfig(() => config && rawReturn(structuredClone(config)));
      })
      .catch((err) => logger.error("Failed to load settings: ", err));
  }, [setConfig]);

  useEffect(() => {
    if (!config) return;

    if (skipNextSave.current) {
      skipNextSave.current = false;
      return;
    }

    const timer = setTimeout(
      () =>
        SettingsService.SaveSettings(config)
          .then(() => logger.info("Settings saved successfully"))
          .catch((err) => {
            logger.error("Failed to save settings: ", err);
            toast.error("Failed to save settings", {
              description: (err as Error).message,
            });
          }),
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
    <SettingsContext.Provider value={{ config, setConfig }}>
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
