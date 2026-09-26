import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import { getCurrentUserQueryKey } from "@/client/@tanstack/react-query.gen";
import type { ProfileInput, User } from "@/client/types.gen";
import { Route as editRoute } from "@/routes/{-$locale}/account/profile/edit";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

const user: User = {
  id: 1,
  firstName: "Dora",
  lastName: null,
  name: "Dora",
  email: "dora@example.com",
  admin: false,
  canAccessAdmin: false,
  assistantMessagesCount: 0,
  createdAt: "2026-09-01T00:00:00Z",
  createdAtAsTimestamp: null,
  type: "user",
};

function openProfile(visitor: User | null = user) {
  worker.use(http.get("*/api/account/profile/edit", () => HttpResponse.json(user)));
  return renderRoute(editRoute, {
    path: "/{-$locale}/account/profile/edit",
    initialPath: "/account/profile/edit",
    user: visitor,
  });
}

test("a guest is sent to sign in first", async () => {
  const { router } = await openProfile(null);

  await expect.poll(() => router.state.location.pathname).toBe("/session/new");
});

test("the profile saves the edited names", async () => {
  const received: ProfileInput[] = [];
  worker.use(
    http.patch("*/api/account/profile", async ({ request }) => {
      const body = (await request.json()) as ProfileInput;
      received.push(body);
      return HttpResponse.json({ ...user, ...body });
    }),
  );

  await openProfile();

  await expect.element(page.getByLabelText("First name")).toHaveValue("Dora");
  await page.getByLabelText("Last name").fill("Explorer");
  await page.getByRole("button", { name: "Save" }).click();

  await expect.element(page.getByText("Данные успешно обновлены")).toBeVisible();
  expect(received).toEqual([{ firstName: "Dora", lastName: "Explorer" }]);
});

test("a name the API would refuse is refused before it is sent", async () => {
  let submitted = false;
  worker.use(
    http.patch("*/api/account/profile", () => {
      submitted = true;
      return HttpResponse.json(user);
    }),
  );

  await openProfile();

  await page.getByLabelText("First name").fill("a".repeat(41));
  await page.getByLabelText("Last name").fill("Dora@home");
  await page.getByRole("button", { name: "Save" }).click();

  await expect.element(page.getByText(/40/)).toBeVisible();
  expect(submitted).toBe(false);
});

test("deleting the account, once confirmed, signs out and goes home", async () => {
  let deleted = false;
  worker.use(
    http.delete("*/api/account/profile", () => {
      deleted = true;
      return new HttpResponse(null, { status: 204 });
    }),
  );

  const { router, queryClient } = await openProfile();

  await page.getByRole("button", { name: "Delete account" }).click();
  await page.getByRole("button", { name: "Yes" }).click();

  await expect.poll(() => router.state.location.pathname).toBe("/");
  expect(deleted).toBe(true);
  expect(queryClient.getQueryData(getCurrentUserQueryKey())).toEqual({ user: null });
  await expect.element(page.getByText("Аккаунт успешно удален")).toBeVisible();
});
