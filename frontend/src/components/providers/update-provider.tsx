import { logger } from "@/lib/logger";
import { Release, Service as UpdaterService } from "@bindings/updater";
import { useMutation } from "@tanstack/react-query";
import { Events, Updater } from "@wailsio/runtime";
import { createContext, useContext, useEffect, useState } from "react";
import { toast } from "sonner";

interface UpdateContextType {
  release: Release | null;
  dialogOpen: boolean;
  setDialogOpen: (open: boolean) => void;
  check: ReturnType<typeof useCheckMutation>;
  install: ReturnType<typeof useInstallMutation>;
}

const UpdateContext = createContext<UpdateContextType | null>(null);

function useCheckMutation(
  setRelease: (release: Release | null) => void,
  setDialogOpen: (open: boolean) => void,
) {
  return useMutation({
    mutationFn: UpdaterService.CheckForUpdates,
    onSuccess: (data) => {
      if (data) {
        setRelease(data);
        setDialogOpen(true);
      } else {
        toast.info("You are using the latest version.");
      }
    },
    onError: (err) => {
      logger.error("Error checking for updates:", err);
      toast.error("Error checking for updates. Please try again later.");
    },
  });
}

function useInstallMutation(setDialogOpen: (open: boolean) => void) {
  return useMutation({
    mutationFn: async () => {
      await UpdaterService.DownloadUpdate();
      await UpdaterService.InstallUpdate();
    },
    onError: (err) => {
      logger.error("Error installing update:", err);
      toast.error("Error installing update. Please try again later.");
      setDialogOpen(false);
    },
  });
}

export function UpdateProvider({ children }: { children: React.ReactNode }) {
  const [release, setRelease] = useState<Release | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);

  const check = useCheckMutation(setRelease, setDialogOpen);
  const install = useInstallMutation(setDialogOpen);

  useEffect(() => {
    const unsubscribe = Events.On(Updater.Events.UpdateAvailable, (event) => {
      setRelease(event.data);
      setDialogOpen(true);
    });
    return unsubscribe;
  }, []);

  return (
    <UpdateContext.Provider value={{ release, dialogOpen, setDialogOpen, check, install }}>
      {children}
    </UpdateContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export function useUpdate() {
  const context = useContext(UpdateContext);
  if (!context) {
    throw new Error("useUpdate must be used within an UpdateProvider");
  }
  return context;
}
