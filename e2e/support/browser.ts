import { expect } from "@playwright/test";
import type { Locator, Page } from "@playwright/test";

/**
 * The page and what every concept needs from it. The step files in this
 * directory drive it; a spec never sees it, because `Site` keeps it in a
 * private field and `test` hands a scenario only the Site.
 */
export class Browser {
  readonly page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  /**
   * Opens a page and waits for it to settle. The server renders every page, but
   * the forms and menus only work once the client bundle has hydrated; a field
   * filled before that would be overwritten by the hydrating form. Network idle
   * after the page's last script has run is how the suite knows that is over.
   */
  async visit(path: string): Promise<void> {
    await this.page.goto(path);
    await this.settled();
  }

  /** Waits out a navigation the page started itself, such as the language switch. */
  async settled(): Promise<void> {
    await this.page.waitForLoadState("networkidle");
  }

  /** The site header: the signed-in user's name, or the sign-in link, lives here. */
  get header(): Locator {
    return this.page.getByRole("banner");
  }

  async expectPath(path: string): Promise<void> {
    await expect(this.page).toHaveURL(path);
  }

  async expectHeading(name: string | RegExp, level?: number): Promise<void> {
    await expect(this.page.getByRole("heading", { level, name })).toBeVisible();
  }

  /** Opens an address the site has no page for and sees the not-found page. */
  async expectNotFound(path: string): Promise<void> {
    const response = await this.page.goto(path);
    expect(response?.status()).toBe(404);
    await this.expectHeading("Page Not Found");
  }
}
