import type { Notice } from "@bindings/notice";
import { Browser, Events } from "@wailsio/runtime";
import {
  CircleCheckBigIcon,
  CircleIcon,
  DownloadIcon,
  ExternalLinkIcon,
  FileTextIcon,
  InboxIcon,
  Loader2Icon,
  PaperclipIcon,
  PinIcon,
  PinOffIcon,
} from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";

import { AppTooltip } from "@/components/app-tooltip";
import { HorizontalScroll } from "@/components/horizontal-scroll";
import { NoticeActionBar } from "@/components/notices/notice-action-bar";
import { NoticeList } from "@/components/notices/notice-list";
import { NoticeListToolbar } from "@/components/notices/notice-list-toolbar";
import { NoticeSelectionProvider, useNoticeActive } from "@/components/providers/notice-provider";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from "@/components/ui/resizable";
import { Separator } from "@/components/ui/separator";
import { Toggle } from "@/components/ui/toggle";
import { useDelayedLoading } from "@/hooks/use-delayed-loading";
import { NoticeFiltersContext } from "@/hooks/use-notice-filters";
import { useNoticeMutations } from "@/hooks/use-notice-mutation";
import { useNoticeDetail, useNoticeList, type NoticeFilters } from "@/hooks/use-notices";
import { logger } from "@/lib/logger";
import { CATEGORY_STYLES, formatDate, type AltCategory } from "@/lib/notices";
import { cn } from "@/lib/utils";

const INITIAL_FILTERS: NoticeFilters = {
  search: "",
  category: "all",
  urgent: false,
  pinned: false,
  unread: false,
};

export default function NoticesPage() {
  return (
    <NoticeSelectionProvider>
      <NoticesPageInner />
    </NoticeSelectionProvider>
  );
}

function NoticesPageInner() {
  const { setSelectedId } = useNoticeActive();

  useEffect(() => {
    const unsubscribe = Events.On("notice:open", (event) => setSelectedId(event.data));
    return unsubscribe;
  }, [setSelectedId]);

  return (
    <ResizablePanelGroup orientation="horizontal" className="flex h-full">
      <ResizablePanel defaultSize="22rem" minSize="20rem" className="flex flex-col bg-card">
        <NoticeListPanel />
      </ResizablePanel>

      <ResizableHandle withHandle />

      <ResizablePanel defaultSize="60%" minSize="30%" className="min-w-0">
        <DetailPanel />
      </ResizablePanel>
    </ResizablePanelGroup>
  );
}

function NoticeListPanel() {
  const [filters, setFilters] = useState<NoticeFilters>(INITIAL_FILTERS);
  const { query } = useNoticeList(filters);

  const notices = query.data ?? [];
  const unreadCount = notices.filter((notice) => !notice.isRead).length;

  return (
    <div className="relative flex flex-1 flex-col overflow-hidden">
      <NoticeFiltersContext
        value={{
          filters,
          setFilters: (patch) => setFilters({ ...filters, ...patch }),
          clearFilters: () => setFilters(INITIAL_FILTERS),
        }}
      >
        <NoticeListToolbar
          noticeCount={notices.length}
          unreadCount={unreadCount}
          loading={query.isLoading}
        />
        <NoticeList
          notices={notices}
          loading={query.isLoading}
          error={query.error}
          onRetry={() => void query.refetch()}
        />
      </NoticeFiltersContext>
      <NoticeActionBar notices={notices} />
    </div>
  );
}

