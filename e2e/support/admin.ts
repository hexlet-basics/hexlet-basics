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
      this.#browser.page.getByRole("row", { name: /javascript/i }).first(),
    ).toBeVisible();
  }

  async openUsers(): Promise<void> {
    await this.#browser.page.getByRole("link", { name: "Users", exact: true }).click();
  }

  async expectUserListed(email: string): Promise<void> {
    await expect(this.#browser.page.getByRole("cell", { name: email })).toBeVisible();
  }
}
