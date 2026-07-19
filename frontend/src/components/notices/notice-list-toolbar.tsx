import { DotIcon, FilterIcon, FilterXIcon, RefreshCwIcon } from "lucide-react";

import { AppTooltip } from "@/components/app-tooltip";
import { HorizontalScroll } from "@/components/horizontal-scroll";
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
import type { NoticeFilters } from "@/hooks/use-notices";
import { CATEGORIES } from "@/lib/notices";
import { cn } from "@/lib/utils";

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

export function NoticeListToolbar({
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
          onValueChange={(v) => {
            onFilterChange({ ...filters, search: v });
          }}
          placeholder="Search notices..."
        />

        <HorizontalScroll className="mt-2 flex gap-1 px-1 outline-none">
          {CATEGORIES.map((cat) => (
            <Badge
              key={cat}
              onClick={() => {
                onFilterChange({ ...filters, category: cat });
              }}
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

      {/*list header*/}
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
          <DropdownMenu>
            <AppTooltip content="Filter">
              <DropdownMenuTrigger
                render={
                  <Button variant="ghost" className="h-auto p-0.5">
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
                    onCheckedChange={(v) => {
                      onFilterChange({ ...filters, [key]: v });
                    }}
                  >
                    {key.charAt(0).toUpperCase() + key.slice(1)}
                  </DropdownMenuCheckboxItem>
                ))}
                <DropdownMenuSeparator />
                <DropdownMenuItem disabled={!hasActiveFilters} onClick={onFilterClear}>
                  Clear
                </DropdownMenuItem>
              </DropdownMenuGroup>
            </DropdownMenuContent>
          </DropdownMenu>

          <AppTooltip content="Sync notices">
            <Button variant="ghost" onClick={onSync} disabled={syncing} className="h-auto p-0.5">
              <RefreshCwIcon className={cn("size-3.5", syncing && "animate-spin")} />
            </Button>
          </AppTooltip>
        </div>
      </div>
    </>
  );
}
