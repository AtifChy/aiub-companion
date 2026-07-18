import { createContext, use } from "react";

interface NoticeSelectionContextType {
  selectedId: string | null;
  select: (id: string) => void;
}

export const NoticeSelectionContext = createContext<NoticeSelectionContextType | null>(null);

export function useNoticeSelection() {
  const context = use(NoticeSelectionContext);
  if (!context) {
    throw new Error("useNoticeSelection must be used within a NoticeSelectionProvider");
  }
  return context;
}
