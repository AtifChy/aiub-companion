import { logger } from "@/lib/logger";
import { Button } from "@base-ui/react";
import { Service as ConfigService, type Config } from "@bindings/config";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Loader2Icon } from "lucide-react";
import { rawReturn } from "mutative";
import { createContext, use, useEffect, useRef, type ReactNode } from "react";
import { toast } from "sonner";
import { useMutative, type Updater } from "use-mutative";

interface SettingsContextType {
  config: Config;
  setConfig: Updater<Config | undefined>;
  resetConfig: () => Promise<void>;
}

const SettingsContext = createContext<SettingsContextType | null>(null);

export function SettingsProvider({ children }: { children: ReactNode }) {
  const queryClient = useQueryClient();

  const [config, setConfig] = useMutative<Config | undefined>(undefined);
  const initialized = useRef(false);

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["config"],
    queryFn: ConfigService.GetConfig,
    staleTime: Infinity,
    gcTime: Infinity,
    refetchOnWindowFocus: false,
  });

  useEffect(() => {
    if (data) {
      setConfig(() => rawReturn(structuredClone(data)));
    } else if (error) {
      logger.error("Failed to fetch settings", error);
      toast.error("Error fetching settings");
    }
  }, [data, error, setConfig]);

  useEffect(() => {
    if (!config) return;

    if (!initialized.current) {
      initialized.current = true;
      return;
    }

    const timer = setTimeout(() => {
      ConfigService.SaveConfig(config)
        .then(() => {
          queryClient.setQueryData(["config"], config);
          logger.info("Settings saved successfully");
          toast.success("Settings saved successfully");
        })
        .catch((err) => {
          logger.error("Failed to save settings", err);
          toast.error("Error saving settings");
        });
    }, 500);

    return () => clearTimeout(timer);
  }, [config, queryClient]);

  const resetMutation = useMutation({
    mutationFn: ConfigService.ResetConfig,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["config"] });
      toast.success("Settings reset successfully");
    },
    onError: (err) => {
      logger.error("Failed to reset settings", err);
      toast.error("Error resetting settings");
    },
  });

  if (error) {
    return (
      <div className="flex h-screen w-screen flex-col items-center justify-center gap-3">
        <p className="text-sm text-muted-foreground">Failed to fetch settings</p>
        <Button onClick={() => void refetch()}>Retry</Button>
      </div>
    );
  }

  if (isLoading || !config) {
    return (
      <div className="flex h-screen w-screen items-center justify-center">
        <Loader2Icon className="size-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <SettingsContext.Provider value={{ config, setConfig, resetConfig: resetMutation.mutateAsync }}>
      {children}
    </SettingsContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export function useSettings() {
  const context = use(SettingsContext);
  if (!context) throw new Error("useSettings must be used within a SettingsProvider");
  return context;
}
