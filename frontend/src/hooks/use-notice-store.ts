import { create } from "zustand";

import type { NoticeFilters } from "@/hooks/use-notices";

const INITIAL_FILTERS: NoticeFilters = {
  search: "",
  category: "all",
  urgent: false,
  pinned: false,
  unread: false,
};

interface NoticeState {
  // Selection State
  selectedId: string | null;
  selectionMode: boolean;
  checkedIds: Set<string>;

  setSelectedId: (id: string | null) => void;
  setSelectionMode: (enabled: boolean) => void;
  toggleChecked: (id: string) => void;
  selectAll: (ids: string[]) => void;
  clearChecked: () => void;

  // Filters State
  filters: NoticeFilters;

  setFilters: (patch: Partial<NoticeFilters>) => void;
  clearFilters: () => void;
}

export const useNoticeStore = create<NoticeState>((set) => ({
  // Selection Initial State
  selectedId: null,
  selectionMode: false,
  checkedIds: new Set(),

  // Filters Initial State
  filters: INITIAL_FILTERS,

  // Selection Actions
  setSelectedId: (id) => set({ selectedId: id }),

  setSelectionMode: (enabled) =>
    set((_) => {
      if (!enabled) {
        return { selectionMode: false, checkedIds: new Set() };
      }
      return { selectionMode: true };
    }),

  toggleChecked: (id) =>
    set((state) => {
      const next = new Set(state.checkedIds);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return { checkedIds: next };
    }),

  selectAll: (ids) => set({ checkedIds: new Set(ids) }),

  clearChecked: () => set({ checkedIds: new Set() }),

  // Filters Actions
  setFilters: (patch) =>
    set((state) => ({
      filters: { ...state.filters, ...patch },
    })),

  clearFilters: () => set({ filters: INITIAL_FILTERS }),
}));
