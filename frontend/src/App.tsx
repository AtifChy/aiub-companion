import Layout from "@/Layout";
import { SettingsProvider, useSettings } from "@/components/settings-provider";
import { ThemeProvider } from "@/components/theme-provider";
import { SidebarProvider } from "@/components/ui/sidebar";
import { TooltipProvider } from "@/components/ui/tooltip";
import { routes } from "@/lib/routes";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { lazy } from "react";
import { HashRouter, Navigate, Route, Routes } from "react-router";

const AboutPage = lazy(() => import("@/pages/about"));

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
        <ShellProviders>
          <AppRouter />
        </ShellProviders>
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
  return (
    <HashRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route index element={<Navigate to={routes[0].path} />} />
          {routes.map((route) => (
            <Route
              key={route.path}
              path={route.path}
              element={<route.component />}
            />
          ))}
        </Route>
        <Route path="/about" element={<AboutPage />} />
      </Routes>
    </HashRouter>
  );
}
