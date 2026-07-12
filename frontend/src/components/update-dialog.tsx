import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogMedia,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { cn } from "@/lib/utils";
import { Events, Updater } from "@wailsio/runtime";
import { DownloadIcon } from "lucide-react";
import { marked } from "marked";
import { useEffect, useState } from "react";
import { useUpdate } from "@/components/providers/update-provider";

export function UpdateDialog() {
  const { release, dialogOpen, setDialogOpen, check, install } = useUpdate();

  if (!release) return null;

  const downloading = install.isPending;

  return (
    <AlertDialog open={dialogOpen} onOpenChange={setDialogOpen}>
      <AlertDialogContent size="sm" className="min-w-md">
        <AlertDialogHeader>
          <AlertDialogMedia className="bg-primary/10 text-primary">
            <DownloadIcon />
          </AlertDialogMedia>
          <AlertDialogTitle>Update Available ({release.version})</AlertDialogTitle>

          <div className="flex flex-col gap-2 text-left">
            <p className="text-sm text-muted-foreground">
              Would you like to download and install it now?
            </p>
          </div>
        </AlertDialogHeader>
        {release.notes && (
          <div className="mt-2 max-h-80 scrollbar-thin scrollbar-thumb-accent overflow-y-auto rounded-md border bg-background/50 p-3 text-sm text-muted-foreground">
            <span className="mb-1 block font-semibold">Release Notes:</span>
            <Markdown>{release.notes}</Markdown>
          </div>
        )}
        <AlertDialogFooter>
          <AlertDialogCancel disabled={downloading} variant="outline">
            Later
          </AlertDialogCancel>
          <DownlaodAction downloading={downloading} onClick={() => check.mutate()} />
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

interface DownlaodActionProps {
  downloading: boolean;
  onClick: () => void;
}

function DownlaodAction({ downloading, onClick }: DownlaodActionProps) {
  const { percent } = useDownloadProgress();
  return (
    <AlertDialogAction
      disabled={downloading}
      variant={downloading ? "outline" : "default"}
      onClick={onClick}
      className="relative overflow-hidden"
    >
      {downloading && (
        <span
          className="absolute inset-y-0 left-0 bg-primary transition-all duration-100"
          style={{ width: `${percent}%` }}
        />
      )}
      <span className="relative z-10">{downloading || "Download & Install"}</span>
    </AlertDialogAction>
  );
}

interface DownloadProgress {
  written: number;
  total: number;
  rate: number;
}

function useDownloadProgress() {
  const [progress, setProgress] = useState<DownloadProgress | null>(null);

  useEffect(() => {
    const unsubscribe = Events.On(Updater.Events.DownloadProgress, (event) => {
      setProgress(event.data);
    });
    return unsubscribe;
  }, []);

  const percent = progress && progress.total > 0 ? (progress.written / progress.total) * 100 : 0;

  return { percent };
}

function Markdown({ children }: { children: string }) {
  const html = marked.parse(children, { async: false });

  return (
    <div
      className={cn(
        "prose prose-sm prose-custom select-text dark:prose-invert",
        "prose-a:no-underline prose-code:rounded-sm prose-code:bg-primary/10 prose-code:px-1 [&_code]:before:content-none [&_code]:after:content-none",
      )}
      // react-doctor-disable-next-line react-doctor/dangerous-html-sink
      dangerouslySetInnerHTML={{ __html: html }}
    />
  );
}
