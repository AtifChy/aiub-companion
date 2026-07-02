import { logger } from "@/lib/logger";
import { Service as ConfigService, type Config } from "@bindings/config";
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
  config: Config;
  setConfig: Updater<Config | undefined>;
  resetConfig: () => Promise<void>;
}

const SettingsContext = createContext<SettingsContextType | null>(null);

export function SettingsProvider({ children }: { children: ReactNode }) {
  const [config, setConfig] = useMutative<Config | undefined>(undefined);
  const skipNextSave = useRef(false);

  useEffect(() => {
    ConfigService.GetConfig()
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
  }, [config]);

  const resetConfig = async () => {
    try {
      await ConfigService.ResetConfig();

      const newConfig = await ConfigService.GetConfig();
      if (newConfig) {
        skipNextSave.current = true;
        setConfig(() => rawReturn(structuredClone(newConfig)));
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
