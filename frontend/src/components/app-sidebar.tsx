import { Service as DesktopService } from "@bindings/desktop";
import { BuildInfo, Service as MetaService } from "@bindings/meta";
import { ChevronsUpDownIcon, CommandIcon, CopyrightIcon } from "lucide-react";
import React, { useEffect, useState } from "react";

import { NavGroup } from "@/components/nav-group";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenuButton,
} from "@/components/ui/sidebar";
import { AboutPage, sections } from "@/lib/routes";

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const [appInfo, setAppInfo] = useState<BuildInfo | null>(null);
  const year = new Date().getFullYear();

  useEffect(() => {
    MetaService.GetBuildInfo()
      .then(setAppInfo)
      .catch((err) => console.error("Failed to get app info: ", err));
  }, []);

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <SidebarMenuButton
          onClick={() => void DesktopService.ShowAboutWindow()}
          onMouseEnter={() => void AboutPage.preload?.()}
          onFocus={() => void AboutPage.preload?.()}
          size="lg"
          className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
        >
          <div className="flex size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
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
        <NavGroup title="others" items={sections.others} />
      </SidebarContent>
      <SidebarFooter className="flex items-center truncate group-data-[collapsible=icon]:hidden">
        <div className="flex items-center gap-1 text-xs text-muted-foreground">
          <CopyrightIcon className="size-3" />
          {year} AIUB Companion
        </div>
      </SidebarFooter>
    </Sidebar>
  );
}
