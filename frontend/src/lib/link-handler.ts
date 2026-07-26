import { Browser } from "@wailsio/runtime";
import type React from "react";
import { toast } from "sonner";

/*
 * Opens a URL in user's default browser using Wails Browser API.
 */
export async function openExternalLink(url: string): Promise<void> {
  try {
    return await Browser.OpenURL(url);
  } catch (err) {
    toast.error("Error opening link", { description: (err as Error).message });
  }
}

export function handleExternalLinkClick(event: React.MouseEvent<HTMLElement>): void {
  const target = event.target as HTMLElement | null;
  const anchor = target?.closest("a");
  if (!anchor) return;

  const href = anchor.getAttribute("href");
  if (!href || href.startsWith("#") || href.startsWith("javascript:")) return;

  event.preventDefault();

  void openExternalLink(href);
}
