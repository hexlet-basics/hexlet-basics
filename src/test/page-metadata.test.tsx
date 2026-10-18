import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
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
import { DESKTOP_VIEWPORT } from "@/test/viewport";

const THEORY_REPEATS = 10;

// What a search engine and a shared link see for the Course page and a Lesson:
// the head each route composes from its payload, read off the loaded match —
// the same tags the server renders into the HTML, since head() runs on the
// loader's data before the page is streamed.

const cover = "https://cdn.example.com/covers/javascript-list.png";

const course: Course = {
  categoryId: null,
  coverListVariant: cover,
  coverThumbVariant: null,
  createdAt: "2026-01-01T00:00:00Z",
  // The build's own name opens the lesson description, as legacy's did.
  currentVersion: {
    createdAt: "2026-01-01T00:00:00Z",
    id: 99,
    name: "JavaScript",
    result: null,
    state: "built",
  },
  currentVersionId: 99,
  enrollmentsCount: 0,
  hexletProgramLandingPage: null,
  id: 1,
  learnAs: null,
  lessonsCount: 1,
  name: "javascript",
  ratingCount: 0,
  ratingValue: 0,
  readiness: null,
  repositoryUrl: null,
  slug: "javascript",
};

const landingPage: CourseLandingPage = {
  courseId: 1,
  courseSlug: "javascript",
  createdAt: "2026-01-01T00:00:00Z",
  description: "Learn JavaScript in the browser",
  duration: 10,
  enrollmentsCount: 0,
  footer: true,
  footerName: null,
  header: "JavaScript course",
  id: 10,
  landingPageToRedirectId: null,
  listed: true,
  main: true,
  metaDescription: "Learn JavaScript from scratch in the browser",
  metaTitle: "Free JavaScript course",
  name: "JavaScript for beginners",
  order: null,
  outcomesDescription: null,
  outcomesHeader: null,
  outcomesImage: null,
  slug: "javascript",
  state: "published",
  usedInDescription: null,
  usedInHeader: null,
};

const progress: CourseProgress = {
  completion: 0,
  nextLesson: { id: 1002, name: "Variables", slug: "variables" },
  state: null,
};

const lessons = [
  {
    description: null,
    finished: false,
    id: 1002,
    name: "Variables",
    position: 1,
    slug: "variables",
  },
];

// Longer than the description's bound, with the line breaks and runs of spaces
// raw markdown has, so both the cut and the squash are visible.
const theory = `A variable is a name bound to a value.\n\n\`\`\`js\nlet   greeting = 'hello';\n\`\`\`\n\n${"Variables hold values that change. ".repeat(THEORY_REPEATS)}`;

const lessonView: CourseLessonView = {
  adjacent: { next: null, previous: null },
  landingPage,
  lesson: {
    course,
    createdAt: "2026-01-01T00:00:00Z",
    definitions: [],
    description: "Storing values",
    exerciseExtension: null,
    id: 1002,
    instructions: "Assign the string `hello` to `greeting`.",
    locale: "en",
    name: "Variables",
    naturalOrder: 1,
    originalCode: null,
    preparedCode: "let greeting = '';\n",
    slug: "variables",
    sourceCodeUrl: null,
    testCode: null,
    testExtension: null,
    theory,
    tips: [],
    version: 99,
    versionId: 99,
  },
  lessons,
  progress,
};

const courseView: CourseView = {
  course,
  enrollment: null,
  landingPage,
  lessons,
  modules: [],
  movedTo: null,
  progress,
  qnaItems: [],
};

// The player lays out by the viewport, and below the sm breakpoint it mounts
// the theory twice (its own pane and a workspace tab), so the lesson is read
// at a desktop viewport, restored when the test is done.
const startingViewport = { height: window.innerHeight, width: window.innerWidth };

describe("page metadata", () => {
  it("a lesson is titled, described and linked the way legacy did it", async () => {
    expect.hasAssertions();
    await page.viewport(DESKTOP_VIEWPORT.width, DESKTOP_VIEWPORT.height);
    try {
      worker.use(
        http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView)),
      );

      const { router } = await renderRoute(lessonRoute, {
        initialPath: "/languages/javascript/lessons/variables",
        path: "/{-$locale}/languages/$slug/lessons/$lessonSlug",
        wrap: (element) => <div style={{ height: "800px", width: "1200px" }}>{element}</div>,
      });

      await expect
        .element(page.getByText("A variable is a name bound to a value.", { exact: true }))
        .toBeVisible();
      const match = router.state.matches.at(-1);
      const meta = match?.meta;
      const url = `${globalThis.location.origin}/languages/javascript/lessons/variables`;

      // The lesson and the course's landing name, behind the site name; version,
      // lesson and raw theory, cut mid-word to 217 characters and an ellipsis,
      // then whitespace squashed, which is why it ends up shorter than 220; and
      // Open Graph as legacy emitted it for a lesson: an article, with no
      // og:description and no Twitter card.
      expect(meta).toStrictEqual(
        expect.arrayContaining([
          { title: "CodeBasics | Variables | JavaScript for beginners" },
          {
            content:
              "[JavaScript] — Variables — A variable is a name bound to a value. ```js let greeting = 'hello'; ``` " +
              "Variables hold values that change. Variables hold values that change. Variables hold values that change. Variable...",
            name: "description",
          },
          { content: "article", property: "og:type" },
          { content: "en", property: "og:locale" },
          { content: "Variables | JavaScript for beginners", property: "og:title" },
          { content: url, property: "og:url" },
          { content: cover, property: "og:image" },
        ]),
      );
      expect(meta).not.toStrictEqual(
        expect.arrayContaining([expect.objectContaining({ property: "og:description" })]),
      );
      const twitterName: unknown = expect.stringMatching(/^twitter:/u);
      expect(meta).not.toStrictEqual(
        expect.arrayContaining([expect.objectContaining({ name: twitterName })]),
      );
      expect(match?.links).toStrictEqual(
        expect.arrayContaining([
          { href: url, rel: "canonical" },
          { href: cover, rel: "image_src" },
        ]),
      );
    } finally {
      await page.viewport(startingViewport.width, startingViewport.height);
    }
  }, 60_000);

  it("the course page carries its landing page's meta title, description and canonical", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView)));

    const { router } = await renderRoute(courseRoute, {
      initialPath: "/languages/javascript",
      path: "/{-$locale}/languages/$slug",
    });

    await expect.element(page.getByRole("heading", { name: "JavaScript course" })).toBeVisible();
    const match = router.state.matches.at(-1);
    const meta = match?.meta;
    const url = `${globalThis.location.origin}/languages/javascript`;

    expect(meta).toStrictEqual(
      expect.arrayContaining([
        { title: "CodeBasics | Free JavaScript course" },
        { content: "Learn JavaScript from scratch in the browser", name: "description" },
        { content: "Free JavaScript course", property: "og:title" },
        { content: "Learn JavaScript from scratch in the browser", property: "og:description" },
        { content: "website", property: "og:type" },
        { content: "en", property: "og:locale" },
        { content: url, property: "og:url" },
        { content: cover, property: "og:image" },
        { content: "summary", name: "twitter:card" },
        { content: "@hexlet_io", name: "twitter:site" },
      ]),
    );
    expect(match?.links).toStrictEqual(
      expect.arrayContaining([
        { href: url, rel: "canonical" },
        { href: cover, rel: "image_src" },
      ]),
    );
  }, 60_000);
});
