import { HttpResponse, http } from "msw";
import { describe, expect, it } from "vitest";
import { page } from "vitest/browser";
import type { ReviewInput } from "@/client/types.gen";
import { Route } from "@/routes/{-$locale}/admin/reviews/new";
import { admin } from "@/test/fixtures";
import { worker } from "@/test/msw";
import { renderRoute } from "@/test/renderRoute";
import { zReviewInput } from "@/client/zod.gen";

// Positive create flow for the reviews resource — also exercises the engine's
// number field type (courseId): a filled number submits as a number, an empty
// one as null.
describe("create review", () => {
  it("creates a review", async () => {
    expect.hasAssertions();
    let received: ReviewInput | null = null;
    worker.use(
      http.post("*/admin/reviews", async ({ request }) => {
        received = zReviewInput.parse(await request.json());
        return HttpResponse.json({
          body: "Great course",
          course: null,
          courseId: 42,
          createdAt: "2026-01-01T00:00:00Z",
          firstName: "Eva",
          fullName: "Eva",
          id: 1,
          lastName: null,
          locale: "en",
          pinned: false,
          state: "draft",
          user: null,
          userId: 0,
        });
      }),
    );

    const { router } = await renderRoute(Route, {
      initialPath: "/admin/reviews/new",
      path: "/{-$locale}/admin/reviews/new",
      user: admin,
    });

    await page.getByLabelText("Course").fill("42");
    await page.getByLabelText("First name").fill("Eva");
    await page.getByLabelText("Review *").fill("Great course");
    await page.getByRole("button", { name: "Create" }).click();

    await expect.element(page.getByText("Record created")).toBeVisible();
    expect(received).toMatchObject({
      body: "Great course",
      courseId: 42,
      firstName: "Eva",
      userId: null,
    });
    // A new review is a draft: the state has its own operations.
    expect(received).not.toHaveProperty("state");
    await expect.poll(() => router.state.location.pathname).toBe("/admin/reviews");
  }, 60_000);
});
