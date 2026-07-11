import { logger } from "@/lib/logger";
import { Release, Service as UpdaterService } from "@bindings/updater";
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { toast } from "sonner";

export function useUpdater() {
  const [release, setRelease] = useState<Release | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);

  const check = useMutation({
    mutationFn: () => UpdaterService.CheckForUpdates(),
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

  const install = useMutation({
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

  return {
    release,
    dialogOpen,
    setDialogOpen,
    check,
    install,
  };
}
