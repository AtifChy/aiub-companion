import { AlertDialogDestructive } from "@/components/alert-dialog-destructive";
import {
  SettingRow,
  SettingsCard,
  SettingSelect,
} from "@/components/settings/settings";
import { useTheme, type Theme } from "@/components/theme-provider";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import type { Settings } from "@bindings/settings";
import { GetSettings, SaveSettings } from "@bindings/settings/service";
import { Loader2Icon } from "lucide-react";
import { rawReturn } from "mutative";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { useMutative, type Updater } from "use-mutative";

type Items<T extends string | number> = { value: T; label: string }[];

const themeItems: Items<Theme> = ["light", "dark", "system"].map((v) => ({
  value: v as Theme,
  label: v.charAt(0).toUpperCase() + v.slice(1),
}));

const logLevelItems: Items<string> = ["DEBUG", "INFO", "WARNING", "ERROR"].map(
  (v) => ({
    value: v,
    label: v.charAt(0) + v.slice(1).toLowerCase(),
  }),
);

const syncIntervalItems: Items<number> = [30, 60, 120, 180, 360].map((v) => ({
  value: v,
  label: v >= 60 ? `${v / 60} hour${v === 60 ? "" : "s"}` : `${v} minutes`,
}));

const fetchCountItems: Items<number> = [10, 20, 30, 50].map((v) => ({
  value: v,
  label: `${v} notices`,
}));

export default function SettingsPage() {
  const { setTheme } = useTheme();

  const [loading, setLoading] = useState(true);
  const [showLoader, setShowLoader] = useState(false);

  const [config, updateConfig] = useMutative<Settings | undefined>(undefined);

  useEffect(() => {
    const timer = setTimeout(() => setShowLoader(true), 300);

    GetSettings()
      .then((data) => {
        updateConfig(() => rawReturn(JSON.parse(JSON.stringify(data))));
        setLoading(false);
      })
      .catch((err) => {
        toast.error("Failed to load settings", {
          description: err instanceof Error ? err.message : String(err),
        });
        setLoading(false);
      })
      .finally(() => clearTimeout(timer));

    return () => clearTimeout(timer);
  }, [updateConfig]);

  useEffect(() => {
    if (!config) return;

    const timer = setTimeout(async () => {
      try {
        await SaveSettings(config);
        setTheme(config.theme as Theme);
      } catch (err) {
        const message = err instanceof Error ? err.message : String(err);
        toast.error("Failed to save settings", { description: message });
      }
    }, 200);

    return () => clearTimeout(timer);
  }, [config, setTheme]);

  if (loading) {
    if (!showLoader) return null;
    return (
      <div className="flex h-full items-center justify-center">
        <Loader2Icon className="text-muted-foreground size-8 animate-spin" />
      </div>
    );
  }

  if (!config) {
    return (
      <div className="flex h-full items-center justify-center">
        <p className="text-muted-foreground">Failed to load settings data.</p>
      </div>
    );
  }

  return <SettingsView config={config} updateConfig={updateConfig} />;
}

