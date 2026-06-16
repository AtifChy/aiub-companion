import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

interface LazyTooltipProps {
  children: React.ReactElement;
  content: React.ReactNode;
}

export function AppTooltip({ children, content }: LazyTooltipProps) {
  return (
    <Tooltip>
      <TooltipTrigger render={children} />
      <TooltipContent className="bg-accent text-foreground shadow-2xl [&>div]:bg-inherit [&>div]:shadow-inherit">
        <p className="font-medium">{content}</p>
      </TooltipContent>
    </Tooltip>
  );
}
