import { Link, useLocation } from "react-router";

import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import type { RouteItem } from "@/lib/routes";

interface NavGroupProps {
  title: string;
  items: RouteItem[];
}

export function NavGroup({ title, items }: NavGroupProps) {
  // const { isMobile } = useSidebar();
  const location = useLocation();

  return (
    <SidebarGroup className="overflow-hidden pt-0 pb-2 text-ellipsis whitespace-nowrap group-data-[collapsible=icon]:pb-0">
      <SidebarGroupLabel className="capitalize">{title}</SidebarGroupLabel>
      <SidebarMenu>
        {items.map((item) => (
          <SidebarMenuItem key={item.label}>
            <SidebarMenuButton
              isActive={location.pathname === item.path}
              onMouseEnter={() => void item.component.preload?.()}
              onFocus={() => void item.component.preload?.()}
              render={
                <Link to={item.path} className="flex items-center gap-2">
                  <item.icon />
                  <span>{item.label}</span>
                </Link>
              }
              className="mb-1 group-data-[collapsible=icon]:mt-1"
            />
          </SidebarMenuItem>
        ))}
      </SidebarMenu>
    </SidebarGroup>
  );
}
