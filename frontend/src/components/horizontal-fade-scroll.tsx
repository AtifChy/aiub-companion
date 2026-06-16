import { cn } from "@/lib/utils";
import { useCallback, useLayoutEffect, useRef, useState } from "react";

type ScrollState = {
  left: boolean;
  right: boolean;
};

interface HorizontalFadeScrollProps {
  children: React.ReactNode;
  className?: string;
}

export function HorizontalFadeScroll({
  children,
  className,
}: HorizontalFadeScrollProps) {
  const [scrollState, setScrollState] = useState<ScrollState>({
    left: false,
    right: false,
  });
  const scrollRef = useRef<HTMLDivElement>(null);

  const checkOverflow = useCallback(() => {
    const container = scrollRef.current;
    if (container) {
      const { scrollLeft, scrollWidth, clientWidth } = container;
      const left = scrollLeft > 0;
      const right = scrollLeft < scrollWidth - clientWidth;
      setScrollState((prev) =>
        prev.left === left && prev.right === right ? prev : { left, right },
      );
    }
  }, []);

  useLayoutEffect(() => {
    const container = scrollRef.current;
    if (!container) return;

    checkOverflow();

    const resizeObserver = new ResizeObserver(() => checkOverflow());
    resizeObserver.observe(container);

    return () => resizeObserver.disconnect();
  }, [children]);

  const handleWheel = useCallback((e: React.WheelEvent) => {
    const container = scrollRef.current;
    if (!container) return;

    e.preventDefault();
    container.scrollLeft += e.deltaY / 2;
  }, []);

  return (
    <div
      data-overflow-left={scrollState.left}
      data-overflow-right={scrollState.right}
      data-overflow-both={scrollState.left && scrollState.right}
      className={cn(
        "relative w-full overflow-hidden",
        "data-[overflow-both=true]:mask-[linear-gradient(to_right,transparent_0%,black_40px,black_calc(100%-40px),transparent_100%)]",
        "data-[overflow-left=true]:data-[overflow-both=false]:mask-[linear-gradient(to_right,transparent_0%,black_40px,black_100%)]",
        "data-[overflow-right=true]:data-[overflow-both=false]:mask-[linear-gradient(to_right,black_0%,black_calc(100%-40px),transparent_100%)]",
      )}
    >
      <div
        ref={scrollRef}
        onScroll={checkOverflow}
        onWheel={handleWheel}
        className={className}
      >
        {children}
      </div>
    </div>
  );
}
