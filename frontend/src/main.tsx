import App from "@/App";
import "@/style.css";
import React from "react";
import { createRoot } from "react-dom/client";

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
