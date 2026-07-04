import Layout from "@/Layout";
import { SettingsProvider, useSettings } from "@/components/settings-provider";
import { ThemeProvider } from "@/components/theme-provider";
import { SidebarProvider } from "@/components/ui/sidebar";
import { Toaster } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { AboutPage, routes } from "@/lib/routes";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { HashRouter, Navigate, Route, Routes } from "react-router";

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
          <Toaster position="top-center" richColors />
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
  const initialRoute = routes[0];
  return (
    <HashRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route index element={<Navigate to={initialRoute.path} />} />
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
