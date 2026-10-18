import { expect } from "@playwright/test";
import type { Browser } from "./browser";

/** The course catalog, a course's page and the way from one to the other. */
export class Catalog {
  readonly #browser: Browser;

  public constructor(browser: Browser) {
    this.#browser = browser;
  }

  public async open(): Promise<void> {
    await this.#browser.visit("/languages");
  }

  public async openCoursePage(path: string): Promise<void> {
    await this.#browser.visit(path);
  }

  /** A course's card is a link named by its title and blurb; `title` matches its start. */
  public async expectCourseCard(title: string): Promise<void> {
    await expect(this.#browser.page.getByRole("link", { name: courseCard(title) })).toBeVisible();
  }

  public async openCourseCard(title: string): Promise<void> {
    await this.#browser.page.getByRole("link", { name: courseCard(title) }).click();
  }

  public async openLesson(name: string): Promise<void> {
    await this.#browser.page.getByRole("link", { name }).click();
  }
}

function courseCard(title: string): RegExp {
  return new RegExp(`^${title}`, "u");
}
