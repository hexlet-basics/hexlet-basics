import { http, HttpResponse } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import type { ReviewInput } from "@/client/types.gen";
import { NewReview } from "@/routes/{-$locale}/admin/reviews/new";
import { worker } from "@/test/msw";
import { renderWithRouter } from "@/test/renderWithRouter";

// Positive create flow for the reviews resource — also exercises the engine's
// number field type (courseId): a filled number submits as a number, an empty
// one as null.
test("creates a review", async () => {
  let received: ReviewInput | undefined;
  worker.use(
    http.post("*/admin/reviews", async ({ request }) => {
      received = (await request.json()) as ReviewInput;
      return HttpResponse.json({
        id: 1,
        user: null,
        course: null,
        userId: 0,
        courseId: 42,
        body: "Great course",
        firstName: "Eva",
        lastName: null,
        fullName: "Eva",
        state: "draft",
        locale: "en",
        pinned: false,
        createdAt: "2026-01-01T00:00:00Z",
      });
    }),
  );

  const { navigate } = await renderWithRouter(<NewReview />);

  await page.getByLabelText("Course").fill("42");
  await page.getByLabelText("First name").fill("Eva");
  await page.getByLabelText("Review").fill("Great course");
  await page.getByRole("button", { name: "Create" }).click();

  await expect.element(page.getByText("Record created")).toBeVisible();
  expect(received).toMatchObject({
    state: "draft",
    courseId: 42,
    userId: null,
    body: "Great course",
    firstName: "Eva",
  });
  expect(navigate).toHaveBeenCalledWith(
    expect.objectContaining({ to: "/{-$locale}/admin/reviews" }),
  );
});
