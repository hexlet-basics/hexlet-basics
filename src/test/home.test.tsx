import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { Course, CourseCatalogItem } from "@/client/types.gen";
import { Route as homeRoute } from "@/routes/{-$locale}/index";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The home page through its real route, inside the real site layout: the course
// catalog a visitor picks a course from, and the ru-only block of Hexlet programs.

function catalogItem(id: number, slug: string, name: string): CourseCatalogItem {
  const course: Course = {
    id,
    slug,
    name,
    learnAs: "first_language",
    readiness: "completed",
    categoryId: null,
    currentVersionId: null,
    currentVersion: null,
    createdAt: "2026-01-01T00:00:00Z",
    enrollmentsCount: 0,
    lessonsCount: 10,
    ratingCount: 0,
    ratingValue: 0,
    repositoryUrl: null,
    hexletProgramLandingPage: null,
    coverListVariant: null,
    coverThumbVariant: null,
  };
  return {
    id,
    slug,
    header: name,
    name,
    locale: "en",
    enrollmentsCount: 0,
    duration: 12,
    coverUrl: null,
    course,
  };
}

const courses = [catalogItem(1, "javascript", "JavaScript"), catalogItem(2, "python", "Python")];

test("lists every course the API returns, each linking to its page", async () => {
  worker.use(http.get("*/languages", () => HttpResponse.json(courses)));

  await renderRoute(homeRoute, { path: "/{-$locale}/", initialPath: "/" });

  const heading = page.getByRole("heading", { name: "Programming language courses", level: 1 });
  await expect.element(heading).toBeVisible();

  const main = page.getByRole("main");
  const javascript = main.getByRole("link", { name: /JavaScript/ }).first();
  await expect.element(javascript).toHaveAttribute("href", "/languages/javascript");
  const python = main.getByRole("link", { name: /Python/ }).first();
  await expect.element(python).toHaveAttribute("href", "/languages/python");
});

test("offers the Hexlet programs on the ru site only", async () => {
  worker.use(http.get("*/languages", () => HttpResponse.json(courses)));

  await renderRoute(homeRoute, { path: "/{-$locale}/", initialPath: "/ru" });

  const python = page.getByRole("link", { name: "Python", exact: true });
  await expect.element(python.first()).toBeVisible();
  const goProgram = page.getByRole("link", { name: /Go/ });
  await expect
    .element(goProgram.first())
    .toHaveAttribute("href", "https://ru.hexlet.io/programs/go");
});

test("answers a catalog that cannot be loaded with the server error page", async () => {
  worker.use(http.get("*/languages", () => HttpResponse.json({ status: 500 }, { status: 500 })));

  await renderRoute(homeRoute, { path: "/{-$locale}/", initialPath: "/" });

  const failure = page.getByRole("heading", { name: "Server Error" });
  await expect.element(failure).toBeVisible();
});
