import { expect, test } from "@playwright/test";
import { header, settled, visit } from "../support/app";

// The catalog, not a course page: a course page is addressed by a landing page
// slug, and a landing page belongs to one locale, so its URL has no twin in the
// other languages to switch to.
test("switching the language keeps the page and prefixes the route", async ({ page }) => {
  await visit(page, "/languages");
  await expect(page.getByRole("heading", { level: 1 })).toBeVisible();

  await header(page).getByRole("button", { name: "English" }).click();
  await page.getByRole("menuitem", { name: "Русский" }).click();
  await expect(page).toHaveURL("/ru/languages");
  await settled(page);
  await expect(page.locator("html")).toHaveAttribute("lang", "ru");
  await expect(page.getByRole("heading", { level: 1 })).toBeVisible();

  await header(page).getByRole("button", { name: "Русский" }).click();
  await page.getByRole("menuitem", { name: "Español" }).click();
  await expect(page).toHaveURL("/es/languages");
  await settled(page);
  await expect(page.locator("html")).toHaveAttribute("lang", "es");

  await header(page).getByRole("button", { name: "Español" }).click();
  await page.getByRole("menuitem", { name: "English" }).click();
  await expect(page).toHaveURL("/languages");
  await expect(page.locator("html")).toHaveAttribute("lang", "en");
});
