import { lazyWithPreload, type PreloadableComponent } from "@/lib/lazy-with-preload";
import {
  Bell,
  Calculator,
  Calendar,
  GraduationCap,
  HelpCircle,
  LineChart,
  Settings,
  type LucideIcon,
} from "lucide-react";
import type { ComponentType } from "react";

const NoticesPage = lazyWithPreload(() => import("@/pages/notices"));
const CGPAPage = lazyWithPreload(() => import("@/pages/cgpa"));
const GPATrendPage = lazyWithPreload(() => import("@/pages/gpa-trend"));
const RoutinePage = lazyWithPreload(() => import("@/pages/routine"));
const SemesterPage = lazyWithPreload(() => import("@/pages/semester"));
const SettingsPage = lazyWithPreload(() => import("@/pages/settings"));
const HelpPage = lazyWithPreload(() => import("@/pages/help"));

export const AboutPage = lazyWithPreload(() => import("@/pages/about"));

export type RouteItem = {
  label: string;
  path: string;
  icon: LucideIcon;
  component: PreloadableComponent<ComponentType<unknown>>;
};

type Section = "workspace" | "tools" | "others";

export const sections: Record<Section, RouteItem[]> = {
  workspace: [
    {
      label: "Notices",
      path: "/notices",
      icon: Bell,
      component: NoticesPage,
    },
    {
      label: "CGPA Calculator",
      path: "/cgpa",
      icon: Calculator,
      component: CGPAPage,
    },
    {
      label: "GPA Trend",
      path: "/gpa-trend",
      icon: LineChart,
      component: GPATrendPage,
    },
  ],
  tools: [
    {
      label: "Routine",
      path: "/routine",
      icon: Calendar,
      component: RoutinePage,
    },
    {
      label: "Semester",
      path: "/semester",
      icon: GraduationCap,
      component: SemesterPage,
    },
  ],
  others: [
    {
      label: "Settings",
      path: "/settings",
      icon: Settings,
      component: SettingsPage,
    },
    {
      label: "Help",
      path: "/help",
      icon: HelpCircle,
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
