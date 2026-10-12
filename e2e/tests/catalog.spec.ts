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
  await visit(page, "/ru/languages/javascript");

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
