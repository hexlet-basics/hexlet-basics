import { expect } from "@playwright/test";
import type { Site } from "./site";

/** The language switch in the header. */
export class Locale {
  readonly #site: Site;

  constructor(site: Site) {
    this.#site = site;
  }

  /** Picks `to` in the menu whose button shows the current language `from`. */
  async switch(from: string, to: string): Promise<void> {
    await this.#site.header.getByRole("button", { name: from }).click();
    await this.#site.page.getByRole("menuitem", { name: to }).click();
  }

  async expectLang(code: string): Promise<void> {
    await expect(this.#site.page.locator("html")).toHaveAttribute("lang", code);
  }
}
