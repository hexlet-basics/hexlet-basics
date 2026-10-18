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
  nextLesson: { id: 1002, slug: "variables", name: "Variables" },
};

const lessons = [
  {
    id: 1002,
    name: "Variables",
    description: null,
    slug: "variables",
    position: 1,
    finished: false,
  },
];

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
    exerciseExtension: null,
    testExtension: null,
    sourceCodeUrl: null,
    createdAt: "2026-01-01T00:00:00Z",
  },
  landingPage,
  lessons,
  adjacent: { previous: null, next: null },
  progress,
};

const courseView: CourseView = {
  course,
  landingPage,
  movedTo: null,
  lessons,
  modules: [],
  qnaItems: [],
  enrollment: null,
  progress,
};

// The player lays out by the viewport, and below the sm breakpoint it mounts
// the theory twice (its own pane and a workspace tab), so the lesson is read
// at a desktop viewport, restored when the test is done.
const startingViewport = { width: window.innerWidth, height: window.innerHeight };

describe("page metadata", () => {
  it("a lesson is titled, described and linked the way legacy did it", async () => {
    expect.hasAssertions();
    await page.viewport(1200, 800);
    try {
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
      const meta = match?.meta;
      const url = `${window.location.origin}/languages/javascript/lessons/variables`;

      // The lesson and the course's landing name, behind the site name; version,
      // lesson and raw theory, cut mid-word to 217 characters and an ellipsis,
      // then whitespace squashed, which is why it ends up shorter than 220; and
      // Open Graph as legacy emitted it for a lesson: an article, with no
      // og:description and no Twitter card.
      expect(meta).toStrictEqual(
        expect.arrayContaining([
          { title: "CodeBasics | Variables | JavaScript for beginners" },
          {
            name: "description",
            content:
              "[JavaScript] — Variables — A variable is a name bound to a value. ```js let greeting = 'hello'; ``` " +
              "Variables hold values that change. Variables hold values that change. Variables hold values that change. Variable...",
          },
          { property: "og:type", content: "article" },
          { property: "og:locale", content: "en" },
          { property: "og:title", content: "Variables | JavaScript for beginners" },
          { property: "og:url", content: url },
          { property: "og:image", content: cover },
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
          { rel: "canonical", href: url },
          { rel: "image_src", href: cover },
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
      path: "/{-$locale}/languages/$slug",
      initialPath: "/languages/javascript",
    });

    await expect.element(page.getByRole("heading", { name: "JavaScript course" })).toBeVisible();
    const match = router.state.matches.at(-1);
    const meta = match?.meta;
    const url = `${window.location.origin}/languages/javascript`;

    expect(meta).toStrictEqual(
      expect.arrayContaining([
        { title: "CodeBasics | Free JavaScript course" },
        { name: "description", content: "Learn JavaScript from scratch in the browser" },
        { property: "og:title", content: "Free JavaScript course" },
        { property: "og:description", content: "Learn JavaScript from scratch in the browser" },
        { property: "og:type", content: "website" },
        { property: "og:locale", content: "en" },
        { property: "og:url", content: url },
        { property: "og:image", content: cover },
        { name: "twitter:card", content: "summary" },
        { name: "twitter:site", content: "@hexlet_io" },
      ]),
    );
    expect(match?.links).toStrictEqual(
      expect.arrayContaining([
        { rel: "canonical", href: url },
        { rel: "image_src", href: cover },
      ]),
    );
  }, 60_000);
});
