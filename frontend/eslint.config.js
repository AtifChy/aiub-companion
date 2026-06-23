import js from "@eslint/js";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import reactYouMightNotNeedAnEffect from "eslint-plugin-react-you-might-not-need-an-effect";
import { defineConfig, globalIgnores } from "eslint/config";
import tseslint from "typescript-eslint";

export default defineConfig([
  globalIgnores([
    "dist/**",
    "node_modules/**",
    "bindings/**",
    "vite.config.ts",
    "src/components/ui/**/*{.ts,.tsx}",
    "src/hooks/use-mobile.ts",
  ]),
  {
    files: ["**/*{.ts,.tsx}"],
    extends: [
      js.configs.recommended,
      tseslint.configs.recommended,
      reactHooks.configs.flat.recommended,
      reactRefresh.configs.vite,
      reactYouMightNotNeedAnEffect.configs.recommended,
    ],
    languageOptions: {
      ecmaVersion: 2020,
      parserOptions: {
        projectService: true,
      },
    },
  },
]);
