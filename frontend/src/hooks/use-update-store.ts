import type { Release } from "@bindings/updater";
import { Events, Updater } from "@wailsio/runtime";
import { create } from "zustand";

export const COUNTDOWN_SECONDS = 10;

interface DownloadProgress {
  written: number;
  total: number;
  rate: number;
}

export type PhaseState =
  | { phase: "idle" }
  | { phase: "downloading"; progress: DownloadProgress | null }
  | { phase: "verifying" }
  | { phase: "ready"; countdown: number }
  | { phase: "installing" };

export type UpdatePhase = PhaseState["phase"];

interface UpdateStore {
  release: Release | null;
  open: boolean;
  phaseState: PhaseState;

  openDialog: (release: Release) => void;
  closeDialog: () => void;
}

export const useUpdateStore = create<UpdateStore>((set) => ({
  release: null,
  open: false,
  phaseState: { phase: "idle" },

  openDialog: (release) => set({ release, open: true, phaseState: { phase: "idle" } }),
  closeDialog: () => set({ open: false }),
}));

export interface UpdateEventHandlers {
  /**
   * Backend reports the update is downloaded and verified,
   * and start the countdown to install the update.
   */
  onReady: () => void;
  /**
   * Any updater error, mid-flow or otherwise abandon whatever was in progress
   */
  onError: () => void;
}

export function initUpdateListeners({ onReady, onError }: UpdateEventHandlers) {
  const unsubs = [
    Events.On(Updater.Events.UpdateAvailable, (event) => {
      const release = event.data;
      if (!release) return;
      useUpdateStore
        .getState()
        .openDialog({ version: release.version, notes: release.notes ?? "" });
    }),
    Events.On(Updater.Events.DownloadStarted, () =>
      useUpdateStore.setState({ phaseState: { phase: "downloading", progress: null } }),
    ),
    Events.On(Updater.Events.DownloadProgress, (event) =>
      useUpdateStore.setState((s) =>
        s.phaseState.phase === "downloading"
          ? { phaseState: { phase: "downloading", progress: event.data } }
          : {},
      ),
    ),
    Events.On(Updater.Events.DownloadComplete, () =>
      useUpdateStore.setState({ phaseState: { phase: "verifying" } }),
    ),
    Events.On(Updater.Events.UpdateReady, () => onReady()),
    Events.On(Updater.Events.Error, () => onError()),
  ];
  return () => unsubs.forEach((unsub) => unsub());
}
