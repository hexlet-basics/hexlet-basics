import { HttpResponse, http } from "msw";
import { describe, expect, it, onTestFinished } from "vitest";
import { page } from "vitest/browser";
import { getCourseLessonQueryKey } from "@/client/@tanstack/react-query.gen";
import type {
  Course,
  CourseLandingPage,
  CourseLessonView,
  CourseProgress,
  CourseView,
  LearnerLesson,
  LessonRef,
} from "@/client/types.gen";
import type { AuthUser } from "@/lib/auth";
import { Route as courseRoute } from "@/routes/{-$locale}/languages/$slug/index";
import { Route as lessonRoute } from "@/routes/{-$locale}/languages/$slug/lessons/$lessonSlug";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";
import { DESKTOP_VIEWPORT } from "@/test/viewport";
import type { RenderedRoute } from "@/test/renderRoute";

// A lesson page is headed with the course it belongs to.
const COURSE = "JavaScript";

const HELLO_WORLD_ID = 1001;
const VARIABLES_ID = 1002;
const STRINGS_ID = 1003;
const FINISHED_LESSONS = 3;

// The Course page, driven through its real route with the API faked at the HTTP
// boundary. What is asserted is what a visitor sees — the landing copy, the
// marks in the lesson list, the one button — and the requests pressing it made.

const course: Course = {
  categoryId: null,
  coverListVariant: null,
  coverThumbVariant: null,
  createdAt: "2026-01-01T00:00:00Z",
  currentVersion: null,
  currentVersionId: 99,
  enrollmentsCount: 1234,
  hexletProgramLandingPage: null,
  id: 1,
  learnAs: null,
  lessonsCount: 3,
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
  metaDescription: "Learn JavaScript",
  metaTitle: "JavaScript course",
  name: "JavaScript",
  order: null,
  outcomesDescription: null,
  outcomesHeader: null,
  outcomesImage: null,
  slug: "javascript-ru",
  state: "published",
  usedInDescription: "Browsers, servers and everything in between",
  usedInHeader: "Where JavaScript is used",
};

// One standing, as the server sends it: the course's own progress, and which of
// its lessons carry a check. The page is expected to render them as given.
interface Standing {
  progress: CourseProgress;
  finished: string[];
}

const lessonRef = (id: number, slug: string, name: string): LessonRef => ({ id, name, slug });

// The course's lessons in course order, each marked or not.
const lessonRows = (finished: string[]): LearnerLesson[] => [
  {
    description: null,
    finished: finished.includes("hello-world"),
    id: 1001,
    name: "Hello, World!",
    position: 1,
    slug: "hello-world",
  },
  {
    description: null,
    finished: finished.includes("variables"),
    id: 1002,
    name: "Variables",
    position: 2,
    slug: "variables",
  },
  {
    description: null,
    finished: finished.includes("strings"),
    id: 1003,
    name: "Strings",
    position: 3,
    slug: "strings",
  },
];

// Three payloads the server can send, each a fact about where the visitor
// stands.

// Has finished the first lesson: the second is next.
const midCourse: Standing = {
  finished: ["hello-world"],
  progress: {
    completion: 33,
    nextLesson: lessonRef(VARIABLES_ID, "variables", "Variables"),
    state: "started",
  },
};

// Has finished everything: there is no Next Lesson.
const finished: Standing = {
  finished: ["hello-world", "variables", "strings"],
  progress: { completion: 100, nextLesson: null, state: "finished" },
};

// Has never been here: nothing finished, the first lesson is next.
const firstVisit: Standing = {
  finished: [],
  progress: {
    completion: 0,
    nextLesson: lessonRef(HELLO_WORLD_ID, "hello-world", "Hello, World!"),
    state: null,
  },
};

// `enrolled` is a learner who has started the course: only they have a completion
// to show on its landing page.
function courseView(standing: Standing, enrolled = false): CourseView {
  return {
    course,
    enrollment: enrolled ? { courseId: course.id, id: 1, state: "started", userId: 7 } : null,
    landingPage,
    lessons: lessonRows(standing.finished),
    modules: [],
    movedTo: null,
    progress: standing.progress,
    qnaItems: [],
  };
}

const learner: AuthUser = {
  admin: false,
  assistantMessagesCount: 0,
  canAccessAdmin: false,
  createdAt: "2026-01-01T00:00:00Z",
  createdAtAsTimestamp: null,
  email: "ada@example.com",
  firstName: "Ada",
  id: 7,
  lastName: null,
  name: "Ada",
  type: "user",
};

