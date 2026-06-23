import type { Notice } from "@bindings/notice";
import {
  GetNoticeDetails,
  GetNotices,
  SyncNotices,
  ToggleNoticePinned,
  ToggleNoticeRead,
} from "@bindings/notice/service";
import { GetSettings } from "@bindings/settings/service";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Events } from "@wailsio/runtime";
import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";
import { useDebounce } from "./use-debounce";

export type Category =
  | "all"
  | "general"
  | "admission"
  | "exam"
  | "registration"
  | "internship"
  | "scholarship"
  | "payment"
  | "holiday";

export type NoticeFilters = {
  search: string;
  category: Category | "all";
  urgent: boolean;
  pinned: boolean;
  unread: boolean;
};

export function useNotices(filter: NoticeFilters, selectedId: string | null) {
  const queryClient = useQueryClient();
  const debouncedSearch = useDebounce(filter.search, 300);

  const listQuery = useQuery({
    queryKey: ["notices", { ...filter, search: debouncedSearch }],
    queryFn: () =>
      GetNotices({
        Category: filter.category === "all" ? "" : filter.category,
        Search: debouncedSearch,
        Urgent: filter.urgent || null,
        Pinned: filter.pinned || null,
        Unread: filter.unread || null,
      }),
  });

  const detailQuery = useQuery({
    queryKey: ["noticeDetails", selectedId],
    queryFn: () => GetNoticeDetails(selectedId!),
    enabled: !!selectedId,
  });

  const invalidate = useCallback(
    async (id?: string) => {
      await queryClient.invalidateQueries({ queryKey: ["notices"] });
      if (!id) return;
      await queryClient.invalidateQueries({ queryKey: ["noticeDetails", id] });
    },
    [queryClient],
  );

  useEffect(() => {
    return Events.On("notices:synced", async (count) => {
      void invalidate();
      toast.success(
        `${count.data} new notice${count.data !== 1 ? "s" : ""} synced`,
      );
    });
  }, [invalidate]);

  const toggleRead = async (id: string, next: boolean) => {
    queryClient.setQueriesData<Notice[]>(
      { queryKey: ["notices"] },
      (old) =>
        old?.map((n) => (n.id === id ? { ...n, isRead: next } : n)) ?? old,
    );
    queryClient.setQueriesData<Notice>(
      { queryKey: ["noticeDetails", id] },
      (old) => (old ? { ...old, isRead: next } : old),
    );
    try {
      await ToggleNoticeRead(id, next);
    } catch (err) {
      await invalidate(id);
      toast.error("Failed to toggle read status", {
        description: (err as Error).message,
      });
    }
  };

  const togglePin = async (id: string, next: boolean) => {
    queryClient.setQueriesData<Notice[]>(
      { queryKey: ["notices"] },
      (old) =>
        old?.map((n) => (n.id === id ? { ...n, isPinned: next } : n)) ?? old,
    );
    queryClient.setQueriesData<Notice>(
      { queryKey: ["noticeDetails", id] },
      (old) => (old ? { ...old, isPinned: next } : old),
    );
    try {
      await ToggleNoticePinned(id, next);
    } catch (err) {
      await invalidate(id);
      toast.error("Failed to toggle pin status", {
        description: (err as Error).message,
      });
    }
  };

  return {
    listQuery,
    detailQuery,
    toggleRead,
    togglePin,
  };
}

export function useSync() {
  const queryClient = useQueryClient();
  const [syncing, setSyncing] = useState(false);

  const sync = async () => {
    setSyncing(true);
    try {
      const config = await GetSettings();
      const count = await SyncNotices(config?.sync.fetch_count ?? 20);
      if (count > 0) {
        toast.success(`${count} new notice${count !== 1 ? "s" : ""} synced`);
      } else {
        toast.info("No new notices to sync");
      }
      await queryClient.invalidateQueries({ queryKey: ["notices"] });
    } catch (err) {
      toast.error("Failed to sync notices", {
        description: (err as Error).message,
      });
    } finally {
      setSyncing(false);
    }
  };

  return { syncing, sync };
}
