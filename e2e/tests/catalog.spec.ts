import { expect, test } from "@playwright/test";
import { visit } from "../support/app";

test("the catalog lists the fixture courses", async ({ page }) => {
  await visit(page, "/languages");

  await expect(
    page.getByRole("heading", { level: 1, name: "Programming language courses" }),
  ).toBeVisible();
  await expect(page.getByRole("link", { name: /^Ruby Course/ })).toBeVisible();
});

test("a course page shows its landing copy and syllabus, and a lesson opens from it", async ({
  page,
}) => {
  await visit(page, "/ru/languages/javascript-ru");

  await expect(page.getByRole("heading", { level: 1, name: "Курс JavaScript" })).toBeVisible();
  await expect(
    page.getByRole("heading", { level: 2, name: "Программа обучения JavaScript" }),
  ).toBeVisible();

  await page.getByRole("link", { name: "Привет, мир!" }).click();

  await expect(page).toHaveURL("/ru/languages/javascript/lessons/hello-world");
  await expect(
    page.getByRole("heading", { level: 1, name: "JavaScript: Привет, мир!" }),
  ).toBeVisible();
});

// The one test that follows a link the site itself renders: the API resolves a
// course page by its landing page's slug, so a link built from anything else
// (the course's own slug) would be a 500 that no test with a faked API sees.
test("a course opens from its card in the catalog", async ({ page }) => {
  await visit(page, "/languages");

  await page.getByRole("link", { name: /^Ruby Course/ }).click();

  await expect(page).toHaveURL("/languages/ruby-en");
  await expect(page.getByRole("heading", { level: 1, name: "Ruby Course" })).toBeVisible();
});

test("a landing page slug of another locale is not found", async ({ page }) => {
  const response = await page.goto("/languages/ruby-ru");

  expect(response?.status()).toBe(404);
  await expect(page.getByRole("heading", { name: "Page Not Found" })).toBeVisible();
});

test("an archived landing page redirects to the course's main page", async ({ page }) => {
  await visit(page, "/ru/languages/javascript-ru-archived");

  await expect(page).toHaveURL("/ru/languages/javascript-ru");
  await expect(page.getByRole("heading", { level: 1, name: "Курс JavaScript" })).toBeVisible();
});
