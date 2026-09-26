import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import { Route as postRoute } from "@/routes/{-$locale}/blog_posts/$slug";
import { Route as forTeachersRoute } from "@/routes/{-$locale}/cases/for_teachers";
import { Route as casesRoute } from "@/routes/{-$locale}/cases/index";
import { Route as pageRoute } from "@/routes/{-$locale}/pages/$id";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The data-free pages (ADR-0015) through their real routes: a static page, the
// ru-only cases and the error pages. Asserted is what a visitor and a crawler
// read: the body, the head, and which page answers.

function renderPage(initialPath: string, locale?: "ru" | "es") {
  return renderRoute(pageRoute, { path: "/{-$locale}/pages/$id", initialPath, locale });
}

test("the about page shows its body under its title and a canonical to itself", async () => {
  const { router } = await renderPage("/pages/about");

  await expect.element(page.getByRole("heading", { name: "About", level: 1 })).toBeVisible();
  await expect.element(page.getByRole("heading", { name: "What is Code Basics?" })).toBeVisible();

  const match = router.state.matches.at(-1);
  expect(match?.meta).toContainEqual({ title: "CodeBasics | About" });
  expect(match?.meta).toContainEqual({ name: "twitter:card", content: "summary" });
  expect(match?.links).toContainEqual({
    rel: "canonical",
    href: `${window.location.origin}/pages/about`,
  });
});

test("a legal page carries no canonical", async () => {
  const { router } = await renderPage("/ru/pages/tos", "ru");

  await expect
    .element(page.getByRole("heading", { name: "Условия использования образовательной платформы" }))
    .toBeVisible();
  expect(router.state.matches.at(-1)?.links ?? []).toEqual([]);
});

test("an es page reads the ru body under the en title, as legacy did", async () => {
  await renderPage("/es/pages/about", "es");

  await expect.element(page.getByRole("heading", { name: "About", level: 1 })).toBeVisible();
  await expect.element(page.getByRole("heading", { name: "Что такое Code Basics" })).toBeVisible();
});

test("an unknown page is the 404 page", async () => {
  await renderPage("/pages/nope");

  await expect.element(page.getByRole("heading", { name: "Page Not Found" })).toBeVisible();
});

test("the teachers' case asks a visitor to sign up and is its own canonical", async () => {
  const { router } = await renderRoute(forTeachersRoute, {
    path: "/{-$locale}/cases/for_teachers",
    initialPath: "/ru/cases/for_teachers",
    locale: "ru",
  });

  await expect
    .element(
      page.getByRole("heading", {
        name: "Эффективное обучение школьников и студентов программированию",
      }),
    )
    .toBeVisible();
  await expect
    .element(page.getByRole("link", { name: "Зарегистрироваться" }))
    .toHaveAttribute("href", "/ru/users/new");

  const match = router.state.matches.at(-1);
  expect(match?.meta).toContainEqual({
    title: "CodeBasics | Обучение программированию школьников и студентов",
  });
  expect(match?.links).toContainEqual({
    rel: "canonical",
    href: `${window.location.origin}/ru/cases/for_teachers`,
  });
});

test("the cases list links to the teachers' case", async () => {
  await renderRoute(casesRoute, {
    path: "/{-$locale}/cases/",
    initialPath: "/ru/cases/",
    locale: "ru",
  });

  await expect
    .element(page.getByRole("link", { name: "Перейти" }))
    .toHaveAttribute("href", "/ru/cases/for_teachers");
});

test("the cases answer only in ru", async () => {
  await renderRoute(casesRoute, { path: "/{-$locale}/cases/", initialPath: "/cases/" });

  await expect.element(page.getByRole("heading", { name: "Page Not Found" })).toBeVisible();
});

test("a load that fails is the 500 page", async () => {
  worker.use(
    http.get("*/api/blog_posts/hello-world", () =>
      HttpResponse.json({ status: 500, title: "Internal Server Error" }, { status: 500 }),
    ),
  );

  await renderRoute(postRoute, {
    path: "/{-$locale}/blog_posts/$slug",
    initialPath: "/blog_posts/hello-world",
  });

  await expect.element(page.getByRole("heading", { name: "Server Error" })).toBeVisible();
});
