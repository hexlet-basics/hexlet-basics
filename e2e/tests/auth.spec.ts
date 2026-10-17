import { expect, test } from "@playwright/test";
import { header, password, signIn, uniqueEmail, visit } from "../support/app";

test("a visitor signs up, signs out and signs in again", async ({ page }) => {
  const email = uniqueEmail();
  const secret = "s3cret-pass-e2e";

  await visit(page, "/users/new");
  await page.getByRole("textbox", { name: "First name" }).fill("Ada");
  await page.getByRole("textbox", { name: "Email" }).fill(email);
  await page.getByRole("textbox", { name: "Password" }).fill(secret);
  await page.getByRole("button", { name: "Create an Account" }).click();

  // Signing up also signs in: the header shows the new account.
  await expect(header(page).getByRole("button", { name: "Ada" })).toBeVisible();

  await header(page).getByRole("button", { name: "Ada" }).click();
  await page.getByRole("menuitem", { name: "Sign out" }).click();
  await expect(header(page).getByRole("link", { name: "Sign in" })).toBeVisible();

  await visit(page, "/session/new");
  await page.getByRole("textbox", { name: "Email" }).fill(email);
  await page.getByRole("textbox", { name: "Password" }).fill(secret);
  await page.getByRole("button", { name: "Sign In" }).click();
  await expect(header(page).getByRole("button", { name: "Ada" })).toBeVisible();
});

test("a wrong password is refused and nobody is signed in", async ({ page }) => {
  await visit(page, "/session/new");
  await page.getByRole("textbox", { name: "Email" }).fill("staff-editor@example.com");
  await page.getByRole("textbox", { name: "Password" }).fill(`${password}-wrong`);
  await page.getByRole("button", { name: "Sign In" }).click();

  await expect(page.getByRole("alert")).toBeVisible();
  await expect(page).toHaveURL("/session/new");
  await expect(header(page).getByRole("link", { name: "Sign in" })).toBeVisible();
});

test("a fixture learner signs in with the seeded password", async ({ page }) => {
  await signIn(page, "staff-editor@example.com", "Staff Editor");
  await expect(page).toHaveURL("/");
});
