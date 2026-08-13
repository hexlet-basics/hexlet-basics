import { http, HttpResponse } from "msw";
import { afterEach, expect, test } from "vitest";
import { page, userEvent } from "vitest/browser";
import type { Course, CourseLandingPage, CourseLesson, CourseLessonView } from "@/client/types.gen";
import { Route as lessonRoute } from "@/routes/{-$locale}/languages/$slug/lessons/$lessonSlug";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The lesson player, driven through its real route with the API faked at the
// HTTP boundary — the seam the admin screens are tested at. What is asserted is
// what a learner can see: the words on the page, the marks in the lesson list,
// and the requests the page did (and did not) make.

const course: Course = {
  id: 1,
  slug: "javascript",
  // The course's own name, which is NOT what the page is titled with.
  name: "javascript",
  learnAs: null,
  readiness: null,
  categoryId: null,
  currentVersionId: 99,
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

const landingPage: CourseLandingPage = {
  id: 10,
  courseId: 1,
  courseSlug: "javascript",
  createdAt: "2026-01-01T00:00:00Z",
  slug: "javascript-ru",
  // The marketing name a learner recognises — the one the page is titled with.
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
  usedInHeader: null,
  usedInDescription: null,
  outcomesHeader: null,
  outcomesDescription: null,
  outcomesImage: null,
  duration: 10,
  enrollmentsCount: 0,
};

const lesson: CourseLesson = {
  course,
  id: 1002,
  name: "Variables",
  slug: "variables",
  locale: "en",
  naturalOrder: 2,
  versionId: 99,
  version: 99,
  description: "Storing values",
  instructions: "Assign the string `hello` to `greeting`.",
  theory: "A variable is a name bound to a value.\n\n```js\nlet greeting = 'hello';\n```",
  definitions: [],
  tips: ["Names are case sensitive"],
  preparedCode: "let greeting = '';\n",
  originalCode: null,
  testCode: null,
  sourceCodeUrl: "https://github.com/hexlet-basics/exercises-javascript",
  createdAt: "2026-01-01T00:00:00Z",
};

// A learner who has finished the first lesson: the second is theirs to take, the
// third is still behind the gate.
function lessonView(overrides: Partial<CourseLessonView> = {}): CourseLessonView {
  return {
    lesson,
    landingPage,
    lessons: [
      { id: 1001, name: "Hello, World!", description: null, slug: "hello-world" },
      { id: 1002, name: "Variables", description: null, slug: "variables" },
      { id: 1003, name: "Strings", description: null, slug: "strings" },
    ],
    progress: {
      state: "started",
      completion: 33,
      nextLessonSlug: "variables",
      furthestFinishedPosition: 1,
      lessons: [
        { slug: "hello-world", position: 1, finished: true, available: true },
        { slug: "variables", position: 2, finished: false, available: true },
        { slug: "strings", position: 3, finished: false, available: false },
      ],
    },
    ...overrides,
  };
}

// The real route at a real URL, so the loader and the route's chrome run as they
// do for a visitor.
//
// The player needs a container with a real height: the two panes divide the
// space they are given, and a zero-height box renders nothing a learner sees.
function renderPlayer(lessonSlug: string) {
  return renderRoute(lessonRoute, {
    path: "/{-$locale}/languages/$slug/lessons/$lessonSlug",
    initialPath: `/languages/javascript/lessons/${lessonSlug}`,
    wrap: (element) => <div style={{ height: "800px", width: "1200px" }}>{element}</div>,
  });
}

// What monaco has painted, with the non-breaking spaces it renders text with
// turned back into ordinary ones so assertions read like the code does.
function editorText() {
  return (document.querySelector(".view-lines")?.textContent ?? "").replaceAll("\u00a0", " ");
}

// The editor is a lazy chunk of a few hundred modules, and the first test to
// open it waits for vite to transform all of them, so waiting on it gets a
// budget of its own rather than the default second. The budget is sized for a
// cold CI runner, not for a warm laptop.
const editorLoad = { timeout: 30_000 };

// Everything the player keeps between visits lives here; a test that seeds it
// must not leak into the next one.
afterEach(() => localStorage.clear());

test("renders the lesson a learner reads, titled from the course's landing copy", async () => {
  worker.use(
    http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
  );

  await renderPlayer("variables");

  // The landing page's marketing name, not the course's own `name`.
  await expect.element(page.getByRole("heading", { name: "JavaScript: Variables" })).toBeVisible();
  await expect.element(page.getByText("A variable is a name bound to a value.")).toBeVisible();
  await expect.element(page.getByRole("heading", { name: "Instructions" })).toBeVisible();
  await expect.element(page.getByText("Assign the string")).toBeVisible();
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
});

test("lists every lesson in course order, marking finished and locked ones", async () => {
  worker.use(
    http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
  );

  await renderPlayer("variables");
  await page.getByRole("tab", { name: "Navigation" }).click();

  const links = page.getByRole("link").elements();
  const lessonLinks = links.filter((el) => el.getAttribute("href")?.includes("/lessons/"));
  expect(lessonLinks.map((el) => el.textContent?.trim())).toEqual([
    "Hello, World!",
    "Variables",
    "Strings",
  ]);

  await expect.element(page.getByLabelText("Finished")).toBeVisible();
  await expect.element(page.getByLabelText("Locked")).toBeVisible();

  // Where the learner is now.
  await expect
    .element(page.getByRole("link", { name: "Variables" }))
    .toHaveAttribute("aria-current", "page");

  // Locked stays clickable: theory is public, only the exercise is gated.
  await expect
    .element(page.getByRole("link", { name: "Strings" }))
    .toHaveAttribute("href", "/languages/javascript/lessons/strings");
});

test("serves a visitor a locked lesson in full", async () => {
  const locked = lessonView({
    progress: {
      state: null,
      completion: 0,
      nextLessonSlug: "hello-world",
      furthestFinishedPosition: 0,
      lessons: [
        { slug: "hello-world", position: 1, finished: false, available: true },
        { slug: "variables", position: 2, finished: false, available: false },
        { slug: "strings", position: 3, finished: false, available: false },
      ],
    },
  });
  worker.use(http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(locked)));

  await renderPlayer("variables");

  // Readable in full — this is the page a search engine indexes.
  await expect.element(page.getByRole("heading", { name: "JavaScript: Variables" })).toBeVisible();
  await expect.element(page.getByText("A variable is a name bound to a value.")).toBeVisible();

  await page.getByRole("tab", { name: "Navigation" }).click();
  await expect.element(page.getByLabelText("Locked").first()).toBeVisible();
});

