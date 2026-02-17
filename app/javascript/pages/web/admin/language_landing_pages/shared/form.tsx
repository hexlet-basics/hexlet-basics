import type { Method } from "@inertiajs/core";
import {
  Box,
  Button,
  Checkbox,
  Fieldset,
  FileInput,
  Select,
  Textarea,
  TextInput,
} from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import { arrayToSelectData } from "@/lib/utils";
import type {
  Language,
  LanguageLandingPage,
  LanguageLandingPageCrud,
  LanguageLandingPageQnaItemCrud,
} from "@/types";

type Props = {
  data: LanguageLandingPageCrud;
  url: string;
  method?: Method;
  languages: Language[];
  landingPages: LanguageLandingPage[];
};

// const locales = [
//   { name: "Russian", code: "ru" },
//   { name: "English", code: "en" },
// ];

export default function Form({
  data,
  landingPages,
  url,
  method,
  languages,
}: Props) {
  const { t } = useTranslation();

  const payload = data;
  const stateEventsSelectData = arrayToSelectData(
    data.meta?.state_events ?? [],
    "key",
    "value",
  );
  const languagesSelectData = arrayToSelectData(languages, "id", "slug");
  const landingPagesSelectData = arrayToSelectData(
    landingPages,
    "id",
    "header",
  );

  const { onSubmit, processing, form } = useAppForm(payload, {
    url,
    method: method ?? "post",
  });

  const qnaCollection = form.useCollection<LanguageLandingPageQnaItemCrud>(
    "qna_items_attributes",
  );
  const defaultQna: LanguageLandingPageQnaItemCrud = {
    id: null,
    question: "",
    answer: "",
    _destroy: false,
    meta: { model: "", relations: {} },
  };

  return (
    <form onSubmit={onSubmit}>
      <Checkbox {...form.getCheckboxProps("main")} />
      <Checkbox {...form.getCheckboxProps("listed")} />
      <Checkbox {...form.getCheckboxProps("footer")} />
      <Select {...form.getSelectProps("state", stateEventsSelectData)} />
      <Select {...form.getSelectProps("language_id", languagesSelectData)} />
      <Select
        {...form.getSelectProps(
          "landing_page_to_redirect_id",
          landingPagesSelectData,
        )}
      />
      <TextInput {...form.getInputProps("slug")} />
      <TextInput {...form.getInputProps("order")} />
      <TextInput {...form.getInputProps("meta_title")} />
      <Textarea {...form.getInputProps("meta_description")} rows={3} />
      <TextInput {...form.getInputProps("name")} />
      <TextInput {...form.getInputProps("header")} />
      <Textarea {...form.getInputProps("description")} rows={5} />
      <TextInput {...form.getInputProps("used_in_header")} />
      <Textarea {...form.getInputProps("used_in_description")} rows={5} />
      <FileInput {...form.getFileInputProps("outcomes_image")} />
      <TextInput {...form.getInputProps("outcomes_header")} />
      <Textarea {...form.getInputProps("outcomes_description")} rows={5} />
      <Fieldset>
        {qnaCollection.forms.map((qnaForm) => (
          <Box key={`${qnaForm.index}-${qnaForm.data.id ?? "new"}`}>
            <input type="hidden" {...qnaForm.getInputProps("id")} />
            <TextInput {...qnaForm.getInputProps("question")} />
            <Textarea {...qnaForm.getInputProps("answer")} rows={5} />
            <Checkbox {...qnaForm.getCheckboxProps("_destroy")} />
            {!qnaForm.data.id && (
              <Button
                type="button"
                variant="outline"
                color="red"
                mt="xs"
                onClick={() => qnaCollection.remove(qnaForm.index)}
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
          onClick={() => qnaCollection.add(defaultQna)}
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
