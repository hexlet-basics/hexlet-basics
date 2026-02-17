import type { Method } from "@inertiajs/core";
import { Button, Checkbox, Select, Textarea, TextInput } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import { arrayToSelectData } from "@/lib/utils";
import * as Routes from "@/routes.js";
import type { Language, ReviewCrud } from "@/types";

type Props = {
  data: ReviewCrud;
  url: string;
  courses: Language[];
  method?: Method;
};

export default function Form({ courses, data, url, method }: Props) {
  const { t } = useTranslation();

  const payload = data;
  const statesSelectData = arrayToSelectData(
    data.meta?.states ?? [],
    "key",
    "value",
  );
  const coursesSelectData = arrayToSelectData(courses, "id", "slug");

  const { onSubmit, processing, form } = useAppForm(payload, {
    url,
    method: method ?? "post",
  });

  return (
    <form onSubmit={onSubmit}>
      <Select {...form.getSelectProps("state", statesSelectData)} />
      <Checkbox {...form.getCheckboxProps("pinned")} />
      <Select {...form.getSelectProps("language_id", coursesSelectData)} />
      {/* Для XAutocomplete пока ставим TextInput (можно позже сделать кастомный autocomplete-хук) */}
      <TextInput {...form.getInputProps("user_id")} />
      <TextInput {...form.getInputProps("first_name")} autoComplete="name" />
      <TextInput {...form.getInputProps("last_name")} autoComplete="name" />
      <Textarea {...form.getInputProps("body")} rows={8} />
      <Button type="submit" loading={processing}>
        {t(($) => $.helpers.submit.save)}
      </Button>
    </form>
  );
}
