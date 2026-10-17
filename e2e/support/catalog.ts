import { expect } from "@playwright/test";
import type { Site } from "./site";

/** The course catalog, a course's page and the way from one to the other. */
export class Catalog {
  readonly #site: Site;

  constructor(site: Site) {
    this.#site = site;
  }

  async open(): Promise<void> {
    await this.#site.visit("/languages");
  }

  async openCoursePage(path: string): Promise<void> {
    await this.#site.visit(path);
  }

  /** A course's card is a link named by its title and blurb; `title` matches its start. */
  async expectCourseCard(title: string): Promise<void> {
    await expect(this.#site.page.getByRole("link", { name: courseCard(title) })).toBeVisible();
  }

  async openCourseCard(title: string): Promise<void> {
    await this.#site.page.getByRole("link", { name: courseCard(title) }).click();
  }

  async openLesson(name: string): Promise<void> {
    await this.#site.page.getByRole("link", { name }).click();
  }
}

function courseCard(title: string): RegExp {
  return new RegExp(`^${title}`);
}
