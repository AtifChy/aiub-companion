import { useSettings } from "@/components/settings-provider";
import { useDebounce } from "@/hooks/use-debounce";
import { logger } from "@/lib/logger";
import { Service as NoticeService, type Notice } from "@bindings/notice";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Events } from "@wailsio/runtime";
import { useCallback, useEffect } from "react";
import { toast } from "sonner";

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
      NoticeService.GetNotices({
        Category: filter.category === "all" ? "" : filter.category,
        Search: debouncedSearch,
        Urgent: filter.urgent || null,
        Pinned: filter.pinned || null,
        Unread: filter.unread || null,
      }),
  });

  const detailQuery = useQuery({
    queryKey: ["noticeDetails", selectedId],
    queryFn: () => NoticeService.GetNoticeDetails(selectedId!),
    enabled: !!selectedId,
  });

  const invalidate = useCallback(
    async (id?: string) => {
      if (!id) {
        await queryClient.invalidateQueries({ queryKey: ["notices"] });
        return;
      }
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ["notices"] }),
        queryClient.invalidateQueries({ queryKey: ["noticeDetails", id] }),
      ]);
    },
    [queryClient],
  );

  useEffect(() => {
    const unsub = Events.On("notices:synced", (count) => {
      void invalidate();
      toast.success(
        `${count.data} new notice${count.data !== 1 ? "s" : ""} synced`,
      );
    });
    return () => unsub();
  }, [invalidate]);

  const toggleRead = useToggleField("isRead", NoticeService.ToggleNoticeRead);
  const togglePin = useToggleField(
    "isPinned",
    NoticeService.ToggleNoticePinned,
  );

  return {
    listQuery,
    detailQuery,
    toggleRead,
    togglePin,
  };
}

function useToggleField(
  field: "isRead" | "isPinned",
  mutationFn: (id: string, next: boolean) => Promise<void>,
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, next }: { id: string; next: boolean }) =>
      mutationFn(id, next),

    onMutate: async ({ id, next }) => {
      await queryClient.cancelQueries({ queryKey: ["notices"] });
      await queryClient.cancelQueries({ queryKey: ["noticeDetails", id] });

      const oldList = queryClient.getQueriesData<Notice[]>({
        queryKey: ["notices"],
      });
      const oldDetail = queryClient.getQueryData<Notice>(["noticeDetails", id]);

      queryClient.setQueriesData<Notice[]>(
        { queryKey: ["notices"] },
        (old) =>
          old?.map((n) => (n.id === id ? { ...n, [field]: next } : n)) ?? old,
      );
      queryClient.setQueriesData<Notice>(
        { queryKey: ["noticeDetails", id] },
        (old) => (old ? { ...old, [field]: next } : old),
      );

      return { oldList, oldDetail };
    },

    onError: (err, { id }, context) => {
      context?.oldList?.forEach(([key, data]) => {
        queryClient.setQueryData(key, data);
      });
      if (context?.oldDetail !== undefined) {
        queryClient.setQueryData(["noticeDetails", id], context.oldDetail);
      }

      logger.error(`Failed to toggle ${field} status: `, err);
      toast.error(`Failed to toggle ${field} status`, {
        description: err.message,
      });
    },

    onSettled: (_data, _err, { id }) => {
      void queryClient.invalidateQueries({ queryKey: ["notices"] });
      void queryClient.invalidateQueries({ queryKey: ["noticeDetails", id] });
    },
  });
}

export function useSync() {
  const queryClient = useQueryClient();
  const { config } = useSettings();

  const syncMutation = useMutation({
    mutationFn: () => NoticeService.SyncNotices(config.sync.fetch_count),
    onSuccess: (count) => {
      if (count > 0) {
        toast.success(`${count} new notice${count !== 1 ? "s" : ""} synced`);
      } else {
        toast.info("No new notices to sync");
      }
      void queryClient.invalidateQueries({ queryKey: ["notices"] });
    },
    onError: (err) => {
      logger.error("Failed to sync notices", err);
      toast.error("Failed to sync notices", {
        description: err.message,
      });
    },
  });

  return { syncing: syncMutation.isPending, sync: () => syncMutation.mutate() };
}
