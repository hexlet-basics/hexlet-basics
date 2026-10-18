import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { Sitemap } from "@/client/types.gen";
import { Route } from "@/routes/{-$locale}/map";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The HTML sitemap, driven through its real route with the API faked at the
// HTTP boundary. Asserted is what a visitor sees: a section per site, ru then
// en, each in its own language and linking into its own locale, and that only
// the ru site has the page.

const sitemap: Sitemap = {
  blogPosts: [
    { id: 3, locale: "ru", name: "Привет, мир", slug: "hello-ru" },
    { id: 4, locale: "en", name: "Hello world", slug: "hello-en" },
    { id: 5, locale: "es", name: "Hola mundo", slug: "hola" },
  ],
  categories: [
    {
      createdAt: "2026-01-01T00:00:00Z",
      description: null,
      header: "Программирование",
      id: 6,
      locale: "ru",
      name: "Программирование",
      slug: "programming-ru",
    },
    {
      createdAt: "2026-01-01T00:00:00Z",
      description: null,
      header: "Programming",
      id: 7,
      locale: "en",
      name: "Programming",
      slug: "programming-en",
    },
  ],
  landingPages: [
    { courseId: 1, header: "Курс Ruby", id: 1, locale: "ru", slug: "ruby-ru" },
    { courseId: 1, header: "Ruby Course", id: 2, locale: "en", slug: "ruby-en" },
  ],
};

function renderMap(initialPath: string) {
  worker.use(http.get("*/api/map", () => HttpResponse.json(sitemap)));
  return renderRoute(Route, { initialPath, path: "/{-$locale}/map" });
}

describe("sitemap", () => {
  it("the ru site lists both sites, each linking into its own locale", async () => {
    expect.hasAssertions();
    const { router } = await renderMap("/ru/map");

    await expect
      .element(page.getByRole("heading", { level: 1, name: "Карта сайта" }))
      .toBeVisible();
    await expect
      .element(page.getByRole("link", { name: "Главная" }))
      .toHaveAttribute("href", "/ru");
    await expect.element(page.getByRole("link", { name: "Home" })).toHaveAttribute("href", "/");

    await page.getByRole("button", { name: "Курсы по программированию" }).click();
    await page.getByRole("button", { exact: true, name: "Courses" }).click();
    await expect
      .element(page.getByRole("link", { name: "Курс Ruby" }))
      .toHaveAttribute("href", "/ru/languages/ruby-ru");
    await expect
      .element(page.getByRole("link", { name: "Ruby Course" }))
      .toHaveAttribute("href", "/languages/ruby-en");

    await page.getByRole("button", { name: "Блог Code Basics" }).click();
    await page.getByRole("button", { exact: true, name: "Blog" }).click();
    await expect
      .element(page.getByRole("link", { name: "Привет, мир" }))
      .toHaveAttribute("href", "/ru/blog_posts/hello-ru");
    await expect
      .element(page.getByRole("link", { name: "Hello world" }))
      .toHaveAttribute("href", "/blog_posts/hello-en");
    expect(page.getByText("Hola mundo").query()).toBeNull();

    await page.getByRole("button", { name: "Категории курсов по программированию" }).click();
    await expect
      .element(page.getByRole("link", { name: "Программирование" }))
      .toHaveAttribute("href", "/ru/language_categories/programming-ru");

    expect(router.state.matches.at(-1)?.meta).toStrictEqual([
      { title: "CodeBasics | Карта сайта" },
    ]);
  }, 60_000);

  it("the en site has no sitemap", async () => {
    expect.hasAssertions();
    await renderMap("/map");

    await expect.element(page.getByRole("heading", { name: "Page Not Found" })).toBeVisible();
  }, 60_000);

  it("the es site has no sitemap", async () => {
    expect.hasAssertions();
    await renderMap("/es/map");

    await expect.element(page.getByRole("heading", { name: "Page Not Found" })).toBeVisible();
  }, 60_000);
});
