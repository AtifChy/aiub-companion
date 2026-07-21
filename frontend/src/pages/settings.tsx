import { Service as LogService } from "@bindings/log";
import { System } from "@wailsio/runtime";
import { useTheme } from "next-themes";
import { toast } from "sonner";

import { AlertDialogDestructive } from "@/components/alert-dialog-destructive";
import { useSettings } from "@/components/providers/settings-provider";
import { type Theme } from "@/components/providers/theme-provider";
import { useUpdate } from "@/components/providers/update-provider";
import { SettingsCard } from "@/components/settings/settings-card";
import { SettingRow } from "@/components/settings/settings-row";
import { SettingSelect } from "@/components/settings/settings-select";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";

type Items<T extends string | number> = { value: T; label: string }[];

const themeItems: Items<Theme> = ["light", "dark", "system"].map((v) => ({
  value: v as Theme,
  label: v.charAt(0).toUpperCase() + v.slice(1),
}));

const logLevelItems: Items<string> = ["DEBUG", "INFO", "WARNING", "ERROR"].map((v) => ({
  value: v,
  label: v.charAt(0) + v.slice(1).toLowerCase(),
}));

const syncIntervalItems: Items<number> = [30, 60, 120, 180, 360].map((v) => ({
  value: v,
  label: v >= 60 ? `${String(v / 60)} hour${v === 60 ? "" : "s"}` : `${String(v)} minutes`,
}));

const updateIntervalItems: Items<string> = ["daily", "weekly", "monthly", "never"].map((v) => ({
  value: v,
  label: v.charAt(0).toUpperCase() + v.slice(1),
}));

const fetchCountItems: Items<number> = [10, 20, 30, 50].map((v) => ({
  value: v,
  label: `${String(v)} notices`,
}));

export default function SettingsPage() {
  const { config, setConfig, resetConfig } = useSettings();
  const { setTheme } = useTheme();
  const { check } = useUpdate();

  return (
    <div className="flex h-full animate-in flex-col duration-200 fade-in-10">
      <div className="mr-0.5 min-h-0 flex-1 scrollbar-thin scrollbar-thumb-accent scrollbar-gutter-both space-y-8 overflow-y-auto p-6 lg:p-10">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Settings</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Manage your app preferences and configurations
          </p>
        </div>

        <div className="grid gap-6">
          {/*Appearance*/}
          <SettingsCard title="Appearance">
            <SettingRow label="Theme" description="Select your prefered color theme">
              <SettingSelect
                items={themeItems}
                value={config.appearance.theme}
                onValueChange={(v) => {
                  setConfig((draft) => {
                    draft.appearance.theme = v;
                    setTheme(v);
                  });
                }}
              />
            </SettingRow>
          </SettingsCard>

          {/*Notifications*/}
          <SettingsCard title="Notifications">
            <SettingRow label="Enable" description="Show desktop alert for new notices">
              <Switch
                checked={config.notifications.enabled}
                onCheckedChange={(v) => {
                  setConfig((draft) => (draft.notifications.enabled = v));
                }}
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
                value={config.sync.interval}
                onValueChange={(v) => {
                  setConfig((draft) => (draft.sync.interval = v));
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
                  setConfig((draft) => (draft.sync.fetch_count = v));
                }}
              />
            </SettingRow>

            <SettingRow
              label="Sync on Startup"
              description="Automatically sync notices when app starts"
            >
              <Switch
                checked={config.sync.on_startup}
                onCheckedChange={(v) => {
                  setConfig((draft) => (draft.sync.on_startup = v));
                }}
                className="cursor-pointer"
              />
            </SettingRow>
          </SettingsCard>

          {/*Launch & Systray*/}
          <SettingsCard title="Launch & Systray">
            <SettingRow label="Auto Start" description="Launch the application on system startup">
              <Switch
                checked={config.launch.auto_start}
                onCheckedChange={(v) => {
                  setConfig((draft) => (draft.launch.auto_start = v));
                }}
                className="cursor-pointer"
              />
            </SettingRow>

            <SettingRow
              label="Start Minimized"
              description="Launch the application directly to the system tray"
            >
              <Switch
                checked={config.launch.start_minimized}
                onCheckedChange={(v) => {
                  setConfig((draft) => (draft.launch.start_minimized = v));
                }}
                className="cursor-pointer"
              />
            </SettingRow>

            <SettingRow
              label="Close to Tray"
              description="Minimize the application to the system tray instead of closing it"
            >
              <Switch
                checked={config.launch.close_to_tray}
                onCheckedChange={(v) => {
                  setConfig((draft) => (draft.launch.close_to_tray = v));
                }}
                className="cursor-pointer"
              />
            </SettingRow>

            <SettingRow
              label="Keep Alive"
              description="Continue running the application in the background even when the window is closed"
            >
              <Switch
                disabled={!config.launch.close_to_tray}
                checked={config.launch.keep_alive}
                onCheckedChange={(v) => {
                  setConfig((draft) => (draft.launch.keep_alive = v));
                }}
                className="cursor-pointer"
              />
            </SettingRow>

            <SettingRow
              label="Restore Window"
              description="Restore window size and position on app start"
            >
              <Switch
                checked={config.launch.restore_window}
                onCheckedChange={(v) => {
                  setConfig((draft) => (draft.launch.restore_window = v));
                }}
                className="cursor-pointer"
              />
            </SettingRow>

            <SettingRow label="Sidebar Open" description="Keep the sidebar open on app start">
              <Switch
                checked={config.launch.sidebar_open}
                onCheckedChange={(v) => {
                  setConfig((draft) => (draft.launch.sidebar_open = v));
                }}
                className="cursor-pointer"
              />
            </SettingRow>
          </SettingsCard>

          {System.IsWindows() && (
            <SettingsCard title="Updates">
              <SettingRow label="Update Interval" description="How often to check for updates">
                <SettingSelect
                  items={updateIntervalItems}
                  value={config.updates.interval}
                  onValueChange={(v) => {
                    setConfig((draft) => (draft.updates.interval = v));
                  }}
                />
              </SettingRow>

              <SettingRow label="Check for Updates" description="Manually check for updates">
                <Button variant="outline" disabled={check.isPending} onClick={() => check.mutate()}>
                  {check.isPending ? "Checking..." : "Check Now"}
                </Button>
              </SettingRow>
            </SettingsCard>
          )}

          {/*Advanced*/}
          <SettingsCard title="Advanced">
            <SettingRow
              label="Log Level"
              description="Amount of information logged by the application"
            >
              <SettingSelect
                items={logLevelItems}
                value={config.logging.level}
                onValueChange={(v) => {
                  setConfig((draft) => (draft.logging.level = v));
                }}
              />
            </SettingRow>

            <SettingRow
              label="Open Log File"
              description="Open the log file in default text editor"
            >
              <Button
                variant="outline"
                onClick={() => {
                  LogService.OpenLogFile().catch((err: unknown) => {
                    toast.error("Failed to open log file", {
                      description: (err as Error).message,
                    });
                  });
                }}
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
                onClick={() => void resetConfig()}
              />
            </SettingRow>
          </SettingsCard>
        </div>
      </div>
    </div>
  );
}
