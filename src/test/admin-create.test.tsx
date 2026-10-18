import type { AnyRoute } from "@tanstack/react-router";
import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import { Route as blogPostRoute } from "@/routes/{-$locale}/admin/blog_posts/new";
import { Route as courseLandingPageRoute } from "@/routes/{-$locale}/admin/course_landing_pages/new";
import { Route as courseRoute } from "@/routes/{-$locale}/admin/courses/new";
import { Route as roleRoute } from "@/routes/{-$locale}/admin/management/roles/new";
import { Route as staffMemberRoute } from "@/routes/{-$locale}/admin/management/staff_members/new";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The admin "new record" pages through their real routes, under the real admin
// guard: an admin fills the form, the API receives what was typed, and the page
// goes back to the resource's list. The forms are one engine (CrudForm) fed by a
// field spec per resource, so each case covers a resource's spec and its route.

interface Case {
  name: string;
  route: AnyRoute;
  resource: string;
  // Form fields to fill, as [label, value].
  fields: [string, string][];
  // What the API must receive of it.
  sent: object;
}

const cases: Case[] = [
  {
    name: "a course",
    route: courseRoute,
    resource: "courses",
    fields: [["Slug *", "ruby"]],
    sent: { slug: "ruby" },
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
    sent: { name: "Hello", slug: "hello" },
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
    sent: { courseId: 5, name: "Ruby", slug: "ruby-landing" },
  },
];

describe("admin create", () => {
  it.each(cases)(
    "an admin creates $name and is taken back to the list",
    async (c) => {
      expect.hasAssertions();
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
      // A new record is a draft: the state has its own operations.
      expect(received).not.toHaveProperty("state");
      // A new course is a draft too: its readiness has its own operations.
      expect(received).not.toHaveProperty("readiness");
      await expect.poll(() => router.state.location.pathname).toBe(`/admin/${c.resource}`);
    },
    60_000,
  );

  it("a visitor who is not an admin is sent away from the admin pages", async () => {
    expect.hasAssertions();
    const { router } = await renderRoute(courseRoute, {
      path: "/{-$locale}/admin/courses/new",
      initialPath: "/admin/courses/new",
      user: { ...admin, admin: false, canAccessAdmin: false },
    });

    await expect.poll(() => router.state.location.pathname).toBe("/");
  }, 60_000);
});
