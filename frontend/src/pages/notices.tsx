import { AppTooltip } from "@/components/app-tooltip";
import { HorizontalFadeScroll } from "@/components/horizontal-fade-scroll";
import { SearchInput } from "@/components/search-input";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import { Separator } from "@/components/ui/separator";
import { Toggle } from "@/components/ui/toggle";
import { useDelayedLoading } from "@/hooks/use-delayed-loading";
import {
  useNotices,
  useSync,
  type Category,
  type NoticeFilters,
} from "@/hooks/use-notices";
import { logger } from "@/lib/logger";
import { cn } from "@/lib/utils";
import type { Notice } from "@bindings/notice";
import { Browser } from "@wailsio/runtime";
import {
  CircleCheckBigIcon,
  CircleIcon,
  DotIcon,
  DownloadIcon,
  ExternalLinkIcon,
  FileTextIcon,
  FilterIcon,
  FilterXIcon,
  InboxIcon,
  Loader2Icon,
  PaperclipIcon,
  PinIcon,
  PinOffIcon,
  RefreshCwIcon,
} from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";

const CATEGORIES: Category[] = [
  "all",
  "general",
  "admission",
  "exam",
  "registration",
  "internship",
  "scholarship",
  "payment",
  "holiday",
];

type AltCategory = Exclude<Category, "all">;
const categoryStyles: Record<AltCategory, string> = {
  admission: "border-blue-500/20 bg-blue-500/10 text-blue-500",
  exam: "border-red-500/20 bg-red-500/10 text-red-500",
  registration: "border-violet-500/20 bg-violet-500/10 text-violet-500",
  internship: "border-cyan-500/20 bg-cyan-500/10 text-cyan-500",
  scholarship: "border-yellow-500/20 bg-yellow-500/10 text-yellow-500",
  payment: "border-green-500/20 bg-green-500/10 text-green-500",
  holiday: "border-orange-500/20 bg-orange-500/10 text-orange-500",
  general: "border-zinc-500/20 bg-zinc-500/10 text-zinc-500",
};

const INITIAL_FILTERS: NoticeFilters = {
  search: "",
  category: "all",
  urgent: false,
  pinned: false,
  unread: false,
};

export default function NoticesPage() {
  const [filters, setFilters] = useState<NoticeFilters>(INITIAL_FILTERS);
  const [selectedId, setSelectedId] = useState<string | null>(null);

  const autoReadRef = useRef<string | null>(null);

  const { listQuery, detailQuery, toggleRead, togglePin } = useNotices(
    filters,
    selectedId,
  );
  const { syncing, sync } = useSync();

  const notices = listQuery.data ?? [];
  const unreadCount = notices.filter((notice) => !notice.isRead).length;

  const detail = detailQuery.data ?? null;

  useEffect(() => {
    if (!detail || detailQuery.isLoading) return;

    if (autoReadRef.current !== detail.id) {
      autoReadRef.current = detail.id;

      if (!detail.isRead) {
        toggleRead.mutate({ id: detail.id, next: true });
      }
    }
  }, [detail, detailQuery.isLoading, toggleRead]);

  const handleSelect = (notice: Notice) => {
    if (notice.id === selectedId) return;
    setSelectedId(notice.id);
  };

  return (
    <ResizablePanelGroup orientation="horizontal" className="flex h-full">
      <ResizablePanel
        defaultSize="40%"
        minSize="35%"
        className="bg-card flex flex-col"
      >
        <NoticeListToolbar
          filters={filters}
          onFilterChange={setFilters}
          onFilterClear={() =>
            setFilters({
              ...filters,
              urgent: false,
              pinned: false,
              unread: false,
            })
          }
          syncing={syncing}
          onSync={sync}
          noticeCount={notices.length}
          unreadCount={unreadCount}
          loading={listQuery.isLoading}
        />

        <NoticeList
          notices={notices}
          loading={listQuery.isLoading}
          error={listQuery.error}
          selectedId={selectedId}
          onSelect={handleSelect}
          onTogglePin={(id, next) => togglePin.mutate({ id, next })}
          onRetry={() => void listQuery.refetch()}
          onClearFilters={() => setFilters(INITIAL_FILTERS)}
          hasActiveFilters={filters.urgent || filters.pinned || filters.unread}
        />
      </ResizablePanel>

      <ResizableHandle withHandle />

      <ResizablePanel defaultSize="60%" minSize="30%" className="min-w-0">
        <DetailView
          notice={detail}
          loading={detailQuery.isLoading}
          onTogglePin={() =>
            detail &&
            togglePin.mutate({ id: detail.id, next: !detail.isPinned })
          }
          onToggleRead={() =>
            detail && toggleRead.mutate({ id: detail.id, next: !detail.isRead })
          }
        />
      </ResizablePanel>
    </ResizablePanelGroup>
  );
}

