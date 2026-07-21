import type { Notice } from "@bindings/notice";
import { useQueryClient } from "@tanstack/react-query";
import { createContext, use, useEffect, useRef, useState, type ReactNode } from "react";

import { useNoticeMutations } from "@/hooks/use-notice-mutation";
import { noticeKeys } from "@/hooks/use-notices";

interface NoticeActiveContextType {
  selectedId: string | null;
  setSelectedId: (id: string) => void;
}

const NoticeActiveContext = createContext<NoticeActiveContextType | null>(null);

// oxlint-disable-next-line react/only-export-components
export function useNoticeActive() {
  const context = use(NoticeActiveContext);
  if (!context) {
    throw new Error("useNoticeActive must be used within a NoticeSelectionProvider");
  }
  return context;
}

interface NoticeBulkContextType {
  selectionMode: boolean;
  setSelectionMode: (on: boolean) => void;
  checkedIds: Set<string>;
  toggleChecked: (id: string) => void;
  selectAll: (ids: string[]) => void;
  clearChecked: () => void;
}

const NoticeBulkContext = createContext<NoticeBulkContextType | null>(null);

// oxlint-disable-next-line react/only-export-components
export function useNoticeBulk() {
  const context = use(NoticeBulkContext);
  if (!context) {
    throw new Error("useNoticeBulk must be used within a NoticeSelectionProvider");
  }
  return context;
}

export function NoticeSelectionProvider({ children }: { children: ReactNode }) {
  // Active State
  const [selectedId, setSelectedId] = useState<string | null>(null);

  // Bulk State
  const [selectionMode, setSelectionMode] = useState(false);
  const [checkedIds, setCheckedIds] = useState<Set<string>>(() => new Set());

  // Keep a ref to the selectionMode state to avoid constantly re-rendering the context provider when selectionMode changes. This is useful for performance optimization in components that consume this context
  const selectionModeRef = useRef(selectionMode);
  selectionModeRef.current = selectionMode;

  const queryClient = useQueryClient();
  const { toggleRead } = useNoticeMutations();

  // Track the previous ID so we can mark it read when navigating away
  const previousIdRef = useRef<string | null>(null);

  const toggleChecked = (id: string) => {
    setCheckedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  const selectAll = (ids: string[]) => setCheckedIds(new Set(ids));

  const clearChecked = () => setCheckedIds(new Set());

  const handleSelectionMode = (on: boolean) => {
    setSelectionMode(on);
    if (!on) clearChecked();
  };

  // When a notice is selected, we want to mark the previous notice as read if it was unread.
  const handleSetSelectedId = (id: string) => {
    if (selectionModeRef.current) {
      toggleChecked(id);
      return;
    }

    const prevId = previousIdRef.current;
    if (prevId && prevId !== id) {
      // Mark the previous notice as read if it was unread
      const prevNotice = queryClient.getQueryData<Notice>(noticeKeys.details(prevId));
      if (prevNotice && !prevNotice.isRead) {
        toggleRead({ id: prevId, next: true });
      }
    }

    previousIdRef.current = id;
    setSelectedId(id);
  };

  // Mark the last selected notice as read when the component unmounts
  useEffect(() => {
    return () => {
      const currId = previousIdRef.current;
      if (currId) {
        const currNotice = queryClient.getQueryData<Notice>(noticeKeys.details(currId));
        if (currNotice && !currNotice.isRead) {
          toggleRead({ id: currId, next: true });
        }
      }
    };
  }, [queryClient, toggleRead]);

  return (
    <NoticeActiveContext value={{ selectedId, setSelectedId: handleSetSelectedId }}>
      <NoticeBulkContext
        value={{
          selectionMode,
          setSelectionMode: handleSelectionMode,
          checkedIds,
          toggleChecked,
          selectAll,
          clearChecked,
        }}
      >
        {children}
      </NoticeBulkContext>
    </NoticeActiveContext>
  );
}
