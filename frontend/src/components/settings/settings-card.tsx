import { cva, type VariantProps } from "class-variance-authority";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { cn } from "@/lib/utils";

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

export function SettingsCard({ title, variant = "default", children }: SettingsCardProps) {
  return (
    <Card size="sm" className="shadow-xs">
      <CardHeader className={cn("border-b", cardVariants({ variant }))}>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-4">{children}</CardContent>
    </Card>
  );
}
