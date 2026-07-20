import { createContext, use, useState, type ReactNode } from "react";

import type { NoticeFilters } from "@/hooks/use-notices";

interface NoticeFiltersContextType {
  filters: NoticeFilters;
  setFilters: (patch: Partial<NoticeFilters>) => void;
  clearFilters: () => void;
}

const NoticeFiltersContext = createContext<NoticeFiltersContextType | null>(null);

const INITIAL_FILTERS: NoticeFilters = {
  search: "",
  category: "all",
  urgent: false,
  pinned: false,
  unread: false,
};

export function NoticeFiltersProvider({ children }: { children: ReactNode }) {
  const [filters, setFilters] = useState<NoticeFilters>(INITIAL_FILTERS);

  const updateFilters = (patch: Partial<NoticeFilters>) =>
    setFilters((prev) => ({ ...prev, ...patch }));

  const clearFilters = () => setFilters(INITIAL_FILTERS);

  return (
    <NoticeFiltersContext value={{ filters, setFilters: updateFilters, clearFilters }}>
      {children}
    </NoticeFiltersContext>
  );
}

// oxlint-disable-next-line react/only-export-components
export const useNoticeFilters = () => {
  const context = use(NoticeFiltersContext);
  if (!context) {
    throw new Error("useNoticeFilters must be used within a NoticeFiltersProvider");
  }
  return context;
};
