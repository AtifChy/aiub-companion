import { createContext, use, useState } from "react";

interface NoticeSelectionContextType {
  selectedId: string | null;
  select: (id: string) => void;

  selectionMode: boolean;
  setSelectionMode: (on: boolean) => void;
  checkedIds: Set<string>;
  toggleChecked: (id: string) => void;
  selectAll: (ids: string[]) => void;
  clearChecked: () => void;
}

export const NoticeSelectionContext = createContext<NoticeSelectionContextType | null>(null);

export function useNoticeSelection() {
  const context = use(NoticeSelectionContext);
  if (!context) {
    throw new Error("useNoticeSelection must be used within a NoticeSelectionProvider");
  }
  return context;
}

export function useNoticeSelectionState() {
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

  return {
    selectionMode,
    setSelectionMode: updateSelectionMode,
    checkedIds,
    toggleChecked,
    selectAll,
    clearChecked,
  };
}
