import { createContext, use, useState, type ReactNode } from "react";

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

  const selectAll = (ids: string[]) => {
    setCheckedIds(new Set(ids));
  };

  const clearChecked = () => {
    setCheckedIds(new Set());
  };

  const updateSelectionMode = (on: boolean) => {
    setSelectionMode(on);
    if (!on) clearChecked();
  };

  return (
    <NoticeActiveContext value={{ selectedId, setSelectedId }}>
      <NoticeBulkContext
        value={{
          selectionMode,
          setSelectionMode: updateSelectionMode,
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
