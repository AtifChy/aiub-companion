import { Loader2Icon } from "lucide-react";
import { useEffect, type ReactNode } from "react";
import { useShallow } from "zustand/shallow";

import { Button } from "@/components/ui/button";
import { useSettingsStore } from "@/hooks/use-settings-store";

export function SettingsProvider({ children }: { children: ReactNode }) {
  const { status, config, fetchConfig } = useSettingsStore(
    useShallow((s) => ({
      status: s.status,
      config: s.config,
      fetchConfig: s.fetchConfig,
    })),
  );

  useEffect(() => void fetchConfig(), [fetchConfig]);

  if (status === "error") {
    return (
      <div className="flex h-screen w-screen flex-col items-center justify-center gap-3">
        <p className="text-sm text-muted-foreground">Failed to load settings</p>
        <Button onClick={() => void fetchConfig()}>Retry</Button>
      </div>
    );
  }

  if (!config) {
    return (
      <div className="flex h-screen w-screen items-center justify-center">
        <Loader2Icon className="size-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return <>{children}</>;
}

// eslint-disable-next-line react-refresh/only-export-components
export function useSettings() {
  const { config, setConfig, resetConfig } = useSettingsStore(
    useShallow((s) => ({
      config: s.config,
      setConfig: s.updateConfig,
      resetConfig: s.resetConfig,
    })),
  );
  if (!config) {
    throw new Error("useSettings must be used within a SettingsProvider");
  }
  return { config, setConfig, resetConfig };
}
