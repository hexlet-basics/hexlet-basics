import type { Banner, BannerInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { bannerFields } from "@/components/admin/generated/banner";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

// Banner form values ARE the generated request body (BannerInput): the datetime
// fields carry RFC3339 strings (the engine's DateTimeField hides the Mantine
// picker format), so the form validates against the generated `zBannerInput`
// directly — no hand-written schema, no per-resource date conversion.

export const emptyBanner: BannerInput = {
  state: "draft",
  background: "cta_gradient",
  locale: "en",
  body: "",
  url: "",
  startsAt: null,
  finishesAt: null,
};

// Seed the form from an existing row: pick only the writable fields (drops
// id/createdAt so they never leak into the request body) and keep `url` a string
// so the text input stays controlled.
export function bannerToForm(banner: Banner): BannerInput {
  return {
    state: banner.state,
    background: banner.background,
    locale: banner.locale,
    body: banner.body,
    url: banner.url ?? "",
    startsAt: banner.startsAt,
    finishesAt: banner.finishesAt,
  };
}

// Fields come from the contract (generated/banner.ts); what is listed here is
// only what the schema cannot say.
export function useBannerFields(): CrudFieldSpec<BannerInput>[] {
  return useGeneratedFields<BannerInput>(bannerFields, {
    set: { body: { type: "textarea" } },
  });
}
