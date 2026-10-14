import { expect, type Locator, type Page } from "@playwright/test";

/**
 * The public password scripts/e2e_db.sh gives the fixture users below. The
 * fixtures themselves carry no password digests.
 */
export const password = "e2e-password";

/** A fixture user who is an admin. */
export const admin = { email: "alice@example.com", name: "Alice Anderson" };

/** A fixture user who is not (an editor on the staff, so the back office is open to them). */
export const learner = { email: "bob@example.com", name: "Bob Brown" };

/** A fixture user with neither the admin flag nor a staff record. */
export const outsider = { email: "carol@example.com", name: "Carol Clark" };

/**
 * Opens a page and waits for it to settle. The server renders every page, but
 * the forms and menus only work once the client bundle has hydrated; a field
 * filled before that would be overwritten by the hydrating form. Network idle
 * after the page's last script has run is how the suite knows that is over.
 */
export async function visit(page: Page, path: string): Promise<void> {
  await page.goto(path);
  await settled(page);
}

/** Waits out a navigation the page started itself, such as the language switch. */
export async function settled(page: Page): Promise<void> {
  await page.waitForLoadState("networkidle");
}

/** The site header: the signed-in user's name, or the sign-in link, lives here. */
export function header(page: Page): Locator {
  return page.getByRole("banner");
}

/** Signs in through the login form, as a person would. */
export async function signIn(page: Page, email: string, name: string): Promise<void> {
  await visit(page, "/session/new");
  await page.getByRole("textbox", { name: "Email" }).fill(email);
  await page.getByRole("textbox", { name: "Password" }).fill(password);
  await page.getByRole("button", { name: "Sign In" }).click();
  await expect(header(page).getByRole("button", { name })).toBeVisible();
}

/** An email nobody has registered, so journeys that sign up never collide. */
export function uniqueEmail(): string {
  return `e2e-${Date.now()}-${Math.random().toString(36).slice(2, 8)}@example.com`;
}
