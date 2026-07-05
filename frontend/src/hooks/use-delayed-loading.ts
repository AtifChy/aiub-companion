import { useEffect, useState } from "react";

export function useDelayedLoading(value: boolean, delay: number = 300) {
  const [delayed, setDelayed] = useState(false);

  if (!value && delayed) {
    setDelayed(false);
  }

  useEffect(() => {
    if (!value) return;
    const timer = setTimeout(() => setDelayed(true), delay);
    return () => clearTimeout(timer);
  }, [delay, value]);

  return delayed;
}
