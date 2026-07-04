import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

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
