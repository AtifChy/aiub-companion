import { Service as ConfigService, type Config } from "@bindings/config";
import { type Draft, create as mutativeCreate } from "mutative";
import { toast } from "sonner";
import { create } from "zustand";

interface SettingsState {
  config: Config | null;
  status: "idle" | "loading" | "error";
  saveTimer: ReturnType<typeof setTimeout> | null;

  fetchConfig: () => Promise<void>;
  updateConfig: (recipe: (draft: Draft<Config>) => void) => void;
  resetConfig: () => Promise<void>;
}

export const useSettingsStore = create<SettingsState>((set, get) => ({
  config: null,
  status: "idle",
  saveTimer: null,

  fetchConfig: async () => {
    set({ status: "loading" });
    try {
      const data = await ConfigService.GetConfig();
      set({ config: structuredClone(data), status: "idle" });
    } catch (err) {
      toast.error("Error fetching settings", { description: (err as Error).message });
      set({ status: "error" });
    }
  },

  updateConfig: (recipe) => {
    const { config, saveTimer } = get();
    if (!config) return;

    const next = mutativeCreate(config, recipe);
    if (saveTimer) clearTimeout(saveTimer);

    const timer = setTimeout(() => {
      ConfigService.SaveConfig(next)
        .then(() => toast.success("Settings saved successfully"))
        .catch((err: unknown) =>
          toast.error("Error saving settings", { description: (err as Error).message }),
        );
    }, 500);

    set({ config: next, saveTimer: timer });
  },

  resetConfig: async () => {
    try {
      await ConfigService.ResetConfig();
      await get().fetchConfig();
      toast.success("Settings reset to default");
    } catch (err) {
      toast.error("Error resetting settings", { description: (err as Error).message });
    }
  },
}));
