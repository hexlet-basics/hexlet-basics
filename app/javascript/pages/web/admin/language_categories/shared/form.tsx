import { useTranslation } from "react-i18next";

import { XForm, XInput } from "@/components/forms";
import { type HTTPVerb, Submit } from "use-inertia-form";

import type { SharedProps } from "@/types";
import type { LanguageCategoryCrud } from "@/types/serializers";
import { usePage } from "@inertiajs/react";

type Props = {
  data: LanguageCategoryCrud;
  url: string;
  method?: HTTPVerb;
};

// const locales = [
//   { name: "Russian", code: "ru" },
//   { name: "English", code: "en" },
// ];

export default function Form({ data, url, method }: Props) {
  const { t } = useTranslation();
  const { courseCategories, landingPagesForLists } =
    usePage<SharedProps>().props;
  const { t: tHelpers } = useTranslation("helpers");
  // const { t: tEnums } = useTranslation("enumerize");
  // const languageProgressEnum = tEnums("language.progress", {
  //   returnObjects: true,
  // });
  // const languageProgressEnumOptions = enumToOptions(languageProgressEnum);
  //
  // const languageLearnAsEnum = tEnums("language.learn_as", {
  //   returnObjects: true,
  // });
  // const languageLearnAsEnumOptions = enumToOptions(languageLearnAsEnum);

  return (
    <XForm method={method} model="language_category" data={data} to={url}>
      <XInput name="name" />
      <XInput name="slug" />

      <Submit className="btn w-100 btn-lg btn-primary mb-3">
        {tHelpers("submit.save")}
      </Submit>
    </XForm>
  );
}
