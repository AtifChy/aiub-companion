import type { Notice } from "@bindings/notice";
import { CircleIcon, InboxIcon, Loader2Icon, PinIcon, PinOffIcon } from "lucide-react";
import { useShallow } from "zustand/shallow";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { useNoticeMutations } from "@/hooks/use-notice-mutation";
import { useNoticeStore } from "@/hooks/use-notice-store";
import { CATEGORY_STYLES, formatDate, type AltCategory } from "@/lib/notices";
import { cn } from "@/lib/utils";

interface NoticeListProps {
  notices: Notice[];
  loading: boolean;
  error: Error | null;
  onRetry: () => void;
}

export function NoticeList({ notices, loading, error, onRetry }: NoticeListProps) {
  const { filters, clearFilters } = useNoticeStore(
    useShallow((s) => ({ filters: s.filters, clearFilters: s.clearFilters })),
  );
  const hasActiveFilters = filters.urgent || filters.pinned || filters.unread;

  if (loading) {
    return (
      <div className="flex h-32 items-center justify-center">
        <Loader2Icon className="size-5 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center gap-1 px-4 py-16 text-center text-muted-foreground">
        <InboxIcon className="size-8 opacity-20" />
        <p className="text-xs font-semibold text-destructive">Failed to load notices</p>
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
      <div className="flex flex-col items-center gap-2 px-4 py-16 text-center text-muted-foreground">
        <InboxIcon className="size-9 opacity-20" />
        <p className="text-xs">No notices found</p>
        {hasActiveFilters && (
          <Button
            variant="ghost"
            size="xs"
            onClick={clearFilters}
            className="text-primary underline-offset-2 hover:underline"
          >
            Clear filters
          </Button>
        )}
      </div>
    );
  }

  return (
    <div className="flex-1 scrollbar-thin scrollbar-thumb-accent overflow-y-auto">
      {notices.map((notice) => (
        <NoticeListItem key={notice.id} notice={notice} />
      ))}
    </div>
  );
}

interface NoticeListItemProps {
  notice: Notice;
}

function NoticeListItem({ notice }: NoticeListItemProps) {
  const { selected, checked, selectionMode, setSelectedId, toggleChecked } = useNoticeStore(
    useShallow((s) => ({
      selected: s.selectedId === notice.id,
      checked: s.checkedIds.has(notice.id),
      selectionMode: s.selectionMode,
      setSelectedId: s.setSelectedId,
      toggleChecked: s.toggleChecked,
    })),
  );

  const { togglePin } = useNoticeMutations();

  return (
    <div
      onClick={() => setSelectedId(notice.id)}
      onKeyDown={(e) => {
        if (e.key === "Enter") {
          e.preventDefault();
          setSelectedId(notice.id);
        }
      }}
      tabIndex={0}
      className={cn(
        "flex items-start",
        "group w-full animate-in cursor-pointer border-b p-3 text-left transition-colors duration-300 slide-in-from-bottom-5 fade-in",
        "outline-none focus-visible:rounded focus-visible:border-ring focus-visible:border-b-transparent focus-visible:border-l-transparent focus-visible:ring-3 focus-visible:ring-ring/50 focus-visible:ring-inset",
        selected
          ? "border-l-2 border-l-primary bg-primary/10"
          : "border-l-2 border-l-transparent hover:bg-muted/50",
      )}
    >
      {selectionMode && (
        <Checkbox
          checked={checked}
          onClick={(e) => e.preventDefault()}
          onCheckedChange={() => toggleChecked(notice.id)}
          className="mt-0.5 mr-2.5 size-3.5"
        />
      )}

      <div className="min-w-0 flex-1">
        <div className="flex items-center justify-between gap-2">
          <div className="mb-1 flex items-center gap-1.5">
            {!notice.isRead && <CircleIcon className="size-1.5 fill-primary text-primary" />}
            <Badge
              className={cn(
                "h-4 px-1.5 py-0 text-[0.6rem] font-semibold tracking-wider uppercase",
                CATEGORY_STYLES[notice.category as AltCategory],
              )}
            >
              {notice.category}
            </Badge>
            {notice.isUrgent && (
              <Badge
                variant="destructive"
                className="h-4 border border-destructive/20 px-1.5 py-0 text-[0.6rem] font-semibold tracking-wider uppercase"
              >
                urgent
              </Badge>
            )}
            <Button
              variant="ghost"
              onClick={() => togglePin({ id: notice.id, next: !notice.isPinned })}
              tabIndex={notice.isPinned ? 0 : -1}
              className={cn(
                "h-auto p-0.5 opacity-0 group-hover:opacity-100",
                notice.isPinned
                  ? "text-chart-4 opacity-100 hover:text-chart-4/50"
                  : "text-muted-foreground hover:text-foreground",
              )}
            >
              {notice.isPinned ? (
                <PinIcon className="size-3 fill-foreground/20" />
              ) : (
                <PinOffIcon className="size-3 fill-foreground/20" />
              )}
            </Button>
          </div>

          <span className="text-[0.65rem] text-muted-foreground">
            {formatDate(notice.postedDate)}
          </span>
        </div>

        <p
          className={cn(
            "line-clamp-1 text-sm leading-snug",
            notice.isRead ? "font-medium text-muted-foreground" : "font-semibold text-foreground",
          )}
        >
          {notice.fullTitle || notice.title}
        </p>

        {notice.summary && (
          <p className="line-clamp-2 text-[0.7rem] text-muted-foreground">{notice.summary}</p>
        )}
      </div>
    </div>
  );
}
