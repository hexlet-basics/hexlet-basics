import type { Method } from "@inertiajs/core";
import { Button, FileInput, Select, Textarea, TextInput } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import { arrayToSelectData, enumToOptions } from "@/lib/utils";
import type { BlogPostCreate, BlogPostUpdate } from "@/types";

type Props = {
  data: BlogPostCreate | BlogPostUpdate;
  url: string;
  method?: Method;
};

export default function Form({ data, url, method }: Props) {
  const { t } = useTranslation();

  const statesEnum = t(($) => $.models.attributes.blog_post["state/values"], {
    returnObjects: true,
  });
  const statesSelectData = arrayToSelectData(
    enumToOptions(statesEnum),
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
        label={t(($) => $.models.attributes.blog_post.state)}
        {...form.getSelectProps("state", statesSelectData)}
      />
      <TextInput
        label={t(($) => $.models.attributes.base.name)}
        {...form.getInputProps("name")}
      />
      <FileInput
        label={t(($) => $.models.attributes.blog_post.cover)}
        {...form.getFileInputProps("cover")}
        name="cover"
      />
      <TextInput
        label={t(($) => $.models.attributes.base.slug)}
        {...form.getInputProps("slug")}
      />
      <Textarea
        label={t(($) => $.models.attributes.base.description)}
        {...form.getInputProps("description")}
        rows={5}
      />
      <Textarea
        label={t(($) => $.models.attributes.blog_post.body)}
        {...form.getInputProps("body")}
        rows={12}
      />
      <Button type="submit" mt="xl" loading={processing}>
        {t(($) => $.helpers.submit.save)}
      </Button>
    </form>
  );
}
