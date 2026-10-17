import { expect, type Locator, type Page } from "@playwright/test";
import { Admin } from "./admin";
import { Catalog } from "./catalog";
import { Lesson } from "./lesson";
import { Locale } from "./locale";
import { Session } from "./session";

/**
 * The site as a person uses it, and the only handle a scenario gets: a spec
 * says what the person does and sees through the steps below and never touches
 * the page, a selector or a wait. The steps live in a file per concept
 * (`session.ts`, `catalog.ts`, `lesson.ts`, `admin.ts`, `locale.ts`) and share
 * this class's page, so a change in how the site is built is fixed in one place.
 */
export class Site {
  readonly session: Session;
  readonly catalog: Catalog;
  readonly lesson: Lesson;
  readonly admin: Admin;
  readonly locale: Locale;
  readonly #page: Page;

  constructor(page: Page) {
    this.#page = page;
    this.session = new Session(this);
    this.catalog = new Catalog(this);
    this.lesson = new Lesson(this);
    this.admin = new Admin(this);
    this.locale = new Locale(this);
  }

  /**
   * The page behind the steps. It is public so the concept files can drive it,
   * and out of reach of a spec by type: `test` in `test.ts` hands a scenario
   * only the Site.
   */
  get page(): Page {
    return this.#page;
  }

  /**
   * Opens a page and waits for it to settle. The server renders every page, but
   * the forms and menus only work once the client bundle has hydrated; a field
   * filled before that would be overwritten by the hydrating form. Network idle
   * after the page's last script has run is how the suite knows that is over.
   */
  async visit(path: string): Promise<void> {
    await this.#page.goto(path);
    await this.settled();
  }

  /** Waits out a navigation the page started itself, such as the language switch. */
  async settled(): Promise<void> {
    await this.#page.waitForLoadState("networkidle");
  }

  /** The site header: the signed-in user's name, or the sign-in link, lives here. */
  get header(): Locator {
    return this.#page.getByRole("banner");
  }

  async expectPath(path: string): Promise<void> {
    await expect(this.#page).toHaveURL(path);
  }

  async expectHeading(name: string | RegExp, level?: number): Promise<void> {
    await expect(this.#page.getByRole("heading", { level, name })).toBeVisible();
  }

  /** Opens an address the site has no page for and sees the not-found page. */
  async expectNotFound(path: string): Promise<void> {
    const response = await this.#page.goto(path);
    expect(response?.status()).toBe(404);
    await this.expectHeading("Page Not Found");
  }
}
