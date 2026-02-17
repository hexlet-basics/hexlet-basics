import type { Locale } from "@/types";
import { fromWindow } from "./utils";

export type GonData = {
  // biome-ignore lint/suspicious/noExplicitAny: false positive
  [key: string]: any;
  suffix: "ru" | null;
  locale: Locale;
};

const isBrowser = typeof window !== "undefined";
const gonData = fromWindow("gon");

if (isBrowser && !gonData) {
  throw new Error("gon is not initialized");
}

export const gon = window.gon as GonData;
