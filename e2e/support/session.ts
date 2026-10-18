import { expect } from "@playwright/test";
import { password } from "./people";
import type { Browser } from "./browser";

/** Signing up, in and out, and who the header says is signed in. */
export class Session {
  readonly #browser: Browser;

  public constructor(browser: Browser) {
    this.#browser = browser;
  }

  public async signUp(account: {
    firstName: string;
    email: string;
    password: string;
  }): Promise<void> {
    const { page } = this.#browser;
    await this.#browser.visit("/users/new");
    await page.getByRole("textbox", { name: "First name" }).fill(account.firstName);
    await page.getByRole("textbox", { name: "Email" }).fill(account.email);
    await page.getByRole("textbox", { name: "Password" }).fill(account.password);
    await page.getByRole("button", { name: "Create an Account" }).click();
  }

  /** Signs in through the login form, as a person would. */
  public async signIn(email: string, secret: string): Promise<void> {
    const { page } = this.#browser;
    await this.#browser.visit("/session/new");
    await page.getByRole("textbox", { name: "Email" }).fill(email);
    await page.getByRole("textbox", { name: "Password" }).fill(secret);
    await page.getByRole("button", { name: "Sign In" }).click();
  }

  /** Signs a fixture user in with the seeded password and sees their name in the header. */
  public async signInAs(user: { email: string; name: string }): Promise<void> {
    await this.signIn(user.email, password);
    await this.expectSignedInAs(user.name);
  }

  public async signOut(name: string): Promise<void> {
    await this.#browser.header.getByRole("button", { name }).click();
    await this.#browser.page.getByRole("menuitem", { name: "Sign out" }).click();
  }

  public async expectSignedInAs(name: string): Promise<void> {
    await expect(this.#browser.header.getByRole("button", { name })).toBeVisible();
  }

  public async expectSignedOut(): Promise<void> {
    await expect(this.#browser.header.getByRole("link", { name: "Sign in" })).toBeVisible();
  }

  /** The form refused the credentials: an alert, the same page, nobody signed in. */
  public async expectSignInRefused(): Promise<void> {
    await expect(this.#browser.page.getByRole("alert")).toBeVisible();
    await this.#browser.expectPath("/session/new");
    await this.expectSignedOut();
  }
}
