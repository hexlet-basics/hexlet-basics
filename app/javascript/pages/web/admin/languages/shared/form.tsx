import { useTranslation } from "react-i18next";

import { XFile, XForm, XInput, XSelect } from "@/components/forms";
import { type HTTPVerb, Submit } from "use-inertia-form";

import { enumToOptions } from "@/lib/utils";
import type LanguageCrud from "@/types/serializers/LanguageCrud";

type Props = {
  data: LanguageCrud;
  url: string;
  method?: HTTPVerb;
};

const locales = [
  { name: "Russian", code: "ru" },
  { name: "English", code: "en" },
];

export default function Form({ data, url, method }: Props) {
  const { t: tHelpers } = useTranslation("helpers");
  const { t: tEnums } = useTranslation("enumerize");
  const languageProgressEnum = tEnums("language.progress", {
    returnObjects: true,
  });
  const languageProgressEnumOptions = enumToOptions(languageProgressEnum);

  const languageLearnAsEnum = tEnums("language.learn_as", {
    returnObjects: true,
  });
  const languageLearnAsEnumOptions = enumToOptions(languageLearnAsEnum);

  return (
    <XForm method={method} model="language" data={data} to={url}>
      <XSelect
        name="progress"
        labelField="name"
        valueField="id"
        items={languageProgressEnumOptions}
      />
      <XSelect
        name="learn_as"
        labelField="name"
        valueField="id"
        items={languageLearnAsEnumOptions}
      />
      <XInput name="slug" />
      <XInput name="hexlet_program_landing_page" />
      <XInput name="openai_assistant_id" />
      <XFile metaName="cover_thumb_url" name="cover" />

      <Submit className="btn w-100 btn-lg btn-primary mb-3">
        {tHelpers("submit.save")}
      </Submit>
    </XForm>
  );
}
