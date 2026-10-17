import { test } from "../support/test";

// The catalog, not a course page: a course page is addressed by a landing page
// slug, and a landing page belongs to one locale, so its URL has no twin in the
// other languages to switch to.
test("switching the language keeps the page and prefixes the route", async ({ site }) => {
  await site.catalog.open();
  await site.expectHeading(/./, 1);

  await site.locale.switch("English", "Русский");
  await site.expectPath("/ru/languages");
  await site.settled();
  await site.locale.expectLang("ru");
  await site.expectHeading(/./, 1);

  await site.locale.switch("Русский", "Español");
  await site.expectPath("/es/languages");
  await site.settled();
  await site.locale.expectLang("es");

  await site.locale.switch("Español", "English");
  await site.expectPath("/languages");
  await site.locale.expectLang("en");
});
