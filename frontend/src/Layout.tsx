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
import { sections } from "@/lib/routes";
import { Window } from "@wailsio/runtime";
import { Suspense, useEffect, useState } from "react";
import { matchPath, Outlet, useLocation } from "react-router";
import { toast } from "sonner";

export default function Layout() {
  const location = useLocation();

  const breadcrumb = (() => {
    for (const [key, items] of Object.entries(sections)) {
      const match = items.find((item) =>
        matchPath({ path: item.path, end: true }, location.pathname),
      );
      if (match) {
        return {
          section: key as keyof typeof sections,
          label: match.label,
        };
      }
    }
    return null;
  })();

  const { section, label } = breadcrumb ?? {
    section: "Workspace",
    label: "Unknown",
  };

  return (
    <div className="flex h-screen w-full">
      <AppSidebar />
      <SidebarInset className="overflow-hidden">
        <Header section={section} label={label} />
        <main className="min-h-0 flex-1 overflow-hidden">
          <Suspense fallback={<Loading />}>
            <Outlet />
          </Suspense>
        </main>
      </SidebarInset>
    </div>
  );
}

function Header({ section, label }: { section: string; label: string }) {
  const [maximized, setMaximized] = useState(false);

  useEffect(() => {
    const handleResize = () => {
      Window.IsMaximised()
        .then(setMaximized)
        .catch((err) => {
          toast.error("Failed to get window state", {
            description: (err as Error).message,
          });
        });
    };

    handleResize();
    window.addEventListener("resize", handleResize);

    return () => window.removeEventListener("resize", handleResize);
  }, []);

  const toggleMaximize = async () => {
    if (maximized) {
      await Window.UnMaximise();
    } else {
      await Window.Maximise();
    }
    setMaximized(!maximized);
  };

  return (
    <header
      onDoubleClick={() => void toggleMaximize()}
      className="wails-drag sticky top-0 z-50 flex h-11 shrink-0 items-center gap-2 border-b p-2 transition-[width,height] ease-linear"
    >
      <SidebarTrigger />
      <Separator
        orientation="vertical"
        className="mr-2 data-[orientation=vertical]:h-4 data-[orientation=vertical]:self-center"
      />
      <Breadcrumb className="pointer-events-none select-none">
        <BreadcrumbList className="min-w-xs">
          <BreadcrumbItem>
            <BreadcrumbLink className="capitalize">{section}</BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{label}</BreadcrumbPage>
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
