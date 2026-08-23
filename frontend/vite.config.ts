import path from "node:path";

import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import wails from "@wailsio/runtime/plugins/vite";
import { defineConfig } from "vite";

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env["WAILS_VITE_PORT"]) || 9245,
    strictPort: true,
  },
  plugins: [react({ compiler: true }), tailwindcss(), wails("./bindings")],
  resolve: {
    alias: {
      "@": path.resolve(import.meta.dirname, "src"),
      "@bindings": path.resolve(import.meta.dirname, "bindings/aiub-companion/internal"),
    },
  },
  build: {
    rolldownOptions: {
      output: {
        codeSplitting: {
          groups: [
            {
              name: "wails",
              test: /[\\/]node_modules[\\/]@wailsio[\\/]/,
              priority: 10,
            },
            {
              name: "react-vendor",
              test: /[\\/]node_modules[\\/](react|react-dom|scheduler)[\\/]/,
              priority: 5,
            },
          ],
        },
      },
    },
  },
});
