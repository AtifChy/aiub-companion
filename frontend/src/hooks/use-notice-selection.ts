import { createContext, use } from "react";

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
