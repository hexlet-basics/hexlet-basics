import type { Method } from "@inertiajs/core";
import {
  Button,
  Fieldset,
  MultiSelect,
  Textarea,
  TextInput,
} from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import { arrayToSelectData } from "@/lib/utils";
import type {
  LanguageCategoryCreate,
  LanguageCategoryQnaItem,
  LanguageCategoryUpdate,
  LanguageLandingPageForLists,
} from "@/types";
import QnaItemsModal from "./qna_items_modal";

type Props = {
  data: LanguageCategoryCreate | LanguageCategoryUpdate;
  qnaItems: LanguageCategoryQnaItem[];
  landingPagesForCategories: LanguageLandingPageForLists[];
  url: string;
  method?: Method;
};

export default function Form({
  data,
  qnaItems,
  landingPagesForCategories,
  url,
  method,
}: Props) {
  const { t } = useTranslation();

  const landingPagesSelectData = arrayToSelectData(
    landingPagesForCategories,
    "id",
    "header",
  );

  const { onSubmit, processing, form } = useAppForm(data, {
    url,
    method: method ?? "post",
  });

  return (
    <form onSubmit={onSubmit}>
      <Fieldset p="lg" mb="xl">
        <legend>{t(($) => $.admin.language_categories.form.main)}</legend>
        <TextInput
          label={t(($) => $.models.attributes.base.name)}
          {...form.getInputProps("name")}
          autoFocus
        />
        <TextInput
          label={t(($) => $.models.attributes.base.header)}
          {...form.getInputProps("header")}
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
      </Fieldset>
      <Fieldset p="lg" mb="xl">
        <legend>{t(($) => $.admin.language_categories.form.items)}</legend>
        <MultiSelect
          label={t(
            ($) =>
              $.models.attributes.language_category_item
                .language_landing_page_id,
          )}
          {...form.getMultiSelectProps(
            "language_landing_page_ids",
            landingPagesSelectData,
          )}
          searchable
          clearable
        />
      </Fieldset>
      {data.id && (
        <QnaItemsModal categoryId={data.id} initialItems={qnaItems} />
      )}
      <Button type="submit" mt="xl" loading={processing}>
        {t(($) => $.helpers.submit.save)}
      </Button>
    </form>
  );
}
