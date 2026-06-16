import babel from "@rolldown/plugin-babel";
import tailwindcss from "@tailwindcss/vite";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import wails from "@wailsio/runtime/plugins/vite";
import path from "node:path";
import { defineConfig } from "vite";

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
    watch: {
      ignored: ["**/bindings/**", "**/.bindings-tmp-*/**"],
    },
  },
  plugins: [
    react(),
    babel({ presets: [reactCompilerPreset()] }),
    tailwindcss(),
    wails("./bindings"),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
      "@bindings": path.resolve(__dirname, "bindings/aiub-companion/internal"),
    },
  },
  build: {
    rolldownOptions: {
      output: {
        codeSplitting: {
          groups: [
            { name: "react", test: /node_modules[\\/]react/, priority: 20 },
            {
              name: "ui",
              test: /node_modules[\\/](antd|radix-ui|base-ui|shadcn|lucide-react)/,
              priority: 15,
            },
            { name: "vendor", test: /node_modules/, priority: 10 },
          ],
        },
      },
    },
  },
});
