import { logger } from "@/lib/logger";
import { Service as ConfigService, type Config } from "@bindings/config";
import { deepEqual } from "fast-equals";
import { Loader2Icon } from "lucide-react";
import { rawReturn } from "mutative";
import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import { toast } from "sonner";
import { useMutative, type Updater } from "use-mutative";

interface SettingsContextType {
  config: Config;
  setConfig: Updater<Config | undefined>;
  resetConfig: () => Promise<void>;
}

const SettingsContext = createContext<SettingsContextType | null>(null);

export function SettingsProvider({ children }: { children: ReactNode }) {
  const [config, setConfig] = useMutative<Config | undefined>(undefined);
  const [lastSaved, setLastSaved] = useState<Config | null>(null);

  useEffect(() => {
    ConfigService.GetConfig()
      .then((config) => {
        setConfig(() => config && rawReturn(structuredClone(config)));
        setLastSaved(structuredClone(config));
      })
      .catch((err) => logger.error("Failed to load settings: ", err));
  }, [setConfig]);

  useEffect(() => {
    if (!config) return;
    if (deepEqual(config, lastSaved)) return; // Don't save if the config hasn't changed

    const timer = setTimeout(
      () =>
        ConfigService.SaveConfig(config)
          .then(() => {
            logger.info("Settings saved successfully");
            toast.success("Settings saved successfully");
          })
          .catch((err) => {
            logger.error("Failed to save settings: ", err);
            toast.error("Failed to save settings", {
              description: (err as Error).message,
            });
          }),
      500,
    );
    return () => clearTimeout(timer);
  }, [config, lastSaved]);

  const resetConfig = async () => {
    try {
      await ConfigService.ResetConfig();

      const newConfig = await ConfigService.GetConfig();
      if (newConfig) {
        setConfig(() => rawReturn(structuredClone(newConfig)));
        setLastSaved(structuredClone(newConfig));
      }

      logger.info("Settings reset successfully");
      toast.success("Settings reset successfully");
    } catch (err) {
      logger.error("Failed to reset settings: ", err);
      toast.error("Failed to reset settings", {
        description: (err as Error).message,
      });
    }
  };

  if (!config) {
    return (
      <div className="flex h-screen w-screen items-center justify-center">
        <Loader2Icon className="text-muted-foreground size-8 animate-spin" />
      </div>
    );
  }

  return (
    <SettingsContext.Provider value={{ config, setConfig, resetConfig }}>
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