test("starts nothing by loading the page", async () => {
  let started = 0;
  worker.use(
    http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
    http.post("*/lessons/:id/start", () => {
      started += 1;
      return HttpResponse.json({});
    }),
  );

  await renderPlayer("variables");
  await expect.element(page.getByRole("heading", { name: "JavaScript: Variables" })).toBeVisible();
  await page.getByRole("tab", { name: "Navigation" }).click();
  await expect.element(page.getByRole("link", { name: "Strings" })).toBeVisible();

  // The router preloads on hover, so a read that started a lesson would enroll a
  // learner in every lesson they pointed at (ADR-0012).
  expect(started).toBe(0);
});

test("divides the two panes where the learner left the divider", async () => {
  worker.use(
    http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
  );
  // What a previous visit to this lesson stored — a reload reads it back rather
  // than resetting to the default split.
  localStorage.setItem("lesson-panes-javascript-variables", JSON.stringify(["70%", "30%"]));

  await renderPlayer("variables");
  await expect.element(page.getByRole("heading", { name: "JavaScript: Variables" })).toBeVisible();

  // The handle is what makes the divider draggable at all.
  expect(document.querySelectorAll(".mantine-Splitter-handle").length).toBe(1);

  const panes = document.querySelectorAll<HTMLElement>(".mantine-Splitter-pane");
  expect(panes.length).toBe(2);
  expect(panes[0]?.style.cssText).toContain("70%");

  localStorage.removeItem("lesson-panes-javascript-variables");
});

test("says so when the lesson is not there", async () => {
  worker.use(
    http.get("*/languages/javascript/lessons/nope", () =>
      HttpResponse.json({ title: "Not Found", status: 404 }, { status: 404 }),
    ),
  );

  await renderPlayer("nope");

  await expect
    .element(page.getByText("Lesson not found. Please, try another lesson."))
    .toBeVisible();
});

test("opens the editor with the lesson's starter code, after the theory", async () => {
  worker.use(
    http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
  );

  await renderPlayer("variables");

  // The theory does not wait for the editor: it renders from the payload, while
  // monaco arrives on its own lazy import afterwards.
  await expect.element(page.getByRole("heading", { name: "JavaScript: Variables" })).toBeVisible();

  await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();
  await expect.poll(editorText, editorLoad).toContain("let greeting = '';");

  // Monaco is this application's own copy, not a CDN's.
  const remote = performance
    .getEntriesByType("resource")
    .filter((entry) => !entry.name.startsWith(location.origin));
  expect(remote).toEqual([]);
});

test("finds yesterday's work still there, and keeps it per lesson", async () => {
  const strings = lessonView({
    lesson: { ...lesson, id: 1003, name: "Strings", slug: "strings", preparedCode: "// strings\n" },
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
});

test("keeps what the learner typed when they come back to the lesson", async () => {
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
});

test("asks before resetting, and restores the starter code", async () => {
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
  await expect.element(page.getByText("You want to reset the exercise progress.")).toBeVisible();
  await page.getByRole("button", { name: "No", exact: true }).click();
  expect(editorText()).toContain("let greeting = 'mess';");

  await page.getByRole("button", { name: "Reset" }).click();
  await page.getByRole("button", { name: "Yes", exact: true }).click();

  await expect.poll(editorText).toContain("let greeting = '';");
});

test("tells the learner about autocomplete once", async () => {
  worker.use(
    http.get("*/languages/javascript/lessons/variables", () => HttpResponse.json(lessonView())),
  );

  const { screen } = await renderPlayer("variables");
  await expect
    .element(page.getByText("The editor suggests commands as you type"), editorLoad)
    .toBeVisible();

  await page.getByRole("button", { name: "Dismiss" }).click();
  await expect
    .element(page.getByText("The editor suggests commands as you type"))
    .not.toBeInTheDocument();

  // And it stays dismissed for the next lesson, and the next visit.
  await screen.unmount();
  await renderPlayer("variables");
  await expect.element(page.getByLabelText("Code editor"), editorLoad).toBeVisible();
  await expect
    .element(page.getByText("The editor suggests commands as you type"))
    .not.toBeInTheDocument();
});
