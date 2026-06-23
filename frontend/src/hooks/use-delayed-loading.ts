import { useEffect, useState } from "react";

export function useDelayedLoading(value: boolean, delay: number = 300) {
  const [delayed, setDelayed] = useState(false);

  useEffect(() => {
    if (!value) {
      // eslint-disable-next-line react-hooks/set-state-in-effect, react-you-might-not-need-an-effect/no-adjust-state-on-prop-change
      setDelayed(false);
      return;
    }
    const timer = setTimeout(() => setDelayed(true), delay);
    return () => clearTimeout(timer);
  }, [delay, value]);

  return delayed;
}
