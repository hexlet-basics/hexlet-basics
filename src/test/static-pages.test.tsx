import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import { ForbiddenPage } from "@/components/ErrorPage";
import { Route as postRoute } from "@/routes/{-$locale}/blog_posts/$slug";
import { Route as forTeachersRoute } from "@/routes/{-$locale}/cases/for_teachers";
import { Route as casesRoute } from "@/routes/{-$locale}/cases/index";
import { Route as pageRoute } from "@/routes/{-$locale}/pages/$id";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";
import { renderWithProviders } from "@/test/renderWithProviders";

const MIN_BODY_LENGTH = 500;

// The data-free pages (ADR-0015) through their real routes: a static page, the
// ru-only cases and the error pages. Asserted is what a visitor and a crawler
// read: the body, the head, and which page answers.

function renderPage(initialPath: string, locale?: "ru" | "es") {
  return renderRoute(pageRoute, { path: "/{-$locale}/pages/$id", initialPath, locale });
}

describe("static pages", () => {
  it("the about page shows its body under its title and a canonical to itself", async () => {
    expect.hasAssertions();
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
  }, 60_000);

  it("a legal page carries no canonical", async () => {
    expect.hasAssertions();
    const { router } = await renderPage("/ru/pages/tos", "ru");

    await expect
      .element(
        page.getByRole("heading", { name: "Условия использования образовательной платформы" }),
      )
      .toBeVisible();
    expect(router.state.matches.at(-1)?.links).toBeOneOf([undefined, []]);
  }, 60_000);

  it("an es page reads the ru body under the en title, as legacy did", async () => {
    expect.hasAssertions();
    await renderPage("/es/pages/about", "es");

    await expect.element(page.getByRole("heading", { name: "About", level: 1 })).toBeVisible();
    await expect
      .element(page.getByRole("heading", { name: "Что такое Code Basics" }))
      .toBeVisible();
  }, 60_000);

  it("an unknown page is the 404 page", async () => {
    expect.hasAssertions();
    await renderPage("/pages/nope");

    await expect.element(page.getByRole("heading", { name: "Page Not Found" })).toBeVisible();
  }, 60_000);

  it("the 403 page says access is denied", async () => {
    expect.hasAssertions();
    await renderWithProviders(<ForbiddenPage />);

    await expect.element(page.getByRole("heading", { name: "Access Denied" })).toBeVisible();
    await expect.element(page.getByText("403", { exact: true })).toBeVisible();
  }, 60_000);

  it("the teachers' case asks a visitor to sign up and is its own canonical", async () => {
    expect.hasAssertions();
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
  }, 60_000);

  it("the cases list links to the teachers' case", async () => {
    expect.hasAssertions();
    await renderRoute(casesRoute, {
      path: "/{-$locale}/cases/",
      initialPath: "/ru/cases/",
      locale: "ru",
    });

    await expect
      .element(page.getByRole("link", { name: "Перейти" }))
      .toHaveAttribute("href", "/ru/cases/for_teachers");
  }, 60_000);

  it("the cases answer only in ru", async () => {
    expect.hasAssertions();
    await renderRoute(casesRoute, { path: "/{-$locale}/cases/", initialPath: "/cases/" });

    await expect.element(page.getByRole("heading", { name: "Page Not Found" })).toBeVisible();
  }, 60_000);

  it("a load that fails is the 500 page", async () => {
    expect.hasAssertions();
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
  }, 60_000);

  // Every page body, in both languages it is written in, through the real route:
  // the title heading and a body that is not empty. The legal texts are long
  // hardcoded markup, so a broken tag or an import that no longer resolves is
  // what this catches.
  const PAGE_IDS = ["about", "authors", "cookie_policy", "privacy", "tos"];
  const LOCALES = [
    { locale: undefined, prefix: "" },
    { locale: "ru" as const, prefix: "/ru" },
  ];

  it.each(PAGE_IDS.flatMap((id) => LOCALES.map((l) => ({ id, ...l }))))(
    "the $id page renders its body in $prefix",
    async ({ id, locale, prefix }) => {
      expect.hasAssertions();
      await renderPage(`${prefix}/pages/${id}`, locale);

      const title = page.getByRole("heading", { level: 1 });
      await expect.element(title).toBeVisible();
      await expect.poll(() => document.body.textContent.length).toBeGreaterThan(MIN_BODY_LENGTH);
    },
    60_000,
  );
});
