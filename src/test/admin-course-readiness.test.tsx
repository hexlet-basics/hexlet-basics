import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { Course, CourseReadiness } from "@/client/types.gen";
import { Route as courseEditRoute } from "@/routes/{-$locale}/admin/courses/$id";
import { Route as courseListRoute } from "@/routes/{-$locale}/admin/courses/index";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// A Course's Readiness is changed by action buttons, one per operation, never by
// the form (the same pattern as the Publication State panels, see
// admin-publication-state.test.tsx). Through the real routes under the real admin
// guard: the buttons offered are the values other than the current one, a click
// calls that operation and the course shown is the one it answered with.

const course = (readiness: CourseReadiness | null): Course =>
  ({
    id: 7,
    slug: "ruby",
    learnAs: "first_language",
    readiness,
    hexletProgramLandingPage: null,
    repositoryUrl: null,
    currentVersion: null,
    lessonsCount: 3,
    createdAt: "2026-01-01T00:00:00Z",
  }) as unknown as Course;

const edit = () =>
  renderRoute(courseEditRoute, {
    path: "/{-$locale}/admin/courses/$id",
    initialPath: "/admin/courses/7",
    user: admin,
  });

const button = (name: string) => page.getByRole("button", { name, exact: true });

describe("the course admin screens", () => {
  it.each<[CourseReadiness, string[], string[]]>([
    ["draft", ["Mark as in development", "Mark as completed"], ["Mark as draft"]],
    ["in_development", ["Mark as draft", "Mark as completed"], ["Mark as in development"]],
    ["completed", ["Mark as draft", "Mark as in development"], ["Mark as completed"]],
  ])(
    "a %s course offers every readiness but its own",
    async (readiness, offered, hidden) => {
      expect.hasAssertions();
      worker.use(http.get("*/admin/courses/7", () => HttpResponse.json(course(readiness))));

      await edit();

      await Promise.all([
        ...offered.map(async (name) => {
          await expect.element(button(name)).toBeVisible();
        }),
        ...hidden.map(async (name) => {
          await expect.element(button(name)).not.toBeInTheDocument();
        }),
      ]);
    },
    60_000,
  );

  it("the form has no readiness field", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/admin/courses/7", () => HttpResponse.json(course("draft"))));

    await edit();

    await expect.element(page.getByLabelText("Slug *")).toBeVisible();
    await expect.element(page.getByLabelText("Course readiness")).not.toBeInTheDocument();
  }, 60_000);

  it("saving the form sends no readiness", async () => {
    expect.hasAssertions();
    let sent: Record<string, unknown> | null = null;
    worker.use(
      http.get("*/admin/courses/7", () => HttpResponse.json(course("draft"))),
      http.put("*/admin/courses/7", async ({ request }) => {
        sent = (await request.json()) as Record<string, unknown>;
        return HttpResponse.json(course("draft"));
      }),
    );

    await edit();
    await page.getByRole("button", { name: "Save" }).click();

    await expect.poll(() => sent).toBeDefined();
    expect(sent).not.toHaveProperty("readiness");
  }, 60_000);

  it("marking a course calls that operation and shows the new readiness", async () => {
    expect.hasAssertions();
    const called: string[] = [];
    let readiness: CourseReadiness = "draft";
    worker.use(
      http.get("*/admin/courses/7", () => HttpResponse.json(course(readiness))),
      http.post("*/admin/courses/7/mark_completed", () => {
        called.push("mark_completed");
        readiness = "completed";
        return HttpResponse.json(course(readiness));
      }),
      http.post("*/admin/courses/7/mark_draft", () => {
        called.push("mark_draft");
        readiness = "draft";
        return HttpResponse.json(course(readiness));
      }),
    );

    await edit();
    await button("Mark as completed").click();
    await expect.element(page.getByText("Readiness changed").first()).toBeVisible();
    await expect.element(button("Mark as completed")).not.toBeInTheDocument();
    await button("Mark as draft").click();
    await expect.element(button("Mark as draft")).not.toBeInTheDocument();

    expect(called).toStrictEqual(["mark_completed", "mark_draft"]);
  }, 60_000);

  it("a refused change tells the admin and keeps the readiness", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/admin/courses/7", () => HttpResponse.json(course("draft"))),
      http.post("*/admin/courses/7/mark_in_development", () =>
        HttpResponse.json({ type: "about:blank", title: "Conflict", status: 409 }, { status: 409 }),
      ),
    );

    await edit();
    await button("Mark as in development").click();

    await expect.element(page.getByText("Could not change the readiness")).toBeVisible();
    await expect.element(button("Mark as in development")).toBeVisible();
  }, 60_000);

  it("the list shows each course's readiness in words", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/admin/courses", () =>
        HttpResponse.json({ items: [course("in_development")], total: 1, page: 1, perPage: 25 }),
      ),
    );

    await renderRoute(courseListRoute, {
      path: "/{-$locale}/admin/courses/",
      initialPath: "/admin/courses",
      user: admin,
    });

    await expect.element(page.getByRole("cell", { name: "In development" })).toBeVisible();
  }, 60_000);
});
