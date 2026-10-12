import { expect, test } from "@playwright/test";
import { header, settled, visit } from "../support/app";

test("switching the language keeps the page and prefixes the route", async ({ page }) => {
  await visit(page, "/languages/javascript");
  await expect(page.getByRole("heading", { level: 1, name: "Курс JavaScript" })).toBeVisible();

  await header(page).getByRole("button", { name: "English" }).click();
  await page.getByRole("menuitem", { name: "Русский" }).click();
  await expect(page).toHaveURL("/ru/languages/javascript");
  await settled(page);
  await expect(page.locator("html")).toHaveAttribute("lang", "ru");
  await expect(page.getByRole("heading", { level: 1, name: "Курс JavaScript" })).toBeVisible();

  await header(page).getByRole("button", { name: "Русский" }).click();
  await page.getByRole("menuitem", { name: "Español" }).click();
  await expect(page).toHaveURL("/es/languages/javascript");
  await settled(page);
  await expect(page.locator("html")).toHaveAttribute("lang", "es");

  await header(page).getByRole("button", { name: "Español" }).click();
  await page.getByRole("menuitem", { name: "English" }).click();
  await expect(page).toHaveURL("/languages/javascript");
  await expect(page.locator("html")).toHaveAttribute("lang", "en");
});
