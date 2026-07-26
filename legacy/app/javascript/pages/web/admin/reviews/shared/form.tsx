import type { Method } from "@inertiajs/core";
import { Button, Checkbox, Select, Textarea, TextInput } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import { arrayToSelectData, enumToOptions } from "@/lib/utils";
import type { Language, ReviewCreate, ReviewUpdate } from "@/types";

type Props = {
  data: ReviewCreate | ReviewUpdate;
  url: string;
  courses: Language[];
  method?: Method;
};

export default function Form({ courses, data, url, method }: Props) {
  const { t } = useTranslation();

  const statesEnum = t(($) => $.models.attributes.review["state/values"], {
    returnObjects: true,
  });
  const statesSelectData = arrayToSelectData(
    enumToOptions(statesEnum),
    "id",
    "name",
  );
  const coursesSelectData = arrayToSelectData(courses, "id", "slug");

  const { onSubmit, processing, form } = useAppForm(data, {
    url,
    method: method ?? "post",
  });

  return (
    <form onSubmit={onSubmit}>
      <Select
        label={t(($) => $.models.attributes.review.state)}
        {...form.getSelectProps("state", statesSelectData)}
      />
      <Checkbox
        label={t(($) => $.models.attributes.review.pinned)}
        {...form.getCheckboxProps("pinned")}
      />
      <Select
        label={t(($) => $.models.attributes.review.language_id)}
        {...form.getSelectProps("language_id", coursesSelectData)}
      />
      {/* Для XAutocomplete пока ставим TextInput (можно позже сделать кастомный autocomplete-хук) */}
      <TextInput
        label={t(($) => $.models.attributes.review.user_id)}
        {...form.getInputProps("user_id")}
      />
      <TextInput
        label={t(($) => $.models.attributes.review.first_name)}
        {...form.getInputProps("first_name")}
        autoComplete="name"
      />
      <TextInput
        label={t(($) => $.models.attributes.review.last_name)}
        {...form.getInputProps("last_name")}
        autoComplete="name"
      />
      <Textarea
        label={t(($) => $.models.attributes.review.body)}
        {...form.getInputProps("body")}
        rows={8}
      />
      <Button type="submit" loading={processing}>
        {t(($) => $.helpers.submit.save)}
      </Button>
    </form>
  );
}