function SettingsView({
  config,
  updateConfig,
}: {
  config: Settings;
  updateConfig: Updater<Settings | undefined>;
}) {
  return (
    <div className="animate-in fade-in-10 flex h-full flex-col duration-200">
      <div className="scrollbar-thumb-accent min-h-0 flex-1 scrollbar-thin scrollbar-gutter-both space-y-8 overflow-y-auto p-6 lg:p-10">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Settings</h1>
          <p className="text-muted-foreground mt-2 text-sm">
            Manage your app preferences and configurations
          </p>
        </div>

        <div className="grid gap-6">
          {/*Appearance*/}
          <SettingsCard title="Appearance">
            <SettingRow
              label="Theme"
              description="Select your prefered color theme"
            >
              <SettingSelect
                items={themeItems}
                value={config.theme}
                onValueChange={(v) => {
                  updateConfig((draft) => {
                    if (draft) draft.theme = v;
                  });
                }}
              />
            </SettingRow>
          </SettingsCard>

          {/*Notifications*/}
          <SettingsCard title="Notifications">
            <SettingRow
              label="Enable"
              description="Show desktop alert for new notices"
            >
              <Switch
                checked={config.notifications.enabled}
                onCheckedChange={(v) =>
                  updateConfig((draft) => {
                    if (draft) draft.notifications.enabled = v;
                  })
                }
                className="cursor-pointer"
              />
            </SettingRow>
          </SettingsCard>

          {/*Sync & Data*/}
          <SettingsCard title="Sync & Data">
            <SettingRow
              label="Sync Interval"
              description="How often to automatically check for new notices"
            >
              <SettingSelect
                items={syncIntervalItems}
                value={config.sync.interval_minutes}
                onValueChange={(v) => {
                  updateConfig((draft) => {
                    if (draft) draft.sync.interval_minutes = v;
                  });
                }}
              />
            </SettingRow>

            <SettingRow
              label="Notice Per Sync"
              description="The number of recent notices to fetch per sync cycle"
            >
              <SettingSelect
                items={fetchCountItems}
                value={config.sync.fetch_count}
                onValueChange={(v) => {
                  updateConfig((draft) => {
                    if (draft) draft.sync.fetch_count = v;
                  });
                }}
              />
            </SettingRow>

            <SettingRow
              label="Sync on Startup"
              description="Automatically sync notices when app starts"
            >
              <Switch
                checked={config.sync.on_startup}
                onCheckedChange={(v) =>
                  updateConfig((draft) => {
                    if (draft) draft.sync.on_startup = v;
                  })
                }
                className="cursor-pointer"
              />
            </SettingRow>
          </SettingsCard>

          {/*Launch & Systray*/}
          <SettingsCard title="Launch & Systray">
            <SettingRow
              label="Start Minimized"
              description="Launch the application directly to the system tray"
            >
              <Switch
                checked={config.launch.start_minimized}
                onCheckedChange={(v) =>
                  updateConfig((draft) => {
                    if (draft) draft.launch.start_minimized = v;
                  })
                }
                className="cursor-pointer"
              />
            </SettingRow>

            <SettingRow
              label="Close to Tray"
              description="Keep the application running in the background when window closed"
            >
              <Switch
                checked={config.launch.close_to_tray}
                onCheckedChange={(v) =>
                  updateConfig((draft) => {
                    if (draft) draft.launch.close_to_tray = v;
                  })
                }
                className="cursor-pointer"
              />
            </SettingRow>

            <SettingRow
              label="Restore Window"
              description="Restore window size and position on app start"
            >
              <Switch
                checked={config.launch.restore_window}
                onCheckedChange={(v) =>
                  updateConfig((draft) => {
                    if (draft) draft.launch.restore_window = v;
                  })
                }
                className="cursor-pointer"
              />
            </SettingRow>
          </SettingsCard>

          {/*Advanced*/}
          <SettingsCard title="Advanced">
            <SettingRow
              label="Log Level"
              description="Ammount of information logged by the application"
            >
              <SettingSelect
                items={logLevelItems}
                value={config.log_level}
                onValueChange={(v) => {
                  updateConfig((draft) => {
                    if (draft) draft.log_level = v;
                  });
                }}
              />
            </SettingRow>

            <SettingRow
              label="Open Log File"
              description="Open the log file in default text editor"
            >
              <Button
                variant="outline"
                onClick={() => toast.error("Not implemented yet")}
              >
                Open Logs
              </Button>
            </SettingRow>
          </SettingsCard>

          {/* Danger Zone */}
          <SettingsCard title="Danger Zone" variant="destructive">
            <SettingRow
              label="Reset All Settings"
              description="Reset all settings to their default values. This action cannot be undone."
            >
              <AlertDialogDestructive
                label="Reset Settings"
                description="Are you sure you want to reset all settings to their default values? This action cannot be undone."
                onClick={() => toast.error("Not implemented yet")}
              />
            </SettingRow>
          </SettingsCard>
        </div>
      </div>
    </div>
  );
}
