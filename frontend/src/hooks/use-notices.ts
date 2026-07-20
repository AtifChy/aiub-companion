import { Service as NoticeService } from "@bindings/notice";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Events } from "@wailsio/runtime";
import { useEffect } from "react";
import { toast } from "sonner";

import { useDebounce } from "@/hooks/use-debounce";

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

  useEffect(() => {
    const unsubscribe = Events.On("notice:synced", (count) => {
      void queryClient.invalidateQueries({ queryKey: ["notices"] });
      toast.success(`${String(count.data)} new notice${count.data !== 1 ? "s" : ""} synced`);
    });
    return unsubscribe;
  }, [queryClient]);

  return { query };
}

export function useNoticeDetail(selectedId: string | null, readOnOpen: (id: string) => void) {
  const options = (id: string) => ({
    queryKey: ["noticeDetails", id],
    queryFn: () => NoticeService.GetNoticeDetails(id),
  });

  const query = useQuery({
    ...options(selectedId ?? ""),
    enabled: !!selectedId,
  });

  const queryClient = useQueryClient();

  const readNotice = (id: string) => {
    void queryClient.fetchQuery(options(id)).then((detail) => {
      if (!detail?.isRead) {
        readOnOpen(id);
      }
      return detail;
    });
  };

  return { query, readNotice };
}
