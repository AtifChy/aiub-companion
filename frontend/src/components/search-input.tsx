import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { Loader2Icon, SearchIcon, SlashIcon, XIcon } from "lucide-react";
import React, {
  forwardRef,
  useEffect,
  useImperativeHandle,
  useRef,
} from "react";

type SearchInputProps = React.ComponentPropsWithoutRef<"div"> & {
  value: string;
  onValueChange: (value: string) => void;
  placeholder?: string;
  searching?: boolean;
  showClearButton?: boolean;
};

export const SearchInput = forwardRef<HTMLInputElement, SearchInputProps>(
  (
    {
      value,
      onValueChange,
      placeholder,
      searching,
      showClearButton,
      className,
      ...props
    },
    ref,
  ) => {
    const innerRef = useRef<HTMLInputElement>(null);

    useImperativeHandle(ref, () => innerRef.current!, []);

    useEffect(() => {
      const handler = (event: KeyboardEvent) => {
        if (event.key === "/" && document.activeElement !== innerRef.current) {
          event.preventDefault();
          innerRef.current?.focus();
        }
      };
      window.addEventListener("keydown", handler);
      return () => window.removeEventListener("keydown", handler);
    }, []);

    return (
      <div className={cn("relative", className)} {...props}>
        <SearchIcon className="text-muted-foreground absolute top-1/2 left-2.5 size-3.5 -translate-y-1/2" />
        <Input
          ref={innerRef}
          value={value}
          onChange={(e) => onValueChange(e.target.value)}
          placeholder={placeholder ?? "Search..."}
          autoComplete="off"
          className="pr-7 pl-8 text-xs"
        />
        {searching ? (
          <Loader2Icon className="text-muted-foreground absolute top-1/2 right-2 size-3.5 -translate-y-1/2 animate-spin" />
        ) : value ? (
          (showClearButton ?? true) && (
            <Button
              variant="ghost"
              onClick={() => onValueChange("")}
              className="text-muted-foreground hover:text-foreground absolute top-1/2 right-2 h-auto -translate-y-1/2 p-0.5"
            >
              <XIcon strokeWidth="2.5" className="size-3.5" />
            </Button>
          )
        ) : (
          <div className="text-secondary-foreground/40 border-secondary/10 bg-secondary/20 absolute top-1/2 right-2 -translate-y-1/2 rounded-sm border p-1">
            <SlashIcon strokeWidth="4" className="size-2.5" />
          </div>
        )}
      </div>
    );
  },
);

SearchInput.displayName = "SearchInput";
