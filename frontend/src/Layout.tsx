import { AppSidebar } from "@/components/app-sidebar";
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
import { Toaster } from "@/components/ui/sonner";
import { sections } from "@/lib/routes";
import { Loader2Icon } from "lucide-react";
import { Suspense, useMemo } from "react";
import { matchPath, Outlet, useLocation } from "react-router";

export default function Layout() {
  const location = useLocation();

  const breadcrumb = useMemo(() => {
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
  }, [location.pathname]);

  const { section, label } = breadcrumb ?? {
    section: "Workspace",
    label: "Unknown",
  };

  return (
    <div className="flex h-screen w-full">
      <AppSidebar />
      <SidebarInset>
        <header className="bg-background sticky top-0 z-50 flex h-12 shrink-0 items-center gap-2 border-b px-4 transition-[width,height] ease-linear">
          <SidebarTrigger className="-ml-1" />
          <Separator
            orientation="vertical"
            className="mr-2 data-[orientation=vertical]:h-4 data-[orientation=vertical]:self-center"
          />
          <Breadcrumb>
            <BreadcrumbList className="min-w-xs">
              <BreadcrumbItem>
                <BreadcrumbLink className="capitalize">
                  {section}
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbPage>{label}</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </header>
        <main className="min-h-0 flex-1 overflow-hidden">
          <Suspense fallback={<Fallback />}>
            <Outlet />
          </Suspense>
        </main>
        <Toaster position="top-center" />
      </SidebarInset>
    </div>
  );
}

const Fallback = () => {
  return (
    <div className="flex h-full items-center justify-center">
      <Loader2Icon className="text-muted-foreground size-8 animate-spin" />
    </div>
  );
};
