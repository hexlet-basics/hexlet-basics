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

// The Course page, driven through its real route with the API faked at the HTTP
// boundary. What is asserted is what a visitor sees — the landing copy, the
// marks in the lesson list, the one button — and the requests pressing it made.

const course: Course = {
  id: 1,
  slug: "javascript",
  name: "javascript",
  learnAs: null,
  readiness: null,
  categoryId: null,
  currentVersionId: 99,
  currentVersion: null,
  createdAt: "2026-01-01T00:00:00Z",
  enrollmentsCount: 1234,
  lessonsCount: 3,
  ratingCount: 0,
  ratingValue: 0,
  repositoryUrl: null,
  hexletProgramLandingPage: null,
  coverListVariant: null,
  coverThumbVariant: null,
};

const landingPage: CourseLandingPage = {
  id: 10,
  courseId: 1,
  courseSlug: "javascript",
  createdAt: "2026-01-01T00:00:00Z",
  slug: "javascript-ru",
  name: "JavaScript",
  main: true,
  listed: true,
  state: "published",
  order: null,
  footer: true,
  footerName: null,
  landingPageToRedirectId: null,
  metaTitle: "JavaScript course",
  metaDescription: "Learn JavaScript",
  header: "JavaScript course",
  description: "Learn JavaScript in the browser",
  usedInHeader: "Where JavaScript is used",
  usedInDescription: "Browsers, servers and everything in between",
  outcomesHeader: null,
  outcomesDescription: null,
  outcomesImage: null,
  duration: 10,
  enrollmentsCount: 0,
};

// One standing, as the server sends it: the course's own progress, and which of
// its lessons carry a check. The page is expected to render them as given.
interface Standing {
  progress: CourseProgress;
  finished: string[];
}

const lessonRef = (id: number, slug: string, name: string): LessonRef => ({ id, slug, name });

// The course's lessons in course order, each marked or not.
const lessonRows = (finished: string[]): LearnerLesson[] => [
  {
    id: 1001,
    name: "Hello, World!",
    description: null,
    slug: "hello-world",
    position: 1,
    finished: finished.includes("hello-world"),
  },
  {
    id: 1002,
    name: "Variables",
    description: null,
    slug: "variables",
    position: 2,
    finished: finished.includes("variables"),
  },
  {
    id: 1003,
    name: "Strings",
    description: null,
    slug: "strings",
    position: 3,
    finished: finished.includes("strings"),
  },
];

// Three payloads the server can send, each a fact about where the visitor
// stands.

// Has finished the first lesson: the second is next.
const midCourse: Standing = {
  progress: {
    state: "started",
    completion: 33,
    nextLesson: lessonRef(1002, "variables", "Variables"),
  },
  finished: ["hello-world"],
};

// Has finished everything: there is no Next Lesson.
const finished: Standing = {
  progress: { state: "finished", completion: 100, nextLesson: null },
  finished: ["hello-world", "variables", "strings"],
};

// Has never been here: nothing finished, the first lesson is next.
const firstVisit: Standing = {
  progress: {
    state: null,
    completion: 0,
    nextLesson: lessonRef(1001, "hello-world", "Hello, World!"),
  },
  finished: [],
};

// `enrolled` is a learner who has started the course: only they have a completion
// to show on its landing page.
function courseView(standing: Standing, enrolled = false): CourseView {
  return {
    course,
    landingPage,
    movedTo: null,
    lessons: lessonRows(standing.finished),
    modules: [],
    qnaItems: [],
    enrollment: enrolled ? { id: 1, userId: 7, courseId: course.id, state: "started" } : null,
    progress: standing.progress,
  };
}

