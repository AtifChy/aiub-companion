import { logger } from "@/lib/logger";
import { type BuildInfo, Service as ConfigService } from "@bindings/config";
import { useEffect, useState } from "react";
import { toast } from "sonner";

import { Separator } from "@/components/ui/separator";
import { WindowControls } from "@/components/window-controls";
import { Window } from "@wailsio/runtime";
import { CopyrightIcon } from "lucide-react";
import appicon from "../assets/appicon.png";

export default function AboutPage() {
  const [appInfo, setAppInfo] = useState<BuildInfo | null>(null);

  useEffect(() => {
    ConfigService.GetBuildInfo()
      .then(setAppInfo)
      .catch((err) => {
        logger.error("Failed to fetch build info:", err);
        toast.error("Failed to fetch build info. Please try again later.");
      });
  }, []);

  return (
    <div className="flex h-screen w-full flex-col">
      <div className="wails-drag sticky flex h-10 w-full items-center gap-2 border-b px-1.5 py-2">
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

          <Separator className="via-primary/30 mt-1 mb-1 bg-transparent bg-linear-to-r from-transparent to-transparent" />

          <div className="bg-muted border-border flex gap-2 rounded border px-3 py-1 text-sm">
            <span className="text-primary font-bold">
              {appInfo?.version ?? "N/A"}
            </span>
          </div>

          <span className="text-muted-foreground text-xs">
            Built on{" "}
            {isNaN(new Date(appInfo?.build_time ?? "").getTime())
              ? "N/A"
              : new Date(appInfo?.build_time ?? "").toLocaleString()}
          </span>

          <div className="text-muted-foreground mt-4 flex items-center gap-1 text-xs">
            <CopyrightIcon className="size-3" />
            {new Date().getFullYear()} {appInfo?.name}. All rights reserved.
          </div>
        </div>
      </div>
    </div>
  );
}
