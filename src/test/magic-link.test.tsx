import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import { getCurrentUserQueryKey } from "@/client/@tanstack/react-query.gen";
import type { EmailInput, User } from "@/client/types.gen";
import { Route as consumeRoute } from "@/routes/{-$locale}/magic_links/$token";
import { Route as requestRoute } from "@/routes/{-$locale}/magic_links/new";
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
  return renderRoute(consumeRoute, {
    path: "/{-$locale}/magic_links/$token",
    initialPath: `/ru/magic_links/${token}`,
  });
}

function openRequestForm(initialPath = "/magic_links/new") {
  return renderRoute(requestRoute, { path: "/{-$locale}/magic_links/new", initialPath });
}

test("a followed link signs the visitor in and takes them home", async () => {
  worker.use(http.get("*/magic_links/good", () => HttpResponse.json(user)));

  const { router, queryClient } = await openLink("good");

  await expect.poll(() => router.state.location.pathname).toBe("/ru");
  expect(queryClient.getQueryData(getCurrentUserQueryKey())).toEqual({ user });
});

test("a refused link sends the visitor back to the request form", async () => {
  worker.use(
    http.get("*/magic_links/stale", () =>
      HttpResponse.json({ message: "invalid" }, { status: 404 }),
    ),
  );

  const { router, queryClient } = await openLink("stale");

  await expect.poll(() => router.state.location.pathname).toBe("/ru/magic_links/new");
  expect(router.state.location.search).toEqual({ invalid: true });
  expect(queryClient.getQueryData(getCurrentUserQueryKey())).toBeUndefined();
});

test("the request form says a refused link is no longer valid", async () => {
  await openRequestForm("/magic_links/new?invalid=true");

  await expect
    .element(page.getByText("This sign-in link is invalid or has expired."))
    .toBeVisible();
});

test("the request form confirms the same way for any email", async () => {
  const received: EmailInput[] = [];
  worker.use(
    http.post("*/magic_links", async ({ request }) => {
      received.push((await request.json()) as EmailInput);
      return new HttpResponse(null, { status: 204 });
    }),
  );

  await openRequestForm();

  await page.getByLabelText("Email *").fill("nobody@example.com");
  await page.getByRole("button", { name: "Send the link" }).click();

  await expect
    .element(page.getByText("If an account exists for this email, we've sent a sign-in link."))
    .toBeVisible();
  expect(received).toEqual([{ email: "nobody@example.com" }]);
});