interface NoticeListToolbarProps {
  filters: NoticeFilters;
  onFilterChange: (value: NoticeFilters) => void;
  onFilterClear: () => void;

  syncing: boolean;
  onSync: () => void;

  noticeCount: number;
  unreadCount: number;
  loading: boolean;
}

function NoticeListToolbar({
  filters,
  onFilterChange,
  onFilterClear,
  syncing,
  onSync,
  noticeCount,
  unreadCount,
  loading,
}: NoticeListToolbarProps) {
  const hasActiveFilters = filters.urgent || filters.pinned || filters.unread;

  return (
    <>
      {/*search + categories*/}
      <div className="border-b px-3 py-2.5">
        <SearchInput
          value={filters.search}
          onValueChange={(v) => onFilterChange({ ...filters, search: v })}
          placeholder="Search notices..."
        />

        <HorizontalFadeScroll className="mt-2 flex scrollbar-none gap-1 overflow-x-auto scroll-smooth px-1 outline-none">
          {CATEGORIES.map((cat) => (
            <Badge
              key={cat}
              onClick={() => onFilterChange({ ...filters, category: cat })}
              variant={filters.category === cat ? "default" : "outline"}
              className={cn(
                "cursor-pointer font-semibold capitalize",
                filters.category !== cat &&
                  "text-muted-foreground hover:text-foreground hover:border-foreground/20",
              )}
            >
              {cat}
            </Badge>
          ))}
        </HorizontalFadeScroll>
      </div>

      {/*list header*/}
      <div className="text-muted-foreground flex items-center justify-between border-b px-2 py-1.5 text-[0.7rem]">
        <span>
          {loading ? (
            "Loading…"
          ) : (
            <div className="flex h-0 items-center">
              <span className="text-foreground font-semibold">
                {noticeCount} notices
              </span>
              {unreadCount > 0 && (
                <>
                  <DotIcon />
                  <span className="text-primary font-semibold">
                    {unreadCount} unread
                  </span>
                </>
              )}
            </div>
          )}
        </span>

        <div className="flex gap-1.5">
          <DropdownMenu>
            <AppTooltip content="Filter">
              <DropdownMenuTrigger
                render={
                  <Button variant="ghost" className="h-auto p-0.5">
                    {hasActiveFilters ? (
                      <FilterXIcon className="fill-primary/20 text-primary size-3.5" />
                    ) : (
                      <FilterIcon className="fill-foreground/20 size-3.5" />
                    )}
                  </Button>
                }
              />
            </AppTooltip>
            <DropdownMenuContent>
              <DropdownMenuGroup>
                <DropdownMenuLabel>Filter</DropdownMenuLabel>
                {(["urgent", "pinned", "unread"] as const).map((key) => (
                  <DropdownMenuCheckboxItem
                    key={key}
                    checked={filters[key]}
                    onCheckedChange={(v) =>
                      onFilterChange({ ...filters, [key]: v })
                    }
                  >
                    {key.charAt(0).toUpperCase() + key.slice(1)}
                  </DropdownMenuCheckboxItem>
                ))}
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  disabled={!hasActiveFilters}
                  onClick={onFilterClear}
                >
                  Clear
                </DropdownMenuItem>
              </DropdownMenuGroup>
            </DropdownMenuContent>
          </DropdownMenu>

          <AppTooltip content="Sync notices">
            <Button
              variant="ghost"
              onClick={onSync}
              disabled={syncing}
              className="h-auto p-0.5"
            >
              <RefreshCwIcon
                className={cn("size-3.5", syncing && "animate-spin")}
              />
            </Button>
          </AppTooltip>
        </div>
      </div>
    </>
  );
}

interface NoticeListProps {
  notices: Notice[];
  loading: boolean;
  error: Error | null;
  selectedId: string | null;
  onSelect: (notice: Notice) => void;
  onTogglePin: (id: string, next: boolean) => void;
  onRetry: () => void;
  onClearFilters: () => void;
  hasActiveFilters: boolean;
}

