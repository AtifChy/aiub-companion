import { NavGroup } from "@/components/nav-group";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenuButton,
} from "@/components/ui/sidebar";
import { sections } from "@/lib/routes";
import { Service as ConfigService, type BuildInfo } from "@bindings/config";
import { Service as DesktopService } from "@bindings/desktop";
import {
  ChevronsUpDownIcon,
  CommandIcon,
  CopyrightIcon,
  type LucideIcon,
} from "lucide-react";
import React, { useEffect, useState } from "react";

export type NavItem = {
  label: string;
  path: string;
  icon: LucideIcon;
  component?: React.ComponentType;
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const [appInfo, setAppInfo] = useState<BuildInfo | null>(null);

  useEffect(() => {
    ConfigService.GetBuildInfo()
      .then(setAppInfo)
      .catch((err) => console.error("Failed to get app info: ", err));
  }, []);

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <SidebarMenuButton
          onClick={() => void DesktopService.ShowAboutWindow()}
          size="lg"
          className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
        >
          <div className="bg-sidebar-primary text-sidebar-primary-foreground flex size-8 items-center justify-center rounded-lg">
            <CommandIcon size="6" />
          </div>
          <div className="flex flex-col gap-0.5 truncate leading-none group-data-[collapsible=icon]:hidden">
            <span>{appInfo?.name ?? "N/A"}</span>
            <span>{appInfo?.version ?? "N/A"}</span>
          </div>
          <ChevronsUpDownIcon className="ml-auto group-data-[collapsible=icon]:hidden" />
        </SidebarMenuButton>
      </SidebarHeader>
      <SidebarContent>
        <NavGroup title="workspace" items={sections.workspace} />
        <NavGroup title="tools" items={sections.tools} />
        <NavGroup title="others" items={sections.others} />
      </SidebarContent>
      <SidebarFooter className="flex items-center truncate group-data-[collapsible=icon]:hidden">
        <div className="text-muted-foreground flex items-center gap-1 text-xs">
          <CopyrightIcon className="size-3" />
          {new Date().getFullYear()} AIUB Companion
        </div>
      </SidebarFooter>
    </Sidebar>
  );
}
