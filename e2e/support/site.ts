import type { Page } from "@playwright/test";
import { Admin } from "./admin";
import { Browser } from "./browser";
import { Catalog } from "./catalog";
import { Lesson } from "./lesson";
import { Locale } from "./locale";
import { Session } from "./session";

/**
 * The site as a person uses it, and the only handle a scenario gets: a spec
 * says what the person does and sees through the steps below and never touches
 * the page, a selector or a wait. The steps live in a file per concept
 * (`session.ts`, `catalog.ts`, `lesson.ts`, `admin.ts`, `locale.ts`) and share
 * one `Browser`, so a change in how the site is built is fixed in one place.
 *
 * The browser sits in a `#private` field: a spec cannot reach the page through
 * the Site, by type or at run time.
 */
export class Site {
  readonly session: Session;
  readonly catalog: Catalog;
  readonly lesson: Lesson;
  readonly admin: Admin;
  readonly locale: Locale;
  readonly #browser: Browser;

  constructor(page: Page) {
    this.#browser = new Browser(page);
    this.session = new Session(this.#browser);
    this.catalog = new Catalog(this.#browser);
    this.lesson = new Lesson(this.#browser);
    this.admin = new Admin(this.#browser);
    this.locale = new Locale(this.#browser);
  }

  /** The address the person is on. */
  async expectPath(path: string): Promise<void> {
    await this.#browser.expectPath(path);
  }

  /** A heading the person can see, of a level when the page's outline matters. */
  async expectHeading(name: string | RegExp, level?: number): Promise<void> {
    await this.#browser.expectHeading(name, level);
  }

  /** The person follows an address the site has no page for. */
  async expectNotFound(path: string): Promise<void> {
    await this.#browser.expectNotFound(path);
  }
}