// The real route at a real URL. Start commands are recorded from before the
// first render, together with where the visitor was when each arrived — so a
// test can tell "nothing started on load" and "started, then navigated" apart
// from their opposites.
async function renderCourse(
  user: AuthUser | null = null,
  initialPath = "/languages/javascript",
  { withPlayer = false } = {},
): Promise<RenderedRoute & { starts: { id: string; pathname: string }[] }> {
  const starts: { id: string; pathname: string }[] = [];
  let pathname = (): string => initialPath;
  worker.use(
    http.post("*/lessons/:id/start", ({ params }) => {
      starts.push({ id: String(params.id), pathname: pathname() });
      return new HttpResponse(null, { status: 204 });
    }),
  );

  const rendered = await renderRoute(courseRoute, {
    initialPath,
    path: "/{-$locale}/languages/$slug",
    user,
    // The lesson player beside it, for a test that follows the button into the
    // lesson; it divides a real height, so it gets one.
    ...(withPlayer && {
      also: [{ path: "/{-$locale}/languages/$slug/lessons/$lessonSlug", route: lessonRoute }],
      wrap: (element) => <div style={{ height: "800px", width: "1200px" }}>{element}</div>,
    }),
  });
  pathname = () => rendered.router.state.location.pathname;

  return { ...rendered, starts };
}

// The lesson entries in the list, in the order a visitor reads them. The marks
// are labelled icons, so they are asserted by accessible name, not here.
function lessonNames(): string[] {
  return page
    .getByRole("link")
    .elements()
    .filter((el) => el.getAttribute("href")?.includes("/lessons/"))
    .map((el) => el.textContent?.trim());
}

// The Next Lesson as the lesson player reads it, with the standing it carries:
// the marks are what the pass the start followed has changed.
function variablesView(standing: Standing): CourseLessonView {
  return {
    adjacent: {
      next: lessonRef(STRINGS_ID, "strings", "Strings"),
      previous: lessonRef(HELLO_WORLD_ID, "hello-world", "Hello, World!"),
    },
    landingPage,
    lesson: {
      course,
      createdAt: "2026-01-01T00:00:00Z",
      definitions: [],
      description: null,
      exerciseExtension: null,
      id: 1002,
      instructions: "Assign a variable.",
      locale: "en",
      name: "Variables",
      naturalOrder: 2,
      originalCode: "",
      preparedCode: "",
      slug: "variables",
      sourceCodeUrl: null,
      testCode: "",
      testExtension: null,
      theory: "A variable is a name bound to a value.",
      tips: [],
      version: 99,
      versionId: 99,
    },
    lessons: lessonRows(standing.finished),
    progress: standing.progress,
  };
}

