import { lazy, type ComponentType, type LazyExoticComponent } from "react";

export type PreloadableComponent<T extends ComponentType<unknown>> = LazyExoticComponent<T> & {
  preload: () => Promise<{ default: T }>;
};

export function lazyWithPreload<T extends ComponentType<unknown>>(
  factory: () => Promise<{ default: T }>,
): PreloadableComponent<T> {
  const Component = lazy(factory) as PreloadableComponent<T>;
  Component.preload = factory;
  return Component;
}
