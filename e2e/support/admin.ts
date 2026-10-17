import { expect } from "@playwright/test";
import type { Site } from "./site";

/** The back office. */
export class Admin {
  readonly #site: Site;

  constructor(site: Site) {
    this.#site = site;
  }

  async open(): Promise<void> {
    await this.#site.visit("/admin");
  }

  async expectCourseList(): Promise<void> {
    await this.#site.expectHeading("Courses");
    await expect(this.#site.page.getByRole("row", { name: /javascript/i }).first()).toBeVisible();
  }

  async openUsers(): Promise<void> {
    await this.#site.page.getByRole("link", { name: "Users", exact: true }).click();
  }

  async expectUserListed(email: string): Promise<void> {
    await expect(this.#site.page.getByRole("cell", { name: email })).toBeVisible();
  }
}
