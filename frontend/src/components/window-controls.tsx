import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { CopyIcon, MinusIcon, SquareIcon, XIcon } from "lucide-react";
import { useEffect, useState } from "react";

interface WindowControlsProps {
  maximized?: boolean;
  onMinimize?: () => void;
  onMaximize?: () => void;
  onClose: () => void;
}

export function WindowControls({
  maximized,
  onMinimize,
  onMaximize,
  onClose,
}: WindowControlsProps) {
  const [focused, setFocused] = useState(() => document.hasFocus());

  useEffect(() => {
    const handleFocus = () => setFocused(true);
    const handleBlur = () => setFocused(false);

    window.addEventListener("focus", handleFocus);
    window.addEventListener("blur", handleBlur);

    return () => {
      window.removeEventListener("focus", handleFocus);
      window.removeEventListener("blur", handleBlur);
    };
  }, []);

  return (
    <div className="ml-auto flex items-center gap-1 wails-no-drag">
      {onMinimize && (
        <Button
          onClick={onMinimize}
          variant="ghost"
          className={cn(
            "flex size-7 items-center justify-center rounded-sm transition-colors dark:hover:bg-primary/20",
            focused ? "text-foreground" : "text-muted-foreground",
          )}
        >
          <MinusIcon strokeWidth={1.5} className="size-4" />
        </Button>
      )}
      {onMaximize && (
        <Button
          onClick={onMaximize}
          variant="ghost"
          className={cn(
            "flex size-7 items-center justify-center rounded-sm transition-colors dark:hover:bg-primary/20",
            focused ? "text-foreground" : "text-muted-foreground",
          )}
        >
          {maximized ? (
            <CopyIcon strokeWidth={1.5} className="size-3.5 -scale-x-100" />
          ) : (
            <SquareIcon strokeWidth={1.5} className="size-3.5" />
          )}
        </Button>
      )}
      <Button
        onClick={onClose}
        variant="destructive"
        className={cn(
          "flex size-7 items-center justify-center rounded-sm transition-colors",
          "bg-transparent hover:text-destructive dark:bg-transparent dark:hover:text-destructive",
          focused ? "text-foreground" : "text-muted-foreground",
        )}
      >
        <XIcon strokeWidth={1.5} className="size-4" />
      </Button>
    </div>
  );
}
