import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { HashRouter, Navigate, Route, Routes } from "react-router";

import { SettingsProvider, useSettings } from "@/components/providers/settings-provider";
import { ThemeProvider } from "@/components/providers/theme-provider";
import { UpdateProvider } from "@/components/providers/update-provider";
import { SidebarProvider } from "@/components/ui/sidebar";
import { Toaster } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { UpdateDialog } from "@/components/update-dialog";
import Layout from "@/Layout";
import { AboutPage, routes } from "@/lib/routes";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: false,
    },
  },
});

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <SettingsProvider>
        <UpdateProvider>
          <ShellProviders>
            <AppRouter />
            <Toaster position="top-center" richColors />
            <UpdateDialog />
          </ShellProviders>
        </UpdateProvider>
      </SettingsProvider>
    </QueryClientProvider>
  );
}

function ShellProviders({ children }: { children: React.ReactNode }) {
  const { config } = useSettings();
  return (
    <ThemeProvider>
      <SidebarProvider
        defaultOpen={config.launch.sidebar_open}
        style={
          {
            "--sidebar-width": "14rem",
            "--sidebar-width-mobile": "14rem",
          } as React.CSSProperties
        }
      >
        <TooltipProvider>{children}</TooltipProvider>
      </SidebarProvider>
    </ThemeProvider>
  );
}

function AppRouter() {
  const initialRoute = routes[0];
  return (
    <HashRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route index element={<Navigate to={initialRoute.path} />} />
          {routes.map((route) => (
            <Route key={route.path} path={route.path} element={<route.component />} />
          ))}
        </Route>

        <Route path="/about" element={<AboutPage />} />
      </Routes>
    </HashRouter>
  );
}
