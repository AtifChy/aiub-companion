import { Notice, Service as NoticeService } from "@bindings/notice";
import { useMutation, useQueryClient, type QueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { useSettings } from "@/components/providers/settings-provider";
import { logger } from "@/lib/logger";

type ToggleableField = "isRead" | "isPinned";

const noticeKeys = {
  all: ["notices"] as const,
  details: (id: string) => ["noticeDetails", id] as const,
};

function patchNoticeCaches(
  queryClient: QueryClient,
  field: ToggleableField,
  ids: string[],
  next: boolean,
) {
  const idSet = new Set(ids);
  const oldList = queryClient.getQueriesData<Notice[]>({ queryKey: noticeKeys.all });
  const oldDetail = ids.map(
    (id) => [id, queryClient.getQueryData<Notice>(noticeKeys.details(id))] as const,
  );

  queryClient.setQueriesData<Notice[]>({ queryKey: noticeKeys.all }, (old) => {
    return old?.map((notice) => {
      if (idSet.has(notice.id)) {
        return Object.assign(new Notice(), notice, { [field]: next });
      }
      return notice;
    });
  });

  for (const id of ids) {
    queryClient.setQueryData<Notice>(noticeKeys.details(id), (old) =>
      old ? Object.assign(new Notice(), old, { [field]: next }) : old,
    );
  }

  return { oldList, oldDetail };
}

function rollbackNoticeCaches(
  queryClient: QueryClient,
  context: ReturnType<typeof patchNoticeCaches>,
) {
  context.oldList.forEach(([key, data]) => {
    queryClient.setQueryData(key, data);
  });
  context.oldDetail.forEach(([id, data]) => {
    if (data !== undefined) queryClient.setQueryData(noticeKeys.details(id), data);
  });
}

function fieldVerb(field: ToggleableField, next: boolean) {
  switch (field) {
    case "isRead":
      return next ? "read" : "unread";
    case "isPinned":
      return next ? "pin" : "unpin";
  }
}

function useFieldMutation(
  field: ToggleableField,
  mutationFn: (id: string, next: boolean) => Promise<void>,
) {
  const queryClient = useQueryClient();

  const toggle = useMutation({
    mutationFn: ({ id, next }: { id: string; next: boolean }) => mutationFn(id, next),

    onMutate: async ({ id, next }) => {
      await queryClient.cancelQueries({ queryKey: noticeKeys.all });
      await queryClient.cancelQueries({ queryKey: noticeKeys.details(id) });

      return patchNoticeCaches(queryClient, field, [id], next);
    },

    onError: (err, _, context) => {
      if (context) rollbackNoticeCaches(queryClient, context);
      toast.error(`Failed to toggle ${field} status`, { description: err.message });
    },

    onSettled: (_data, _err, { id }) => {
      void queryClient.invalidateQueries({ queryKey: noticeKeys.all });
      void queryClient.invalidateQueries({ queryKey: noticeKeys.details(id) });
    },
  });

  const bulkToggle = useMutation({
    mutationFn: async ({ ids, next }: { ids: string[]; next: boolean }) => {
      const result = await Promise.allSettled(ids.map((id) => mutationFn(id, next)));
      const failedIds = ids.filter((_, index) => result[index]?.status === "rejected");
      return { failedIds };
    },

    onMutate: async ({ ids, next }) => {
      await queryClient.cancelQueries({ queryKey: noticeKeys.all });
      await Promise.all(
        ids.map((id) => queryClient.cancelQueries({ queryKey: noticeKeys.details(id) })),
      );
      return patchNoticeCaches(queryClient, field, ids, next);
    },

    onSuccess: ({ failedIds }, { ids, next }) => {
      if (failedIds.length > 0) {
        toast.warning(
          `${String(failedIds.length)} out of ${String(ids.length)} failed to toggle ${field} status`,
        );
      } else {
        toast.success(
          `Marked ${String(ids.length)} notice${ids.length !== 1 ? "s" : ""} as ${fieldVerb(field, next)}`,
        );
      }
    },

    onError: (err, _, context) => {
      if (context) rollbackNoticeCaches(queryClient, context);
      toast.error(`Failed to toggle ${field} status`, { description: err.message });
    },

    onSettled: (_data, _err, { ids }) => {
      void queryClient.invalidateQueries({ queryKey: noticeKeys.all });
      void Promise.all(
        ids.map((id) => queryClient.invalidateQueries({ queryKey: noticeKeys.details(id) })),
      );
    },
  });

  return { toggle, bulkToggle };
}

export function useNoticeMutations() {
  const { toggle: toggleReadMutation, bulkToggle: bulkToggleReadMutation } = useFieldMutation(
    "isRead",
    NoticeService.ToggleNoticeRead,
  );
  const { toggle: togglePinMutation, bulkToggle: bulkTogglePinMutation } = useFieldMutation(
    "isPinned",
    NoticeService.ToggleNoticePinned,
  );

  return {
    toggleRead: toggleReadMutation.mutate,
    togglePin: togglePinMutation.mutate,
    bulkToggleRead: bulkToggleReadMutation.mutate,
    bulkTogglePin: bulkTogglePinMutation.mutate,
  };
}

export function useSync() {
  const queryClient = useQueryClient();
  const { config } = useSettings();

  const syncMutation = useMutation({
    mutationFn: () => NoticeService.SyncNotices(config.sync.fetch_count),
    onSuccess: (notices) => {
      const count = notices.length;
      if (count > 0) {
        toast.success(`${String(count)} new notice${count !== 1 ? "s" : ""} synced`);
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

  return { syncing: syncMutation.isPending, sync: syncMutation.mutate };
}
