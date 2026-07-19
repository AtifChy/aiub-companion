import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";

interface AppTooltipProps {
  children: React.ReactElement;
  content: React.ReactNode;
  side?: "top" | "right" | "bottom" | "left";
}

export function AppTooltip({ children, content, side = "top" }: AppTooltipProps) {
  return (
    <Tooltip>
      <TooltipTrigger render={children} />
      <TooltipContent
        side={side}
        className="bg-accent text-foreground shadow-2xl [&>div]:bg-inherit [&>div]:shadow-inherit"
      >
        <p className="font-medium">{content}</p>
      </TooltipContent>
    </Tooltip>
  );
}
