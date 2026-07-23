import { BuildInfo, Service as MetaService } from "@bindings/meta";
import { Window } from "@wailsio/runtime";
import { CopyrightIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";

import appicon from "@/assets/appicon.png";
import { Separator } from "@/components/ui/separator";
import { WindowControls } from "@/components/window-controls";
import { logger } from "@/lib/logger";

const year = () => new Date().getFullYear();

export default function AboutPage() {
  const [appInfo, setAppInfo] = useState<BuildInfo | null>(null);

  useEffect(() => {
    MetaService.GetBuildInfo()
      .then(setAppInfo)
      .catch((err: unknown) => {
        logger.error("Failed to fetch build info:", err);
        toast.error("Failed to fetch build info. Please try again later.");
      });
  }, []);

  return (
    <div className="flex h-screen w-full flex-col">
      <div className="sticky flex h-10 w-full items-center gap-2 border-b px-1.5 py-2 wails-drag">
        <img src={appicon} alt="App Icon" className="ml-1 size-4" />
        <span className="text-sm">About {appInfo?.name}</span>
        <WindowControls onClose={() => void Window.Close()} />
      </div>

      <div className="flex flex-1 flex-col items-center justify-center p-6 text-center">
        <div className="flex flex-col items-center gap-2">
          <div className="mb-1 size-24 overflow-hidden">
            <img src={appicon} alt="App Icon" />
          </div>

          <h1 className="text-xl font-bold">{appInfo?.name}</h1>

          <Separator className="mt-1 mb-1 bg-transparent bg-linear-to-r from-transparent via-primary/30 to-transparent" />

          <div className="flex gap-2 rounded border border-border bg-muted px-3 py-1 text-sm">
            <span className="font-bold text-primary">{appInfo?.version ?? "N/A"}</span>
          </div>

          <span className="text-xs text-muted-foreground">
            Built on{" "}
            {isNaN(new Date(appInfo?.build_time ?? "").getTime())
              ? "N/A"
              : new Date(appInfo?.build_time ?? "").toLocaleString(undefined, {
                  dateStyle: "medium",
                  timeStyle: "short",
                })}
          </span>

          <div className="mt-4 flex items-center gap-1 text-xs text-muted-foreground">
            <CopyrightIcon className="size-3" />
            {year()} {appInfo?.name}. All rights reserved.
          </div>
        </div>
      </div>
    </div>
  );
}
