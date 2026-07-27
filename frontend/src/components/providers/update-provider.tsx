import { Service as UpdaterService } from "@bindings/updater";
import { useMutation } from "@tanstack/react-query";
import { createContext, use, useEffect, useRef } from "react";
import { toast } from "sonner";

import { COUNTDOWN_SECONDS, initUpdateListeners, useUpdateStore } from "@/hooks/use-update-store";
import { logger } from "@/lib/logger";

interface UpdateContextType {
  check: ReturnType<typeof useCheckMutation>;
  download: ReturnType<typeof useDownloadMutation>;
  install: ReturnType<typeof useInstallMutation>;
  runInstall: () => void;
  cancelCountdown: () => void;
}

const UpdateContext = createContext<UpdateContextType | null>(null);

function useCheckMutation() {
  return useMutation({
    mutationFn: UpdaterService.CheckForUpdates,
    onSuccess: (data) => {
      if (data) useUpdateStore.getState().openDialog(data);
      else toast.info("You are using the latest version.");
    },
    onError: (err) => {
      logger.error("Error checking for updates:", err);
      toast.error("Error checking for updates. Please try again later.");
    },
  });
}

function useDownloadMutation() {
  return useMutation({
    mutationFn: UpdaterService.DownloadUpdate,
    onError: (err) => {
      logger.error("Error downloading update:", err);
      toast.error("Error downloading update. Please try again later.");
    },
  });
}

function useInstallMutation() {
  return useMutation({
    mutationFn: UpdaterService.InstallUpdate,
    onError: (err) => {
      logger.error("Error installing update:", err);
      toast.error("Error installing update. Please try again later.");
      useUpdateStore.getState().closeDialog();
    },
  });
}

export function UpdateProvider({ children }: { children: React.ReactNode }) {
  const check = useCheckMutation();
  const download = useDownloadMutation();
  const install = useInstallMutation();

  const countdownTimer = useRef<ReturnType<typeof setInterval>>(null);

  const stopCountdown = () => {
    if (countdownTimer.current) {
      clearInterval(countdownTimer.current);
      countdownTimer.current = null;
    }
  };

  const runInstall = () => {
    stopCountdown();
    useUpdateStore.setState({ phaseState: { phase: "installing" } });
    install.mutate();
  };

  const startCountdown = () => {
    stopCountdown();
    useUpdateStore.setState({ phaseState: { phase: "ready", countdown: COUNTDOWN_SECONDS } });

    countdownTimer.current = setInterval(() => {
      const current = useUpdateStore.getState().phaseState;
      if (current.phase !== "ready") {
        stopCountdown();
        return;
      }
      if (current.countdown <= 1) runInstall();
      else
        useUpdateStore.setState({
          phaseState: { phase: "ready", countdown: current.countdown - 1 },
        });
    }, 1000);
  };

  const cancelCountdown = () => {
    stopCountdown();
    useUpdateStore.setState({ phaseState: { phase: "idle" } });
  };

  // oxlint-disable-next-line react-hooks/exhaustive-deps
  useEffect(() => initUpdateListeners({ onReady: startCountdown, onError: cancelCountdown }), []);

  useEffect(
    () =>
      useUpdateStore.subscribe((state) => {
        if (!state.open || state.phaseState.phase !== "ready") stopCountdown();
      }),
    [],
  );

  const value = () => ({ check, download, install, runInstall, cancelCountdown });

  return <UpdateContext value={value()}>{children}</UpdateContext>;
}

// eslint-disable-next-line react-refresh/only-export-components
export function useUpdate() {
  const context = use(UpdateContext);
  if (!context) throw new Error("useUpdate must be used within an UpdateProvider");
  return context;
}
