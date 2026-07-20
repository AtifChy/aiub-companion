import reactX from "eslint-plugin-react-x";
import reactYouMightNotNeedAnEffect from "eslint-plugin-react-you-might-not-need-an-effect";
import { defineConfig } from "oxlint";
import eslintRecommended from "oxlint-config-presets/@eslint/recommended.json" with { type: "json" };
import tsStrictTypeChecked from "oxlint-config-presets/@typescript-eslint/strict-type-checked.json" with { type: "json" };
import tsStylisticTypeChecked from "oxlint-config-presets/@typescript-eslint/stylistic-type-checked.json" with { type: "json" };
import vitestRecommended from "oxlint-config-presets/@vitest/recommended.json" with { type: "json" };
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
    vitestRecommended,
  ],
  plugins: ["eslint", "typescript", "oxc", "react", "react-perf", "vitest"],
  categories: {
    correctness: "error",
  },
  jsPlugins: ["eslint-plugin-react-x", "eslint-plugin-react-you-might-not-need-an-effect"],
  rules: {
    "typescript/no-confusing-void-expression": "off",

    "react/react-in-jsx-scope": "off",
    "react/no-unescaped-entities": "off",

    ...reactX.configs["recommended-typescript"].rules,

    ...reactYouMightNotNeedAnEffect.configs.recommended.rules,
  },
  env: { builtin: true },
  options: {
    typeAware: true,
    typeCheck: true,
  },
});