const learner: AuthUser = {
  id: 7,
  firstName: "Ada",
  lastName: null,
  name: "Ada",
  email: "ada@example.com",
  admin: false,
  canAccessAdmin: false,
  assistantMessagesCount: 0,
  createdAt: "2026-01-01T00:00:00Z",
  createdAtAsTimestamp: null,
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
) {
  const starts: { id: string; pathname: string }[] = [];
  let pathname = () => initialPath;
  worker.use(
    http.post("*/lessons/:id/start", ({ params }) => {
      starts.push({ id: String(params.id), pathname: pathname() });
      return new HttpResponse(null, { status: 204 });
    }),
  );

  const rendered = await renderRoute(courseRoute, {
    path: "/{-$locale}/languages/$slug",
    initialPath,
    user,
    // The lesson player beside it, for a test that follows the button into the
    // lesson; it divides a real height, so it gets one.
    ...(withPlayer && {
      also: [{ route: lessonRoute, path: "/{-$locale}/languages/$slug/lessons/$lessonSlug" }],
      wrap: (element) => <div style={{ height: "800px", width: "1200px" }}>{element}</div>,
    }),
  });
  pathname = () => rendered.router.state.location.pathname;

  return { ...rendered, starts };
}

// The lesson entries in the list, in the order a visitor reads them. The marks
// are labelled icons, so they are asserted by accessible name, not here.
function lessonNames() {
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
    lesson: {
      course,
      id: 1002,
      name: "Variables",
      slug: "variables",
      locale: "en",
      naturalOrder: 2,
      versionId: 99,
      version: 99,
      description: null,
      instructions: "Assign a variable.",
      theory: "A variable is a name bound to a value.",
      definitions: [],
      tips: [],
      preparedCode: "",
      originalCode: "",
      testCode: "",
      exerciseExtension: null,
      testExtension: null,
      sourceCodeUrl: null,
      createdAt: "2026-01-01T00:00:00Z",
    },
    landingPage,
    lessons: lessonRows(standing.finished),
    adjacent: {
      previous: lessonRef(1001, "hello-world", "Hello, World!"),
      next: lessonRef(1003, "strings", "Strings"),
    },
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
    expect(page.getByLabelText("Finished").elements()).toHaveLength(3);
    expect(page.getByRole("button", { name: "Continue Learning" }).elements()).toHaveLength(0);
    expect(page.getByRole("button", { name: "Try It" }).elements()).toHaveLength(0);
  }, 60_000);

  it("lays the program out by module, the first one open, and answers the landing page's questions", async () => {
    expect.hasAssertions();
    const view: CourseView = {
      ...courseView(midCourse),
      modules: [
        {
          id: 1,
          name: "Basics",
          description: "The first program and variables.",
          lessons: lessonRows(midCourse.finished).slice(0, 2),
        },
        {
          id: 3,
          name: "Text",
          description: "Working with strings.",
          lessons: lessonRows(midCourse.finished).slice(2),
        },
      ],
      qnaItems: [{ id: 8001, question: "Is it free?", answer: "Yes, **entirely**." }],
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
      utm_source: "code-basics",
      utm_medium: "referral",
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
    const startingViewport = { width: window.innerWidth, height: window.innerHeight };
    await page.viewport(1200, 800);
    onTestFinished(() => page.viewport(startingViewport.width, startingViewport.height));
    worker.use(http.get("*/languages/javascript", () => HttpResponse.json(courseView(midCourse))));
    // The fresh read is held until the test lets it go, so what the learner is
    // shown while it is in flight can be looked at.
    let release!: () => void;
    const held = new Promise<void>((resolve) => {
      release = resolve;
    });
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
    expect(page.getByRole("heading", { name: "JavaScript: Variables" }).elements()).toHaveLength(0);
    await expect.element(page.getByRole("button", { name: "Continue Learning" })).toBeVisible();

    release();
    await expect
      .element(page.getByRole("heading", { name: "JavaScript: Variables" }))
      .toBeVisible();
    await page.getByRole("tab", { name: "Navigation" }).click();
    await expect.element(page.getByRole("link", { name: "Finished Hello, World!" })).toBeVisible();
  }, 60_000);
});
