import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

interface SettingRowProps {
  label: string;
  description: string;
  orientation?: "horizontal" | "vertical";
  className?: string;
  children: React.ReactNode;
}

export function SettingRow({
  label,
  description,
  orientation = "horizontal",
  className,
  children,
}: SettingRowProps) {
  return (
    <div
      className={cn(
        "flex",
        orientation === "horizontal"
          ? "flex-row items-center justify-between gap-4"
          : "flex-col gap-3",
        className,
      )}
    >
      <div className="space-y-1">
        <Label>{label}</Label>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
      {children}
    </div>
  );
}
