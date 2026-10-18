import { HttpResponse, http } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { UserInput } from "@/client/types.gen";
import { Route } from "@/routes/{-$locale}/admin/users/new";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// Positive create flow for the users resource — also exercises the engine's
// Checkbox field type (admin flag) end to end.
test("creates a user", async () => {
  let received: UserInput | undefined;
  worker.use(
    http.post("*/admin/api/users", async ({ request }) => {
      received = (await request.json()) as UserInput;
      return HttpResponse.json({
        id: 1,
        email: "dora@example.com",
        firstName: "Dora",
        lastName: null,
        admin: true,
      });
    }),
  );

  const { router } = await renderRoute(Route, {
    path: "/{-$locale}/admin/users/new",
    initialPath: "/admin/users/new",
    user: admin,
  });

  await page.getByLabelText("Email *").fill("dora@example.com");
  await page.getByLabelText("First name").fill("Dora");
  await page.getByLabelText("Admin?").click();
  await page.getByRole("button", { name: "Create" }).click();

  await expect.element(page.getByText("Record created")).toBeVisible();
  expect(received).toMatchObject({
    email: "dora@example.com",
    firstName: "Dora",
    admin: true,
  });
  await expect.poll(() => router.state.location.pathname).toBe("/admin/users");
});

// A 422 from the API lands on the field it names, through the same CrudForm
// Every resource shares, and the form stays put for another try.
test("shows a refused field of a new user at that field", async () => {
  worker.use(
    http.post("*/admin/api/users", () =>
      HttpResponse.json({ errors: { email: ["has already been taken"] } }, { status: 422 }),
    ),
  );

  const { router } = await renderRoute(Route, {
    path: "/{-$locale}/admin/users/new",
    initialPath: "/admin/users/new",
    user: admin,
  });

  await page.getByLabelText("Email *").fill("dora@example.com");
  await page.getByRole("button", { name: "Create" }).click();

  await expect.element(page.getByText("has already been taken")).toBeVisible();
  expect(router.state.location.pathname).toBe("/admin/users/new");
});
