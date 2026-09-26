import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
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
  landingPages: [
    { id: 1, courseId: 1, slug: "ruby-ru", header: "Курс Ruby", locale: "ru" },
    { id: 2, courseId: 1, slug: "ruby-en", header: "Ruby Course", locale: "en" },
  ],
  blogPosts: [
    { id: 3, name: "Привет, мир", slug: "hello-ru", locale: "ru" },
    { id: 4, name: "Hello world", slug: "hello-en", locale: "en" },
    { id: 5, name: "Hola mundo", slug: "hola", locale: "es" },
  ],
  categories: [
    {
      id: 6,
      slug: "programming-ru",
      name: "Программирование",
      header: "Программирование",
      description: null,
      locale: "ru",
      createdAt: "2026-01-01T00:00:00Z",
    },
    {
      id: 7,
      slug: "programming-en",
      name: "Programming",
      header: "Programming",
      description: null,
      locale: "en",
      createdAt: "2026-01-01T00:00:00Z",
    },
  ],
};

function renderMap(initialPath: string) {
  worker.use(http.get("*/api/map", () => HttpResponse.json(sitemap)));
  return renderRoute(Route, { path: "/{-$locale}/map", initialPath });
}

test("the ru site lists both sites, each linking into its own locale", async () => {
  const { router } = await renderMap("/ru/map");

  await expect.element(page.getByRole("heading", { name: "Карта сайта", level: 1 })).toBeVisible();
  await expect.element(page.getByRole("link", { name: "Главная" })).toHaveAttribute("href", "/ru");
  await expect.element(page.getByRole("link", { name: "Home" })).toHaveAttribute("href", "/");

  await page.getByRole("button", { name: "Курсы по программированию" }).click();
  await page.getByRole("button", { name: "Courses", exact: true }).click();
  await expect
    .element(page.getByRole("link", { name: "Курс Ruby" }))
    .toHaveAttribute("href", "/ru/languages/ruby-ru");
  await expect
    .element(page.getByRole("link", { name: "Ruby Course" }))
    .toHaveAttribute("href", "/languages/ruby-en");

  await page.getByRole("button", { name: "Блог Code Basics" }).click();
  await page.getByRole("button", { name: "Blog", exact: true }).click();
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

  expect(router.state.matches.at(-1)?.meta).toEqual([{ title: "CodeBasics | Карта сайта" }]);
});

test("the en site has no sitemap", async () => {
  await renderMap("/map");

  await expect.element(page.getByRole("heading", { name: "Page Not Found" })).toBeVisible();
});

test("the es site has no sitemap", async () => {
  await renderMap("/es/map");

  await expect.element(page.getByRole("heading", { name: "Page Not Found" })).toBeVisible();
});
