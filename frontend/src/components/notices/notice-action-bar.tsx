import type { Notice } from "@bindings/notice";
import { CircleCheckBigIcon, CircleIcon, PinIcon, PinOffIcon, XIcon } from "lucide-react";

import { useNoticeMutations } from "@/hooks/use-notice-mutation";
import { useNoticeSelection } from "@/hooks/use-notice-selection";
import { cn } from "@/lib/utils";

import { AppTooltip } from "../app-tooltip";
import { Button } from "../ui/button";
import { Checkbox } from "../ui/checkbox";
import { Label } from "../ui/label";
import { Separator } from "../ui/separator";

function deriveNext(
  checkedIds: Set<string>,
  notices: Notice[],
  field: "isRead" | "isPinned",
): boolean {
  const selected = notices.filter((notices) => checkedIds.has(notices.id));
  return !selected.every((notice) => notice[field]);
}

interface NoticeActionBarProps {
  notices: Notice[];
}

export function NoticeActionBar({ notices }: NoticeActionBarProps) {
  const { selectionMode, setSelectionMode, checkedIds, selectAll, clearChecked } =
    useNoticeSelection();
  const { bulkToggleRead, bulkTogglePin } = useNoticeMutations();

  const allIds = notices.map((notice) => notice.id);
  const count = checkedIds.size;
  const totalCount = notices.length;
  const allSelected = count === totalCount && totalCount > 0;
  const someSelected = count > 0 && !allSelected;

  const readNext = deriveNext(checkedIds, notices, "isRead");
  const pinNext = deriveNext(checkedIds, notices, "isPinned");

  type BulkAction = typeof bulkToggleRead;

  const handleBulkToggle = (action: BulkAction, next: boolean) => {
    if (count === 0) return;
    action({ ids: [...checkedIds], next });
  };

  return (
    <div className="absolute inset-x-0 bottom-0 z-10 flex justify-center">
      <div
        className={cn(
          "mx-3 mb-3 flex w-fit items-center justify-center gap-2 rounded-xl",
          "border bg-background/80 px-3 py-2 shadow-lg backdrop-blur-md",
          "transition-all duration-200 ease-in-out",
          selectionMode
            ? "translate-y-0 opacity-100"
            : "pointer-events-none translate-y-4 opacity-0",
        )}
        aria-hidden={!selectionMode}
      >
        <div className="flex items-center gap-2">
          <Checkbox
            id="select-all"
            indeterminate={someSelected}
            checked={allSelected}
            onCheckedChange={(checked) => (checked ? selectAll(allIds) : clearChecked())}
          />
          <Label
            htmlFor="select-all"
            className="min-w-0 cursor-pointer text-xs text-nowrap text-muted-foreground"
          >
            {count > 0 ? `${String(count)} of ${String(totalCount)}` : "Select all"}
          </Label>
        </div>

        <Separator orientation="vertical" />

        <AppTooltip content={`Mark as ${readNext ? "read" : "unread"}`}>
          <Button
            variant="ghost"
            size="xs"
            disabled={count === 0}
            onClick={() => handleBulkToggle(bulkToggleRead, readNext)}
            className="dark:hover:bg-primary/20"
          >
            {readNext ? (
              <>
                <CircleCheckBigIcon className="size-3.5" />
                Read
              </>
            ) : (
              <>
                <CircleIcon className="size-3.5" />
                Unread
              </>
            )}
          </Button>
        </AppTooltip>

        <Separator orientation="vertical" />

        <AppTooltip content={`Mark as ${pinNext ? "pinned" : "unpinned"}`}>
          <Button
            variant="ghost"
            size="xs"
            disabled={count === 0}
            onClick={() => handleBulkToggle(bulkTogglePin, pinNext)}
            className="dark:hover:bg-primary/20"
          >
            {pinNext ? (
              <>
                <PinIcon className="size-3.5" />
                Pin
              </>
            ) : (
              <>
                <PinOffIcon className="size-3.5" />
                Unpin
              </>
            )}
          </Button>
        </AppTooltip>

        <Separator orientation="vertical" />

        <AppTooltip content="Exit selection">
          <Button
            variant="destructive"
            size="icon-xs"
            onClick={() => setSelectionMode(false)}
            className="bg-transparent dark:bg-transparent"
          >
            <XIcon className="size-4" />
          </Button>
        </AppTooltip>
      </div>
    </div>
  );
}
