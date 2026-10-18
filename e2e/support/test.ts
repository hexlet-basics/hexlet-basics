import { test as base } from "@playwright/test";
import type { TestType } from "@playwright/test";
import { Site } from "./site";

/**
 * Playwright's `test` with one fixture, `site`, and the type narrowed to it: a
 * scenario that destructures `page` fails `tsc -b`, so the spec stays free of
 * selectors and waits and only says what the person does and sees.
 */
const withSite = base.extend<{ site: Site }>({
  site: async ({ page }, use) => {
    await use(new Site(page));
  },
});

export const test = withSite as unknown as TestType<{ site: Site }, object>;
