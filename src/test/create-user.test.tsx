import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { UserInput } from "@/client/types.gen";
import { Route } from "@/routes/{-$locale}/admin/users/new";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { validationProblem } from "@/test/problem";
import { renderRoute } from "@/test/renderRoute";
import { zUserInput } from "@/client/zod.gen";

// Positive create flow for the users resource — also exercises the engine's
// checkbox field type (admin flag) end to end.
describe("create user", () => {
  it("creates a user", async () => {
    expect.hasAssertions();
    let received: UserInput | null = null;
    worker.use(
      http.post("*/admin/api/users", async ({ request }) => {
        received = zUserInput.parse(await request.json());
        return HttpResponse.json({
          admin: true,
          email: "dora@example.com",
          firstName: "Dora",
          id: 1,
          lastName: null,
        });
      }),
    );

    const { router } = await renderRoute(Route, {
      initialPath: "/admin/users/new",
      path: "/{-$locale}/admin/users/new",
      user: admin,
    });

    await page.getByLabelText("Email *").fill("dora@example.com");
    await page.getByLabelText("First name").fill("Dora");
    await page.getByLabelText("Admin?").click();
    await page.getByRole("button", { name: "Create" }).click();

    await expect.element(page.getByText("Record created")).toBeVisible();
    expect(received).toMatchObject({
      admin: true,
      email: "dora@example.com",
      firstName: "Dora",
    });
    await expect.poll(() => router.state.location.pathname).toBe("/admin/users");
  }, 60_000);

  // A 422 from the API lands on the field it names, through the same CrudForm
  // every resource shares, and the form stays put for another try.
  it("shows a refused field of a new user at that field", async () => {
    expect.hasAssertions();
    worker.use(
      http.post("*/admin/api/users", () =>
        validationProblem({ email: ["has already been taken"] }),
      ),
    );

    const { router } = await renderRoute(Route, {
      initialPath: "/admin/users/new",
      path: "/{-$locale}/admin/users/new",
      user: admin,
    });

    await page.getByLabelText("Email *").fill("dora@example.com");
    await page.getByRole("button", { name: "Create" }).click();

    await expect.element(page.getByText("has already been taken")).toBeVisible();
    expect(router.state.location.pathname).toBe("/admin/users/new");
  }, 60_000);
});
