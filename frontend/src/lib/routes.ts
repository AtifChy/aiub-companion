import {
  BellIcon,
  CalendarIcon,
  GraduationCapIcon,
  HelpCircleIcon,
  SettingsIcon,
  type LucideIcon,
} from "lucide-react";
import type { ComponentType } from "react";

import { lazyWithPreload, type PreloadableComponent } from "@/lib/lazy-with-preload";

const NoticesPage = lazyWithPreload(() => import("@/pages/notices"));
const RoutinePage = lazyWithPreload(() => import("@/pages/routine"));
const SemesterPage = lazyWithPreload(() => import("@/pages/semester"));
const SettingsPage = lazyWithPreload(() => import("@/pages/settings"));
const HelpPage = lazyWithPreload(() => import("@/pages/help"));

export const AboutPage = lazyWithPreload(() => import("@/pages/about"));

export interface RouteItem {
  label: string;
  path: string;
  icon: LucideIcon;
  component: PreloadableComponent<ComponentType<unknown>>;
}

type Section = "workspace" | "tools" | "others";

export const sections: Record<Section, RouteItem[]> = {
  workspace: [
    {
      label: "Notices",
      path: "/notices",
      icon: BellIcon,
      component: NoticesPage,
    },
    {
      label: "Semester",
      path: "/semester",
      icon: GraduationCapIcon,
      component: SemesterPage,
    },
  ],
  tools: [
    {
      label: "Routine",
      path: "/routine",
      icon: CalendarIcon,
      component: RoutinePage,
    },
  ],
  others: [
    {
      label: "Settings",
      path: "/settings",
      icon: SettingsIcon,
      component: SettingsPage,
    },
    {
      label: "Help",
      path: "/help",
      icon: HelpCircleIcon,
      component: HelpPage,
    },
  ],
};

export const routes = Object.values(sections)
  .flat()
  .filter((item): item is RouteItem & { component: React.ComponentType } => !!item.component) as [
  RouteItem & { component: React.ComponentType },
  ...(RouteItem & { component: React.ComponentType })[],
];
