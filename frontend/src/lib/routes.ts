import {
  Bell,
  Calculator,
  Calendar,
  HelpCircle,
  LineChart,
  Settings,
  Timer,
  type LucideIcon,
} from "lucide-react";
import type React from "react";
import { lazy } from "react";

const NoticesPage = lazy(() => import("@/pages/notices"));
const CGPAPage = lazy(() => import("@/pages/cgpa"));
const GPATrendPage = lazy(() => import("@/pages/gpa-trend"));
const RoutinePage = lazy(() => import("@/pages/routine"));
const ExamCountdownPage = lazy(() => import("@/pages/exam-countdown"));
const SettingsPage = lazy(() => import("@/pages/settings"));
const HelpPage = lazy(() => import("@/pages/help"));

export type RouteItem = {
  label: string;
  path: string;
  icon: LucideIcon;
  component: React.LazyExoticComponent<React.ComponentType>;
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
      label: "Exam Countdown",
      path: "/exam-countdown",
      icon: Timer,
      component: ExamCountdownPage,
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
  .filter(
    (item): item is RouteItem & { component: React.ComponentType } =>
      !!item.component,
  ) as [
  RouteItem & { component: React.ComponentType },
  ...(RouteItem & { component: React.ComponentType })[],
];
