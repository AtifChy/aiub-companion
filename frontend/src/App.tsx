import Layout from "@/Layout";
import { SettingsProvider } from "@/components/settings-provider";
import { ThemeProvider } from "@/components/theme-provider";
import { SidebarProvider } from "@/components/ui/sidebar";
import { TooltipProvider } from "@/components/ui/tooltip";
import { routes } from "@/lib/routes";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { HashRouter, Navigate, Route, Routes } from "react-router";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <SettingsProvider>
        <ThemeProvider>
          <SidebarProvider defaultOpen={false}>
            <TooltipProvider>
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
                </Routes>
              </HashRouter>
            </TooltipProvider>
          </SidebarProvider>
        </ThemeProvider>
      </SettingsProvider>
    </QueryClientProvider>
  );
}

export default App;
