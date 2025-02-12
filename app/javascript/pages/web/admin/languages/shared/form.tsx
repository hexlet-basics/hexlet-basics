import { useTranslation } from "react-i18next";

import { XForm, XInput, XSelect } from "@/components/forms";
import type { OriginalLanguage } from "@/types/serializers";
// import { Col, Row } from "react-bootstrap";
import { type HTTPVerb, Submit } from "use-inertia-form";

// import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";
import { usePage } from "@inertiajs/react";
import { enumToOptions } from "@/lib/utils";

type Props = {
  data: OriginalLanguage;
  url: string;
  method?: HTTPVerb;
};

// const locales = [
//   { name: "Russian", code: "ru" },
//   { name: "English", code: "en" },
// ];

export default function Form({ data, url, method }: Props) {
  // const { t } = useTranslation();
  const { courseCategories } = usePage<SharedProps>().props;
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
    <XForm method={method} model="language" data={{ language: data }} to={url}>
      <XSelect
        name="category_id"
        labelField="name"
        valueField="id"
        items={courseCategories}
      />
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

      <Submit className="btn w-100 btn-lg btn-primary mb-3">
        {tHelpers("submit.save")}
      </Submit>
    </XForm>
  );
}
