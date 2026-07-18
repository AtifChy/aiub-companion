import * as React from "react";

const MOBILE_BREAKPOINT = 768;

const query = `(max-width: ${String(MOBILE_BREAKPOINT - 1)}px)`;

export function useIsMobile() {
  return React.useSyncExternalStore(
    (callback) => {
      const mql = window.matchMedia(query);
      mql.addEventListener("change", callback);
      return () => {
        mql.removeEventListener("change", callback);
      };
    },
    () => window.matchMedia(query).matches,
    () => false,
  );
}
