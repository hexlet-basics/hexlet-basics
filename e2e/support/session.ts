import { expect } from "@playwright/test";
import { password } from "./people";
import type { Site } from "./site";

/** Signing up, in and out, and who the header says is signed in. */
export class Session {
  readonly #site: Site;

  constructor(site: Site) {
    this.#site = site;
  }

  async signUp(account: { firstName: string; email: string; password: string }): Promise<void> {
    const { page } = this.#site;
    await this.#site.visit("/users/new");
    await page.getByRole("textbox", { name: "First name" }).fill(account.firstName);
    await page.getByRole("textbox", { name: "Email" }).fill(account.email);
    await page.getByRole("textbox", { name: "Password" }).fill(account.password);
    await page.getByRole("button", { name: "Create an Account" }).click();
  }

  /** Signs in through the login form, as a person would. */
  async signIn(email: string, secret: string): Promise<void> {
    const { page } = this.#site;
    await this.#site.visit("/session/new");
    await page.getByRole("textbox", { name: "Email" }).fill(email);
    await page.getByRole("textbox", { name: "Password" }).fill(secret);
    await page.getByRole("button", { name: "Sign In" }).click();
  }

  /** Signs a fixture user in with the seeded password and sees their name in the header. */
  async signInAs(user: { email: string; name: string }): Promise<void> {
    await this.signIn(user.email, password);
    await this.expectSignedInAs(user.name);
  }

  async signOut(name: string): Promise<void> {
    await this.#site.header.getByRole("button", { name }).click();
    await this.#site.page.getByRole("menuitem", { name: "Sign out" }).click();
  }

  async expectSignedInAs(name: string): Promise<void> {
    await expect(this.#site.header.getByRole("button", { name })).toBeVisible();
  }

  async expectSignedOut(): Promise<void> {
    await expect(this.#site.header.getByRole("link", { name: "Sign in" })).toBeVisible();
  }

  /** The form refused the credentials: an alert, the same page, nobody signed in. */
  async expectSignInRefused(): Promise<void> {
    await expect(this.#site.page.getByRole("alert")).toBeVisible();
    await this.#site.expectPath("/session/new");
    await this.expectSignedOut();
  }
}
