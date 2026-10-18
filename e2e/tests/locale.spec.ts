import { test } from "../support/test";

// The catalog, not a course page: a course page is addressed by a landing page
// slug, and a landing page belongs to one locale, so its URL has no twin in the
// other languages to switch to.
test("switching the language keeps the page and prefixes the route", async ({ site }) => {
  await site.catalog.open();
  await site.expectHeading(/./, 1);

  await site.locale.switch("English", "Русский", { path: "/ru/languages", lang: "ru" });
  await site.expectHeading(/./, 1);

  await site.locale.switch("Русский", "Español", { path: "/es/languages", lang: "es" });

  await site.locale.switch("Español", "English", { path: "/languages", lang: "en" });
});
