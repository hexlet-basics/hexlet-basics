import { Cookies } from "react-cookie";
import { HttpResponse, http } from "msw";
import { describe, expect, it, onTestFinished } from "vitest";
import { page, userEvent } from "vitest/browser";
import type {
  AdjacentLessons,
  Course,
  CourseLandingPage,
  CourseLesson,
  CourseLessonView,
  CourseView,
  EnrollmentState,
  LearnerLesson,
  LessonAssistantChat,
  LessonAssistantMessage,
  LessonCheckingResponse,
  LessonRef,
} from "@/client/types.gen";
import type { AuthUser } from "@/lib/auth";
import { getCourseQueryKey } from "@/client/@tanstack/react-query.gen";
import { Route as lessonRoute } from "@/routes/{-$locale}/languages/$slug/lessons/$lessonSlug";
import { Route as successRoute } from "@/routes/{-$locale}/languages/$slug/success";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";
import { sleep } from "@/test/wait";
import type { Locator } from "vitest/browser";
import type { DefaultBodyType } from "msw";
import type { RenderedRoute } from "@/test/renderRoute";

// A lesson page is headed with the course it belongs to.
const cookies = new Cookies();

const COURSE = "JavaScript";

const FIRST_LESSON_ID = 1001;
// How long a check takes to answer, in the tests that look at the player while it runs.
const SLOW_CHECK_MS = 600;
const QUICK_CHECK_MS = 300;
const LONG_LINE_LENGTH = 200;

// The lesson player, driven through its real route with the API faked at the
// HTTP boundary — the seam the admin screens are tested at. What is asserted is
// what a learner can see: the words on the page, the marks in the lesson list,
// and the requests the page did (and did not) make.

const course: Course = {
  categoryId: null,
  coverListVariant: null,
  coverThumbVariant: null,
  createdAt: "2026-01-01T00:00:00Z",
  currentVersion: null,
  currentVersionId: 99,
  enrollmentsCount: 0,
  hexletProgramLandingPage: null,
  id: 1,
  learnAs: null,
  lessonsCount: 3,
  // The course's own name, which is NOT what the page is titled with.
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
  // The marketing name a learner recognises — the one the page is titled with.
  name: "JavaScript",
  order: null,
  outcomesDescription: null,
  outcomesHeader: null,
  outcomesImage: null,
  slug: "javascript-ru",
  state: "published",
  usedInDescription: null,
  usedInHeader: null,
};

const lesson: CourseLesson = {
  course,
  createdAt: "2026-01-01T00:00:00Z",
  definitions: [],
  description: "Storing values",
  exerciseExtension: "js",
  id: 1002,
  instructions: "Assign the string `hello` to `greeting`.",
  locale: "en",
  name: "Variables",
  naturalOrder: 2,
  originalCode: "let greeting = 'hello'; // the author's answer\n",
  preparedCode: "let greeting = '';\n",
  slug: "variables",
  sourceCodeUrl: "https://github.com/hexlet-basics/exercises-javascript",
  testCode: "assert.equal(greeting, 'hello');\n",
  testExtension: "js",
  theory: "A variable is a name bound to a value.\n\n```js\nlet greeting = 'hello';\n```",
  tips: ["Names are case sensitive"],
  version: 99,
  versionId: 99,
};

// The course's lessons in course order, each marked or not.
function lessonRows(finished: string[] = []): LearnerLesson[] {
  return [
    { name: "Hello, World!", slug: "hello-world" },
    { name: "Variables", slug: "variables" },
    { name: "Strings", slug: "strings" },
  ].map(({ slug, name }, index) => ({
    description: null,
    finished: finished.includes(slug),
    id: FIRST_LESSON_ID + index,
    name,
    position: index + 1,
    slug,
  }));
}

const lessonRef = (slug: string): LessonRef => {
  const row = lessonRows().find((item) => item.slug === slug);
  if (!row) {
    throw new Error(`no lesson row ${slug}`);
  }
  return { id: row.id, name: row.name, slug };
};

// The lessons on either side of one, in course order: the server's answer.
function adjacentTo(slug: string): AdjacentLessons {
  const slugs = lessonRows().map((item) => item.slug);
  const index = slugs.indexOf(slug);
  const previous = slugs[index - 1];
  const next = slugs[index + 1];
  return {
    next: next ? lessonRef(next) : null,
    previous: previous ? lessonRef(previous) : null,
  };
}

// A learner who has finished the first lesson: the second is next.
function lessonView(overrides: Partial<CourseLessonView> = {}): CourseLessonView {
  return {
    adjacent: adjacentTo((overrides.lesson ?? lesson).slug),
    landingPage,
    lesson,
    lessons: lessonRows(["hello-world"]),
    progress: { completion: 33, nextLesson: lessonRef("variables"), state: "started" },
    ...overrides,
  };
}

// A signed-in learner, for the tests that need one. Everything else here is a
// guest, which is the harder case and the one legacy under-tested.
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

// The real route at a real URL, so the loader and the route's chrome run as they
// do for a visitor. The lesson layout gives the panes the viewport to divide.
function renderPlayer(lessonSlug: string, user: AuthUser | null = null): Promise<RenderedRoute> {
  return renderRoute(lessonRoute, {
    initialPath: `/languages/javascript/lessons/${lessonSlug}`,
    path: "/{-$locale}/languages/$slug/lessons/$lessonSlug",
    user,
  });
}

// The layout switches on the viewport, not on the container, so each test sets
// the viewport it means: a desktop unless it says otherwise, and the viewport
// the suite was started with once the file is done.
const desktop = { height: 800, width: 1200 };
const phone = { height: 667, width: 375 };
const startingViewport = { height: window.innerHeight, width: window.innerWidth };
// What monaco has painted, with the non-breaking spaces it renders text with
// turned back into ordinary ones so assertions read like the code does.
function editorText(): string {
  return (document.querySelector(".view-lines")?.textContent ?? "").replaceAll("\u00A0", " ");
}

