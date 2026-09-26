import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type {
  Course,
  CourseLandingPage,
  CourseLessonView,
  CourseProgress,
  CourseView,
} from "@/client/types.gen";
import { Route as courseRoute } from "@/routes/{-$locale}/languages/$slug/index";
import { Route as lessonRoute } from "@/routes/{-$locale}/languages/$slug/lessons/$lessonSlug";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// What a search engine and a shared link see for the Course page and a Lesson:
// the head each route composes from its payload, read off the loaded match —
// the same tags the server renders into the HTML, since head() runs on the
// loader's data before the page is streamed.

const cover = "https://cdn.example.com/covers/javascript-list.png";

const course: Course = {
  id: 1,
  slug: "javascript",
  name: "javascript",
  learnAs: null,
  readiness: null,
  categoryId: null,
  currentVersionId: 99,
  // The build's own name opens the lesson description, as legacy's did.
  currentVersion: {
    id: 99,
    name: "JavaScript",
    result: null,
    state: "built",
    createdAt: "2026-01-01T00:00:00Z",
  },
  createdAt: "2026-01-01T00:00:00Z",
  enrollmentsCount: 0,
  lessonsCount: 1,
  ratingCount: 0,
  ratingValue: 0,
  repositoryUrl: null,
  hexletProgramLandingPage: null,
  coverListVariant: cover,
  coverThumbVariant: null,
};

const landingPage: CourseLandingPage = {
  id: 10,
  courseId: 1,
  courseSlug: "javascript",
  createdAt: "2026-01-01T00:00:00Z",
  slug: "javascript",
  name: "JavaScript for beginners",
  main: true,
  listed: true,
  state: "published",
  order: null,
  footer: true,
  footerName: null,
  landingPageToRedirectId: null,
  metaTitle: "Free JavaScript course",
  metaDescription: "Learn JavaScript from scratch in the browser",
  header: "JavaScript course",
  description: "Learn JavaScript in the browser",
  usedInHeader: null,
  usedInDescription: null,
  outcomesHeader: null,
  outcomesDescription: null,
  outcomesImage: null,
  duration: 10,
  enrollmentsCount: 0,
};

const progress: CourseProgress = {
  state: null,
  completion: 0,
  nextLessonSlug: "variables",
  furthestFinishedPosition: 0,
  lessons: [{ slug: "variables", position: 1, finished: false, available: true }],
};

const lessons = [{ id: 1002, name: "Variables", description: null, slug: "variables" }];

// Longer than the description's bound, with the line breaks and runs of spaces
// raw markdown has, so both the cut and the squash are visible.
const theory = `A variable is a name bound to a value.\n\n\`\`\`js\nlet   greeting = 'hello';\n\`\`\`\n\n${"Variables hold values that change. ".repeat(10)}`;

const lessonView: CourseLessonView = {
  lesson: {
    course,
    id: 1002,
    name: "Variables",
    slug: "variables",
    locale: "en",
    naturalOrder: 1,
    versionId: 99,
    version: 99,
    description: "Storing values",
    instructions: "Assign the string `hello` to `greeting`.",
    theory,
    definitions: [],
    tips: [],
    preparedCode: "let greeting = '';\n",
    originalCode: null,
    testCode: null,
    sourceCodeUrl: null,
    createdAt: "2026-01-01T00:00:00Z",
  },
  landingPage,
  lessons,
  progress,
};

const courseView: CourseView = {
  course,
  landingPage,
  lessons,
  enrollment: null,
  progress,
};

test("a lesson is titled, described and linked the way legacy did it", async () => {
  worker.use(
    http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView)),
  );

  const { router } = await renderRoute(lessonRoute, {
    path: "/{-$locale}/languages/$slug/lessons/$lessonSlug",
    initialPath: "/languages/javascript/lessons/variables",
    wrap: (element) => <div style={{ height: "800px", width: "1200px" }}>{element}</div>,
  });

  await expect
    .element(page.getByText("A variable is a name bound to a value.", { exact: true }))
    .toBeVisible();
  const match = router.state.matches.at(-1);
  const meta = match?.meta ?? [];
  const url = `${window.location.origin}/languages/javascript/lessons/variables`;

  // The lesson and the course's landing name, behind the site name.
  expect(meta).toContainEqual({ title: "CodeBasics | Variables | JavaScript for beginners" });

  // Version, lesson, raw theory — cut mid-word to 217 characters and an
  // ellipsis, then whitespace squashed, which is why it ends up shorter than 220.
  expect(meta).toContainEqual({
    name: "description",
    content:
      "[JavaScript] — Variables — A variable is a name bound to a value. ```js let greeting = 'hello'; ``` " +
      "Variables hold values that change. Variables hold values that change. Variables hold values that change. Variable...",
  });

  // Open Graph as legacy emitted it for a lesson: an article, with no
  // og:description and no Twitter card.
  expect(meta).toContainEqual({ property: "og:type", content: "article" });
  expect(meta).toContainEqual({ property: "og:locale", content: "en" });
  expect(meta).toContainEqual({
    property: "og:title",
    content: "Variables | JavaScript for beginners",
  });
  expect(meta).toContainEqual({ property: "og:url", content: url });
  expect(meta).toContainEqual({ property: "og:image", content: cover });
  expect(meta.some((tag) => tag?.property === "og:description")).toBe(false);
  expect(meta.some((tag) => tag?.name?.startsWith("twitter:"))).toBe(false);

  expect(match?.links).toContainEqual({ rel: "canonical", href: url });
  expect(match?.links).toContainEqual({ rel: "image_src", href: cover });
});

test("the course page carries its landing page's meta title, description and canonical", async () => {
  worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView)));

  const { router } = await renderRoute(courseRoute, {
    path: "/{-$locale}/languages/$slug",
    initialPath: "/languages/javascript",
  });

  await expect.element(page.getByRole("heading", { name: "JavaScript course" })).toBeVisible();
  const match = router.state.matches.at(-1);
  const meta = match?.meta ?? [];
  const url = `${window.location.origin}/languages/javascript`;

  expect(meta).toContainEqual({ title: "CodeBasics | Free JavaScript course" });
  expect(meta).toContainEqual({
    name: "description",
    content: "Learn JavaScript from scratch in the browser",
  });
  expect(meta).toContainEqual({ property: "og:title", content: "Free JavaScript course" });
  expect(meta).toContainEqual({
    property: "og:description",
    content: "Learn JavaScript from scratch in the browser",
  });
  expect(meta).toContainEqual({ property: "og:type", content: "website" });
  expect(meta).toContainEqual({ property: "og:locale", content: "en" });
  expect(meta).toContainEqual({ property: "og:url", content: url });
  expect(meta).toContainEqual({ property: "og:image", content: cover });
  expect(meta).toContainEqual({ name: "twitter:card", content: "summary" });
  expect(meta).toContainEqual({ name: "twitter:site", content: "@hexlet_io" });

  expect(match?.links).toContainEqual({ rel: "canonical", href: url });
  expect(match?.links).toContainEqual({ rel: "image_src", href: cover });
});
