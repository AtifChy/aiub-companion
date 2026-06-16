import Layout from "@/Layout";
import { ThemeProvider } from "@/components/theme-provider";
import { SidebarProvider } from "@/components/ui/sidebar";
import { TooltipProvider } from "@/components/ui/tooltip";
import { IS_DESKTOP } from "@/lib/env";
import { routes } from "@/lib/routes";
import {
  BrowserRouter,
  HashRouter,
  Navigate,
  Route,
  Routes,
} from "react-router";

const Router = IS_DESKTOP ? HashRouter : BrowserRouter;

function App() {
  return (
    <ThemeProvider storageKey="vite-ui-theme">
      <SidebarProvider>
        <TooltipProvider>
          <Router>
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
          </Router>
        </TooltipProvider>
      </SidebarProvider>
    </ThemeProvider>
  );
}

export default App;
