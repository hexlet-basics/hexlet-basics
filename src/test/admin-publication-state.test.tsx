import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { Banner, CourseLandingPage, Review } from "@/client/types.gen";
import type { PublicationState } from "@/lib/publicationState";
import { Route as bannerEditRoute } from "@/routes/{-$locale}/admin/banners/$id";
import { Route as bannerListRoute } from "@/routes/{-$locale}/admin/banners/index";
import { Route as landingEditRoute } from "@/routes/{-$locale}/admin/course_landing_pages/$id";
import { Route as landingListRoute } from "@/routes/{-$locale}/admin/course_landing_pages/index";
import { Route as reviewEditRoute } from "@/routes/{-$locale}/admin/reviews/$id";
import { Route as reviewListRoute } from "@/routes/{-$locale}/admin/reviews/index";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";
import type { Locator } from "vitest/browser";

// The Publication State of a banner, a review and a landing page is changed by
// action buttons, one per operation, never by the form (the same pattern as the
// blog post, admin-blog-post-state.test.tsx). Through the real routes under the
// real admin guard: the buttons offered follow the current state, a click calls
// that operation and the record shown is the one it answered with.

const banner = (state: PublicationState): Banner => ({
  background: "cta_gradient",
  body: "Promo",
  createdAt: "2026-01-01T00:00:00Z",
  finishesAt: null,
  id: 7,
  locale: "en",
  startsAt: null,
  state,
  url: null,
});

// The embedded account and course are only read for their name and slug.
const review = (state: PublicationState): Review =>
  ({
    body: "Great course",
    course: { slug: "ruby" },
    courseId: 2,
    createdAt: "2026-01-01T00:00:00Z",
    firstName: "Eva",
    fullName: "Eva",
    id: 7,
    lastName: null,
    locale: "en",
    pinned: false,
    state,
    user: { firstName: "Eva", id: 1, lastName: null, name: "Eva" },
    userId: 1,
  }) as unknown as Review;

const landing = (state: PublicationState): CourseLandingPage => ({
  courseId: 2,
  courseSlug: "ruby",
  createdAt: "2026-01-01T00:00:00Z",
  description: "About Ruby",
  duration: 10,
  enrollmentsCount: 0,
  footer: false,
  footerName: null,
  header: "Ruby",
  id: 7,
  landingPageToRedirectId: null,
  listed: true,
  main: false,
  metaDescription: "Ruby course",
  metaTitle: "Ruby",
  name: "Ruby",
  order: null,
  outcomesDescription: null,
  outcomesHeader: null,
  outcomesImage: null,
  slug: "ruby-ru",
  state,
  usedInDescription: null,
  usedInHeader: null,
});

interface Subject {
  name: string;
  resource: string;
  record: (state: PublicationState) => Banner | Review | CourseLandingPage;
  edit: () => ReturnType<typeof renderRoute>;
  list: () => ReturnType<typeof renderRoute>;
  // What the list's cell shows for the state "archived" and the form's own label.
  nameLabel: string;
}

const subjects: Subject[] = [
  {
    edit: () =>
      renderRoute(bannerEditRoute, {
        initialPath: "/admin/banners/7",
        path: "/{-$locale}/admin/banners/$id",
        user: admin,
      }),
    list: () =>
      renderRoute(bannerListRoute, {
        initialPath: "/admin/banners",
        path: "/{-$locale}/admin/banners/",
        user: admin,
      }),
    name: "banner",
    nameLabel: "Text (Markdown) *",
    record: banner,
    resource: "banners",
  },
  {
    edit: () =>
      renderRoute(reviewEditRoute, {
        initialPath: "/admin/reviews/7",
        path: "/{-$locale}/admin/reviews/$id",
        user: admin,
      }),
    list: () =>
      renderRoute(reviewListRoute, {
        initialPath: "/admin/reviews",
        path: "/{-$locale}/admin/reviews/",
        user: admin,
      }),
    name: "review",
    nameLabel: "Review *",
    record: review,
    resource: "reviews",
  },
  {
    edit: () =>
      renderRoute(landingEditRoute, {
        initialPath: "/admin/course_landing_pages/7",
        path: "/{-$locale}/admin/course_landing_pages/$id",
        user: admin,
      }),
    list: () =>
      renderRoute(landingListRoute, {
        initialPath: "/admin/course_landing_pages",
        path: "/{-$locale}/admin/course_landing_pages/",
        user: admin,
      }),
    name: "course landing page",
    nameLabel: "Name *",
    record: landing,
    resource: "course_landing_pages",
  },
];

