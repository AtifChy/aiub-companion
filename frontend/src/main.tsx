import React from "react";

import "@/style.css";
import { createRoot } from "react-dom/client";

import App from "@/App";

const container = document.getElementById("root");
const root = createRoot(container!);

if (import.meta.env.DEV) {
  void import("react-scan").then(({ scan }) => scan());
}

root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
