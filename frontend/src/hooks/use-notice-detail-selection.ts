import { Events } from "@wailsio/runtime";
import { useCallback, useEffect, useState } from "react";

import { useNoticeMutations } from "./use-notice-mutation";
import { useNoticeDetail } from "./use-notices";

export function useNoticeDetailSelection() {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const { toggleRead, togglePin } = useNoticeMutations();

  const { query, readNotice } = useNoticeDetail(selectedId, (id: string) =>
    toggleRead({ id, next: true }),
  );

  const select = useCallback(
    (id: string) => {
      setSelectedId((prev) => {
        if (prev === id) return prev;
        readNotice(id);
        return id;
      });
    },
    [readNotice],
  );

  useEffect(() => {
    const unsubscribe = Events.On("notice:open", (event) => select(event.data));
    return unsubscribe;
  }, [select]);

  return {
    selectedId,
    select,
    notice: query.data,
    isLoading: query.isLoading,
    toggleRead,
    togglePin,
  };
}