// The edit screens of the landing page also load its FAQ and nothing else.
const noQna = http.get("*/admin/course_landing_pages/7/qna_items", () => HttpResponse.json([]));

const button = (name: string): Locator => page.getByRole("button", { exact: true, name });

describe.each(subjects)("the $name admin screens", (subject) => {
  it.each<[PublicationState, string[], string[]]>([
    ["draft", ["Publish", "Archive"], ["Unpublish"]],
    ["published", ["Unpublish", "Archive"], ["Publish"]],
    ["archived", ["Publish"], ["Unpublish", "Archive"]],
  ])(
    "a %s record offers only the transitions legal for it",
    async (state, offered, hidden) => {
      expect.hasAssertions();
      worker.use(
        noQna,
        http.get(`*/admin/${subject.resource}/7`, () => HttpResponse.json(subject.record(state))),
      );

      await subject.edit();

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

  it("the form has no state field", async () => {
    expect.hasAssertions();
    worker.use(
      noQna,
      http.get(`*/admin/${subject.resource}/7`, () => HttpResponse.json(subject.record("draft"))),
    );

    await subject.edit();

    await expect.element(page.getByLabelText(subject.nameLabel)).toBeVisible();
    await expect.element(page.getByLabelText("State")).not.toBeInTheDocument();
  }, 60_000);

  it("publishing calls the publish operation and shows the new state", async () => {
    expect.hasAssertions();
    let state: PublicationState = "draft";
    let published = false;
    worker.use(
      noQna,
      http.get(`*/admin/${subject.resource}/7`, () => HttpResponse.json(subject.record(state))),
      http.post(`*/admin/${subject.resource}/7/publish`, () => {
        published = true;
        state = "published";
        return HttpResponse.json(subject.record(state));
      }),
    );

    await subject.edit();
    await button("Publish").click();

    await expect.element(page.getByText("State changed")).toBeVisible();
    expect(published).toBe(true);
    await expect.element(button("Unpublish")).toBeVisible();
    await expect.element(button("Publish")).not.toBeInTheDocument();
  }, 60_000);

  it("archiving and unpublishing call their own operations", async () => {
    expect.hasAssertions();
    const called: string[] = [];
    let state: PublicationState = "published";
    worker.use(
      noQna,
      http.get(`*/admin/${subject.resource}/7`, () => HttpResponse.json(subject.record(state))),
      http.post(`*/admin/${subject.resource}/7/unpublish`, () => {
        called.push("unpublish");
        state = "draft";
        return HttpResponse.json(subject.record(state));
      }),
      http.post(`*/admin/${subject.resource}/7/archive`, () => {
        called.push("archive");
        state = "archived";
        return HttpResponse.json(subject.record(state));
      }),
    );

    await subject.edit();
    await button("Unpublish").click();
    await expect.element(button("Publish")).toBeVisible();
    await button("Archive").click();
    await expect.element(button("Unpublish")).not.toBeInTheDocument();

    expect(called).toStrictEqual(["unpublish", "archive"]);
  }, 60_000);

  it("a refused transition tells the admin and keeps the state", async () => {
    expect.hasAssertions();
    worker.use(
      noQna,
      http.get(`*/admin/${subject.resource}/7`, () => HttpResponse.json(subject.record("draft"))),
      http.post(`*/admin/${subject.resource}/7/publish`, () =>
        HttpResponse.json({ status: 409, title: "Conflict", type: "about:blank" }, { status: 409 }),
      ),
    );

    await subject.edit();
    await button("Publish").click();

    await expect.element(page.getByText("Could not change the state")).toBeVisible();
    await expect.element(button("Publish")).toBeVisible();
  }, 60_000);

  it("the list shows each record's state", async () => {
    expect.hasAssertions();
    worker.use(
      http.get(`*/admin/${subject.resource}`, () =>
        HttpResponse.json({
          items: [subject.record("archived")],
          page: 1,
          perPage: 25,
          total: 1,
        }),
      ),
    );

    await subject.list();

    await expect.element(page.getByRole("cell", { name: "Archived" })).toBeVisible();
  }, 60_000);
});
