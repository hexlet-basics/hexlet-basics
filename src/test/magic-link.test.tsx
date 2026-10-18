import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import { getCurrentUserQueryKey } from "@/client/@tanstack/react-query.gen";
import type { EmailInput, User } from "@/client/types.gen";
import { Route as consumeRoute } from "@/routes/{-$locale}/magic_links/$token";
import { Route as requestRoute } from "@/routes/{-$locale}/magic_links/new";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";
import type { RenderedRoute } from "@/test/renderRoute";
import { zEmailInput } from "@/client/zod.gen";
import { searchOf } from "@/test/search";

const user: User = {
  admin: false,
  assistantMessagesCount: 0,
  canAccessAdmin: false,
  createdAt: "2026-09-01T00:00:00Z",
  createdAtAsTimestamp: null,
  email: "dora@example.com",
  firstName: "Dora",
  id: 1,
  lastName: null,
  name: "Dora",
  type: "user",
};

function openLink(token: string): Promise<RenderedRoute> {
  return renderRoute(consumeRoute, {
    initialPath: `/ru/magic_links/${token}`,
    path: "/{-$locale}/magic_links/$token",
  });
}

function openRequestForm(initialPath = "/magic_links/new"): Promise<RenderedRoute> {
  return renderRoute(requestRoute, { initialPath, path: "/{-$locale}/magic_links/new" });
}

describe("magic link", () => {
  it("a followed link signs the visitor in and takes them home", async () => {
    expect.hasAssertions();
    worker.use(http.get("*/magic_links/good", () => HttpResponse.json(user)));

    const { router, queryClient } = await openLink("good");

    await expect.poll(() => router.state.location.pathname).toBe("/ru");
    expect(queryClient.getQueryData(getCurrentUserQueryKey())).toStrictEqual({ user });
  }, 60_000);

  it("a refused link sends the visitor back to the request form", async () => {
    expect.hasAssertions();
    worker.use(
      http.get("*/magic_links/stale", () =>
        HttpResponse.json({ message: "invalid" }, { status: 404 }),
      ),
    );

    const { router, queryClient } = await openLink("stale");

    await expect.poll(() => router.state.location.pathname).toBe("/ru/magic_links/new");
    expect(searchOf(router)).toStrictEqual({
      invalid: true,
    });
    expect(queryClient.getQueryData(getCurrentUserQueryKey())).toBeUndefined();
  }, 60_000);

  it("the request form says a refused link is no longer valid", async () => {
    expect.hasAssertions();
    await openRequestForm("/magic_links/new?invalid=true");

    await expect
      .element(page.getByText("This sign-in link is invalid or has expired."))
      .toBeVisible();
  }, 60_000);

  it("the request form confirms the same way for any email", async () => {
    expect.hasAssertions();
    const received: EmailInput[] = [];
    worker.use(
      http.post("*/magic_links", async ({ request }) => {
        received.push(zEmailInput.parse(await request.json()));
        return new HttpResponse(null, { status: 204 });
      }),
    );

    await openRequestForm();

    await page.getByLabelText("Email *").fill("nobody@example.com");
    await page.getByRole("button", { name: "Send the link" }).click();

    await expect
      .element(page.getByText("If an account exists for this email, we've sent a sign-in link."))
      .toBeVisible();
    expect(received).toStrictEqual([{ email: "nobody@example.com" }]);
  }, 60_000);
});
