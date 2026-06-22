import { AppTooltip } from "@/components/app-tooltip";
import { HorizontalFadeScroll } from "@/components/horizontal-fade-scroll";
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
import { Input } from "@/components/ui/input";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import { Separator } from "@/components/ui/separator";
import { Toggle } from "@/components/ui/toggle";
import { useDebounce } from "@/hooks/use-debounce";
import { cn } from "@/lib/utils";
import { Notice } from "@bindings/notice";
import {
  GetNoticeDetails,
  GetNotices,
  SyncNotices,
  ToggleNoticePinned,
  ToggleNoticeRead,
} from "@bindings/notice/service";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Browser, Events } from "@wailsio/runtime";
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
  SearchIcon,
  SlashIcon,
  XIcon,
} from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";

const PAGE_SIZE = 30;

const CATEGORIES: readonly string[] = [
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

type Category = Exclude<(typeof CATEGORIES)[number], "all">;

const categoryStyles: Record<Category, string> = {
  admission: "border-blue-500/20 bg-blue-500/10 text-blue-500",
  exam: "border-red-500/20 bg-red-500/10 text-red-500",
  registration: "border-violet-500/20 bg-violet-500/10 text-violet-500",
  internship: "border-cyan-500/20 bg-cyan-500/10 text-cyan-500",
  scholarship: "border-yellow-500/20 bg-yellow-500/10 text-yellow-500",
  payment: "border-green-500/20 bg-green-500/10 text-green-500",
  holiday: "border-orange-500/20 bg-orange-500/10 text-orange-500",
  general: "border-zinc-500/20 bg-zinc-500/10 text-zinc-500",
};

export default function NoticesPage() {
  const queryClient = useQueryClient();

  const [search, setSearch] = useState("");
  const [activeCategory, setActiveCategory] = useState<Category | "all">("all");
  const [filter, setFilter] = useState({
    urgent: false,
    pinned: false,
    unread: false,
  });

  const [syncing, setSyncing] = useState(false);
  const [selectedId, setSelectedId] = useState<string | null>(null);

  const debouncedSearch = useDebounce(search, 300);
  const searchRef = useRef<HTMLInputElement>(null);

  const noticesQuery = useQuery({
    queryKey: ["notices", activeCategory, debouncedSearch, filter],
    queryFn: () =>
      GetNotices({
        Category: activeCategory === "all" ? "" : activeCategory,
        Search: debouncedSearch,
        Urgent: filter.urgent || null,
        Pinned: filter.pinned || null,
        Unread: filter.unread || null,
        Limit: PAGE_SIZE,
        Offset: 0,
      }),
  });

  const notices = noticesQuery.data ?? [];
  const loadingList = noticesQuery.isLoading;
  const fetchError = noticesQuery.error;

  const noticeDetailsQuery = useQuery({
    queryKey: ["noticeDetails", selectedId],
    queryFn: () => GetNoticeDetails(selectedId!),
    enabled: !!selectedId,
  });

  const detail = noticeDetailsQuery.data ?? null;
  const loadingDetail = noticeDetailsQuery.isLoading && !detail;

  const unreadCount = notices.filter((n) => !n.isRead).length;

  // search on "/" key
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === "/" && document.activeElement !== searchRef.current) {
        e.preventDefault();
        searchRef.current?.focus();
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, []);

  // listen for notices synced event
  useEffect(() => {
    const cleanup = Events.On("notices:synced", async (count) => {
      toast.success(`${count.data} new notice` + (count.data !== 1 ? "s" : ""));
      queryClient.invalidateQueries({ queryKey: ["notices"] });
    });
    return cleanup;
  }, [queryClient]);

  // select notice
  const handleSelect = async (notice: Notice) => {
    if (notice.id === selectedId) return;

    setSelectedId(notice.id);

    if (notice.isRead) return;

    try {
      await ToggleNoticeRead(notice.id, true);
      queryClient.invalidateQueries({ queryKey: ["notices"] });
    } catch (err) {
      toast.error("Failed to mark notice as read: " + err);
    }
  };

  // toggle read
  const handleToggleRead = async () => {
    if (!detail) return;
    const next = !detail.isRead;
    try {
      await ToggleNoticeRead(detail.id, next);
      queryClient.invalidateQueries({ queryKey: ["notices"] });
    } catch (err) {
      toast.error("Failed to toggle read: " + err);
    }
  };

  // toggle pin
  const handleTogglePin = async (
    id: string,
    current: boolean,
    e?: React.MouseEvent,
  ) => {
    e?.stopPropagation();
    const next = !current;
    try {
      await ToggleNoticePinned(id, next);
      queryClient.invalidateQueries({ queryKey: ["notices"] });
    } catch (err) {
      toast.error("Failed to toggle pin: " + err);
    }
  };

  // sync notices
  const handleSync = async () => {
    setSyncing(true);
    try {
      const count = await SyncNotices(PAGE_SIZE);
      if (count > 0) {
        toast.success(`Synced ${count} notice${count !== 1 ? "s" : ""}.`, {
          position: "top-center",
        });
      } else {
        toast.info("No new notices to sync.");
      }
      await queryClient.invalidateQueries({ queryKey: ["notices"] });
    } catch (err) {
      toast.error("Failed to sync notices: " + err);
    } finally {
      setSyncing(false);
    }
  };

  return (
    <ResizablePanelGroup
      orientation="horizontal"
      className="flex h-full overflow-hidden"
    >
      <ResizablePanel defaultSize="35%" className="bg-card flex flex-col">
        {/*list toolbar*/}
        <div className="border-b px-3 py-2.5">
          {/*search*/}
          <div className="relative">
            <SearchIcon className="text-muted-foreground absolute top-1/2 left-2.5 size-3.5 -translate-y-1/2" />
            <Input
              ref={searchRef}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search notices…"
              autoComplete="off"
              className="pr-7 pl-8 text-xs select-none"
            />
            {search ? (
              <Button
                variant="ghost"
                onClick={() => setSearch("")}
                className="text-muted-foreground hover:text-foreground absolute top-1/2 right-2 h-auto -translate-y-1/2 p-0.5"
              >
                <XIcon strokeWidth="2.5" className="size-3.5" />
              </Button>
            ) : (
              <div className="text-secondary-foreground/40 border-secondary/10 bg-secondary/20 absolute top-1/2 right-2 -translate-y-1/2 rounded-sm border p-1">
                <SlashIcon strokeWidth="4" className="size-2.5" />
              </div>
            )}
          </div>

          {/*category tabs*/}
          <HorizontalFadeScroll className="mt-2 flex scrollbar-none gap-1 overflow-x-auto scroll-smooth px-1 outline-none">
            {CATEGORIES.map((cat) => (
              <Badge
                key={cat}
                onClick={() => setActiveCategory(cat)}
                variant={activeCategory === cat ? "default" : "outline"}
                className={cn(
                  "cursor-pointer font-semibold capitalize",
                  activeCategory !== cat &&
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
            {loadingList ? (
              "Loading…"
            ) : (
              <div className="flex h-0 items-center">
                <span className="text-foreground font-semibold">
                  {notices.length} notices
                </span>
                {unreadCount >= 0 && (
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
                      {!filter.urgent && !filter.pinned && !filter.unread ? (
                        <FilterIcon className="fill-foreground/20 size-3.5" />
                      ) : (
                        <FilterXIcon className="fill-primary/20 text-primary size-3.5" />
                      )}
                    </Button>
                  }
                />
              </AppTooltip>
              <DropdownMenuContent>
                <DropdownMenuGroup>
                  <DropdownMenuLabel>Filter</DropdownMenuLabel>
                  <DropdownMenuCheckboxItem
                    checked={filter.urgent}
                    onCheckedChange={(v) =>
                      setFilter((prev) => ({ ...prev, urgent: v }))
                    }
                  >
                    Urgent
                  </DropdownMenuCheckboxItem>
                  <DropdownMenuCheckboxItem
                    checked={filter.pinned}
                    onCheckedChange={(v) =>
                      setFilter((prev) => ({ ...prev, pinned: v }))
                    }
                  >
                    Pinned
                  </DropdownMenuCheckboxItem>
                  <DropdownMenuCheckboxItem
                    checked={filter.unread}
                    onCheckedChange={(v) =>
                      setFilter((prev) => ({ ...prev, unread: v }))
                    }
                  >
                    Unread
                  </DropdownMenuCheckboxItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    disabled={
                      !filter.urgent && !filter.pinned && !filter.unread
                    }
                    onClick={() =>
                      setFilter({ urgent: false, pinned: false, unread: false })
                    }
                  >
                    Clear
                  </DropdownMenuItem>
                </DropdownMenuGroup>
              </DropdownMenuContent>
            </DropdownMenu>
            <AppTooltip content="Sync notices">
              <Button
                variant="ghost"
                onClick={handleSync}
                disabled={syncing}
                className="h-auto p-0.5"
              >
                <RefreshCwIcon
                  className={cn("size-3.5", syncing ? "animate-spin" : "")}
                />
              </Button>
            </AppTooltip>
          </div>
        </div>

        {/*notice list*/}
        <div className="scrollbar-thumb-accent flex-1 scrollbar-thin overflow-y-auto">
          {loadingList ? (
            <div className="flex h-32 items-center justify-center">
              <Loader2Icon className="text-muted-foreground size-5 animate-spin" />
            </div>
          ) : fetchError ? (
            <div className="text-muted-foreground flex flex-col items-center gap-1 px-4 py-16 text-center">
              <InboxIcon className="size-8 opacity-20" />
              <p className="text-destructive text-xs font-semibold">
                Failed to load notices
              </p>
              <p className="text-[0.65rem] break-all">{fetchError.message}</p>
              <Button
                variant="ghost"
                size="xs"
                onClick={() =>
                  queryClient.invalidateQueries({ queryKey: ["notices"] })
                }
                className="text-primary underline-offset-2 hover:underline"
              >
                Retry
              </Button>
            </div>
          ) : notices.length === 0 ? (
            <div className="text-muted-foreground flex flex-col items-center gap-2 px-4 py-16 text-center">
              <InboxIcon className="size-9 opacity-20" />
              <p className="text-xs">No notices found</p>
              {(search || activeCategory !== "all") && (
                <Button
                  variant="ghost"
                  size="xs"
                  onClick={() => {
                    setSearch("");
                    setActiveCategory("all");
                  }}
                  className="text-primary underline-offset-2 hover:underline"
                >
                  Clear filters
                </Button>
              )}
            </div>
          ) : (
            notices.map((notice) => (
              <NoticeListItem
                key={notice.id}
                notice={notice}
                isActive={selectedId === notice.id}
                onClick={() => handleSelect(notice)}
                onTogglePin={(e) =>
                  handleTogglePin(notice.id, notice.isPinned, e)
                }
              />
            ))
          )}
        </div>
      </ResizablePanel>

      <ResizableHandle withHandle />

      {/*detail panel*/}
      <ResizablePanel defaultSize="65%" minSize="30%" className="min-w-0">
        <DetailPanel
          notice={detail}
          loading={loadingDetail}
          onTogglePin={() =>
            detail && handleTogglePin(detail.id, detail.isPinned)
          }
          onToggleRead={handleToggleRead}
        />
      </ResizablePanel>
    </ResizablePanelGroup>
  );
}

interface NoticeListItemProps {
  notice: Notice;
  isActive: boolean;
  onClick: () => void;
  onTogglePin: (e: React.MouseEvent) => void;
}

function NoticeListItem({
  notice,
  isActive,
  onClick,
  onTogglePin,
}: NoticeListItemProps) {
  return (
    <button
      onClick={onClick}
      className={cn(
        "group animate-in fade-in slide-in-from-bottom-5 w-full cursor-pointer border-b p-3 text-left transition-colors duration-300",
        "focus-visible:border-ring focus-visible:ring-ring/50 outline-none focus-visible:rounded focus-visible:border-b-transparent focus-visible:border-l-transparent focus-visible:ring-3 focus-visible:ring-inset",
        isActive
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
              categoryStyles[notice.category],
            )}
          >
            {notice.category}
          </Badge>
          {notice.isUrgent && (
            <Badge
              variant="destructive"
              className={cn(
                "border-destructive/20 h-4 border px-1.5 py-0 text-[0.6rem] font-semibold tracking-wider uppercase",
              )}
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

        <div className="flex items-center gap-1">
          <span className="text-muted-foreground text-[0.65rem]">
            {formatDate(notice.postedDate)}
          </span>
        </div>
      </div>

      {/*title*/}
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

      {/*summary*/}
      {notice.summary && (
        <p className="text-muted-foreground line-clamp-2 text-[0.7rem]">
          {notice.summary}
        </p>
      )}
    </button>
  );
}

interface DetailPanelProps {
  notice: Notice | null;
  loading: boolean;
  onTogglePin: () => void;
  onToggleRead: () => void;
}

function DetailPanel({
  notice,
  loading,
  onTogglePin,
  onToggleRead,
}: DetailPanelProps) {
  if (loading && !notice?.isCached) {
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
      {/*detail header*/}
      <div className="bg-card border-b px-6 py-4">
        {/*badges*/}
        <div className="flex flex-wrap items-center gap-2">
          <Badge
            className={cn(
              "px-2 text-[0.65rem] font-bold tracking-wider uppercase",
              categoryStyles[notice.category] ?? categoryStyles["default"],
            )}
          >
            {notice.category}
          </Badge>
          {notice.isUrgent && (
            <Badge
              variant="destructive"
              className={cn(
                "border-destructive/20 px-2 text-[0.65rem] font-bold tracking-wider uppercase",
                categoryStyles["urgent"],
              )}
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

        {/*title*/}
        <h1 className="my-3 line-clamp-1 font-serif text-xl leading-snug font-bold">
          {notice.fullTitle || notice.title}
        </h1>

        {/*meta + actions*/}
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
                onClick={() => {
                  Browser.OpenURL("https://aiub.edu/" + notice.id).catch(
                    (err) => toast.error("Failed to open URL: " + err),
                  );
                }}
              >
                <ExternalLinkIcon className="size-3.5" />
              </Button>
            </AppTooltip>
          </div>
        </div>
      </div>

      {/*scrollable body*/}
      <div className="scrollbar-thumb-accent scrollbar-thin overflow-y-auto px-6 py-5">
        {/*content*/}
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

        {/*attachments*/}
        {notice.attachments?.length > 0 && (
          <div className="mt-4">
            <Separator className="mb-4" />
            <h3 className="mb-2.5 flex items-center gap-1.5 text-xs tracking-wider uppercase">
              <PaperclipIcon className="text-primary size-3.5" />
              Attachments
            </h3>
            <div className="flex flex-col gap-1.5">
              {notice.attachments.map((att) => (
                <div
                  key={att.id}
                  className="bg-muted/40 border-border hover:bg-muted flex justify-between rounded border px-3 py-2.5"
                >
                  <div className="flex min-w-0 items-center gap-2">
                    <FileTextIcon className="text-primary size-3.5" />
                    <span className="truncate text-xs font-medium">
                      {att.label || "Attachment"}
                    </span>
                  </div>
                  <div className="ml-3 flex shrink-0 items-center">
                    {att.localPath ? (
                      <AppTooltip content="Open local file">
                        <Button
                          variant="ghost"
                          className="hover:text-muted-foreground h-auto p-0.5"
                          onClick={() => {
                            if (att.url) {
                              Browser.OpenURL(att.url).catch((err) =>
                                toast.error("Failed to open URL: " + err),
                              );
                            }
                          }}
                        >
                          <ExternalLinkIcon className="size-3.5" />
                        </Button>
                      </AppTooltip>
                    ) : (
                      <AppTooltip content="Download">
                        <Button
                          variant="ghost"
                          className="hover:text-muted-foreground h-auto p-0.5"
                          onClick={() => {
                            if (att.url) {
                              Browser.OpenURL(att.url).catch((err) =>
                                toast.error("Failed to open URL: " + err),
                              );
                            }
                          }}
                        >
                          <DownloadIcon className="size-3.5" />
                        </Button>
                      </AppTooltip>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
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
