import type { Method } from "@inertiajs/core";
import { Button, FileInput, Select, TextInput } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import { arrayToSelectData, enumToOptions } from "@/lib/utils";
import type { LanguageCreate, LanguageUpdate } from "@/types/serializers";

type Props = {
  data: LanguageCreate | LanguageUpdate;
  url: string;
  method?: Method;
};

// const locales = [
//   { name: 'Russian', code: 'ru' },
//   { name: 'English', code: 'en' },
// ];

export default function Form({ data, url, method }: Props) {
  const { t } = useTranslation();

  const languageProgressEnum = t(($) => $.enums.language.progress, {
    returnObjects: true,
  });
  const languageProgressEnumOptions = enumToOptions(languageProgressEnum);

  const languageLearnAsEnum = t(($) => $.enums.language.learn_as, {
    returnObjects: true,
  });
  const languageLearnAsEnumOptions = enumToOptions(languageLearnAsEnum);
  const payload = data;
  const languageProgressSelectData = arrayToSelectData(
    languageProgressEnumOptions,
    "id",
    "name",
  );
  const languageLearnAsSelectData = arrayToSelectData(
    languageLearnAsEnumOptions,
    "id",
    "name",
  );

  const { onSubmit, processing, form } = useAppForm(payload, {
    url,
    method: method ?? "post",
  });

  return (
    <form onSubmit={onSubmit}>
      <Select
        label={t(($) => $.models.attributes.language.progress)}
        {...form.getSelectProps("progress", languageProgressSelectData)}
      />
      <Select
        label={t(($) => $.models.attributes.language.learn_as)}
        {...form.getSelectProps("learn_as", languageLearnAsSelectData)}
      />
      <TextInput
        label={t(($) => $.models.attributes.language.slug)}
        {...form.getInputProps("slug")}
      />
      <TextInput
        label={t(
          ($) => $.models.attributes.language.hexlet_program_landing_page,
        )}
        {...form.getInputProps("hexlet_program_landing_page")}
      />
      <TextInput
        label={t(($) => $.models.attributes.language.openai_assistant_id)}
        {...form.getInputProps("openai_assistant_id")}
      />
      <FileInput
        label={t(($) => $.models.attributes.language.cover)}
        {...form.getFileInputProps("cover")}
        name="cover"
      />
      <Button type="submit" loading={processing}>
        {t(($) => $.helpers.submit.save)}
      </Button>
    </form>
  );
}
