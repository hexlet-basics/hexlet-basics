import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import { getCurrentUserQueryKey } from "@/client/@tanstack/react-query.gen";
import type { EmailInput, ResetPasswordInput, User } from "@/client/types.gen";
import { Route as editRoute } from "@/routes/{-$locale}/password/$token/edit";
import { Route as requestRoute } from "@/routes/{-$locale}/remind_password/new";
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

function openLink(token: string) {
  return renderRoute(editRoute, {
    path: "/{-$locale}/password/$token/edit",
    initialPath: `/ru/password/${token}/edit`,
  });
}

function openRequestForm(initialPath = "/remind_password/new") {
  return renderRoute(requestRoute, { path: "/{-$locale}/remind_password/new", initialPath });
}

function acceptLink(token: string) {
  worker.use(http.get(`*/password/${token}/edit`, () => new HttpResponse(null, { status: 204 })));
}

test("a valid link signs the visitor in with the new password and takes them home", async () => {
  const received: ResetPasswordInput[] = [];
  acceptLink("good");
  worker.use(
    http.patch("*/password/good", async ({ request }) => {
      received.push((await request.json()) as ResetPasswordInput);
      return HttpResponse.json(user);
    }),
  );

  const { router, queryClient } = await openLink("good");

  // The six-character rule is on the field before anything is submitted.
  await expect.element(page.getByText("At least 6 characters")).toBeVisible();
  await page.getByLabelText("New password *").fill("new-password");
  await page.getByRole("button", { name: "Save password" }).click();

  await expect.poll(() => router.state.location.pathname).toBe("/ru");
  expect(received).toEqual([{ password: "new-password" }]);
  expect(queryClient.getQueryData(getCurrentUserQueryKey())).toEqual({ user });
});

test("a refused link sends the visitor back to the request form", async () => {
  worker.use(
    http.get("*/password/stale/edit", () =>
      HttpResponse.json({ message: "invalid" }, { status: 404 }),
    ),
  );

  const { router } = await openLink("stale");

  await expect.poll(() => router.state.location.pathname).toBe("/ru/remind_password/new");
  expect(router.state.location.search).toEqual({ invalid: true });
});

test("a short password is refused before it is sent", async () => {
  let submitted = false;
  acceptLink("short");
  worker.use(
    http.patch("*/password/short", () => {
      submitted = true;
      return HttpResponse.json(user);
    }),
  );

  const { router } = await openLink("short");

  await page.getByLabelText("New password *").fill("12345");
  await page.getByRole("button", { name: "Save password" }).click();

  await expect
    .element(page.getByLabelText("New password *"))
    .toHaveAttribute("aria-invalid", "true");
  expect(submitted).toBe(false);
  expect(router.state.location.pathname).toBe("/ru/password/short/edit");
});

test("the request form says a refused link is no longer valid", async () => {
  await openRequestForm("/remind_password/new?invalid=true");

  await expect
    .element(page.getByText("This password reset link is invalid or has expired."))
    .toBeVisible();
});

test("the request form confirms the same way for any email", async () => {
  const received: EmailInput[] = [];
  worker.use(
    http.post("*/remind_password", async ({ request }) => {
      received.push((await request.json()) as EmailInput);
      return new HttpResponse(null, { status: 204 });
    }),
  );

  await openRequestForm();

  await page.getByLabelText("Email *").fill("nobody@example.com");
  await page.getByRole("button", { name: "Reset password" }).click();

  await expect
    .element(
      page.getByText("If an account exists for this email, we've sent a password reset link."),
    )
    .toBeVisible();
  expect(received).toEqual([{ email: "nobody@example.com" }]);
});
