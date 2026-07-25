import { Notice, Service as NoticeService } from "@bindings/notice";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Events } from "@wailsio/runtime";
import { useEffect } from "react";
import { toast } from "sonner";

import { useDebounce } from "@/hooks/use-debounce";

export const noticeKeys = {
  all: ["notices"] as const,
  details: (id: string) => ["noticeDetails", id] as const,
};

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

export function useNoticeList(filter: NoticeFilters) {
  const queryClient = useQueryClient();
  const debouncedSearch = useDebounce(filter.search, 300);

  const query = useQuery({
    queryKey: [...noticeKeys.all, { ...filter, search: debouncedSearch }],
    queryFn: () =>
      NoticeService.GetNotices({
        Category: filter.category === "all" ? "" : filter.category,
        Search: debouncedSearch,
        Urgent: filter.urgent || null,
        Pinned: filter.pinned || null,
        Unread: filter.unread || null,
      }),
  });

  useEffect(() => {
    const unsubscribe = Events.On("notice:synced", (count) => {
      void queryClient.invalidateQueries({ queryKey: noticeKeys.all });
      toast.success(`${String(count.data)} new notice${count.data !== 1 ? "s" : ""} synced`);
    });
    return unsubscribe;
  }, [queryClient]);

  return { query };
}

export function useNoticeDetail(selectedId: string | null) {
  const queryClient = useQueryClient();

  const query = useQuery({
    queryKey: noticeKeys.details(selectedId ?? ""),
    queryFn: () => NoticeService.GetNoticeDetails(selectedId ?? ""),
    enabled: !!selectedId,
  });

  useEffect(() => {
    const notice = query.data;
    if (!notice) return;

    queryClient.setQueriesData<Notice[]>({ queryKey: noticeKeys.all }, (old) =>
      old?.map((n) =>
        n.id === notice.id && n.fullTitle !== notice.fullTitle
          ? Object.assign(new Notice(), n, {
              fullTitle: notice.fullTitle,
              isCached: notice.isCached,
            })
          : n,
      ),
    );
  }, [query.data, queryClient]);

  return { query };
}