// A completed run, as the server reports one.
function checkResult(overrides: Partial<LessonCheckingResponse> = {}): LessonCheckingResponse {
  return {
    courseHasBeenFinished: false,
    lessonHasBeenFinished: true,
    output: "1 example, 0 failures",
    passed: true,
    result: "passed",
    status: 0,
    ...overrides,
  };
}

// A markup course: the exercise is a page, and the learner watches it take shape.
function htmlLessonView(): CourseLessonView {
  return lessonView({
    landingPage: { ...landingPage, courseSlug: "html", name: "HTML", slug: "html-ru" },
    lesson: {
      ...lesson,
      course: { ...course, name: "html", slug: "html" },
      name: "Headings",
      preparedCode: "<h1>Hello</h1>\n",
      slug: "headings",
    },
  });
}

// The last lesson of the course, which the learner has reached but not passed.
function lastLessonView(): CourseLessonView {
  return lessonView({
    lesson: { ...lesson, id: 1003, name: "Strings", slug: "strings" },
    lessons: lessonRows(["hello-world", "variables"]),
    progress: { completion: 66, nextLesson: lessonRef("strings"), state: "started" },
  });
}

// The player at a phone's viewport, in a container of the same size.
async function renderOnPhone(
  lessonSlug: string,
  user: AuthUser | null = null,
): Promise<RenderedRoute> {
  await page.viewport(phone.width, phone.height);
  return renderPlayer(lessonSlug, user);
}

// The two panes — theory and navigation, then the workspace — as locators, and
// how wide each is drawn: on a phone one fills the screen and the other is
// folded to nothing.
function panes(): { theory: Locator; workspace: Locator } {
  const [theory, workspace] = Array.from(
    document.querySelectorAll<HTMLElement>(".mantine-Splitter-pane"),
    (el) => page.elementLocator(el),
  );
  if (!theory || !workspace) {
    throw new Error("the player has no panes");
  }
  return { theory, workspace };
}

function paneWidths(): number[] {
  return Array.from(document.querySelectorAll<HTMLElement>(".mantine-Splitter-pane"), (el) =>
    Math.round(el.getBoundingClientRect().width),
  );
}

function assistantMessage(overrides: Partial<LessonAssistantMessage>): LessonAssistantMessage {
  return {
    content: "",
    courseLessonName: "Variables",
    courseLessonSlug: "variables",
    courseSlug: "javascript",
    createdAt: "2026-01-01T00:00:00Z",
    id: 1,
    role: "assistant",
    userId: null,
    ...overrides,
  };
}

// The answer as the server sends it: plain text, a chunk at a time.
function streamedAnswer(chunks: string[]): HttpResponse<DefaultBodyType> {
  const encoder = new TextEncoder();
  const body = new ReadableStream({
    start(controller) {
      for (const chunk of chunks) {
        controller.enqueue(encoder.encode(chunk));
      }
      controller.close();
    },
  });
  return new HttpResponse(body, { headers: { "Content-Type": "text/plain; charset=utf-8" } });
}

// What the learner sees underlined: monaco paints a marker as a squiggle.
function squiggles(): number {
  return document.querySelectorAll(".squiggly-error").length;
}

// The course read for a learner whose Enrollment is in `state` — what the
// completion page checks before it congratulates anyone.
function courseView(state: EnrollmentState): CourseView {
  const progress = lastLessonView().progress ?? null;
  return {
    course,
    enrollment: progress && { courseId: course.id, id: 5, state, userId: learner.id },
    landingPage,
    lessons: [],
    modules: [],
    movedTo: null,
    progress,
    qnaItems: [],
  };
}

function assistantChat(overrides: Partial<LessonAssistantChat> = {}): LessonAssistantChat {
  return {
    messages: [
      assistantMessage({ content: "What is a variable?", id: 1, role: "user", userId: 7 }),
      assistantMessage({ content: "A name bound to a value.", id: 2, role: "assistant" }),
    ],
    quotaExceeded: false,
    ...overrides,
  };
}

// What each control under the editor says, by its label or else its text.
function controlLabels(bar: HTMLElement | null): string[] {
  return [...(bar?.children ?? [])].map(
    (el) => el.getAttribute("aria-label") ?? el.textContent?.trim(),
  );
}

// Lays the page out at a desktop viewport for one test. Everything the player
// keeps between visits lives in localStorage, and a test that seeds it must not
// leak into the next one, so both are put back when the test is done.
async function atDesktop(): Promise<void> {
  onTestFinished(async () => {
    localStorage.clear();
    await page.viewport(startingViewport.width, startingViewport.height);
  });
  await page.viewport(desktop.width, desktop.height);
}

