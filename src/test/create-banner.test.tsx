import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { BannerInput } from "@/client/types.gen";
import { Route } from "@/routes/{-$locale}/admin/banners/new";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";

// Positive full-cycle use-case exercising the select + datetime field types: the
// enum defaults ride along in the payload, blank dates serialize as null, and the
// app confirms + navigates back.
describe("create banner", () => {
  it("creates a banner with the default enum selections", async () => {
    expect.hasAssertions();
    let received: BannerInput | null = null;
    worker.use(
      http.post("*/admin/banners", async ({ request }) => {
        received = (await request.json()) as BannerInput;
        return HttpResponse.json({
          background: "cta_gradient",
          body: "Promo",
          createdAt: "2026-01-01T00:00:00Z",
          finishesAt: null,
          id: 1,
          locale: "en",
          startsAt: null,
          state: "draft",
          url: null,
        });
      }),
    );

    const { router } = await renderRoute(Route, {
      initialPath: "/admin/banners/new",
      path: "/{-$locale}/admin/banners/new",
      user: admin,
    });

    await page.getByLabelText("Text (Markdown) *").fill("Promo");
    await page.getByRole("button", { name: "Create" }).click();

    await expect.element(page.getByText("Record created")).toBeVisible();
    expect(received).toMatchObject({
      background: "cta_gradient",
      body: "Promo",
      finishesAt: null,
      locale: "en",
      startsAt: null,
    });
    // A new banner is a draft: the state has its own operations.
    expect(received).not.toHaveProperty("state");
    await expect.poll(() => router.state.location.pathname).toBe("/admin/banners");
  }, 60_000);
});
