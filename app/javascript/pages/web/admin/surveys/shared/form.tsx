import type { Method } from "@inertiajs/core";
import {
  Box,
  Button,
  Fieldset,
  Select,
  Textarea,
  TextInput,
} from "@mantine/core";
import { useTranslation } from "react-i18next";
import { enums } from "@/generated/enums";
import { useAppForm } from "@/hooks/useAppForm";
import { enumToSelectData } from "@/lib/utils";
import type { SurveyCrud, SurveyItemCrud } from "@/types";

type Props = {
  data: SurveyCrud;
  url: string;
  method?: Method;
};

export default function Form({ data, url, method }: Props) {
  const { t } = useTranslation();

  const payload = data;

  const { onSubmit, processing, form } = useAppForm(payload, {
    url,
    method: method ?? "post",
  });

  const itemsCollection =
    form.useCollection<SurveyItemCrud>("items_attributes");
  const defaultItem: SurveyItemCrud = {
    id: null,
    survey_id: null,
    value: "",
    tag_list: "",
    state: null,
    order: null,
    value_for_select: null,
    _destroy: false,
    meta: { model: "", relations: {} },
  };

  return (
    <form onSubmit={onSubmit}>
      <TextInput {...form.getInputProps("question")} />
      <TextInput {...form.getInputProps("slug")} />
      <Textarea {...form.getInputProps("description")} rows={8} />
      <Fieldset>
        {itemsCollection.forms.map((itemForm) => (
          <Box key={`${itemForm.index}-${itemForm.data.id ?? "new"}`}>
            <TextInput {...itemForm.getInputProps("value")} />
            <TextInput {...itemForm.getInputProps("tag_list")} />
            <Select
              {...itemForm.getSelectProps(
                "state",
                enumToSelectData(enums.surveyItemState),
              )}
            />
            <TextInput {...itemForm.getInputProps("order")} />
            <Button
              type="button"
              variant="outline"
              color="red"
              mt="xs"
              onClick={() => itemsCollection.remove(itemForm.index)}
            >
              {t(($) => $.helpers.crud.remove)}
            </Button>
          </Box>
        ))}
        <Button
          type="button"
          variant="light"
          mt="sm"
          onClick={() => itemsCollection.add(defaultItem)}
        >
          {t(($) => $.helpers.crud.add)}
        </Button>
      </Fieldset>
      <Button type="submit" loading={processing}>
        {t(($) => $.helpers.submit.save)}
      </Button>
    </form>
  );
}
