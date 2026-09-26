import { http, HttpResponse } from "msw";
import { I18nextProvider } from "react-i18next";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type {
  Course,
  CourseCatalogItem,
  CourseCategory,
  CourseCategoryView,
  User,
} from "@/client/types.gen";
import { createI18n } from "@/lib/i18n";
import { Route as showRoute } from "@/routes/{-$locale}/language_categories/$slug";
import { Route as indexRoute } from "@/routes/{-$locale}/language_categories/index";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The public course categories, driven through their real routes with the API
// faked at the HTTP boundary. Asserted is what a visitor sees: the categories,
// a category's courses and questions, whether the consultation form is offered,
// and the page's head.

const user: User = {
  id: 1,
  firstName: "Alice",
  lastName: null,
  name: "Alice",
  email: "alice@example.com",
  admin: false,
  canAccessAdmin: false,
  assistantMessagesCount: null,
  createdAt: "2026-01-01T00:00:00Z",
  createdAtAsTimestamp: null,
  type: "user",
};

const programming: CourseCategory = {
  id: 1,
  slug: "programming",
  name: "Programming",
  header: "Programming",
  description: "Learn to program from scratch",
  locale: "en",
  createdAt: "2026-01-01T00:00:00Z",
};

const course: Course = {
  id: 1,
  slug: "ruby",
  name: "ruby",
  learnAs: null,
  readiness: null,
  categoryId: null,
  currentVersionId: null,
  currentVersion: null,
  createdAt: "2026-01-01T00:00:00Z",
  enrollmentsCount: 0,
  lessonsCount: 3,
  ratingCount: 0,
  ratingValue: 0,
  repositoryUrl: null,
  hexletProgramLandingPage: null,
  coverListVariant: null,
  coverThumbVariant: null,
};

const rubyLanding: CourseCatalogItem = {
  id: 100,
  slug: "ruby-en",
  header: "Ruby course",
  name: "Ruby",
  locale: "en",
  enrollmentsCount: 42,
  duration: 3,
  coverUrl: null,
  course,
};

const view: CourseCategoryView = {
  category: programming,
  url: "https://code-basics.com/language_categories/programming",
  landingPages: [rubyLanding],
  qnaItems: [{ id: 1, question: "Is experience required?", answer: "No, **none**." }],
};

function serveCategory() {
  worker.use(http.get("*/api/language_categories/programming", () => HttpResponse.json(view)));
}

// The ru copy, for the form legacy offered only on the ru site: the router's
// own i18n stays English, so this wraps the page in a Russian one.
function inRussian() {
  const ru = createI18n();
  void ru.changeLanguage("ru");
  return (element: React.ReactNode) => <I18nextProvider i18n={ru}>{element}</I18nextProvider>;
}

test("the index links every category the API returns", async () => {
  worker.use(http.get("*/api/language_categories", () => HttpResponse.json([programming])));

  const { router } = await renderRoute(indexRoute, {
    path: "/{-$locale}/language_categories/",
    initialPath: "/language_categories/",
  });

  await expect
    .element(page.getByRole("heading", { name: "Course categories", level: 1 }))
    .toBeVisible();
  await expect
    .element(page.getByRole("link", { name: /Programming/ }))
    .toHaveAttribute("href", "/language_categories/programming");

  const meta = router.state.matches.at(-1)?.meta;
  expect(meta).toContainEqual({ title: "CodeBasics | Course categories" });
  expect(meta).toContainEqual({ name: "twitter:site", content: "@hexlet_io" });
});

test("a category shows its courses and questions, with the legacy head", async () => {
  serveCategory();

  const { router } = await renderRoute(showRoute, {
    path: "/{-$locale}/language_categories/$slug",
    initialPath: "/language_categories/programming",
  });

  await expect.element(page.getByRole("heading", { name: "Programming", level: 1 })).toBeVisible();
  await expect.element(page.getByText("Learn to program from scratch")).toBeVisible();
  await expect.element(page.getByRole("heading", { name: "Ruby" })).toBeVisible();
  await expect.element(page.getByText("Is experience required?")).toBeVisible();
  await expect.element(page.getByText("none", { exact: true })).toBeVisible();

  const match = router.state.matches.at(-1);
  expect(match?.meta).toContainEqual({ title: "CodeBasics | Courses in the Programming category" });
  expect(match?.meta).toContainEqual({
    name: "description",
    content: expect.stringContaining("in the Programming category"),
  });
  expect(match?.links).toContainEqual({ rel: "canonical", href: view.url });
});

test("a signed-in visitor on the ru site is offered a consultation", async () => {
  serveCategory();

  await renderRoute(showRoute, {
    path: "/{-$locale}/language_categories/$slug",
    initialPath: "/ru/language_categories/programming",
    user,
    wrap: inRussian(),
  });

  await expect.element(page.getByText("Нужна помощь? Оставьте заявку, мы поможем")).toBeVisible();
  await expect.element(page.getByRole("button", { name: "Отправить" })).toBeVisible();
});

test("a guest is not offered a consultation", async () => {
  serveCategory();

  await renderRoute(showRoute, {
    path: "/{-$locale}/language_categories/$slug",
    initialPath: "/ru/language_categories/programming",
    wrap: inRussian(),
  });

  await expect.element(page.getByRole("heading", { name: "Programming", level: 1 })).toBeVisible();
  expect(page.getByText("Нужна помощь? Оставьте заявку, мы поможем").query()).toBeNull();
});

test("a category that is not in this locale is not found", async () => {
  worker.use(
    http.get("*/api/language_categories/programming", () =>
      HttpResponse.json({ status: 404, title: "Not Found" }, { status: 404 }),
    ),
  );

  await renderRoute(showRoute, {
    path: "/{-$locale}/language_categories/$slug",
    initialPath: "/language_categories/programming",
  });

  await expect.element(page.getByText("Not Found")).toBeVisible();
});
