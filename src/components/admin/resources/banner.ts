import type { Banner, BannerInput } from "@/client/types.gen";
import { zBannerInput } from "@/client/zod.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { useSchemaFields } from "@/components/admin/schemaFields";

// Banner form values ARE the generated request body (BannerInput): the datetime
// fields carry RFC3339 strings (the engine's DateTimeField hides the Mantine
// picker format), so the form validates against the generated `zBannerInput`
// directly and its field list is derived from that same schema — nothing here
// restates the shape.

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

export function useBannerFields(): CrudFieldSpec<BannerInput>[] {
  return useSchemaFields<BannerInput>(zBannerInput, {
    namespace: "banner",
    textarea: ["body"],
  });
}
