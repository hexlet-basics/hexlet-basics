import type { AnyRoute } from "@tanstack/react-router";
import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { User } from "@/client/types.gen";
import { Route as blogPostRoute } from "@/routes/{-$locale}/admin/blog_posts/new";
import { Route as courseLandingPageRoute } from "@/routes/{-$locale}/admin/course_landing_pages/new";
import { Route as courseRoute } from "@/routes/{-$locale}/admin/courses/new";
import { Route as roleRoute } from "@/routes/{-$locale}/admin/management/roles/new";
import { Route as staffMemberRoute } from "@/routes/{-$locale}/admin/management/staff_members/new";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The admin "new record" pages through their real routes, under the real admin
// guard: an admin fills the form, the API receives what was typed, and the page
// goes back to the resource's list. The forms are one engine (CrudForm) fed by a
// field spec per resource, so each case covers a resource's spec and its route.

const admin: User = {
  id: 1,
  firstName: "Ada",
  lastName: null,
  name: "Ada",
  email: "ada@example.com",
  admin: true,
  canAccessAdmin: true,
  assistantMessagesCount: 0,
  createdAt: "2026-01-01T00:00:00Z",
  createdAtAsTimestamp: null,
  type: "user",
};

type Case = {
  name: string;
  route: AnyRoute;
  resource: string;
  // Form fields to fill, as [label, value].
  fields: [string, string][];
  // What the API must receive of it.
  sent: object;
};

const cases: Case[] = [
  {
    name: "a course",
    route: courseRoute,
    resource: "courses",
    fields: [["Slug *", "ruby"]],
    sent: { slug: "ruby", readiness: "draft" },
  },
  {
    name: "a role",
    route: roleRoute,
    resource: "management/roles",
    fields: [["Name *", "Editor"]],
    sent: { name: "Editor" },
  },
  {
    name: "a staff member",
    route: staffMemberRoute,
    resource: "management/staff_members",
    fields: [
      ["User *", "7"],
      ["Role *", "3"],
    ],
    sent: { userId: 7, roleId: 3 },
  },
  {
    name: "a blog post",
    route: blogPostRoute,
    resource: "blog_posts",
    fields: [
      ["Name *", "Hello"],
      ["Slug *", "hello"],
    ],
    sent: { name: "Hello", slug: "hello", state: "draft" },
  },
  {
    name: "a course landing page",
    route: courseLandingPageRoute,
    resource: "course_landing_pages",
    fields: [
      ["Course *", "5"],
      ["Name *", "Ruby"],
      ["Slug *", "ruby-landing"],
    ],
    sent: { courseId: 5, name: "Ruby", slug: "ruby-landing", state: "draft" },
  },
];

test.each(cases)("an admin creates $name and is taken back to the list", async (c) => {
  let received: unknown;
  worker.use(
    http.post(`*/admin/${c.resource}`, async ({ request }) => {
      received = await request.json();
      return HttpResponse.json({ id: 1 });
    }),
  );

  const { router } = await renderRoute(c.route, {
    path: `/{-$locale}/admin/${c.resource}/new`,
    initialPath: `/admin/${c.resource}/new`,
    user: admin,
  });

  for (const [label, value] of c.fields) {
    const field = page.getByLabelText(label);
    await field.fill(value);
  }
  const create = page.getByRole("button", { name: "Create" });
  await create.click();

  const confirmation = page.getByText("Record created");
  await expect.element(confirmation).toBeVisible();
  expect(received).toMatchObject(c.sent);
  await expect.poll(() => router.state.location.pathname).toBe(`/admin/${c.resource}`);
});

test("a visitor who is not an admin is sent away from the admin pages", async () => {
  const { router } = await renderRoute(courseRoute, {
    path: "/{-$locale}/admin/courses/new",
    initialPath: "/admin/courses/new",
    user: { ...admin, admin: false, canAccessAdmin: false },
  });

  await expect.poll(() => router.state.location.pathname).toBe("/");
});