describe("course page", () => {
  it("a learner mid-course sees the landing copy, the marked lessons, and Continue", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/languages/javascript", () => HttpResponse.json(courseView(midCourse, true))),
    );

    const { router, starts } = await renderCourse(learner);

    await expect.element(page.getByRole("heading", { name: "JavaScript course" })).toBeVisible();
    await expect.element(page.getByText("Learn JavaScript in the browser")).toBeVisible();
    await expect
      .element(page.getByRole("heading", { name: "Where JavaScript is used" }))
      .toBeVisible();

    // Every current lesson, in course order.
    expect(lessonNames()).toStrictEqual(["Hello, World!", "Variables", "Strings"]);
    await expect.element(page.getByRole("link", { name: "Finished Hello, World!" })).toBeVisible();
    await expect.element(page.getByRole("link", { name: "Strings" })).toBeVisible();
    expect(page.getByLabelText("Finished").elements()).toHaveLength(1);
    // Lessons may be taken in any order, so every entry is a plain link.
    await expect
      .element(page.getByRole("link", { name: "Strings" }))
      .toHaveAttribute("href", "/languages/javascript/lessons/strings");

    // Completion as the server reported it.
    await expect.element(page.getByText("Completed 33%")).toBeVisible();

    // Nothing started by merely loading the page.
    expect(starts).toStrictEqual([]);

    // Continue enters the Next Lesson — started first, then navigated to.
    await page.getByRole("button", { name: "Continue Learning" }).click();
    await expect
      .poll(() => router.state.location.pathname)
      .toBe("/languages/javascript/lessons/variables");
    expect(starts).toStrictEqual([{ id: "1002", pathname: "/languages/javascript" }]);
  }, 60_000);

  it("a finished learner is told so and offered no lesson", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/languages/javascript", () => HttpResponse.json(courseView(finished, true))),
    );

    await renderCourse(learner);

    await expect.element(page.getByText("Course finished!")).toBeVisible();
    await expect.element(page.getByText("Completed 100%")).toBeVisible();
    expect(page.getByLabelText("Finished").elements()).toHaveLength(FINISHED_LESSONS);
    expect(page.getByRole("button", { name: "Continue Learning" }).elements()).toHaveLength(0);
    expect(page.getByRole("button", { name: "Try It" }).elements()).toHaveLength(0);
  }, 60_000);

  it("lays the program out by module, the first one open, and answers the landing page's questions", async () => {
    expect.hasAssertions();
    const view: CourseView = {
      ...courseView(midCourse),
      modules: [
        {
          description: "The first program and variables.",
          id: 1,
          lessons: lessonRows(midCourse.finished).slice(0, 2),
          name: "Basics",
        },
        {
          description: "Working with strings.",
          id: 3,
          lessons: lessonRows(midCourse.finished).slice(2),
          name: "Text",
        },
      ],
      qnaItems: [{ answer: "Yes, **entirely**.", id: 8001, question: "Is it free?" }],
    };
    worker.use(http.get("*/languages/javascript", () => HttpResponse.json(view)));

    await renderCourse(learner);

    // The first module is open: its lessons, marked, beside its description.
    await expect.element(page.getByRole("heading", { name: "Basics" })).toBeVisible();
    await expect.element(page.getByText("The first program and variables.")).toBeVisible();
    await expect.element(page.getByRole("link", { name: "Finished Hello, World!" })).toBeVisible();
    await expect.element(page.getByRole("link", { name: "Variables" })).toBeVisible();
    // Only the open module's lessons are shown; the next is folded until asked
    // for, and opening it folds the first, as legacy's accordion did.
    expect(lessonNames()).toStrictEqual(["Hello, World!", "Variables"]);
    await page.getByRole("button", { name: "Text" }).click();
    await expect.element(page.getByRole("link", { name: "Strings" })).toBeVisible();
    await expect.element(page.getByText("Working with strings.")).toBeVisible();
    await expect.poll(lessonNames).toEqual(["Strings"]);

    await expect
      .element(page.getByRole("heading", { name: "Sorting out the questions" }))
      .toBeVisible();
    await expect.element(page.getByText("Is it free?")).toBeVisible();
    await expect.element(page.getByText("entirely")).toBeVisible();
  }, 60_000);

  it("a finished learner is pointed on to the Hexlet program, tagged as a referral", async () => {
    expect.hasAssertions();
    const view = courseView(finished);
    view.course = {
      ...course,
      hexletProgramLandingPage: "https://ru.hexlet.io/programs/js?ref=1",
    };
    worker.use(http.get("*/languages/javascript", () => HttpResponse.json(view)));

    await renderCourse(learner);

    await expect
      .element(page.getByRole("link", { name: "Profession and employment" }))
      .toBeVisible();
    const href = page
      .getByRole("link", { name: "Profession and employment" })
      .element()
      .getAttribute("href");
    const url = new URL(String(href));
    expect(`${url.origin}${url.pathname}`).toBe("https://ru.hexlet.io/programs/js");
    expect(Object.fromEntries(url.searchParams)).toStrictEqual({
      ref: "1",
      utm_medium: "referral",
      utm_source: "code-basics",
    });
  }, 60_000);

  it("a guest carrying progress sees the same page and the same action", async () => {
    expect.hasAssertions();
    // The same payload a signed-in learner gets: the guest's position is the
    // cookie, and the page has no reason to know the difference.
    worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView(midCourse))));

    const { router, starts } = await renderCourse(null);

    // A landing page shows no completion to someone who has not started.
    await expect.element(page.getByRole("button", { name: "Continue Learning" })).toBeVisible();
    expect(page.getByText("Completed 33%").elements()).toHaveLength(0);
    expect(lessonNames()).toStrictEqual(["Hello, World!", "Variables", "Strings"]);
    await expect.element(page.getByRole("link", { name: "Finished Hello, World!" })).toBeVisible();
    await expect.element(page.getByRole("link", { name: "Strings" })).toBeVisible();
    expect(page.getByLabelText("Finished").elements()).toHaveLength(1);

    await page.getByRole("button", { name: "Continue Learning" }).click();
    await expect
      .poll(() => router.state.location.pathname)
      .toBe("/languages/javascript/lessons/variables");
    expect(starts).toStrictEqual([{ id: "1002", pathname: "/languages/javascript" }]);
  }, 60_000);

  it("a first-time visitor sees every lesson as a plain link and Try It", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView(firstVisit))));

    const { router, starts } = await renderCourse(null);

    await expect.element(page.getByRole("button", { name: "Try It" })).toBeVisible();
    expect(page.getByText("Completed 0%").elements()).toHaveLength(0);
    expect(page.getByLabelText("Finished").elements()).toHaveLength(0);
    expect(lessonNames()).toStrictEqual(["Hello, World!", "Variables", "Strings"]);
    await expect.element(page.getByRole("link", { name: "Variables" })).toBeVisible();
    await expect.element(page.getByRole("link", { name: "Strings" })).toBeVisible();

    await page.getByRole("button", { name: "Try It" }).click();
    await expect
      .poll(() => router.state.location.pathname)
      .toBe("/languages/javascript/lessons/hello-world");
    expect(starts).toStrictEqual([{ id: "1001", pathname: "/languages/javascript" }]);
  }, 60_000);

  it("a landing page slug that is not served here is not found", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/languages/javascript-en", () =>
        HttpResponse.json({ status: 404, title: "Not Found" }, { status: 404 }),
      ),
    );

    await renderCourse(null, "/languages/javascript-en");

    await expect.element(page.getByRole("heading", { name: "Page Not Found" })).toBeVisible();
  }, 60_000);

  it("an archived landing page sends the visitor on to the page it moved to", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/languages/javascript-old", () =>
        HttpResponse.json({ ...courseView(firstVisit), movedTo: "javascript-ru" }),
      ),
      http.get("*/languages/javascript-ru", () => HttpResponse.json(courseView(firstVisit))),
    );

    const { router } = await renderCourse(null, "/languages/javascript-old");

    await expect.poll(() => router.state.location.pathname).toBe("/languages/javascript-ru");
    await expect
      .element(page.getByRole("heading", { level: 1, name: "JavaScript course" }))
      .toBeVisible();
  }, 60_000);

  it("keeps the URL locale when entering the course", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView(firstVisit))));

    const { router } = await renderCourse(null, "/ru/languages/javascript");

    await page.getByRole("button", { name: "Начать обучение" }).click();
    await expect
      .poll(() => router.state.location.pathname)
      .toBe("/ru/languages/javascript/lessons/hello-world");
  }, 60_000);

  it("stays on the course page when the start command fails", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView(firstVisit))));

    const { router } = await renderCourse(null);
    // Registered after the recorder, so it is the handler that answers.
    worker.use(http.post("*/lessons/:id/start", () => HttpResponse.error()));

    await page.getByRole("button", { name: "Try It" }).click();
    await expect
      .element(page.getByText("There was a network problem", { exact: false }))
      .toBeVisible();
    expect(router.state.location.pathname).toBe("/languages/javascript");
  }, 60_000);

  it("enters the lesson as it stands after the start, not as a hover cached it", async () => {
    expect.hasAssertions();
    // The player lays out by the viewport, so it gets the desktop one it expects.
    const startingViewport = { height: window.innerHeight, width: window.innerWidth };
    await page.viewport(DESKTOP_VIEWPORT.width, DESKTOP_VIEWPORT.height);
    onTestFinished(() => page.viewport(startingViewport.width, startingViewport.height));
    worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView(midCourse))));
    // The fresh read is held until the test lets it go, so what the learner is
    // shown while it is in flight can be looked at.
    const { promise: held, resolve: release } = Promise.withResolvers<null>();
    let reads = 0;
    worker.use(
      http.get("*/languages/javascript/lessons/variables", async () => {
        reads += 1;
        await held;
        return HttpResponse.json(variablesView(midCourse));
      }),
    );

    const { starts, queryClient } = await renderCourse(learner, "/languages/javascript", {
      withPlayer: true,
    });
    // What a hover over the list would have preloaded before the pass: the same
    // lesson, with nothing finished yet.
    queryClient.setQueryData(
      getCourseLessonQueryKey({ path: { courseSlug: "javascript", slug: "variables" } }),
      variablesView(firstVisit),
    );

    await page.getByRole("button", { name: "Continue Learning" }).click();
    await expect.poll(() => reads).toBe(1);
    expect(starts).toStrictEqual([{ id: "1002", pathname: "/languages/javascript" }]);

    // While the lesson is read afresh the learner is still on the course page:
    // the cached copy, marks and all, is never put in front of them.
    expect(page.getByRole("heading", { name: `${COURSE}: Variables` }).elements()).toHaveLength(0);
    await expect.element(page.getByRole("button", { name: "Continue Learning" })).toBeVisible();

    release(null);
    await expect.element(page.getByRole("heading", { name: `${COURSE}: Variables` })).toBeVisible();
    await page.getByRole("tab", { name: "Navigation" }).click();
    await expect.element(page.getByRole("link", { name: "Finished Hello, World!" })).toBeVisible();
  }, 60_000);
});
