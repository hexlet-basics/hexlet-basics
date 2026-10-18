import { http, HttpResponse } from "msw";
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

// The Publication State of a banner, a review and a landing page is changed by
// action buttons, one per operation, never by the form (the same pattern as the
// blog post, admin-blog-post-state.test.tsx). Through the real routes under the
// real admin guard: the buttons offered follow the current state, a click calls
// that operation and the record shown is the one it answered with.

const banner = (state: PublicationState): Banner => ({
  id: 7,
  locale: "en",
  body: "Promo",
  url: null,
  background: "cta_gradient",
  state,
  startsAt: null,
  finishesAt: null,
  createdAt: "2026-01-01T00:00:00Z",
});

// The embedded account and course are only read for their name and slug.
const review = (state: PublicationState) =>
  ({
    id: 7,
    user: { id: 1, firstName: "Eva", lastName: null, name: "Eva" },
    course: { slug: "ruby" },
    userId: 1,
    courseId: 2,
    body: "Great course",
    firstName: "Eva",
    lastName: null,
    fullName: "Eva",
    state,
    locale: "en",
    pinned: false,
    createdAt: "2026-01-01T00:00:00Z",
  }) as unknown as Review;

const landing = (state: PublicationState): CourseLandingPage => ({
  id: 7,
  courseId: 2,
  courseSlug: "ruby",
  createdAt: "2026-01-01T00:00:00Z",
  slug: "ruby-ru",
  name: "Ruby",
  main: false,
  listed: true,
  state,
  order: null,
  footer: false,
  footerName: null,
  landingPageToRedirectId: null,
  metaTitle: "Ruby",
  metaDescription: "Ruby course",
  header: "Ruby",
  description: "About Ruby",
  usedInHeader: null,
  usedInDescription: null,
  outcomesHeader: null,
  outcomesDescription: null,
  outcomesImage: null,
  duration: 10,
  enrollmentsCount: 0,
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
    name: "banner",
    resource: "banners",
    record: banner,
    edit: () =>
      renderRoute(bannerEditRoute, {
        path: "/{-$locale}/admin/banners/$id",
        initialPath: "/admin/banners/7",
        user: admin,
      }),
    list: () =>
      renderRoute(bannerListRoute, {
        path: "/{-$locale}/admin/banners/",
        initialPath: "/admin/banners",
        user: admin,
      }),
    nameLabel: "Text (Markdown) *",
  },
  {
    name: "review",
    resource: "reviews",
    record: review,
    edit: () =>
      renderRoute(reviewEditRoute, {
        path: "/{-$locale}/admin/reviews/$id",
        initialPath: "/admin/reviews/7",
        user: admin,
      }),
    list: () =>
      renderRoute(reviewListRoute, {
        path: "/{-$locale}/admin/reviews/",
        initialPath: "/admin/reviews",
        user: admin,
      }),
    nameLabel: "Review *",
  },
  {
    name: "course landing page",
    resource: "course_landing_pages",
    record: landing,
    edit: () =>
      renderRoute(landingEditRoute, {
        path: "/{-$locale}/admin/course_landing_pages/$id",
        initialPath: "/admin/course_landing_pages/7",
        user: admin,
      }),
    list: () =>
      renderRoute(landingListRoute, {
        path: "/{-$locale}/admin/course_landing_pages/",
        initialPath: "/admin/course_landing_pages",
        user: admin,
      }),
    nameLabel: "Name *",
  },
];

// The edit screens of the landing page also load its FAQ and nothing else.
const noQna = http.get("*/admin/course_landing_pages/7/qna_items", () => HttpResponse.json([]));

const button = (name: string) => page.getByRole("button", { name, exact: true });

describe.each(subjects)("the $name admin screens", (subject) => {
  it.each<[PublicationState, string[], string[]]>([
    ["draft", ["Publish", "Archive"], ["Unpublish"]],
    ["published", ["Unpublish", "Archive"], ["Publish"]],
    ["archived", ["Publish"], ["Unpublish", "Archive"]],
  ])("a %s record offers only the transitions legal for it", async (state, offered, hidden) => {
    worker.use(
      noQna,
      http.get(`*/admin/${subject.resource}/7`, () => HttpResponse.json(subject.record(state))),
    );

    await subject.edit();

    for (const name of offered) {
      await expect.element(button(name)).toBeVisible();
    }
    for (const name of hidden) {
      await expect.element(button(name)).not.toBeInTheDocument();
    }
  });

  it("the form has no state field", async () => {
    worker.use(
      noQna,
      http.get(`*/admin/${subject.resource}/7`, () => HttpResponse.json(subject.record("draft"))),
    );

    await subject.edit();

    await expect.element(page.getByLabelText(subject.nameLabel)).toBeVisible();
    await expect.element(page.getByLabelText("State")).not.toBeInTheDocument();
  });

  it("publishing calls the publish operation and shows the new state", async () => {
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
  });

  it("archiving and unpublishing call their own operations", async () => {
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
  });

  it("a refused transition tells the admin and keeps the state", async () => {
    worker.use(
      noQna,
      http.get(`*/admin/${subject.resource}/7`, () => HttpResponse.json(subject.record("draft"))),
      http.post(`*/admin/${subject.resource}/7/publish`, () =>
        HttpResponse.json({ type: "about:blank", title: "Conflict", status: 409 }, { status: 409 }),
      ),
    );

    await subject.edit();
    await button("Publish").click();

    await expect.element(page.getByText("Could not change the state")).toBeVisible();
    await expect.element(button("Publish")).toBeVisible();
  });

  it("the list shows each record's state", async () => {
    worker.use(
      http.get(`*/admin/${subject.resource}`, () =>
        HttpResponse.json({
          items: [subject.record("archived")],
          total: 1,
          page: 1,
          perPage: 25,
        }),
      ),
    );

    await subject.list();

    await expect.element(page.getByRole("cell", { name: "Archived" })).toBeVisible();
  });
});