function DetailPanel() {
  const { selectedId } = useNoticeActive();
  const { query } = useNoticeDetail(selectedId);
  const { toggleRead, togglePin } = useNoticeMutations();

  const showLoader = useDelayedLoading(query.isLoading);

  const notice = query.data ?? null;

  if (query.isLoading) {
    if (!showLoader) return null;
    return (
      <div className="flex h-full items-center justify-center">
        <Loader2Icon className="size-10 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!notice) {
    return (
      <div className="flex h-full flex-col items-center justify-center gap-3 text-muted-foreground">
        <InboxIcon className="size-12 opacity-20" />
        <p className="text-sm">Select a notice to read it</p>
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col overflow-hidden">
      <div className="border-b bg-card px-6 py-4">
        <div className="flex flex-wrap items-center gap-2">
          <Badge
            className={cn(
              "px-2 text-[0.65rem] font-bold tracking-wider uppercase",
              CATEGORY_STYLES[notice.category as AltCategory],
            )}
          >
            {notice.category}
          </Badge>
          {notice.isUrgent && (
            <Badge
              variant="destructive"
              className="border-destructive/20 px-2 text-[0.65rem] font-bold tracking-wider uppercase"
            >
              urgent
            </Badge>
          )}
          {notice.attachments.length > 0 && (
            <Badge variant="outline" className="px-2 text-xs font-semibold">
              <PaperclipIcon />
              {notice.attachments.length}
            </Badge>
          )}
        </div>

        <HorizontalScroll>
          <h1 className="my-3 font-serif text-xl leading-snug font-bold text-nowrap">
            {notice.fullTitle || notice.title}
          </h1>
        </HorizontalScroll>

        <div className="flex flex-wrap items-center justify-between gap-3">
          <span className="text-xs text-muted-foreground">{formatDate(notice.postedDate)}</span>

          <div className="flex items-center">
            <AppTooltip content={notice.isRead ? "Mark as unread" : "Mark as read"}>
              <Toggle
                pressed={notice.isRead}
                onPressedChange={() => toggleRead({ id: notice.id, next: !notice.isRead })}
                size="sm"
                className="aria-pressed:bg-transparent aria-pressed:hover:bg-muted"
              >
                <div className="flex items-center gap-1 group-aria-pressed/toggle:text-chart-3 hover:text-chart-3/50">
                  {notice.isRead ? (
                    <CircleCheckBigIcon className="size-3.5" />
                  ) : (
                    <CircleIcon className="size-3.5" />
                  )}
                </div>
              </Toggle>
            </AppTooltip>
            <AppTooltip content={notice.isPinned ? "Unpin" : "Pin"}>
              <Toggle
                pressed={notice.isPinned}
                onPressedChange={() => togglePin({ id: notice.id, next: !notice.isPinned })}
                size="sm"
                className="aria-pressed:bg-transparent aria-pressed:hover:bg-muted"
              >
                <div className="flex items-center gap-1 *:fill-foreground/20 group-aria-pressed/toggle:text-chart-4 hover:text-chart-4/50">
                  {notice.isPinned ? (
                    <PinOffIcon className="size-3.5" />
                  ) : (
                    <PinIcon className="size-3.5" />
                  )}
                </div>
              </Toggle>
            </AppTooltip>
            <AppTooltip content="Open in browser">
              <Button
                variant="ghost"
                className="hover:text-primary"
                onClick={() =>
                  void Browser.OpenURL("https://aiub.edu/" + notice.id).catch((err: unknown) => {
                    logger.error("Failed to open URL", err);
                    toast.error("Failed to open URL", {
                      description: (err as Error).message,
                    });
                  })
                }
              >
                <ExternalLinkIcon className="size-3.5" />
              </Button>
            </AppTooltip>
          </div>
        </div>
      </div>

      <div className="mr-0.5 scrollbar-thin scrollbar-thumb-accent overflow-y-auto px-6 py-5">
        {notice.content ? (
          <div
            className="mx-auto prose prose-sm text-center prose-custom dark:prose-invert prose-li:text-left"
            // react-doctor-disable-next-line react-doctor/dangerous-html-sink
            dangerouslySetInnerHTML={{ __html: notice.content }}
          />
        ) : (
          <p className="text-sm text-muted-foreground italic">No content available</p>
        )}

        {notice.attachments.length > 0 && (
          <div className="mt-4">
            <Separator className="mb-4" />
            <h3 className="mb-2.5 flex items-center gap-1.5 text-xs tracking-wider uppercase">
              <PaperclipIcon className="size-3.5 text-primary" />
              Attachments
            </h3>
            <div className="flex flex-col gap-1.5">
              {notice.attachments.map((att) => (
                <AttachmentItem key={att.id} attachment={att} />
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

interface AttachmentItemProps {
  attachment: Notice["attachments"][number];
}

function AttachmentItem({ attachment }: AttachmentItemProps) {
  const openUrl = () => {
    if (!attachment.url) return;
    Browser.OpenURL(attachment.url).catch((err: unknown) => {
      toast.error("Failed to open attachment URL", {
        description: (err as Error).message,
      });
    });
  };

  return (
    <div className="flex justify-between rounded border border-border bg-muted/40 px-3 py-2.5 hover:bg-muted">
      <div className="flex min-w-0 items-center gap-2">
        <FileTextIcon className="size-3.5 text-primary" />
        <span className="truncate text-xs font-medium">{attachment.label || "Attachment"}</span>
      </div>
      <div className="ml-3 flex shrink-0 items-center">
        <AppTooltip content={attachment.localPath ? "Open local file" : "Download"}>
          <Button
            variant="ghost"
            className="h-auto p-0.5 hover:text-muted-foreground"
            onClick={openUrl}
          >
            {attachment.localPath ? (
              <ExternalLinkIcon className="size-3.5" />
            ) : (
              <DownloadIcon className="size-3.5" />
            )}
          </Button>
        </AppTooltip>
      </div>
    </div>
  );
}
