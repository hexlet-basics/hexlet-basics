import dayjs from "dayjs";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import type { Banner, BannerInput } from "@/client/types.gen";
import {
  zBannerBackground,
  zBannerLocale,
  zBannerState,
} from "@/client/zod.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";

// Mantine's DateTimePicker emits/consumes this local format; the API speaks
// RFC3339. Conversion lives here so the engine's date field stays format-agnostic.
const PICKER_FORMAT = "YYYY-MM-DD HH:mm:ss";

// The banner form model. Enums reuse the generated validators (so values are
// checked and typed as the contract literals); dates are the picker's string form
// and `url` is an empty string when blank — both normalized on submit. This is a
// distinct shape from `BannerInput`, so the form validates itself rather than the
// request body (which categories could share but banners can't).
const bannerFormSchema = z.object({
  state: zBannerState,
  background: zBannerBackground,
  locale: zBannerLocale,
  body: z.string().min(1),
  url: z.string(),
  startsAt: z.string(),
  finishesAt: z.string(),
});

export type BannerFormValues = z.infer<typeof bannerFormSchema>;

export { bannerFormSchema };

export const emptyBanner: BannerFormValues = {
  state: "draft",
  background: "cta_gradient",
  locale: "en",
  body: "",
  url: "",
  startsAt: "",
  finishesAt: "",
};

const toPicker = (iso: string | null): string =>
  iso ? dayjs(iso).format(PICKER_FORMAT) : "";

const toIso = (value: string): string | null =>
  value ? dayjs(value).toISOString() : null;

// Seed the form from an existing row (edit). RFC3339 timestamps become the
// picker's local string; null columns become empty strings so inputs stay
// controlled.
export function bannerToForm(banner: Banner): BannerFormValues {
  return {
    state: banner.state,
    background: banner.background,
    locale: banner.locale,
    body: banner.body,
    url: banner.url ?? "",
    startsAt: toPicker(banner.startsAt),
    finishesAt: toPicker(banner.finishesAt),
  };
}

// Map form values to the request body: picker strings back to RFC3339, blank
// optionals to null.
export function bannerToInput(values: BannerFormValues): BannerInput {
  return {
    state: values.state,
    background: values.background,
    locale: values.locale,
    body: values.body,
    url: values.url || null,
    startsAt: toIso(values.startsAt),
    finishesAt: toIso(values.finishesAt),
  };
}

// Field descriptors driving CrudForm. Enum option labels reuse the legacy
// `*/values` i18n maps, so admin copy matches the Rails back office.
export function useBannerFields(): CrudFieldSpec<BannerFormValues>[] {
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
          label: t(
            ($) => $.models.attributes.banner["background/values"].cta_gradient,
          ),
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
      autoFocus: true,
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
