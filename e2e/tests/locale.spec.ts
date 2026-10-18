import { test } from "../support/test";

// The catalog, not a course page: a course page is addressed by a landing page
// slug, and a landing page belongs to one locale, so its URL has no twin in the
// other languages to switch to.
test("switching the language keeps the page and prefixes the route", async ({ site }) => {
  await site.catalog.open();
  await site.expectHeading(/./u, 1);

  await site.locale.switch("English", "Русский", { lang: "ru", path: "/ru/languages" });
  await site.expectHeading(/./u, 1);

  await site.locale.switch("Русский", "Español", { lang: "es", path: "/es/languages" });

  await site.locale.switch("Español", "English", { lang: "en", path: "/languages" });
});
