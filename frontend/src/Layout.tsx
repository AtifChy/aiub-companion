import { Window } from "@wailsio/runtime";
import { Suspense, useEffect, useState } from "react";
import { matchPath, Outlet, useLocation } from "react-router";

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

export default function Layout() {
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

  const [maximized, setMaximized] = useState(false);
  useEffect(() => void Window.IsMaximised().then(setMaximized), []);

  const toggleMaximize = async () => {
    await Window.ToggleMaximise();
    setMaximized(await Window.IsMaximised());
  };

  return (
    <header
      onDoubleClick={() => void toggleMaximize()}
      className="sticky top-0 z-50 flex h-11 shrink-0 items-center gap-2 border-b p-2 transition-[width,height] ease-linear wails-drag"
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
        maximized={maximized}
        onMinimize={() => void Window.Minimise()}
        onMaximize={() => void toggleMaximize()}
        onClose={() => void Window.Close()}
      />
    </header>
  );
}
