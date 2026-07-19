import { useEffect, useRef } from "react";

import { cn } from "@/lib/utils";

interface HorizontalScrollProps {
  children: React.ReactNode;
  className?: string;
}

export function HorizontalScroll({ children, className }: HorizontalScrollProps) {
  const scrollRef = useRef<HTMLDivElement>(null);

  const targetScrollLeft = useRef(0);
  const animationFrame = useRef<number | null>(null);

  useEffect(() => {
    const container = scrollRef.current;
    if (!container) return;

    const animate = () => {
      const diff = targetScrollLeft.current - container.scrollLeft;

      if (Math.abs(diff) < 1) {
        container.scrollLeft = targetScrollLeft.current;
        animationFrame.current = null;
        return;
      }

      // Ease towards the target scroll position
      container.scrollLeft += diff * 0.2;
      animationFrame.current = requestAnimationFrame(animate);
    };

    const handleWheel = (e: WheelEvent) => {
      if (Math.abs(e.deltaX) >= Math.abs(e.deltaY)) return;

      e.preventDefault();

      const delta = e.deltaMode === WheelEvent.DOM_DELTA_LINE ? e.deltaY * 16 : e.deltaY;

      if (animationFrame.current === null) {
        targetScrollLeft.current = container.scrollLeft + delta;
      }

      const maxScroll = container.scrollWidth - container.clientWidth;
      targetScrollLeft.current = Math.max(0, Math.min(maxScroll, targetScrollLeft.current + delta));

      animationFrame.current ??= requestAnimationFrame(animate);

      container.scrollLeft += delta;
    };

    container.addEventListener("wheel", handleWheel, { passive: false });

    return () => {
      container.removeEventListener("wheel", handleWheel);

      if (animationFrame.current !== null) {
        cancelAnimationFrame(animationFrame.current);
      }
    };
  }, []);

  return (
    <div ref={scrollRef} className={cn("scroll-fade-x scrollbar-none overflow-x-auto", className)}>
      {children}
    </div>
  );
}
