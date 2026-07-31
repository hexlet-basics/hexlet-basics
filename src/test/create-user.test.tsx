import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { UserInput } from "@/client/types.gen";
import { NewUser } from "@/routes/{-$locale}/admin/users/new";
import { worker } from "@/test/msw";
import { renderWithRouter } from "@/test/renderWithRouter";

// Positive create flow for the users resource — also exercises the engine's
// checkbox field type (admin flag) end to end.
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

  const { navigate } = await renderWithRouter(<NewUser />);

  await page.getByLabelText("Email").fill("dora@example.com");
  await page.getByLabelText("First name").fill("Dora");
  await page.getByLabelText("Admin?").click();
  await page.getByRole("button", { name: "Create" }).click();

  await expect.element(page.getByText("Record created")).toBeVisible();
  expect(received).toMatchObject({
    email: "dora@example.com",
    firstName: "Dora",
    admin: true,
  });
  expect(navigate).toHaveBeenCalledWith(expect.objectContaining({ to: "/{-$locale}/admin/users" }));
});
