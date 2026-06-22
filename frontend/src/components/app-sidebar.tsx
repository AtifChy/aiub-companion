import { NavGroup } from "@/components/nav-group";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { sections } from "@/lib/routes";
import { Command, Copyright, type LucideIcon } from "lucide-react";
import React, { useEffect, useState } from "react";
import type { Info as AppInfo } from "@bindings/meta";
import { GetInfo as GetAppInfo } from "@bindings/meta/service";

export type NavItem = {
  label: string;
  path: string;
  icon: LucideIcon;
  component?: React.ComponentType;
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const [appInfo, setAppInfo] = useState<AppInfo | null>(null);

  useEffect(() => {
    GetAppInfo()
      .then(setAppInfo)
      .catch((err) => console.error("Failed to get app info: ", err));
  }, []);

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <div className="bg-sidebar-primary text-sidebar-primary-foreground flex size-8 items-center justify-center rounded-lg">
                <Command size="6" />
              </div>
              <div className="flex flex-col gap-0.5 truncate leading-none group-data-[collapsible=icon]:hidden">
                <span>{appInfo?.name ?? "MyApp"}</span>
                <span>{appInfo?.version ?? "v0.0.0"}</span>
              </div>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavGroup title="workspace" items={sections.workspace} />
        <NavGroup title="tools" items={sections.tools} />
        <NavGroup title="others" items={sections.others} />
      </SidebarContent>
      <SidebarFooter className="flex items-center truncate group-data-[collapsible=icon]:hidden">
        <div className="text-muted-foreground flex items-center gap-1 text-xs">
          <Copyright className="size-3" />
          {new Date().getFullYear()} AIUB Companion
        </div>
      </SidebarFooter>
    </Sidebar>
  );
}
