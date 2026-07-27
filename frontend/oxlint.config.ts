import tanstackQuery from "@tanstack/eslint-plugin-query";
import reactYouMightNotNeedAnEffect from "eslint-plugin-react-you-might-not-need-an-effect";
import { defineConfig } from "oxlint";
import eslintRecommended from "oxlint-config-presets/@eslint/recommended.json" with { type: "json" };
import tsStrictTypeChecked from "oxlint-config-presets/@typescript-eslint/strict-type-checked.json" with { type: "json" };
import tsStylisticTypeChecked from "oxlint-config-presets/@typescript-eslint/stylistic-type-checked.json" with { type: "json" };
import reactHooksRecommended from "oxlint-config-presets/react-hooks/recommended-latest.json" with { type: "json" };
import reactRefreshRecommended from "oxlint-config-presets/react-refresh/recommended.json" with { type: "json" };
import reactRecommended from "oxlint-config-presets/react/recommended.json" with { type: "json" };

export default defineConfig({
  ignorePatterns: ["dist/**", "node_modules/**", "bindings/**", "src/components/ui/**/*{.ts,.tsx}"],
  extends: [
    eslintRecommended,
    tsStrictTypeChecked,
    tsStylisticTypeChecked,
    reactRecommended,
    reactHooksRecommended,
    reactRefreshRecommended,
  ],
  plugins: ["eslint", "typescript", "oxc", "react", "react-perf"],
  categories: {
    correctness: "error",
  },
  jsPlugins: ["eslint-plugin-react-you-might-not-need-an-effect", "@tanstack/eslint-plugin-query"],
  rules: {
    "typescript/no-confusing-void-expression": "off",

    "react/react-in-jsx-scope": "off",
    "react/no-unescaped-entities": "off",

    ...reactYouMightNotNeedAnEffect.configs.strict.rules,

    ...tanstackQuery.configs.recommended.rules,
  },
  env: { builtin: true },
  options: {
    typeAware: true,
    typeCheck: true,
  },
});