function NoticeList({
  notices,
  loading,
  error,
  selectedId,
  onSelect,
  onTogglePin,
  onRetry,
  onClearFilters,
  hasActiveFilters,
}: NoticeListProps) {
  if (loading) {
    return (
      <div className="flex h-32 items-center justify-center">
        <Loader2Icon className="text-muted-foreground size-5 animate-spin" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-muted-foreground flex flex-col items-center gap-1 px-4 py-16 text-center">
        <InboxIcon className="size-8 opacity-20" />
        <p className="text-destructive text-xs font-semibold">
          Failed to load notices
        </p>
        <p className="text-[0.65rem] break-all">{error.message}</p>
        <Button
          variant="ghost"
          size="xs"
          onClick={onRetry}
          className="text-primary underline-offset-2 hover:underline"
        >
          Retry
        </Button>
      </div>
    );
  }

  if (notices.length === 0) {
    return (
      <div className="text-muted-foreground flex flex-col items-center gap-2 px-4 py-16 text-center">
        <InboxIcon className="size-9 opacity-20" />
        <p className="text-xs">No notices found</p>
        {hasActiveFilters && (
          <Button
            variant="ghost"
            size="xs"
            onClick={onClearFilters}
            className="text-primary underline-offset-2 hover:underline"
          >
            Clear filters
          </Button>
        )}
      </div>
    );
  }

  return (
    <div className="scrollbar-thumb-accent flex-1 scrollbar-thin overflow-y-auto">
      {notices.map((notice) => (
        <NoticeListItem
          key={notice.id}
          notice={notice}
          selected={selectedId === notice.id}
          onSelect={() => onSelect(notice)}
          onTogglePin={(e) => {
            e.stopPropagation();
            onTogglePin(notice.id, !notice.isPinned);
          }}
        />
      ))}
    </div>
  );
}

interface NoticeListItemProps {
  notice: Notice;
  selected: boolean;
  onSelect: () => void;
  onTogglePin: (e: React.MouseEvent) => void;
}

function NoticeListItem({
  notice,
  selected,
  onSelect,
  onTogglePin,
}: NoticeListItemProps) {
  return (
    <div
      onClick={onSelect}
      className={cn(
        "group animate-in fade-in slide-in-from-bottom-5 w-full cursor-pointer border-b p-3 text-left transition-colors duration-300",
        "focus-visible:border-ring focus-visible:ring-ring/50 outline-none focus-visible:rounded focus-visible:border-b-transparent focus-visible:border-l-transparent focus-visible:ring-3 focus-visible:ring-inset",
        selected
          ? "bg-primary/10 border-l-primary border-l-2"
          : "hover:bg-muted/50 border-l-2 border-l-transparent",
      )}
    >
      <div className="flex items-center justify-between gap-2">
        <div className="mb-1 flex items-center gap-1.5">
          {!notice.isRead && (
            <CircleIcon className="fill-primary text-primary size-1.5" />
          )}
          <Badge
            className={cn(
              "h-4 px-1.5 py-0 text-[0.6rem] font-semibold tracking-wider uppercase",
              categoryStyles[notice.category as AltCategory],
            )}
          >
            {notice.category}
          </Badge>
          {notice.isUrgent && (
            <Badge
              variant="destructive"
              className="border-destructive/20 h-4 border px-1.5 py-0 text-[0.6rem] font-semibold tracking-wider uppercase"
            >
              urgent
            </Badge>
          )}
          <Button
            variant="ghost"
            onClick={onTogglePin}
            className={cn(
              "h-auto p-0.5 opacity-0 group-hover:opacity-100",
              notice.isPinned
                ? "text-chart-4 hover:text-chart-4/50 opacity-100"
                : "text-muted-foreground hover:text-foreground",
            )}
          >
            {notice.isPinned ? (
              <PinIcon className="fill-foreground/20 size-3" />
            ) : (
              <PinOffIcon className="fill-foreground/20 size-3" />
            )}
          </Button>
        </div>

        <span className="text-muted-foreground text-[0.65rem]">
          {formatDate(notice.postedDate)}
        </span>
      </div>

      <p
        className={cn(
          "line-clamp-1 text-sm leading-snug",
          notice.isRead
            ? "text-muted-foreground font-medium"
            : "text-foreground font-semibold",
        )}
      >
        {notice.fullTitle || notice.title}
      </p>

      {notice.summary && (
        <p className="text-muted-foreground line-clamp-2 text-[0.7rem]">
          {notice.summary}
        </p>
      )}
    </div>
  );
}

