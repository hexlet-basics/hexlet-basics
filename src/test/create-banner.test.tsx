import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { BannerInput } from "@/client/types.gen";
import { NewBanner } from "@/routes/{-$locale}/admin/banners/new";
import { worker } from "@/test/msw";
import { renderWithRouter } from "@/test/renderWithRouter";

// Positive full-cycle use-case exercising the select + datetime field types: the
// enum defaults ride along in the payload, blank dates serialize as null, and the
// app confirms + navigates back.
test("creates a banner with the default enum selections", async () => {
  let received: BannerInput | undefined;
  worker.use(
    http.post("*/admin/banners", async ({ request }) => {
      received = (await request.json()) as BannerInput;
      return HttpResponse.json({
        id: 1,
        locale: "en",
        body: "Promo",
        url: null,
        background: "cta_gradient",
        state: "draft",
        startsAt: null,
        finishesAt: null,
        createdAt: "2026-01-01T00:00:00Z",
      });
    }),
  );

  const { navigate } = await renderWithRouter(<NewBanner />);

  await page.getByLabelText("Text (Markdown)").fill("Promo");
  await page.getByRole("button", { name: "Create" }).click();

  await expect.element(page.getByText("Record created")).toBeVisible();
  expect(received).toMatchObject({
    body: "Promo",
    state: "draft",
    background: "cta_gradient",
    locale: "en",
    startsAt: null,
    finishesAt: null,
  });
  expect(navigate).toHaveBeenCalledWith(
    expect.objectContaining({ to: "/{-$locale}/admin/banners" }),
  );
});
