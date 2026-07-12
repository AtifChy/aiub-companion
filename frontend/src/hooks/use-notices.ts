import { useSettings } from "@/components/providers/settings-provider";
import { useDebounce } from "@/hooks/use-debounce";
import { logger } from "@/lib/logger";
import { Notice, Service as NoticeService } from "@bindings/notice";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Events } from "@wailsio/runtime";
import { useEffect } from "react";
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

export interface NoticeFilters {
  search: string;
  category: Category | "all";
  urgent: boolean;
  pinned: boolean;
  unread: boolean;
}

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

  const detailOptions = (id: string) => ({
    queryKey: ["noticeDetails", id],
    queryFn: () => NoticeService.GetNoticeDetails(id),
  });

  const detailQuery = useQuery({
    ...detailOptions(selectedId ?? ""),
    enabled: !!selectedId,
  });

  useEffect(() => {
    const unsubscribe = Events.On("notices:synced", (count) => {
      void queryClient.invalidateQueries({ queryKey: ["notices"] });
      toast.success(`${count.data} new notice${count.data !== 1 ? "s" : ""} synced`);
    });
    return unsubscribe;
  }, [queryClient]);

  const toggleRead = useToggleField("isRead", NoticeService.ToggleNoticeRead);
  const togglePin = useToggleField("isPinned", NoticeService.ToggleNoticePinned);

  const readNotice = (id: string) => {
    void queryClient.fetchQuery(detailOptions(id)).then((detail) => {
      if (!detail?.isRead) {
        toggleRead.mutate({ id, next: true });
      }
      return detail;
    });
  };

  return {
    listQuery,
    detailQuery,
    toggleRead,
    togglePin,
    readNotice,
  };
}

function useToggleField(
  field: "isRead" | "isPinned",
  mutationFn: (id: string, next: boolean) => Promise<void>,
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, next }: { id: string; next: boolean }) => mutationFn(id, next),

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
          old?.map((n) => (n.id === id ? Object.assign(new Notice(), n, { [field]: next }) : n)) ??
          old,
      );
      queryClient.setQueriesData<Notice>({ queryKey: ["noticeDetails", id] }, (old) =>
        old ? Object.assign(new Notice(), old, { [field]: next }) : old,
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
