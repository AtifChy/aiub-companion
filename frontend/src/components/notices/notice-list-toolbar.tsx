import {
  CheckSquareIcon,
  DotIcon,
  FilterIcon,
  FilterXIcon,
  RefreshCwIcon,
  SquareIcon,
} from "lucide-react";

import { AppTooltip } from "@/components/app-tooltip";
import { HorizontalScroll } from "@/components/horizontal-scroll";
import { useNoticeFilters } from "@/components/providers/notice-filters-provider";
import { useNoticeBulk } from "@/components/providers/notice-provider";
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
import { Toggle } from "@/components/ui/toggle";
import { useSync } from "@/hooks/use-notice-mutation";
import { CATEGORIES } from "@/lib/notices";
import { cn } from "@/lib/utils";

interface NoticeListToolbarProps {
  noticeCount: number;
  unreadCount: number;
  loading: boolean;
}

export function NoticeListToolbar({ noticeCount, unreadCount, loading }: NoticeListToolbarProps) {
  return (
    <>
      {/*search + categories*/}
      <NoticeListToolbarSearch />

      {/*list header*/}
      <NoticeListToolbarActions
        noticeCount={noticeCount}
        unreadCount={unreadCount}
        loading={loading}
      />
    </>
  );
}

function NoticeListToolbarSearch() {
  const { filters, setFilters } = useNoticeFilters();

  return (
    <div className="border-b px-3 py-2.5">
      <SearchInput
        value={filters.search}
        onValueChange={(v) => setFilters({ search: v })}
        placeholder="Search notices..."
      />

      <HorizontalScroll className="mt-2 flex gap-1 px-1 outline-none">
        {CATEGORIES.map((cat) => (
          <Badge
            key={cat}
            onClick={() => setFilters({ category: cat })}
            variant={filters.category === cat ? "default" : "outline"}
            className={cn(
              "cursor-pointer font-semibold capitalize",
              filters.category !== cat &&
                "text-muted-foreground hover:border-foreground/20 hover:text-foreground",
            )}
          >
            {cat}
          </Badge>
        ))}
      </HorizontalScroll>
    </div>
  );
}

function NoticeListToolbarActions({ noticeCount, unreadCount, loading }: NoticeListToolbarProps) {
  const { syncing, sync } = useSync();
  const { selectionMode, setSelectionMode } = useNoticeBulk();

  return (
    <div className="flex items-center justify-between border-b px-2 py-1.5 text-[0.7rem] text-muted-foreground">
      <span>
        {loading ? (
          "Loading…"
        ) : (
          <div className="flex h-0 items-center">
            <span className="font-semibold text-foreground">{noticeCount} notices</span>
            {unreadCount > 0 && (
              <>
                <DotIcon />
                <span className="font-semibold text-primary">{unreadCount} unread</span>
              </>
            )}
          </div>
        )}
      </span>

      <div className="flex gap-1.5">
        <AppTooltip content={selectionMode ? "Exit selection" : "Select notices"}>
          <Toggle
            pressed={selectionMode}
            onPressedChange={() => setSelectionMode(!selectionMode)}
            size={null}
            className="h-auto rounded-md p-0.5 data-[state=on]:text-primary"
          >
            {selectionMode ? (
              <CheckSquareIcon className="fill size-3.5 fill-primary/10 text-primary" />
            ) : (
              <SquareIcon className="size-3.5 fill-foreground/10" />
            )}
          </Toggle>
        </AppTooltip>

        <NoticeListToolbarFilters />

        <AppTooltip content="Sync notices">
          <Button
            variant="ghost"
            onClick={() => sync()}
            disabled={syncing}
            className="h-auto rounded-md border-none p-0.5"
          >
            <RefreshCwIcon className={cn("size-3.5", syncing && "animate-spin")} />
          </Button>
        </AppTooltip>
      </div>
    </div>
  );
}

function NoticeListToolbarFilters() {
  const { filters, setFilters, clearFilters } = useNoticeFilters();
  const hasActiveFilters = filters.urgent || filters.pinned || filters.unread;

  return (
    <DropdownMenu>
      <AppTooltip content="Filter">
        <DropdownMenuTrigger
          render={
            <Button variant="ghost" className="h-auto rounded-md border-none p-0.5">
              {hasActiveFilters ? (
                <FilterXIcon className="size-3.5 fill-primary/20 text-primary" />
              ) : (
                <FilterIcon className="size-3.5 fill-foreground/20" />
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
              onCheckedChange={(v) => setFilters({ [key]: v })}
            >
              {key.charAt(0).toUpperCase() + key.slice(1)}
            </DropdownMenuCheckboxItem>
          ))}
          <DropdownMenuSeparator />
          <DropdownMenuItem disabled={!hasActiveFilters} onClick={clearFilters}>
            Clear
          </DropdownMenuItem>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
