import DOMPurify from "dompurify";
import {
  ArrowDownToLineIcon,
  CircleCheckBigIcon,
  DownloadIcon,
  Loader2Icon,
  ShieldCheckIcon,
} from "lucide-react";
import { marked } from "marked";
import { useShallow } from "zustand/shallow";

import { useUpdate } from "@/components/providers/update-provider";
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
import { Progress } from "@/components/ui/progress";
import { COUNTDOWN_SECONDS, useUpdateStore, type UpdatePhase } from "@/hooks/use-update-store";
import { handleExternalLinkClick } from "@/lib/link-handler";
import { cn, formatBytes, formatSpeed } from "@/lib/utils";

export function UpdateDialog() {
  const { download, runInstall, cancelCountdown } = useUpdate();

  const { release, dialogOpen, phaseState, closeDialog } = useUpdateStore(
    useShallow((s) => ({
      release: s.release,
      dialogOpen: s.open,
      phaseState: s.phaseState,
      closeDialog: s.closeDialog,
    })),
  );

  if (!release) return null;

  const phase = phaseState.phase;
  const progress = phaseState.phase === "downloading" ? phaseState.progress : null;
  const countdown = phaseState.phase === "ready" ? phaseState.countdown : COUNTDOWN_SECONDS;

  const handleDownload = () => {
    useUpdateStore.setState({ phaseState: { phase: "downloading", progress: null } });
    download.mutate();
  };

  const percent =
    progress && progress.total > 0 ? Math.min(100, (progress.written / progress.total) * 100) : 0;

  const isBusy = phase !== "idle" && phase !== "ready";

  return (
    <AlertDialog open={dialogOpen} onOpenChange={isBusy ? undefined : closeDialog}>
      <AlertDialogContent size="sm" className="min-w-md">
        <AlertDialogHeader>
          <AlertDialogMedia
            className={cn(
              "transition-colors duration-300",
              phase === "ready"
                ? "bg-emerald-500/10 text-emerald-500"
                : "bg-primary/10 text-primary",
            )}
          >
            <StateIcon state={phase} />
          </AlertDialogMedia>

          <AlertDialogTitle>
            {phase === "idle" && `Update Available (${release.version})`}
            {phase === "downloading" && `Downloading`}
            {phase === "verifying" && `Verifying`}
            {phase === "ready" && `Update Ready to Install`}
            {phase === "installing" && `Installing`}
          </AlertDialogTitle>

          <div className="flex flex-col gap-2 text-left">
            <p className="w-64 text-center text-sm text-muted-foreground">
              {phase === "idle" &&
                "A new version is available. Would you like to download and install it now?"}
              {phase === "downloading" && "Downloading the latest version…"}
              {phase === "verifying" && "Verifying the downloaded file…"}
              {phase === "ready" && `Download complete! Installing…`}
              {phase === "installing" && "Applying update and restarting the application…"}
            </p>
          </div>
        </AlertDialogHeader>

        {phase === "idle" && release.notes && (
          <div className="mt-2 max-h-80 scrollbar-thin scrollbar-thumb-accent overflow-y-auto rounded-md border bg-background/50 p-3 text-sm text-muted-foreground">
            <span className="mb-1 block font-semibold">Release Notes:</span>
            <Markdown>{release.notes}</Markdown>
          </div>
        )}

        {phase === "downloading" && (
          <div className="mt-2 flex flex-col gap-1.5 rounded-lg border bg-background p-3">
            <div className="flex items-center justify-between text-xs font-medium text-muted-foreground">
              <span>
                {progress
                  ? `${formatBytes(progress.written)} / ${formatBytes(progress.total)}`
                  : "Preparing…"}
              </span>
              <span>
                {progress ? `${percent.toFixed(1)}% · ${formatSpeed(progress.rate)}` : "0%"}
              </span>
            </div>
            <Progress value={percent} className="mt-2 h-2 *:bg-primary/20" />
          </div>
        )}

        {phase === "verifying" && (
          <div className="mt-2 flex items-center gap-2 rounded-lg border bg-background p-3 text-xs font-medium text-muted-foreground">
            <Loader2Icon className="size-4 animate-spin" />
            Verifying integrity…
          </div>
        )}

        {phase === "ready" && (
          <div className="mt-2 rounded-lg bg-background">
            <div className="flex flex-col gap-1.5 rounded-lg border border-emerald-500/20 bg-emerald-500/5 p-3">
              <div className="flex items-center justify-between text-xs font-medium text-emerald-600 dark:text-emerald-400">
                <span>Auto install starts in…</span>
                <span>{countdown}s remaining</span>
              </div>
              <Progress
                value={(countdown / COUNTDOWN_SECONDS) * 100}
                className="mt-2 h-2 *:bg-emerald-500/10 [&>[data-slot=progress-track]>[data-slot=progress-indicator]]:bg-emerald-500"
              />
            </div>
          </div>
        )}

        <AlertDialogFooter className="*:only:col-span-2">
          {phase === "idle" && (
            <>
              <AlertDialogCancel variant="outline">Later</AlertDialogCancel>
              <AlertDialogAction onClick={handleDownload}>Download</AlertDialogAction>
            </>
          )}

          {(phase === "downloading" || phase === "verifying") && (
            <AlertDialogCancel disabled variant="outline">
              {phase === "downloading" ? "Downloading…" : "Verifying…"}
            </AlertDialogCancel>
          )}

          {phase === "ready" && (
            <>
              <AlertDialogCancel variant="outline" onClick={cancelCountdown}>
                Cancel
              </AlertDialogCancel>
              <AlertDialogAction
                onClick={runInstall}
                className="bg-emerald-600 text-background hover:bg-emerald-700"
              >
                Install Now
              </AlertDialogAction>
            </>
          )}

          {phase === "installing" && (
            <AlertDialogAction disabled>
              <Loader2Icon className="mr-2 size-4 animate-spin" />
              Installing…
            </AlertDialogAction>
          )}
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

function StateIcon({ state }: { state: UpdatePhase }) {
  switch (state) {
    case "idle":
      return <DownloadIcon />;
    case "downloading":
      return <ArrowDownToLineIcon className="translate-y-1 animate-bounce" />;
    case "verifying":
      return <ShieldCheckIcon />;
    case "ready":
      return <CircleCheckBigIcon />;
    case "installing":
      return <Loader2Icon className="animate-spin" />;
  }
}

const renderer = new marked.Renderer();
const githubCompareRegex = /^https:\/\/github\.com\/[^/]+\/[^/]+\/compare\/([^/]+)\.\.\.([^/]+)$/;

renderer.link = function ({ href, title, tokens }): string {
  const text = this.parser.parseInline(tokens);

  const match = githubCompareRegex.exec(href);

  if (match && text === href) {
    return `<a href="${href}"${title ? ` title="${title}"` : ""}><code>${String(match[1])}...${String(match[2])}</code></a>`;
  }

  return `<a href="${href}"${title ? ` title="${title}"` : ""}>${text}</a>`;
};

function Markdown({ children }: { children: string }) {
  const html = DOMPurify.sanitize(marked.parse(children, { async: false, renderer }));

  return (
    <div
      role="document"
      className={cn(
        "select-text selection:bg-primary/20",
        "prose prose-sm prose-custom dark:prose-invert",
        "prose-a:no-underline prose-code:rounded-sm prose-code:bg-primary/10 prose-code:px-1 [&_code]:before:content-none [&_code]:after:content-none",
        "prose-blockquote:bg-primary/5 prose-blockquote:py-px prose-blockquote:text-muted-foreground prose-blockquote:not-italic prose-blockquote:[&_p]:before:content-none",
      )}
      onClick={handleExternalLinkClick}
      dangerouslySetInnerHTML={{ __html: html }}
    />
  );
}