interface DetailViewProps {
  notice: Notice | null;
  loading: boolean;
  onTogglePin: () => void;
  onToggleRead: () => void;
}

function DetailView({
  notice,
  loading,
  onTogglePin,
  onToggleRead,
}: DetailViewProps) {
  const showLoader = useDelayedLoading(loading);

  if (loading) {
    if (!showLoader) return null;
    return (
      <div className="flex h-full items-center justify-center">
        <Loader2Icon className="text-muted-foreground size-10 animate-spin" />
      </div>
    );
  }

  if (!notice) {
    return (
      <div className="text-muted-foreground flex h-full flex-col items-center justify-center gap-3">
        <InboxIcon className="size-12 opacity-20" />
        <p className="text-sm">Select a notice to read it</p>
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col overflow-hidden">
      <div className="bg-card border-b px-6 py-4">
        <div className="flex flex-wrap items-center gap-2">
          <Badge
            className={cn(
              "px-2 text-[0.65rem] font-bold tracking-wider uppercase",
              categoryStyles[notice.category as AltCategory],
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
          {notice.attachments?.length > 0 && (
            <Badge variant="outline" className="px-2 text-xs font-semibold">
              <PaperclipIcon />
              {notice.attachments.length}
            </Badge>
          )}
        </div>

        <h1 className="my-3 line-clamp-1 font-serif text-xl leading-snug font-bold">
          {notice.fullTitle || notice.title}
        </h1>

        <div className="flex flex-wrap items-center justify-between gap-3">
          <span className="text-muted-foreground text-xs">
            {formatDate(notice.postedDate)}
          </span>

          <div className="flex items-center">
            <AppTooltip
              content={notice.isRead ? "Mark as unread" : "Mark as read"}
            >
              <Toggle
                pressed={notice.isRead}
                onPressedChange={onToggleRead}
                size="sm"
                className="aria-pressed:hover:bg-muted aria-pressed:bg-transparent"
              >
                <div className="group-aria-pressed/toggle:text-chart-3 hover:text-chart-3/50 flex items-center gap-1">
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
                onPressedChange={onTogglePin}
                size="sm"
                className="aria-pressed:hover:bg-muted aria-pressed:bg-transparent"
              >
                <div className="*:fill-foreground/20 group-aria-pressed/toggle:text-chart-4 hover:text-chart-4/50 flex items-center gap-1">
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
                  void Browser.OpenURL("https://aiub.edu/" + notice.id).catch(
                    (err) => {
                      logger.error("Failed to open URL", err);
                      toast.error("Failed to open URL", {
                        description: (err as Error).message,
                      });
                    },
                  )
                }
              >
                <ExternalLinkIcon className="size-3.5" />
              </Button>
            </AppTooltip>
          </div>
        </div>
      </div>

      <div className="scrollbar-thumb-accent m-0.5 scrollbar-thin overflow-y-auto px-6 py-5">
        {notice.content ? (
          <div
            className="prose prose-sm prose-custom dark:prose-invert prose-li:text-left mx-auto text-center"
            dangerouslySetInnerHTML={{ __html: notice.content }}
          />
        ) : (
          <p className="text-muted-foreground text-sm italic">
            No content available
          </p>
        )}

        {notice.attachments?.length > 0 && (
          <div className="mt-4">
            <Separator className="mb-4" />
            <h3 className="mb-2.5 flex items-center gap-1.5 text-xs tracking-wider uppercase">
              <PaperclipIcon className="text-primary size-3.5" />
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
    Browser.OpenURL(attachment.url).catch((err) => {
      toast.error("Failed to open attachment URL", {
        description: (err as Error).message,
      });
    });
  };

  return (
    <div className="bg-muted/40 border-border hover:bg-muted flex justify-between rounded border px-3 py-2.5">
      <div className="flex min-w-0 items-center gap-2">
        <FileTextIcon className="text-primary size-3.5" />
        <span className="truncate text-xs font-medium">
          {attachment.label || "Attachment"}
        </span>
      </div>
      <div className="ml-3 flex shrink-0 items-center">
        <AppTooltip
          content={attachment.localPath ? "Open local file" : "Download"}
        >
          <Button
            variant="ghost"
            className="hover:text-muted-foreground h-auto p-0.5"
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

function formatDate(raw: string): string {
  if (!raw) return "";
  const date = new Date(raw);
  if (isNaN(date.getTime())) return raw;
  return date.toLocaleDateString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
}
