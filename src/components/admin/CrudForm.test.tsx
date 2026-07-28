import { expect, test, vi } from "vitest";
import { page } from "vitest/browser";
import type { BannerInput } from "@/client/types.gen";
import { zBannerInput } from "@/client/zod.gen";
import { renderWithProviders } from "@/test/renderWithProviders";
import { CrudForm } from "./CrudForm";
import { emptyBanner, useBannerFields } from "./resources/banner";

// Exercises the shared form engine through the banner resource: schema-validated
// submit, enum/date defaults carried in the payload, and required-field blocking.
function BannerFormHarness({
  onSubmit,
}: {
  onSubmit: (values: BannerInput) => void;
}) {
  const fields = useBannerFields();
  return (
    <CrudForm
      fields={fields}
      schema={zBannerInput}
      defaultValues={emptyBanner}
      onSubmit={onSubmit}
      submitLabel="Save"
      onCancel={() => {}}
    />
  );
}

test("submits the typed values including enum/date defaults", async () => {
  const onSubmit = vi.fn();
  renderWithProviders(<BannerFormHarness onSubmit={onSubmit} />);

  await page.getByLabelText("Text (Markdown)").fill("Promo");
  await page.getByRole("button", { name: "Save" }).click();

  expect(onSubmit).toHaveBeenCalledTimes(1);
  expect(onSubmit).toHaveBeenCalledWith(
    expect.objectContaining({
      body: "Promo",
      state: "draft",
      background: "cta_gradient",
      locale: "en",
      url: "",
      startsAt: null,
      finishesAt: null,
    }),
  );
});

test("blocks submit when a required field is empty", async () => {
  const onSubmit = vi.fn();
  renderWithProviders(<BannerFormHarness onSubmit={onSubmit} />);

  // `body` is required (zBannerInput: min length 1); submitting empty must not
  // reach the handler.
  await page.getByRole("button", { name: "Save" }).click();

  expect(onSubmit).not.toHaveBeenCalled();
});