describe("lesson player", () => {
  // The editor is a lazy chunk of a few hundred modules, and the first test to
  // open it waits for vite to transform all of them, so waiting on it gets a
  // budget of its own rather than the default second. The budget is sized for a
  // cold CI runner, not for a warm laptop.
  const editorLoad = { timeout: 30_000 };

  it("renders the lesson a learner reads, titled from the course's landing copy", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );

    await renderPlayer("variables");

    // The landing page's marketing name, not the course's own `name`.
    await expect.element(page.getByRole("heading", { name: `${COURSE}: Variables` })).toBeVisible();
    await expect.element(page.getByText("A variable is a name bound to a value.")).toBeVisible();
    await expect.element(page.getByRole("heading", { name: "Instructions" })).toBeVisible();
    await expect.element(page.getByText("Assign the string hello to greeting.")).toBeVisible();
    await expect.element(page.getByRole("heading", { name: "Tips" })).toBeVisible();
    await expect.element(page.getByText("Names are case sensitive")).toBeVisible();

    // The questions everyone asks, and the way to fix what you just read.
    await expect
      .element(page.getByText("The exercise doesn't pass checking. What to do? 😶"))
      .toBeVisible();
    await expect
      .element(page.getByRole("link", { name: "Lesson source on GitHub" }))
      .toHaveAttribute("href", "https://github.com/hexlet-basics/exercises-javascript");

    // The workspace the learner works in, tab by tab.
    await expect.element(page.getByRole("tab", { name: "Editor" })).toBeVisible();
    await expect.element(page.getByRole("tab", { name: "Output" })).toBeVisible();
    await expect.element(page.getByRole("tab", { name: "Tests" })).toBeVisible();
    await expect.element(page.getByRole("tab", { name: "Solution" })).toBeVisible();
  }, 60_000);

  it("lists every lesson in course order, marking the finished ones", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );

    await renderPlayer("variables");
    await page.getByRole("tab", { name: "Navigation" }).click();

    // The list itself — the control bar's Previous links to a lesson too.
    const links = page.getByRole("tabpanel", { name: "Navigation" }).getByRole("link").elements();
    const lessonLinks = links.filter((el) => el.getAttribute("href")?.includes("/lessons/"));
    expect(lessonLinks.map((el) => el.textContent?.trim())).toStrictEqual([
      "Hello, World!",
      "Variables",
      "Strings",
    ]);

    await expect.element(page.getByLabelText("Finished")).toBeVisible();

    // Where the learner is now.
    await expect
      .element(page.getByRole("link", { name: "Variables" }))
      .toHaveAttribute("aria-current", "page");

    // Any lesson may be taken in any order, so every entry is a plain link.
    await expect
      .element(page.getByRole("link", { name: "Strings" }))
      .toHaveAttribute("href", "/languages/javascript/lessons/strings");
  }, 60_000);

  it("serves a visitor any lesson in full", async () => {
    expect.hasAssertions();
    await atDesktop();
    const visitor = lessonView({
      lessons: lessonRows(),
      progress: { completion: 0, nextLesson: lessonRef("hello-world"), state: null },
    });
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(visitor)),
    );

    await renderPlayer("variables");

    // Readable in full — this is the page a search engine indexes.
    await expect.element(page.getByRole("heading", { name: `${COURSE}: Variables` })).toBeVisible();
    await expect.element(page.getByText("A variable is a name bound to a value.")).toBeVisible();
  }, 60_000);

  it("starts nothing by loading the page", async () => {
    expect.hasAssertions();
    await atDesktop();
    let started = 0;
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.post("*/lessons/:id/start", () => {
        started += 1;
        return HttpResponse.json({});
      }),
    );

    await renderPlayer("variables");
    await expect.element(page.getByRole("heading", { name: `${COURSE}: Variables` })).toBeVisible();
    await page.getByRole("tab", { name: "Navigation" }).click();
    await expect.element(page.getByRole("link", { name: "Strings" })).toBeVisible();

    // The router preloads on hover, so a read that started a lesson would enroll a
    // learner in every lesson they pointed at (ADR-0012).
    expect(started).toBe(0);
  }, 60_000);

  it("divides the two panes where the learner left the divider", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );
    // What a previous visit stored — a reload reads it back rather than resetting
    // to the default split.
    cookies.set("lesson_panes", "70", { path: "/" });

    await renderPlayer("variables");
    await expect.element(page.getByRole("heading", { name: `${COURSE}: Variables` })).toBeVisible();

    // The handle is what makes the divider draggable at all.
    expect(document.querySelectorAll(".mantine-Splitter-handle")).toHaveLength(1);

    const paneElements = document.querySelectorAll<HTMLElement>(".mantine-Splitter-pane");
    expect(paneElements).toHaveLength(2);
    expect(paneElements[0]?.style.cssText).toContain("70%");

    cookies.remove("lesson_panes", { path: "/" });
  }, 60_000);

  it("says so when the lesson is not there", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/nope", () =>
        HttpResponse.json({ status: 404, title: "Not Found" }, { status: 404 }),
      ),
    );

    await renderPlayer("nope");

    await expect
      .element(page.getByText("Lesson not found. Please, try another lesson."))
      .toBeVisible();
  }, 60_000);

  it("opens the editor with the lesson's starter code, after the theory", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );

    await renderPlayer("variables");

    // The theory does not wait for the editor: it renders from the payload, while
    // monaco arrives on its own lazy import afterwards.
    await expect.element(page.getByRole("heading", { name: `${COURSE}: Variables` })).toBeVisible();

    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();
    await expect.poll(editorText, editorLoad).toContain("let greeting = '';");

    // Monaco is this application's own copy, not a CDN's: every resource, the
    // API's (answered by MSW) included, comes from the page's own origin.
    const remote = performance
      .getEntriesByType("resource")
      .filter((entry) => !entry.name.startsWith(location.origin));
    expect(remote).toStrictEqual([]);
  }, 60_000);

  it("finds yesterday's work still there, and keeps it per lesson", async () => {
    expect.hasAssertions();
    await atDesktop();
    const strings = lessonView({
      lesson: {
        ...lesson,
        id: 1003,
        name: "Strings",
        preparedCode: "// strings\n",
        slug: "strings",
      },
    });
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.get("*/languages/javascript/lessons/strings", () => HttpResponse.json(strings)),
    );
    localStorage.setItem(
      "lesson-code-javascript-variables",
      JSON.stringify("let greeting = 'yesterday';"),
    );

    const { screen } = await renderPlayer("variables");
    await expect.poll(editorText, editorLoad).toContain("let greeting = 'yesterday';");

    // The next lesson opens on its own starter code, not on the work done in the
    // one before it.
    await screen.unmount();
    await renderPlayer("strings");
    await expect.poll(editorText, editorLoad).toContain("// strings");
    expect(editorText()).not.toContain("yesterday");
  }, 60_000);

  it("keeps what the learner typed when they come back to the lesson", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );

    const { screen } = await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();

    // The code area of the editor pane — the theory has code blocks of its own —
    // and not monaco's hidden input, because that is what a learner clicks.
    await page.getByRole("tabpanel", { name: "Editor" }).getByRole("code").click();
    await userEvent.keyboard("// mine");
    await expect.poll(editorText).toContain("// mine");

    await screen.unmount();
    await renderPlayer("variables");

    await expect.poll(editorText, editorLoad).toContain("// mine");
  }, 60_000);

  it("asks before resetting, and restores the starter code", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );
    localStorage.setItem(
      "lesson-code-javascript-variables",
      JSON.stringify("let greeting = 'mess';"),
    );

    await renderPlayer("variables");
    await expect.poll(editorText, editorLoad).toContain("let greeting = 'mess';");

    await page.getByRole("button", { name: "Reset" }).click();

    // One misclick must not cost a learner their work.
    await expect
      .element(
        page.getByText(
          "You want to reset the exercise progress. The current code version will not be saved. We hope you’ve already copied it. Continue resetting?",
        ),
      )
      .toBeVisible();
    await page.getByRole("button", { exact: true, name: "No" }).click();
    expect(editorText()).toContain("let greeting = 'mess';");

    await page.getByRole("button", { name: "Reset" }).click();
    await page.getByRole("button", { exact: true, name: "Yes" }).click();

    await expect.poll(editorText).toContain("let greeting = '';");
  }, 60_000);

  // The visitor here is a guest — no session, progress carried in a cookie the
  // browser never reads. The page is the same page a signed-in learner gets, which
  // is the point: nothing on this screen branches on whether there is an account.
  it("runs the solution and says it passed", async () => {
    expect.hasAssertions();
    await atDesktop();
    const submitted: unknown[] = [];
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.post("*/lessons/1002/check", async ({ request }) => {
        submitted.push(await request.json());
        return HttpResponse.json(checkResult());
      }),
    );

    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();

    await page.getByRole("button", { name: "Run" }).click();

    // The learner is looking at the output pane by the time the answer lands.
    await expect.element(page.getByText("Tests passed")).toBeVisible();
    await expect.element(page.getByText("1 example, 0 failures")).toBeVisible();

    // What was sent: the buffer, against the version the page was loaded with.
    expect(submitted).toStrictEqual([{ code: "let greeting = '';\n", versionId: 99 }]);
  }, 60_000);

  it("reports a failing solution without taking the lesson away", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.post("*/lessons/1002/check", () =>
        HttpResponse.json(
          checkResult({
            lessonHasBeenFinished: false,
            output: "expected 'hello', got ''",
            passed: false,
            result: "failed",
            status: 1,
          }),
        ),
      ),
    );

    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();
    await page.getByRole("button", { name: "Run" }).click();

    await expect.element(page.getByText("Tests Failed", { exact: false })).toBeVisible();
    await expect.element(page.getByText("expected 'hello', got ''")).toBeVisible();
    await expect.element(page.getByRole("button", { name: "Run" })).toBeEnabled();
  }, 60_000);

  it("reports a solution that never terminated as its own thing", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.post("*/lessons/1002/check", () =>
        HttpResponse.json(
          checkResult({
            lessonHasBeenFinished: false,
            output: "",
            passed: false,
            result: "failed-infinity",
            status: 124,
          }),
        ),
      ),
    );

    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();
    await page.getByRole("button", { name: "Run" }).click();

    await expect.element(page.getByText("Infinity Loop", { exact: false })).toBeVisible();
  }, 60_000);

  it("says when the request itself failed, rather than leaving it to be read as a wrong answer", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.post("*/lessons/1002/check", () => HttpResponse.error()),
    );

    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();
    await page.getByRole("button", { name: "Run" }).click();

    await expect
      .element(page.getByText("There was a network problem", { exact: false }))
      .toBeVisible();
    await expect.element(page.getByRole("button", { name: "Run" })).toBeEnabled();
    await expect.element(page.getByText("Tests passed")).not.toBeInTheDocument();
  }, 60_000);

  // Some courses (TypeScript) start from an empty file. The contract refuses an
  // empty submission, so running it would come back as a 400 and read as a network
  // failure; the page says what is missing instead and asks the server nothing.
  it("asks for a solution rather than sending an empty one", async () => {
    expect.hasAssertions();
    await atDesktop();
    let checks = 0;
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () =>
        HttpResponse.json(lessonView({ lesson: { ...lesson, preparedCode: "" } })),
      ),
      http.post("*/lessons/1002/check", () => {
        checks += 1;
        return HttpResponse.json(checkResult());
      }),
    );

    // An empty buffer leaves monaco's input element zero-wide, so "visible" would
    // never hold for it.
    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeInTheDocument();
    await page.getByRole("button", { name: "Run" }).click();

    await expect
      .element(page.getByText("Write your solution in the editor before running the check"))
      .toBeVisible();
    await expect
      .element(page.getByText("There was a network problem", { exact: false }))
      .not.toBeInTheDocument();
    await expect.element(page.getByRole("button", { name: "Run" })).toBeEnabled();
    expect(checks).toBe(0);
  }, 60_000);

  // The learner is moved to the output pane the moment they press run, and the
  // check can take seconds (a container starts); the pane must say so rather than
  // sit empty.
  it("shows that the check is running while it is", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.post("*/lessons/1002/check", async () => {
        await sleep(SLOW_CHECK_MS);
        return HttpResponse.json(checkResult());
      }),
    );

    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();
    await page.getByRole("button", { name: "Run" }).click();

    await expect.element(page.getByText("Running the tests…")).toBeVisible();

    await expect.element(page.getByText("Tests passed")).toBeVisible();
    await expect.element(page.getByText("Running the tests…")).not.toBeInTheDocument();
  }, 60_000);

  it("tells the learner what the output pane is for before the first run", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );

    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();
    await page.getByRole("tab", { name: "Output" }).click();

    await expect.element(page.getByText("Press Run to see the tests output")).toBeVisible();
  }, 60_000);

  // Test runners colour their report; the colours show, the escape codes do not.
  it("renders the runner's colours instead of their escape codes", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.post("*/lessons/1002/check", () =>
        HttpResponse.json(
          checkResult({
            lessonHasBeenFinished: false,
            output: "\u001B[31mexpected hello\u001B[39m",
            passed: false,
            result: "failed",
            status: 1,
          }),
        ),
      ),
    );

    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();
    await page.getByRole("button", { name: "Run" }).click();

    await expect.element(page.getByText("expected hello")).toBeVisible();
    await expect.element(page.getByText("[31m", { exact: false })).not.toBeInTheDocument();
  }, 60_000);

  it("refuses a second run while one is in flight", async () => {
    expect.hasAssertions();
    await atDesktop();
    let checks = 0;
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.post("*/lessons/1002/check", async () => {
        checks += 1;
        await sleep(QUICK_CHECK_MS);
        return HttpResponse.json(checkResult());
      }),
    );

    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();

    const run = page.getByRole("button", { name: "Run" });
    await run.click();
    await expect.element(run).toBeDisabled();

    await expect.element(page.getByText("Tests passed")).toBeVisible();
    expect(checks).toBe(1);
  }, 60_000);

  it("runs on ctrl+enter, without leaving the keyboard", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.post("*/lessons/1002/check", () => HttpResponse.json(checkResult())),
    );

    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();

    await page.getByRole("tabpanel", { name: "Editor" }).getByRole("code").click();
    await userEvent.keyboard("{Control>}{Enter}{/Control}");

    await expect.element(page.getByText("Tests passed")).toBeVisible();
  }, 60_000);

  it("shows the lesson's tests to anyone who opens them", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );

    await renderPlayer("variables");
    await page.getByRole("tab", { name: "Tests" }).click();

    await expect.element(page.getByText("assert.equal(greeting, 'hello');")).toBeVisible();
  }, 60_000);

  it("keeps the reference solution behind a wait until the lesson is passed", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.post("*/lessons/1002/check", () => HttpResponse.json(checkResult())),
    );

    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();
    await page.getByRole("tab", { name: "Solution" }).click();

    // How long is left of the wait, so a stuck learner knows whether to keep at it.
    await expect.element(page.getByText(/^\d\d:\d\d$/u)).toBeVisible();
    await expect
      .element(page.getByText("let greeting = 'hello'; // the author's answer"))
      .not.toBeInTheDocument();

    // Passing opens it, next to the learner's own code.
    await page.getByRole("button", { name: "Run" }).click();
    await expect.element(page.getByText("Tests passed")).toBeVisible();
    await page.getByRole("tab", { name: "Solution" }).click();
    await expect
      .element(page.getByText("let greeting = 'hello'; // the author's answer"))
      .toBeVisible();
    // The learner's own code, beside the author's, so the comparison is on one
    // screen — and it is the buffer, not the editor's copy of it.
    await expect
      .element(page.getByRole("tabpanel", { name: "Solution" }).getByText("let greeting = '';"))
      .toBeVisible();
  }, 60_000);

  it("opens the reference solution straight away on a lesson already finished", async () => {
    expect.hasAssertions();
    await atDesktop();
    const finished = lessonView({
      lessons: lessonRows(["hello-world", "variables"]),
      progress: { completion: 66, nextLesson: lessonRef("strings"), state: "started" },
    });
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(finished)),
    );

    await renderPlayer("variables");
    await page.getByRole("tab", { name: "Solution" }).click();

    await expect
      .element(page.getByText("let greeting = 'hello'; // the author's answer"))
      .toBeVisible();
  }, 60_000);

  it("gives a signed-in learner the same run a guest gets", async () => {
    expect.hasAssertions();
    await atDesktop();
    const submitted: unknown[] = [];
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.post("*/lessons/1002/check", async ({ request }) => {
        submitted.push(await request.json());
        return HttpResponse.json(checkResult());
      }),
    );

    await renderPlayer("variables", learner);
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();
    await page.getByRole("button", { name: "Run" }).click();

    await expect.element(page.getByText("Tests passed")).toBeVisible();
    expect(submitted).toStrictEqual([{ code: "let greeting = '';\n", versionId: 99 }]);
  }, 60_000);

  it("previews a markup course's page under the editor, as the learner writes it", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/html/lessons/headings", () => HttpResponse.json(htmlLessonView())),
    );

    await renderRoute(lessonRoute, {
      initialPath: "/languages/html/lessons/headings",
      path: "/{-$locale}/languages/$slug/lessons/$lessonSlug",
      wrap: (element) => <div style={{ height: "800px", width: "1200px" }}>{element}</div>,
    });

    const preview = page.getByTitle("Preview");
    await expect.element(preview, editorLoad).toBeVisible();
    await expect.element(preview).toHaveAttribute("srcdoc", "<h1>Hello</h1>\n");

    // The learner's markup cannot reach the page around it: no scripts, and an
    // opaque origin whose document the page cannot open either.
    const frame = preview.element();
    expect(frame).toBeInstanceOf(HTMLIFrameElement);
    expect(frame.getAttribute("sandbox")).toBe("");
    expect(Reflect.get(frame, "contentDocument")).toBeNull();

    // It follows the buffer, not the starter code it opened with.
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();
    await page.getByRole("tabpanel", { name: "Editor" }).getByRole("code").click();
    await userEvent.keyboard("{Control>}{End}{/Control}mine");

    await expect.poll(() => frame.getAttribute("srcdoc")).toContain("mine");
  }, 60_000);

  it("leaves a program course's editor without a preview", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );

    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();

    expect(document.querySelector("iframe")).toBeNull();
  }, 60_000);

  // ---- Progression -----------------------------------------------------------

  it("carries reset, Previous, Run and Next under the editor, in legacy's order", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );

    await renderPlayer("variables", learner);
    await expect.element(page.getByRole("button", { name: "Run" })).toBeVisible();

    const bar = page.getByRole("button", { name: "Run" }).element().parentElement;
    expect(controlLabels(bar)).toStrictEqual(["Reset", "← Previous", "Run", "Next →"]);
  }, 60_000);

  it("goes back with a plain link, and has nowhere to go back to on the first lesson", async () => {
    expect.hasAssertions();
    await atDesktop();
    const first = lessonView({
      lesson: { ...lesson, id: 1001, name: "Hello, World!", slug: "hello-world" },
    });
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.get("*/languages/javascript/lessons/hello-world", () => HttpResponse.json(first)),
    );

    const { screen } = await renderPlayer("variables", learner);
    await expect
      .element(page.getByRole("link", { name: "← Previous" }))
      .toHaveAttribute("href", "/languages/javascript/lessons/hello-world");

    await screen.unmount();
    await renderPlayer("hello-world", learner);
    await expect.element(page.getByRole("button", { name: "← Previous" })).toBeDisabled();
  }, 60_000);

  it("moves a learner on: Next waits for the pass, then starts the next lesson before going there", async () => {
    expect.hasAssertions();
    await atDesktop();
    const strings = lessonView({
      lesson: {
        ...lesson,
        id: 1003,
        name: "Strings",
        preparedCode: "// strings\n",
        slug: "strings",
      },
    });
    const requests: string[] = [];
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.post("*/lessons/1002/check", () => HttpResponse.json(checkResult())),
      http.post("*/lessons/:id/start", ({ params }) => {
        requests.push(`start ${String(params.id)}`);
        return new HttpResponse(null, { status: 204 });
      }),
      http.get("*/languages/javascript/lessons/strings", () => {
        requests.push("read strings");
        return HttpResponse.json(strings);
      }),
    );

    const { router } = await renderPlayer("variables", learner);
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();

    // The button says what is expected: pass first.
    const next = page.getByRole("button", { name: "Next →" });
    await expect.element(next).toBeDisabled();

    await page.getByRole("button", { name: "Run" }).click();
    await expect.element(page.getByText("Tests passed")).toBeVisible();
    await expect.element(next).toBeEnabled();

    await next.click();

    await expect.element(page.getByRole("heading", { name: `${COURSE}: Strings` })).toBeVisible();
    expect(router.state.location.pathname).toBe("/languages/javascript/lessons/strings");
    // Started exactly once, and before the page it leads to was read (ADR-0012).
    expect(requests).toStrictEqual(["start 1003", "read strings"]);
  }, 60_000);

  it("lets a learner revisiting a finished lesson move on straight away", async () => {
    expect.hasAssertions();
    await atDesktop();
    const { lessons, progress } = lastLessonView();
    const finished = lessonView({ lessons, progress });
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(finished)),
    );

    await renderPlayer("variables", learner);

    // Revisiting is not the same as being stuck: no run needed.
    await expect.element(page.getByRole("button", { name: "Next →" })).toBeEnabled();
  }, 60_000);

  it("says so on the last lesson, and leads to the completion page", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/strings", () => HttpResponse.json(lastLessonView())),
      http.post("*/lessons/1003/check", () => HttpResponse.json(checkResult())),
      // The pass has finished the Enrollment on the server. The course is read by
      // its landing page's slug, which is what the completion page is addressed by.
      http.get("*/api/languages/javascript-ru", () => HttpResponse.json(courseView("finished"))),
    );

    const { router, queryClient } = await renderRoute(lessonRoute, {
      also: [{ path: "/{-$locale}/languages/$slug/success", route: successRoute }],
      initialPath: "/languages/javascript/lessons/strings",
      path: "/{-$locale}/languages/$slug/lessons/$lessonSlug",
      user: learner,
      wrap: (element) => (
        <div style={{ height: `${desktop.height}px`, width: `${desktop.width}px` }}>{element}</div>
      ),
    });
    // A course read left from an earlier visit to the Course page, taken before
    // the pass: served as it is, it would bounce the learner back as unfinished.
    queryClient.setQueryData(
      getCourseQueryKey({ path: { slug: "javascript-ru" } }),
      courseView("started"),
    );
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();

    await expect.element(page.getByRole("button", { name: "Finish" })).toBeDisabled();
    await expect.element(page.getByRole("button", { name: "Next →" })).not.toBeInTheDocument();

    await page.getByRole("button", { name: "Run" }).click();
    await expect.element(page.getByText("Tests passed")).toBeVisible();

    const finish = page.getByRole("link", { name: "Finish" });
    await expect.element(finish).toHaveAttribute("href", "/languages/javascript-ru/success");
    await finish.click();

    await expect
      .element(page.getByRole("heading", { name: "Congratulations, you completed the course!" }))
      .toBeVisible();
    expect(router.state.location.pathname).toBe("/languages/javascript-ru/success");
  }, 60_000);

  it("offers a guest sign-up rather than Finish on the last lesson", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/strings", () => HttpResponse.json(lastLessonView())),
      http.post("*/lessons/1003/check", () => HttpResponse.json(checkResult())),
    );

    await renderPlayer("strings");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();

    await page.getByRole("button", { name: "Run" }).click();
    await expect.element(page.getByText("Tests passed")).toBeVisible();

    // The completion page needs a session, so a guest is asked to sign up first.
    await expect.element(page.getByRole("link", { name: "Finish" })).not.toBeInTheDocument();
    await expect
      .element(page.getByRole("link", { exact: true, name: "Next →" }))
      .toHaveAttribute("href", expect.stringContaining("/users/new"));
  }, 60_000);

  it("offers a guest sign-up in place of Next, and tells them it keeps their progress", async () => {
    expect.hasAssertions();
    await atDesktop();
    let started = 0;
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.post("*/lessons/1002/check", () => HttpResponse.json(checkResult())),
      http.post("*/lessons/:id/start", () => {
        started += 1;
        return HttpResponse.json({});
      }),
    );

    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();

    const prompt = "Be sure to register to ensure you don't lose the results you've achieved";
    await expect.element(page.getByRole("button", { name: "Next →" })).toBeDisabled();
    await expect.element(page.getByText(prompt)).not.toBeInTheDocument();

    await page.getByRole("button", { name: "Run" }).click();
    await expect.element(page.getByText("Tests passed")).toBeVisible();
    await expect.element(page.getByText(prompt)).toBeVisible();

    // Both lead to sign-up and carry the way back to this lesson. The merge that
    // credits what was passed is the server's, and is covered by its tests.
    for (const name of ["Next →", "register"]) {
      const href = page.getByRole("link", { exact: true, name }).element().getAttribute("href");
      const url = new URL(String(href), location.origin);
      expect(url.pathname).toBe("/users/new");
      expect(url.searchParams.get("redirect")).toBe("/languages/javascript/lessons/variables");
    }

    // A guest's Next starts nothing: it leads to sign-up, not to the next lesson.
    expect(started).toBe(0);
  }, 60_000);

  it("tells a guest nothing on a lesson they finished before, but lets them sign up", async () => {
    expect.hasAssertions();
    await atDesktop();
    const { lessons, progress } = lastLessonView();
    const finished = lessonView({ lessons, progress });
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(finished)),
    );

    await renderPlayer("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();

    // The way on is open, as it is for a learner revisiting a finished lesson…
    await expect
      .element(page.getByRole("link", { exact: true, name: "Next →" }))
      .toHaveAttribute("href", expect.stringContaining("/users/new"));
    // …but the prompt is about a pass just achieved, and there has been none.
    await expect
      .element(
        page.getByText("Be sure to register to ensure you don't lose the results you've achieved"),
      )
      .not.toBeInTheDocument();
  }, 60_000);

  // ---- On a phone ------------------------------------------------------------

  it("on a phone, fills the screen with one pane and swaps them with the burger", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );

    await renderOnPhone("variables");
    await expect.element(page.getByRole("tab", { name: "Editor" })).toBeVisible();

    // The workspace first: it is where the learner acts.
    await expect.poll(paneWidths).toEqual([0, phone.width]);

    await panes().workspace.getByRole("button", { name: "Navigation" }).click();
    await expect.poll(paneWidths).toEqual([phone.width, 0]);
    await expect
      .element(panes().theory.getByRole("heading", { name: `${COURSE}: Variables` }))
      .toBeInViewport();

    // The theory pane has a burger of its own, which swaps back.
    await panes().theory.getByRole("button", { name: "Navigation" }).click();
    await expect.poll(paneWidths).toEqual([0, phone.width]);
    await expect.element(page.getByRole("tab", { name: "Editor" })).toBeInViewport();
  }, 60_000);

  it("on a phone, keeps what the learner typed across a pane swap and the theory tab", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );

    await renderOnPhone("variables");
    await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();
    const editor = document.querySelector(".monaco-editor");

    await page.getByRole("tabpanel", { name: "Editor" }).getByRole("code").click();
    await userEvent.keyboard("// mine");
    await expect.poll(editorText).toContain("// mine");

    // To the theory pane and back.
    await panes().workspace.getByRole("button", { name: "Navigation" }).click();
    await expect.poll(paneWidths).toEqual([phone.width, 0]);
    await panes().theory.getByRole("button", { name: "Navigation" }).click();
    await expect.poll(paneWidths).toEqual([0, phone.width]);
    await expect.poll(editorText).toContain("// mine");

    // The theory without leaving the editor's side: the workspace's own tab.
    await panes().workspace.getByRole("tab", { name: "Lesson" }).click();
    await expect
      .element(panes().workspace.getByRole("heading", { name: `${COURSE}: Variables` }))
      .toBeInViewport();
    await panes().workspace.getByRole("tab", { name: "Editor" }).click();
    await expect.element(page.getByLabelText("Code editor")).toBeVisible();
    await expect.poll(editorText).toContain("// mine");

    // The same editor all along, never remounted.
    expect(document.querySelector(".monaco-editor")).toBe(editor);
  }, 60_000);

  it("on a phone, keeps Run, Previous and Next on the screen", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );

    await renderOnPhone("variables", learner);

    await expect.element(page.getByRole("button", { name: "Run" })).toBeInViewport({ ratio: 1 });
    await expect
      .element(page.getByRole("link", { name: "← Previous" }))
      .toBeInViewport({ ratio: 1 });
    await expect.element(page.getByRole("button", { name: "Next →" })).toBeInViewport({ ratio: 1 });
  }, 60_000);

  it("on a phone, has no divider to drag", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    );

    await renderOnPhone("variables");
    await expect.poll(paneWidths).toEqual([0, phone.width]);

    // No grip to take hold of, and the keys that move a desktop's divider move
    // nothing here.
    expect(document.querySelector(".mantine-Splitter-thumb")).toBeNull();
    document.querySelector<HTMLElement>(".mantine-Splitter-handle")?.focus();
    await userEvent.keyboard("{ArrowRight}{ArrowRight}{ArrowLeft}");
    expect(paneWidths()).toStrictEqual([0, phone.width]);
  }, 60_000);

  it("on a phone, fits the theory to the screen, wide code and all", async () => {
    expect.hasAssertions();
    await atDesktop();
    const wide = lessonView({
      lesson: {
        ...lesson,
        theory: `A long line:\n\n\`\`\`js\nconst x = '${"a".repeat(LONG_LINE_LENGTH)}';\n\`\`\``,
      },
    });
    worker.use(http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(wide)));

    await renderOnPhone("variables");
    await expect.poll(paneWidths).toEqual([0, phone.width]);
    await panes().workspace.getByRole("tab", { name: "Lesson" }).click();
    await expect.element(panes().workspace.getByText("A long line:")).toBeInViewport();

    // The page itself never scrolls sideways: the wide block scrolls inside its
    // own box, so the prose around it reads at the screen's width.
    expect(document.documentElement.scrollWidth).toBeLessThanOrEqual(phone.width);
  }, 60_000);

  // ---- The assistant -----------------------------------------------------------

  it("asks a visitor to sign in before the assistant helps", async () => {
    expect.hasAssertions();
    await atDesktop();
    let listed = 0;
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.get("*/api/ai/lessons/1002/messages", () => {
        listed += 1;
        return HttpResponse.json(assistantChat());
      }),
    );

    await renderPlayer("variables");
    await page.getByRole("tab", { name: "Assistant" }).click();

    await expect.element(page.getByText(/please sign up or log in/u)).toBeVisible();
    expect(listed).toBe(0);
  }, 60_000);

  it("shows the learner's chat and streams the answer to a new question", async () => {
    expect.hasAssertions();
    await atDesktop();
    const asked: unknown[] = [];
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.get("*/api/ai/lessons/1002/messages", () => HttpResponse.json(assistantChat())),
      http.post("*/api/ai/lessons/1002/messages", async ({ request }) => {
        asked.push(await request.json());
        return streamedAnswer(["Use ", "`let` ", "and assign it."]);
      }),
    );

    await renderPlayer("variables", learner);
    await page.getByRole("tab", { name: "Assistant" }).click();

    await expect.element(page.getByText("What is a variable?")).toBeVisible();
    await expect.element(page.getByText("A name bound to a value.")).toBeVisible();

    await page.getByRole("textbox", { name: "Your question" }).fill("How do I declare one?");
    await page.getByRole("button", { name: "Send" }).click();

    await expect.element(page.getByText("How do I declare one?")).toBeVisible();
    await expect.element(page.getByText(/and assign it\./u)).toBeVisible();

    // The question goes out with the editor's buffer; nothing has been run yet.
    expect(asked).toStrictEqual([
      { message: "How do I declare one?", output: "", userCode: "let greeting = '';\n" },
    ]);
  }, 60_000);

  it("switches the assistant off once today's questions are used up", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.get("*/api/ai/lessons/1002/messages", () =>
        HttpResponse.json(assistantChat({ quotaExceeded: true })),
      ),
    );

    await renderPlayer("variables", learner);
    await page.getByRole("tab", { name: "Assistant" }).click();

    await expect.element(page.getByText(/reached today's message limit/u)).toBeVisible();
    await expect
      .element(page.getByRole("link", { name: "Telegram community" }))
      .toHaveAttribute("href", "https://t.me/hexletcommunity");
    expect(page.getByRole("button", { name: "Send" }).elements()).toHaveLength(0);
  }, 60_000);

  it("switches the assistant off when the server refuses a question over the quota", async () => {
    expect.hasAssertions();
    await atDesktop();
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
      http.get("*/api/ai/lessons/1002/messages", () => HttpResponse.json(assistantChat())),
      http.post("*/api/ai/lessons/1002/messages", () =>
        HttpResponse.json(
          { status: 429, title: "Too Many Requests", type: "about:blank" },
          { headers: { "Content-Type": "application/problem+json" }, status: 429 },
        ),
      ),
    );

    await renderPlayer("variables", learner);
    await page.getByRole("tab", { name: "Assistant" }).click();
    await page.getByRole("textbox", { name: "Your question" }).fill("One more?");
    await page.getByRole("button", { name: "Send" }).click();

    await expect.element(page.getByText(/reached today's message limit/u)).toBeVisible();
  }, 60_000);

  it("underlines a syntax error in a language monaco has no service of", async () => {
    expect.hasAssertions();
    await atDesktop();
    // C has no language service in monaco, so the underline can only come from the
    // syntax check.
    const cLesson = lessonView({
      lesson: {
        ...lesson,
        exerciseExtension: "c",
        preparedCode: "int main(void) { return 0; }\n",
      },
    });
    worker.use(
      http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(cLesson)),
    );

    await renderPlayer("variables");
    await expect.poll(editorText, editorLoad).toContain("int main(void)");
    expect(squiggles()).toBe(0);

    // A declaration with no value. No bracket, so monaco's auto-closing stays out of it.
    await page.getByRole("tabpanel", { name: "Editor" }).getByRole("code").click();
    await userEvent.keyboard("{Control>}{End}{/Control}int x = ;");
    await expect.poll(squiggles, editorLoad).toBeGreaterThan(0);

    // Giving it one clears the underline.
    await userEvent.keyboard("{ArrowLeft}1");
    await expect.poll(squiggles, editorLoad).toBe(0);
  }, 60_000);
});
