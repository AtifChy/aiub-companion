import { Label } from "@/components/ui/label";

interface SettingRowProps {
  label: string;
  description: string;
  children: React.ReactNode;
}

export function SettingRow({ label, description, children }: SettingRowProps) {
  return (
    <div className="flex items-center justify-between gap-4">
      <div className="space-y-1">
        <Label>{label}</Label>
        <p className="text-muted-foreground text-sm">{description}</p>
      </div>
      {children}
    </div>
  );
}
