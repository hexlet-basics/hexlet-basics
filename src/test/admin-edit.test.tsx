import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type {
  CourseCategory,
  QnaItem,
  RolePermissionsInput,
  StaffRoleDetail,
} from "@/client/types.gen";
import { Route as categoryRoute } from "@/routes/{-$locale}/admin/course_categories/$id";
import { Route as roleRoute } from "@/routes/{-$locale}/admin/management/roles/$id";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// The admin edit pages that carry a nested editor under the record's own form,
// through their real routes under the real admin guard: the role's permissions
// and a category's FAQ. The API is faked at the HTTP boundary; asserted is
// what an admin sees and what the API is sent.

const role: StaffRoleDetail = {
  description: null,
  grants: ["blog_posts.read"],
  id: 5,
  name: "Editor",
};

const catalog = {
  permissions: [
    { capability: true, kind: "read", name: "audit.read" },
    { capability: false, kind: "read", name: "banners.read" },
    { capability: false, kind: "write", name: "blog_posts.create" },
    { capability: false, kind: "read", name: "blog_posts.read" },
  ],
};

function openCategory() {
  return renderRoute(categoryRoute, {
    initialPath: "/admin/course_categories/7",
    path: "/{-$locale}/admin/course_categories/$id",
    user: admin,
  });
}

describe("admin edit", () => {
  it("a role's permissions are shown by resource and saved as the admin ticked them", async () => {
    expect.hasAssertions();
    const received: RolePermissionsInput[] = [];
    worker.use(
      http.get("*/admin/management/roles/5", () => HttpResponse.json(role)),
      http.get("*/admin/management/permissions", () => HttpResponse.json(catalog)),
      http.get("*/admin/management/role_permissions/5", () => HttpResponse.json(role)),
      http.put("*/admin/management/role_permissions/5", async ({ request }) => {
        received.push((await request.json()) as RolePermissionsInput);
        return HttpResponse.json(role);
      }),
    );

    await renderRoute(roleRoute, {
      initialPath: "/admin/management/roles/5",
      path: "/{-$locale}/admin/management/roles/$id",
      user: admin,
    });

    // Permissions are listed under the resource they belong to.
    await expect.element(page.getByText("banners", { exact: true })).toBeVisible();
    await expect.element(page.getByText("blog_posts", { exact: true })).toBeVisible();

    // A capability is offered with its own group and a note that it is explicit only.
    await expect.element(page.getByText("audit", { exact: true })).toBeVisible();
    await expect.element(page.getByRole("checkbox", { name: "audit.read" })).not.toBeChecked();
    await expect.element(page.getByText("Explicit grant only")).toBeVisible();

    // The role reads blog posts and may not create them yet.
    const read = page.getByRole("checkbox", { name: "blog_posts.read" });
    await expect.element(read).toBeChecked();
    const create = page.getByRole("checkbox", { name: "blog_posts.create" });
    await expect.element(create).not.toBeChecked();

    await create.click();
    await read.click();
    const save = page.getByRole("button", { name: "Save" }).last();
    await save.click();

    const confirmation = page.getByText("Record updated");
    await expect.element(confirmation).toBeVisible();
    expect(received[0]?.grants).toStrictEqual(["blog_posts.create"]);
  }, 60_000);

  const category: CourseCategory = {
    createdAt: "2026-01-01T00:00:00Z",
    description: null,
    header: "Programming",
    id: 7,
    locale: "en",
    name: "Programming",
    slug: "programming",
  };

  const faq: QnaItem[] = [{ answer: "Yes.", id: 1, question: "Is it free?" }];

  it("an admin adds a question to a category's FAQ", async () => {
    expect.hasAssertions();
    let created: unknown = null;
    worker.use(
      http.get("*/admin/course_categories/7", () => HttpResponse.json(category)),
      http.get("*/admin/course_categories/7/qna_items", () => HttpResponse.json(faq)),
      http.post("*/admin/course_categories/7/qna_items", async ({ request }) => {
        created = await request.json();
        return HttpResponse.json({ answer: "Because.", id: 2, question: "Why?" });
      }),
    );

    await openCategory();

    const existing = page.getByLabelText("Question").first();
    await expect.element(existing).toHaveValue("Is it free?");

    const question = page.getByLabelText("Question").last();
    await question.fill("Why?");
    const answer = page.getByLabelText("Answer").last();
    await answer.fill("Because.");
    const create = page.getByRole("button", { name: "Create" }).last();
    await create.click();

    const confirmation = page.getByText("Record created");
    await expect.element(confirmation).toBeVisible();
    expect(created).toStrictEqual({ answer: "Because.", question: "Why?" });
  }, 60_000);

  it("an admin edits a question and is told it was saved", async () => {
    expect.hasAssertions();
    let updated: unknown = null;
    worker.use(
      http.get("*/admin/course_categories/7", () => HttpResponse.json(category)),
      http.get("*/admin/course_categories/7/qna_items", () => HttpResponse.json(faq)),
      http.put("*/admin/course_categories/7/qna_items/1", async ({ request }) => {
        updated = await request.json();
        return HttpResponse.json({ answer: "Always.", id: 1, question: "Is it free?" });
      }),
    );

    await openCategory();

    const answer = page.getByLabelText("Answer").first();
    await answer.fill("Always.");
    const save = page.getByRole("button", { name: "Save" }).nth(1);
    await save.click();

    const confirmation = page.getByText("Record updated");
    await expect.element(confirmation).toBeVisible();
    expect(updated).toStrictEqual({ answer: "Always.", question: "Is it free?" });
  }, 60_000);

  it("an admin deletes a question after confirming", async () => {
    expect.hasAssertions();
    let deleted = false;
    worker.use(
      http.get("*/admin/course_categories/7", () => HttpResponse.json(category)),
      http.get("*/admin/course_categories/7/qna_items", () => HttpResponse.json(faq)),
      http.delete("*/admin/course_categories/7/qna_items/1", () => {
        deleted = true;
        return new HttpResponse(null, { status: 204 });
      }),
    );

    await openCategory();

    const remove = page.getByRole("button", { name: "Delete" }).first();
    await remove.click();
    const confirm = page.getByRole("dialog").getByRole("button", { name: "Delete" });
    await confirm.click();

    const confirmation = page.getByText("Record deleted");
    await expect.element(confirmation).toBeVisible();
    expect(deleted).toBe(true);
  }, 60_000);
});
