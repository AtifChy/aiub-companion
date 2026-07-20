import { createContext, use } from "react";

import type { NoticeFilters } from "@/hooks/use-notices";

interface NoticeFiltersContextType {
  filters: NoticeFilters;
  setFilters: (patch: Partial<NoticeFilters>) => void;
  clearFilters: () => void;
}

export const NoticeFiltersContext = createContext<NoticeFiltersContextType | null>(null);

export const useNoticeFilters = () => {
  const context = use(NoticeFiltersContext);
  if (!context) {
    throw new Error("useNoticeFilters must be used within a NoticeFiltersProvider");
  }
  return context;
};
