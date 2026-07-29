import { useTranslation } from "react-i18next";
import type { Banner, BannerInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";

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

// Field descriptors driving CrudForm. Enum option labels reuse the legacy
// `*/values` i18n maps, so admin copy matches the Rails back office.
export function useBannerFields(): CrudFieldSpec<BannerInput>[] {
  const { t } = useTranslation();
  return [
    {
      name: "state",
      label: t(($) => $.models.attributes.banner.state),
      type: "select",
      required: true,
      options: [
        {
          value: "draft",
          label: t(($) => $.models.attributes.banner["state/values"].draft),
        },
        {
          value: "published",
          label: t(($) => $.models.attributes.banner["state/values"].published),
        },
        {
          value: "archived",
          label: t(($) => $.models.attributes.banner["state/values"].archived),
        },
      ],
    },
    {
      name: "background",
      label: t(($) => $.models.attributes.banner.background),
      type: "select",
      required: true,
      options: [
        {
          value: "cta_gradient",
          label: t(($) => $.models.attributes.banner["background/values"].cta_gradient),
        },
        {
          value: "dark",
          label: t(($) => $.models.attributes.banner["background/values"].dark),
        },
        {
          value: "blue",
          label: t(($) => $.models.attributes.banner["background/values"].blue),
        },
      ],
    },
    {
      name: "locale",
      label: t(($) => $.models.attributes.banner.locale),
      type: "select",
      required: true,
      options: [
        {
          value: "en",
          label: t(($) => $.models.attributes.banner["locale/values"].en),
        },
        {
          value: "ru",
          label: t(($) => $.models.attributes.banner["locale/values"].ru),
        },
      ],
    },
    {
      name: "body",
      label: t(($) => $.models.attributes.banner.body),
      type: "textarea",
      required: true,
    },
    {
      name: "url",
      label: t(($) => $.models.attributes.banner.url),
    },
    {
      name: "startsAt",
      label: t(($) => $.models.attributes.banner.starts_at),
      type: "datetime",
    },
    {
      name: "finishesAt",
      label: t(($) => $.models.attributes.banner.finishes_at),
      type: "datetime",
    },
  ];
}
