import { expect, test } from "@playwright/test";
import { admin, header, learner, signIn, visit } from "../support/app";

test("a guest is sent to sign in and returns to the admin afterwards", async ({ page }) => {
  await visit(page, "/admin");

  await expect(page).toHaveURL("/session/new?redirect=%2Fadmin");
  await expect(page.getByRole("heading", { level: 1, name: "Sign In" })).toBeVisible();
});

test("a signed-in user who is not an admin is refused", async ({ page }) => {
  await signIn(page, learner.email, learner.name);

  await visit(page, "/admin");

  await expect(page).toHaveURL("/");
  await expect(header(page).getByRole("button", { name: learner.name })).toBeVisible();
});

test("an admin sees the back office and its lists", async ({ page }) => {
  await signIn(page, admin.email, admin.name);

  await visit(page, "/admin");
  await expect(page).toHaveURL("/admin/courses");
  await expect(page.getByRole("heading", { name: "Courses" })).toBeVisible();
  await expect(page.getByRole("row", { name: /javascript/i }).first()).toBeVisible();

  await page.getByRole("link", { name: "Users", exact: true }).click();
  await expect(page).toHaveURL("/admin/users");
  await expect(page.getByRole("heading", { name: "Users" })).toBeVisible();
  await expect(page.getByRole("cell", { name: admin.email })).toBeVisible();
  await expect(page.getByRole("cell", { name: learner.email })).toBeVisible();
});
