import type { Method } from "@inertiajs/core";
import { Button, Select, Textarea, TextInput } from "@mantine/core";
import { DateTimePicker } from "@mantine/dates";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import { arrayToSelectData, enumToOptions } from "@/lib/utils";
import type { BannerCreate, BannerUpdate } from "@/types";

type Props = {
  data: BannerCreate | BannerUpdate;
  url: string;
  method?: Method;
};

export default function Form({ data, url, method }: Props) {
  const { t } = useTranslation();

  const statesEnum = t(($) => $.models.attributes.banner["state/values"], {
    returnObjects: true,
  });
  const statesSelectData = arrayToSelectData(
    enumToOptions(statesEnum),
    "id",
    "name",
  );

  const backgroundsEnum = t(
    ($) => $.models.attributes.banner["background/values"],
    { returnObjects: true },
  );
  const backgroundsSelectData = arrayToSelectData(
    enumToOptions(backgroundsEnum),
    "id",
    "name",
  );

  const localesEnum = t(($) => $.models.attributes.banner["locale/values"], {
    returnObjects: true,
  });
  const localesSelectData = arrayToSelectData(
    enumToOptions(localesEnum),
    "id",
    "name",
  );

  const { onSubmit, processing, form } = useAppForm(data, {
    url,
    method: method ?? "post",
  });

  return (
    <form onSubmit={onSubmit}>
      <Select
        label={t(($) => $.models.attributes.banner.locale)}
        {...form.getSelectProps("locale", localesSelectData)}
      />
      <Textarea
        label={t(($) => $.models.attributes.banner.body)}
        description={t(($) => $.admin.banners.form.body_hint)}
        {...form.getInputProps("body")}
        rows={4}
      />
      <TextInput
        label={t(($) => $.models.attributes.banner.url)}
        {...form.getInputProps("url")}
      />
      <Select
        label={t(($) => $.models.attributes.banner.background)}
        {...form.getSelectProps("background", backgroundsSelectData)}
      />
      <Select
        label={t(($) => $.models.attributes.banner.state)}
        {...form.getSelectProps("state", statesSelectData)}
      />
      <DateTimePicker
        clearable
        label={t(($) => $.models.attributes.banner.starts_at)}
        {...form.getDateInputProps("starts_at")}
      />
      <DateTimePicker
        clearable
        label={t(($) => $.models.attributes.banner.finishes_at)}
        {...form.getDateInputProps("finishes_at")}
      />
      <Button type="submit" loading={processing} mt="md">
        {t(($) => $.helpers.submit.save)}
      </Button>
    </form>
  );
}
