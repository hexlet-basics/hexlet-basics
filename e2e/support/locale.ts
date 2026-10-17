import { expect } from "@playwright/test";
import type { Browser } from "./browser";

/** The language switch in the header. */
export class Locale {
  readonly #browser: Browser;

  constructor(browser: Browser) {
    this.#browser = browser;
  }

  /**
   * Picks `to` in the menu whose button shows the current language `from`, and
   * sees the site arrive in it: the route under `path` and the page's `lang`.
   * The address is checked before the page is awaited, because the switch
   * navigates on its own and network idle right after the click could
   * otherwise pass before the navigation has started.
   */
  async switch(from: string, to: string, arrived: { path: string; lang: string }): Promise<void> {
    await this.#browser.header.getByRole("button", { name: from }).click();
    await this.#browser.page.getByRole("menuitem", { name: to }).click();
    await this.#browser.expectPath(arrived.path);
    await this.#browser.settled();
    await expect(this.#browser.page.locator("html")).toHaveAttribute("lang", arrived.lang);
  }
}
