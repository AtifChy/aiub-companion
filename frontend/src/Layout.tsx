import { Events, System, Window } from "@wailsio/runtime";
import { Suspense, useEffect, useState } from "react";
import { matchPath, Outlet, useLocation, useNavigate } from "react-router";

import { AppSidebar } from "@/components/app-sidebar";
import Loading from "@/components/loading";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Separator } from "@/components/ui/separator";
import { SidebarInset, SidebarTrigger } from "@/components/ui/sidebar";
import { WindowControls } from "@/components/window-controls";
import { routes } from "@/lib/routes";
import { cn } from "@/lib/utils";

export default function Layout() {
  const navigate = useNavigate();

  useEffect(() => {
    return Events.On("notice:open", () => {
      void navigate("/notices");
    });
  }, [navigate]);

  return (
    <div className="flex h-screen w-full">
      <AppSidebar />
      <SidebarInset className="overflow-hidden">
        <Header />
        <main className="min-h-0 flex-1 overflow-hidden">
          <Suspense fallback={<Loading />}>
            <Outlet />
          </Suspense>
        </main>
      </SidebarInset>
    </div>
  );
}

function Header() {
  const location = useLocation();
  const breadcrumb = routes.find((route) =>
    matchPath({ path: route.path, end: true }, location.pathname),
  );

  const [isMaximized, setIsMaximized] = useState(false);
  useEffect(() => {
    void Window.IsMaximised().then(setIsMaximized);
    return Events.On("window:maximized", (e) => setIsMaximized(e.data));
  }, []);

  return (
    <header
      onDoubleClick={() => void Window.ToggleMaximise()}
      className={cn(
        "sticky top-0 z-50 flex h-11 shrink-0 items-center gap-2 border-b p-2 transition-[width,height] ease-linear",
        System.IsWindows() ? "wails-app-drag" : "wails-drag",
      )}
    >
      <SidebarTrigger />
      <Separator
        orientation="vertical"
        className="mr-2 data-[orientation=vertical]:h-4 data-[orientation=vertical]:self-center"
      />
      <Breadcrumb className="pointer-events-none select-none">
        <BreadcrumbList className="min-w-xs">
          <BreadcrumbItem>
            <BreadcrumbLink className="capitalize">{breadcrumb?.section}</BreadcrumbLink>
          </BreadcrumbItem>
          {breadcrumb?.label && <BreadcrumbSeparator />}
          <BreadcrumbItem>
            <BreadcrumbPage>{breadcrumb?.label}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <WindowControls
        maximized={isMaximized}
        onMinimize={() => void Window.Minimise()}
        onMaximize={() => void Window.ToggleMaximise()}
        onClose={() => void Window.Close()}
      />
    </header>
  );
}
