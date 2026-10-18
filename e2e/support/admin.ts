import { expect } from "@playwright/test";
import type { Browser } from "./browser";

/** The back office. */
export class Admin {
  readonly #browser: Browser;

  constructor(browser: Browser) {
    this.#browser = browser;
  }

  async open(): Promise<void> {
    await this.#browser.visit("/admin");
  }

  async expectCourseList(): Promise<void> {
    await this.#browser.expectHeading("Courses");
    await expect(
      this.#browser.page.getByRole("row", { name: /javascript/iu }).first(),
    ).toBeVisible();
  }

  async openUsers(): Promise<void> {
    await this.#browser.page.getByRole("link", { exact: true, name: "Users" }).click();
  }

  async expectUserListed(email: string): Promise<void> {
    // The list is paginated, newest first, so an old fixture user is not on its
    // first page: find the user the way an administrator does, by searching.
    await this.#browser.page.getByRole("textbox").fill(email);
    await expect(this.#browser.page.getByRole("cell", { name: email })).toBeVisible();
  }
}
