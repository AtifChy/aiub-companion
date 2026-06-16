export const IS_WAILS = typeof window !== "undefined" && "wails" in window;
export const IS_DESKTOP = IS_WAILS;
export const IS_WEB = !IS_WAILS;
