import type { MantineColorSchemeManager } from "@mantine/core";
import Cookies from "js-cookie";

const colorSchemeCookieKey = "codebasics-color-scheme";
const cookieOptions = {
  expires: 365,
  path: "/",
  sameSite: "lax" as const,
};

export function cookieColorSchemeManager(): MantineColorSchemeManager {
  return {
    get: (defaultValue) => {
      const colorScheme = Cookies.get(colorSchemeCookieKey);

      return colorScheme === "dark" || colorScheme === "light"
        ? colorScheme
        : defaultValue;
    },
    set: (value) => {
      Cookies.set(
        colorSchemeCookieKey,
        value === "dark" ? "dark" : "light",
        cookieOptions,
      );
    },
    subscribe: () => {},
    unsubscribe: () => {},
    clear: () => {
      Cookies.remove(colorSchemeCookieKey, { path: "/" });
    },
  };
}
