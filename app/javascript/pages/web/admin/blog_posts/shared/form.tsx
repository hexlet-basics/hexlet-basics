import type { Method } from "@inertiajs/core";
import { Button, FileInput, Select, Textarea, TextInput } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import { arrayToSelectData } from "@/lib/utils";
import type { BlogPostCrud } from "@/types";

type Props = {
  data: BlogPostCrud;
  url: string;
  method?: Method;
};

const locales = [
  { name: "Russian", code: "ru" },
  { name: "English", code: "en" },
];

export default function Form({ data, url, method }: Props) {
  const { t } = useTranslation();

  const payload = data;
  const statesSelectData = arrayToSelectData(
    data.meta?.states ?? [],
    "key",
    "value",
  );

  const { onSubmit, processing, form } = useAppForm(payload, {
    url,
    method: method ?? "post",
  });

  return (
    <form onSubmit={onSubmit}>
      <Select {...form.getSelectProps("state", statesSelectData)} />
      <TextInput {...form.getInputProps("name")} />
      <FileInput {...form.getFileInputProps("cover")} name="cover" />
      <TextInput {...form.getInputProps("slug")} />
      <Textarea {...form.getInputProps("description")} rows={5} />
      <Textarea {...form.getInputProps("body")} rows={12} />
      <Button type="submit" mt="xl" loading={processing}>
        {t(($) => $.helpers.submit.save)}
      </Button>
    </form>
  );
}
