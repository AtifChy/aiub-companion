import { useTheme, type Theme } from "@/components/theme-provider";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";
import { GetSettings, SaveSettings } from "@bindings/settings/service";
import type { Settings } from "@bindings/settings";
import { Loader2Icon, LucideSave } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";

const themeItems: Record<Theme, string> = {
  light: "Light",
  dark: "Dark",
  system: "System",
};

const logLevelItems: { label: string; value: string }[] = [
  { label: "Debug", value: "DEBUG" },
  { label: "Info", value: "INFO" },
  { label: "Warning", value: "WARNING" },
  { label: "Error", value: "ERROR" },
];

export default function SettingsPage() {
  const { setTheme } = useTheme();
  const [loading, setLoading] = useState(true);
  const [showLoader, setShowLoader] = useState(false);
  const [saving, setSaving] = useState(false);
  const [config, setConfig] = useState<Settings | null>(null);

  useEffect(() => {
    const timer = setTimeout(() => setShowLoader(true), 150);

    GetSettings()
      .then((data) => {
        setConfig(data);
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
  }, []);

  const handleSave = async () => {
    if (!config) return;
    setSaving(true);
    try {
      await SaveSettings(config);
      setTheme(config.theme as Theme);
      toast.success("Settings saved successfully");
    } catch (err) {
      toast.error("Failed to save settings", {
        description: err instanceof Error ? err.message : String(err),
      });
    } finally {
      setSaving(false);
    }
  };

  function updateConfig<K extends keyof Settings>(
    key: K,
    value: Settings[K],
  ): void;
  function updateConfig<
    K extends {
      [P in keyof Settings]: Settings[P] extends object ? P : never;
    }[keyof Settings],
    NK extends keyof Settings[K],
  >(section: K, key: NK, value: Settings[K][NK]): void;

  function updateConfig(
    key: unknown,
    valueOrKey: unknown,
    nestedValue?: unknown,
  ) {
    setConfig((prev) => {
      if (!prev) return prev;
      if (nestedValue !== undefined) {
        return {
          ...prev,
          [key as string]: {
            ...prev[key as string],
            [valueOrKey as string]: nestedValue,
          },
        };
      }
      return { ...prev, [key as string]: valueOrKey };
    });
  }

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

  return (
    <div className="animate-in slide-in-from-bottom-5 fade-in flex h-full flex-col duration-200">
      <div className="scrollbar-thumb-accent min-h-0 flex-1 scrollbar-thin space-y-8 overflow-y-auto p-6 lg:p-10">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Settings</h1>
          <p className="text-muted-foreground mt-2 text-sm">
            Manage your app preferences and configurations
          </p>
        </div>

        <div className="grid gap-6">
          {/*Appearance*/}
          <Card className="shadow-xs">
            <CardHeader>
              <CardTitle>Appearance</CardTitle>
            </CardHeader>
            <Separator />
            <CardContent>
              <div className="flex items-center justify-between gap-4">
                <div className="space-y-1">
                  <Label>Theme</Label>
                  <p className="text-muted-foreground">
                    Select your prefered color theme
                  </p>
                </div>
                <Select
                  items={themeItems}
                  value={config.theme}
                  onValueChange={(v) => v !== null && updateConfig("theme", v)}
                >
                  <SelectTrigger className="w-30">
                    <SelectValue placeholder="Select theme" />
                  </SelectTrigger>
                  <SelectContent className="p-1">
                    {Object.entries(themeItems).map((item) => (
                      <SelectItem key={item[0]} value={item[0]}>
                        {item[1]}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </CardContent>
          </Card>

          {/*Sync & Data*/}
          <Card className="shadow-xs">
            <CardHeader>
              <CardTitle>Sync & Data</CardTitle>
            </CardHeader>
            <Separator />
            <CardContent className="flex flex-col gap-4">
              <div className="flex items-center justify-between gap-4">
                <div className="space-y-1">
                  <Label>Sync Interval</Label>
                  <p className="text-muted-foreground text-sm">
                    How often to automatically check for new notices
                  </p>
                </div>
                <Select
                  value={config.sync.interval_minutes}
                  onValueChange={(v) =>
                    v !== null && updateConfig("sync", "interval_minutes", v)
                  }
                >
                  <SelectTrigger className="w-30">
                    <SelectValue placeholder="Select interval" />
                  </SelectTrigger>
                  <SelectContent className="w-50 p-1">
                    {[15, 30, 60, 180].map((interval) => (
                      <SelectItem key={interval} value={interval}>
                        Every {interval} minutes
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="flex items-center justify-between gap-4">
                <div className="space-y-1">
                  <Label>Notice Per Sync</Label>
                  <p className="text-muted-foreground text-sm">
                    The number of recent notices to fetch per sync cycle
                  </p>
                </div>
                <Select
                  value={config.sync.notice_fetch_count}
                  onValueChange={(v) =>
                    v !== null && updateConfig("sync", "notice_fetch_count", v)
                  }
                >
                  <SelectTrigger className="w-30">
                    <SelectValue placeholder="Select count" />
                  </SelectTrigger>
                  <SelectContent className="w-50 p-1">
                    {[10, 20, 30, 50].map((count) => (
                      <SelectItem key={count} value={count}>
                        {count} notices
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="flex items-center justify-between gap-4">
                <div className="space-y-1">
                  <Label>Sync on Startup</Label>
                  <p className="text-muted-foreground text-sm">
                    Automatically sync notices when app starts
                  </p>
                </div>
                <Switch
                  checked={config.sync.on_startup}
                  onCheckedChange={(v) => updateConfig("sync", "on_startup", v)}
                />
              </div>
            </CardContent>
          </Card>

          {/*Launch & Systray*/}
          <Card className="shadow-xs">
            <CardHeader>
              <CardTitle>Launch & Systray</CardTitle>
            </CardHeader>
            <Separator />
            <CardContent className="flex flex-col gap-4">
              <div className="flex items-center justify-between gap-4">
                <div className="space-y-1">
                  <Label>Start Minimized</Label>
                  <p className="text-muted-foreground text-sm">
                    Launch the application directly to the system tray
                  </p>
                </div>
                <Switch
                  checked={config.launch.start_minimized}
                  onCheckedChange={(v) =>
                    updateConfig("launch", "start_minimized", v)
                  }
                />
              </div>

              <div className="flex items-center justify-between gap-4">
                <div className="space-y-1">
                  <Label>Close to Tray</Label>
                  <p className="text-muted-foreground text-sm">
                    Keep the application running in the background when window
                    closed
                  </p>
                </div>
                <Switch
                  checked={config.launch.close_to_tray}
                  onCheckedChange={(v) =>
                    updateConfig("launch", "close_to_tray", v)
                  }
                />
              </div>
            </CardContent>
          </Card>

          {/*Notifications*/}
          <Card className="shadow-xs">
            <CardHeader>
              <CardTitle>Notifications</CardTitle>
            </CardHeader>
            <Separator />
            <CardContent className="flex flex-col gap-4">
              <div className="flex items-center justify-between gap-4">
                <div className="space-y-1">
                  <Label>Enable</Label>
                  <p className="text-muted-foreground text-sm">
                    Show desktop alert for new notices
                  </p>
                </div>
                <Switch
                  checked={config.notifications.enabled}
                  onCheckedChange={(v) =>
                    updateConfig("notifications", "enabled", v)
                  }
                />
              </div>
            </CardContent>
          </Card>

          {/*Advanced*/}
          <Card className="shadow-xs">
            <CardHeader>
              <CardTitle>Advanced</CardTitle>
            </CardHeader>
            <Separator />
            <CardContent>
              <div className="flex items-center justify-between gap-4">
                <div className="space-y-1">
                  <Label>Log Level</Label>
                  <p className="text-muted-foreground">
                    Ammount of information logged by the application
                  </p>
                </div>
                <Select
                  items={logLevelItems}
                  value={config.log_level}
                  onValueChange={(v) =>
                    v !== null && updateConfig("log_level", v)
                  }
                >
                  <SelectTrigger className="w-30">
                    <SelectValue placeholder="Select theme" />
                  </SelectTrigger>
                  <SelectContent className="p-1">
                    {logLevelItems.map((level) => (
                      <SelectItem key={level.value} value={level.value}>
                        {level.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      <div className="border-t p-4">
        <div className="flex justify-end">
          <Button onClick={handleSave} disabled={saving} size="lg">
            {saving ? <Loader2Icon className="animate-spin" /> : <LucideSave />}
            {saving ? "Saving..." : "Save Changes"}
          </Button>
        </div>
      </div>
    </div>
  );
}
