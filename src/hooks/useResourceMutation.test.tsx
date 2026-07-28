import { HttpResponse, http } from "msw";
import { expect, test } from "vitest";
import { page } from "vitest/browser";
import { adminCreateCourseCategoryMutation } from "@/client/@tanstack/react-query.gen";
import type { CourseCategoryInput } from "@/client/types.gen";
import { worker } from "@/test/msw";
import { renderWithProviders } from "@/test/renderWithProviders";
import { useResourceMutation } from "./useResourceMutation";

function CreateHarness() {
  const mutation = useResourceMutation({
    mutation: adminCreateCourseCategoryMutation(),
    invalidate: [],
    successMessage: "Saved!",
    errorMessage: "Failed!",
  });
  return (
    <button
      type="button"
      onClick={() =>
        mutation.mutate({
          body: {
            name: "Rails",
            header: "Rails",
            slug: "rails",
            description: null,
          },
        })
      }
    >
      Create
    </button>
  );
}

test("posts the body through the generated client and toasts success", async () => {
  let received: CourseCategoryInput | undefined;
  worker.use(
    http.post("*/admin/language_categories", async ({ request }) => {
      received = (await request.json()) as CourseCategoryInput;
      return HttpResponse.json({
        id: 1,
        name: "Rails",
        header: "Rails",
        slug: "rails",
        description: null,
        locale: "en",
        createdAt: "2026-01-01T00:00:00Z",
      });
    }),
  );

  renderWithProviders(<CreateHarness />);
  await page.getByRole("button", { name: "Create" }).click();

  await expect.element(page.getByText("Saved!")).toBeVisible();
  expect(received).toMatchObject({ name: "Rails", slug: "rails" });
});
