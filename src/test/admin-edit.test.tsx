import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
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
// through their real routes under the real admin guard: the role's permission
// matrix and a category's FAQ. The API is faked at the HTTP boundary; asserted is
// what an admin sees and what the API is sent.

const role: StaffRoleDetail = {
  id: 5,
  name: "Editor",
  description: null,
  permissions: [
    {
      id: 1,
      roleId: 5,
      resource: "banners",
      canIndex: true,
      canCreate: false,
      canUpdate: false,
      canDestroy: false,
    },
  ],
  grants: ["blog_posts.read"],
};

const catalog = {
  permissions: [
    { name: "blog_posts.create", kind: "write" },
    { name: "blog_posts.read", kind: "read" },
  ],
};

test("a role's permissions are saved as the grid the admin ticked", async () => {
  let received: RolePermissionsInput | undefined;
  worker.use(
    http.get("*/admin/management/roles/5", () => HttpResponse.json(role)),
    http.get("*/admin/management/permissions", () => HttpResponse.json(catalog)),
    http.get("*/admin/management/role_permissions/5", () => HttpResponse.json(role)),
    http.put("*/admin/management/role_permissions/5", async ({ request }) => {
      received = (await request.json()) as RolePermissionsInput;
      return HttpResponse.json(role);
    }),
  );

  await renderRoute(roleRoute, {
    path: "/{-$locale}/admin/management/roles/$id",
    initialPath: "/admin/management/roles/5",
    user: admin,
  });

  // What the role already has is ticked; nothing else is.
  const viewBanners = page.getByRole("checkbox", { name: "banners canIndex" });
  await expect.element(viewBanners).toBeChecked();
  const createBanners = page.getByRole("checkbox", { name: "banners canCreate" });
  await expect.element(createBanners).not.toBeChecked();

  await createBanners.click();
  const save = page.getByRole("button", { name: "Save" }).last();
  await save.click();

  const confirmation = page.getByText("Record updated");
  await expect.element(confirmation).toBeVisible();
  const banners = received?.permissions.find((row) => row.resource === "banners");
  expect(banners).toMatchObject({ canIndex: true, canCreate: true, canUpdate: false });
});

test("a role's blog post permissions are shown and saved as the admin ticked them", async () => {
  let received: RolePermissionsInput | undefined;
  worker.use(
    http.get("*/admin/management/roles/5", () => HttpResponse.json(role)),
    http.get("*/admin/management/permissions", () => HttpResponse.json(catalog)),
    http.get("*/admin/management/role_permissions/5", () => HttpResponse.json(role)),
    http.put("*/admin/management/role_permissions/5", async ({ request }) => {
      received = (await request.json()) as RolePermissionsInput;
      return HttpResponse.json(role);
    }),
  );

  await renderRoute(roleRoute, {
    path: "/{-$locale}/admin/management/roles/$id",
    initialPath: "/admin/management/roles/5",
    user: admin,
  });

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
  expect(received?.grants).toEqual(["blog_posts.create"]);
});

const category: CourseCategory = {
  id: 7,
  slug: "programming",
  name: "Programming",
  header: "Programming",
  description: null,
  locale: "en",
  createdAt: "2026-01-01T00:00:00Z",
};

const faq: QnaItem[] = [{ id: 1, question: "Is it free?", answer: "Yes." }];

function openCategory() {
  return renderRoute(categoryRoute, {
    path: "/{-$locale}/admin/course_categories/$id",
    initialPath: "/admin/course_categories/7",
    user: admin,
  });
}

test("an admin adds a question to a category's FAQ", async () => {
  let created: unknown;
  worker.use(
    http.get("*/admin/course_categories/7", () => HttpResponse.json(category)),
    http.get("*/admin/course_categories/7/qna_items", () => HttpResponse.json(faq)),
    http.post("*/admin/course_categories/7/qna_items", async ({ request }) => {
      created = await request.json();
      return HttpResponse.json({ id: 2, question: "Why?", answer: "Because." });
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
  expect(created).toEqual({ question: "Why?", answer: "Because." });
});

test("an admin edits a question and is told it was saved", async () => {
  let updated: unknown;
  worker.use(
    http.get("*/admin/course_categories/7", () => HttpResponse.json(category)),
    http.get("*/admin/course_categories/7/qna_items", () => HttpResponse.json(faq)),
    http.put("*/admin/course_categories/7/qna_items/1", async ({ request }) => {
      updated = await request.json();
      return HttpResponse.json({ id: 1, question: "Is it free?", answer: "Always." });
    }),
  );

  await openCategory();

  const answer = page.getByLabelText("Answer").first();
  await answer.fill("Always.");
  const save = page.getByRole("button", { name: "Save" }).nth(1);
  await save.click();

  const confirmation = page.getByText("Record updated");
  await expect.element(confirmation).toBeVisible();
  expect(updated).toEqual({ question: "Is it free?", answer: "Always." });
});

test("an admin deletes a question after confirming", async () => {
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
});
