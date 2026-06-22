import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { cva, type VariantProps } from "class-variance-authority";

const cardVariants = cva("", {
  variants: {
    variant: {
      default: "text-card-foreground",
      destructive: "text-destructive",
    },
  },
  defaultVariants: {
    variant: "default",
  },
});

interface SettingsCardProps {
  title: string;
  variant?: VariantProps<typeof cardVariants>["variant"];
  children: React.ReactNode;
}

export function SettingsCard({
  title,
  variant = "default",
  children,
}: SettingsCardProps) {
  return (
    <Card size="sm" className="shadow-xs">
      <CardHeader className={cardVariants({ variant })}>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <Separator />
      <CardContent className="flex flex-col gap-4">{children}</CardContent>
    </Card>
  );
}

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

interface SettingSelectProps<T extends string | number> {
  items: { value: T; label: string }[];
  value: T;
  onValueChange: (value: T) => void;
}

export function SettingSelect<T extends string | number>({
  items,
  value,
  onValueChange,
}: SettingSelectProps<T>) {
  return (
    <Select
      items={items}
      value={value}
      onValueChange={(v) => {
        if (!v) return;
        onValueChange(v);
      }}
    >
      <SelectTrigger>
        <SelectValue placeholder="Select option" />
      </SelectTrigger>
      <SelectContent className="p-1">
        {items.map((item) => (
          <SelectItem key={item.value} value={item.value}>
            {item.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
