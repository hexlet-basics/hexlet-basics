import type { Method } from "@inertiajs/core";
import {
  Box,
  Button,
  Checkbox,
  Fieldset,
  Select,
  Textarea,
  TextInput,
} from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import { arrayToSelectData } from "@/lib/utils";
import type {
  LanguageCategoryCrud,
  LanguageCategoryItemCrud,
  LanguageCategoryQnaItemCrud,
} from "@/types";

type Props = {
  data: LanguageCategoryCrud;
  url: string;
  method?: Method;
};

export default function Form({ data, url, method }: Props) {
  const { t } = useTranslation();

  const payload = data;
  const landingPagesSelectData = arrayToSelectData(
    data.meta?.landing_pages_for_categories ?? [],
    "id",
    "header",
  );

  const { onSubmit, processing, form } = useAppForm(payload, {
    url,
    method: method ?? "post",
  });

  const itemsCollection =
    form.useCollection<LanguageCategoryItemCrud>("items_attributes");
  const qnaItemsCollection = form.useCollection<LanguageCategoryQnaItemCrud>(
    "qna_items_attributes",
  );

  const defaultItemValues: LanguageCategoryItemCrud = {
    id: null,
    language_category_id: null,
    language_landing_page_id: null,
    _destroy: false,
  };
  const defaultQnaValues: LanguageCategoryQnaItemCrud = {
    id: null,
    question: "",
    answer: "",
    _destroy: false,
  };

  return (
    <form onSubmit={onSubmit}>
      <Fieldset p="lg" mb="xl">
        <legend>{t(($) => $.admin.language_categories.form.main)}</legend>
        <TextInput {...form.getInputProps("name")} autoFocus />
        <TextInput {...form.getInputProps("header")} />
        <TextInput {...form.getInputProps("slug")} />
        <Textarea {...form.getInputProps("description")} rows={5} />
      </Fieldset>
      <Fieldset p="lg" mb="xl">
        <legend>{t(($) => $.admin.language_categories.form.items)}</legend>
        {itemsCollection.forms.map((itemForm) => (
          <Box key={`${itemForm.index}-${itemForm.data.id ?? "new"}`} mb="xl">
            <input type="hidden" {...itemForm.getInputProps("id")} />
            <Select
              {...itemForm.getSelectProps(
                "language_landing_page_id",
                landingPagesSelectData,
              )}
            />
            <Checkbox {...itemForm.getCheckboxProps("_destroy")} />
            {!itemForm.data.id && (
              <Button
                type="button"
                variant="outline"
                color="red"
                mt="xs"
                onClick={() => itemsCollection.remove(itemForm.index)}
              >
                {t(($) => $.helpers.crud.remove)}
              </Button>
            )}
          </Box>
        ))}
        <Button
          type="button"
          variant="light"
          mt="sm"
          onClick={() => itemsCollection.add(defaultItemValues)}
        >
          {t(($) => $.helpers.crud.add)}
        </Button>
      </Fieldset>
      <Fieldset p="lg" mb="xl">
        <legend>{t(($) => $.admin.language_categories.form.qna_items)}</legend>
        {qnaItemsCollection.forms.map((qnaItemForm) => (
          <Box
            key={`${qnaItemForm.index}-${qnaItemForm.data.id ?? "new"}`}
            mb="xl"
          >
            <input type="hidden" {...qnaItemForm.getInputProps("id")} />
            <TextInput {...qnaItemForm.getInputProps("question")} />
            <Textarea {...qnaItemForm.getInputProps("answer")} rows={5} />
            <Checkbox {...qnaItemForm.getCheckboxProps("_destroy")} />
            {!qnaItemForm.data.id && (
              <Button
                type="button"
                variant="outline"
                color="red"
                mt="xs"
                onClick={() => qnaItemsCollection.remove(qnaItemForm.index)}
              >
                {t(($) => $.helpers.crud.remove)}
              </Button>
            )}
          </Box>
        ))}
        <Button
          type="button"
          variant="light"
          mt="sm"
          onClick={() => qnaItemsCollection.add(defaultQnaValues)}
        >
          {t(($) => $.helpers.crud.add)}
        </Button>
      </Fieldset>
      <Button type="submit" mt="xl" loading={processing}>
        {t(($) => $.helpers.submit.save)}
      </Button>
    </form>
  );
}
